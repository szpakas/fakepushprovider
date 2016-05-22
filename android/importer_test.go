package android

import (
	"strings"
	"testing"

	a "github.com/stretchr/testify/assert"
)

func Test_Importer_Factory(t *testing.T) {
	i, _, _, closer := tsImporterSetup()
	defer closer()

	a.IsType(t, &JSONImporter{}, i, "Incorrect type")

	a.NotNil(t, i.Storage, "Storage: not initialised")
	a.NotNil(t, i.Mapper, "Mapper: not initialised")
}

func Test_Importer_ImportApps(t *testing.T) {
	i, s, _, closer := tsImporterSetup()
	defer closer()

	i.ImportApps(strings.NewReader(strings.Join([]string{tfAppAExportJSON, tfAppBExportJSON, tfAppCExportJSON}, "\n")))

	for _, aExp := range []App{TFAppA, TFAppB, TFAppC} {
		aGot, err := s.AppLoad(aExp.ID)
		if a.NoError(t, err, "%s: error on load", aExp.ID) {
			a.Equal(t, aExp, *aGot, "%s: mismatch on loaded object", aExp.ID)
		}
	}
}

func Test_Importer_ImportInstances(t *testing.T) {
	i, s, m, closer := tsImporterSetup()
	defer closer()

	// GIVEN: apps are imported
	apps := []App{TFAppA, TFAppB, TFAppC}
	for i, _ := range apps {
		s.AppSave(&apps[i])
	}

	i.ImportInstances(strings.NewReader(strings.Join([]string{
		tfInsAAExportJSON, tfInsABExportJSON, tfInsACExportJSON, tfInsAZExportJSON,
		tfInsBAExportJSON, tfInsBBExportJSON, tfInsBCExportJSON,
	}, "\n")))

	// THEN: instances are persisted in storage
	instances := []Instance{
		TFInsAA, TFInsAB, TFInsAC, TFInsAZ,
		TFInsBA, TFInsBB, TFInsBC,
	}

	for i, _ := range instances {
		iExp := instances[i]
		iGot, err := s.InstanceLoad(iExp.App.ID, iExp.ID)
		if a.NoError(t, err, "%s: error on load", iExp.ID) {
			a.Equal(t, iExp, *iGot, "%s: mismatch on loaded object", iExp.ID)
		}
	}

	// AND: mappings are created
	for i, _ := range instances {
		thAssertInstanceInMapper(t, m, &instances[i])
	}
}
