<template>
  <div class="opBox">
    <h1>Operation</h1>
    <select df-opr>
      <option value="">Select an operation</option>
      <option value="add">Addition</option>
      <option value="sub">Substraction</option>
    </select>
    <div>
      <input type="text" placeholder="Number A" id="a" df-aValue>
      <input type="text" placeholder="Number B" id="b" df-bValue>
    </div>
  </div>
</template>

<style scoped>

.opBox {
  width: 100%;
  height: 100%;
}

.opBox div {
  margin-top: 8px;
}

h1 {
  font-family: Arial, Helvetica, sans-serif;
  font-size: medium;
}

select {
  margin: left 4px right 4px;
  width: 100%;
}

div input {
  position: relative;
  font-family: Arial, Helvetica, sans-serif;
  width: 100%;
}

</style>