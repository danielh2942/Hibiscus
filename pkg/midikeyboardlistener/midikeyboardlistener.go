package midikeyboardlistener

import (
	"log"

	"github.com/karalabe/usb"
)

const (
	vendor  uint16 = 0x2467
	product uint16 = 0x203A
)

func Test() {
	devices, err := usb.EnumerateRaw(vendor, product)
	if err != nil {
		log.Fatalln("FATAL ERROR:", err)
	}
	log.Printf("%d devices found", len(devices))
}
