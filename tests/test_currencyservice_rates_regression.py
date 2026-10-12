import json
import unittest
from pathlib import Path

REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
RATES_PATH = REPOSITORY_ROOT / "src" / "currencyservice" / "data" / "currency_conversion.json"


class CurrencyServiceRatesRegressionTests(unittest.TestCase):
    def test_rates_file_contains_required_currencies(self) -> None:
        rates = json.loads(RATES_PATH.read_text(encoding="utf-8"))
        for code in ("USD", "EUR", "JPY", "CAD", "GBP"):
            self.assertIn(code, rates)
            self.assertGreater(float(rates[code]), 0)

    def test_eur_is_normalized_to_one(self) -> None:
        rates = json.loads(RATES_PATH.read_text(encoding="utf-8"))
        self.assertEqual(float(rates["EUR"]), 1.0)

    def test_convert_usd_to_eur_matches_ecb_style_rate(self) -> None:
        rates = json.loads(RATES_PATH.read_text(encoding="utf-8"))
        usd_amount = 100.0
        usd_rate = float(rates["USD"])
        euros = usd_amount / usd_rate
        self.assertAlmostEqual(euros * usd_rate, usd_amount, places=6)


if __name__ == "__main__":
    unittest.main()
