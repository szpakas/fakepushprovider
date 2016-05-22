package android

import (
	"testing"

	a "github.com/stretchr/testify/assert"
)

func Test_Generator_Factory(t *testing.T) {
	m, closer := tsGeneratorSetup(3, 4, 5)
	defer closer()

	a.IsType(t, &Generator{}, m, "Incorrect type")

	a.Equal(t, 3, m.AppTotal, "AppTotal: not initialised")
	a.Equal(t, 4, m.InstancesPerApp, "InstancesPerApp: not initialised")
	a.Equal(t, 5, m.RegistrationIDPerInstanceMax, "RegistrationIDPerInstanceMax: not initialised")
}

func Test_Generator_Generate(t *testing.T) {
	maxRegIDs := 5
	g, gCloser := tsGeneratorSetup(3, 4, maxRegIDs)
	defer gCloser()
	s, sCloser := tsMemoryStorageSetup()
	defer sCloser()

	g.Generate(s)

	a.Len(t, s.apps, 3, "mismatch on number of apps")
	for aID, _ := range s.instances {
		if a.Contains(t, s.instances, aID, "no instances for app %s", aID) {
			a.Len(t, s.instances[aID], 4, "mismatch on number of instances for app %s", aID)

			for iID, iObj := range s.instances[aID] {
				if len(iObj.RegistrationIDS) == 0 {
					a.Fail(t, "no registation IDs on instance %s for app %s", iID, aID)
				}
				if len(iObj.RegistrationIDS) > maxRegIDs {
					a.Fail(t, "number of registation IDs on instance %s for app %s is too high", iID, aID)
				}
				a.NotEmpty(t, iObj.CanonicalID, "CanonicalID empty for instance %s for app %s", iID, aID)
				a.Contains(t, iObj.RegistrationIDS, iObj.CanonicalID, "CanonicalID not in RegistrationIDs for instance %s for app %s", iID, aID)
			}
		}
	}
}
