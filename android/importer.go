package android

import (
	"bufio"
	"encoding/json"
	"io"
)

type ImporterStorage interface {
	AppLoad(id string) (*App, error)
	AppSave(o *App) error
	InstanceSave(o *Instance) error
}

type JSONImporter struct {
	Storage ImporterStorage
	Mapper  Mapper
}

func NewJSONImporter(s ImporterStorage, m Mapper) *JSONImporter {
	return &JSONImporter{
		Storage: s,
		Mapper:  m,
	}
}

func (i *JSONImporter) ImportApps(r io.Reader) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		// Allocate new element each time. Storage is working on references to objects.
		app := new(App)
		_ = json.Unmarshal(s.Bytes(), app)
		_ = i.Storage.AppSave(app)
	}
}

func (i *JSONImporter) ImportInstances(r io.Reader) ImportInstancesReport {
	rep := ImportInstancesReport{}

	s := bufio.NewScanner(r)
	for s.Scan() {
		insExp := new(InstanceExported)
		_ = json.Unmarshal(s.Bytes(), insExp)

		app, err := i.Storage.AppLoad(insExp.AppID)
		if err == ErrElementNotFound {
			rep.Failed++
			rep.Failures = append(rep.Failures, ImportInstanceFailureReason{
				ID: insExp.ID,
				AppID: insExp.AppID,
				Reason: FailUnknownApp,
			})
			continue
		}

		ins := new(Instance)
		*ins = insExp.Instance
		ins.App = app

		_ = i.Storage.InstanceSave(ins)

		_ = i.Mapper.Add(ins)

		rep.Succeeded++
	}

	return rep
}

type ImportFailureReason string

const (
	FailUnknownApp ImportFailureReason = "unknownApp"
)

type ImportInstanceFailureReason struct {
	ID     string
	AppID  string
	Reason ImportFailureReason
}

type ImportInstancesReport struct {
	Succeeded int
	Failed    int
	Failures  []ImportInstanceFailureReason
}
