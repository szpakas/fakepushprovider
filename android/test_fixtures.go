package android

var TFAppA = App{
	ID:       "appIdA",
	SenderID: "senderIdA",
	ApiKey:   "apiKeyA",
}

var TFAppB = App{
	ID:       "appIdB",
	SenderID: "senderIdB",
	ApiKey:   "apiKeyB",
}
var TFAppC = App{
	ID:       "appIdC",
	SenderID: "senderIdC",
	ApiKey:   "apiKeyC",
}

var TFInsAA = Instance{
	ID:              "instanceAA",
	State:           InstanceStateRegistered,
	App:             &TFAppA,
	RegistrationIDS: []RegistrationID{"RegIdAA1"},
	CanonicalID:     "RegIdAA1",
}

var TFInsAB = Instance{
	ID:              "instanceAB",
	State:           InstanceStateRegistered,
	App:             &TFAppA,
	RegistrationIDS: []RegistrationID{"RegIdAB1", "RegIdAB2"},
	CanonicalID:     "RegIdAB1",
}

var TFInsAC = Instance{
	ID:              "instanceAC",
	State:           InstanceStateRegistered,
	App:             &TFAppA,
	RegistrationIDS: []RegistrationID{"RegIdAC1", "RegIdAC2", "RegIdAC3"},
	CanonicalID:     "RegIdAC3",
}

var TFInsAZ = Instance{
	ID:              "instanceAZ",
	State:           InstanceStateUnregistered,
	App:             &TFAppA,
	RegistrationIDS: []RegistrationID{"RegIdAZ1", "RegIdAZ2"},
}

var TFInsBA = Instance{
	ID:              "instanceBA",
	State:           InstanceStateRegistered,
	App:             &TFAppB,
	RegistrationIDS: []RegistrationID{"RegIdBA1"},
	CanonicalID:     "RegIdBA1",
}

var TFInsBB = Instance{
	ID:              "instanceBB",
	State:           InstanceStateRegistered,
	App:             &TFAppB,
	RegistrationIDS: []RegistrationID{"RegIdBB1", "RegIdBB2"},
	CanonicalID:     "RegIdBB1",
}

var TFInsBC = Instance{
	ID:              "instanceBC",
	State:           InstanceStateRegistered,
	App:             &TFAppB,
	RegistrationIDS: []RegistrationID{"RegIdBC1", "RegIdBC2", "RegIdBC3"},
	CanonicalID:     "RegIdBC3",
}

var tfAppAExportJSON = `{"ID":"appIdA","SenderID":"senderIdA","ApiKey":"apiKeyA"}`
var tfAppBExportJSON = `{"ID":"appIdB","SenderID":"senderIdB","ApiKey":"apiKeyB"}`
var tfAppCExportJSON = `{"ID":"appIdC","SenderID":"senderIdC","ApiKey":"apiKeyC"}`

var tfInsAAExportJSON = `{"ID":"instanceAA","State":1,"RegistrationIDS":["RegIdAA1"],"CanonicalID":"RegIdAA1","AppID":"appIdA"}`
var tfInsABExportJSON = `{"ID":"instanceAB","State":1,"RegistrationIDS":["RegIdAB1","RegIdAB2"],"CanonicalID":"RegIdAB1","AppID":"appIdA"}`
var tfInsACExportJSON = `{"ID":"instanceAC","State":1,"RegistrationIDS":["RegIdAC1","RegIdAC2","RegIdAC3"],"CanonicalID":"RegIdAC3","AppID":"appIdA"}`
var tfInsAZExportJSON = `{"ID":"instanceAZ","State":2,"RegistrationIDS":["RegIdAZ1","RegIdAZ2"],"AppID":"appIdA"}`
var tfInsBAExportJSON = `{"ID":"instanceBA","State":1,"RegistrationIDS":["RegIdBA1"],"CanonicalID":"RegIdBA1","AppID":"appIdB"}`
var tfInsBBExportJSON = `{"ID":"instanceBB","State":1,"RegistrationIDS":["RegIdBB1","RegIdBB2"],"CanonicalID":"RegIdBB1","AppID":"appIdB"}`
var tfInsBCExportJSON = `{"ID":"instanceBC","State":1,"RegistrationIDS":["RegIdBC1","RegIdBC2","RegIdBC3"],"CanonicalID":"RegIdBC3","AppID":"appIdB"}`
