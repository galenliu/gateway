
// http://madreporite.com/insteon/ramprate.htm
const RATES = [
  2.0,
  6 * 60,
  4.5 * 60,
  3.5 * 60,
  2.5 * 60,
  1.5 * 60,
  47,
  38.5,
  32,
  28,
  23.5,
  19,
  6.5,
  2,
  0.3,
  0.1
];

export function rateForDuration(duration) {
  let closestIndex;
  let closestDistance;
  for (let i = 1; i < RATES.length; i++) {
    // The command we'll use this in uses half of the published rate.
    const rate = RATES[i] / 2;

    const distance = Math.abs(rate - duration);
    if (!closestIndex) {
      closestIndex = i;
      closestDistance = distance;
      continue;
    }

    if (distance < closestDistance) {
      closestIndex = i;
      closestDistance = distance;
    }
  }

  return closestIndex;
}
