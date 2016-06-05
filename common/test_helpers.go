package common

import (
	"fmt"

	a "github.com/stretchr/testify/assert"
)

// THAssertReporter is a custom reporter for testify library.
// It allows extra symbol to be passed so table tests are more readable
type THAssertReporter struct {
	Backend *a.Assertions
	Symbol  string
}

// Errorf implements Reporter.Errorf.
func (r *THAssertReporter) Errorf(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	if r.Symbol != "" {
		msg = fmt.Sprintf("--- test case: %s ---\n%s", r.Symbol, msg)
	}
	r.Backend.Fail(msg)
}
