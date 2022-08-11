//go:build windows
// +build windows

package keyboardlistener

// Windows Keyboard Listener
import (
	"fmt"
	"strconv"
	"time"

	"github.com/danielh2942/hibiscus/pkg/jingle"
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
	activeKeys  []string
	instruments []jingle.Instrument
	instrument  int
	octave      int
	running     bool
}

func (kbd *KeyboardListener) StartMonitor() {
	if len(kbd.instruments) == 0 {
		kbd.running = false
	}
	kbd.instrument = 0
	kbd.running = true
	for _, inst := range kbd.instruments {
		inst.Init()
	}
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
		arpRate := 0
		events, stopFn := gowinkey.Listen()
	eventloop:
		for e := range events {
			switch e.State {
			case gowinkey.KeyDown:
				switch e.String() {
				case "A", "W", "E", "T", "Y", "U", "S", "D", "F", "G", "H", "J":
					for _, keys := range kbd.activeKeys {
						if keys == e.String() {
							continue eventloop
						}
					}
					kbd.activeKeys = append(kbd.activeKeys, e.String())
					kbd.instruments[kbd.instrument].AddNote((12 * kbd.octave) + keyVals[e.String()])

				case "Z":
					if kbd.octave > 0 {
						kbd.instruments[kbd.instrument].FlushNotes()
						kbd.activeKeys = make([]string, 0)
						kbd.octave--
					}

				case "X":
					if kbd.octave < 11 {
						kbd.instruments[kbd.instrument].FlushNotes()
						kbd.activeKeys = make([]string, 0)
						kbd.octave++
					}
				case "V":
					fmt.Println("Arp Triggered")
					kbd.instruments[kbd.instrument].ToggleArpeggio()
				case "B":
					fmt.Println("Arp State Toggled")
					switch kbd.instruments[kbd.instrument].GetArpeggioState() {
					case jingle.EArpUp:
						fmt.Println("Arp Set to Down")
						kbd.instruments[kbd.instrument].SetArpeggioState(jingle.EArpDown)
					case jingle.EArpDown:
						fmt.Println("Arp Set to Bidirectional")
						kbd.instruments[kbd.instrument].SetArpeggioState(jingle.EArpBidirectional)
					case jingle.EArpBidirectional:
						fmt.Println("Arp Set to Up")
						kbd.instruments[kbd.instrument].SetArpeggioState(jingle.EArpUp)
					}
				case "M":
					if arpRate < 60 {
						arpRate++
						fmt.Println("Set Arp to", arpRate, "Hz")
						kbd.instruments[kbd.instrument].ArpeggioRate(float64(arpRate))
					}

				case "N":
					if arpRate > 0 {
						arpRate--
						fmt.Println("Set Arp to", arpRate, "Hz")
						kbd.instruments[kbd.instrument].ArpeggioRate(float64(arpRate))
					}
				}
				switch e.VirtualKey.String() {
				case "TAB":
					kbd.running = false
					stopFn()
					return
				case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
					idx, _ := strconv.Atoi(e.VirtualKey.String())
					if idx == 0 {
						idx = 9
					} else {
						idx--
					}
					kbd.instruments[kbd.instrument].FlushNotes()
					kbd.instrument = idx % len(kbd.instruments)
				}

			case gowinkey.KeyUp:
				switch e.VirtualKey.String() {
				case "A", "W", "E", "T", "Y", "U", "S", "D", "F", "G", "H", "J":
					for idx, keys := range kbd.activeKeys {
						if keys == e.String() {
							kbd.activeKeys = append(kbd.activeKeys[:idx], kbd.activeKeys[idx+1:]...)
							kbd.instruments[kbd.instrument].RemoveNote((12 * kbd.octave) + keyVals[e.String()])
						}
					}
				}
			}
		}
	}()
}

func (kbd *KeyboardListener) AddInstrument(inst jingle.Instrument) {
	if len(kbd.instruments) == 0 {
		kbd.instruments = make([]jingle.Instrument, 1)
		kbd.instruments[0] = inst
		return
	}
	kbd.instruments = append(kbd.instruments, inst)
}

func (kbd *KeyboardListener) Err() error {
	return nil
}

func (kbd *KeyboardListener) Stream(samples [][2]float64) (n int, ok bool) {
	if !kbd.running {
		fmt.Println("Off")
		return 0, false
	}
	n, _ = kbd.instruments[kbd.instrument].Stream(samples)
	return n, true
}
