/*
 * Copyright 2018 Google LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

'use strict';

const path = require('path');

const DEFAULT_DATA_PATH = path.join(__dirname, './data/currency_conversion.json');

/**
 * Helper function that handles decimal/fractional carrying
 */
function carry (amount) {
  const fractionSize = Math.pow(10, 9);
  amount.nanos += (amount.units % 1) * fractionSize;
  amount.units = Math.floor(amount.units) + Math.floor(amount.nanos / fractionSize);
  amount.nanos = amount.nanos % fractionSize;
  return amount;
}

function loadCurrencyData (dataPath = DEFAULT_DATA_PATH) {
  return require(dataPath);
}

function getSupportedCurrencyCodes (dataPath = DEFAULT_DATA_PATH) {
  return Object.keys(loadCurrencyData(dataPath));
}

/**
 * Converts a money amount from one currency to another using ECB rate table.
 */
function convertAmount (from, toCode, dataPath = DEFAULT_DATA_PATH) {
  const data = loadCurrencyData(dataPath);
  const fromRate = data[from.currency_code];
  const toRate = data[toCode];

  if (fromRate === undefined || toRate === undefined) {
    throw new Error(`unsupported currency conversion: ${from.currency_code} -> ${toCode}`);
  }

  const euros = carry({
    units: from.units / fromRate,
    nanos: from.nanos / fromRate
  });

  euros.nanos = Math.round(euros.nanos);

  const result = carry({
    units: euros.units * toRate,
    nanos: euros.nanos * toRate
  });

  result.units = Math.floor(result.units);
  result.nanos = Math.floor(result.nanos);
  result.currency_code = toCode;

  return result;
}

module.exports = {
  carry,
  loadCurrencyData,
  getSupportedCurrencyCodes,
  convertAmount,
};
