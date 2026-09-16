package main

import (
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

			screen.PollEvent()
		ev := screen.PollEvent()

		switch ev.(type) {
case *tcell.EventKey:
}
		screen.Fini()
}
