package apns

import (
	"crypto/tls"
	"net/http"
)

// App is a single application in iTunes.
type App struct {
	// ID is and internal id for the app
	ID string

	// Certificate is a SSL certification used to negotiate TLS connection.
	// It's not exported currently as support for certs is missing.
	Certificate tls.Certificate `json:"-"`

	// BundleID represents app in iTunes.
	// It's used to select app when APNS connection is established by matching "apns-topic" or certificate subject.
	BundleID string
}

type Token string

type Instance struct {
	ID    string
	State InstanceState
	App   *App
	Token Token

	// LastSeen is timestamp when device was last seen.
	// Used when device is unregistered.
	LastSeen int64 `json:",omitempty"`
}

type InstanceState int

const (
	InstanceStateRegistered = iota + 1
	InstanceStateUnregistered
)

type MessagePriority int

const (
	PriorityHigh MessagePriority = 10
	PriorityLow  MessagePriority = 5
)

// RemoteNotification holds a single communication to be submitted to the APNs provider.
type RemoteNotification struct {
	// Token represents a recipient device token.
	Token string

	// Topic is the name to which the recipient should subscribe to.
	// Typically it's bundle ID of the app.
	Topic string

	// Expiration identifies the date when the notification is no longer valid and can be discarded.
	// It's a UNIX epoch date expressed in seconds (UTC).
	// If this value is nonzero, APNs stores the notification and tries to deliver it at least once,
	// repeating the attempt as needed if it is unable to deliver the notification the first time.
	// If the value is 0, APNs treats the notification as if it expires immediately
	// and does not store the notification or attempt to redeliver it.
	Expiration int64

	// Priority sets the priority of the push in APNS.
	// Specify one of the following values:
	// - 10: Send the push message immediately.
	//       Notifications with this priority must trigger an alert, sound, or badge on the target device.
	//       It is an error to use this priority for a push notification that contains only the content-available key.
	// -  5: Send the push message at a time that takes into account power considerations for the device.
	//       Notifications with this priority might be grouped and delivered in bursts.
	//       They are throttled, and in some cases are not delivered.
	// If you omit this header, the APNs server sets the priority to 10.
	Priority MessagePriority

	// Payload is a content which is send to the device.
	// It's described in "Remote Notification Payload" section of the APNs docs.
	// It have to be JSON encoded.
	Payload []byte
}

// NotificationPayload is a payload of the remote notification.
// As defined in tables 5-1 and 5-2 in APNS documentation ("The Remote Notification Payload" section)
type NotificationPayload struct {
	Aps struct {
		// Alert triggers system alert or banner.
		// Can be either string or dictionary (defined in NotificationPayloadApsAlert).
		Alert            interface{} `json:"alert"`
		Badge            int         `json:"badge"`
		Sound            string      `json:"sound"`
		ContentAvailable int         `json:"content-available"`
	} `json:"aps"`
}

// NotificationPayloadApsAlert is verbose structure for alert section
// See Table 5-2 in APNS docs.
type NotificationPayloadApsAlert struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Response represent a response from APNS to the single RemoteNotification.
type ErrorResponse struct {
	// Status is an HTTP status code of the response.
	//
	// Possible values as per APNS documentation:
	// 200  Success - no response other than StatusOK is produced
	// 400  Bad request
	// 403  There was an error with the certificate.
	// 405  The request used a bad :method value. Only POST requests are supported.
	// 410  The device token is no longer active for the topic.
	// 413  The notification payload was too large.
	// 429  The server received too many requests for the same device token.
	// 500  Internal server error
	// 503  The server is shutting down and unavailable.
	Status int `json:"-"`

	// Reason is indicating the reason for the failure.
	Reason ResultReason `json:"reason"`

	// If the value of StatusCode is 410, this is the last time at which APNs
	// confirmed that the device token was no longer valid for the topic.
	Timestamp int64 `json:"timestamp,omitempty"`
}

type ResultReason string

func (r *ResultReason) Status() int {
	return http.StatusBadRequest
}

// The possible Reason error codes returned from APNs.
// From table 6-6 in the Apple Local and Remote Notification Programming Guide.
// taken from sideshow/apns2 package (MIT license)
const (
	// The message payload was empty.
	ReasonPayloadEmpty ResultReason = "PayloadEmpty"
	// The message payload was too large. The maximum payload size is 4096 bytes.
	ReasonPayloadTooLarge ResultReason = "PayloadTooLarge"
	// The apns-topic was invalid.
	ReasonBadTopic ResultReason = "BadTopic"
	// Pushing to this topic is not allowed.
	ReasonTopicDisallowed ResultReason = "TopicDisallowed"
	// The apns-id value is bad.
	ReasonBadMessageID ResultReason = "BadMessageId"
	// The apns-expiration value is bad.
	ReasonBadExpirationDate ResultReason = "BadExpirationDate"
	// The apns-priority value is bad.
	ReasonBadPriority ResultReason = "BadPriority"
	// The device token is not specified in the request :path. Verify that the
	// :path header contains the device token.
	ReasonMissingDeviceToken ResultReason = "MissingDeviceToken"
	// The specified device token was bad. Verify that the request contains a valid
	// token and that the token matches the environment.
	ReasonBadDeviceToken ResultReason = "BadDeviceToken"
	// The device token does not match the specified topic.
	ReasonDeviceTokenNotForTopic ResultReason = "DeviceTokenNotForTopic"
	// The device token is inactive for the specified topic.
	ReasonUnregistered ResultReason = "Unregistered"
	// One or more headers were repeated.
	ReasonDuplicateHeaders ResultReason = "DuplicateHeaders"
	// The client certificate was for the wrong environment.
	ReasonBadCertificateEnvironment ResultReason = "BadCertificateEnvironment"
	// The certificate was bad.
	ReasonBadCertificate ResultReason = "BadCertificate"
	// The specified action is not allowed.
	ReasonForbidden ResultReason = "Forbidden"
	// The request contained a bad :path value.
	ReasonBadPath ResultReason = "BadPath"
	// The specified :method was not POST.
	ReasonMethodNotAllowed ResultReason = "MethodNotAllowed"
	// Too many requests were made consecutively to the same device token.
	ReasonTooManyRequests ResultReason = "TooManyRequests"
	// Idle time out.
	ReasonIdleTimeout ResultReason = "IdleTimeout"
	// The server is shutting down.
	ReasonShutdown ResultReason = "Shutdown"
	// An internal server error occurred.
	ReasonInternalServerError ResultReason = "InternalServerError"
	// The service is unavailable.
	ReasonServiceUnavailable ResultReason = "ServiceUnavailable"
	// The apns-topic header of the request was not specified and was required.
	// The apns-topic header is mandatory when the client is connected using a
	// certificate that supports multiple topics.
	ReasonMissingTopic ResultReason = "MissingTopic"
)
