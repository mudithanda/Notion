package main

import (
	"log"
	"os"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// MODEL
type Model struct {
	textInput              textinput.Model // new file name input
	createFileInputVisible bool
	vaultDir               string
	currentFile            *os.File
	textarea               textarea.Model
	list                   list.Model
	listVisible            bool
	confirmDelete          bool
}

func (m Model) Init() tea.Cmd { // if want to do something before application runs
	return nil
}

// closeEditor closes the open file if any is open and clears the editor.
func (m *Model) closeEditor() {
	if m.currentFile != nil {
		m.currentFile.Close()
		m.currentFile = nil
	}
	m.textarea.Reset()
}

// INITIALIZE
func initialModel() Model {
	vaultdir, err := vaultDir()
	if err != nil {
		log.Fatal(err)
	}

	return Model{
		textInput:              newFileNameInput(),
		createFileInputVisible: false,
		vaultDir:               vaultdir,
		textarea:               newEditor(),
		list:                   newNoteList(vaultdir),
	}
}

func newFileNameInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Enter the name of file..."
	ti.SetVirtualCursor(true)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(50)

	s := ti.Styles()
	s.Cursor.Color = cursorColor
	s.Focused.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA"))    // typed text when focused
	s.Focused.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("212")) // placeholder
	s.Focused.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))      // the ">" prompt
	ti.SetStyles(s)

	return ti
}

func newEditor() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Write your note here..."
	ta.SetVirtualCursor(true)
	ta.ShowLineNumbers = false
	ta.SetStyles(textarea.DefaultStyles(true)) // default dark styles.
	ta.Focus()

	return ta
}

func newNoteList(dir string) list.Model {
	l := list.New(ListFiles(dir), list.NewDefaultDelegate(), 60, 20)
	l.Title = "Your Notes"

	return l
}