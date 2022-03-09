<script setup>
const props = defineProps({
    name:String,
    type:String
})
</script>

<template>
  <div v-if="props.type=='assign'">
    <h1>Assignation</h1>
    <input placeholder="Variable name">
  </div>
  <div v-else>
    <h3>{{props.name}}</h3>
    <p v-if="props.type=='add'">A + B</p>
    <p v-else-if="props.type=='sub'">A - B</p>
    <p v-else-if="props.type=='mul'">A * B</p>
    <p v-else-if="props.type=='div'">A / B</p>
    <p v-else>Err</p>
  </div>
</template>