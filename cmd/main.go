package main

import (
	"fmt"
	"log"
	"time"

	"github.com/danielh2942/hibiscus/pkg/audioproc"
	"github.com/danielh2942/hibiscus/pkg/jingletools"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

func main() {
	fmt.Println("Hello World!")
	audioproc.Demo()
	//PlayJingleSave("junk.json")
}

func PlayJingleSave(path string) {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))
	done := make(chan bool)
	// Play Jingle save file
	if q, err := jingletools.Composer2Hib(path, 60); err != nil {
		log.Panic("ERROR:", err)
	} else {
		speaker.Play(beep.Seq(q, beep.Callback(func() {
			done <- true
		})))
	}
	<-done
}
