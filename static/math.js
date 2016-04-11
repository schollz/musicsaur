function meanValues(numbers) {
  // mean of [3, 5, 4, 4, 1, 1, 2, 3] is 2.875
  var total = 0,
    i;
  for (i = 0; i < numbers.length; i += 1) {
    total += numbers[i];
  }
  return total / numbers.length;
}

function medianValues(numbers) {
  // median of [3, 5, 4, 4, 1, 1, 2, 3] = 3
  var median2 = 0,
    numsLen = numbers.length;
  numbers.sort();
  if (numsLen % 2 === 0) { // is even
    // average of two middle numbers
    median2 = (numbers[numsLen / 2 - 1] + numbers[numsLen / 2]) / 2;
  } else { // is odd
    // middle number only
    median2 = numbers[(numsLen - 1) / 2];
  }
  return median2;
}

function modeValues(numbers) {
  // as result can be bimodal or multimodal,
  // the returned result is provided as an array
  // mode of [3, 5, 4, 4, 1, 1, 2, 3] = [1, 3, 4]
  var modes = [],
    count = [],
    i,
    number,
    maxIndex = 0;
  for (i = 0; i < numbers.length; i += 1) {
    number = numbers[i];
    count[number] = (count[number] || 0) + 1;
    if (count[number] > maxIndex) {
      maxIndex = count[number];
    }
  }
  for (i in count)
    if (count.hasOwnProperty(i)) {
      if (count[i] === maxIndex) {
        modes.push(Number(i));
      }
    }
  return modes[0];
}

function range(numbers) {
  // range of [3, 5, 4, 4, 1, 1, 2, 3] is [1, 5]
  numbers.sort();
  return [numbers[0], numbers[numbers.length - 1]];
}

function stdDev(values) {
  var avg = meanValues(values);

  var squareDiffs = values.map(function(value) {
    var diff = value - avg;
    var sqrDiff = diff * diff;
    return sqrDiff;
  });

  var avgSquareDiff = meanValues(squareDiffs);

  var stdDev = Math.sqrt(avgSquareDiff);
  return stdDev;
}

function modulus(val1, val2) {
  return val1 % val2;
}
