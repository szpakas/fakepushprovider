package main

//noinspection SpellCheckingInspection
import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/http2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/zap"
	"github.com/vrischmann/envconfig"

	"github.com/szpakas/fakepushprovider/apns"
	"github.com/szpakas/fakepushprovider/common"
	"github.com/szpakas/fakepushprovider/fcm"
	fhttp "github.com/szpakas/fakepushprovider/fcm/http"
)

const (
	// ConfigAppPrefix prefixes all ENV values used to config the program.
	ConfigAppPrefix = "APP"
)

type config struct {
	// Service sets one of the supported services
	// Valid services: fcm, apns
	Service string

	// AppsFile is path to test data file with apps definition.
	AppsFile string

	// InstancesFile is path to test data file with instances definition.
	// Beware: apps from apps file have to match
	InstancesFile string

	// HTTPHost is address on which HTTP server endpoint is listening.
	HTTPHost string `envconfig:"default=0.0.0.0"`

	// HTTPPort is a port number on which HTTP server endpoint is listening.
	HTTPPort int `envconfig:"default=8080"`

	// MetricsHost is address on which metric HTTP server endpoint is listening.
	MetricsHost string `envconfig:"default=0.0.0.0"`

	// MetricsPort is a port number on which metric HTTP server endpoint is listening.
	MetricsPort int `envconfig:"default=9999"`

	// LogLevel is a minimal log severity required for the message to be logged.
	// Valid levels: [all, debug, info, warn, error, fatal, panic, none].
	LogLevel string `envconfig:"default=info"`

	// APNSCertFile is path to file with APNS SSL cert in PEM format
	APNSCertFile string `envconfig:"optional"`

	// APNSKeyFile is path to file with APNS SSL cert in PEM format
	APNSKeyFile string `envconfig:"optional"`
}

func main() {
	lgr := zap.NewJSON()

	// TODO(szpakas): add possibility to enforce seed
	rand.Seed(time.Now().UTC().UnixNano())

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

	appsFile, err := os.Open(cfg.AppsFile)
	if err != nil {
		lgr.Fatal(err.Error())
	}

	instancesFile, err := os.Open(cfg.InstancesFile)
	if err != nil {
		lgr.Fatal(err.Error())
	}

	go func() {
		listenOn := fmt.Sprintf("%s:%d", cfg.MetricsHost, cfg.MetricsPort)
		lgr.Info("metrics start listening", zap.String("host", cfg.MetricsHost), zap.Int("port", cfg.MetricsPort))
		s := &http.Server{
			Addr:    listenOn,
			Handler: prometheus.Handler(),
		}
		if err := s.ListenAndServe(); err != nil {
			lgr.Fatal(err.Error())
		}
	}()

	switch cfg.Service {
	case "fcm":
		serveFCM(cfg, lgr, appsFile, instancesFile)
	case "apns":
		serveAPNS(cfg, lgr, appsFile, instancesFile)
	default:
		lgr.Fatal("unknown service")
	}
}

func serveFCM(cfg *config, lgr zap.Logger, appsFile, instancesFile *os.File) {
	lgr = lgr.With(zap.String("service", "FCM"))

	storage := fcm.NewMemoryStorage()
	mapper := fcm.NewMemoryMapper()
	importer := fcm.NewJSONImporter(storage, mapper)

	appsFile, err := os.Open(cfg.AppsFile)
	if err != nil {
		lgr.Fatal(err.Error())
	}
	importer.ImportApps(appsFile)
	lgr.Info("import: apps")
	appsFile.Close()

	importRep := importer.ImportInstances(instancesFile)
	lgr.Info("import: instances")
	instancesFile.Close()

	lgr.Debug(fmt.Sprintf("import:instances:report => %+v", importRep))
	lgr.Debug(fmt.Sprintf("storage:report => %+v", storage.Report()))

	listenOn := fmt.Sprintf("%s:%d", cfg.HTTPHost, cfg.HTTPPort)
	lgr.Info("start listening", zap.String("host", cfg.HTTPHost), zap.Int("port", cfg.HTTPPort))

	h := fhttp.NewHandler(storage)
	m := common.LoggingMiddleware{
		Handler: h,
		Logger:  lgr,
	}
	server := &http.Server{Addr: listenOn, Handler: &m}
	if err := server.ListenAndServe(); err != nil {
		lgr.Fatal(err.Error())
	}
}

func serveAPNS(cfg *config, lgr zap.Logger, appsFile, instancesFile *os.File) {
	lgr = lgr.With(zap.String("service", "APNS"))

	if _, err := os.Open(cfg.APNSCertFile); err != nil {
		lgr.Fatal(errors.Wrap(err, "APNSCertFile").Error())
	}
	if _, err := os.Open(cfg.APNSKeyFile); err != nil {
		lgr.Fatal(errors.Wrap(err, "APNSKeyFile").Error())
	}

	listenOn := fmt.Sprintf("%s:%d", cfg.HTTPHost, cfg.HTTPPort)
	lgr.Info("start listening", zap.String("host", cfg.HTTPHost), zap.Int("port", cfg.HTTPPort))

	storage := apns.NewMemoryStorage()
	mapper := apns.NewMemoryMapper()
	importer := apns.NewJSONImporter(storage, mapper)

	importer.ImportApps(appsFile)
	lgr.Info("import: apps")
	appsFile.Close()

	importRep := importer.ImportInstances(instancesFile)
	lgr.Info("import: instances")
	instancesFile.Close()

	lgr.Debug(fmt.Sprintf("import:instances:report => %+v", importRep))
	//lgr.Debug(fmt.Sprintf("storage:report => %+v", storage.Report()))

	h := apns.NewHandler(storage)
	mDelay := common.DelayMiddleware{
		Handler: h,
		DelayFn: func() int {
			d := 20 + rand.Intn(100)
			time.Sleep(time.Millisecond * time.Duration(d))
			return d
		},
	}
	mLog := common.LoggingMiddleware{
		Handler: &mDelay,
		Logger:  lgr,
	}
	mInstr := common.NewInstrMiddleware(&mLog, "apns")
	srv := &http.Server{Addr: listenOn, Handler: mInstr}
	http2.ConfigureServer(srv, &http2.Server{})

	if err := srv.ListenAndServeTLS(cfg.APNSCertFile, cfg.APNSKeyFile); err != nil {
		lgr.Fatal(err.Error())
	}
}
