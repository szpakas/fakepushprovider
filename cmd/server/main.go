package main

//noinspection SpellCheckingInspection
import (
	"fmt"
	"net/http"
	"os"

	"github.com/uber-go/zap"
	"github.com/vrischmann/envconfig"

	"github.com/szpakas/fakepushprovider/android"
	ahttp "github.com/szpakas/fakepushprovider/android/http"
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

	// HTTPHost is address on which HTTP server endpoint is listening.
	HTTPHost string `envconfig:"default=0.0.0.0"`

	// HTTPPort is a port number on which HTTP server endpoint is listening.
	HTTPPort int `envconfig:"default=8080"`

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

	storage := android.NewMemoryStorage()
	mapper := android.NewMemoryMapper()
	importer := android.NewJSONImporter(storage, mapper)

	appsFile, err := os.Open(cfg.AppsFile)
	if err != nil {
		lgr.Fatal(err.Error())
	}
	importer.ImportApps(appsFile)
	lgr.Info("import: apps")
	appsFile.Close()

	instancesFile, err := os.Open(cfg.InstancesFile)
	if err != nil {
		lgr.Fatal(err.Error())
	}
	importRep := importer.ImportInstances(instancesFile)
	lgr.Info("import: instances")
	instancesFile.Close()

	lgr.Debug(fmt.Sprintf("import:instances:report => %+v", importRep))
	lgr.Debug(fmt.Sprintf("storage:report => %+v", storage.Report()))

	listenOn := fmt.Sprintf("%s:%d", cfg.HTTPHost, cfg.HTTPPort)
	lgr.Info("start listening", zap.String("host", cfg.HTTPHost), zap.Int("port", cfg.HTTPPort))

	h := ahttp.NewHandler(storage)
	m := ahttp.LoggingMiddleware{
		Handler: h,
		Logger:  lgr,
	}
	server := &http.Server{Addr: listenOn, Handler: &m}
	if err := server.ListenAndServe(); err != nil {
		lgr.Fatal(err.Error())
	}
}
