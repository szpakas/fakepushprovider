package android

type App struct {
	ID       string
	SenderID string
	ApiKey   string
}

type RegistrationID string

type Instance struct {
	ID              string
	State           InstanceState
	App             *App
	RegistrationIDS []RegistrationID
	CanonicalID     RegistrationID
}

type InstanceState int

const (
	InstanceStateRegistered = iota + 1
	InstanceStateUnregistered
)
