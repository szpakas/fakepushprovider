package apns

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	a "github.com/stretchr/testify/assert"
	ar "github.com/stretchr/testify/require"
)

var _ = spew.Config

func Test_Handler_Factory(t *testing.T) {
	st, h, _, _, closer := tsServerSetup(t, "")
	defer closer()

	a.IsType(t, &Handler{}, h, "Incorrect type")
	ar.NotNil(t, h, "empty handler")
	a.Equal(t, st, h.Storage, "storage not assigned")
}
