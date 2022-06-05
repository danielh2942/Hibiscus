package jingle

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

// LoadInWavetable reads in a jingle format wavetable and loads it into memory
func LoadInWavetable(filename string) (Instrument, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	pathToFile := path.Join(cwd, filename)
	dat, err := os.Open(pathToFile)
	if err != nil {
		return nil, err
	}
	defer dat.Close()
	bytedata, _ := ioutil.ReadAll(dat)

	var wt Wavetable
	err = json.Unmarshal(bytedata, &wt)
	if err != nil {
		return nil, err
	}
	// Save future calculation
	wt.length = len(wt.Audiodata)
	return &wt, nil
}

type Jingle struct {
	sequence    map[int][][3]int
	instruments [4]Instrument
}

func (j *Jingle) LoadSequence(newSequence map[int][][3]int) {
	j.sequence = newSequence
}

func (j *Jingle) LoadInstrument(instrument Instrument, pos uint) error {
	if pos > 3 {
		return errors.New("Jingle only supports up to 4 instrument tracks")
	}
	j.instruments[pos] = instrument
	return nil
}
