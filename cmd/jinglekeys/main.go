package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"time"

	"github.com/danielh2942/hibiscus/pkg/jingle"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

func main() {
	fmt.Println("Press 0 - 9 to change instruments")
	fmt.Println("Press \"tab\" to close")
	fmt.Println("Z - octave down")
	fmt.Println("X - octave up")
	fmt.Println("B - Arp On")
	fmt.Println("N - Arp Freq Down")
	fmt.Println("M - Arp Freq Up")
	//audioproc.Demo()
	//PlayJingleSave("junk.json")
	files, err := ioutil.ReadDir("./jingleinstruments")
	if err != nil {
		log.Panic("ERROR:", err)
	}
	filenames := make([]string, len(files))
	for i, file := range files {
		filenames[i] = path.Join("./jingleinstruments", file.Name())
	}
	PlayJingleSynth(filenames)
	//jingle.KeyboardTest()
}

func PlayJingleSynth(paths []string) {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))
	done := make(chan bool)
	kbd := jingle.KeyboardListener{}
	for _, path := range paths {
		if q, err := jingle.LoadInWavetable(path); err != nil {
			log.Panic("ERROR:", err)
		} else {
			kbd.AddInstrument(q)
		}
	}
	kbd.StartMonitor()
	speaker.Play(beep.Seq(&kbd, beep.Callback(func() {
		done <- true
	})))
	<-done
}

func PlayJingleSave(path string) {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))
	done := make(chan bool)
	// Play Jingle save file
	if q, err := jingle.Composer2Hib(path, 60); err != nil {
		log.Panic("ERROR:", err)
	} else {
		speaker.Play(beep.Seq(q, beep.Callback(func() {
			done <- true
		})))
	}
	<-done
}
