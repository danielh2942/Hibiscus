package jingle

// Keyboard Listener
import (
	"fmt"
	"time"

	"github.com/daspoet/gowinkey"
)

func KeyboardTest() {
	events, stopFn := gowinkey.Listen()

	time.AfterFunc(time.Minute, func() {
		stopFn()
	})

	for e := range events {
		switch e.State {
		case gowinkey.KeyDown:
			fmt.Println("pressed", e)
		case gowinkey.KeyUp:
			fmt.Println("released", e)
		}
	}
}

type KeyboardListener struct {
	activeKeys []string
	instrument Instrument
	octave     int
	running    bool
}

func (kbd *KeyboardListener) StartMonitor() {
	kbd.running = true
	kbd.instrument.Init()
	kbd.octave = 5
	kbd.activeKeys = make([]string, 0)
	keyVals := map[string]int{
		"A": 0,
		"W": 1,
		"S": 2,
		"E": 3,
		"D": 4,
		"F": 5,
		"T": 6,
		"G": 7,
		"Y": 8,
		"H": 9,
		"U": 10,
		"J": 11,
	}
	go func() {
		events, stopFn := gowinkey.Listen()
	eventloop:
		for e := range events {
			switch e.State {
			case gowinkey.KeyDown:
				switch e.String() {
				case "A", "W", "E", "T", "Y", "U", "S", "D", "F", "G", "H", "J":
					fmt.Println(e.String())
					for _, keys := range kbd.activeKeys {
						if keys == e.String() {
							continue eventloop
						}
					}
					kbd.activeKeys = append(kbd.activeKeys, e.String())
					kbd.instrument.AddNote((12 * kbd.octave) + keyVals[e.String()])

				case "Z":
					if kbd.octave > 0 {
						kbd.instrument.FlushNotes()
						kbd.activeKeys = make([]string, 0)
						kbd.octave--
					}

				case "X":
					if kbd.octave < 11 {
						kbd.instrument.FlushNotes()
						kbd.activeKeys = make([]string, 0)
						kbd.octave++
					}
				}
				switch e.VirtualKey.String() {
				case "1":
					kbd.running = false
					stopFn()
					return
				}

			case gowinkey.KeyUp:
				switch e.VirtualKey.String() {
				case "A", "W", "E", "T", "Y", "U", "S", "D", "F", "G", "H", "J":
					for idx, keys := range kbd.activeKeys {
						if keys == e.String() {
							kbd.activeKeys = append(kbd.activeKeys[:idx], kbd.activeKeys[idx+1:]...)
							kbd.instrument.RemoveNote((12 * kbd.octave) + keyVals[e.String()])
						}
					}
				}
			}
		}
	}()
}

func (kbd *KeyboardListener) SetInstrument(inst Instrument) {
	kbd.instrument = inst
}

func (kbd *KeyboardListener) Err() error {
	return nil
}

func (kbd *KeyboardListener) Stream(samples [][2]float64) (n int, ok bool) {
	if !kbd.running {
		fmt.Println("Off")
		return 0, false
	}
	n, _ = kbd.instrument.Stream(samples)
	return n, true
}
