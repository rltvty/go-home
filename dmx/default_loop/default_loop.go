package default_loop

import (
	"fmt"
	"github.com/rltvty/go-home/dmx/astronomy"
	"os"
	"time"
)

//Dawn
//SunRise
//SunPeak
//SunSet
//Dusk

type Color struct {
	Red byte
	Blue byte
	Green byte
	White byte
	Amber byte
	UV byte
}

func (c Color) String() string {
	return fmt.Sprintf("R: %v\tB: %v\tG: %v\tW: %v\tA: %v\tUV: %v\t", c.Red, c.Blue, c.Green, c.White, c.Amber, c.UV)
}

type Setting struct {
	Color Color
	Time time.Time
}

func nightColor(events astronomy.Events, currentTime time.Time) Color {
	if currentTime.After(events.Dusk) && currentTime.Before(events.Dawn) {
		//in current night
	} else if currentTime.Add(time.Duration.Hours())


}

func Program(events astronomy.Events) Color {
	t = time.Now()

}