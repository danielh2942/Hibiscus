package jingle

import (
	"errors"
	"math"
)

// Instrument interface to make it less taxing for the author to use drumkits or wavetables.
// They are the only instruments supported by Jingle
type Instrument interface {
	// Jingle specific audio engine stuff
	AddNote(note int) error
	RemoveNote(note int) error
	// Arpeggio
	Arpeggio(enabled bool)
	ArpeggioRate(freq float64)
	// beep.Streamers functions
	Err() error
	Stream(samples [][2]float64) (length int, ok bool)
}

func getNoteSteps() [128]float64 {
	return [128]float64{
		0.03125, 0.03311, 0.03508, 0.03717, 0.03938, 0.04172, 0.0442, 0.04683, 0.04961, 0.05256, 0.05569, 0.059,
		0.06251, 0.06622, 0.07016, 0.07433, 0.07875, 0.08344, 0.0884, 0.09365, 0.09922, 0.10512, 0.11137, 0.118,
		0.12501, 0.13245, 0.14032, 0.14867, 0.15751, 0.16687, 0.17679, 0.18731, 0.19844, 0.21024, 0.22275, 0.23599,
		0.25002, 0.26489, 0.28064, 0.29733, 0.31501, 0.33374, 0.35359, 0.37461, 0.39689, 0.42049, 0.44549, 0.47198,
		0.50005, 0.52978, 0.56129, 0.59466, 0.63002, 0.66749, 0.70718, 0.74923, 0.79378, 0.84098, 0.89099, 0.94397,
		1.0, 1.05957, 1.12257, 1.18932, 1.26004, 1.33497, 1.41435, 1.49845, 1.58756, 1.68196, 1.78197, 1.88793,
		2.0002, 2.11913, 2.24514, 2.37865, 2.52009, 2.66994, 2.8287, 2.99691, 3.17511, 3.36391, 3.56394, 3.77587,
		4.00039, 4.23827, 4.49029, 4.75729, 5.04018, 5.33988, 5.65741, 5.99381, 6.35022, 6.72783, 7.12789, 7.55173,
		8.00078, 8.47653, 8.98057, 9.51459, 10.08035, 10.67976, 11.31481, 11.98763, 12.70045, 13.45566, 14.25577,
		5.10346, 16.00156, 16.95307, 17.96115, 19.02917, 20.16071, 21.35952, 22.62963, 23.97526, 25.4009, 26.91131,
		28.51155, 30.20693, 32.00313, 33.90613, 35.9223, 38.05835, 40.32141, 42.71905, 45.25926, 47.95051,
	}
}

// Wavetable, this is the classic Polyphonic, Gated Jingle instrument wrapper
type Wavetable struct {
	Audiodata          []float64 `json:"audiodata"`
	length             int
	presentNotes       []int
	noteframePositions map[int]float64
	arpeggio           bool
	arpRate            float64
}

func (wt *Wavetable) AddNote(note int) error {
	if len(wt.presentNotes) == 0 {
		wt.presentNotes = make([]int, 1)
		wt.presentNotes[0] = note

		if len(wt.noteframePositions) == 0 {
			wt.noteframePositions = map[int]float64{note: 0.0}
		}
	}
	for _, x := range wt.presentNotes {
		if x == note {
			return errors.New("note already present")
		}
	}
	wt.presentNotes = append(wt.presentNotes, note)
	wt.noteframePositions[note] = 0
	return nil
}

func (wt *Wavetable) RemoveNote(note int) error {
	for i, x := range wt.presentNotes {
		if x == note {
			wt.presentNotes = append(wt.presentNotes[:i], wt.presentNotes[i+1:]...)
			delete(wt.noteframePositions, note)
			return nil
		}
	}

	return errors.New("note not present, ignoring")
}

// Err returns an error (not ever in this case)
func (wt *Wavetable) Err() error {
	return nil
}

// Arpeggio toggles the use of an arpeggio
func (wt *Wavetable) Arpeggio(enabled bool) {
	wt.arpeggio = enabled
}

// ArpeggioRate takes a frequency and acts accordingly,
// It ignores the sign on the number provided
func (wt *Wavetable) ArpeggioRate(freq float64) {
	// changes each frame :)
	wt.arpRate = math.Abs(freq) / 44100.0
}

func (wt *Wavetable) Stream(samples [][2]float64) (length int, more bool) {
	myNumPresentNotes := len(wt.presentNotes)
	if myNumPresentNotes == 0 {
		return 0, false
	}
	wtlen := float64(wt.length)
	stepTable := getNoteSteps()
	for i := range samples {
		samples[i][0] = 0
		samples[i][1] = 0
		for _, key := range wt.presentNotes {
			val := math.Mod(wt.noteframePositions[key], 1.0)
			sample := (wt.Audiodata[int(wt.noteframePositions[key])] * val) + (wt.Audiodata[int(wt.noteframePositions[key]+1)] * (1 - val))
			samples[i][0] += sample
			samples[i][1] += sample
			wt.noteframePositions[key] = math.Mod(wt.noteframePositions[key]+stepTable[key], wtlen-1)
		}
		samples[i][0] /= float64(len(wt.presentNotes))
		samples[i][1] /= float64(len(wt.presentNotes))
	}
	return len(samples), true
}
