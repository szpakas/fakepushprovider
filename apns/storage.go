package apns

import (
	"errors"
)

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
		if apiKey == s.apps[appID].BundleID {
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

func (s *MemoryStorage) InstanceFind(appID string, token Token) (*Instance, error) {
	for insID, _ := range s.instances[appID] {
		if token == s.instances[appID][insID].Token {
			return s.instances[appID][insID], nil
		}
	}
	return nil, ErrElementNotFound
}

func (s *MemoryStorage) InstancesTotal() int {
	var out int
	for appID, _ := range s.instances {
		out += len(s.instances[appID])
	}
	return out
}

//// Report is producing report on the state of the storage.
//// It's useful as debugging and monitoring tool.
//func (s *MemoryStorage) Report() map[string]interface{} {
//	out := make(map[string]interface{})
//	out["apps:total"] = len(s.apps)
//	for appID, app := range s.apps {
//		out[fmt.Sprintf("apps:id=%s:id", appID)] = app.ID
//		out[fmt.Sprintf("apps:id=%s:apiKey", appID)] = app.ApiKey
//		out[fmt.Sprintf("apps:id=%s:senderId", appID)] = app.SenderID
//	}
//
//	out["instances:total:count"] = s.InstancesTotal()
//	for appID, _ := range s.instances {
//		out[fmt.Sprintf("apps:id=%s:instances:total", appID)] = len(s.instances[appID])
//	}
//
//	return out
//}
