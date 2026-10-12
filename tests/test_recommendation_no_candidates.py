import sys
import unittest
from pathlib import Path
from unittest.mock import MagicMock

REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
RECOMMENDATION_SRC = REPOSITORY_ROOT / "src" / "recommendationservice"
sys.path.insert(0, str(RECOMMENDATION_SRC))

import demo_pb2  # noqa: E402
import recommendation_server  # noqa: E402
from recommendation_server import RecommendationService  # noqa: E402


class RecommendationNoCandidatesTests(unittest.TestCase):
    def test_list_recommendations_empty_when_every_product_excluded(self) -> None:
        catalog_stub = MagicMock()
        catalog_stub.ListProducts.return_value = demo_pb2.ListProductsResponse(
            products=[
                demo_pb2.Product(id="0PUK6V6EV0"),
                demo_pb2.Product(id="1YMWWN1N4O"),
            ]
        )
        recommendation_server.product_catalog_stub = catalog_stub

        service = RecommendationService()
        request = demo_pb2.ListRecommendationsRequest(
            product_ids=["0PUK6V6EV0", "1YMWWN1N4O"]
        )
        response = service.ListRecommendations(request, None)

        self.assertEqual(list(response.product_ids), [])
        catalog_stub.ListProducts.assert_called_once()


if __name__ == "__main__":
    unittest.main()
