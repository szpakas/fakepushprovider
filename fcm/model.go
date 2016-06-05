package fcm

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

	// CanonicalID is the primary RegistrationID associated with instance at the given point in time.
	// It's defined only for registered instances.
	CanonicalID RegistrationID `json:",omitempty"`
}

type InstanceState int

const (
	InstanceStateRegistered = iota + 1
	InstanceStateUnregistered
)
