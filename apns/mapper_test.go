package apns

import (
	"testing"

	a "github.com/stretchr/testify/assert"
	ar "github.com/stretchr/testify/require"
)

func Test_Mapper_Factory(t *testing.T) {
	m, closer := tsMapperSetup()
	defer closer()

	a.IsType(t, &MemoryMapper{}, m, "Incorrect type")

	a.NotNil(t, m.tokens, "regIDs: not initialised")
}

func Test_Mapper_Add_Success(t *testing.T) {
	m, closer := tsMapperSetup()

	ar.NoError(t, m.Add(&TFInsAA), "Add: unexpected error, instanceID: %s", TFInsAA.ID)
	thAssertInstanceInMapper(t, m, &TFInsAA)
	closer()
}

func Test_Mapper_Add_Error_AlreadyExists(t *testing.T) {
	m, closer := tsMapperSetup()
	defer closer()

	ar.NoError(t, m.Add(&TFInsAA), "Add: unexpected error on 1st instance")
	ar.NoError(t, m.Add(&TFInsAB), "Add: unexpected error on 2nd instance")

	err := m.Add(&TFInsAA)
	ar.Error(t, err, "Add: expected error on 3rd instance")
	errMapped, asrOk := err.(AlreadyMappedError)
	ar.True(t, asrOk, "Add: error asertion failed - unknown error type")

	cExp := map[Token]*Instance{
		TFInsAA.Token: &TFInsAA,
	}
	a.Equal(t, cExp, errMapped.Collisions, "Mismatch on collsions in error")
	thAssertInstanceInMapper(t, m, &TFInsAA) // so that it's not removed/remapped
	thAssertInstanceInMapper(t, m, &TFInsAB)
}
