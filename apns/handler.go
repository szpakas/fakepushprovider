package apns

import (
	"encoding/json"
	"net/http"
	"regexp"
)

const (
	headerAPNSID    = "apns-id"
	headerAPNSTopic = "apns-topic"
)

var rPath = regexp.MustCompile(`^/3/device/(.+)$`)

type Storage interface {
	AppFind(bundleID string) (*App, error)
	InstanceFind(appID string, token Token) (*Instance, error)
}

type Handler struct {
	Storage Storage
}

func NewHandler(storage Storage) *Handler {
	return &Handler{
		Storage: storage,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	apnsID := r.Header.Get(headerAPNSID)
	w.Header().Set(headerAPNSID, apnsID)

	reqPathMatch := rPath.FindAllStringSubmatch(r.URL.Path, -1)
	if reqPathMatch == nil {
		h.writeError(w, r, ReasonMissingDeviceToken)
		return
	}

	var pl NotificationPayload
	_ = json.NewDecoder(r.Body).Decode(&pl)

	if (NotificationPayload{}) == pl {
		h.writeError(w, r, ReasonPayloadEmpty)
		return
	}

	topic := r.Header.Get(headerAPNSTopic)
	if topic == "" {
		h.writeError(w, r, ReasonMissingTopic)
		return
	}

	app, err := h.Storage.AppFind(topic)
	if err == ErrElementNotFound {
		h.writeError(w, r, ReasonBadTopic)
		return
	}

	// token is second element from from 1st matching group
	instance, err := h.Storage.InstanceFind(app.ID, Token(reqPathMatch[0][1]))
	if instance.State == InstanceStateUnregistered {
		h.writeErrorUnregistered(w, r, instance.LastSeen)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) writeErrorUnregistered(w http.ResponseWriter, r *http.Request, lastSeen int64) {
	reason := ReasonUnregistered
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(reason.Status())
	response := map[string]interface{}{
		"reason":    string(reason),
		"timestamp": lastSeen,
	}
	_ = json.NewEncoder(w).Encode(&response)
}

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, reason ResultReason) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(reason.Status())
	response := map[string]interface{}{
		"reason": string(reason),
	}
	_ = json.NewEncoder(w).Encode(&response)
}
