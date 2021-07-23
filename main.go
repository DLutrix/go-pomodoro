package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gen2brain/beeep"
)

func notify(title, message, appIcon string) {
	err := beeep.Alert(title, message, appIcon)
	if err != nil {
		panic(err)
	}
}

func pomodoro(focusInSecond, breakInSecond, session int, isDone, isBreak chan<- bool, isStart chan<- int) {
	for i := 0; i < session; i++ {

		isStart <- i + 1

		fmt.Print("\033[41m")

		for i := focusInSecond; i > 0; i-- {
			fmt.Printf("\033[2K\r%d:%d - %s", i/60, i%60, "Focus")
			time.Sleep(1 * time.Second)
		}

		if i == session-1 {
			break
		}

		isBreak <- true

		fmt.Print("\033[42m")
		fmt.Print("\033[30m")

		for i := breakInSecond; i > 0; i-- {
			fmt.Printf("\033[2K\r%d:%d - %s", i/60, i%60, "Break")
			time.Sleep(1 * time.Second)
		}
	}

	isDone <- true
}

func main() {
	focusTime := flag.Int(`f`, 25, `put your timer`)
	breakTime := flag.Int(`b`, 5, `put yout break time`)
	session := flag.Int(`s`, 1, `put your session`)

	flag.Parse()

	if *focusTime == 0 || *breakTime == 0 || *session == 0 {
		log.Fatalln(`Flag cannot be zero`)
	}

	focusInMinute := time.Duration(*focusTime) * time.Minute
	focusInSecond := int(focusInMinute / time.Second)
	breakInMinute := time.Duration(*breakTime) * time.Minute
	breakInSecond := int(breakInMinute / time.Second)

	isDone := make(chan bool)
	isBreak := make(chan bool)
	isStart := make(chan int)

	go pomodoro(focusInSecond, breakInSecond, *session, isDone, isBreak, isStart)

receiveChannels:
	for {
		select {
		case <-isBreak:
			notify("GO POMODORO", "Time to break", "assets/logo.png")
		case v := <-isStart:
			notify("GO POMODORO", fmt.Sprintf("Session %d start", v), "assets/logo.png")
		case <-isDone:
			fmt.Print("\033[49m")
			fmt.Print("\033[H\033[2J")
			notify("GO POMODORO", "Done", "assets/logo.png")
			break receiveChannels
		}
	}
}
