package service

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/sshelll/clstelegraph/api"
	"github.com/sshelll/menuscreen"
	"github.com/sshelll/sinfra/tview/txtview"
)

type Core struct{}

func (core *Core) Start() error {
reset:
	// refresh
	telegraphList, err := core.fetch()
	if err != nil {
		return err
	}

	for {
	entry:
		// select
		telegraph, lastTime, reset_, err := core.selectTelegraph(telegraphList)
		if err != nil {
			return err
		}

		// reset
		if reset_ {
			goto reset
		}

		// fetch more or exit
		if telegraph == nil {
			if lastTime < 1 {
				return nil
			}
			telegraphList, err = core.fetchMore(lastTime)
			if err != nil {
				return err
			}
			goto entry
		}

		// display
		if err := core.displayTelegraph(telegraph); err != nil {
			return err
		}
	}
}

func (core *Core) fetch() ([]*api.Telegraph, error) {
	refreshResp, err := api.RefreshTelegraphList()
	if err != nil {
		return nil, err
	}

	telegraphList := make([]*api.Telegraph, 0, 8)
	for k := range refreshResp.L {
		t := refreshResp.L[k]
		if t.Title == "" && t.Brief == "" && t.Content == "" {
			continue
		}
		telegraphList = append(telegraphList, t)
	}

	core.sortTelegraphs(telegraphList)

	return telegraphList, nil
}

func (core *Core) fetchMore(lastTime int64) ([]*api.Telegraph, error) {
	rollResp, err := api.RollTelegraphList(lastTime, 10)
	if err != nil {
		return nil, err
	}

	telegraphList := core.filterEmptyTelegraphs(rollResp.Data.RollData)
	core.sortTelegraphs(telegraphList)

	return telegraphList, nil
}

func (core *Core) filterEmptyTelegraphs(telegraphList []*api.Telegraph) []*api.Telegraph {
	filtered := make([]*api.Telegraph, 0, len(telegraphList))
	for _, t := range telegraphList {
		if t.Title == "" && t.Brief == "" && t.Content == "" {
			continue
		}
		filtered = append(filtered, t)
	}
	return filtered
}

func (core *Core) sortTelegraphs(telegraphList []*api.Telegraph) {
	sort.Slice(telegraphList, func(i, j int) bool {
		return telegraphList[i].CTime > telegraphList[j].CTime
	})
}

func (core *Core) selectTelegraph(telegraphList []*api.Telegraph) (*api.Telegraph, int64, bool, error) {
	menu := core.buildMenu()
	defer menu.Fini()

	menu.SetTitle("财联社电报")
	for _, t := range telegraphList {
		item := fmt.Sprintf("%s ", time.Unix(t.CTime, 0).Format("2006-01-02 15:04:05"))
		if t.Title != "" {
			item += t.Title
		} else {
			item += t.Brief
		}
		menu.AppendLines(item)
	}
	menu.AppendLines("➡ More...")
	menu.AppendLines("➡ Reset")
	menu.Start()

	idx, _, ok := menu.ChosenLine()
	if !ok {
		return nil, -1, false, nil
	}

	if idx == len(telegraphList) {
		return nil, telegraphList[len(telegraphList)-1].CTime - 1, false, nil
	}

	if idx == len(telegraphList)+1 {
		return nil, -1, true, nil
	}

	return telegraphList[idx], -1, false, nil
}

func (core *Core) displayTelegraph(telegraph *api.Telegraph) error {
	txtViewer := core.buildTxtViewer(telegraph)
	go func() {
		// show title if any
		if telegraph.Title != "" {
			fmt.Fprintf(txtViewer, "[red]Title:[white]\n%s\n", telegraph.Title)
		}

		// show time if any
		if telegraph.CTime > 0 {
			fmt.Fprintf(txtViewer, "[red]Time:[white]\n%s\n", time.Unix(telegraph.CTime, 0).Format("2006-01-02 15:04:05"))
		}

		// show tags if any
		if len(telegraph.Subjects) > 0 {
			fmt.Fprintf(txtViewer, "[red]Tags:[white]\n")
			for _, s := range telegraph.Subjects {
				fmt.Fprintf(txtViewer, " %s", s.SubjectName)
				time.Sleep(50 * time.Millisecond)
			}
			fmt.Fprintf(txtViewer, "\n")
		}

		// show brief if no content
		if len(telegraph.Content) == 0 {
			if len(telegraph.Brief) > 0 {
				fmt.Fprintf(txtViewer, "[red]Brief:[white]\n%s\n", telegraph.Brief)
				return
			}
			fmt.Fprintf(txtViewer, "[grey]No content[white]")
			return
		}

		// show content
		fmt.Fprintf(txtViewer, "[red]Content:[white]\n")
		for _, c := range telegraph.Content {
			fmt.Fprintf(txtViewer, "%s", string(c))
			time.Sleep(30 * time.Millisecond)
		}
	}()
	return txtViewer.Run()
}

func (core *Core) buildMenu() *menuscreen.MenuScreen {
	screen, err := menuscreen.NewMenuScreen()
	if err != nil {
		log.Fatalln(err)
	}
	return screen
}

func (core *Core) buildTxtViewer(telegraph *api.Telegraph) *txtview.Viewer {
	opts := txtview.NewDefaultOpts()

	title := telegraph.Title
	content := telegraph.Content
	brief := telegraph.Brief

	if title != "" {
		opts.Title = &title
	}

	opts.FullScreen = false
	opts.Cols = len(title)
	if opts.Cols < 50 {
		opts.Cols = 50
	}

	if len(content) > len(brief) {
		opts.Rows = len(content)/opts.Cols + 10
	} else {
		opts.Rows = len(brief)/opts.Cols + 10
	}

	opts.DoneFunc = func(k tcell.Key, v *txtview.Viewer) {
		if k == tcell.KeyEnter {
			v.Stop()
		}
	}
	return txtview.NewViewer(opts)
}
