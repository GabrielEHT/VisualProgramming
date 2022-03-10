<script setup>
const props = defineProps({
    type:String
})
</script>

<template>
  <h1>
    <slot name="title"></slot>
  </h1>
  <input v-if="props.type=='num'" placeholder="Number" df-value>
  <input v-else-if="props.type=='assign'" placeholder="Variable name" df-value>
  <p v-else-if="props.type=='add'">A + B</p>
  <p v-else-if="props.type=='sub'">A - B</p>
  <p v-else-if="props.type=='mul'">A * B</p>
  <p v-else-if="props.type=='div'">A / B</p>
  <p v-else>Err</p>
</template>

<style scoped>
h1 {
  font-size: medium;
}

input:hover {
  background-color: gray;
  color: aliceblue;
}
</style>