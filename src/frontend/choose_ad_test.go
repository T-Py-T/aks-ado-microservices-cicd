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
	"testing"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"
)

func TestPickRandomAd_emptyListReturnsNil(t *testing.T) {
	if ad := pickRandomAd(nil); ad != nil {
		t.Fatalf("pickRandomAd(nil) = %v, want nil", ad)
	}
	if ad := pickRandomAd([]*pb.Ad{}); ad != nil {
		t.Fatalf("pickRandomAd([]) = %v, want nil", ad)
	}
}

func TestPickRandomAd_nonEmptyReturnsMember(t *testing.T) {
	ads := []*pb.Ad{
		{RedirectUrl: "https://example.com/a"},
		{RedirectUrl: "https://example.com/b"},
	}
	ad := pickRandomAd(ads)
	if ad == nil {
		t.Fatal("pickRandomAd expected an ad, got nil")
	}
	found := false
	for _, candidate := range ads {
		if candidate == ad {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("pickRandomAd returned unexpected ad %v", ad)
	}
}
