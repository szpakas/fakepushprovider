package apns

import (
	"fmt"
	"math/rand"
	"time"
)

type Generator struct {
	AppTotal            int
	InstancesPerApp     int
	UnregisteredPercent float64
}

func NewGenerator(at, ipa int, up float64) *Generator {
	return &Generator{
		AppTotal:            at,
		InstancesPerApp:     ipa,
		UnregisteredPercent: up,
	}
}

func (g *Generator) Generate(s Storer) {
	var state InstanceState
	var totalCnt, unregisteredCnt float64
	pE := g.UnregisteredPercent / 100.0

	for aCnt := 0; aCnt < g.AppTotal; aCnt++ {
		app := App{
			ID:       fmt.Sprintf("appId-%d", aCnt+1),
			BundleID: fmt.Sprintf("bundleId-%d", aCnt+1),
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
				ID:    fmt.Sprintf("%s-instanceId-%d", app.ID, iCnt+1),
				State: state,
				App:   &app,
			}
			if state == InstanceStateUnregistered {
				// random time in last 8-72 hours
				ins.LastSeen = time.Now().Unix() - (rand.Int63n(72-8)+8)*3600
			}
			ins.Token = Token(fmt.Sprintf("%s-token", ins.ID))

			s.InstanceSave(&ins)
		}
	}
}
