<script setup>
import Drawflow from 'drawflow'
import "drawflow/dist/drawflow.min.css"
import { shallowRef, ref, h, render, onMounted } from 'vue'
import * as components from './components/nodes.js'

const name = ref("")
const code = ref(null)
const dialog = ref(null)
const editor = shallowRef({})
const warn = ref({error:false})
const nodeData = ref([
  {name:'assignation', type:'assign', class:'Value', in:1},
  {name:'number', type:'num', class:'Value'},
  {name:'operation', type:'operations', class:'Operation', in:2},
  {name:'if-else block', type:'flowcon', class:'Conditional', in:1, out:2},
  {name:'for loop', type:'flowloop', class:'Loop', in:1}
])
var nodeList = [];
var tempSave = {};
var script;

function showWarning(text) {
  warn.value.text = text
  warn.value.error = true
  setTimeout(() => {warn.value.error = false}, 5000)
  //clearTimeout()
}

// Usar ref para nameLabel?
function setName() {
  let nameLabel = document.getElementById('scriptName')
  nameLabel.innerHTML = name.value
  dialog.value.close()
  saveCode()
}

// Hacer función "overwriteCode"
function saveCode() {
  tempSave.name = name.value
  tempSave.list = nodeList
  script.arrayBuffer()
  .then((data) => {tempSave.code = data})
  tempSave.nodes = editor.value.export()
  /*const http = new XMLHttpRequest()
  console.log(tempSave)
  http.open('POST', 'http://localhost:8080/scripts')
  http.addEventListener('load', () => {console.log(http.response)})
  http.send(JSON.stringify(tempSave))*/
}

function loadCode() {
  let nameLabel = document.getElementById('scriptName')
  name.vaule = tempSave.name
  nameLabel.innerHTML = name.value
  nodeList = tempSave.list
  code.value.data = createScript([tempSave.code])
  editor.value.import(tempSave.nodes)
}

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
  script = new Blob(data, {type:"text/plain;charset=utf-8"})
  let scriptUrl = window.URL.createObjectURL(script)
  return scriptUrl
}

function generateCode(lastNode) {
  console.log(lastNode)
  var codeLine;
  var formated = false;
  var nodeInfo;
  var inputs;

  if (typeof lastNode == 'string') {
    nodeInfo = editor.value.getNodeFromId(lastNode);
  } else {
    nodeInfo = editor.value.getNodeFromId(lastNode.id);
  }

  inputs = nodeInfo.inputs

  if (nodeInfo.name == 'assignation') {
    codeLine = nodeInfo.data.val + ' = '
  } else if (nodeInfo.name == 'operation') {
    let symbol;
    if (nodeInfo.data.val == 'add') {
      symbol = '+'
    } else if (nodeInfo.data.val == 'sub') {
      symbol = '-'
    } else if (nodeInfo.data.val == 'mul') {
      symbol = '*'
    } else if (nodeInfo.data.val == 'div') {
      symbol = '/'
    }
    let aNum = generateCode(inputs.input_1.connections[0].node)
    let bNum = generateCode(inputs.input_2.connections[0].node)
    codeLine = `${aNum} ${symbol} ${bNum}`
    formated = true;
  } else {
    codeLine = nodeInfo.data.val
  }

  if (formated == false && nodeInfo.name != 'number') {
    for (let input in inputs) {
      for (let node of nodeList) {
        if (node.id == inputs[input].connections[0].node) {
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

function addConnection(output_id) {
  for (let node of nodeList) {
    if (output_id == node.id) {
      node.output = true
    }
  }
}

function addNode(data) {
  var vars = {}
  if (data.class == 'Conditional') {
    vars = {'val':'', 'con':''}
  } else {
    vars = {'val':''}
  }
  editor.value.addNode(
    data.name,
    data.in? data.in : 0,
    data.out? data.out : 1,
    0,
    0,
    data.class,
    vars,
    data.type,
    'vue'
  )
}

onMounted(() => {
  // Initialices Drawflow
  let id = document.getElementById("drawflow");
  let Vue = { version: 3, h, render };
  editor.value = new Drawflow(id, Vue);

  // Registers all nodes
  for (let node of nodeData.value) {
    var comp;
    var props = {};
    if (node.class == 'Operation') {
      comp = components.operations
    } else if (node.class == 'Value') {
      comp = components.datatypes
      props = {'type':node.type}
    } else if (node.class == 'Conditional') {
      comp = components.conditional
    } else if (node.class == 'Loop') {
      comp = components.loop
    } else {
      comp = components.datatypes
    }
    editor.value.registerNode(
      node.type,
      comp,
      props,
      {}
    )
  }

  // Keeps track of created nodes
  editor.value.on('nodeCreated', (id) => {
    console.log('New node:', id);
    console.log(editor.value.getNodeFromId(id));
    nodeList.push({'id':id, 'output':false})
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

    if (output.data.val == '') {
      editor.value.removeSingleConnection(data.output_id, data.input_id, data.output_class, data.input_class)
      showWarning('The node you are connecting doesn\'t has value!')
    } else if (input.class == 'Operation') {
      if (input.data.val == '') {
        editor.value.removeSingleConnection(data.output_id, data.input_id, data.output_class, data.input_class)
        showWarning('You have to choose an operation!')
      } else {
        addConnection(data.output_id)
      }
    } else if (input.class == 'Value') {
      if (input.data.val == '') {
        editor.value.removeSingleConnection(data.output_id, data.input_id, data.output_class, data.input_class)
        showWarning('Missing name in assignation node!')
      } else {
        addConnection(data.output_id)
      }
    }
  })

  editor.value.start();
})

</script>

<template>
  <div class="box">
    <div v-if="warn.error">
      <dialog open id="warn-box">
        <p>{{warn.text}}</p>
      </dialog>
    </div>
    <dialog ref="dialog">
      <p>Enter a name</p>
      <input type="text" placeholder="Name" v-model="name">
      <button @click="dialog.close()">Cancel</button>
      <button @click="setName()">Accept</button>
    </dialog>
    <p id="script-name">Unsaved</p>
    <div class="left-panel">
      <h3>Blocks</h3>
      <ul>
        <li v-for="data in nodeData">
          <button @click="addNode(data)">New {{data.name}}</button>
        </li>
      </ul>
    </div>
    <div id="drawflow"></div>
    <div class="right-panel">
      <button @click="renderCode" id="generate">Generate code</button>
      <button @click="sendExecTree()" id="execute">Execute code</button>
      <div id="list"><object ref="code" width=200 height=400></object></div>
      <button @click="name==''? dialog.showModal() : saveCode()" class="database" id="save">Save code</button>
      <button @click="loadCode()" class="database" id="load">Load code</button>
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

#warn-box {
  top: -15px;
  border: 2px solid black;
  border-radius: 20%;
}

#script-name {
  position: absolute;
  left: 238px;
  top: 5px;
}

.left-panel {
  height: 100%;
  width: 25%;
  top: 0px;
  left: 0px;
  background: rgb(218, 230, 233);
  padding-left: 15px;
}

.left-panel * {
  font-size: medium;
  font-family: Arial, Helvetica, sans-serif;
}

.left-panel h3 {
  position: relative;
  left: 60px;
  border: 5px solid black;
  width: 52px;
  padding: 10px;
  background-color: white;
}

.left-panel ul {
  list-style-type: none;
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
  position: relative;
  height: 30px;
}

.right-panel .database {
  position: absolute;
  bottom: 0px;
}

#generate{
  left: 5px;
  top: 15px;
}

#execute {
  top: 15px;
  right: -48px;
}

#save {
  right: 170px;
}

#load {
  right: 5px;
}

.right-panel div {
  position: relative;
  top: 30px;
}

#list object {
  font: Arial;
  font-size: 14px;
  background-color: white;
  width: 97%;
  height: 512px;
  margin-left: 4px;
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
  top: 34px;
}

.drawflow .drawflow-node.Operation .output {
  top: 49px;
}

.drawflow .drawflow-node.Value .input,
.drawflow .drawflow-node.Value .output{
  top: 22px;
}

.drawflow .drawflow-node.Loop .output {
  top: 42px;
  right: -4px;
}

</style>
