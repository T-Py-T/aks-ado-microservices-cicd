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

import (
	"context"
	"testing"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/productcatalogservice/genproto"
)

func TestSearchProductsEmptyQueryReturnsNoResults(t *testing.T) {
	catalog := &productCatalog{
		catalog: pb.ListProductsResponse{
			Products: []*pb.Product{
				{Id: "p1", Name: "Widget", Description: "A widget"},
				{Id: "p2", Name: "Gadget", Description: "A gadget"},
			},
		},
	}

	resp, err := catalog.SearchProducts(context.Background(), &pb.SearchProductsRequest{Query: ""})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Results) != 0 {
		t.Fatalf("empty query returned %d products, want 0", len(resp.Results))
	}

	resp, err = catalog.SearchProducts(context.Background(), &pb.SearchProductsRequest{Query: "   "})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Results) != 0 {
		t.Fatalf("whitespace-only query returned %d products, want 0", len(resp.Results))
	}
}
