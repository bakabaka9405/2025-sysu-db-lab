<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui'
import { login } from '../api/user'
import { useUserStore } from '../stores/user'

const router = useRouter()
const message = useMessage()
const userStore = useUserStore()

const formRef = ref(null)
const loginLoading = ref(false)

const formValue = ref({
  usernameOrPhone: '',
  password: ''
})

/**
 * 表单验证规则
 * 支持用户名或手机号登录
 */
const rules = {
  usernameOrPhone: {
    required: true,
    message: '请输入用户名或手机号',
    trigger: ['blur', 'input'],
    validator: (_rule: any, value: string) => {
      if (!value) return false
      // 验证手机号格式
      const phoneRegex = /^1[3-9]\d{9}$/
      if (phoneRegex.test(value)) return true
      // 验证用户名格式
      const usernameRegex = /^[a-zA-Z0-9_]{3,20}$/
      if (usernameRegex.test(value)) return true
      return false
    }
  },
  password: {
    required: true,
    message: '请输入密码',
    trigger: ['blur', 'input']
  }
}

/**
 * 处理用户登录
 */
const handleLogin = async () => {
  if (loginLoading.value) return

  try {
    loginLoading.value = true
    const res = await login(formValue.value)

    if (!res.data) {
      throw new Error('未知错误');
    }

    // 保存token和用户信息到store
    userStore.setToken(res.data.accessToken)
    userStore.setUserInfo(res.data.user)

    message.success('登录成功')
    router.push('/')
  } catch (error: any) {
    message.error(error.response?.data?.message || '用户名或密码错误')
  } finally {
    setTimeout(() => {
      loginLoading.value = false
    }, 1000)
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-content">
      <div class="login-left">
        <NCard class="login-card" :bordered="false">
          <div class="login-header">
            <h2 class="login-title">欢迎回来</h2>
            <p class="login-subtitle">请登录您的账号</p>
          </div>

          <NForm ref="formRef" :model="formValue" label-placement="left" label-width="0"
            require-mark-placement="right-hanging" :rules="rules">
            <NFormItem path="usernameOrPhone">
              <NInput v-model:value="formValue.usernameOrPhone" placeholder="用户名或手机号" size="large" />
            </NFormItem>

            <NFormItem path="password">
              <NInput v-model:value="formValue.password" type="password" placeholder="密码" size="large"
                show-password-on="click" @keyup.enter="handleLogin" />
            </NFormItem>

            <div class="form-actions">
              <NButton type="primary" size="large" block :loading="loginLoading" @click="handleLogin">
                {{ loginLoading ? '登录中...' : '登录' }}
              </NButton>
            </div>

            <div class="form-links">
              <p class="link-text">
                还没有账号？
                <router-link to="/register" class="link">立即注册</router-link>
              </p>
            </div>
          </NForm>
        </NCard>
      </div>

      <div class="login-right">
        <div class="brand-content">
          <h1 class="brand-title">校园物流系统</h1>
          <p class="brand-description">便捷、高效、安全的校园物流服务平台</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #36ad6a 0%, #18a058 100%);
}

.login-content {
  display: flex;
  width: 1000px;
  height: 600px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  overflow: hidden;
}

.login-left {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px;
}

.login-card {
  width: 100%;
  max-width: 400px;
  background: transparent;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-title {
  font-size: 24px;
  color: #333;
  margin: 0 0 8px;
  font-weight: 600;
}

.login-subtitle {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.form-actions {
  margin-top: 32px;
}

.form-links {
  margin-top: 20px;
  text-align: center;
}

.link-text {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.link {
  color: #18a058;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.3s ease;
}

.link:hover {
  color: #0c7a43;
  text-decoration: underline;
}

.login-right {
  flex: 1;
  background: linear-gradient(135deg, #36ad6a 0%, #18a058 100%);
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px;
  color: white;
  text-align: center;
}

.brand-title {
  font-size: 48px;
  font-weight: bold;
  margin: 0 0 16px;
  letter-spacing: 2px;
}

.brand-description {
  font-size: 18px;
  margin: 0;
  opacity: 0.95;
}
</style>
