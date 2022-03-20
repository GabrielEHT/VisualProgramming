<script setup>
import Drawflow from 'drawflow'
import "drawflow/dist/drawflow.min.css"
import { shallowRef, ref, h, render, onMounted } from 'vue'
import * as components from './components/nodes.js'

const code = ref(null)
const editor = shallowRef({})
// Faltan nodos
const nodeData = ref([
  {name:'number', type:'num', class:'Value', vars:{'val':''}},
  {name:'assignation', type:'assign', class:'Value', in:1, vars:{'val':''}},
  {name:'operation', type:'operations', class:'Operation', in:2, vars:{'val':'','aNum':'', 'bNum':''}},
])
var nodeList = [];

function getLastNode() {
  var lastNode;
  for (let node of nodeList) {
    if (node.output == false && lastNode == undefined) {
      lastNode = node
    } else if (node.output == false) {
      return 'err'
    }
  }
  return lastNode
}

function renderCode() {
  if (nodeList.length) {
    var lastNode = getLastNode();
    if (lastNode != 'err') {
      console.log(nodeList)
      code.value.data = createScript([generateCode(lastNode)])
    } else {
      alert('You left some unconnected nodes!')
    }
  } else {
    alert('You haven\'t created any nodes!');
  }
}

function createScript(data) {
  let script = new Blob(data, {type:"text/plain;charset=utf-8"})
  let scriptUrl = window.URL.createObjectURL(script)
  return scriptUrl
}

function generateCode(lastNode) {
  var codeLine;
  var formated = false;
  var nodeInfo;
  var inputs;

  if (typeof lastNode == 'string') {
    nodeInfo = editor.value.getNodeFromId(lastNode);
    inputs = []
  } else {
    nodeInfo = editor.value.getNodeFromId(lastNode.id);
    inputs = lastNode.input_from.split('')
    if (inputs.length > 1) {
      inputs.splice(1, 1)
    }
  }

  if (nodeInfo.name == 'assignation') {
    codeLine = nodeInfo.data.val + ' = '
  } else if (nodeInfo.name == 'operation') {
    let symbol;
    if (nodeInfo.data.val == 'add') {
      symbol = '+'
    } else if (nodeInfo.data.val == 'sub') {
      symbol = '-'
    } else if (nodeInfo.data.val == 'sub') {
      symbol = '*'
    } else if (nodeInfo.data.val == 'sub') {
      symbol = '/'
    }
    codeLine = `${generateCode(inputs[0])} + ${generateCode(inputs[1])}`
    formated = true;
  } else {
    codeLine = nodeInfo.data.val
  }

  if (formated == false) {
    for (let input of inputs) {
      for (let node of nodeList) {
        if (node.id == input) {
          codeLine += generateCode(node)
        }
      }
    }
  }

  return codeLine;
}

function sendExecTree() {
  var lastNode = getLastNode()
  if (lastNode != 'err') {
    let execTree = generateExecTree(lastNode)
    console.log(execTree)
    sendData(execTree)
  } else {
    alert('Unexpected error happened')
  }
}

function generateExecTree(lastNode) {
  let nodeInfo = editor.value.getNodeFromId(lastNode.id);
  var nodeExec = {[nodeInfo.name + ':' + nodeInfo.data.val]:[]};
  let inputs = lastNode.input_from.split('')
  if (inputs.length > 1) {
    inputs.splice(1, 1)
  }
  for (let input of inputs) {
    for (let node of nodeList) {
      if (node.id == input && node.input_from != 'none') {
        nodeExec[nodeInfo.name + ':' + nodeInfo.data.val].push(generateExecTree(node))
      } else if (node.id == input) {
        nodeExec[nodeInfo.name + ':' + nodeInfo.data.val].push(editor.value.getNodeFromId(input).data.val)
      }
    }
  }
  return nodeExec;
}

async function sendData(data) {
  const http = new XMLHttpRequest()
  http.open('POST', 'http://localhost:8080/', true)
  http.addEventListener('load', () => {
    let r = http.response
    if (typeof r == "string") {
      console.log(r)
    } else {
      console.log(JSON.parse(r))
    }
  })
  await http.send(JSON.stringify(data))
}


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

function addConnection(output_id, input_id) {
  for (let node of nodeList) {
    if (output_id == node.id) {
      node.output = true
    } else if (input_id == node.id) {
      if (node.input_from == 'none') {
        node.input_from = output_id
      } else {
        node.input_from += ':' + output_id
      }
    }
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
    nodeList.push({'id':id, 'output':false, 'input_from':'none'})
  })

  // Keeps track of deleted nodes
  editor.value.on('nodeRemoved', (id) => {
    nodeList = nodeList.filter(node => node.id!=id)
    console.log('Removed id:', id, nodeList)
  })

  //Checks if the created connection is valid
  editor.value.on('connectionCreated', (data) => {
    let input = editor.value.getNodeFromId(data.input_id);
    let output = editor.value.getNodeFromId(data.output_id);

    if (input.class == 'Operation') {
      if (input.data.val != '' && output.data.val != '') {
        addConnection(data.output_id, data.input_id)
      } else {
        editor.value.removeSingleConnection(data.output_id, data.input_id, data.output_class, data.input_class)
        alert('Something is not defined')
      }
    } else if (input.class == 'Value') {
      if (input.data.val) {
        addConnection(data.output_id, data.input_id)
      } else{
        editor.value.removeSingleConnection(data.output_id, data.input_id, data.output_class, data.input_class)
        alert('Missing name in assignation node!')
      }
    }
  })

  //Updates connection info on removed connection event
  editor.value.on('connectionRemoved', (data) => {
    for (let node of nodeList) {
      if (data.output_id == node.id) {
        node.output = false
      } else if (data.input_id == node.id) {
        let inputs = node.input_from.split('')
        if (inputs.length == 1) {
          node.input_from = 'none'
        } else if (inputs.indexOf(data.output_id) < 1) {
          node.input_from = inputs[2]
        } else {
          node.input_from = inputs[0]
        }
      }
    }
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
    </div>
    <div id="drawflow"></div>
    <div class="right-panel">
      <button @click="renderCode" id="generate">generate code</button>
      <div id="list"><object ref="code" width=200 height=400></object></div>
      <button @click="sendExecTree()" id="execute">Execute code</button>
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
  background-color: lightcyan;
  width: 30%;
  height: 100%;
}

.right-panel button {
  position: absolute;
  right: 0px;
}

#generate{
  top: 0px;
}

#execute {
  bottom: 0px;
}

.right-panel div {
  position: relative;
  top: 30px;
}

#list object {
  font: Arial;
  font-size: 14px;
  background-color: white;
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

.drawflow .drawflow-node.Operation .output {
  top: 49px;
}

.drawflow .drawflow-node.Value .input,
.drawflow .drawflow-node.Value .output{
  top: 22px;
}

</style>
