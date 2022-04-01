<script setup>
import Drawflow from 'drawflow'
import "drawflow/dist/drawflow.min.css"
import { shallowRef, ref, h, render, onMounted } from 'vue'
import * as components from './components/nodes.js'

// Crear "helpBox"
const name = ref("")
const code = ref(null)
const dialog = ref(null)
const editor = shallowRef({})
const warn = ref({error:false})
const nodeData = ref([
  {name:'assignation', type:'assign', class:'Value', in:2, out:2}, //flow input and output, value input and output
  {name:'number', type:'num', class:'Value'}, //no inputs, value output
  {name:'operation', type:'operations', class:'Operation', in:2}, //two value inputs, value output
  {name:'if-else block', type:'flowcon', class:'Conditional', in:3, out:2}, //flow input, two value inputs and two flow outputs
  {name:'for loop', type:'flowloop', class:'Loop', in:2} //flow input, value input and flow output
])
var script;
var nodeList = [];
var tempSave = {};
var coords = {x:100, y:100}

function requestExecution() {
  const http = new XMLHttpRequest()
  http.open('POST', 'http://localhost:8080/exec')
  http.addEventListener('load', () => {
    console.log(http.response)
  })
  script.arrayBuffer()
  .then((buffer) => {
    let bytes = new Int8Array(buffer)
    http.send(JSON.stringify({data:bytes}))
  })
}

// Añadir animación
function showWarning(text) {
  warn.value.text = text
  warn.value.error = true
  setTimeout(() => {warn.value.error = false}, 5000)
}

// usar ref para nameLabel?
function setName() {
  let nameLabel = document.getElementById('script-name')
  nameLabel.innerHTML = name.value
  dialog.value.close()
  saveCode()
}

// hacer función "overwriteCode"
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
  let nameLabel = document.getElementById('script-name')
  name.vaule = tempSave.name
  nameLabel.innerHTML = name.value
  nodeList = tempSave.list
  code.value.data = createScript([tempSave.code])
  editor.value.import(tempSave.nodes)
}

function getRootNode() {
  let rootNode;
  for (let node of nodeList) {
    if (node.inputs.input_1 != undefined && node.inputs.input_1.connections.length == 0) {
      if (rootNode) {
        rootNode = 'err'
        break
      } else {
        rootNode = node
      }
    }
  }
  return rootNode
}

function getNodeFromId(id) {
  for (let node of nodeList) {
    if(id == node.id) {
      return node
    }
  }
}

function updateNodeData() {
  let i = 0;
  for (let node of nodeList) {
    let nodeInfo = editor.value.getNodeFromId(node.id)
    nodeInfo.flow_inputs = node.flow_inputs
    nodeInfo.flow_outputs = node.flow_outputs
    nodeList[i] = nodeInfo
    i++;
  }
}

function renderCode() {
  if (nodeList.length) {
    updateNodeData()
    console.log(nodeList)
    let validator = true;
    let message;
    for (let node of nodeList) {
      if (node.data.val == '') {
        if (node.html == 'assign') { // cambiar para usar class
          message = 'There is an assignation node without name!'
        } else if (node.html == 'num') {
          message = 'There is a number node without a value!'
        } else if (node.class == 'Operation') {
          message = 'You have to select a operation for all operation nodes!'
        } else if (node.class == 'Conditional') {
          message = 'You have to define a comparison for all if-else nodes!'
        }
        validator = false
        break
      }
    }
    if (validator) {
      let execTree = generateExecTree()
      console.log(execTree)
      let text = generateCode(execTree)
      console.log(text)
      code.value.data = createScript(text)
      // generar codigo
    } else { // hacer que el editor enfoque al nodo del error
      showWarning(message)
    }
  } else {
    showWarning('You haven\'t created any nodes!');
  }
}

function createScript(data) {
  script = new Blob(data, {type:"text/plain;charset=utf-8"})
  let scriptUrl = window.URL.createObjectURL(script)
  return scriptUrl
}

function resolveValueNodes(id) {
  let node = editor.value.getNodeFromId(id)
  let result
  switch (node.class) {
    case 'Value':
      result = node.data.val
      break
    case 'Operation':
      let a = resolveValueNodes(node.inputs.input_1.connections[0].node)
      let b = resolveValueNodes(node.inputs.input_2.connections[0].node)
      result = a + ' ' + node.data.val + ' ' + b
      break
    default:
      console.log('Unknown class')
  }
  return result
}

function generateCode(execTree, indentLevel) {
  let codeText = [];
  if (indentLevel == undefined) {
    indentLevel = 0;
  }
  let spaces = ' '.repeat(indentLevel * 4)
  for (let line of execTree) {
    if (line.hasOwnProperty('assign')) {
      let node = editor.value.getNodeFromId(line.assign)
      let result = resolveValueNodes(node.inputs.input_2.connections[0].node)
      codeText.push(spaces + node.data.val + ' = ' + result + '\n')
    } else if (line.hasOwnProperty('if')) {
      let node = editor.value.getNodeFromId(line.id)
      let a = resolveValueNodes(node.inputs.input_2.connections[0].node)
      let b = resolveValueNodes(node.inputs.input_3.connections[0].node)
      codeText.push(spaces + 'if ' + a + ' ' + node.data.val + ' ' + b + ':\n')
      let ifBlock = generateCode(line['if'], indentLevel + 1)
      for (let indLine of ifBlock) {
        codeText.push(indLine)
      }
      codeText.push(spaces + 'else:\n')
      let elseBlock = generateCode(line['else'], indentLevel + 1)
      for (let indLine of elseBlock) {
        codeText.push(indLine)
      }
    } else if (line.hasOwnProperty('for')) {
      let node = editor.value.getNodeFromId(line.id)
      let result = resolveValueNodes(node.inputs.input_2.connections[0].node)
      codeText.push(spaces + 'for i in range(' + result + '):\n')
      let forBlock = generateCode(line['for'], indentLevel + 1)
      for (let indLine of forBlock) {
        codeText.push(indLine)
      }
    }
  }
  return codeText
  /*if (nodeInfo.name == 'assignation') {
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
  }*/
}

function generateExecTree(rootNode, execTree) {
  if (execTree == undefined) {
    execTree = []
  }
  if (rootNode == undefined) {
    rootNode = getRootNode()
    console.log(rootNode)
  }
  let nextNode;
  if (rootNode.class == 'Conditional') {
    let conditional = {id:rootNode.id}
    nextNode = getNodeFromId(rootNode.outputs.output_1.connections[0].node);
    conditional['if'] = generateExecTree(nextNode, [])
    if (rootNode.outputs.output_2.connections.length > 0) {
      nextNode = getNodeFromId(rootNode.outputs.output_2.connections[0].node);
      conditional['else'] = generateExecTree(nextNode, [])
    }
    execTree.push(conditional)
  } else if (rootNode.class == 'Loop') {
    let loop = {id:rootNode.id}
    nextNode = getNodeFromId(rootNode.outputs.output_1.connections[0].node);
    loop['for'] = generateExecTree(nextNode, [])
    execTree.push(loop)
  } else {
    execTree.push({'assign':rootNode.id})
    if (rootNode.outputs.output_1.connections.length > 0) {
      nextNode = getNodeFromId(rootNode.outputs.output_1.connections[0].node);
      execTree = generateExecTree(nextNode, execTree)
    }
  }
  return execTree
}

// mover a otro módulo?
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
  let vars = {}
  if (data.class == 'Conditional') {
    vars = {'val':'', 'con':''}
  } else if (data.class == 'Loop') {
    vars = {'val':'i'}
  } else {
    vars = {'val':''}
  }
  editor.value.addNode(
    data.name,
    data.in? data.in : 0,
    data.out? data.out : 1,
    coords.x,
    coords.y,
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

  // Defines flow inputs and outputs of new nodes and adds them to the node list
  editor.value.on('nodeCreated', (id) => {
    console.log('New node:', id);
    let node = editor.value.getNodeFromId(id)
    let flow_inputs = [];
    let flow_outputs = [];
    if (node.class != 'Operation' && node.inputs.input_1) {
      flow_inputs = ['input_1']
    }
    if (node.class == 'Conditional') {
      flow_outputs = ['output_1', 'output_2']
    } else if (node.class == 'Loop' || node.html == 'assign') { // cambiar para usar class en lugar de html
      flow_outputs = ['output_1']
    }
    node.flow_inputs = flow_inputs
    node.flow_outputs = flow_outputs
    console.log(node)
    nodeList.push(node)
  })

  // Removes the deleted node from the node list
  editor.value.on('nodeRemoved', (id) => {
    nodeList = nodeList.filter(node => node.id!=id)
    console.log('Removed id:', id, nodeList)
  })

  // Checks if the created connection is valid
  editor.value.on('connectionCreated', (data) => {
    let input = getNodeFromId(data.input_id);
    let output = getNodeFromId(data.output_id);
    let output_type = 'value';
    let input_type = 'value';

    for (let flow of output.flow_outputs) {
      if (data.output_class == flow) {
        output_type = 'flow'
      }
    }

    for (let flow of input.flow_inputs) {
      if (data.input_class == flow) {
        input_type = 'flow'
      }
    }

    if (output_type == 'value') {
      if (input_type == 'flow') {
        editor.value.removeSingleConnection(data.output_id, data.input_id, data.output_class, data.input_class)
        showWarning('You can\'t connect a value output to a flow input!')
      }
    } else {
      if (input_type == 'value') {
        editor.value.removeSingleConnection(data.output_id, data.input_id, data.output_class, data.input_class)
        showWarning('You can\'t connect a flow output to a value input!')
      }
    }
  })

  // arreglar
  // Updates connection state of output node on removed connection
  editor.value.on('connectionRemoved', (data) => {
    let output = editor.value.getNodeFromId(data.output_id)
    let disconnected = false
    if (output.class == 'Conditional') {
      if (output.outputs.output_1.connections.length == 0 && output.outputs.output_2.connections.length == 0) {
        disconnected = true
      }
    } else {
      disconnected = true
    }
    for (let node of nodeList) {
      if (node.id == data.output_id && disconnected) {
        node.output = false
      }
    }
  })

  // Keeps track of the screen position
  editor.value.on('translate', (pos) => {
    coords.x = pos.x * -1 + 100
    coords.y = pos.y * -1 + 100
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
      <button @click="editor.clear()">Clear editor</button>
    </div>
    <div id="drawflow"></div>
    <div class="right-panel">
      <button @click="renderCode()" id="generate">Generate code</button>
      <button @click="requestExecution()" id="execute">Execute code</button>
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
  z-index: 1;
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
