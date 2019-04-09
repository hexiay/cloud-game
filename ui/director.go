package ui

import (
	"image"
	"log"
	"time"

	"github.com/giongto35/cloud-game/nes"
)

type View interface {
	Enter()
	Exit()
	Update(t, dt float64)
}

type Director struct {
	// audio        *Audio
	view          View
	timestamp     float64
	imageChannel  chan *image.RGBA
	inputChannel  chan int
	closedChannel chan bool
	roomID        string
}

func NewDirector(roomID string, imageChannel chan *image.RGBA, inputChannel chan int, closedChannel chan bool) *Director {
	director := Director{}
	// director.audio = audio
	director.imageChannel = imageChannel
	director.inputChannel = inputChannel
	director.closedChannel = closedChannel
	director.roomID = roomID
	return &director
}

func (d *Director) SetView(view View) {
	if d.view != nil {
		d.view.Exit()
	}
	d.view = view
	if d.view != nil {
		d.view.Enter()
	}
	d.timestamp = float64(time.Now().Nanosecond()) / float64(time.Second)
}

func (d *Director) Step() {
	timestamp := float64(time.Now().Nanosecond()) / float64(time.Second)
	dt := timestamp - d.timestamp
	d.timestamp = timestamp
	if d.view != nil {
		d.view.Update(timestamp, dt)
	}
}

func (d *Director) Start(paths []string) {
	if len(paths) == 1 {
		d.PlayGame(paths[0])
	}
	d.Run()
}

func (d *Director) Run() {
L:
	for {
		// quit game

		select {
		// if there is event from close channel => the game is ended
		case <-d.closedChannel:
			break L
		default:
		}

		d.Step()
	}
	d.SetView(nil)
}

func (d *Director) PlayGame(path string) {
	// Generate hash that is indentifier of a room (game path + ropomID)
	hash, err := hashFile(path, d.roomID)
	if err != nil {
		log.Fatalln(err)
	}
	console, err := nes.NewConsole(path)
	if err != nil {
		log.Fatalln(err)
	}
	// Set GameView as current view
	d.SetView(NewGameView(d, console, path, hash, d.imageChannel, d.inputChannel))
}
