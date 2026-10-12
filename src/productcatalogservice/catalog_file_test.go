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
	"os"
	"path/filepath"
	"testing"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/productcatalogservice/genproto"
)

func TestReadCatalogFileLoadsBundledProducts(t *testing.T) {
	catalogPath := filepath.Join("products.json")
	if _, err := os.Stat(catalogPath); err != nil {
		t.Fatalf("products.json must exist in service directory: %v", err)
	}

	var catalog pb.ListProductsResponse
	if err := readCatalogFile(&catalog); err != nil {
		t.Fatalf("readCatalogFile failed: %v", err)
	}
	if len(catalog.Products) == 0 {
		t.Fatal("expected at least one product in catalog")
	}

	ids := make(map[string]struct{})
	for _, product := range catalog.Products {
		if product.Id == "" {
			t.Fatal("catalog product missing id")
		}
		if _, seen := ids[product.Id]; seen {
			t.Fatalf("duplicate product id %q", product.Id)
		}
		ids[product.Id] = struct{}{}
		if product.PriceUsd == nil || product.PriceUsd.CurrencyCode != "USD" {
			t.Fatalf("product %q has unexpected price", product.Id)
		}
	}
}

func TestReadCatalogFileRejectsInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	badFile := filepath.Join(dir, "products.json")
	if err := os.WriteFile(badFile, []byte("{not-json"), 0o644); err != nil {
		t.Fatalf("write temp catalog: %v", err)
	}

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		_ = os.Chdir(origWD)
	}()

	var catalog pb.ListProductsResponse
	if err := readCatalogFile(&catalog); err == nil {
		t.Fatal("expected error for invalid catalog JSON")
	}
}
