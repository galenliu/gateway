import { readFileSync } from 'fs';
import { join } from 'path';

const database = JSON.parse(readFileSync(join(__dirname, 'devices.json')));

export function findProductInfo(category, subcategory) {
  const categoryInfo = database[category];
  if (!categoryInfo) {
    return null;
  }

  let productInfo = categoryInfo[subcategory];
  if (!productInfo) {
    productInfo = categoryInfo[0];
  }

  return productInfo;
}
