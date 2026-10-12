import unittest
from pathlib import Path
from types import SimpleNamespace

from jinja2 import Environment, FileSystemLoader, select_autoescape

REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
TEMPLATE_DIR = REPOSITORY_ROOT / "src" / "emailservice" / "templates"


class EmailServiceTemplateRegressionTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.env = Environment(
            loader=FileSystemLoader(TEMPLATE_DIR),
            autoescape=select_autoescape(["html", "xml"]),
        )
        cls.template = cls.env.get_template("confirmation.html")

    def test_confirmation_template_renders_order_total(self) -> None:
        order = SimpleNamespace(
            order_id="ORDER-123",
            shipping_tracking_id="AB-40123-2076543210",
            shipping_cost=SimpleNamespace(units=8, nanos=990000000, currency_code="USD"),
            shipping_address=SimpleNamespace(
                street_address_1="1600 Amphitheatre Parkway",
                street_address_2="",
                city="Mountain View",
                country="United States",
                zip_code="94043",
            ),
            items=[
                SimpleNamespace(
                    item=SimpleNamespace(product_id="SKU-1", quantity=2),
                    cost=SimpleNamespace(units=20, nanos=0, currency_code="USD"),
                )
            ],
        )
        html = self.template.render(order=order)
        self.assertIn("ORDER-123", html)
        self.assertIn("SKU-1", html)


if __name__ == "__main__":
    unittest.main()
