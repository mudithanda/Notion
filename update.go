package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// UPDATE
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-5)

	case tea.KeyPressMsg:
		var (
			keyCmd  tea.Cmd
			handled bool
		)
		m, keyCmd, handled = m.handleKey(msg)
		if handled {
			return m, keyCmd
		}
		// not handled by us: fall through and forward the key to the active widget
	}

	// File name input
	if m.createFileInputVisible {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	// Textarea
	if m.currentFile != nil {
		m.textarea, cmd = m.textarea.Update(msg)
	}

	// List
	if m.listVisible {
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

// handleKey runs our own shortcuts. The bool reports whether the key was fully handled
// if false, Update forwards the key to the active widget
func (m Model) handleKey(msg tea.KeyPressMsg) (Model, tea.Cmd, bool) {
	if m.confirmDelete && msg.String() != "ctrl+d" {
		m.confirmDelete = false // any other key cancels
	}

	switch msg.String() {

	case "ctrl+c":
		return m, tea.Quit, true

	case "ctrl+n":
		m.createFileInputVisible = true
		return m, nil, true

	case "esc":
		return m.handleEsc()

	case "ctrl+l":
		m.list.SetItems(ListFiles(m.vaultDir)) // refresh every time it opens
		m.listVisible = true
		return m, nil, true

	case "ctrl+s":
		return m.handleSave()

	case "ctrl+d":
		return m.handleDelete()

	case "enter":
		return m.handleEnter()
	}

	return m, nil, false
}

// Esc: go back from whatever screen is open.
func (m Model) handleEsc() (Model, tea.Cmd, bool) {
	if m.createFileInputVisible {
		m.textInput.Reset()
		m.createFileInputVisible = false
	}

	m.closeEditor()

	if m.listVisible {
		if m.list.FilterState() == list.Filtering {
			return m, nil, false // let the list cancel its own filter
		}
		m.listVisible = false
	}

	return m, nil, true
}

// Ctrl+S: write the textarea to the open file, then close it.
func (m Model) handleSave() (Model, tea.Cmd, bool) {
	if m.currentFile == nil {
		return m, nil, false
	}

	if err := saveFile(m.currentFile, m.textarea.Value()); err != nil {
		fmt.Println("Cannot save the file :(")
		return m, nil, true
	}

	if err := m.currentFile.Close(); err != nil {
		fmt.Println("Cannot close file :(")
	}

	m.currentFile = nil
	m.textarea.Reset()

	return m, nil, true
}

// Ctrl+D: delete the selected file from list or the open file
// asking for a second Ctrl+D to confirm.
func (m Model) handleDelete() (Model, tea.Cmd, bool) {

	// work out which file this applies to
	var path string

	switch {
	case m.listVisible:
		sel, ok := m.list.SelectedItem().(item)
		if !ok {
			return m, nil, true
		}
		path = filepath.Join(m.vaultDir, sel.title)

	case m.currentFile != nil:
		path = m.currentFile.Name()

	default:
		return m, nil, true // nothing to delete
	}

	// first press -> ask for confirmation
	if !m.confirmDelete {
		m.confirmDelete = true
		return m, nil, true
	}

	// second press -> delete
	m.confirmDelete = false

	// close the handle first if the target is the open file
	if m.currentFile != nil && m.currentFile.Name() == path {
		m.closeEditor()
	}

	if err := os.Remove(path); err != nil {
		log.Printf("Error deleting file: %v", err)
		return m, nil, true
	}

	m.list.SetItems(ListFiles(m.vaultDir))
	return m, nil, true
}

// Enter -> open the selected note (list) or create a new one (file name input).
func (m Model) handleEnter() (Model, tea.Cmd, bool) {
	if m.listVisible {
		return m.openSelected(), nil, true
	}

	return m.createNote()
}

// openSelected loads the highlighted note into the editor.
func (m Model) openSelected() Model {
	sel, ok := m.list.SelectedItem().(item)
	if !ok {
		return m
	}

	path := filepath.Join(m.vaultDir, sel.title)

	content, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return m
	}

	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return m
	}

	m.closeEditor() // don't leak a previously open file
	m.textarea.SetValue(string(content))
	m.currentFile = f
	m.listVisible = false

	return m
}

// createNote creates <name>.md from the file name input and opens it.
func (m Model) createNote() (Model, tea.Cmd, bool) {
	filename := strings.TrimSpace(m.textInput.Value())
	if filename == "" {
		return m, nil, false // let Enter reach the textarea
	}

	path := filepath.Join(m.vaultDir, filename+".md")

	if _, err := os.Stat(path); err == nil { // file exists
		return m, nil, true
	}

	f, err := os.Create(path)
	if err != nil {
		log.Printf("Could not create file: %v", err)
		return m, nil, true
	}

	m.currentFile = f
	m.textInput.Reset()
	m.createFileInputVisible = false

	return m, nil, true
}