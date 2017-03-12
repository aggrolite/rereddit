package main

import (
	"fmt"
	"log"

	"github.com/aggrolite/rereddit/manager"
	"github.com/jroimartin/gocui"
	"github.com/jzelinskie/geddit"
)

func main() {
	// New "root" GUI.
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	// And this "root" GUI is managed by the reddit function, which draws the views.
	g.SetManagerFunc(reddit)

	// Key bindings.
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func link(g *gocui.Gui) error {
	return nil
}

func reddit(g *gocui.Gui) error {
	// Create the header view.
	maxX, _ := g.Size()
	if v, err := g.SetView("header", -1, -1, maxX, 1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		//v.BgColor = gocui.ColorWhite
		v.FgColor = gocui.ColorRed
		fmt.Fprintln(v, "reddit")
	}

	m, err := manager.NewReRedditManager()
	if err != nil {
		return err
	}

	// Fetch current user's front page.
	links, err := m.API.Frontpage(geddit.DefaultPopularity, geddit.ListingOptions{Limit: 10})
	if err != nil {
		return err
	}

	m.Views = make([]*gocui.View, len(links))

	// Create a new view for each link found on the front page.
	y := 2
	for _, link := range links {
		if v, err := g.SetView("link"+link.ID, 0, y, maxX-(maxX/4), y+3); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Frame = true
			v.FgColor = gocui.ColorBlue
			v.Title = link.Subreddit
			v.FgColor = gocui.ColorBlue
			v.SelBgColor = gocui.ColorWhite
			fmt.Fprintf(v, "%s\n%c %d %c %c %d %c %s\n", link.Title, '\u21e7', link.Score, '\u00b7', '\U0001f4ac', link.NumComments, '\u00b7', link.Domain)
		}
		y += 4
	}

	me, err := m.API.Me()
	if err != nil {
		return err
	}

	// Create view for logged in user.
	if v, err := g.SetView("user", maxX-(maxX/4), -1, maxX, 1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintf(v, "%s %c\n", me.Name, '\u2709')
		v.Frame = false
		//v.BgColor = gocui.ColorWhite
		v.FgColor = gocui.ColorRed
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
