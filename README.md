# Notion

A small note-taking app that runs in your terminal. I built it in Go with [Bubble Tea](https://github.com/charmbracelet/bubbletea) because I have never used a terminal note taking app and I wondered what one would feel like in the terminal. So I built it, and learned Bubble Tea along the way. 
 
## What it does
 
- Create a new note by typing a name
- Browse all your notes in a list, with the last-modified time next to each one
- Filter the list by typing `/`
- Open a note, edit it, save it
- Delete a note (it asks you to confirm first)
  
That's about it. Give me some suggestions on how to improve it.
 
## Running it
 
You need Go installed. Any recent version should be fine, since Bubble Tea v2 wants a fairly new one.
 
```bash
git clone <repo-url>
cd Notion
go run .
```
 
Use `go run .` and not `go run main.go`. The code is split across several files, and `go run main.go` only compiles that one file, so you'll get an `undefined: initialModel` error.
 
## Keys
 
| Key | What it does |
|---|---|
| `Ctrl+n` | New note. Type a name and press Enter |
| `Ctrl+l` | Open the list of notes |
| `Enter` | In the list: open the selected note. In the name box: create the note |
| `Ctrl+s` | Save the note and close the editor |
| `Ctrl+d` | Delete the selected note (from the list) or the open one. Press it twice to confirm |
| `Esc` | Go back. This does not save, so use `Ctrl+s` first |
| `Ctrl+c` | Quit |
 
## Where your notes live
 
Notes are stored in a `.notion` folder inside your home directory:
 
- Windows: `C:\Users\<you>\.notion`
- macOS / Linux: `~/.notion`
The folder is created the first time you run the app. New notes get a `.md` extension automatically, so typing `notes` gives you `notes.md`.
 
## Structure
 
Everything is in `package main`, split by job:
 
| File | What's in it |
|---|---|
| `main.go` | Starts the program |
| `model.go` | The model (the app's state) and the setup for the input, editor and list |
| `update.go` | Key handling. One small function per shortcut |
| `view.go` | Draws whichever screen is active |
| `styles.go` | Colors and Lip Gloss styles |
| `vault.go` | Finding the notes folder, listing files, saving to disk |
 
If you haven't used Bubble Tea before, the short version is that it follows the Elm pattern. The model holds the state, `Update` reacts to key presses and returns the new state, and `View` turns that state into text. Nothing draws to the screen directly.
 
## Built with
 
- [Bubble Tea v2](https://github.com/charmbracelet/bubbletea) for the app loop
- [Bubbles v2](https://github.com/charmbracelet/bubbles) for the text input, text area and list
- [Lip Gloss v2](https://github.com/charmbracelet/lipgloss) for styling
