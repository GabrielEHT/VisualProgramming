<template>
  <div class="opBox">
    <h1>Operation</h1>
    <select df-val>
      <option value="">Select an operation</option>
      <option value="+">Addition</option>
      <option value="-">Substraction</option>
      <option value="*">Multiplication</option>
      <option value="/">Division</option>
    </select>
    <h2>A value</h2>
    <h2>B value</h2>
  </div>
</template>

<style scoped>

.opBox {
  font-family: Arial, Helvetica, sans-serif;
  width: 100%;
  height: 100%;
}

.opBox div {
  display: grid;
  margin-top: 8px;
}

h1 {
  font-size: medium;
}

h2 {
  font-size: small;
}

div input {
  position: relative;
}

</style>