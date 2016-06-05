package fcm

import (
	"testing"

	a "github.com/stretchr/testify/assert"
	ar "github.com/stretchr/testify/require"
)

func Test_MemoryStorage_Factory(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	a.IsType(t, &MemoryStorage{}, s, "Incorrect type")

	a.NotNil(t, s.apps, "apps: not initialised")
	a.NotNil(t, s.instances, "instances: not initialised")
}

func Test_MemoryStorage_AppSave_Success(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	ar.NoError(t, s.AppSave(&TFAppA), "AppSave: unexpected error")
	ar.Contains(t, s.apps, TFAppA.ID, "App ID not in storage")
	a.Equal(t, s.apps[TFAppA.ID], &TFAppA, "App from storage does not match")
}

func Test_MemoryStorage_AppLoad_Exists(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	s.apps[TFAppA.ID] = &TFAppA
	o, err := s.AppLoad(TFAppA.ID)
	ar.NoError(t, err, "AppLoad: unexpected error")
	a.Equal(t, &TFAppA, o, "App from storage does not match")
}

func Test_MemoryStorage_AppLoad_NotFound(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	_, err := s.AppLoad(TFAppA.ID)
	ar.Equal(t, ErrElementNotFound, err, "AppLoad: expected NotFound error")
}

func Test_MemoryStorage_AppFind_Exists(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	s.apps[TFAppA.ID] = &TFAppA
	o, err := s.AppFind(TFAppA.ApiKey)
	ar.NoError(t, err, "AppFind: unexpected error")
	a.Equal(t, &TFAppA, o, "App from storage does not match")
}

func Test_MemoryStorage_AppFind_NotFound(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	_, err := s.AppFind(TFAppA.ApiKey)
	ar.Equal(t, ErrElementNotFound, err, "AppFind: expected NotFound error")
}

func Test_MemoryStorage_InstanceSave_Single_Success(t *testing.T) {
	s, closer := tsMemoryStorageWitAppsSetup()
	defer closer()

	ar.NoError(t, s.InstanceSave(&TFInsAA), "InstanceSave: unexpected error")
	ar.Contains(t, s.instances, TFInsAA.App.ID, "storage on app level not initiated")
	ar.Contains(t, s.instances[TFInsAA.App.ID], TFInsAA.ID, "Instance ID not in storage")
	a.Equal(t, s.instances[TFInsAA.App.ID][TFInsAA.ID], &TFInsAA, "Instance from storage does not match")
}

func Test_MemoryStorage_InstanceSave_Multiple_Success(t *testing.T) {
	s, closer := tsMemoryStorageWitAppsSetup()
	defer closer()

	ar.NoError(t, s.InstanceSave(&TFInsAA), "InstanceSave %s: unexpected error", "AA")
	ar.NoError(t, s.InstanceSave(&TFInsAB), "InstanceSave %s: unexpected error", "AB")
	ar.NoError(t, s.InstanceSave(&TFInsAC), "InstanceSave %s: unexpected error", "AC")
	ar.NoError(t, s.InstanceSave(&TFInsBA), "InstanceSave %s: unexpected error", "BA")
	ar.NoError(t, s.InstanceSave(&TFInsBB), "InstanceSave %s: unexpected error", "BB")
	ar.NoError(t, s.InstanceSave(&TFInsBC), "InstanceSave %s: unexpected error", "BC")

	a.Equal(t, s.instances[TFInsAA.App.ID][TFInsAA.ID], &TFInsAA, "Instance %s: from storage does not match", "AA")
	a.Equal(t, s.instances[TFInsAB.App.ID][TFInsAB.ID], &TFInsAB, "Instance %s: from storage does not match", "AB")
	a.Equal(t, s.instances[TFInsAC.App.ID][TFInsAC.ID], &TFInsAC, "Instance %s: from storage does not match", "AC")
}

func Test_MemoryStorage_InstanceLoad_Exists(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	s.instances[TFInsAA.App.ID] = make(map[string]*Instance)
	s.instances[TFInsAA.App.ID][TFInsAA.ID] = &TFInsAA
	o, err := s.InstanceLoad(TFInsAA.App.ID, TFInsAA.ID)
	ar.NoError(t, err, "InstanceLoad: unexpected error")
	a.Equal(t, &TFInsAA, o, "Instance from storage does not match")
}

func Test_MemoryStorage_InstanceLoad_NotFound_AppExists(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	s.instances[TFInsAA.App.ID] = make(map[string]*Instance)
	_, err := s.InstanceLoad(TFInsAA.App.ID, TFInsAA.ID)
	a.Equal(t, ErrElementNotFound, err, "InstanceLoad: expected NotFound error")
}

func Test_MemoryStorage_InstanceLoad_NotFound_AppDoesNotExists(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	_, err := s.InstanceLoad(TFInsAA.App.ID, TFInsAA.ID)
	a.Equal(t, ErrElementNotFound, err, "InstanceLoad: expected NotFound error")
}

func Test_MemoryStorage_InstanceFind_ByCanonicalID_Exists(t *testing.T) {
	s, closer := TSMemoryStorageWitAppsAndInstancesSetup()
	defer closer()

	o, err := s.InstanceFind(TFInsAA.App.ID, TFInsAA.CanonicalID)
	ar.NoError(t, err, "InstanceFind: unexpected error")
	a.Equal(t, &TFInsAA, o, "Instance from storage does not match")
}

func Test_MemoryStorage_InstanceFind_ByNotCanonicalID_Exists(t *testing.T) {
	s, closer := TSMemoryStorageWitAppsAndInstancesSetup()
	defer closer()

	o, err := s.InstanceFind(TFInsAC.App.ID, TFInsAC.RegistrationIDS[1])
	ar.NotEqual(t, TFInsAC.CanonicalID, TFInsAC.RegistrationIDS[1], "selected ID is a canonical one")
	ar.NoError(t, err, "InstanceFind: unexpected error")
	a.Equal(t, &TFInsAC, o, "Instance from storage does not match")
}

func Test_MemoryStorage_InstanceFind_Error_NotFound_AppExists(t *testing.T) {
	s, closer := TSMemoryStorageWitAppsAndInstancesSetup()
	defer closer()

	_, err := s.InstanceFind(TFInsAA.App.ID, "FakeRegID")
	ar.Equal(t, ErrElementNotFound, err, "InstanceFind: expected error")
}

func Test_MemoryStorage_InstanceFind_Error_NotFound_AppDoesNotExists(t *testing.T) {
	s, closer := TSMemoryStorageWitAppsAndInstancesSetup()
	defer closer()

	_, err := s.InstanceFind("FakeAppID", "FakeRegID")
	ar.Equal(t, ErrElementNotFound, err, "InstanceFind: expected error")
}

//func Test_MemoryStorage_AppFind_NotFound(t *testing.T) {
//	s, closer := tsMemoryStorageSetup(); defer closer()
//
//	_, err := s.AppFind(tfAppA.ID)
//	ar.Equal(t, ErrElementNotFound, err, "AppFind: expected NotFound error")
//}
