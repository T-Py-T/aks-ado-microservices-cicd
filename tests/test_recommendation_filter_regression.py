import random
import unittest


def list_recommendation_ids(catalog_ids: list[str], requested_ids: list[str], max_responses: int = 5) -> list[str]:
    """Pure reimplementation of RecommendationService.ListRecommendations selection logic."""
    filtered_products = list(set(catalog_ids) - set(requested_ids))
    num_products = len(filtered_products)
    num_return = min(max_responses, num_products)
    indices = random.sample(range(num_products), num_return)
    return [filtered_products[i] for i in indices]


class RecommendationFilterRegressionTests(unittest.TestCase):
    def test_excludes_requested_products(self) -> None:
        catalog = ["a", "b", "c", "d", "e", "f"]
        random.seed(0)
        result = list_recommendation_ids(catalog, ["a", "c"], max_responses=5)
        self.assertTrue(result)
        self.assertTrue(set(result).issubset({"b", "d", "e", "f"}))
        self.assertEqual(len(result), 4)

    def test_caps_at_max_responses(self) -> None:
        catalog = [str(i) for i in range(20)]
        random.seed(1)
        result = list_recommendation_ids(catalog, [], max_responses=5)
        self.assertEqual(len(result), 5)

    def test_empty_catalog_returns_empty(self) -> None:
        random.seed(2)
        result = list_recommendation_ids([], ["x"], max_responses=5)
        self.assertEqual(result, [])


if __name__ == "__main__":
    unittest.main()
