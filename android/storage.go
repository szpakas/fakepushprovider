package android

import "github.com/pkg/errors"

var (
	ErrElementNotFound = errors.New("Storage: element not found")
)

type MemoryStorage struct {
	apps      map[string]*App
	instances map[string]map[string]*Instance
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		apps:      make(map[string]*App),
		instances: make(map[string]map[string]*Instance),
	}
}

func (s *MemoryStorage) AppSave(o *App) error {
	s.apps[o.ID] = o
	return nil
}

func (s *MemoryStorage) AppLoad(id string) (*App, error) {
	o, found := s.apps[id]
	if !found {
		return nil, ErrElementNotFound
	}
	return o, nil
}

func (s *MemoryStorage) AppFind(apiKey string) (*App, error) {
	for appID, _ := range s.apps {
		if apiKey == s.apps[appID].ApiKey {
			return s.apps[appID], nil
		}
	}
	return nil, ErrElementNotFound
}

func (s *MemoryStorage) InstanceSave(o *Instance) error {
	_, found := s.instances[o.App.ID]
	if !found {
		s.instances[o.App.ID] = make(map[string]*Instance)
	}
	s.instances[o.App.ID][o.ID] = o
	return nil
}

func (s *MemoryStorage) InstanceLoad(aID, iID string) (*Instance, error) {
	if _, found := s.instances[aID]; !found {
		return nil, ErrElementNotFound
	}

	o, found := s.instances[aID][iID]
	if !found {
		return nil, ErrElementNotFound
	}
	return o, nil
}

func (s *MemoryStorage) InstanceFind(appID string, registrationID RegistrationID) (*Instance, error) {
	for insID, _ := range s.instances[appID] {
		if registrationID == s.instances[appID][insID].CanonicalID {
			return s.instances[appID][insID], nil
		}
		for _, regID := range s.instances[appID][insID].RegistrationIDS {
			if registrationID == regID {
				return s.instances[appID][insID], nil
			}
		}
	}
	return nil, ErrElementNotFound
}
