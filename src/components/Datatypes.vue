<script setup>
  const props = defineProps({
    type:''
  })
</script>

<template>
  <div v-if="props.type=='num'">
    <h1>Number</h1>
    <input type="text" placeholder="Enter a number" df-val>
  </div>
  <div v-else-if="props.type=='assign'">
    <h1>Assignation</h1>
    <input type="text" placeholder="Enter variable name" df-val>
  </div>
  <div v-else>
    <h1>Error</h1>
  </div>
</template>

<style scoped>
  h1 {
    font-family: Arial, Helvetica, sans-serif;
    font-size: medium;
  }

  input {
    font-family: Arial, Helvetica, sans-serif;
  }
</style>