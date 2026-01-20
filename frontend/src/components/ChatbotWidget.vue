<template>
  <div class="chatbot-widget">
    <!-- Chat Button (collapsed state) -->
    <div v-if="!isOpen" class="chat-button" @click="toggleChat">
      <span class="chat-icon">🍷</span>
      <span class="chat-bubble">Need wine help?</span>
    </div>

    <!-- Chat Panel (expanded state) -->
    <div v-else class="chat-panel">
      <!-- Header -->
      <div class="chat-header">
        <div class="header-info">
          <div class="bot-avatar">🍷</div>
          <div class="bot-details">
            <span class="bot-name">Wine Sommelier</span>
            <span class="bot-status">● Online</span>
          </div>
        </div>
        <button class="close-btn" @click="toggleChat">×</button>
      </div>

      <!-- Messages -->
      <div class="chat-messages" ref="messagesContainer">
        <div 
          v-for="(msg, index) in messages" 
          :key="index" 
          :class="['message', msg.sender]"
        >
          <div class="message-content">
            <p v-html="msg.text"></p>
            <div v-if="msg.products" class="product-cards">
              <div 
                v-for="product in msg.products" 
                :key="product.ID" 
                class="product-card"
                @click="goToProduct(product.ID)"
              >
                <img :src="product.image_url || defaultImage" :alt="product.name" />
                <div class="product-info">
                  <span class="name">{{ product.name }}</span>
                  <span class="price">${{ product.price }}</span>
                </div>
              </div>
            </div>
          </div>
          <span class="message-time">{{ msg.time }}</span>
        </div>

        <!-- Typing indicator -->
        <div v-if="isTyping" class="message bot">
          <div class="typing-indicator">
            <span></span><span></span><span></span>
          </div>
        </div>
      </div>

      <!-- Quick Actions (shown when no conversation started) -->
      <div v-if="messages.length === 1" class="quick-actions">
        <button @click="askQuestion('red')">🍷 Red Wines</button>
        <button @click="askQuestion('white')">🥂 White Wines</button>
        <button @click="askQuestion('rose')">🌸 Rosé Wines</button>
        <button @click="askQuestion('budget')">💰 Under $40</button>
        <button @click="askQuestion('premium')">⭐ Premium Pick</button>
      </div>

      <!-- Input -->
      <div class="chat-input">
        <input 
          v-model="userInput" 
          @keyup.enter="sendMessage"
          placeholder="Ask about wines..." 
        />
        <button @click="sendMessage" :disabled="!userInput.trim()">
          <span>➤</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../services/api'

const router = useRouter()
const isOpen = ref(false)
const isTyping = ref(false)
const userInput = ref('')
const messagesContainer = ref(null)
const defaultImage = 'https://res.cloudinary.com/djxwpckvm/image/upload/v1768553510/wine-shop/products/hwarjsfnqqhkvtg8dw74.jpg'

const messages = ref([
  {
    sender: 'bot',
    text: 'Hi there! 👋 I\'m your Wine Sommelier. I can help you find the perfect wine. What are you looking for today?',
    time: getCurrentTime()
  }
])

function getCurrentTime() {
  return new Date().toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
}

function toggleChat() {
  isOpen.value = !isOpen.value
}

function goToProduct(id) {
  router.push(`/products/${id}`)
  isOpen.value = false
}

async function scrollToBottom() {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

async function askQuestion(type) {
  let query = ''
  switch(type) {
    case 'red':
      query = 'Show me red wines'
      break
    case 'white':
      query = 'Show me white wines'
      break
    case 'rose':
      query = 'Show me rosé wines'
      break
    case 'budget':
      query = 'Wines under $40'
      break
    case 'premium':
      query = 'What\'s your best premium wine?'
      break
  }
  userInput.value = query
  await sendMessage()
}

async function sendMessage() {
  if (!userInput.value.trim()) return

  const text = userInput.value.trim()
  
  // Add user message
  messages.value.push({
    sender: 'user',
    text: text,
    time: getCurrentTime()
  })
  
  userInput.value = ''
  await scrollToBottom()
  
  // Show typing indicator
  isTyping.value = true
  await scrollToBottom()

  // Process the message and get response
  await processMessage(text)
  
  isTyping.value = false
  await scrollToBottom()
}

async function processMessage(text) {
  const lowerText = text.toLowerCase()
  let response = { text: '', products: null }

  try {
    // Fetch products for recommendations
    const params = {}
    
    // Detect category
    if (lowerText.includes('red')) {
      params.category = 'Red'
    } else if (lowerText.includes('white')) {
      params.category = 'White'
    } else if (lowerText.includes('rosé') || lowerText.includes('rose')) {
      params.category = 'Rosé'
    }

    // Fetch products
    const res = await api.get('/products', { params: { ...params, limit: 50 } })
    let products = res.data.data || []

    // Filter based on query
    if (lowerText.includes('under $40') || lowerText.includes('budget') || lowerText.includes('cheap')) {
      products = products.filter(p => p.price < 40)
      response.text = '💰 Here are some great wines under $40:'
    } else if (lowerText.includes('premium') || lowerText.includes('best') || lowerText.includes('expensive')) {
      products = products.sort((a, b) => b.price - a.price).slice(0, 3)
      response.text = '⭐ Our premium selections for you:'
    } else if (params.category) {
      response.text = `🍷 Here are our ${params.category.toLowerCase()} wines:`
    } else if (lowerText.includes('recommend') || lowerText.includes('suggest')) {
      products = products.sort(() => Math.random() - 0.5).slice(0, 3)
      response.text = '✨ Based on your taste, I recommend:'
    } else if (lowerText.includes('hello') || lowerText.includes('hi') || lowerText.includes('hey')) {
      response.text = 'Hello! 👋 I\'d love to help you find the perfect wine. Are you looking for red, white, or rosé? Or would you like me to recommend something special?'
    } else if (lowerText.includes('thank')) {
      response.text = 'You\'re welcome! 🍷 Enjoy your wine! Let me know if you need anything else.'
    } else {
      // General search
      if (products.length > 0) {
        products = products.slice(0, 4)
        response.text = 'Here are some wines you might enjoy:'
      } else {
        response.text = 'I\'d be happy to help! Try asking about red wines, white wines, or wines under $40. You can also ask for my premium recommendations!'
      }
    }

    if (products.length > 0 && !response.text.includes('Hello') && !response.text.includes('welcome')) {
      response.products = products.slice(0, 4)
    }

  } catch (error) {
    response.text = 'I\'m having trouble connecting right now. Please try again in a moment! 🍷'
  }

  // Simulate typing delay
  await new Promise(resolve => setTimeout(resolve, 800))

  messages.value.push({
    sender: 'bot',
    text: response.text,
    products: response.products,
    time: getCurrentTime()
  })
}
</script>

<style scoped>
.chatbot-widget {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 9999;
  font-family: 'Inter', sans-serif;
}

/* Chat Button */
.chat-button {
  display: flex;
  align-items: center;
  gap: 10px;
  background: linear-gradient(135deg, #722F37 0%, #8B4513 100%);
  color: white;
  padding: 14px 20px;
  border-radius: 50px;
  cursor: pointer;
  box-shadow: 0 4px 20px rgba(114, 47, 55, 0.4);
  transition: all 0.3s ease;
}

.chat-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 24px rgba(114, 47, 55, 0.5);
}

.chat-icon {
  font-size: 1.5rem;
}

.chat-bubble {
  font-size: 0.95rem;
  font-weight: 500;
}

/* Chat Panel */
.chat-panel {
  width: 380px;
  height: 520px;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Header */
.chat-header {
  background: linear-gradient(135deg, #722F37 0%, #8B4513 100%);
  color: white;
  padding: 16px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.bot-avatar {
  width: 40px;
  height: 40px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.3rem;
}

.bot-details {
  display: flex;
  flex-direction: column;
}

.bot-name {
  font-weight: 600;
  font-size: 1rem;
}

.bot-status {
  font-size: 0.75rem;
  opacity: 0.9;
  color: #90EE90;
}

.close-btn {
  background: none;
  border: none;
  color: white;
  font-size: 1.8rem;
  cursor: pointer;
  padding: 0 8px;
  opacity: 0.8;
  transition: opacity 0.2s;
}

.close-btn:hover {
  opacity: 1;
}

/* Messages */
.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: #f8f4f0;
}

.message {
  display: flex;
  flex-direction: column;
  max-width: 85%;
}

.message.bot {
  align-self: flex-start;
}

.message.user {
  align-self: flex-end;
}

.message-content {
  padding: 12px 16px;
  border-radius: 16px;
  font-size: 0.95rem;
  line-height: 1.4;
}

.message.bot .message-content {
  background: white;
  border: 1px solid #e8e0d8;
  border-radius: 16px 16px 16px 4px;
}

.message.user .message-content {
  background: linear-gradient(135deg, #722F37 0%, #8B4513 100%);
  color: white;
  border-radius: 16px 16px 4px 16px;
}

.message-time {
  font-size: 0.7rem;
  color: #999;
  margin-top: 4px;
  padding: 0 4px;
}

.message.user .message-time {
  text-align: right;
}

/* Product Cards in Chat */
.product-cards {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.product-card {
  display: flex;
  align-items: center;
  gap: 10px;
  background: #f8f4f0;
  padding: 10px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
  width: 100%;
}

.product-card:hover {
  background: #f0e8e0;
  transform: translateX(4px);
}

.product-card img {
  width: 45px;
  height: 45px;
  object-fit: cover;
  border-radius: 8px;
}

.product-info {
  display: flex;
  flex-direction: column;
}

.product-info .name {
  font-weight: 600;
  font-size: 0.85rem;
  color: #333;
}

.product-info .price {
  font-size: 0.8rem;
  color: #722F37;
  font-weight: 500;
}

/* Typing Indicator */
.typing-indicator {
  display: flex;
  gap: 4px;
  padding: 12px 16px;
  background: white;
  border-radius: 16px;
  border: 1px solid #e8e0d8;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  background: #999;
  border-radius: 50%;
  animation: bounce 1.4s infinite ease-in-out;
}

.typing-indicator span:nth-child(1) { animation-delay: 0s; }
.typing-indicator span:nth-child(2) { animation-delay: 0.2s; }
.typing-indicator span:nth-child(3) { animation-delay: 0.4s; }

@keyframes bounce {
  0%, 80%, 100% { transform: scale(0.6); opacity: 0.5; }
  40% { transform: scale(1); opacity: 1; }
}

/* Quick Actions */
.quick-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 12px 16px;
  background: #f8f4f0;
  border-top: 1px solid #e8e0d8;
}

.quick-actions button {
  background: white;
  border: 1px solid #e8e0d8;
  padding: 8px 14px;
  border-radius: 20px;
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.2s;
}

.quick-actions button:hover {
  background: #722F37;
  color: white;
  border-color: #722F37;
}

/* Input */
.chat-input {
  display: flex;
  gap: 10px;
  padding: 16px;
  background: white;
  border-top: 1px solid #e8e0d8;
}

.chat-input input {
  flex: 1;
  border: 1px solid #ddd;
  padding: 12px 16px;
  border-radius: 24px;
  font-size: 0.95rem;
  outline: none;
  transition: border-color 0.2s;
}

.chat-input input:focus {
  border-color: #722F37;
}

.chat-input button {
  width: 44px;
  height: 44px;
  border: none;
  background: linear-gradient(135deg, #722F37 0%, #8B4513 100%);
  color: white;
  border-radius: 50%;
  cursor: pointer;
  font-size: 1.2rem;
  transition: transform 0.2s;
}

.chat-input button:hover:not(:disabled) {
  transform: scale(1.05);
}

.chat-input button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Mobile responsiveness */
@media (max-width: 480px) {
  .chat-panel {
    width: calc(100vw - 32px);
    height: calc(100vh - 100px);
    position: fixed;
    bottom: 16px;
    right: 16px;
  }
}
</style>
