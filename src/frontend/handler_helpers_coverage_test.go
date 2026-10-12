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
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"
)

func TestCartIDsExtractsProductIds(t *testing.T) {
	items := []*pb.CartItem{
		{ProductId: "abc", Quantity: 2},
		{ProductId: "xyz", Quantity: 1},
	}
	got := cartIDs(items)
	want := []string{"abc", "xyz"}
	if len(got) != len(want) {
		t.Fatalf("cartIDs length: got %d want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("cartIDs[%d]: got %q want %q", i, got[i], want[i])
		}
	}
}

func TestCartSizeSumsQuantities(t *testing.T) {
	items := []*pb.CartItem{
		{ProductId: "a", Quantity: 2},
		{ProductId: "b", Quantity: 5},
	}
	if got, want := cartSize(items), 7; got != want {
		t.Errorf("cartSize: got %d want %d", got, want)
	}
	if cartSize(nil) != 0 {
		t.Errorf("cartSize(nil) should be 0")
	}
}

func TestRenderMoneyUsesCurrencyLogoAndNanos(t *testing.T) {
	m := pb.Money{CurrencyCode: "EUR", Units: 12, Nanos: 500000000}
	if got, want := renderMoney(m), "€12.50"; got != want {
		t.Errorf("renderMoney: got %q want %q", got, want)
	}
}

func TestRenderCurrencyLogoFallsBackToDollar(t *testing.T) {
	if got, want := renderCurrencyLogo("CHF"), "$"; got != want {
		t.Errorf("unknown currency logo: got %q want %q", got, want)
	}
	if got, want := renderCurrencyLogo("GBP"), "£"; got != want {
		t.Errorf("GBP logo: got %q want %q", got, want)
	}
}

func TestCurrentCurrencyReadsCookieOrDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if got, want := currentCurrency(req), defaultCurrency; got != want {
		t.Errorf("without cookie: got %q want %q", got, want)
	}

	req.AddCookie(&http.Cookie{Name: cookieCurrency, Value: "CAD"})
	if got, want := currentCurrency(req), "CAD"; got != want {
		t.Errorf("with cookie: got %q want %q", got, want)
	}
}

func TestSessionIDFromContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if sessionID(req) != "" {
		t.Errorf("expected empty session without context value")
	}

	ctx := context.WithValue(req.Context(), ctxKeySessionID{}, "sess-123")
	req = req.WithContext(ctx)
	if got, want := sessionID(req), "sess-123"; got != want {
		t.Errorf("sessionID: got %q want %q", got, want)
	}
}

func TestStringInSlice(t *testing.T) {
	slice := []string{"local", "gcp", "azure"}
	if !stringinSlice(slice, "gcp") {
		t.Error("expected gcp to be in slice")
	}
	if stringinSlice(slice, "onprem") {
		t.Error("expected onprem not to be in slice")
	}
}
