package main

import "charm.land/lipgloss/v2"

var (
	cursorColor = lipgloss.Color("212")

	// margin around the notes list
	docStyle = lipgloss.NewStyle().Margin(1, 2)

	// "WELCOME TO NOTION" banner
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#a0209e")).
			PaddingRight(2).
			PaddingLeft(2)
)