<template>
  <div class="chat-container">
    <div class="sidebar">
      <h2 class="section-title">Profile</h2>
      <div class="profile-box">
        {{ myUsername }}
      </div>
      <h2 class="section-title">Online users</h2>
      <div
  v-for="user in users"
  :key="user"
  class="user-item"
  :class="{ selected: selected === user }"
  @click="selectUser(user)"
>
  <span class="status-dot"></span>
  {{ user }}
</div>

    </div>
    <div class="chat-main">
      <h1 class="chat-title">
        Chatting with <span class="chat-name">{{ selected || "..." }}</span>
      </h1>
      <div class="chat-window" ref="scrollArea">
        <div
          v-for="msg in messages"
          :key="msg.id"
          class="chat-message"
          :class="msg.from === myUsername ? 'me' : 'them'"
        >
          <div class="bubble">
            <strong class="from">{{ msg.from }}</strong><br />
            {{ msg.text }}
          </div>
        </div>
      </div>
      <form class="chat-input-area" @submit.prevent="sendMessage">
        <input
          v-model="newMessage"
          placeholder="Type a message..."
          class="chat-input"
        />
        <button class="chat-send">Send</button>
      </form>
    </div>

  </div>
</template>


<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, nextTick } from "vue";

interface Message {
  id: number;
  from: string;
  to: string;
  text: string;
}

interface PresenceMessage {
  type: "presence";
  users: string[];
}

interface ChatMessage {
  type: "message";
  from: string;
  to: string;
  text: string;
}

const config = useRuntimeConfig();
const token = localStorage.getItem("token") || "";
const myUsername = localStorage.getItem("username") || "";

const ws = ref<WebSocket | null>(null);
const users = ref<string[]>([]);
const selected = ref<string>("");
const messages = ref<Message[]>([]);
const newMessage = ref("");
const scrollArea = ref<HTMLElement | null>(null);

let msgId = 0;

function addMessage(from: string, to: string, text: string) {
  msgId++;
  messages.value.push({ id: msgId, from, to, text });
  nextTick(() => {
    if (scrollArea.value) {
      scrollArea.value.scrollTop = scrollArea.value.scrollHeight;
    }
  });
}

onMounted(() => {
  if (!token) {
    window.location.href = "/login";
    return;
  }

  ws.value = new WebSocket(`${config.public.wsBase ?? "ws://localhost:8080"}/ws?token=${token}`);

  ws.value.onopen = () => {
    console.log("WebSocket connected");
  };

  ws.value.onmessage = (event: MessageEvent) => {
    try {
      const data: PresenceMessage | ChatMessage = JSON.parse(event.data);
      if ((data as PresenceMessage).type === "presence") {
        const presence = data as PresenceMessage;
        users.value = (presence.users || []).filter(u => u !== myUsername);
      } else if ((data as ChatMessage).type === "message") {
        const chat = data as ChatMessage;
        if (chat.to === myUsername || chat.from === myUsername) {
          addMessage(chat.from, chat.to, chat.text);
        }
      }
    } catch (e) {
      console.error("ws parse error", e);
    }
  };

  ws.value.onclose = () => {
    console.log("WebSocket closed");
  };
});

onBeforeUnmount(() => {
  ws.value?.close();
});

const selectUser = (user: string) => {
  selected.value = user;
};

const sendMessage = () => {
  if (!newMessage.value || !selected.value) return;
  const payload = { to: selected.value, text: newMessage.value };
  ws.value?.send(JSON.stringify(payload));
  addMessage(myUsername, selected.value, newMessage.value);
  newMessage.value = "";
};
</script>

<style scoped>
/* ===== Main layout ===== */
.chat-container {
  display: grid;
  grid-template-columns: 260px 1fr;
  gap: 20px;
  max-width: 1100px;
  margin: 20px auto;
  padding: 0 15px;
}

/* ===== Sidebar ===== */
.sidebar {
  border-right: 1px solid #ddd;
  padding-right: 15px;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 8px;
}

.profile-box {
  padding: 10px;
  background: #eef7ff;
  border-radius: 8px;
  margin-bottom: 20px;
  font-weight: 500;
}

.user-item {
  padding: 10px;
  border-radius: 8px;
  cursor: pointer;
  margin-bottom: 6px;
  border: 1px solid #ddd;
  transition: 0.2s;
}

.user-item:hover {
  background: #f3f3f3;
}

.user-item.selected {
  background: #dcecff;
  border-color: #8bb8ff;
}

/* ===== Chat area ===== */
.chat-main {
  padding-left: 10px;
}

.chat-title {
  font-size: 26px;
  font-weight: bold;
  margin-bottom: 15px;
}

.chat-name {
  color: #2563eb;
}

.chat-window {
  border: 1px solid #ddd;
  height: 420px;
  padding: 12px;
  overflow-y: auto;
  border-radius: 10px;
  background: #fafafa;
}

.chat-message {
  display: flex;
  margin-bottom: 10px;
}

.chat-message.them {
  justify-content: flex-start;
}

.chat-message.me {
  justify-content: flex-end;
}

.bubble {
  max-width: 65%;
  padding: 10px 14px;
  border-radius: 14px;
  line-height: 1.4;
  font-size: 15px;
  background: #e3eaff;
  border: 1px solid #c7d6ff;
}

.chat-message.me .bubble {
  background: #d1fadf;
  border-color: #9eeabf;
}

/* ===== Sender name =====*/
.from {
  font-size: 13px;
  color: #555;
}

.chat-input-area {
  margin-top: 10px;
  display: flex;
  gap: 10px;
}

.chat-input {
  flex: 1;
  padding: 12px;
  border-radius: 10px;
  border: 1px solid #ccc;
  font-size: 15px;
}

.chat-input:focus {
  border-color: #3b82f6;
  outline: none;
}

.chat-send {
  padding: 12px 20px;
  background: #008000;
  color: white;
  border: none;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: 0.2s;
}

.chat-send:hover {
  background: #00b496;
}
.status-dot {
  display: inline-block;
  width: 10px;
  height: 10px;
  background: #22c55e; /* ногоон */
  border-radius: 50%;
  margin-right: 8px;
  border: 2px solid white;
  box-shadow: 0 0 3px rgba(0,0,0,0.3);
  vertical-align: middle;
}


/* ===== Responsive for Mobile/Tablet ===== */
@media (max-width: 768px) {
  .chat-container {
    grid-template-columns: 1fr;
    gap: 10px;
  }

  .sidebar {
    border-right: none;
    border-bottom: 1px solid #ddd;
    padding-bottom: 10px;
    margin-bottom: 10px;
  }

  .chat-main {
    padding-left: 0;
  }

  .chat-title {
    font-size: 22px;
    text-align: center;
  }

  .bubble {
    max-width: 80%;
    font-size: 14px;
  }

  .chat-window {
    height: 350px;
  }

  .chat-input-area {
    flex-direction: column;
  }

  .chat-send {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .bubble {
    max-width: 90%;
  }

  .chat-window {
    height: 300px;
  }

  .chat-title {
    font-size: 20px;
  }

  .chat-input {
    padding: 10px;
  }
}
</style>
