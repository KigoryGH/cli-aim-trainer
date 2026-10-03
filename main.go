package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
)

var (
	timeLeft int
	duration int
)

func drawText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, ch := range []rune(text) {
		screen.SetContent(x+i, y, ch, nil, style)
	}
}

func countdown(screen tcell.Screen, dur int, stop chan bool) {
	for i := dur; i >= 0; i-- {
		select {
		case <-stop:
			return
		default:
			timeLeft = i
			width, _ := screen.Size()
			timeStr := "Time:" + strconv.Itoa(i)
			startX := width - len([]rune(timeStr))
			if i <= 6 {
				fmt.Print("\a")
				drawText(screen, startX, 1, timeStr, tcell.StyleDefault.Foreground(tcell.ColorRed))
			} else {
				drawText(screen, startX, 1, timeStr, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			}
			screen.Show()
			time.Sleep(time.Second)
		}
	}
}

func drawTarget(screen tcell.Screen, x, y int, target []string) {
	for dy, row := range target {
		for dx, ch := range []rune(row) {
			if ch != ' ' {
				screen.SetContent(x+dx, y+dy, ch, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
			}
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "cli-aim-trainer",
	Short: "A CLI aim trainer game",
	Run: func(cmd *cobra.Command, args []string) {
		timeLeft = duration

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
		x := rand.Intn(width - 4)
		y := rand.Intn(height - 4)

		target := []string{
			" ●● ",
			"●●●●",
			" ●● ",
		}

		score := 0
		confirming := false

		redraw := func() {
			screen.Clear()
			scoreStr := "Score:" + strconv.Itoa(score)
			drawText(screen, width-len([]rune(scoreStr)), 0, scoreStr, tcell.StyleDefault.Foreground(tcell.ColorYellow))
			timeStr := "Time:" + strconv.Itoa(timeLeft)
			drawText(screen, width-len([]rune(timeStr)), 1, timeStr, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			drawTarget(screen, x, y, target)
			screen.Show()
		}

		stop := make(chan bool)
		go countdown(screen, duration, stop)
		redraw()

		for {
			if timeLeft <= 0 {
				screen.Fini()
				showOverview(score, width, height, duration)
				return
			}
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventMouse:
				if !confirming {
					mx, my := ev.Position()
					if ev.Buttons() == tcell.Button1 {
						if mx >= x && mx < x+4 && my >= y && my < y+4 {
							x = rand.Intn(width - 4)
							y = rand.Intn(height - 4)
							score++
							redraw()
						}
					}
				}
			case *tcell.EventKey:
				if confirming {
					if ev.Rune() == 'y' {
						stop <- true
						score = 0
						timeLeft = duration
						x = rand.Intn(width - 4)
						y = rand.Intn(height - 4)
						stop = make(chan bool)
						go countdown(screen, duration, stop)
						confirming = false
						redraw()
					} else if ev.Rune() == 'n' {
						confirming = false
						redraw()
					}
				} else {
					if ev.Rune() == 'q' {
						screen.Fini()
						showOverview(score, width, height, duration)
						return
					}
					if ev.Rune() == 'r' {
						confirming = true
						resetMsg := "Reset?"
						yesMsg := "(y)es"
						noMsg := "(n)o"
						drawText(screen, width/2-len(resetMsg)/2, 0, resetMsg, tcell.StyleDefault.Foreground(tcell.ColorYellow))
						drawText(screen, width/2-len(yesMsg)/2, 1, yesMsg, tcell.StyleDefault.Foreground(tcell.ColorGreen))
						drawText(screen, width/2-len(noMsg)/2, 2, noMsg, tcell.StyleDefault.Foreground(tcell.ColorRed))
						screen.Show()
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.Flags().IntVar(&duration, "time", 60, "game duration in seconds")
}
