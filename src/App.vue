<script setup>
import Drawflow from 'drawflow'
import "drawflow/dist/drawflow.min.css"
import { shallowRef, ref, h, render, onMounted } from 'vue'
import addComp from './components/add.vue';
import assingComp from './components/assign.vue'
import GeneratorVue from './components/Generator.vue';

const showGenerator = ref(false)
const Vue = { version: 3, h, render };
const editor = shallowRef({})

function newNode(data) {
  editor.value.addNode(data.name, data.in, 1, 0, 0, data.class, {}, data.comp, 'vue')
}

function clicking() {
  console.log('Clicked!')
  showGenerator.value = !showGenerator.value
}

function hola() {
  console.log('Hey')
}

onMounted(() => {
  let id = document.getElementById("drawflow");
  console.log(id)
  editor.value = new Drawflow(id, Vue);
  editor.value.registerNode('assign', assingComp, {}, {})
  editor.value.registerNode('add', addComp, {}, {})
  editor.value.start();
})
</script>

<template>
  <div class="box">
    <div class="left-panel">
      <h3 class="test">Blocks</h3>
      <button>New number</button>
      <button @click="newNode({ name:'Assignation', in:1, class:'Flow', comp:'assign'})">New assignation</button>
      <button @click="newNode({ name:'Addition', in:2, class:'Operation', comp:'add'})">New addition</button>
      <button>substraction</button>
      <button>multiplication</button>
      <button>division</button>
      <button>new function</button>
      <button>add if-else</button>
      <button>add for loop</button>
      <button @click="clicking">generate code</button>
    </div>
    <div id="drawflow"></div>
    <div v-if="showGenerator" class="right-panel">
      <GeneratorVue>Here goes the code</GeneratorVue>
    </div>
  </div>
</template>

<style scoped>
.box {
  position: absolute;
  display: flex;
  height: 100%;
  width: 100%;
  left: 0px;
  top: 0px;
}

.left-panel {
  height: 100%;
  width: 18%;
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

#drawflow {
  width: 100%;
  height: 100%;
  border: 1px solid red;
  text-align: initial;
}

.right-panel {
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
