package android

import (
	"testing"

	a "github.com/stretchr/testify/assert"
	ar "github.com/stretchr/testify/require"
)

func Test_MemoryStorage_Factory(t *testing.T) {
	s, closer := tsMemoryStorageSetup(); defer closer()

	a.IsType(t, &MemoryStorage{}, s, "Incorrect type")

	a.NotNil(t, s.apps, "apps: not initialised")
	a.NotNil(t, s.instances, "instances: not initialised")
}

func Test_MemoryStorage_AppSave_Success(t *testing.T) {
	s, closer := tsMemoryStorageSetup(); defer closer()

	ar.NoError(t, s.AppSave(&tfAppA), "AppSave: unexpected error")
	ar.Contains(t, s.apps, tfAppA.ID, "App ID not in storage")
	a.Equal(t, s.apps[tfAppA.ID], &tfAppA, "App from storage does not match")
}

func Test_MemoryStorage_AppLoad_Exists(t *testing.T) {
	s, closer := tsMemoryStorageSetup(); defer closer()

	s.apps[tfAppA.ID] = &tfAppA
	o, err := s.AppLoad(tfAppA.ID)
	ar.NoError(t, err, "AppLoad: unexpected error")
	a.Equal(t, &tfAppA, o, "App from storage does not match")
}

func Test_MemoryStorage_AppLoad_NotFound(t *testing.T) {
	s, closer := tsMemoryStorageSetup(); defer closer()

	_, err := s.AppLoad(tfAppA.ID)
	ar.Equal(t, ErrElementNotFound, err, "AppLoad: expected NotFound error")
}

func Test_MemoryStorage_InstanceSave_Single_Success(t *testing.T) {
	s, closer := tsMemoryStorageWitAppsSetup(); defer closer()

	ar.NoError(t, s.InstanceSave(&tfInsAA), "InstanceSave: unexpected error")
	ar.Contains(t, s.instances, tfInsAA.App.ID, "storage on app level not initiated")
	ar.Contains(t, s.instances[tfInsAA.App.ID], tfInsAA.ID, "Instance ID not in storage")
	a.Equal(t, s.instances[tfInsAA.App.ID][tfInsAA.ID], &tfInsAA, "Instance from storage does not match")
}

func Test_MemoryStorage_InstanceSave_Multiple_Success(t *testing.T) {
	s, closer := tsMemoryStorageWitAppsSetup(); defer closer()

	ar.NoError(t, s.InstanceSave(&tfInsAA), "InstanceSave %s: unexpected error", "AA")
	ar.NoError(t, s.InstanceSave(&tfInsAB), "InstanceSave %s: unexpected error", "AB")
	ar.NoError(t, s.InstanceSave(&tfInsAC), "InstanceSave %s: unexpected error", "AC")
	ar.NoError(t, s.InstanceSave(&tfInsBA), "InstanceSave %s: unexpected error", "BA")
	ar.NoError(t, s.InstanceSave(&tfInsBB), "InstanceSave %s: unexpected error", "BB")
	ar.NoError(t, s.InstanceSave(&tfInsBC), "InstanceSave %s: unexpected error", "BC")

	a.Equal(t, s.instances[tfInsAA.App.ID][tfInsAA.ID], &tfInsAA, "Instance %s: from storage does not match", "AA")
	a.Equal(t, s.instances[tfInsAB.App.ID][tfInsAB.ID], &tfInsAB, "Instance %s: from storage does not match", "AB")
	a.Equal(t, s.instances[tfInsAC.App.ID][tfInsAC.ID], &tfInsAC, "Instance %s: from storage does not match", "AC")
}

func Test_MemoryStorage_InstanceLoad_Exists(t *testing.T) {
	s, closer := tsMemoryStorageSetup(); defer closer()

	s.instances[tfInsAA.App.ID] = make(map[string]*Instance)
	s.instances[tfInsAA.App.ID][tfInsAA.ID] = &tfInsAA
	o, err := s.InstanceLoad(tfInsAA.App.ID, tfInsAA.ID)
	ar.NoError(t, err, "InstanceLoad: unexpected error")
	a.Equal(t, &tfInsAA, o, "Instance from storage does not match")
}

func Test_MemoryStorage_InstanceLoad_NotFound_AppExists(t *testing.T) {
	s, closer := tsMemoryStorageSetup(); defer closer()

	s.instances[tfInsAA.App.ID] = make(map[string]*Instance)
	_, err := s.InstanceLoad(tfInsAA.App.ID, tfInsAA.ID)
	a.Equal(t, ErrElementNotFound, err, "InstanceLoad: expected NotFound error")
}

func Test_MemoryStorage_InstanceLoad_NotFound_AppDoesNotExists(t *testing.T) {
	s, closer := tsMemoryStorageSetup(); defer closer()

	_, err := s.InstanceLoad(tfInsAA.App.ID, tfInsAA.ID)
	a.Equal(t, ErrElementNotFound, err, "InstanceLoad: expected NotFound error")
}
