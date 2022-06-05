package jingle

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
	"path"
	"sort"

	"github.com/danielh2942/hibiscus/pkg/audioproc"
)

// Composer2Hib converts a file stored in jingle format to one that is understood by Hibuscus
// Jingle was my previous adventure in audio processing so this will not be useful for anyone else
func Composer2Hib(filename string, bpm int) (queue *audioproc.Queue, err error) {
	framesPerSecond := float64(bpm/60) * float64(32)
	var mQueue audioproc.Queue
	filepath := path.Join(".", filename)
	dat, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer dat.Close()
	bytedata, _ := ioutil.ReadAll(dat)

	var composerData map[int][][3]int
	json.Unmarshal(bytedata, &composerData)
	hibMap := make(map[int]int)
	// Currently this only understands single instrument/monophonic Composer tracks
	// I will add poly support eventually
	keys := make([]int, 0)
	for k := range composerData {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		v := composerData[k]
		if v[0][2] == 1 { // Start note
			hibMap[v[0][1]] = k
		}
		if v[0][2] == 0 { // End Note
			start := hibMap[v[0][1]]
			len := float64(k-start) / framesPerSecond
			freq := 440.0 * math.Pow(2, (float64(v[0][1])-69)/12.0)
			mQueue.Add(audioproc.SineTone(freq, len))
			delete(composerData, k) // Clear key value
		}
	}
	return &mQueue, nil
}
