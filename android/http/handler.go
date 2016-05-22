package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/szpakas/fakepushprovider/android"
)

var _ = spew.Config

type Storage interface {
	AppFind(apiKey string) (*android.App, error)
	InstanceFind(appID string, registrationID android.RegistrationID) (*android.Instance, error)
}

type handler struct {
	storage Storage
}

func newHandler(storage Storage) *handler {
	return &handler{
		storage: storage,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqB := new(DownstreamMessage)
	_ = json.NewDecoder(r.Body).Decode(reqB)

	// Authorization: key=$api_key
	authorizationValue := r.Header.Get("Authorization")
	apiKey := strings.Split(authorizationValue, "=")[1]

	var regIDs []android.RegistrationID
	for _, regID := range reqB.RegistrationIDS {
		regIDs = append(regIDs, regID)
	}
	if reqB.To != "" {
		regIDs = append(regIDs, reqB.To)
	}

	app, _ := h.storage.AppFind(apiKey)

	resB := DownstreamResponse{
		MulticastID: 1234567890,
	}

	for _, regID := range regIDs {
		ins, err := h.storage.InstanceFind(app.ID, regID)
		if err != nil || ins.State == android.InstanceStateUnregistered {
			resB.Failure++
			resB.Results = append(resB.Results, MessageResult{Error: DeviceUnregistered})
			continue
		}

		resB.Success++
		resB.Results = append(resB.Results, MessageResult{MessageID: "m:1234"})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resB)
}