package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	a "github.com/stretchr/testify/assert"
	ar "github.com/stretchr/testify/require"

	"github.com/szpakas/fakepushprovider/common"
	"github.com/szpakas/fakepushprovider/fcm"
)

func Test_Handler_Factory(t *testing.T) {
	_, h, _, _, closer := tsServerSetup(t, "")
	defer closer()

	a.IsType(t, &Handler{}, h, "Incorrect type")
	ar.NotNil(t, h, "empty handler")
	a.NotNil(t, h.Storage, "storage not defined")
}

func Test_Handler_Response_Success(t *testing.T) {
	n := thsNotification{
		Title: "Portugal vs. Denmark",
		Body:  "5 to 1",
	}
	tests := map[string]struct {
		req                            thsPushReq
		success, failure, canonicalIDS int
		res                            []thsMessageResult
	}{
		"single, unicast, success": {
			thsPushReq{
				apiKey:       fcm.TFAppA.ApiKey,
				To:           fcm.TFInsAA.CanonicalID,
				Notification: n,
			},
			1, 0, 0,
			[]thsMessageResult{{Success: true}},
		},
		"single, multicast, success": {
			thsPushReq{
				apiKey:          fcm.TFAppA.ApiKey,
				RegistrationIDS: []fcm.RegistrationID{fcm.TFInsAB.CanonicalID},
				Notification:    n,
			},
			1, 0, 0,
			[]thsMessageResult{{Success: true}},
		},
		"multiple, multicast, success": {
			thsPushReq{
				apiKey: fcm.TFAppA.ApiKey,
				RegistrationIDS: []fcm.RegistrationID{
					fcm.TFInsAA.CanonicalID, fcm.TFInsAB.CanonicalID, fcm.TFInsAC.CanonicalID,
				},
				Notification: n,
			},
			3, 0, 0,
			[]thsMessageResult{{Success: true}, {Success: true}, {Success: true}},
		},
		"single, unicast, error, app exists, registration unknown": {
			thsPushReq{
				apiKey:       fcm.TFAppA.ApiKey,
				To:           "fakeRegID",
				Notification: n,
			},
			0, 1, 0,
			[]thsMessageResult{{Error: DeviceUnregistered}},
		},
		"single, unicast, error, app exists, instance unregistered": {
			thsPushReq{
				apiKey:       fcm.TFAppA.ApiKey,
				To:           fcm.TFInsAZ.RegistrationIDS[0],
				Notification: n,
			},
			0, 1, 0,
			[]thsMessageResult{{Error: DeviceUnregistered}},
		},
	}

	for symbol, td := range tests {
		_, _, _, e, closer := tsServerSetup(t, symbol)

		req := td.req
		req.e = e
		resp := req.Req().Expect()

		thAssertSuccessResponse(t, resp, td.success, td.failure, td.canonicalIDS, td.res)

		closer()
	}
}

// -- setup

func tsServerSetup(t *testing.T, symbol string) (*fcm.MemoryStorage, *Handler, *httptest.Server, *httpexpect.Expect, func()) {
	// -- storage
	st, stCloser := fcm.TSMemoryStorageWitAppsAndInstancesSetup()

	// -- handler
	h := NewHandler(st)

	srv := httptest.NewServer(h)

	// -- client test helper
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  srv.URL,
		Client:   http.DefaultClient,
		Reporter: &common.THAssertReporter{a.New(t), symbol},
	})

	closer := func() {
		stCloser()
		srv.Close()
	}

	return st, h, srv, e, closer
}

// -- helpers

type thsMessageResult struct {
	Success        bool
	RegistrationID fcm.RegistrationID
	Error          DownstreamError
}

func thAssertSuccessResponse(t *testing.T, resp *httpexpect.Response, success, failure, canonicalIDS int, mResExp []thsMessageResult) {
	resp.Status(http.StatusOK)
	resp.ContentTypeJSON()

	respObj := resp.JSON().Object()

	respObj.Keys().ContainsOnly("multicast_id", "success", "failure", "canonical_ids", "results")

	respObj.Value("multicast_id").Number().Gt(0)

	respObj.ValueEqual("success", success)
	respObj.ValueEqual("failure", failure)
	respObj.ValueEqual("canonical_ids", canonicalIDS)

	// match responses to tokens
	mRes := respObj.Value("results").Array()

	ar.EqualValues(t, len(mResExp), mRes.Length().Raw(), "incorrect number of results returned")
	for i, mR := range mResExp {
		if mR.Success {
			mRes.Element(i).Object().Value("message_id").String().NotEmpty()
			continue
		}
	}
}

type thsNotification struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

type thsPushReq struct {
	e      *httpexpect.Expect `json:"-"`
	apiKey string             `json:"-"`

	To              fcm.RegistrationID   `json:"to,omitempty"`
	RegistrationIDS []fcm.RegistrationID `json:"registration_ids,omitempty"`
	Notification    thsNotification      `json:"notification,omitempty"`
}

func (r *thsPushReq) Req() *httpexpect.Request {
	req := r.e.POST("/gcm/send").
		WithHeader("Authorization", fmt.Sprintf("key=%s", r.apiKey)).
		WithJSON(r)

	return req
}
