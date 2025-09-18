<template>
  <div class="chat-page">
    <h2>OpenIM Chat Demo</h2>

    <div v-if="!isLogin" class="login-box">
      <input v-model="userID" placeholder="输入你的用户ID" />
      <!-- <input v-model="token" placeholder="输入 token" /> -->
      <button @click="login">登录</button>
    </div>

    <div v-else class="chat-box">
      <div class="chat-header">
        <span>已登录用户：{{ userID }}</span>
        <button @click="logout">退出</button>
      </div>

      <div class="receiver">
        <input v-model="receiverID" placeholder="输入对方用户ID" />
      </div>

      <div class="messages">
        <div
          v-for="(msg, idx) in messages"
          :key="idx"
          :class="{ self: msg.sender === userID }"
        >
          <b>{{ msg.sender }}:</b> {{ msg.content }}
        </div>
      </div>

      <div class="send-box">
        <input
          v-model="newMessage"
          placeholder="输入消息"
          @keyup.enter="sendMessage"
        />
        <button @click="sendMessage">发送</button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { getSDK,CbEvents } from '@openim/client-sdk';

const IMSDK = getSDK();

IMSDK.on(CbEvents.OnConnectSuccess, () => console.log("WS 已连接"));
IMSDK.on(CbEvents.OnConnectFailed, (evt) => console.error("WS 连接失败", evt));

const isLogin = ref(false);
const userID = ref("4319292610");
const platformID = ref(5);
const token = ref("");
const receiverID = ref("5927646776");
const newMessage = ref("1234567489");
const messages = ref<{ sender: string; content: string }[]>([]);

async function login() {
  try {

    // 1. 调用后端获取 token
    const res = await fetch("http://localhost:8081/token", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ userID: userID.value, platformID: platformID.value }),
    });
    const data = await res.json();
    token.value = data.token;

    const config = {
        userID: userID.value,       // IM 用户 userID
        token: token.value,        // IM 用户令牌
        platformID: platformID.value,   // 当前登录平台号
        apiAddr: "http://192.168.0.100:10002",   // IM api 地址，一般为`http://your-server-ip:10002`或`https://your-server-ip/api
        wsAddr: "ws://192.168.0.100:10001/ws",    // IM ws 地址，一般为`ws://your-server-ip:10001`或`wss://your-server-ip/msg_gateway
    }
    IMSDK.login(config)
    .then(async (res) => {
        // 登录成功
        console.log("登录成功", res);
    })
    .catch(({ errCode, errMsg }) => {
        // 登录失败
        console.error("登录失败", errCode, errMsg);
    });

    isLogin.value = true;
  } catch (err) {
    console.error("登录失败", err);
  }
}

async function sendMessage() {
  if (!receiverID.value || !newMessage.value) return;

   try {
    // 1. 创建文本消息
    const { data: msg } = await IMSDK.createTextMessage(newMessage.value);
    if (!msg) {
      console.error("消息创建失败，msg 为 null");
      return;
    }
    console.log("消息创建成功", msg);

    // 2. 发送消息
    const req = {
      recvID: receiverID.value,
      groupID: "",
      message: msg,
    }
    console.log("消息发送请求", req);
    const { data: sendRes } = await IMSDK.sendMessage(req);
    console.log("消息发送成功", sendRes);

    // 3. 本地显示
    messages.value.push({
      sender: userID.value,
      content: newMessage.value
    });

    newMessage.value = "";
  } catch (err: any) {
    console.error("消息发送失败", err.errCode, err.errMsg || err.message);
  }
}
</script>

<style scoped>
.chat-page { width: 400px; margin: 20px auto; font-family: Arial, sans-serif; }
.login-box, .chat-box { border: 1px solid #ddd; padding: 15px; border-radius: 6px; }
.messages { height: 200px; overflow-y: auto; border: 1px solid #eee; margin: 10px 0; padding: 5px; }
.messages div { margin: 4px 0; }
.messages div.self { text-align: right; color: blue; }
</style>
