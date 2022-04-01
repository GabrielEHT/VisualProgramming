<template>
  <h1>If else block</h1>
  <h2>if:</h2>
  <div>
    <label for="conditions">value</label>
    <select df-val name="conditions">
      <option value="">Select a condition</option>
      <option value="==">Is equal than</option>
      <option value="!=">Is different than</option>
      <option value="<">Is minus than</option>
      <option value="<=">Is minus or equal than</option>
      <option value=">">Is more than</option>
      <option value=">=">Is more or equal than</option>
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