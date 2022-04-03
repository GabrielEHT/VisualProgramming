<template>
  <h1>If else block</h1>
  <h2>if:</h2>
  <div>
    <label for="conditions">value</label>
    <select df-val name="conditions">
      <option value="">Select a condition</option>
      <option value="==">Is equal to</option>
      <option value="!=">Is unequal to</option>
      <option value=">">Is greater than</option>
      <option value="<">Is less than</option>
      <option value=">=">Is greater than or equal to</option>
      <option value="<=">Is less than or equal to</option>
    </select>
    <input type="number" placeholder="value" df-con>
  </div>
  <h2 id="out">else:</h2>
</template>

<style scoped>
  h1 {
    font-size: medium;
  }

  h2 {
    font-size: small;
  }

  #out {
    position: relative;
    right: -42px;
  }

  label {
    font-size: small;
  }

  div {
    display: grid;
  }

</style>