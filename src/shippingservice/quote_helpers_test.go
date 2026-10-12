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
		name   string
		value  float64
		dollar uint32
		cents  uint32
	}{
		{"flat shipping rate", 8.99, 8, 99},
		{"whole dollars", 12.0, 12, 0},
		{"fraction rounds down at cent boundary", 0.999, 0, 99},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateQuoteFromFloat(tt.value)
			if got.Dollars != tt.dollar || got.Cents != tt.cents {
				t.Fatalf("CreateQuoteFromFloat(%v) = %+v, want dollars=%d cents=%d",
					tt.value, got, tt.dollar, tt.cents)
			}
		})
	}
}

func TestCreateQuoteFromCountIgnoresCount(t *testing.T) {
	for _, count := range []int{0, 1, 25} {
		got := CreateQuoteFromCount(count)
		want := CreateQuoteFromFloat(8.99)
		if got != want {
			t.Fatalf("CreateQuoteFromCount(%d) = %+v, want %+v", count, got, want)
		}
	}
}

func TestQuoteString(t *testing.T) {
	got := Quote{Dollars: 8, Cents: 99}.String()
	if got != "$8.99" {
		t.Fatalf("Quote.String() = %q, want %q", got, "$8.99")
	}
}
