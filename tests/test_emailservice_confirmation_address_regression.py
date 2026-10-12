import unittest
from pathlib import Path
from types import SimpleNamespace

from jinja2 import Environment, FileSystemLoader, select_autoescape

REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
TEMPLATE_DIR = REPOSITORY_ROOT / "src" / "emailservice" / "templates"


class EmailConfirmationAddressRegressionTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.env = Environment(
            loader=FileSystemLoader(TEMPLATE_DIR),
            autoescape=select_autoescape(["html", "xml"]),
        )
        cls.template = cls.env.get_template("confirmation.html")

    def test_confirmation_template_renders_proto_address_fields(self) -> None:
        """Address in demo.proto uses street_address and state, not street_address_1/2."""
        order = SimpleNamespace(
            order_id="ORDER-456",
            shipping_tracking_id="AB-40123-2076543210",
            shipping_cost=SimpleNamespace(units=8, nanos=990000000, currency_code="USD"),
            shipping_address=SimpleNamespace(
                street_address="1600 Amphitheatre Parkway",
                city="Mountain View",
                state="CA",
                country="United States",
                zip_code=94043,
            ),
            items=[
                SimpleNamespace(
                    item=SimpleNamespace(product_id="SKU-1", quantity=1),
                    cost=SimpleNamespace(units=20, nanos=0, currency_code="USD"),
                )
            ],
        )
        html = self.template.render(order=order)
        self.assertIn("1600 Amphitheatre Parkway", html)
        self.assertIn("Mountain View", html)
        self.assertIn("CA", html)
        self.assertIn("94043", html)
        self.assertNotIn("street_address_1", html)


if __name__ == "__main__":
    unittest.main()
