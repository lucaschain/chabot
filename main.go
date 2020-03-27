package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
)

func findTibiaPid() int32 {
	windowName := "client"
	names, err := robotgo.FindIds(windowName)

	if err != nil {
		log.Printf("error getting window names: %w", err)
		os.Exit(1)
	}

	tibiaPid := int32(0)
	for _, pid := range names {
		name, _ := robotgo.FindName(pid)

		if name == windowName {
			tibiaPid = pid
			break
		}
	}

	return tibiaPid
}

func activateTibia(pid int32) {
	log.Println("activate tibia")

	robotgo.ActivePID(pid)
}

func eatFood() {
	log.Println("eating food")
	robotgo.KeyTap("f5")
}

func makeRune() {
	log.Println("making rune")
	robotgo.KeyTap("f1")
}

func sleepRand(start, end int) {
	amount := time.Duration(start+rand.Intn(end)) * time.Millisecond

	log.Println("sleeping for", amount)

	time.Sleep(amount)
}

func foodAndRune() {
	sleepRand(15000, 21000)
	eatFood()
	sleepRand(1500, 2000)
	makeRune()
}

func main() {
	tibiaPid := findTibiaPid()
	activateTibia(tibiaPid)

	c1 := make(chan bool, 1)

	listenForEsc := func() {
		ev := robotgo.AddEvent("esc")
		if ev {
			log.Println("press q to leave")
			ev = robotgo.AddEvent("q")
			if ev {
				log.Println("exiting")
			}
		}

		c1 <- false
	}

	go listenForEsc()
	for {
		go func() {
			foodAndRune()
			c1 <- true
		}()

		cont := <-c1

		if cont != true {
			break
		}
	}
}
