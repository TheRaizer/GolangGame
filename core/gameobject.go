package core

import (
	"log"
	"os/exec"
)

type GameObject interface {
	OnInit()
	OnUpdate(dt uint64)
	GetID() string
}

type BaseGameObject struct{}

func (obj BaseGameObject) GetID() string {
	uuid, err := exec.Command("uuidgen").Output()

	if err != nil {
		log.Fatal(err)
	}

	return string(uuid)
}
