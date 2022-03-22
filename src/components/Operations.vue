<template>
  <div class="opBox">
    <h1>Operation</h1>
    <select df-val>
      <option value="">Select an operation</option>
      <option value="add">Addition</option>
      <option value="sub">Substraction</option>
      <option value="mul">Multiplication</option>
      <option value="div">Division</option>
    </select>
    <div>
      <input type="text" placeholder="Number A" id="a" df-aNum>
      <input type="text" placeholder="Number B" id="b" df-bNum>
    </div>
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

select {
  font-family: Arial, Helvetica, sans-serif;
}

div input {
  position: relative;
  font-family: Arial, Helvetica, sans-serif;
}

</style>