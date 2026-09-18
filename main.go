package main

import (
	tcell "github.com/gdamore/tcell/v2"
	"math/rand"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	err = screen.Init()
	if err != nil {
		panic(err)
	}

	screen.Clear()

	width, height := screen.Size()
	x := rand.Intn(width)
	y := rand.Intn(height)
	screen.SetContent(x, y, 'X', nil, tcell.StyleDefault)

	screen.Show()

	switch ev := screen.PollEvent().(type) {
	case *tcell.EventKey:
		if ev.Rune() == 'q' {
			screen.Fini()
			return
		}
	}
}
