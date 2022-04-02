<script setup>
  import { ref, watch } from 'vue'

  const num = ref('')

  watch(num, (val) => {
    console.log('Actual value', parseInt(val))
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