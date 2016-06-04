package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/szpakas/fakepushprovider/android"
)

// TODO: add configuration via env
func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var (
		appTotal                             = 2
		instancesPerApp                      = 150
		unregisteredPercent          float64 = 10
		registrationIDPerInstanceMax         = 10
	)

	appF, err := os.Create("tmp/apps.json")
	check(err)
	defer appF.Close()

	insF, err := os.Create("tmp/instances.json")
	check(err)
	defer insF.Close()

	e := android.NewJSONExporter(appF, insF)

	g := android.NewGenerator(appTotal, instancesPerApp, unregisteredPercent, registrationIDPerInstanceMax)
	g.Generate(e)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
