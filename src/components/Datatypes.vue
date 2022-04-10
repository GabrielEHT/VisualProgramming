<script setup>
  import { ref, watch } from 'vue'

  const num = ref('')
  const varname = ref('')

  watch(num, (val) => {
    console.log('Actual value', parseInt(val))
  })

  watch(varname, (val) => {
    console.log('Value:', varname)
  })

  const props = defineProps({
    type:''
  })

</script>

<template>
  <div v-if="props.type=='num'">
    <h1>Number</h1>
    <input v-model="num" type="text" placeholder="Number" df-val>
  </div>
  <div v-else-if="props.type=='var'">
    <h1>Variable</h1>
    <input v-model="varname" type="text" placeholder="Variable name" df-val>
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