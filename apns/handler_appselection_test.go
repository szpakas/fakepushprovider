package apns

import (
	"fmt"
	"net/http"
	"testing"

	//a "github.com/stretchr/testify/assert"
	//ar "github.com/stretchr/testify/require"

	"github.com/davecgh/go-spew/spew"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
)

var _ = spew.Config

func Test_Handler_Responses(t *testing.T) {
	pay := map[string]interface{}{
		"aps": map[string]interface{}{
			"alert": map[string]interface{}{
				"title": "Title A",
				"body":  "Body A",
			},
		},
	}
	tests := map[string]struct {
		req     thsPushReq
		status  int
		eReason ResultReason
		// lastSeen is timestamp when device was last seen
		lastSeen int64
	}{
		"success": {
			req: thsPushReq{
				BundleID: TFAppA.BundleID,
				Token:    TFInsAA.Token,
				Payload:  pay,
			},
			status: http.StatusOK,
		},
		"failure, unknown app": {
			req: thsPushReq{
				BundleID: "FakeBundleID",
				Token:    TFInsAA.Token,
				Payload:  pay,
			},
			eReason: ReasonBadTopic,
			status:  http.StatusBadRequest,
		},
		"failure, empty payload": {
			req: thsPushReq{
				BundleID: TFAppA.BundleID,
				Token:    TFInsAA.Token,
			},
			eReason: ReasonPayloadEmpty,
			status:  http.StatusBadRequest,
		},
		"failure, missing token": {
			req: thsPushReq{
				BundleID: TFAppA.BundleID,
				Payload:  pay,
			},
			eReason: ReasonMissingDeviceToken,
			status:  http.StatusBadRequest,
		},
		"failure, unregistered": {
			req: thsPushReq{
				BundleID: TFAppA.BundleID,
				Token:    TFInsAZ.Token,
				Payload:  pay,
			},
			eReason:  ReasonUnregistered,
			status:   http.StatusGone,
			lastSeen: TFInsAZ.LastSeen,
		},
		"failure, missing topic (bundleID)": {
			req: thsPushReq{
				Token:   TFInsAA.Token,
				Payload: pay,
			},
			eReason: ReasonMissingTopic,
			status:  http.StatusBadRequest,
		},
	}

	for sym, tc := range tests {
		_, _, _, e, closer := tsServerSetup(t, sym)

		apnsID := uuid.NewV1().String()

		req := tc.req
		req.APNSID = apnsID
		req.e = e

		resp := req.Req().Expect()

		if tc.eReason == ResultReason("") {
			thAssertSuccessResponse(t, resp, apnsID)
		} else {
			thAssertErrorResponse(t, resp, apnsID, tc.eReason, tc.lastSeen)
		}

		closer()
	}
}

// -- helpers

func thAssertSuccessResponse(t *testing.T, resp *httpexpect.Response, apnsID string) {
	resp.Status(http.StatusOK)
	resp.Body().Empty()
	resp.Header("apns-id").Equal(apnsID)
	//resp.Header(":status").
}

func thAssertErrorResponse(t *testing.T, resp *httpexpect.Response, apnsID string, reason ResultReason, lastSeen int64) {
	resp.Header("apns-id").Equal(apnsID)

	var status int
	switch reason {
	default:
		status = http.StatusBadRequest
	}

	resp.Status(status)
	resp.ContentTypeJSON()

	respObj := resp.JSON().Object()

	keys := []interface{}{"reason"}
	if lastSeen != int64(0) {
		keys = append(keys, "timestamp")
	}
	respObj.Keys().ContainsOnly(keys...)
	respObj.Value("reason").String().Equal(string(reason))

	if lastSeen != int64(0) {
		respObj.Value("timestamp").Number().Equal(lastSeen)
	}
}

type thsPushReq struct {
	e *httpexpect.Expect

	BundleID string
	APNSID   string
	Token    Token
	Payload  map[string]interface{}
}

func (r *thsPushReq) Req() *httpexpect.Request {
	req := r.e.POST(fmt.Sprintf("/3/device/%s", string(r.Token))).
		WithHeader("apns-topic", r.BundleID).
		WithHeader("apns-id", r.APNSID).
		WithJSON(r.Payload)

	return req
}
