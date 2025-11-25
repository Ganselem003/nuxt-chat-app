<template>
  <div class="register-container">
    <h1 class="title">Register</h1>

    <form class="form" @submit.prevent="register">
      <input v-model="username" placeholder="Username" class="input" />

      <input
        v-model="password"
        type="password"
        placeholder="Password"
        class="input"
      />

      <button class="btn">Register</button>
    </form>

    <p v-if="message" class="success">{{ message }}</p>
    <p v-if="error" class="error">{{ error }}</p>

    <p class="subtext">
      <NuxtLink to="/login" class="link">Back to login</NuxtLink>
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";

const username = ref("");
const password = ref("");
const message = ref("");
const error = ref("");
const router = useRouter();

const register = async () => {
  try {
    const config = useRuntimeConfig();

    await $fetch(`${config.public.apiBase}/api/auth/register`, {
      method: "POST",
      body: { username: username.value, password: password.value },
    });

    message.value = "Registered. Redirecting to login...";
    setTimeout(() => router.push("/login"), 800);
  } catch (e: any) {
    error.value = e.data?.message || "Registration failed";
  }
};
</script>

<style scoped>
/* ===== Layout ===== */
.register-container {
  max-width: 380px;
  margin: 80px auto;
  text-align: center;
}

/* ===== Title ===== */
.title {
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 24px;
}

/* ===== Form ===== */
.form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* ===== Inputs ===== */
.input {
  padding: 12px;
  border: 1px solid #ccc;
  border-radius: 8px;
  font-size: 16px;
  width: 100%;
  transition: border 0.2s;
}

.input:focus {
  border-color: #3b82f6;
  outline: none;
}

/* ===== Button ===== */
.btn {
  background: #008000;
  color: white;
  padding: 12px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: 0.2s;
  border: none;
}

.btn:hover {
  background: #00b496;
}

/* ===== Feedback messages ===== */
.success {
  color: #22c55e;
  margin-top: 10px;
  font-size: 14px;
}

.error {
  color: #ef4444;
  margin-top: 10px;
  font-size: 14px;
}

/* ===== Back link ===== */
.subtext {
  margin-top: 14px;
  font-size: 14px;
}

.link {
  color: #2563eb;
  text-decoration: underline;
}
</style>
