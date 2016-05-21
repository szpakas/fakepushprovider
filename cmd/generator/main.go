package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/szpakas/fakepushprovider/android"
)

type omit *struct{}

type InstanceWrapped struct {
	android.Instance

	AppID string
	App   omit `json:"App,omitempty"`
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var (
		appTotal                     = 2
		instancesPerApp              = 3
		registrationIDPerInstanceMax = 10
	)

	appF, err := os.Create("tmp/apps.json")
	check(err)
	defer appF.Close()
	ae := json.NewEncoder(appF)

	insF, err := os.Create("tmp/instances.json")
	check(err)
	defer insF.Close()
	ie := json.NewEncoder(insF)

	for aCnt := 0; aCnt < appTotal; aCnt++ {
		app := android.App{
			ID:       fmt.Sprintf("appId-%d", aCnt+1),
			SenderID: fmt.Sprintf("senderId-%d", aCnt+1),
			ApiKey:   fmt.Sprintf("apiKey-%d", aCnt+1),
		}

		_ = ae.Encode(&app)

		for iCnt := 0; iCnt < instancesPerApp; iCnt++ {
			ins := InstanceWrapped{
				Instance: android.Instance{
					ID:              fmt.Sprintf("%s-instanceId-%d", app.ID, iCnt+1),
					State:           android.InstanceStateRegistered,
					RegistrationIDS: make([]android.RegistrationID, rand.Intn(registrationIDPerInstanceMax-1)+1),
				},
				AppID: app.ID,
			}
			for rCnt := 0; rCnt < cap(ins.RegistrationIDS); rCnt++ {
				ins.RegistrationIDS[rCnt] = android.RegistrationID(fmt.Sprintf("%s-regId-%d", ins.ID, rCnt+1))
			}
			ins.CanonicalID = ins.RegistrationIDS[rand.Intn(len(ins.RegistrationIDS))]

			_ = ie.Encode(&ins)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
