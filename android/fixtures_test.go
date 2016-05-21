package android

var tfAppA = App{
	ID: "appIdA",
	SenderID: "senderIdA",
	ApiKey: "apiKeyA",
}

var tfAppB = App{
	ID: "appIdB",
	SenderID: "senderIdB",
	ApiKey: "apiKeyB",
}
var tfAppC = App{
	ID: "appIdC",
	SenderID: "senderIdC",
	ApiKey: "apiKeyC",
}

var tfInsAA = Instance{
	ID: "instanceAA",
	State: InstanceStateRegistered,
	App: &tfAppA,
	RegistrationIDS: []RegistrationID{"RegIdAA1"},
}

var tfInsAB = Instance{
	ID: "instanceAB",
	State: InstanceStateRegistered,
	App: &tfAppA,
	RegistrationIDS: []RegistrationID{"RegIdAB1", "RegIdAB2"},
	CanonicalID: "RegIdAB1",
}

var tfInsAC = Instance{
	ID: "instanceAC",
	State: InstanceStateRegistered,
	App: &tfAppA,
	RegistrationIDS: []RegistrationID{"RegIdAC1", "RegIdAC2", "RegIdAC3"},
	CanonicalID: "RegIdAC3",
}

var tfInsAZ = Instance{
	ID: "instanceAZ",
	State: InstanceStateUnregistered,
	App: &tfAppA,
	RegistrationIDS: []RegistrationID{"RegIdAZ1", "RegIdAZ2"},
}

var tfInsBA = Instance{
	ID: "instanceBA",
	State: InstanceStateRegistered,
	App: &tfAppB,
	RegistrationIDS: []RegistrationID{"RegIdBA1"},
}

var tfInsBB = Instance{
	ID: "instanceBB",
	State: InstanceStateRegistered,
	App: &tfAppB,
	RegistrationIDS: []RegistrationID{"RegIdBB1", "RegIdBB2"},
	CanonicalID: "RegIdBB1",
}

var tfInsBC = Instance{
	ID: "instanceBC",
	State: InstanceStateRegistered,
	App: &tfAppB,
	RegistrationIDS: []RegistrationID{"RegIdBC1", "RegIdBC2", "RegIdBC3"},
	CanonicalID: "RegIdBC3",
}
