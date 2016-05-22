package android

import (
	"fmt"
	"math/rand"
)

type Storer interface {
	AppSave(o *App) error
	InstanceSave(o *Instance) error
}

type Generator struct {
	AppTotal                     int
	InstancesPerApp              int
	RegistrationIDPerInstanceMax int
}

func NewGenerator(at, ipa, rpi int) *Generator {
	return &Generator{
		AppTotal:                     at,
		InstancesPerApp:              ipa,
		RegistrationIDPerInstanceMax: rpi,
	}
}

func (g *Generator) Generate(s Storer) {
	for aCnt := 0; aCnt < g.AppTotal; aCnt++ {
		app := App{
			ID:       fmt.Sprintf("appId-%d", aCnt+1),
			SenderID: fmt.Sprintf("senderId-%d", aCnt+1),
			ApiKey:   fmt.Sprintf("apiKey-%d", aCnt+1),
		}

		s.AppSave(&app)

		for iCnt := 0; iCnt < g.InstancesPerApp; iCnt++ {
			ins := Instance{
				ID:              fmt.Sprintf("%s-instanceId-%d", app.ID, iCnt+1),
				State:           InstanceStateRegistered,
				RegistrationIDS: make([]RegistrationID, rand.Intn(g.RegistrationIDPerInstanceMax-1)+1),
				App:             &app,
			}
			for rCnt := 0; rCnt < cap(ins.RegistrationIDS); rCnt++ {
				ins.RegistrationIDS[rCnt] = RegistrationID(fmt.Sprintf("%s-regId-%d", ins.ID, rCnt+1))
			}
			ins.CanonicalID = ins.RegistrationIDS[rand.Intn(len(ins.RegistrationIDS))]

			s.InstanceSave(&ins)
		}
	}
}
