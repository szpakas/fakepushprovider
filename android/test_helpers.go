package android

import (
	"bytes"
	"testing"

	a "github.com/stretchr/testify/assert"
)

// -- helpers
func thAssertInstanceInMapper(t *testing.T, m *MemoryMapper, i *Instance) {
	for _, r := range i.RegistrationIDS {
		if !a.Contains(t, m.regIDs, r, "missing registrationID: %s, instanceID: %s", r, i.ID) {
			continue
		}
		a.Equal(t, i, m.regIDs[r], "instance does not match on registrationID: %s, instanceID: %s", r, i.ID)
	}
}

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

func TSMemoryStorageWitAppsAndInstancesSetup() (*MemoryStorage, func()) {
	s, closer := tsMemoryStorageSetup()
	s.AppSave(&TFAppA)
	s.AppSave(&TFAppB)
	s.AppSave(&TFAppC)

	s.InstanceSave(&TFInsAA)
	s.InstanceSave(&TFInsAB)
	s.InstanceSave(&TFInsAC)
	s.InstanceSave(&TFInsAZ)
	s.InstanceSave(&TFInsBA)
	s.InstanceSave(&TFInsBB)
	s.InstanceSave(&TFInsBC)

	return s, closer
}

func tsMapperSetup() (*MemoryMapper, func()) {
	m := NewMemoryMapper()
	closer := func() {}
	return m, closer
}

func tsGeneratorSetup(at, ipa, rpi int) (*Generator, func()) {
	m := NewGenerator(at, ipa, rpi)
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
