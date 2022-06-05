package main

import (
	"log"
	"time"

	"github.com/danielh2942/hibiscus/pkg/jingle"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

func main() {
	//audioproc.Demo()
	//PlayJingleSave("junk.json")
	PlayJingleSynth("enya.json")
}

func PlayJingleSynth(path string) {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))
	done := make(chan bool)
	if q, err := jingle.LoadInWavetable(path); err != nil {
		log.Panic("ERROR:", err)
	} else {
		q.AddNote(60)
		q.AddNote(57)
		q.AddNote(53)
		speaker.Play(beep.Seq(q, beep.Callback(func() {
			done <- true
		})))
	}
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
