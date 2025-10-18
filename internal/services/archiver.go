package services

import (
	"fmt"
	"time"
)

type Archiver struct {
	state    string
	progress float32
}

func NewArchiver() *Archiver {
	return &Archiver{
		state:    "Waiting",
		progress: 0,
	}
}

func (a *Archiver) Status() string {
	return a.state
}

func (a *Archiver) Progress() float32 {
	return a.progress
}

func (a *Archiver) Start() {
	go func() {
		if a.state != "Running" {
			a.state = "Running"
			a.progress = 0
		}

		ticker := time.NewTicker(300 * time.Millisecond)
		defer ticker.Stop()

		for {
			<-ticker.C
			a.progress += 0.01
			fmt.Println(a.progress)
			if a.progress >= 1.0 {
				a.state = "Finished"
				return
			}
		}
	}()
}
