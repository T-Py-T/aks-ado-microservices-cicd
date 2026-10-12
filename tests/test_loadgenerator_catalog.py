"""Keep load-generator product IDs aligned with the catalog shipped by productcatalogservice."""

import json
import re
import unittest
from pathlib import Path

REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
LOCUSTFILE = REPOSITORY_ROOT / "src" / "loadgenerator" / "locustfile.py"
CATALOG = REPOSITORY_ROOT / "src" / "productcatalogservice" / "products.json"


class LoadGeneratorCatalogTests(unittest.TestCase):
    def test_locust_products_exist_in_catalog(self) -> None:
        locust_source = LOCUSTFILE.read_text()
        match = re.search(r"products\s*=\s*\[(.*?)\]", locust_source, re.DOTALL)
        self.assertIsNotNone(match, "locustfile products list not found")

        locust_ids = re.findall(r"'([^']+)'", match.group(1))
        self.assertGreater(len(locust_ids), 0, "expected at least one locust product id")

        catalog = json.loads(CATALOG.read_text())
        catalog_ids = {item["id"] for item in catalog["products"]}

        missing = sorted(set(locust_ids) - catalog_ids)
        self.assertEqual(missing, [], f"locust products missing from catalog: {missing}")

    def test_catalog_has_unique_ids(self) -> None:
        catalog = json.loads(CATALOG.read_text())
        ids = [item["id"] for item in catalog["products"]]
        self.assertEqual(len(ids), len(set(ids)), "catalog product ids must be unique")


if __name__ == "__main__":
    unittest.main()
