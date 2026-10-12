// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestResponseRecorderDefaultsStatusOnWrite(t *testing.T) {
	rr := httptest.NewRecorder()
	rec := &responseRecorder{w: rr}

	rec.Write([]byte("ok"))
	if rec.status != http.StatusOK {
		t.Errorf("status after Write: got %d want %d", rec.status, http.StatusOK)
	}
	if rec.b != 2 {
		t.Errorf("bytes written: got %d want 2", rec.b)
	}
}

func TestResponseRecorderHonorsWriteHeader(t *testing.T) {
	rr := httptest.NewRecorder()
	rec := &responseRecorder{w: rr}
	rec.WriteHeader(http.StatusNotFound)
	rec.Write([]byte("missing"))
	if rec.status != http.StatusNotFound {
		t.Errorf("status: got %d want %d", rec.status, http.StatusNotFound)
	}
}

func TestLogHandlerInjectsRequestID(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)

	var capturedID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID, _ = r.Context().Value(ctxKeyRequestID{}).(string)
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	(&logHandler{log: logger, next: next}).ServeHTTP(rr, req)

	if capturedID == "" {
		t.Fatal("expected non-empty request id in context")
	}
	if rr.Code != http.StatusNoContent {
		t.Errorf("response code: got %d want %d", rr.Code, http.StatusNoContent)
	}
}

func TestEnsureSessionIDSetsCookieWhenMissing(t *testing.T) {
	t.Setenv("ENABLE_SINGLE_SHARED_SESSION", "")

	var sessionFromContext string
	handler := ensureSessionID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionFromContext, _ = r.Context().Value(ctxKeySessionID{}).(string)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if sessionFromContext == "" {
		t.Fatal("expected session id in context")
	}
	cookies := rr.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != cookieSessionID {
		t.Fatalf("expected one session cookie, got %#v", cookies)
	}
	if cookies[0].Value != sessionFromContext {
		t.Errorf("cookie value %q does not match context session %q", cookies[0].Value, sessionFromContext)
	}
}

func TestEnsureSessionIDUsesSharedSessionWhenEnabled(t *testing.T) {
	t.Setenv("ENABLE_SINGLE_SHARED_SESSION", "true")
	const shared = "12345678-1234-1234-1234-123456789123"

	var sessionFromContext string
	handler := ensureSessionID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionFromContext, _ = r.Context().Value(ctxKeySessionID{}).(string)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if sessionFromContext != shared {
		t.Errorf("shared session: got %q want %q", sessionFromContext, shared)
	}
}
