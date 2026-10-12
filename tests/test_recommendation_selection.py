"""Tests recommendation list filtering logic used by recommendationservice."""

import random
import unittest


def list_recommendations(
    requested_product_ids: list[str],
    catalog_product_ids: list[str],
    *,
    max_responses: int = 5,
    rng: random.Random | None = None,
) -> list[str]:
    """Pure-Python mirror of RecommendationService.ListRecommendations selection."""
    rng = rng or random.Random()
    filtered = list(set(catalog_product_ids) - set(requested_product_ids))
    num_return = min(max_responses, len(filtered))
    if num_return == 0:
        return []
    indices = rng.sample(range(len(filtered)), num_return)
    return [filtered[i] for i in indices]


class RecommendationSelectionTests(unittest.TestCase):
    def test_excludes_requested_products(self) -> None:
        catalog = ["A", "B", "C", "D", "E"]
        requested = ["A", "B"]
        result = list_recommendations(requested, catalog, rng=random.Random(0))
        self.assertTrue(result)
        for product_id in result:
            self.assertNotIn(product_id, requested)

    def test_caps_at_five_recommendations(self) -> None:
        catalog = [f"id-{i}" for i in range(30)]
        result = list_recommendations([], catalog, max_responses=5, rng=random.Random(1))
        self.assertEqual(len(result), 5)

    def test_returns_empty_when_nothing_left(self) -> None:
        catalog = ["only"]
        result = list_recommendations(["only"], catalog)
        self.assertEqual(result, [])

    def test_deduplicates_catalog_ids(self) -> None:
        catalog = ["dup", "dup", "other"]
        result = list_recommendations([], catalog, max_responses=5, rng=random.Random(2))
        self.assertEqual(len(result), len(set(result)))


if __name__ == "__main__":
    unittest.main()
