# Copyright 2018 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import random


def list_recommendations(product_ids, catalog_product_ids, max_responses=5, rng=None):
  """Return up to max_responses product ids not already in the request."""
  if rng is None:
    rng = random

  filtered_products = list(set(catalog_product_ids) - set(product_ids))
  num_products = len(filtered_products)
  num_return = min(max_responses, num_products)
  if num_return == 0:
    return []

  indices = rng.sample(range(num_products), num_return)
  return [filtered_products[i] for i in indices]
