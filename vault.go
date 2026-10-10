package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/list"
)

// vaultDir returns ~/.notion, creating it if it doesn't exist.
func vaultDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, ".notion")

	if err := os.MkdirAll(dir, 0750); err != nil {
		return "", err
	}

	return dir, nil
}

// item is one entry in the notes list (implements list.Item).
type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// ListFiles returns every file in dir as list items.
func ListFiles(dir string) []list.Item {
	items := make([]list.Item, 0)

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Println("Error reading files :(")
		return items
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		modTime := info.ModTime().Format("02 Jan 2006, 15:04")

		items = append(items, item{
			title: entry.Name(),
			desc:  fmt.Sprintf("Modified: %s", modTime),
		})
	}

	return items
}

// saveFile overwrites f with content.
func saveFile(f *os.File, content string) error {
	if err := f.Truncate(0); err != nil {
		return err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	_, err := f.WriteString(content)
	return err
}