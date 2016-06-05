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

	service := flag.String("s", "", "service for which generate data: apns, fcm (fcm: Firebase Cloud Messaging - former GCM)")
	appTotal := flag.Int("a", 4, "number of apps to generate")
	instancesPerApp := flag.Int("i", 150, "number of instances per app to generate")
	unregisteredPercent := flag.Float64("u", 10.0, "percent of instances with unregistered status")
	registrationIDPerInstanceMax := flag.Int("r", 10, "maximum number of registrationIDs per app (android only)")

	// TODO(szpakas): add custom usage

	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("Two arguments required: apps and instances file.")
	}

	appF, err := os.Create(args[0])
	check(err)
	defer appF.Close()

	insF, err := os.Create(args[1])
	check(err)
	defer insF.Close()

	switch *service {
	case "fcm":
		e := fcm.NewJSONExporter(appF, insF)
		g := fcm.NewGenerator(*appTotal, *instancesPerApp, *unregisteredPercent, *registrationIDPerInstanceMax)
		g.Generate(e)
	case "apns":
		e := apns.NewJSONExporter(appF, insF)
		g := apns.NewGenerator(*appTotal, *instancesPerApp, *unregisteredPercent)
		g.Generate(e)
	default:
		log.Fatalf("Unknown platform requested: %s", *service)
	}
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
