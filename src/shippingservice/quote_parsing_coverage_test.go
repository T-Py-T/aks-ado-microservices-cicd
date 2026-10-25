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

import "testing"

func TestCreateQuoteFromFloatSplitsUnitsAndCents(t *testing.T) {
	q := CreateQuoteFromFloat(8.99)
	if q.Dollars != 8 || q.Cents != 99 {
		t.Errorf("CreateQuoteFromFloat(8.99): got %d.%02d want 8.99", q.Dollars, q.Cents)
	}
}

func TestQuoteStringFormatsDollarAmount(t *testing.T) {
	q := Quote{Dollars: 8, Cents: 99}
	if got, want := q.String(), "$8.99"; got != want {
		t.Errorf("Quote.String: got %q want %q", got, want)
	}
}

func TestCreateQuoteFromCountUsesFlatRate(t *testing.T) {
	one := CreateQuoteFromCount(1)
	many := CreateQuoteFromCount(99)
	if one != many {
		t.Errorf("CreateQuoteFromCount should not vary with count: one=%v many=%v", one, many)
	}
	if one.Dollars != 8 || one.Cents != 99 {
		t.Errorf("expected flat 8.99 quote, got %v", one)
	}
}
