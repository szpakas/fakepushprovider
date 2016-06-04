package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gamegos/gcmlib"
	"github.com/uber-go/zap"
	"github.com/vrischmann/envconfig"

	"github.com/szpakas/fakepushprovider/android"
)

const (
	// ConfigAppPrefix prefixes all ENV values used to config the program.
	ConfigAppPrefix = "APP"
)

type config struct {
	// AppsFile is path to test data file with apps definition.
	AppsFile      string

	// InstancesFile is path to test data file with instances definition.
	// Beware: apps from apps file have to match
	InstancesFile string

	// GCMEndpoint is endpoint of the GCM server (or fake).
	// example: "http://localhost:8080"
	GCMEndpoint string

	// LogLevel is a minimal log severity required for the message to be logged.
	// Valid levels: [all, debug, info, warn, error, fatal, panic, none].
	LogLevel string `envconfig:"default=info"`
}

func main() {
	lgr := zap.NewJSON()

	// - config from env
	cfg := &config{}
	if err := envconfig.InitWithPrefix(&cfg, ConfigAppPrefix); err != nil {
		lgr.Fatal(err.Error())
	}

	// -- logging
	var logLevel zap.Level
	if err := logLevel.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		lgr.Fatal(err.Error())
	}

	lgr.SetLevel(logLevel)
	lgr.Debug(fmt.Sprintf("Parsed config from env => %+v", *cfg))

	lgr.Info("starting")

	// -- apps
	storage := android.NewMemoryStorage()
	mapper := android.NewMemoryMapper()
	importer := android.NewJSONImporter(storage, mapper)

	appsFile, err := os.Open(cfg.AppsFile)
	if err != nil {
		lgr.Fatal(err.Error())
	}
	importer.ImportApps(appsFile)
	lgr.Info("import:apps")
	appsFile.Close()

	// -- instances
	iFile, err := os.Open(cfg.InstancesFile)
	if err != nil {
		lgr.Fatal(err.Error())
	}

	scn := bufio.NewScanner(iFile)
	for scn.Scan() {
		insExp := new(android.InstanceExported)
		_ = json.Unmarshal(scn.Bytes(), insExp)
		app, err := storage.AppLoad(insExp.AppID)
		if err == android.ErrElementNotFound {
			continue
		}

		message := &gcmlib.Message{
			To: string(insExp.CanonicalID),
			Notification: &gcmlib.Notification{
				Title: "title",
				Body:  "body",
			},
		}

		client := gcmlib.NewClient(gcmlib.Config{
			APIKey:       app.ApiKey,
			MaxRetries:   -1,
			SendEndpoint: cfg.GCMEndpoint,
		})

		fields := []zap.Field{
			zap.String("mobile:app:id", app.ID),
			zap.String("mobile:instance:token", string(insExp.CanonicalID)),
		}
		res, gcmErr := client.Send(message)

		if gcmErr != nil {
			lgr.Error(
				"message:failed",
				append(
					fields,
					zap.Int("response:error:code", int(gcmErr.Code())),
					zap.String("response:error:message", gcmErr.Error()),
				)...,
			)
			continue
		}
		lgr.Debug(
			"message:sent",
			append(
				fields,
				zap.Int("response:success:count", int(res.Success)),
				zap.Int("response:failure:count", int(res.Failure)),
			)...,
		)
	}
}
