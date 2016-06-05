package apns

import "time"

var TFAppA = App{
	ID:       "appIdA",
	BundleID: "bundleIdA",
}

var TFAppB = App{
	ID:       "appIdB",
	BundleID: "bundleIdB",
}

var TFAppC = App{
	ID:       "appIdC",
	BundleID: "bundleIdC",
}

var TFInsAA = Instance{
	ID:    "instanceAA",
	State: InstanceStateRegistered,
	App:   &TFAppA,
	Token: "TokenAA",
}

var TFInsAB = Instance{
	ID:    "instanceAB",
	State: InstanceStateRegistered,
	App:   &TFAppA,
	Token: "TokenAB",
}

var TFInsAC = Instance{
	ID:    "instanceAC",
	State: InstanceStateRegistered,
	App:   &TFAppA,
	Token: "TokenAC",
}

var TFInsAD = Instance{
	ID:    "instanceAD",
	State: InstanceStateRegistered,
	App:   &TFAppA,
	Token: "TokenAD",
}

var TFInsAZ = Instance{
	ID:       "instanceAZ",
	State:    InstanceStateUnregistered,
	App:      &TFAppA,
	Token:    "TokenAZ",
	LastSeen: time.Date(2016, 1, 2, 20, 27, 33, 456, time.UTC).Unix(),
}

var TFInsBA = Instance{
	ID:    "instanceBA",
	State: InstanceStateRegistered,
	App:   &TFAppB,
	Token: "TokenBA",
}

var TFInsBB = Instance{
	ID:    "instanceBB",
	State: InstanceStateRegistered,
	App:   &TFAppB,
	Token: "TokenBB",
}

var TFInsBC = Instance{
	ID:    "instanceBC",
	State: InstanceStateRegistered,
	App:   &TFAppB,
	Token: "TokenBC",
}

var tfAppAExportJSON = `{"ID":"appIdA","BundleID":"bundleIdA"}`
var tfAppBExportJSON = `{"ID":"appIdB","BundleID":"bundleIdB"}`
var tfAppCExportJSON = `{"ID":"appIdC","BundleID":"bundleIdC"}`

var tfInsAAExportJSON = `{"ID":"instanceAA","State":1,"Token":"TokenAA","AppID":"appIdA"}`
var tfInsABExportJSON = `{"ID":"instanceAB","State":1,"Token":"TokenAB","AppID":"appIdA"}`
var tfInsACExportJSON = `{"ID":"instanceAC","State":1,"Token":"TokenAC","AppID":"appIdA"}`
var tfInsADExportJSON = `{"ID":"instanceAD","State":1,"Token":"TokenAD","AppID":"appIdA"}`
var tfInsAZExportJSON = `{"ID":"instanceAZ","State":2,"Token":"TokenAZ","AppID":"appIdA","LastSeen":1451766453}`
var tfInsBAExportJSON = `{"ID":"instanceBA","State":1,"Token":"TokenBA","AppID":"appIdB"}`
var tfInsBBExportJSON = `{"ID":"instanceBB","State":1,"Token":"TokenBB","AppID":"appIdB"}`
var tfInsBCExportJSON = `{"ID":"instanceBC","State":1,"Token":"TokenBC","AppID":"appIdB"}`
