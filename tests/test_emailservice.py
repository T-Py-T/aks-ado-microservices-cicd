import os
import sys
import unittest

EMAIL_SERVICE_ROOT = os.path.join(os.path.dirname(__file__), "..", "src", "emailservice")
sys.path.insert(0, EMAIL_SERVICE_ROOT)

import grpc  # noqa: E402
import demo_pb2  # noqa: E402
from grpc_health.v1 import health_pb2  # noqa: E402
from email_server import DummyEmailService  # noqa: E402


class _Context:
    def __init__(self) -> None:
        self.code = None
        self.details = None

    def set_code(self, code) -> None:
        self.code = code

    def set_details(self, details) -> None:
        self.details = details


class EmailServiceTests(unittest.TestCase):
    def test_dummy_send_order_confirmation_succeeds(self) -> None:
        service = DummyEmailService()
        context = _Context()
        request = demo_pb2.SendOrderConfirmationRequest(
            email="user@example.com",
            order=demo_pb2.OrderResult(
                order_id="order-123",
                shipping_cost=demo_pb2.Money(currency_code="USD", units=5, nanos=0),
            ),
        )

        response = service.SendOrderConfirmation(request, context)

        self.assertIsInstance(response, demo_pb2.Empty)
        self.assertIsNone(context.code)

    def test_health_check_reports_serving(self) -> None:
        service = DummyEmailService()
        response = service.Check(None, None)
        self.assertEqual(response.status, health_pb2.HealthCheckResponse.SERVING)

    def test_watch_sets_unimplemented_status(self) -> None:
        service = DummyEmailService()
        context = _Context()
        service.Watch(None, context)
        self.assertEqual(context.code, grpc.StatusCode.UNIMPLEMENTED)


if __name__ == "__main__":
    unittest.main()
