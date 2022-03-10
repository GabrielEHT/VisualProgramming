<script setup>
import Drawflow from 'drawflow'
// Necesario?
import "drawflow/dist/drawflow.min.css"
import { shallowRef, ref, h, render, onMounted } from 'vue'
// Separar en distintos componentes?
import nodes from './components/nodes.vue'

const showGenerator = ref(false)
const editor = shallowRef({})
// Faltan nodos
const nodeData = ref([
  {name:'Number', type:'num', class:'Logic', vars:{'value':''}},
  {name:'Assignation', type:'assign', class:'Logic', in:1, vars:{'value':''}},
  {name:'Addition', type:'add', class:'Operation', in:2},
  {name:'Substraction', type:'sub', class:'Operation', in:2},
  {name:'Multiplication', type:'mul', class:'Operation', in:2},
  {name:'Division', type:'div', class:'Operation', in:2}
])
var nodeCount;

// Cambiar nombre?
function newNode(data) {
  editor.value.addNode(
    data.name,
    data.in? data.in : 0,
    data.out? data.out : 1,
    0,
    0,
    data.class,
    data.vars? data.vars : {},
    data.type,
    'vue')
}

function runGenerator() {
  if (nodeCount) {
    console.log(nodeCount)
    let nodesData = []
    for (let i = 1; i <= nodeCount; i++) {
      nodesData.push(editor.value.getNodeFromId(i))
    }
    console.log(nodesData)
    nodesData.sort((node) => {
      
    })
    showGenerator.value = !showGenerator.value
  } else {
    alert('You haven\'t defined any nodes');
  }
}

onMounted(() => {
  // Initialices Drawflow
  let id = document.getElementById("drawflow");
  let Vue = { version: 3, h, render };
  editor.value = new Drawflow(id, Vue);

  // Registers all nodes
  for (let i = 0; i < nodeData.value.length; i++) {
    editor.value.registerNode(
      nodeData.value[i].type,
      nodes,
      {type:nodeData.value[i].type},
      {title:() => nodeData.value[i].name}
    )
  }

  // Keeps track of new created nodes
  editor.value.on('nodeCreated', (id) => {
    console.log('New node:', id);
    console.log(editor.value.getNodeFromId(id));
    nodeCount? nodeCount++ : nodeCount = 1;
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
      <button @click="runGenerator">generate code</button>
    </div>
    <div id="drawflow"></div>
    <div v-if="showGenerator" class="right-panel">
      <button>X</button>
      <div>
        <p>Generating...</p>
      </div>
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
  padding-left: 15px;
}

.left-panel * {
  font-size: medium;
  font-family: Arial, Helvetica, sans-serif;
}

#drawflow {
  width: 100%;
  height: 100%;
  text-align: initial;
  border: 2px solid black;
}

.right-panel {
  right: 0px;
  top: 0px;
  height: 100%;
  background-color: lightcyan;
  width: 18%;
  height: 100%;
}

.right-panel button {
  position: absolute;
  top: 0px;
  right: 0px;
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
