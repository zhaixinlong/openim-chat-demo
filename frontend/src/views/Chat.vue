<template>
  <div class="chat-page">
    <h2>OpenIM Chat Demo</h2>

    <div v-if="!isLogin" class="login-box">
      <input v-model="userID" placeholder="输入你的用户ID" />
      <input v-model="receiverID" placeholder="输入 receiverID" />
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
IMSDK.on(CbEvents.OnRecvNewMessages, (msgs) => {
    console.log("收到消息", msgs);
    for (const msg of msgs.data) {
      messages.value.push({
        sender: msg.sendID,
        content: msg.textElem.content
      });
    }
});

const isLogin = ref(false);
const userID = ref("");
const platformID = ref(5);
const token = ref("");
const receiverID = ref("");
const newMessage = ref("1234567489");
const messages = ref<{ sender: string; content: string }[]>([]);

const registerUrl = "http://127.0.0.1:8081/user_register"
const tokenUrl = "http://127.0.0.1:8081/token"
const apiAddr = "http://127.0.0.1:10002"
const wsAddr = "ws://127.0.0.1:10001"

async function login() {
  try {

    const resRegist = await fetch(registerUrl, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ userID: userID.value, nickname: `用户${userID.value}`, faceURL: `https://api.dicebear.com/7.x/avataaars/svg?seed=${userID.value}` }),
    });
    console.log("注册结果", resRegist);
    if (resRegist.status != 200) {
      console.error("注册失败", resRegist.status, resRegist.statusText);
      return;
    }
    const resRegistData = await resRegist.json();
    console.log("注册结果 data", resRegistData);
    if (resRegistData.errCode != 0) {
      console.error("注册失败", resRegistData.errCode, resRegistData.errMsg);
      return;
    }

    // 1. 调用后端获取 token
    const res = await fetch(tokenUrl, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ userID: userID.value, platformID: platformID.value }),
    });
    console.log("获取 token 结果 res", res);
    if (res.status != 200) {
      console.error("获取 token", res.status, res.statusText);
      return;
    }
    const data = await res.json();
    console.log("获取 token 结果 data", data);
    if (data.errCode != 0) {
      console.error("获取 token 失败", data.errCode, data.errMsg);
      return;
    }
    token.value = data.data.token;

    const config = {
        userID: userID.value,       // IM 用户 userID
        token: token.value,        // IM 用户令牌
        platformID: platformID.value,   // 当前登录平台号
        apiAddr: apiAddr,   // IM api 地址，一般为`http://your-server-ip:10002`或`https://your-server-ip/api
        wsAddr: wsAddr,    // IM ws 地址，一般为`ws://your-server-ip:10001`或`wss://your-server-ip/msg_gateway
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
