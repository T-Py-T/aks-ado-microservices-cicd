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

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/checkoutservice/genproto"
)

func TestMultiplySlow(t *testing.T) {
	base := mmc(3, 250000000, "USD")

	tests := []struct {
		name string
		n    uint32
		want pb.Money
	}{
		{"multiply by one", 1, base},
		{"multiply by zero uses loop guard", 0, base},
		{"multiply by three", 3, mmc(9, 750000000, "USD")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MultiplySlow(base, tt.n)
			if !AreEquals(got, tt.want) {
				t.Fatalf("MultiplySlow(%v, %d) = %v, want %v", base, tt.n, got, tt.want)
			}
		})
	}
}
