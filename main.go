package main

import (
	"math/rand"
	"strconv"
	"time"

	tcell "github.com/gdamore/tcell/v2"
)

var timeLeft int = 60

func drawText(screen tcell.Screen, x, y int, text string) {
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, nil, tcell.StyleDefault)
	}
}

func countdown(screen tcell.Screen) {
	_, _ = screen.Size()
	for i := 60; i >= 0; i-- {
		timeLeft = i
		width, _ := screen.Size()
		timeStr := "Time:" + strconv.Itoa(i)
		startX := width - len(timeStr)
		drawText(screen, startX, 1, timeStr)
		screen.Show()
		time.Sleep(time.Second)
	}
}

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

	go countdown(screen)

	width, height := screen.Size()
	x := rand.Intn(width)
	y := rand.Intn(height)

	scoreStr := "Score:0"
	drawText(screen, width-len(scoreStr), 0, scoreStr)
	drawText(screen, x, y, "X")
	screen.Show()

	score := 0

	for {
		if timeLeft <= 0 {
			screen.Fini()
			return
		}
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventMouse:
			mx, my := ev.Position()
			if ev.Buttons() == tcell.Button1 {
				if mx == x && my == y {
					x = rand.Intn(width)
					y = rand.Intn(height)
					score++
					screen.Clear()
					drawText(screen, x, y, "X")
					scoreStr := "Score:" + strconv.Itoa(score)
					drawText(screen, width-len(scoreStr), 0, scoreStr)
					timeStr := "Time:" + strconv.Itoa(timeLeft)
					drawText(screen, width-len(timeStr), 1, timeStr)
					screen.Show()
				}
			}
		case *tcell.EventKey:
			if ev.Rune() == 'q' {
				screen.Fini()
				return
			}
		}
	}
}
