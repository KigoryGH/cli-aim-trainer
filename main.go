package main

import (
	"math/rand"
	"strconv"

	tcell "github.com/gdamore/tcell/v2"
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
	screen.EnableMouse()

	width, height := screen.Size()
	x := rand.Intn(width)
	y := rand.Intn(height)
	screen.SetContent(x, y, 'X', nil, tcell.StyleDefault)

	label := "Score:"
	for i, ch := range label {
		screen.SetContent(i, 0, ch, nil, tcell.StyleDefault)
	}
	scoreStr := strconv.Itoa(0)
	for i, ch := range scoreStr {
		screen.SetContent(len(label)+i, 0, ch, nil, tcell.StyleDefault)
	}
	screen.Show()

	score := 0

	for {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventMouse: // for mouse primaryButton
			mx, my := ev.Position()
			if ev.Buttons() == tcell.Button1 {
				if mx == x && my == y {
					x = rand.Intn(width)
					y = rand.Intn(height)
					screen.Clear()
					screen.SetContent(x, y, 'X', nil, tcell.StyleDefault)
					score++
					label := "Score:"
					for i, ch := range label {
						screen.SetContent(i, 0, ch, nil, tcell.StyleDefault)
					}
					scoreStr := strconv.Itoa(score)
					for i, ch := range scoreStr {
						screen.SetContent(len(label)+i, 0, ch, nil, tcell.StyleDefault)
					}
					screen.Show()
				}
			}
		case *tcell.EventKey: // for quit keypress
			if ev.Rune() == 'q' {
				screen.Fini()
				return

			}
		}
	}
}
