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
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func requestWithLogger(method, target, body string) *http.Request {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	log := logrus.New()
	ctx := context.WithValue(req.Context(), ctxKeyLog{}, log)
	return req.WithContext(ctx)
}

func TestSetCurrencyRejectsUnsupportedCodeRegression(t *testing.T) {
	fe := &frontendServer{}
	req := requestWithLogger(http.MethodPost, "/setCurrency", "currency_code=BTC")
	rec := httptest.NewRecorder()

	fe.setCurrencyHandler(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
	for _, c := range rec.Result().Cookies() {
		if c.Name == cookieCurrency {
			t.Fatalf("should not set currency cookie for invalid code, got %q", c.Value)
		}
	}
}

func TestSetCurrencyAcceptsWhitelistedCodeRegression(t *testing.T) {
	fe := &frontendServer{}
	req := requestWithLogger(http.MethodPost, "/setCurrency", "currency_code=EUR")
	rec := httptest.NewRecorder()

	fe.setCurrencyHandler(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusFound)
	}
	var found bool
	for _, c := range rec.Result().Cookies() {
		if c.Name == cookieCurrency && c.Value == "EUR" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected shop_currency=EUR cookie")
	}
}
