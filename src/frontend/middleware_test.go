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
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestEnsureSessionIDSetsCookieWhenMissing(t *testing.T) {
	t.Setenv("ENABLE_SINGLE_SHARED_SESSION", "")

	var seenSession string
	handler := ensureSessionID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenSession = sessionID(r)
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if seenSession == "" {
		t.Fatal("expected session id in request context")
	}
	res := rec.Result()
	cookies := res.Cookies()
	if len(cookies) != 1 || cookies[0].Name != cookieSessionID {
		t.Fatalf("unexpected cookies: %+v", cookies)
	}
	if cookies[0].Value != seenSession {
		t.Fatalf("cookie value %q does not match context session %q", cookies[0].Value, seenSession)
	}
	if cookies[0].MaxAge != cookieMaxAge {
		t.Fatalf("cookie MaxAge = %d; want %d", cookies[0].MaxAge, cookieMaxAge)
	}
}

func TestEnsureSessionIDUsesSharedSessionWhenEnabled(t *testing.T) {
	t.Setenv("ENABLE_SINGLE_SHARED_SESSION", "true")

	var seenSession string
	handler := ensureSessionID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenSession = sessionID(r)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	const shared = "12345678-1234-1234-1234-123456789123"
	if seenSession != shared {
		t.Fatalf("session = %q; want shared %q", seenSession, shared)
	}
}

func TestEnsureSessionIDReusesExistingCookie(t *testing.T) {
	const existing = "existing-session-id"
	handler := ensureSessionID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sessionID(r) != existing {
			t.Fatalf("session = %q; want %q", sessionID(r), existing)
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: cookieSessionID, Value: existing})
	rec := httptest.NewRecorder()
	handler(rec, req)

	if len(rec.Result().Cookies()) != 0 {
		t.Fatal("should not set a new session cookie when one is present")
	}
}

func TestResponseRecorderTracksStatusAndBytes(t *testing.T) {
	base := httptest.NewRecorder()
	rr := &responseRecorder{w: base}

	rr.WriteHeader(http.StatusCreated)
	_, _ = rr.Write([]byte("hello"))

	if rr.status != http.StatusCreated {
		t.Fatalf("status = %d; want %d", rr.status, http.StatusCreated)
	}
	if rr.b != 5 {
		t.Fatalf("bytes = %d; want 5", rr.b)
	}
}

func TestLogHandlerAttachesRequestLogger(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	var hasLogger bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, hasLogger = r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/cart", nil)
	ctx := context.WithValue(req.Context(), ctxKeySessionID{}, "trace-session")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	(&logHandler{log: log, next: next}).ServeHTTP(rec, req)

	if !hasLogger {
		t.Fatal("expected logrus logger in request context")
	}
}
