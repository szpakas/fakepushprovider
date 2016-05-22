package android

import (
	"fmt"
	"testing"

	a "github.com/stretchr/testify/assert"
	ar "github.com/stretchr/testify/require"
)

func Test_Mapper_Factory(t *testing.T) {
	m, closer := tsMapperSetup()
	defer closer()

	a.IsType(t, &MemoryMapper{}, m, "Incorrect type")

	a.NotNil(t, m.regIDs, "regIDs: not initialised")
}

func Test_Mapper_Add_Success(t *testing.T) {
	tests := []int{1, 2, 4, 8, 10, 100}
	for tNo, n := range tests {
		m, closer := tsMapperSetup()

		ins := Instance{
			ID:              fmt.Sprintf("i-%d", tNo),
			App:             &tfAppA,
			RegistrationIDS: make([]RegistrationID, n),
		}

		for i := 0; i < n; i++ {
			ins.RegistrationIDS[i] = RegistrationID(fmt.Sprintf("RegID-%s-%d", ins.ID, i))
		}

		ar.NoError(t, m.Add(&ins), "Add: unexpected error, instanceID: %s", ins.ID)
		thAssertInstanceInMapper(t, m, &ins)
		closer()
	}
}

func Test_Mapper_Add_Error_AlreadyExists(t *testing.T) {
	m, closer := tsMapperSetup()
	defer closer()

	insAA := Instance{
		ID:              "i-AA",
		App:             &tfAppA,
		RegistrationIDS: []RegistrationID{"A", "B", "C"},
	}

	insAB := Instance{
		ID:              "i-AB",
		App:             &tfAppA,
		RegistrationIDS: []RegistrationID{"D", "E", "F"},
	}

	insB := Instance{
		ID:              "i-B",
		App:             &tfAppA,
		RegistrationIDS: []RegistrationID{"G", "A", "H", "I", "F"},
	}

	ar.NoError(t, m.Add(&insAA), "Add: unexpected error on 1st instance")
	ar.NoError(t, m.Add(&insAB), "Add: unexpected error on 2nd instance")

	err := m.Add(&insB)
	ar.Error(t, err, "Add: expected error on 3rd instance")
	errMapped, asrOk := err.(AlreadyMappedError)
	ar.True(t, asrOk, "Add: error asertion failed - unknown error type")

	cExp := map[RegistrationID]*Instance{
		"A": &insAA,
		"F": &insAB,
	}
	a.Equal(t, cExp, errMapped.Collisions, "Mismatch on collsions in error")
	thAssertInstanceInMapper(t, m, &insAA) // so that it's not removed/remapped
	thAssertInstanceInMapper(t, m, &insAB)

	for _, rID := range []RegistrationID{"G", "H", "I"} {
		a.NotContains(t, m.regIDs, rID, "registrationID %s from 3rd instance found in map", rID)
	}
}
