package apns

import (
	"testing"

	a "github.com/stretchr/testify/assert"
)

func Test_Generator_Factory(t *testing.T) {
	m, closer := tsGeneratorSetup(3, 150, 10)
	defer closer()

	a.IsType(t, &Generator{}, m, "Incorrect type")

	a.Equal(t, 3, m.AppTotal, "AppTotal: not initialised")
	a.Equal(t, 150, m.InstancesPerApp, "InstancesPerApp: not initialised")
	a.Equal(t, 10.0, m.UnregisteredPercent, "UnregisteredPercent: not initialised")
}

func Test_Generator_Generate(t *testing.T) {
	instPerApp := 150
	var unregisteredPct float64 = 10
	g, gCloser := tsGeneratorSetup(3, instPerApp, unregisteredPct)
	defer gCloser()
	s, sCloser := tsMemoryStorageSetup()
	defer sCloser()

	g.Generate(s)

	a.Len(t, s.apps, 3, "mismatch on number of apps")

	// -- apps
	for aID, app := range s.apps {
		a.NotEmpty(t, app.ID, "ID empty for app %s", aID)
		a.NotEmpty(t, app.BundleID, "BundleID empty for app %s", aID)
		// TODO(szpakas): add cert test
	}

	// -- instances
	for aID, _ := range s.instances {
		if a.Contains(t, s.instances, aID, "no instances for app %s", aID) {
			a.Len(t, s.instances[aID], instPerApp, "mismatch on number of instances for app %s", aID)

			totalCnt := 0
			unregisteredCnt := 0
			for iID, iObj := range s.instances[aID] {
				totalCnt++
				if iObj.State == InstanceStateUnregistered {
					unregisteredCnt++
					a.NotZero(t, iObj.LastSeen, "LastSeen empty for unregistered instance %s for app %s", iID, aID)
				}
				a.NotEmpty(t, iObj.Token, "Token empty for instance %s for app %s", iID, aID)
			}
			// allow difference of 1 instance to cover round-up errors
			a.InDelta(t, (float64(totalCnt) * (unregisteredPct / 100.0)), float64(unregisteredCnt), 1)
		}
	}
}
