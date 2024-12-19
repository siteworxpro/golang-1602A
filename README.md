# 1602A Programmer Library For GO

Connect and display simple logic for a 1602A screen in a raspberry pi.

define your pinouts currently only 8-bit mode is supported
```go
screen.New(&screen.PinConfig{
		RS: rpio.Pin(2),
		RW: rpio.Pin(3),
		D0: rpio.Pin(14),
		D1: rpio.Pin(15),
		D2: rpio.Pin(18),
		D3: rpio.Pin(23),
		D4: rpio.Pin(24),
		D5: rpio.Pin(25),
		D6: rpio.Pin(8),
		D7: rpio.Pin(7),
		E:  rpio.Pin(21),
	})
```

Simple hello world clock
```go
package main

import (
	"git.s.int/rrise/raspberry-pi/screen/screen"
	"github.com/stianeikeland/go-rpio/v4"
	"time"
)

func main() {
	s := screen.New(&screen.PinConfig{
		RS: rpio.Pin(2),
		RW: rpio.Pin(3),
		D0: rpio.Pin(14),
		D1: rpio.Pin(15),
		D2: rpio.Pin(18),
		D3: rpio.Pin(23),
		D4: rpio.Pin(24),
		D5: rpio.Pin(25),
		D6: rpio.Pin(8),
		D7: rpio.Pin(7),
		E:  rpio.Pin(21),
	})

	s.Init()
	s.SetScreenFormat(true, false)
	s.SetDisplay(true, false, false)
	s.SetCursorPosition(2, false)
	s.WriteString("Hello World!")

	for {
		s.SetCursorPosition(4, true)
		s.WriteString(time.Now().Format(time.TimeOnly))
		time.Sleep(1 * time.Second)
	}

}
```