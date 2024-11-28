package main

import (
    "fmt"
    "bufio"
    "log"
    "os"
    "github.com/mmcdole/gofeed"
    "github.com/charmbracelet/bubbletea"
)

func GetFeeds() ([]string){
    file,err := os.Open("urls")
    if err != nil {
        log.Fatalf("unable to read file: ",err)
    }
    defer file.Close()

    var lines []string
    scanner:= bufio.NewScanner(file)
    for scanner.Scan(){
        lines=append(lines, scanner.Text())
    }
    if err := scanner.Err(); err != nil{
        log.Fatalf("Unable to get lines from file: ",err)
    }
    return lines
}

func ReadFeed(url string) *gofeed.Feed {
    fp := gofeed.NewParser()
    feed, _ := fp.ParseURL(url)
    return feed
}

type model struct {
    choices []string
    cursor int
    selected map[int] struct{}
}

func initialModel() model {
    rss_urls := GetFeeds()
    var sources ([]string)
    for _,url := range rss_urls {
        feed := ReadFeed(url)
            sources = append(sources,feed.Title)
        }
    
    return model{
        choices: sources,
        selected: make(map[int]struct{}),
    }
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) View() string {
    // The header
    s := "What should we buy at the market?\n\n"

    // Iterate over our choices
    for i, choice := range m.choices {

        // Is the cursor pointing at this choice?
        cursor := " " // no cursor
        if m.cursor == i {
            cursor = ">" // cursor!
        }

        // Is this choice selected?
        checked := " " // not selected
        if _, ok := m.selected[i]; ok {
            checked = "x" // selected!
        }

        // Render the row
        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
    }

    // The footer
    s += "\nPress q to quit.\n"

    // Send the UI for rendering
    return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    // Is it a key press?
    case tea.KeyMsg:

        // Cool, what was the actual key pressed?
        switch msg.String() {

        // These keys should exit the program.
        case "ctrl+c", "q":
            return m, tea.Quit

        // The "up" and "k" keys move the cursor up
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }

        // The "down" and "j" keys move the cursor down
        case "down", "j":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }

        // The "enter" key and the spacebar (a literal space) toggle
        // the selected state for the item that the cursor is pointing at.
        case "enter", " ":
            _, ok := m.selected[m.cursor]
            if ok {
                delete(m.selected, m.cursor)
            } else {
                m.selected[m.cursor] = struct{}{}
            }
        }
    }

    // Return the updated model to the Bubble Tea runtime for processing.
    // Note that we're not returning a command.
    return m, nil
}
    
    
func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}
