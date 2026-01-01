<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import {
  NCard,
  NForm,
  NFormItem,
  NInput,
  NButton,
  NSpace,
  NDescriptions,
  NDescriptionsItem,
  NTag,
  useMessage
} from 'naive-ui'
import { getProfile, updateProfile } from '../api/user'
import type { UserProfile } from '../api/types'

const router = useRouter()
const userStore = useUserStore()
const message = useMessage()

const loading = ref(false)
const editMode = ref(false)
const profile = ref<UserProfile | null>(null)

// 编辑表单
const formValue = ref({
  realName: '',
  email: ''
})

// 角色映射
const roleMap: Record<string, { text: string; type: any }> = {
  user: { text: '用户', type: 'default' },
  staff: { text: '员工', type: 'info' },
  admin: { text: '管理员', type: 'success' }
}

// 加载用户信息
const loadProfile = async () => {
  loading.value = true
  const res = await getProfile()
  if (res.code === 0 && res.data) {
    profile.value = res.data
    formValue.value = {
      realName: res.data.realName,
      email: res.data.email || ''
    }
  } else {
    message.error(res.message || '加载失败')
  }
  loading.value = false
}

// 进入编辑模式
const handleEdit = () => {
  if (profile.value) {
    formValue.value = {
      realName: profile.value.realName,
      email: profile.value.email || ''
    }
  }
  editMode.value = true
}

// 取消编辑
const handleCancel = () => {
  editMode.value = false
}

// 保存修改
const handleSave = async () => {
  if (loading.value) return

  if (!formValue.value.realName) {
    message.error('请输入真实姓名')
    return
  }

  loading.value = true
  const res = await updateProfile(formValue.value)
  if (res.code === 0) {
    message.success('更新成功')
    editMode.value = false
    // 更新本地用户信息
    if (profile.value) {
      profile.value.realName = formValue.value.realName
      profile.value.email = formValue.value.email
      userStore.setUserInfo(profile.value)
    }
    loadProfile()
  } else {
    message.error(res.message || '更新失败')
  }
  loading.value = false
}

// 退出登录
const handleLogout = () => {
  userStore.logout()
  message.success('已退出登录')
  router.push('/')
}

onMounted(() => {
  loadProfile()
})
</script>

<template>
  <div class="page-container">
    <NCard title="个人中心" style="max-width: 800px; margin: 0 auto;">
      <div v-if="!editMode && profile">
        <NDescriptions label-placement="left" :column="1" bordered>
          <NDescriptionsItem label="用户名">
            {{ profile.username }}
          </NDescriptionsItem>
          <NDescriptionsItem label="真实姓名">
            {{ profile.realName }}
          </NDescriptionsItem>
          <NDescriptionsItem label="手机号">
            {{ profile.phone }}
          </NDescriptionsItem>
          <NDescriptionsItem label="邮箱">
            {{ profile.email || '-' }}
          </NDescriptionsItem>
          <NDescriptionsItem v-if="profile.idCard" label="身份证号">
            {{ profile.idCard }}
          </NDescriptionsItem>
          <NDescriptionsItem label="角色">
            <NTag :type="roleMap[profile.role]?.type || 'default'">
              {{ roleMap[profile.role]?.text || profile.role }}
            </NTag>
          </NDescriptionsItem>
        </NDescriptions>

        <NSpace style="margin-top: 24px;" justify="center">
          <NButton type="primary" @click="handleEdit">
            编辑资料
          </NButton>
          <NButton @click="handleLogout">
            退出登录
          </NButton>
        </NSpace>
      </div>

      <div v-else>
        <NForm :model="formValue" label-placement="left" label-width="100">
          <NFormItem label="真实姓名" required>
            <NInput v-model:value="formValue.realName" placeholder="请输入真实姓名" />
          </NFormItem>

          <NFormItem label="邮箱">
            <NInput v-model:value="formValue.email" placeholder="请输入邮箱" />
          </NFormItem>
        </NForm>

        <NSpace style="margin-top: 24px;" justify="center">
          <NButton @click="handleCancel">
            取消
          </NButton>
          <NButton type="primary" :loading="loading" @click="handleSave">
            保存
          </NButton>
        </NSpace>
      </div>
    </NCard>
  </div>
</template>

<style scoped>
.page-container {
  padding: 24px;
  min-height: 100vh;
  background-color: #f5f5f5;
}
</style>
