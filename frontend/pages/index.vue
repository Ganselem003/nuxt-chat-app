<template>
  <div class="container">
    <div class="box">
      <h1>Welcome to my Web Chat App</h1>
      <a href="/login" class="btn">Login Here</a>
    </div>
  </div>
</template>

<style scoped>
.container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: white;
}

.box {
  text-align: center;
  color: black;
}

.btn {
  display: inline-block;
  background: #008000;
  color: white;
  padding: 10px 20px;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 600;
  transition: 0.2s;
}

.btn:hover {
  background: #00B496;
}
</style>
