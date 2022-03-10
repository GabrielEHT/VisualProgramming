<script setup>
import Drawflow from 'drawflow'
// Necesario?
import "drawflow/dist/drawflow.min.css"
import { shallowRef, ref, h, render, onMounted } from 'vue'
import nodes from './components/nodes.vue'
import GeneratorVue from './components/Generator.vue';

const showGenerator = ref(false)
const Vue = { version: 3, h, render };
const editor = shallowRef({})
// Faltan nodos
const nodeData = ref([
  {name:'Number', type:'num', class:'Logic', in:0},
  {name:'Assignation', type:'assign', class:'Logic', in:1},
  {name:'Addition', type:'add', class:'Operation', in:2},
  {name:'Substraction', type:'sub', class:'Operation', in:2},
  {name:'Multiplication', type:'mul', class:'Operation', in:2},
  {name:'Division', type:'div', class:'Operation', in:2}
])
var connections = []

function newNode(data) {
  editor.value.addNode(data.name, data.in, 1, 0, 0, data.class, {}, data.type, 'vue')
}

function clicking() {
  console.log(editor.value.getNodeFromId(connections[0].output_id))
  showGenerator.value = !showGenerator.value
}

onMounted(() => {
  let id = document.getElementById("drawflow");
  console.log(id)
  editor.value = new Drawflow(id, Vue);
  for (let i = 0; i < nodeData.value.length; i++) {
    editor.value.registerNode(
      nodeData.value[i].type,
      nodes,
      {type:nodeData.value[i].type},
      {title:() => nodeData.value[i].name}
    )
  }
  editor.value.on('nodeCreated', (id) => {
    console.log('New node:', id)
  })
  editor.value.on('connectionCreated', (data) => {
    connections.push(data)
    console.log(connections)
  })
  editor.value.start();
})
</script>

<template>
  <div class="box">
    <div class="left-panel">
      <h3 class="test">Blocks</h3>
      <ul>
        <li v-for="data in nodeData">
          <button @click="newNode(data)">New {{data.name}}</button>
        </li>
      </ul>
      <button>New function</button>
      <button>New if-else block</button>
      <button>New for loop</button>
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
  border-right: 2px solid black;
  text-align: center;
}

.left-panel * {
  font-size: medium;
  font-family: Arial, Helvetica, sans-serif;
}

#drawflow {
  width: 100%;
  height: 100%;
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

<style>
.drawflow .parent-node .drawflow-node {
    background-color: brown;
    width: auto;
}

.drawflow .parent-node .drawflow-node:hover {
    background-color: lightsalmon;
}

.drawflow .parent-node .drawflow-node .output {
  background-color: aquamarine;
}

.drawflow .parent-node .drawflow-node .output:hover {
  background-color: coral;
}
</style>
