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

import (
	"testing"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/genproto"
)

// Regression: IsPositive/IsNegative must not report sign for invalid Money when units == 0.
func TestIsPositiveNegativeRejectInvalidZeroUnits(t *testing.T) {
	invalidPositive := pb.Money{Units: 0, Nanos: 2_000_000_000}
	if IsValid(invalidPositive) {
		t.Fatal("expected invalid money for out-of-range positive nanos with zero units")
	}
	if IsPositive(invalidPositive) {
		t.Errorf("IsPositive must be false when IsValid is false (got true for %v)", invalidPositive)
	}

	invalidNegative := pb.Money{Units: 0, Nanos: -2_000_000_000}
	if IsValid(invalidNegative) {
		t.Fatal("expected invalid money for out-of-range negative nanos with zero units")
	}
	if IsNegative(invalidNegative) {
		t.Errorf("IsNegative must be false when IsValid is false (got true for %v)", invalidNegative)
	}
}
