package screen

import (
	"github.com/charmbracelet/log"
	"github.com/stianeikeland/go-rpio/v4"
	"strconv"
	"strings"
	"time"
)

type Screen struct {
	pinConfig    *PinConfig
	initialized  bool
	eightBitMode bool
}

type PinConfig struct {
	RS rpio.Pin
	RW rpio.Pin
	D0 rpio.Pin
	D1 rpio.Pin
	D2 rpio.Pin
	D3 rpio.Pin
	D4 rpio.Pin
	D5 rpio.Pin
	D6 rpio.Pin
	D7 rpio.Pin
	E  rpio.Pin
}

const clearDisplay = 0x01
const returnHome = 0x02

const setFunction = 0x20
const setFunction8Bit = 0x10
const setFunction2Line = 0x08
const setFunction5x11 = 0x04

const moveCursor = 0x10
const moveCursorRight = 0x04
const moveCursorLeft = 0x00

const displayFunction = 0x08
const displayOn = 0x04
const displayCursor = 0x02
const displayCursorBlink = 0x01

const setPosition = 0x80
const lineTwoHome = 0x40

var pinMap = map[int]*rpio.Pin{}

func New(config *PinConfig) *Screen {
	return &Screen{
		pinConfig: config,
	}
}

func (s *Screen) SetCursorPosition(cols uint32, secondLine bool) {
	var position uint32 = setPosition

	if secondLine {
		position = position + lineTwoHome
	}

	position = position + cols

	s.sendBits(position, false)
}

func (s *Screen) clearPins() {
	log.Infof("Clearing pins")

	s.pinConfig.E.Low()

	for _, pin := range pinMap {
		pin.Low()
	}
}

func (s *Screen) Init() {
	err := rpio.Open()
	if err != nil {
		panic(err)
	}

	s.pinConfig.D0.Output()
	s.pinConfig.D1.Output()
	s.pinConfig.D2.Output()
	s.pinConfig.D3.Output()
	s.pinConfig.D4.Output()
	s.pinConfig.D5.Output()
	s.pinConfig.D6.Output()
	s.pinConfig.D7.Output()
	s.pinConfig.RS.Output()
	s.pinConfig.RW.Output()
	s.pinConfig.E.Output()

	s.clearPins()
	s.ClearScreen()
	s.Home()

	s.initialized = true
}

func (s *Screen) ClearScreen() {
	s.sendBits(clearDisplay, false)
}

func (s *Screen) Home() {
	s.sendBits(returnHome, false)
}

func (s *Screen) SetScreenFormat(twoLines bool, fiveByEleven bool) {
	var format uint32 = setFunction

	format = format + setFunction8Bit

	if twoLines {
		format = format + setFunction2Line
	}
	if fiveByEleven {
		format = format + setFunction5x11
	}

	log.Infof("Setting format: %v", format)

	s.sendBits(format, false)
}

func (s *Screen) SetDisplay(on bool, cursor bool, cursorBlink bool) {
	var display uint32 = displayFunction

	if on {
		display = display + displayOn
	}

	if cursor {
		display = display + displayCursor
	}

	if cursorBlink {
		display = display + displayCursorBlink
	}

	log.Infof("Setting display: %v", display)

	s.sendBits(display, false)
}

func (s *Screen) CursorRight(number uint8) {
	for i := 0; i < int(number); i++ {
		s.sendBits(moveCursorRight+moveCursor, false)
	}
}

func (s *Screen) CursorLeft(number uint8) {
	for i := 0; i < int(number); i++ {
		s.sendBits(moveCursorLeft+moveCursor, false)
	}
}

func (s *Screen) decToBin(n uint32) string {
	bin := strconv.FormatInt(int64(n), 2)

	bin = strings.Repeat("0", 8-len(bin)) + bin

	return bin
}

func (s *Screen) WriteString(st string) {
	for _, c := range st {
		s.sendBits(uint32(c), true)
	}
}

func (s *Screen) sendBits(bits uint32, writeMode bool) {

	bin := s.decToBin(bits)

	err := rpio.Open()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = rpio.Close()
	}()

	log.Infof("Running command: %s", bin)

	s.pinConfig.RS.Low()
	if writeMode {
		s.pinConfig.RS.High()
	}
	s.pinConfig.RW.Low()

	for i, v := range bin {

		mode := rpio.Low
		if string(v) == "1" {
			mode = rpio.High
		}

		log.Infof("Writing command: %s", bin)

		// pins and bits are backwards
		// d7 d6 d5 d4 d3 d2 d1 d0
		// 0  1  2  3  4  5  6  7
		switch i {
		case 0:
			s.pinConfig.D7.Write(mode)
		case 1:
			s.pinConfig.D6.Write(mode)
		case 2:
			s.pinConfig.D5.Write(mode)
		case 3:
			s.pinConfig.D4.Write(mode)
		case 4:
			s.pinConfig.D3.Write(mode)
		case 5:
			s.pinConfig.D2.Write(mode)
		case 6:
			s.pinConfig.D1.Write(mode)
		case 7:
			s.pinConfig.D0.Write(mode)
		}
	}

	time.Sleep(500 * time.Microsecond)
	s.pinConfig.E.High()
	time.Sleep(500 * time.Nanosecond)
	s.pinConfig.E.Low()
	s.clearPins()
}
