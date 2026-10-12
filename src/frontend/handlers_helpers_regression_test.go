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

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"
)

func TestCartHelpers(t *testing.T) {
	items := []*pb.CartItem{
		{ProductId: "abc", Quantity: 2},
		{ProductId: "def", Quantity: 3},
	}
	if got := cartIDs(items); len(got) != 2 || got[0] != "abc" || got[1] != "def" {
		t.Fatalf("cartIDs returned %#v", got)
	}
	if got := cartSize(items); got != 5 {
		t.Fatalf("cartSize = %d, want 5", got)
	}
	if cartSize(nil) != 0 {
		t.Fatal("cartSize(nil) should be 0")
	}
}

func TestRenderMoneyAndCurrencyLogo(t *testing.T) {
	m := pb.Money{CurrencyCode: "EUR", Units: 12, Nanos: 340000000}
	if got := renderMoney(m); got != "€12.34" {
		t.Fatalf("renderMoney = %q, want €12.34", got)
	}
	if renderCurrencyLogo("UNKNOWN") != "$" {
		t.Fatal("unknown currency should default to $ logo")
	}
}

func TestStringInSlice(t *testing.T) {
	slice := []string{"gcp", "aws", "azure"}
	if !stringinSlice(slice, "aws") {
		t.Fatal("expected aws to be found")
	}
	if stringinSlice(slice, "onprem") {
		t.Fatal("onprem should not be in slice")
	}
}

func TestCurrentCurrencyCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if currentCurrency(req) != defaultCurrency {
		t.Fatalf("default currency = %q, want %q", currentCurrency(req), defaultCurrency)
	}
	req.AddCookie(&http.Cookie{Name: cookieCurrency, Value: "JPY"})
	if currentCurrency(req) != "JPY" {
		t.Fatalf("cookie currency = %q, want JPY", currentCurrency(req))
	}
}

func TestSessionIDFromContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if sessionID(req) != "" {
		t.Fatal("missing session context should return empty string")
	}
	ctx := context.WithValue(req.Context(), ctxKeySessionID{}, "session-123")
	req = req.WithContext(ctx)
	if sessionID(req) != "session-123" {
		t.Fatalf("sessionID = %q, want session-123", sessionID(req))
	}
}

func TestPlatformDetailsMapping(t *testing.T) {
	cases := map[string]string{
		"gcp":     "Google Cloud",
		"aws":     "AWS",
		"azure":   "Azure",
		"onprem":  "On-Premises",
		"alibaba": "Alibaba Cloud",
		"":        "local",
	}
	for env, provider := range cases {
		plat := &platformDetails{}
		plat.setPlatformDetails(env)
		if plat.provider != provider {
			t.Fatalf("env %q provider = %q, want %q", env, plat.provider, provider)
		}
	}
}
