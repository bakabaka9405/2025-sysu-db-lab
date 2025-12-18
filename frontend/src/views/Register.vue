<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui'
import { register } from '../api/user'

const router = useRouter()
const message = useMessage()
const formRef = ref(null)
const registerLoading = ref(false)

const formValue = ref({
  username: '',
  password: '',
  confirmPassword: '',
  realName: '',
  phone: '',
  idCard: '',
  email: ''
})

// 表单验证状态变量
let usernameValid = false
let passwordValid = false
let confirmPasswordValid = false
let realNameValid = false
let phoneValid = false
let idCardValid = true // 可选字段，默认有效
let emailValid = true // 可选字段，默认有效

/**
 * 表单验证规则
 */
const rules = {
  username: {
    required: true,
    message: '请输入用户名（3-20位字母、数字或下划线）',
    trigger: ['blur', 'input'],
    validator: (_rule: any, value: string) => {
      if (!value) return false
      const usernameRegex = /^[a-zA-Z0-9_]{3,20}$/
      return usernameValid = usernameRegex.test(value)
    }
  },
  password: {
    required: true,
    message: '请输入合法的密码（6-20位）',
    trigger: ['blur', 'input'],
    validator: (_rule: any, value: string) => {
      return passwordValid = (value.length >= 6 && value.length <= 20)
    }
  },
  confirmPassword: {
    required: true,
    message: '两次输入的密码不一致',
    trigger: ['blur', 'input'],
    validator: (_rule: any, value: string) => {
      return confirmPasswordValid = (value === formValue.value.password)
    }
  },
  realName: {
    required: true,
    message: '请输入真实姓名',
    trigger: ['blur', 'input'],
    validator: (_rule: any, value: string) => {
      return realNameValid = (value.length >= 2 && value.length <= 20)
    }
  },
  phone: {
    required: true,
    message: '请输入正确的手机号（1开头的11位数字）',
    trigger: ['blur', 'input'],
    validator: (_rule: any, value: string) => {
      const phoneRegex = /^1[3-9]\d{9}$/
      return phoneValid = phoneRegex.test(value)
    }
  },
  idCard: {
    required: false,
    message: '请输入正确的身份证号（18位）',
    trigger: ['blur', 'input'],
    validator: (_rule: any, value: string) => {
      if (!value) {
        idCardValid = true
        return true
      }
      const idCardRegex = /^[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$/
      return idCardValid = idCardRegex.test(value)
    }
  },
  email: {
    required: false,
    message: '请输入正确的邮箱地址',
    trigger: ['blur', 'input'],
    validator: (_rule: any, value: string) => {
      if (!value) {
        emailValid = true
        return true
      }
      const emailRegex = /^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$/
      return emailValid = emailRegex.test(value)
    }
  }
}

/**
 * 处理用户注册
 */
const handleRegister = async () => {
  if (registerLoading.value) return

  // 验证所有字段
  if (!usernameValid || !passwordValid || !confirmPasswordValid || !realNameValid || !phoneValid || !idCardValid || !emailValid) {
    message.error('请检查注册信息是否符合要求')
    return
  }

  try {
    registerLoading.value = true
    await register({
      username: formValue.value.username,
      password: formValue.value.password,
      realName: formValue.value.realName,
      phone: formValue.value.phone,
      idCard: formValue.value.idCard || undefined,
      email: formValue.value.email || undefined
    })
    message.success('注册成功，即将跳转到登录页')
    setTimeout(() => {
      router.push('/login')
    }, 1500)
  } catch (error: any) {
    message.error(error.response?.data?.message || '注册失败，请重试')
  } finally {
    setTimeout(() => {
      registerLoading.value = false
    }, 1000)
  }
}
</script>

<template>
  <div class="register-container">
    <div class="register-content">
      <div class="register-left">
        <NCard class="register-card" :bordered="false">
          <div class="register-header">
            <h2 class="register-title">创建账号</h2>
            <p class="register-subtitle">请填写以下信息完成注册</p>
          </div>

          <NForm
            ref="formRef"
            :model="formValue"
            label-placement="left"
            label-width="0"
            require-mark-placement="right-hanging"
            :rules="rules"
          >
            <NFormItem path="username">
              <NInput
                v-model:value="formValue.username"
                placeholder="用户名"
                size="large"
              />
            </NFormItem>

            <NFormItem path="realName">
              <NInput
                v-model:value="formValue.realName"
                placeholder="真实姓名"
                size="large"
              />
            </NFormItem>

            <NFormItem path="phone">
              <NInput
                v-model:value="formValue.phone"
                placeholder="手机号"
                size="large"
              />
            </NFormItem>

            <NFormItem path="password">
              <NInput
                v-model:value="formValue.password"
                type="password"
                placeholder="密码"
                size="large"
                show-password-on="click"
              />
            </NFormItem>

            <NFormItem path="confirmPassword">
              <NInput
                v-model:value="formValue.confirmPassword"
                type="password"
                placeholder="确认密码"
                size="large"
                show-password-on="click"
              />
            </NFormItem>

            <NFormItem path="idCard">
              <NInput
                v-model:value="formValue.idCard"
                placeholder="身份证号（选填）"
                size="large"
              />
            </NFormItem>

            <NFormItem path="email">
              <NInput
                v-model:value="formValue.email"
                placeholder="邮箱（选填）"
                size="large"
              />
            </NFormItem>

            <div class="form-actions">
              <NButton
                type="primary"
                size="large"
                block
                :loading="registerLoading"
                @click="handleRegister"
              >
                {{ registerLoading ? '注册中...' : '注册' }}
              </NButton>
            </div>

            <div class="form-links">
              <p class="link-text">
                已有账号？
                <router-link to="/login" class="link">立即登录</router-link>
              </p>
            </div>
          </NForm>
        </NCard>
      </div>

      <div class="register-right">
        <div class="brand-content">
          <h1 class="brand-title">校园物流系统</h1>
          <p class="brand-description">便捷、高效、安全的校园物流服务平台</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px 0;
}

.register-content {
  display: flex;
  width: 1000px;
  min-height: 700px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  overflow: hidden;
}

.register-left {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px;
}

.register-card {
  width: 100%;
  max-width: 400px;
  background: transparent;
}

.register-header {
  text-align: center;
  margin-bottom: 24px;
}

.register-title {
  font-size: 24px;
  color: #333;
  margin: 0 0 8px;
  font-weight: 600;
}

.register-subtitle {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.form-actions {
  margin-top: 24px;
}

.form-links {
  margin-top: 16px;
  text-align: center;
}

.link-text {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.link {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.3s ease;
}

.link:hover {
  color: #5568d3;
  text-decoration: underline;
}

.register-right {
  flex: 1;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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

/* 强制输入框文本左对齐 */
:deep(.n-input__input-el) {
  text-align: left !important;
}

:deep(.n-input__textarea-el) {
  text-align: left !important;
}
</style>
