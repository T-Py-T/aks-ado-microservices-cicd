// Copyright 2023 Google LLC
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

func TestSearchProductsCaseInsensitive(t *testing.T) {
	catalog := &productCatalog{
		catalog: pb.ListProductsResponse{
			Products: []*pb.Product{
				{Id: "1", Name: "Vintage Camera", Description: "film"},
				{Id: "2", Name: "Modern Lens", Description: "CAMERA mount"},
			},
		},
	}

	resp, err := catalog.SearchProducts(context.Background(), &pb.SearchProductsRequest{Query: "CaMeRa"})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("got %d results; want 2 matching name and description", len(resp.Results))
	}
}

func TestSearchProductsNoMatch(t *testing.T) {
	catalog := &productCatalog{
		catalog: pb.ListProductsResponse{
			Products: []*pb.Product{
				{Id: "1", Name: "Tea Cup", Description: "ceramic"},
			},
		},
	}

	resp, err := catalog.SearchProducts(context.Background(), &pb.SearchProductsRequest{Query: "camera"})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Results) != 0 {
		t.Fatalf("expected no results, got %d", len(resp.Results))
	}
}

func TestSearchProductsEmptyQueryMatchesAll(t *testing.T) {
	catalog := &productCatalog{
		catalog: pb.ListProductsResponse{
			Products: []*pb.Product{
				{Id: "1", Name: "A"},
				{Id: "2", Name: "B"},
			},
		},
	}

	resp, err := catalog.SearchProducts(context.Background(), &pb.SearchProductsRequest{Query: ""})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("empty query should match every product; got %d results", len(resp.Results))
	}
}
