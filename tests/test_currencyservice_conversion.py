"""Unit tests for currencyservice conversion math (mirrors server.js _carry / convert)."""

import json
import math
import unittest
from pathlib import Path

REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
RATES_PATH = REPOSITORY_ROOT / "src" / "currencyservice" / "data" / "currency_conversion.json"


def carry(amount: dict) -> dict:
    fraction_size = 10**9
    amount = dict(amount)
    amount["nanos"] = amount.get("nanos", 0) + (amount["units"] % 1) * fraction_size
    amount["units"] = math.floor(amount["units"]) + math.floor(amount["nanos"] / fraction_size)
    amount["nanos"] = amount["nanos"] % fraction_size
    return amount


def convert_amount(from_amount: dict, to_code: str, rates: dict) -> dict:
    from_currency = from_amount["currency_code"]
    euros = carry(
        {
            "units": from_amount["units"] / rates[from_currency],
            "nanos": from_amount["nanos"] / rates[from_currency],
        }
    )
    euros["nanos"] = round(euros["nanos"])
    result = carry(
        {
            "units": euros["units"] * rates[to_code],
            "nanos": euros["nanos"] * rates[to_code],
        }
    )
    result["units"] = math.floor(result["units"])
    result["nanos"] = math.floor(result["nanos"])
    result["currency_code"] = to_code
    return result


class CurrencyConversionTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        raw = json.loads(RATES_PATH.read_text())
        cls.rates = {code: float(value) for code, value in raw.items()}

    def test_rates_include_core_currencies(self) -> None:
        for code in ("EUR", "USD", "GBP", "JPY"):
            self.assertIn(code, self.rates)
            self.assertGreater(self.rates[code], 0)

    def test_carry_normalizes_fractional_units(self) -> None:
        amount = carry({"units": 1.25, "nanos": 0})
        self.assertEqual(amount["units"], 1)
        self.assertEqual(amount["nanos"], 250_000_000)

    def test_convert_eur_to_usd(self) -> None:
        result = convert_amount(
            {"currency_code": "EUR", "units": 10, "nanos": 0},
            "USD",
            self.rates,
        )
        self.assertEqual(result["currency_code"], "USD")
        self.assertGreater(result["units"], 0)

    def test_convert_usd_round_trip_is_close(self) -> None:
        original = {"currency_code": "USD", "units": 25, "nanos": 500_000_000}
        eur = convert_amount(original, "EUR", self.rates)
        back = convert_amount(eur, "USD", self.rates)
        self.assertEqual(back["currency_code"], "USD")
        self.assertEqual(back["units"], original["units"])
        self.assertAlmostEqual(back["nanos"], original["nanos"], delta=2_000_000)


if __name__ == "__main__":
    unittest.main()
