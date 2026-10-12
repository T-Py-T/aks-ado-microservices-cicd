import os
import sys
import unittest

RECOMMENDATION_SERVICE_ROOT = os.path.join(
    os.path.dirname(__file__), "..", "src", "recommendationservice"
)
sys.path.insert(0, RECOMMENDATION_SERVICE_ROOT)

from recommendations import list_recommendations


class RecommendationLogicTests(unittest.TestCase):
    def test_excludes_requested_products(self) -> None:
        catalog = ["OLJCESPC7Z", "66VCHSJNUP", "1YMWWN1N4O", "L9ECAV7KIM", "2ZYFJ3GM2N"]
        requested = ["OLJCESPC7Z", "66VCHSJNUP"]
        result = list_recommendations(requested, catalog, rng=__import__("random").Random(0))
        self.assertEqual(len(result), 3)
        for product_id in result:
            self.assertNotIn(product_id, requested)

    def test_caps_results_at_five(self) -> None:
        catalog = [f"product-{index}" for index in range(20)]
        result = list_recommendations([], catalog, max_responses=5, rng=__import__("random").Random(1))
        self.assertEqual(len(result), 5)

    def test_returns_empty_when_catalog_is_fully_excluded(self) -> None:
        catalog = ["only-product"]
        result = list_recommendations(["only-product"], catalog)
        self.assertEqual(result, [])


if __name__ == "__main__":
    unittest.main()
