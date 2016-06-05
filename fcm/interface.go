package fcm

type Storer interface {
	AppSave(o *App) error
	InstanceSave(o *Instance) error
}

type Mapper interface {
	Add(i *Instance) error
}
