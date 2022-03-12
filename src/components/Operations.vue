<script setup>
import { ref } from 'vue'

const opType = ref({
  opName:['Addition', 'Sub', 'Mul', 'Div'],
  opSymbol: ['+', '-', '*', '/']
});
const actOp = ref(0);

</script>

<template>
  <div class="opBox">
    <h1>Operation</h1>
    <button @click="actOp==3? actOp=0:actOp++">{{opType.opName[actOp]}}</button>
    <div>
      <input type="text" placeholder="A" df-aValue>
      <h2>{{opType.opSymbol[actOp]}}</h2>
      <input type="text" placeholder="B" df-bValue>
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

button {
  width: 100%;
  margin: left 4px right 4px;
}

div input {
  position: relative;
  font-family: Arial, Helvetica, sans-serif;
  left: -2px;
}

div h2 {
  font-family: Arial, Helvetica, sans-serif;
  font-size: large;
  padding-left: 50%;
}

div h2+input {
  font-family: Arial, Helvetica, sans-serif;
}
</style>