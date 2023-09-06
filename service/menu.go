package service

import (
	"fmt"
	"log"
	"os"

	"github.com/sshelll/clstelegraph/api"
	"github.com/sshelll/menuscreen"
)

const (
	displayTmpl = `
	Title: %s
	Tag: %v
	Content: %s
	`
)

type Menu struct{}

func (m *Menu) Start() error {
	screen := m.buildScreen()

	refreshResp, err := api.RefreshTelegraphList()
	if err != nil {
		return err
	}

	telegraphList := make([]*api.RefreshTelegraph, 0, 8)
	for k := range refreshResp.L {
		t := refreshResp.L[k]
		if t.Title == "" {
			continue
		}
		telegraphList = append(telegraphList, t)
	}

	screen.SetTitle("财联社电报")
	for _, t := range telegraphList {
		screen.AppendLines(t.Title)
	}
	screen.Start()

	idx, _, ok := screen.ChosenLine()
	screen.Fini()
	if !ok {
		os.Exit(0)
	}

	telegraph := telegraphList[idx]

	// box := tview.NewBox().SetBorder(true).SetTitle("财联社电报")

	screen = m.buildScreen()
	defer screen.Fini()
	screen.AppendLines(
		fmt.Sprintf("Title: %s", telegraph.Title),
		fmt.Sprintf("Tag: %v", telegraph.Subjects),
		fmt.Sprintf("Content: %s", telegraph.Content),
		"Press any key to continue.",
	)
	screen.Start()
	return nil
}

func (m *Menu) buildScreen() *menuscreen.MenuScreen {
	screen, err := menuscreen.NewMenuScreen()
	if err != nil {
		log.Fatalln(err)
	}
	return screen
}
