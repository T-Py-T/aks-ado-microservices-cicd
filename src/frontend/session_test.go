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
	"net/http/httptest"
	"testing"
)

func TestSessionIDFromContext(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	if sessionID(req) != "" {
		t.Fatalf("expected empty session without context value")
	}

	req = req.WithContext(context.WithValue(req.Context(), ctxKeySessionID{}, "ctx-session"))
	if sessionID(req) != "ctx-session" {
		t.Fatalf("sessionID = %q", sessionID(req))
	}
}

func TestWhitelistedCurrencies(t *testing.T) {
	for code := range whitelistedCurrencies {
		if !whitelistedCurrencies[code] {
			t.Fatalf("currency %q should be whitelisted", code)
		}
	}
	if whitelistedCurrencies["CHF"] {
		t.Fatal("CHF should not be whitelisted")
	}
}
