package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/szpakas/fakepushprovider/apns"
	"github.com/szpakas/fakepushprovider/fcm"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	appTotal := flag.Int("a", 4, "number of apps to generate")
	instancesPerApp := flag.Int("i", 150, "number of instances per app to generate")
	unregisteredPercent := flag.Float64("u", 10.0, "percent of instances with unregistered status")
	registrationIDPerInstanceMax := flag.Int("r", 10, "maximum number of registrationIDs per app (android only)")

	// TODO(szpakas): add custom usage

	flag.Parse()
	args := flag.Args()
	if len(args) != 3 {
		log.Fatal("Three arguments required: service (apns or fcm), apps and instances file.")
	}

	appF, err := os.Create(args[1])
	check(err)
	defer appF.Close()

	insF, err := os.Create(args[2])
	check(err)
	defer insF.Close()

	s := args[0]
	switch s {
	case "fcm":
		e := fcm.NewJSONExporter(appF, insF)
		g := fcm.NewGenerator(*appTotal, *instancesPerApp, *unregisteredPercent, *registrationIDPerInstanceMax)
		g.Generate(e)
	case "apns":
		e := apns.NewJSONExporter(appF, insF)
		g := apns.NewGenerator(*appTotal, *instancesPerApp, *unregisteredPercent)
		g.Generate(e)
	default:
		log.Fatalf("Unknown platform requested: %s", s)
	}
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
