package apns

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	a "github.com/stretchr/testify/assert"

	"github.com/szpakas/fakepushprovider/common"
)

// -- suite setup
func tsMemoryStorageSetup() (*MemoryStorage, func()) {
	s := NewMemoryStorage()
	closer := func() {}
	return s, closer
}

func tsMemoryStorageWitAppsSetup() (*MemoryStorage, func()) {
	s, closer := tsMemoryStorageSetup()
	s.AppSave(&TFAppA)
	s.AppSave(&TFAppB)
	s.AppSave(&TFAppC)
	return s, closer
}

func tsMemoryStorageWitAppsAndInstancesSetup() (*MemoryStorage, func()) {
	s, closer := tsMemoryStorageSetup()
	s.AppSave(&TFAppA)
	s.AppSave(&TFAppB)
	s.AppSave(&TFAppC)

	s.InstanceSave(&TFInsAA)
	s.InstanceSave(&TFInsAB)
	s.InstanceSave(&TFInsAC)
	s.InstanceSave(&TFInsAD)
	s.InstanceSave(&TFInsAZ)

	return s, closer
}

func tsMapperSetup() (*MemoryMapper, func()) {
	m := NewMemoryMapper()
	closer := func() {}
	return m, closer
}

func tsGeneratorSetup(at, ipa int, up float64) (*Generator, func()) {
	m := NewGenerator(at, ipa, up)
	closer := func() {}
	return m, closer
}

func tsImporterSetup() (*JSONImporter, *MemoryStorage, *MemoryMapper, func()) {
	s := NewMemoryStorage()
	m := NewMemoryMapper()
	i := NewJSONImporter(s, m)
	return i, s, m, func() {}
}

func tsExporterSetup() (*JSONExporter, *bytes.Buffer, *bytes.Buffer, func()) {
	aw := &bytes.Buffer{}
	iw := &bytes.Buffer{}
	e := NewJSONExporter(aw, iw)
	return e, aw, iw, func() {}
}

// -- helpers
func thAssertInstanceInMapper(t *testing.T, m *MemoryMapper, i *Instance) {
	a.Contains(t, m.tokens, i.Token, "missing token: %s, instanceID: %s", i.Token, i.ID)
}

// -- http helpers
func tsServerSetup(t *testing.T, symbol string) (*MemoryStorage, *Handler, *httptest.Server, *httpexpect.Expect, func()) {
	// -- storage
	st, stCloser := tsMemoryStorageWitAppsAndInstancesSetup()

	// -- handler
	h := NewHandler(st)

	srv := httptest.NewServer(h)

	// -- client test helper
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  srv.URL,
		Client:   http.DefaultClient,
		Reporter: &common.THAssertReporter{a.New(t), symbol},
	})

	closer := func() {
		stCloser()
		srv.Close()
	}

	return st, h, srv, e, closer
}
