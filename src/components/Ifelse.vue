<template>
  <h1>If else block</h1>
  <h2>if:</h2>
  <p>A value is</p>
  <select df-val name="conditions">
    <option value="">Select a condition</option>
    <option value="==">Is equal to</option>
    <option value="!=">Is unequal to</option>
    <option value=">">Is greater than</option>
    <option value="<">Is less than</option>
    <option value=">=">Is greater than or equal to</option>
    <option value="<=">Is less than or equal to</option>
  </select>
  <p>B value</p>
  <h2 id="do">do:</h2>
  <h2 id="else">else:</h2>
</template>

<style scoped>
  select,
  h1,
  h2,
  p {
    font-family: Arial, Helvetica, sans-serif;
  }

  h1 {
    font-size: medium;
  }

  h2,
  p {
    font-size: small;
  }

  #do,
  #else {
    position: relative;
  }

  #do {
    top: -27px;
    left: 176px;
  }

  #else {
    top: -25px;
    left: 166px;
  }

</style>