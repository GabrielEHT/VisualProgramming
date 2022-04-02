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
  width: 100%;
  height: 100%;
}

.opBox div {
  display: grid;
  margin-top: 8px;
}

h1 {
  font-family: Arial, Helvetica, sans-serif;
  font-size: medium;
}

h2 {
  font-family: Arial, Helvetica, sans-serif;
  font-size: small;
}

select {
  font-family: Arial, Helvetica, sans-serif;
}

div input {
  position: relative;
  font-family: Arial, Helvetica, sans-serif;
}

</style>