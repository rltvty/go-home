package default_loop

import (
	"fmt"
	"github.com/rltvty/go-home/dmx/astronomy"
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
	return fmt.Sprintf("R: %v\tG: %v\tB: %v\tW: %v\tA: %v\tUV: %v\t", c.Red, c.Green, c.Blue, c.White, c.Amber, c.UV)
}

type Setting struct {
	Color Color
	Time time.Time
}

func preDawnColor(now time.Time, events astronomy.Events) Color {
	red := 10
	if localSeconds(now) > localSeconds(events.Dawn) - 3600 {
		secondsUntilDawn := localSeconds(events.Dawn) - localSeconds(now)
		ratio := float32(3600 - secondsUntilDawn)/3600.0
		red += int(ratio * 100)
	}

	return Color{
		Red:   byte(red),
		Blue:  0,
		Green: 0,
		White: 0,
		Amber: 0,
		UV:    0,
	}
}

func wakeColor(now time.Time, events astronomy.Events) Color {
	return Color{
		Red:   0,
		Blue:  100,
		Green: 100,
		White: 0,
		Amber: 0,
		UV:    255,
	}
}

func morningColor(now time.Time, events astronomy.Events) Color {
	return Color{
		Red:   0,
		Blue:  255,
		Green: 255,
		White: 0,
		Amber: 0,
		UV:    255,
	}
}

func afternoonColor(now time.Time, events astronomy.Events) Color {
	return Color{
		Red:   0,
		Blue:  255,
		Green: 0,
		White: 255,
		Amber: 0,
		UV:    255,
	}
}

func eveningColor(now time.Time, events astronomy.Events) Color {
	return Color{
		Red:   255,
		Blue:  100,
		Green: 0,
		White: 0,
		Amber: 0,
		UV:    0,
	}
}


func nightColor(now time.Time, events astronomy.Events) Color {
	return Color{
		Red:   150,
		Blue:  0,
		Green: 0,
		White: 0,
		Amber: 0,
		UV:    150,
	}
}

//add some smoothing, like average over the last 1000 seconds or something
func Program(now time.Time, events astronomy.Events) (Color, string) {
	switch {
	case timeBefore(now, events.Dawn):
		return preDawnColor(now, events), "preDawn"
	case timeBefore(now, events.SunRise):
		return wakeColor(now, events), "wake"
	case timeBefore(now, events.SunPeak):
		return morningColor(now, events), "morning"
	case timeBefore(now, events.SunSet):
		return afternoonColor(now, events), "afternoon"
	case timeBefore(now, events.Dusk):
		return eveningColor(now, events), "evening"
	default:
		return nightColor(now, events), "night"
	}
}	

func timeBefore(a time.Time, b time.Time) bool {
	return localSeconds(a) < localSeconds(b)
}

func localSeconds(a time.Time) int {
	a = a.Local()
	return (a.Hour() * 3600) + (a.Minute() * 60) + a.Second()
}