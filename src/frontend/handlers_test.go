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

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"
)

func TestCurrentCurrencyUsesCookieOrDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if currentCurrency(req) != defaultCurrency {
		t.Fatalf("expected default %q, got %q", defaultCurrency, currentCurrency(req))
	}

	req.AddCookie(&http.Cookie{Name: cookieCurrency, Value: "EUR"})
	if currentCurrency(req) != "EUR" {
		t.Fatalf("expected cookie currency EUR, got %q", currentCurrency(req))
	}
}

func TestCartSizeAndCartIDs(t *testing.T) {
	items := []*pb.CartItem{
		{ProductId: "A", Quantity: 2},
		{ProductId: "B", Quantity: 3},
	}
	if cartSize(items) != 5 {
		t.Fatalf("cartSize = %d; want 5", cartSize(items))
	}
	if ids := cartIDs(items); len(ids) != 2 || ids[0] != "A" || ids[1] != "B" {
		t.Fatalf("cartIDs = %v; want [A B]", ids)
	}
	if cartSize(nil) != 0 || len(cartIDs(nil)) != 0 {
		t.Fatalf("nil cart should be empty")
	}
}

func TestRenderMoneyAndCurrencyLogo(t *testing.T) {
	m := pb.Money{CurrencyCode: "EUR", Units: 12, Nanos: 340000000}
	if renderMoney(m) != "€12.34" {
		t.Fatalf("renderMoney = %q; want €12.34", renderMoney(m))
	}
	if renderCurrencyLogo("JPY") != "¥" {
		t.Fatalf("unexpected JPY logo")
	}
	if renderCurrencyLogo("UNKNOWN") != "$" {
		t.Fatalf("unknown currency should default to $")
	}
}

func TestStringInSlice(t *testing.T) {
	slice := []string{"local", "gcp", "azure"}
	if !stringinSlice(slice, "gcp") || stringinSlice(slice, "aws") {
		t.Fatalf("stringinSlice returned unexpected result")
	}
}
