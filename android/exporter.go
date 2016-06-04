package android

import (
	"encoding/json"
	"io"
)

// omit is used to exclude fields during JSON marshalling.
type omit *struct{}

type InstanceExported struct {
	Instance

	AppID string
	App   omit `json:"App,omitempty"`
}

// JSONExporter exports data in JSON format.
// Implements android.Storer interface for generator.
type JSONExporter struct {
	AppWriter      io.Writer
	InstanceWriter io.Writer

	AppEncoder      *json.Encoder
	InstanceEncoder *json.Encoder
}

func NewJSONExporter(appWriter, instanceWriter io.Writer) *JSONExporter {
	return &JSONExporter{
		AppWriter:       appWriter,
		InstanceWriter:  instanceWriter,
		AppEncoder:      json.NewEncoder(appWriter),
		InstanceEncoder: json.NewEncoder(instanceWriter),
	}
}

func (e *JSONExporter) AppSave(o *App) error {
	return e.AppEncoder.Encode(o)
}

func (e *JSONExporter) InstanceSave(o *Instance) error {
	oW := InstanceExported{
		Instance: *o,
		AppID:    o.App.ID,
	}
	return e.InstanceEncoder.Encode(&oW)
}
