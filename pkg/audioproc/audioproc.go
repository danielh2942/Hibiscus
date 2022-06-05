package audioproc

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func Demo() {
	cwd, _ := os.Getwd()
	f, err := os.Open(filepath.Join(cwd, "../testaudio/test.mp3"))
	if err != nil {
		log.Fatalln("ERROR: ", err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}

func PlaySineTone() {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/20))

	done := make(chan bool)

	speaker.Play(beep.Seq(SineTone(880, 5), beep.Callback(func() {
		done <- true
	})))
	<-done
}
