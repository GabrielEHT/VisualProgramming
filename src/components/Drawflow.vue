<script setup>
//import * as Vue from 'vue'
import Drawflow from 'drawflow'
import "drawflow/dist/drawflow.min.css"
import { shallowRef, h, render, onMounted } from 'vue'
import testComponent from './test.vue'

const Vue = { version: 3, h, render };
const editor = shallowRef({})

onMounted(() => {
  let id = document.getElementById("drawflow");
  console.log(id)
  editor.value = new Drawflow(id, Vue);
  editor.value.registerNode('works', testComponent, {}, {})
  editor.value.start();
  editor.value.addNode("test", 1, 1, 150, 50, "test", {}, 'works', 'vue')
})
</script>

<template>
<div id="drawflow"></div>
</template>

<style scoped>
#drawflow {
  width: 100%;
  height: 500px;
  border: 1px solid red;
  text-align: initial;
}
</style>