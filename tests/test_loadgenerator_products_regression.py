import ast
import unittest
from pathlib import Path

REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
LOCUSTFILE = REPOSITORY_ROOT / "src" / "loadgenerator" / "locustfile.py"


class LoadGeneratorProductsRegressionTests(unittest.TestCase):
    def test_product_catalog_is_unique_and_non_empty(self) -> None:
        module = ast.parse(LOCUSTFILE.read_text(encoding="utf-8"))
        products = None
        for node in module.body:
            if isinstance(node, ast.Assign):
                for target in node.targets:
                    if isinstance(target, ast.Name) and target.id == "products":
                        products = ast.literal_eval(node.value)
        self.assertIsNotNone(products)
        self.assertGreater(len(products), 0)
        self.assertEqual(len(products), len(set(products)))
        for product_id in products:
            self.assertRegex(product_id, r"^[A-Z0-9]+$")


if __name__ == "__main__":
    unittest.main()
