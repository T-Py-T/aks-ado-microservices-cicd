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

package money

import "testing"

func TestMultiplySlowScalesByRepeatingSum(t *testing.T) {
	base := mmc(2, 250000000, "USD")
	got := MultiplySlow(base, 3)
	want := mmc(6, 750000000, "USD")
	if !AreEquals(got, want) {
		t.Errorf("MultiplySlow(2.25 USD, 3): got %v want %v", got, want)
	}
}

func TestMultiplySlowIdentityForOne(t *testing.T) {
	base := mmc(10, 0, "EUR")
	got := MultiplySlow(base, 1)
	if !AreEquals(got, base) {
		t.Errorf("MultiplySlow(n, 1) should return input: got %v", got)
	}
}

func TestMultiplySlowZeroMultiplier(t *testing.T) {
	base := mmc(5, 100000000, "USD")
	got := MultiplySlow(base, 0)
	if !AreEquals(got, base) {
		t.Errorf("MultiplySlow(n, 0) returns unchanged value in this implementation: got %v", got)
	}
}
