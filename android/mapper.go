package android

type AlreadyMappedError struct {
	Collisions map[RegistrationID]*Instance
}

func (AlreadyMappedError) Error() string {
	return "Mapper: element already mapped"
}

type Mapper struct {
	regIDs map[RegistrationID]*Instance
}

func NewMapper() *Mapper {
	return &Mapper{
		regIDs: make(map[RegistrationID]*Instance),
	}
}

func (m *Mapper) Add(i *Instance) error {
	collisions := make(map[RegistrationID]*Instance, len(i.RegistrationIDS))
	for _, r := range i.RegistrationIDS {
		if oI, found := m.regIDs[r]; found {
			collisions[r] = oI
		}
	}
	if len(collisions) > 0 {
		return AlreadyMappedError{collisions}
	}

	for _, r := range i.RegistrationIDS {
		m.regIDs[r] = i
	}
	return nil
}
