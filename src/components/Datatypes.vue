<script setup>
  import { ref } from 'vue'

  const num = ref('')
  const varname = ref('')

  function onlyNums(e) {
    console.log(e)
    if (e.data === '-') {
      if (num.value.includes('-', 1) || num.value.length > 1) {
        num.value = num.value.slice(0, -1)
      }
    } else if (isNaN(e.data)) {
      num.value = num.value.slice(0, -1)
    }
  }

  const props = defineProps({
    type:''
  })

</script>

<template>
  <div v-if="props.type=='num'">
    <h1>Number</h1>
    <input v-model="num" type="text" placeholder="Number" @input="onlyNums" maxlength="20" df-val>
  </div>
  <div v-else-if="props.type=='var'">
    <h1>Variable call</h1>
    <input v-model="varname" type="text" placeholder="Variable name" maxlength="20" df-val>
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