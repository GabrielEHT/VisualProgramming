<script setup>
import Drawflow from 'drawflow'
import "drawflow/dist/drawflow.min.css"
import { shallowRef, ref, h, render, onMounted } from 'vue'
import * as components from './components/nodes.js'

const code = ref(null)
const nameLabel = ref(null)
const listDiag = ref(null)
const scriptList = ref([])
const editor = shallowRef({})
const warn = ref({error:false})
const overwriteWarn = ref(null)
const nodeData = ref([
  {name:'Assignation', type:'assign', class:'Assign', in:2, out:2}, // flow input and output, value input and output
  {name:'Number', type:'num', class:'Value'}, // no inputs, value output
  {name:'Variable', type:'var', class:'Value'}, // no inputs, value output
  {name:'Operation', type:'operations', class:'Operation', in:2}, // two value inputs, value output
  {name:'If-else block', type:'flowcon', class:'Conditional', in:3, out:3}, // flow input, two value inputs and three flow outputs
  {name:'For loop', type:'flowloop', class:'Loop', in:2, out:3}, // flow input, value input, value output and two flow outputs
  {name:'Print', type:'misc', class:'Misc', in:2} // flow and value inputs, flow output
])
var script;
var name;
var nodeList = [];
var tempSave = {};
var coords = {x:100, y:100}

function showWarning(text) {
  warn.value.text = text
  warn.value.error = true
  setTimeout(() => {warn.value.error = false}, 5000)
}

function editName(mode) {
  nameLabel.value.readOnly = false
}

function setName() {
  if (nameLabel.value.value.toLowerCase() != 'unsaved') {
    name = nameLabel.value.value
  }
  nameLabel.value.readOnly = true
}

function checkScrìpts() {
  console.log(name)
  for (let userScript of scriptList.value) {
    console.log(userScript.name)
    if (name == userScript.name) {
      return true
    }
  }
  return false
}

function newScript() {
  name = null
  nameLabel.value.value = 'Unsaved'
  nodeList = []
  editor.value.clear()
  script = null
  code.value.data = ""
}

function requestExecution() {
  if (checkNodes()) {
    createScript()
    const http = new XMLHttpRequest()
    http.open('POST', 'http://localhost:8080/exec')
    http.addEventListener('load', () => {
      stout = document.getElementById('stout')
      stout.innerHTML = http.response
    })
    let wholeScript = ''
    for (let line of script) {
      wholeScript += line
    }
    http.send(JSON.stringify({data:wholeScript}))
  }
}

function saveScript() {
  if (checkNodes()) {
    if (name) {
      if (checkScrìpts()) {
        overwriteWarn.value.showModal()
      } else {
        createScript()
        const http = new XMLHttpRequest()
        let data = {}
        data.name = name
        http.open('POST', 'http://localhost:8080/users/Admin')
        data.list = JSON.stringify(nodeList) + '|'
        let wholeScript = ''
        for (let line of script) {
          wholeScript += line + '|'
        }
        data.script = wholeScript.slice(0, -1)
        data.nodes = JSON.stringify(editor.value.export()) + '|'
        console.log(data)
        http.addEventListener('load', () => {
          console.log(http.response)
          getScriptList()
        })
        http.send(JSON.stringify(data))
      }
    } else {
      showWarning('You have to name your script!')
    }
  }
}

function overwriteScript() {
  overwriteWarn.value.close()
  createScript()
  const http = new XMLHttpRequest()
  let data = {}
  http.open('POST', 'http://localhost:8080/users/Admin/' + name.replaceAll(' ', '_'))
  data.list = JSON.stringify(nodeList) + '|'
  let wholeScript = ''
  for (let line of script) {
    wholeScript += line + '|'
  }
  data.script = wholeScript.slice(0, -1)
  data.nodes = JSON.stringify(editor.value.export()) + '|'
  console.log(data)
  http.addEventListener('load', () => {
    console.log(http.response)
    getScriptList()
  })
  http.send(JSON.stringify(data))
}

function getScriptList() {
  const http = new XMLHttpRequest()
  http.open('GET', 'http://localhost:8080/users/Admin')
  http.addEventListener('load', () => {
    if (scriptList.value.err) {
      showWarning('Successfully connected to the server')
      clearInterval(scriptList.value.id)
      scriptList.value = []
    }
    if (http.response != "empty") {
      scriptList.value = JSON.parse(http.response)
    }
  })
  http.addEventListener('error', () => {
    showWarning('Server error, wait a moment or reload the page')
    if (scriptList.value.err == undefined) {
      scriptList.value = {
        err:true,
        id:setInterval(getScriptList, '15000')
      }
    }
  })
  http.send()
}

function loadScript(scriptName) {
  listDiag.value.close()
  const http = new XMLHttpRequest()
  http.open('GET', 'http://localhost:8080/users/Admin/' + scriptName.replaceAll(' ', '_'))
  http.addEventListener('load', () => {
    const resp = JSON.parse(http.response)
    name = resp.name
    nameLabel.value.value = name
    nodeList = JSON.parse(resp.nodeList.slice(0, -1))
    editor.value.import(JSON.parse(resp.drawflow.slice(0, -1)))
    script = resp.code.split("|")
    createScript(script)
  })
  http.send()
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

function checkNodes() {
  if (nodeList.length) {
    updateNodeData()
    console.log(nodeList)
    let message;
    let rootCount = 0;
    for (let node of nodeList) {
      if (node.data.val == '' && node.class != 'Misc') {
        switch (node.class) {
          case 'Assign':
            message = 'There is an assignation node without name!'
            break
          case 'Value':
            message = `There is a ${node.name.toLowerCase()} node without a value!`
            break
          case 'Operation':
            message = 'You have to select a operation for all operation nodes!'
            break
          case 'Conditional':
            message = 'You have to select a condition for all if-else nodes!'
            break
        }
        showWarning(message)
        return false
      } else if (node.inputs.input_2 && node.inputs.input_2.connections.length == 0) {
        showWarning('You have unconnected nodes!')
        return false
      } else if (node.inputs.input_3 && node.inputs.input_3.connections.length == 0) {
        showWarning('You have unconnected nodes!')
        return false
      } else if (node.inputs.input_1 && node.inputs.input_1.connections.length == 0) {
        if (rootCount) {
          showWarning('Too many root nodes') // mejorar
          return false
        } else {
          rootCount++
        }
      }
    }
    return true
  } else {
    showWarning('You haven\'t created any nodes!');
    return false
  }
}

function createScript() {
  let execTree = generateExecTree()
  console.log(execTree)
  script = generateCode(execTree)
  console.log(script)
  let scriptBlob = new Blob(script, {type:"text/plain;charset=utf-8"})
  let scriptUrl = window.URL.createObjectURL(scriptBlob)
  code.value.data = scriptUrl
}

function resolveValueNodes(id) {
  let node = editor.value.getNodeFromId(id)
  let result
  switch (node.class) {
    case 'Assign':
    case 'Loop':
    case 'Value':
      result = node.data.val
      break
    case 'Operation':
      let a = resolveValueNodes(node.inputs.input_1.connections[0].node)
      let b = resolveValueNodes(node.inputs.input_2.connections[0].node)
      result = `${a} ${node.data.val} ${b}`
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
    if (line.hasOwnProperty('assignation')) {
      let node = editor.value.getNodeFromId(line.assignation)
      let result = resolveValueNodes(node.inputs.input_2.connections[0].node)
      codeText.push(spaces + `${node.data.val} = ${result}\n`)
    } else if (line.hasOwnProperty('if')) {
      let node = editor.value.getNodeFromId(line.id)
      let a = resolveValueNodes(node.inputs.input_2.connections[0].node)
      let b = resolveValueNodes(node.inputs.input_3.connections[0].node)
      codeText.push(spaces + `if ${a} ${node.data.val} ${b}:\n`)
      let ifBlock = generateCode(line['if'], indentLevel + 1)
      for (let indLine of ifBlock) {
        codeText.push(indLine)
      }
      if (line['else']) {
        codeText.push(spaces + 'else:\n')
        let elseBlock = generateCode(line['else'], indentLevel + 1)
        for (let indLine of elseBlock) {
          codeText.push(indLine)
        }
      }
    } else if (line.hasOwnProperty('for')) {
      let node = editor.value.getNodeFromId(line.id)
      let result = resolveValueNodes(node.inputs.input_2.connections[0].node)
      codeText.push(spaces + `for i in range(${result}):\n`)
      let forBlock = generateCode(line['for'], indentLevel + 1)
      for (let indLine of forBlock) {
        codeText.push(indLine)
      }
    } else if (line.hasOwnProperty('print')) {
      let node = editor.value.getNodeFromId(line.print)
      let result = resolveValueNodes(node.inputs.input_2.connections[0].node)
      codeText.push(spaces + `print(${result})\n`)
    }
  }
  return codeText
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
  switch (rootNode.class) {
    case 'Assign':
    case 'Misc':
      execTree.push({[rootNode.name.toLowerCase()]:rootNode.id})
      break
    case 'Conditional':
      let conditional = {id:rootNode.id}
      nextNode = getNodeFromId(rootNode.outputs.output_2.connections[0].node);
      conditional['if'] = generateExecTree(nextNode, [])
      if (rootNode.outputs.output_3.connections.length > 0) {
        nextNode = getNodeFromId(rootNode.outputs.output_3.connections[0].node);
        conditional['else'] = generateExecTree(nextNode, [])
      }
      execTree.push(conditional)
      break
    case 'Loop':
      let loop = {id:rootNode.id}
      nextNode = getNodeFromId(rootNode.outputs.output_2.connections[0].node);
      loop['for'] = generateExecTree(nextNode, [])
      execTree.push(loop)
      break
   default:
      console.log('Unknown class')
  }
  if (rootNode.outputs.output_1.connections.length > 0) {
    nextNode = getNodeFromId(rootNode.outputs.output_1.connections[0].node);
    execTree = generateExecTree(nextNode, execTree)
  }
  return execTree
}

function addNode(data) {
  let vars = {}
  switch (data.class) {
    case 'Conditional':
      vars = {'val':'', 'con':''}
      break
    case 'Loop':
      vars = {'val':'i'}
      break
    default:
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
  getScriptList()

  // Initialices Drawflow
  let id = document.getElementById("drawflow");
  let Vue = { version: 3, h, render };
  editor.value = new Drawflow(id, Vue);

  // Registers all nodes
  for (let node of nodeData.value) {
    var comp;
    var props = {};
    switch (node.class) {
      case 'Assign':
        comp = components.assign
        break
      case 'Operation':
        comp = components.operations
        break
      case 'Value':
        comp = components.datatypes
        props = {'type':node.type}
        break
      case 'Conditional':
        comp = components.conditional
        break
      case 'Loop':
        comp = components.loop
        break
      case 'Misc':
        comp = components.misc
        break
      default:
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
      flow_outputs = ['output_1', 'output_2', 'output_3']
    } else if (node.class == 'Loop') {
      flow_outputs = ['output_1', 'output_2']
    } else if (node.class != 'Operation' && node.class != 'Value') {
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
    <dialog ref="overwriteWarn">
      <h1>Warning!</h1>
      <p>This script already exists, do you wish to overwrite it?</p>
      <button @click="overwriteWarn.close()">Cancel</button>
      <button @click="overwriteScript()">Accept</button>
    </dialog>
    <dialog ref="listDiag">
      <button @click="listDiag.close()">X</button>
      <div v-if="scriptList.err">
        <h1>Server error...</h1>
        <p>Could't get your scripts list</p>
      </div>
      <div v-else>
        <h1>Your scripts:</h1>
        <div v-if="scriptList.length == 0">
          <p>You don't have scripts</p>
        </div>
        <div v-else>
          <ul>
            <li v-for="userScript in scriptList">
              <button @click="loadScript(userScript.name)">{{userScript.name}}</button>
            </li>
          </ul>
        </div>
      </div>
    </dialog>
    <input readonly id="script-name" ref="nameLabel" value="Unsaved" @dblclick="editName()" @focusout="setName()">
    <div class="left-panel">
      <h3>Nodes</h3>
      <ul>
        <li v-for="data in nodeData">
          <button @click="addNode(data)">{{data.name}}</button>
        </li>
      </ul>
      <button @click="newScript()">New script</button>
    </div>
    <div id="drawflow"></div>
    <div class="right-panel">
      <button @click="createScript()" id="generate">Generate script</button>
      <button @click="requestExecution()" id="execute">Execute script</button>
      <div id="code">
        <object ref="code"></object>
      </div>
      <textarea id="stout" readonly>Waiting script execution...</textarea>
      <button @click="saveScript()" class="database" id="save">Save script</button>
      <button @click="listDiag.showModal()" class="database" id="load">Load script</button>
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
  border: 3px solid rgba(72, 212, 163, 0.83);
  border-radius: 20px;
  z-index: 1;
  background-color: rgba(83, 245, 142, 0.96);
  animation-name: slideIn, slideOut;
  animation-duration: 1s, 1s;
  animation-delay: 0s, 4s;
}

@keyframes slideIn {
  from {
    top: -90px;
  }
  to {
    top: -15px;
  }
}

@keyframes slideOut {
  from {
    top: -15px;
  }
  to {
    top: -90px;
  }
}

#script-name {
  position: absolute;
  left: 238px;
  top: 5px;
  z-index: 1;
  border: 0px;
}

.left-panel,
.right-panel {
  height: 100%;
  background-color: rgba(83, 169, 245, 0.96);
}

.left-panel {
  width: 25%;
  top: 0px;
  left: 0px;
  padding-left: 15px;
}

.left-panel * {
  font-size: medium;
  font-family: Arial, Helvetica, sans-serif;
}

.left-panel h3 {
  position: relative;
  left: 64px;
  border: 2px solid black;
  width: 54px;
  padding: 10px;
  background-color: rgba(255, 255, 255, 0.70);
  font-size: large;
}

.left-panel ul {
  list-style-type: none;
}

#drawflow {
  width: 100%;
  height: 100%;
  text-align: initial;
  border-left: 2px solid black;
  border-right: 2px solid black;
}

.right-panel {
  width: 30%;
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
  right: -28px;
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

#code {
  height: 400px;
  background-color: white;
  margin-left: 5px;
  margin-right: 5px;
}

#code object {
  font-size: 14px;
  background-color: transparent;
  width: 100%;
  height: 100%;
}

#stout {
  position: relative;
  bottom: -37px;
  margin-left: 5px;
  margin-right: 5px;
  resize: none;
  width: 93%;
  height: 91px;
}

</style>

<style>

/* Nodes colors */

.drawflow .parent-node .drawflow-node {
    background-color: rgb(200, 210, 220);
    width: auto;
}

.drawflow .parent-node .drawflow-node.selected {
    background-color: rgb(220, 230, 240);
}

/* Colors for flow outputs */

.drawflow .parent-node .drawflow-node.Conditional .output,
.drawflow .parent-node .drawflow-node:not(.Value):not(.Operation) .output.output_1,
.drawflow .parent-node .drawflow-node.Loop .output.output_2,
.drawflow .parent-node .drawflow-node:not(.Operation) .input.input_1 {
  background-color: rgba(255, 255, 255, 0);
  border-radius: 0;
  border-left: 20px solid mediumorchid;
  border-right: 0px;
  border-top: 10px solid transparent;
  border-bottom: 10px solid transparent;
  width: 0px;
  height: 0px;
}

/* Colors for value outputs */

.drawflow .parent-node .drawflow-node .output {
  background-color: rgb(255, 167, 33);
}

.drawflow .drawflow-node .output:hover {
  background-color: rgb(255, 194, 102);
}

/* Assignation node styling */

.drawflow .parent-node .drawflow-node.Assign .output.output_1,
.drawflow .parent-node .drawflow-node.Assign .input.input_1 {
  top: -9px;
}

.drawflow .parent-node .drawflow-node.Assign .output.output_2,
.drawflow .parent-node .drawflow-node.Assign .input.input_2 {
  top: 10px;
}

.drawflow .parent-node .drawflow-node.Assign .input.input_1 {
  left: -22px;
}

.drawflow .parent-node .drawflow-node.Assign .output.output_1 {
  left: 8px
}

/* Operation node styling */

.drawflow .drawflow-node.Operation .input {
  top: 34px;
}

.drawflow .drawflow-node.Operation .output {
  top: 49px;
}

/* Value nodes styling */

.drawflow .drawflow-node.Value .input,
.drawflow .drawflow-node.Value .output{
  top: 22px;
}

/* If-else node styling */

.drawflow .drawflow-node.Conditional .drawflow_content_node {
  height: 182px;
}

.drawflow .drawflow-node.Conditional .input.input_1 {
  top: -39px;
}

.drawflow .drawflow-node.Conditional .input.input_2{
  top: -10px;
}

.drawflow .drawflow-node.Conditional .input.input_3{
  top: 23px;
}

.drawflow .drawflow-node.Conditional .output {
  top: 53px;
}

.drawflow .drawflow-node.Conditional .output.output_1 {
  top: -43px;
}

/* For loop node styling */

.drawflow .drawflow-node.Loop .drawflow_content_node {
  height: 125px;
}

.drawflow .drawflow-node.Loop .input.input_1 {
  top: -26px;
  left: -26px;
}

.drawflow .drawflow-node.Loop .input {
  top: 3px;
}

.drawflow .drawflow-node.Loop .output.output_1 {
  top: -12px;
}

.drawflow .drawflow-node.Loop .output {
  top: 31px;
  right: -6px;
}

/* Miscellaneous nodes styling */

.drawflow .drawflow-node.Misc .output {
  top: -12px;
  right: -8px;
}

</style>
