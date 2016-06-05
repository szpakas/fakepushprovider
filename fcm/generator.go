package fcm

import (
	"fmt"
	"math/rand"
)

type Generator struct {
	AppTotal                     int
	InstancesPerApp              int
	UnregisteredPercent          float64
	RegistrationIDPerInstanceMax int
}

func NewGenerator(at, ipa int, up float64, rpi int) *Generator {
	return &Generator{
		AppTotal:                     at,
		InstancesPerApp:              ipa,
		UnregisteredPercent:          up,
		RegistrationIDPerInstanceMax: rpi,
	}
}

func (g *Generator) Generate(s Storer) {
	var state InstanceState
	var totalCnt, unregisteredCnt float64
	pE := g.UnregisteredPercent / 100.0

	for aCnt := 0; aCnt < g.AppTotal; aCnt++ {
		app := App{
			ID:       fmt.Sprintf("appId-%d", aCnt+1),
			SenderID: fmt.Sprintf("senderId-%d", aCnt+1),
			ApiKey:   fmt.Sprintf("apiKey-%d", aCnt+1),
		}

		s.AppSave(&app)

		totalCnt = 0
		unregisteredCnt = 0
		for iCnt := 0; iCnt < g.InstancesPerApp; iCnt++ {
			totalCnt++

			state = InstanceStateRegistered
			// we are looking for exact percentage of unregistered instances
			if (totalCnt*pE - unregisteredCnt) >= 1 {
				state = InstanceStateUnregistered
				unregisteredCnt++
			}

			ins := Instance{
				ID:              fmt.Sprintf("%s-instanceId-%d", app.ID, iCnt+1),
				State:           state,
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
