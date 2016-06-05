package fcm

import (
	"testing"

	a "github.com/stretchr/testify/assert"
)

func Test_Exporter_Factory(t *testing.T) {
	e, _, _, closer := tsExporterSetup()
	defer closer()

	a.IsType(t, &JSONExporter{}, e, "Incorrect type")

	a.NotNil(t, e.AppWriter, "AppWriter: not initialised")
	a.NotNil(t, e.InstanceWriter, "InstanceWriter: not initialised")
	a.NotNil(t, e.AppEncoder, "AppEncoder: not initialised")
	a.NotNil(t, e.InstanceEncoder, "InstanceEncoder: not initialised")
}

func Test_Exporter_AppSave(t *testing.T) {
	tests := []struct {
		app *App
		exp string
	}{
		{&TFAppA, tfAppAExportJSON},
		{&TFAppB, tfAppBExportJSON},
		{&TFAppC, tfAppCExportJSON},
	}
	for _, tc := range tests {
		e, aw, _, closer := tsExporterSetup()
		a.NoError(t, e.AppSave(tc.app), "%s: error on save", tc.app.ID)
		a.JSONEq(t, tc.exp, aw.String(), "%s: mismatch", tc.app.ID)
		closer()
	}
}

func Test_Exporter_InstanceSave(t *testing.T) {
	tests := []struct {
		ins *Instance
		exp string
	}{
		{&TFInsAA, tfInsAAExportJSON},
		{&TFInsAB, tfInsABExportJSON},
		{&TFInsAC, tfInsACExportJSON},
		{&TFInsAZ, tfInsAZExportJSON},
		{&TFInsBA, tfInsBAExportJSON},
		{&TFInsBB, tfInsBBExportJSON},
		{&TFInsBC, tfInsBCExportJSON},
	}
	for _, tc := range tests {
		e, _, iw, closer := tsExporterSetup()
		a.NoError(t, e.InstanceSave(tc.ins), "%s: error on save", tc.ins.ID)
		a.JSONEq(t, tc.exp, iw.String(), "%s: mismatch", tc.ins.ID)
		closer()
	}
}
