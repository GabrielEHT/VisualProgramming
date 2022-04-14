<template>
	<h1>Assign</h1>
  <div>
    <p>In</p>
    <input type="text" placeholder="Variable name" df-val>
    <p>Out</p>
  </div>
</template>

<style scoped>
  h1,
  div {
    font-family: Arial, Helvetica, sans-serif;
  }

  div {
    display: flex;
    align-items: center;
  }

  h1 {
    position: relative;
    font-size: medium;
    top: -5px;
    left: -1px;
  }
  p {
    position: relative;
    font-size: small;
  }
  input {
    height: 20px;
    width: 100px;
    margin-left: 5px;
    margin-right: 5px;
  }  
</style>