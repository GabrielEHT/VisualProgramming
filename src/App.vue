<script setup>
import GeneratorVue from './components/Generator.vue';
import DrawflowVue from './components/Drawflow.vue';
import { ref } from 'vue'

const showGenerator = ref(false)

function clicking() {
  console.log('Clicked!')
  showGenerator.value = !showGenerator.value
}
</script>

<template>
  <div class="left-panel">
    <h3 class="test">Blocks</h3>
    <button @click="clicking">generate code</button>
  </div>
  <div class="central-panel">
    <DrawflowVue/>
  </div>
  <div v-if="showGenerator" class="right-panel">
    <GeneratorVue>Here goes the code</GeneratorVue>
  </div>
</template>

<style scoped>
.left-panel {
  position: absolute;
  height: 100%;
  width: 15%;
  top: 0px;
  left: 0px;
  background: rgb(218, 230, 233);
  border-right: 2px;
  border-right-style: solid;
  border-right-color: black;
}

.left-panel * {
  font-size: medium;
  font-family: Arial, Helvetica, sans-serif;
  left: 15%;
}

.central-panel {
  position: relative;
  left: 15%;
  background-color: blue;
}

.right-panel {
  position: absolute;
  right: 0px;
  top: 0px;
  height: 100%;
  background-color: lightcyan;
  border-left: 1px;
  border-left-style: solid;
  border-left-color: black;
  width: 15%;
  height: 100%;
}
</style>
