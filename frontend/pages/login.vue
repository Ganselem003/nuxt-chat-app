<template>
  <div class="login-container">
    <h1 class="title">Login</h1>

    <form class="form" @submit.prevent="login">
      <input v-model="username" placeholder="Username" class="input" />
      <input
        v-model="password"
        type="password"
        placeholder="Password"
        class="input"
      />

      <button class="btn">Login</button>
    </form>

    <p class="subtext">
      <NuxtLink to="/register" class="link"
        >Create an account</NuxtLink
      >
    </p>

    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const username = ref("");
const password = ref("");
const error = ref("");

interface LoginResponse {
  token: string;
}

const login = async () => {
  try {
    const config = useRuntimeConfig();
    const res = await $fetch<LoginResponse>(
      `${config.public.apiBase}/api/auth/login`,
      {
        method: "POST",
        body: { username: username.value, password: password.value },
      }
    );

    localStorage.setItem("token", res.token);
    localStorage.setItem("username", username.value);

    router.push("/chat");
  } catch (e: any) {
    error.value = e.data?.message || "Login failed";
  }
};
</script>

<style scoped>
/* ===== Layout ===== */
.login-container {
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

/* ===== Links ===== */
.subtext {
  margin-top: 14px;
  font-size: 14px;
}

.link {
  color: #2563eb;
  text-decoration: underline;
}

/* ===== Error text ===== */
.error {
  color: #ef4444;
  margin-top: 10px;
  font-size: 14px;
}
</style>
