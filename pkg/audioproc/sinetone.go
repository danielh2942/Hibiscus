package audioproc

import (
	"math"

	"github.com/faiface/beep"
)

func SineTone(freq float64, durationSeconds float64) beep.Streamer {
	count := 0
	durationFrames := int(math.Floor(44100 * durationSeconds))
	var step float64 = (math.Pi * 2.0 * freq) / 44100
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			if count > durationFrames {
				return i + 1, false
			}
			temp := math.Sin(step * float64(count))
			samples[i][0] = temp
			samples[i][1] = temp
			count += 1
		}
		return len(samples), true
	})
}
