package apns

type MemoryMapper struct {
	tokens map[Token]*Instance
}

func NewMemoryMapper() *MemoryMapper {
	return &MemoryMapper{
		tokens: make(map[Token]*Instance),
	}
}

func (m *MemoryMapper) Add(i *Instance) error {
	if iFound, found := m.tokens[i.Token]; found {
		return AlreadyMappedError{map[Token]*Instance{i.Token: iFound}}
	}

	m.tokens[i.Token] = i
	return nil
}

type AlreadyMappedError struct {
	Collisions map[Token]*Instance
}

func (AlreadyMappedError) Error() string {
	return "Mapper: element already mapped"
}
