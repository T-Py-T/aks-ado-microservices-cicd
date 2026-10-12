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

package money

import (
	"testing"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"
)

func TestMultiplySlowScalesValidMoney(t *testing.T) {
	base := pb.Money{Units: 2, Nanos: 500000000, CurrencyCode: "USD"}
	got := MultiplySlow(base, 3)
	want, err := Sum(base, base)
	if err != nil {
		t.Fatalf("Sum failed: %v", err)
	}
	want, err = Sum(want, base)
	if err != nil {
		t.Fatalf("Sum failed: %v", err)
	}
	if !AreEquals(got, want) {
		t.Fatalf("MultiplySlow = %v, want %v", got, want)
	}
}

func TestMultiplySlowIdentityAndZero(t *testing.T) {
	base := pb.Money{Units: 7, Nanos: 250000000, CurrencyCode: "EUR"}
	if !AreEquals(MultiplySlow(base, 1), base) {
		t.Fatal("multiplying by 1 should return the same amount")
	}
	if !AreEquals(MultiplySlow(base, 0), base) {
		t.Fatal("multiplying by 0 should return the base amount unchanged")
	}
}
