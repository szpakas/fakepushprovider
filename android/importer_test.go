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

func Test_Importer_ImportInstances_Success(t *testing.T) {
	i, s, m, closer := tsImporterSetup()
	defer closer()

	// GIVEN: apps are imported
	apps := []App{TFAppA, TFAppB, TFAppC}
	for i, _ := range apps {
		s.AppSave(&apps[i])
	}
	instancesToImport := []string{
		tfInsAAExportJSON, tfInsABExportJSON, tfInsACExportJSON, tfInsAZExportJSON,
		tfInsBAExportJSON, tfInsBBExportJSON, tfInsBCExportJSON,
	}
	repGot := i.ImportInstances(strings.NewReader(strings.Join(instancesToImport, "\n")))

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

	// AND: only those instances are persisted
	a.Equal(t, len(instances), s.InstancesTotal())

	// AND: mappings are created
	for i, _ := range instances {
		thAssertInstanceInMapper(t, m, &instances[i])
	}

	// AND: report matches expectations
	repExp := ImportInstancesReport{
		Succeeded: len(instancesToImport),
		Failed:    0,
	}
	a.Equal(t, repExp, repGot)
}

func Test_Importer_ImportInstances_Error_UnknownApp(t *testing.T) {
	i, s, m, closer := tsImporterSetup()
	defer closer()

	// GIVEN: apps are imported but TFAppA is NOT imported
	apps := []App{TFAppB, TFAppC}
	for i, _ := range apps {
		s.AppSave(&apps[i])
	}

	// WHEN: instances are imported
	instancesToFail := []string{tfInsAAExportJSON, tfInsABExportJSON}
	instancesToSucceed := []string{tfInsBAExportJSON, tfInsBBExportJSON, tfInsBCExportJSON}

	var instancesToImport []string
	instancesToImport = append(instancesToImport, instancesToSucceed...)
	instancesToImport = append(instancesToImport, instancesToFail...)

	repGot := i.ImportInstances(strings.NewReader(strings.Join(instancesToImport, "\n")))

	// THEN: instances are persisted in storage
	instances := []Instance{
		TFInsBA, TFInsBB, TFInsBC,
	}

	for i, _ := range instances {
		iExp := instances[i]
		iGot, err := s.InstanceLoad(iExp.App.ID, iExp.ID)
		if a.NoError(t, err, "%s: error on load", iExp.ID) {
			a.Equal(t, iExp, *iGot, "%s: mismatch on loaded object", iExp.ID)
		}
	}

	// AND: only those instances are persisted
	a.Equal(t, len(instances), s.InstancesTotal())

	// AND: mappings are created
	for i, _ := range instances {
		thAssertInstanceInMapper(t, m, &instances[i])
	}

	// AND: report matches expectations
	repExp := ImportInstancesReport{
		Succeeded: 3,
		Failed:    2,
		Failures: []ImportInstanceFailureReason{
			{ID: TFInsAA.ID, AppID: TFInsAA.App.ID, Reason: FailUnknownApp},
			{ID: TFInsAB.ID, AppID: TFInsAB.App.ID, Reason: FailUnknownApp},
		},
	}
	a.Equal(t, repExp, repGot)
}
