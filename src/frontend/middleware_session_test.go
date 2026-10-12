// Copyright 2018 Google LLC
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

func TestEnsureSessionIDSetsCookieWhenMissing(t *testing.T) {
	t.Setenv("ENABLE_SINGLE_SHARED_SESSION", "false")

	var captured string
	handler := ensureSessionID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.Context().Value(ctxKeySessionID{}).(string)
	}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler(rr, req)

	if captured == "" {
		t.Fatal("expected generated session id in context")
	}
	cookie := rr.Result().Cookies()[0]
	if cookie.Name != cookieSessionID || cookie.Value != captured {
		t.Fatalf("cookie %#v does not match context session %q", cookie, captured)
	}
	if cookie.MaxAge != cookieMaxAge {
		t.Fatalf("cookie MaxAge = %d, want %d", cookie.MaxAge, cookieMaxAge)
	}
}

func TestEnsureSessionIDReusesExistingCookie(t *testing.T) {
	var captured string
	handler := ensureSessionID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.Context().Value(ctxKeySessionID{}).(string)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: cookieSessionID, Value: "existing-session"})
	rr := httptest.NewRecorder()
	handler(rr, req)

	if captured != "existing-session" {
		t.Fatalf("session id = %q, want existing-session", captured)
	}
	if len(rr.Result().Cookies()) != 0 {
		t.Fatal("should not set a new cookie when one already exists")
	}
}

func TestEnsureSessionIDSharedSessionMode(t *testing.T) {
	t.Setenv("ENABLE_SINGLE_SHARED_SESSION", "true")

	var captured string
	handler := ensureSessionID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.Context().Value(ctxKeySessionID{}).(string)
	}))

	rr := httptest.NewRecorder()
	handler(rr, httptest.NewRequest(http.MethodGet, "/", nil))

	const shared = "12345678-1234-1234-1234-123456789123"
	if captured != shared {
		t.Fatalf("shared session = %q, want %q", captured, shared)
	}
}

func TestResponseRecorderCapturesStatusAndBytes(t *testing.T) {
	base := httptest.NewRecorder()
	rec := &responseRecorder{w: base}
	rec.WriteHeader(http.StatusCreated)
	rec.Write([]byte("hello"))

	if rec.status != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.status, http.StatusCreated)
	}
	if rec.b != 5 {
		t.Fatalf("bytes written = %d, want 5", rec.b)
	}
}

func TestLogHandlerInjectsRequestID(t *testing.T) {
	log := logrus.New()
	var requestID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID = r.Context().Value(ctxKeyRequestID{}).(string)
	})
	handler := &logHandler{log: log, next: next}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/cart", nil))

	if requestID == "" {
		t.Fatal("logHandler should attach a request id to the context")
	}
}
