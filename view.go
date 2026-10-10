package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

const helpText = "Ctrl+n : New File   Ctrl+l : List   Esc : Back   Ctrl+s : Save   Ctrl+d : Delete   Ctrl+c : Quit"

// VIEW
func (m Model) View() tea.View {
	welcome := titleStyle.Render("WELCOME TO NOTION...🧠")

	return tea.NewView(fmt.Sprintf("\n%s\n\n%s\n\n%s\n", welcome, m.body(), helpText))
}

// body renders whichever screen is active.
func (m Model) body() string {
	body := ""

	if m.createFileInputVisible {
		body = m.textInput.View()
	}

	if m.currentFile != nil {
		body = m.textarea.View()
		if m.confirmDelete {
			body += "\n  Press Ctrl+d again to delete this file, any other key to cancel."
		}
	}

	// the list covers the editor if both are active
	if m.listVisible {
		body = m.list.View()
		if m.confirmDelete {
			body += "\n  Press Ctrl+d again to delete, any other key to cancel."
		}
	}

	return body
}