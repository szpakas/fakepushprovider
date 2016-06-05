package http

import "github.com/szpakas/fakepushprovider/fcm"

// DownstreamMessage is a message designed for application on device and send from app server to GCM.
// Most of the documentation on keys taken from GCM and/or APNS docs.
type DownstreamMessage struct {
	// To specifies the recipient of a message.
	// The value must be a registration token or notification key.
	To fcm.RegistrationID `json:"to,omitempty"`

	// RegistrationIDS specifies a list of devices (registration tokens, or IDs) receiving a multicast message.
	// It must contain at least 1 and at most 1000 registration tokens.
	// Use this parameter only for multicast messaging, not for single recipients.
	// Multicast messages (sending to more than 1 registration tokens) are allowed using HTTP JSON format only.
	RegistrationIDS []fcm.RegistrationID `json:"registration_ids,omitempty"`

	// CollapseKey identifies a group of messages (e.g., with collapse_key: "Updates Available")
	// that can be collapsed, so that only the last message gets sent when delivery can be resumed.
	//CollapseKey string `json:"collapse_key,omitempty"`

	// Priority is a priority of the message.
	// Valid values are: normal (default), high
	//Priority string `json:"priority,omitempty"`

	// ContentAvailable represent content-available in the APNS payload.
	// Not used for now.
	// required: false
	//ContentAvailable         bool   `json:"content_available,omitempty"`

	// DelayWhileIdle if set to true, it indicates that the message should not be sent until the device becomes active.
	// required: false (default: false)
	//DelayWhileIdle bool `json:"delay_while_idle,omitempty"`

	// TimeToLive specifies how many seconds message should be kept in GCM storage if the device is offline.
	//
	// If set to 0, than GCM guarantees best effort for messages that must be delivered "now or never."
	// Message in such case will be discarded if can't be delivered immediately.
	//
	// required: false (default: 4 weeks - 2419200 seconds, max: 4 weeks)
	//TimeToLive uint `json:"time_to_live,omitempty"`

	//RestrictedPackageName string `json:"restricted_package_name"`

	// DeliveryReceiptRequested lets the app server request confirmation of message delivery.
	//
	// When this parameter is set to true, CCS sends a delivery receipt when the device confirms that it received the message.
	// The receipt is NOT reliable and triggers lots of false negatives (meaning that lack of receipt
	// should NOT be understood as if the message was not delivered)
	//DeliveryReceiptRequested bool `json:"delivery_receipt_requested"`

	// DryRun allows developers to test a request without actually sending a message.
	//DryRun bool `json:"dry_run,omitempty"`

	// Data is a JSON object (key->value) encoded to string containing extra information to be passed to app.
	//
	// The key should not be a reserved word ("from" or any word starting with "google" or "gcm").
	// Do not use any of the words defined in this table (such as collapse_key).
	// Values in string types are recommended. You have to convert values in objects
	// or other non-string data types (e.g., integers or booleans) to string.
	//Data string `json:"data,omitempty"`

	// Notification is an GCM Notification object
	Notification *Notification `json:"notification,omitempty"`
}

// Notification is a high-level model of message send to providers.
// It's aimed to be used by both APNS and GCM.
type Notification struct {
	// GCM: notification.title
	// APNS: aps.alert.title
	Title string `json:"title,omitempty"`

	// GCM: notification.body
	// APNS: aps.alert.body
	Body string `json:"body,omitempty"`

	// GCM: notification.icon
	// APNS: -
	Icon string `json:"icon,omitempty"`

	// GCM: notification.sound
	// APNS: aps.sound
	Sound string `json:"sound,omitempty"`

	//Badge        string `json:"badge,omitempty"`
	//Tag          string `json:"tag,omitempty"`
	//Color        string `json:"color,omitempty"`
	//ClickAction  string `json:"click_action,omitempty"`
	//BodyLocKey   string `json:"body_loc_key,omitempty"`
	//BodyLocArgs  string `json:"body_loc_args,omitempty"`
	//TitleLocArgs string `json:"title_loc_args,omitempty"`
	//TitleLocKey  string `json:"title_loc_key,omitempty"`
}

type DownstreamResponse struct {
	MulticastID  int             `json:"multicast_id,omitempty"`
	Success      int             `json:"success"`
	Failure      int             `json:"failure"`
	CanonicalIDS int             `json:"canonical_ids"`
	Results      []MessageResult `json:"results"`
}

type MessageResult struct {
	MessageID      string             `json:"message_id,omitempty"`
	RegistrationID fcm.RegistrationID `json:"registration_id,omitempty"`
	Error          DownstreamError    `json:"error,omitempty"`
}

type DownstreamError string

//noinspection GoUnusedConst
const (
	// GCM docs, Downstream message error response codes (Table 9)
	InvalidJSON               DownstreamError = "INVALID_JSON"
	BadRegistration           DownstreamError = "BAD_REGISTRATION"
	DeviceUnregistered        DownstreamError = "DEVICE_UNREGISTERED"
	BadACK                    DownstreamError = "BAD_ACK"
	ServiceUnavailable        DownstreamError = "SERVICE_UNAVAILABLE"
	InternalServerError       DownstreamError = "INTERNAL_SERVER_ERROR"
	DeviceMessageRateExceeded DownstreamError = "DEVICE_MESSAGE_RATE_EXCEEDED"
	TopicMessageRateExceeded  DownstreamError = "TOPICS_MESSAGE_RATE_EXCEEDED"
	ConnectionDraining        DownstreamError = "CONNECTION_DRAINING"
)
