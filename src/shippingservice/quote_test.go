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

import "testing"

func TestCreateQuoteFromFloat(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		dollars uint32
		cents   uint32
	}{
		{"flat shipping rate", 8.99, 8, 99},
		{"whole dollars", 12.00, 12, 0},
		{"fractional cents truncate", 1.999, 1, 99},
		{"zero", 0, 0, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CreateQuoteFromFloat(tc.value)
			if got.Dollars != tc.dollars || got.Cents != tc.cents {
				t.Fatalf("CreateQuoteFromFloat(%v) = %+v; want dollars=%d cents=%d", tc.value, got, tc.dollars, tc.cents)
			}
		})
	}
}

func TestCreateQuoteFromCountUsesFlatRate(t *testing.T) {
	for _, count := range []int{0, 1, 42} {
		got := CreateQuoteFromCount(count)
		want := CreateQuoteFromFloat(8.99)
		if got != want {
			t.Fatalf("CreateQuoteFromCount(%d) = %+v; want flat %+v", count, got, want)
		}
	}
}

func TestQuoteString(t *testing.T) {
	if CreateQuoteFromFloat(8.99).String() != "$8.99" {
		t.Fatalf("unexpected string: %q", CreateQuoteFromFloat(8.99).String())
	}
}
