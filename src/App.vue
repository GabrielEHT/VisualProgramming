<script setup>
import Drawflow from 'drawflow'
import "drawflow/dist/drawflow.min.css"
import { shallowRef, ref, h, render, onMounted } from 'vue'
import * as components from './components/nodes.vue'

const showGenerator = ref(false)
const editor = shallowRef({})
// Faltan nodos
const nodeData = ref([
  {name:'number', type:'num', class:'Value', vars:{'value':''}},
  {name:'assignation', type:'assign', class:'Value', in:1, vars:{'value':''}},
  {name:'operation', type:'operations', class:'Operation', in:2, vars:{'opr':'','aValue':'', 'bValue':''}},
])
var nodeList = [];
var execTree = {};

function addNode(data) {
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
  let nodeCount = nodeList.length;
  if (nodeCount) {
    console.log('Number of nodes:', nodeCount);
    let createdNodes = [];
    for (let i = 1; i <= nodeCount; i++) {
      createdNodes.push(editor.value.getNodeFromId(i))
    }
    console.log(createdNodes)
    console.log(execTree)
    console.log(nodeList)
    // Make a sprting system
    showGenerator.value = true
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
  for (let node of nodeData.value) {
    if (node.class == 'Operation') {
      editor.value.registerNode(
        node.type,
        components.operations,
        {},
        {}
      )
    } else if (node.class == 'Value') {
      editor.value.registerNode(
        node.type,
        components.datatypes,
        {'type':node.type},
        {}
      )
    } else {
      editor.value.registerNode(
        node.type,
        components.datatypes,
        {},
        {}
      )
    }
  }

  // Keeps track of new created nodes
  editor.value.on('nodeCreated', (id) => {
    console.log('New node:', id);
    console.log(editor.value.getNodeFromId(id));
    nodeList.push({'id':id, 'connected':false})
  })

  // Keeps track of deleted nodes
  editor.value.on('nodeRemoved', (id) => {
    nodeList = nodeList.filter(node => node.id!=id)
    console.log('Removed id:', id, nodeList)
  })

  editor.value.on('connectionCreated', (data) => {
    let input = editor.value.getNodeFromId(data.input_id);
    let output = editor.value.getNodeFromId(data.output_id);

    if (input.class == 'Operation') {
      if (input.data.opr != '' && output.data.value != '') {
        let id = 'Operation' + data.input_id;
        let opr = input.data.opr;
        let val = output.data.value;
        execTree[id]? console.log(execTree) : execTree[id] = {[opr]:['X', 'X']}
        data.input_class=='input_1'? execTree[id][opr][0] = val : execTree[id][opr][1] = val;
        for (let node of nodeList) {
          if (node.id == data.output_id) {
            node.connected = true
          }
        }
      } else {
        editor.value.removeSingleConnection(data.output_id, data.input_id, data.output_class, data.input_class)
        alert('Something is not defined')
      }
    } else if (input.class == 'Value') {
      if (input.data.value) {
        if (output.class == 'Operation') {
          let oprid = 'Operation' + data.output_id;
          let id = 'Assignation' + data.input_id;
          execTree[oprid][id] = input.data.value;
        } else {
          // covertir en Array?
          let id = 'Assignation' + data.input_id;
          let name = input.data.value;
          let val = output.data.value;
          execTree[id] = {[name]:val};
        }
      } else{
        editor.value.removeSingleConnection(data.output_id, data.input_id, data.output_class, data.input_class)
        alert('Missing name in assignation node!')
      }
    } else {
      // Algo?
    }
  })

  editor.value.on('connectionRemoved', (data) => {
    // TODO
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
          <button @click="addNode(data)">New {{data.name}}</button>
        </li>
      </ul>
      <button>New function</button>
      <button>New if-else block</button>
      <button>New for loop</button>
      <button @click="runGenerator">generate code</button>
    </div>
    <div id="drawflow"></div>
    <div v-if="showGenerator" class="right-panel">
      <button @click="showGenerator=false">X</button>
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
    background-color: rgb(212, 216, 218);
    width: auto;
}

.drawflow .parent-node .drawflow-node.selected {
    background-color: rgb(190, 207, 216);
}

.drawflow .parent-node .drawflow-node .output {
  background-color: aquamarine;
}

.drawflow .parent-node .drawflow-node .output:hover {
  background-color: coral;
}

.drawflow .drawflow-node.Operation .input {
  top: 38px;
}

</style>
