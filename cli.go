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

func countdown(screen tcell.Screen, dur int) {
	for i := dur; i >= 0; i-- {
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

		go countdown(screen, duration)

		width, height := screen.Size()
		x := rand.Intn(width - 4)
		y := rand.Intn(height - 4)

		target := []string{
			" ●● ",
			"●●●●",
			" ●● ",
		}

		score := 0
		scoreStr := "Score:0"
		drawText(screen, width-len([]rune(scoreStr)), 0, scoreStr, tcell.StyleDefault.Foreground(tcell.ColorYellow))
		drawTarget(screen, x, y, target)
		screen.Show()

		for {
			if timeLeft <= 0 {
				screen.Fini()
				showOverview(score, width, height, duration)
				return
			}
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventMouse:
				mx, my := ev.Position()
				if ev.Buttons() == tcell.Button1 {
					if mx >= x && mx < x+4 && my >= y && my < y+4 {
						x = rand.Intn(width - 4)
						y = rand.Intn(height - 4)
						score++
						screen.Clear()
						drawTarget(screen, x, y, target)
						scoreStr := "Score:" + strconv.Itoa(score)
						drawText(screen, width-len([]rune(scoreStr)), 0, scoreStr, tcell.StyleDefault.Foreground(tcell.ColorYellow))
						timeStr := "Time:" + strconv.Itoa(timeLeft)
						drawText(screen, width-len([]rune(timeStr)), 1, timeStr, tcell.StyleDefault.Foreground(tcell.ColorWhite))
						screen.Show()
					}
				}
			case *tcell.EventKey:
				if ev.Rune() == 'q' {
					screen.Fini()
					showOverview(score, width, height, duration)
					return
				}
			}
		}
	},
}

func init() {
	rootCmd.Flags().IntVar(&duration, "time", 60, "game duration in seconds")
}
