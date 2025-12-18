import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { UserInfo } from '../api/types';

export const useUserStore = defineStore('user', () => {
  // State
  const userInfo = ref<UserInfo | null>(null);
  const accessToken = ref<string>('');

  // Getters
  const isLoggedIn = computed(() => !!accessToken.value);

  // Actions
  const setToken = (token: string) => {
    accessToken.value = token;
    localStorage.setItem('accessToken', token);
  };

  const setUserInfo = (info: UserInfo) => {
    userInfo.value = info;
    localStorage.setItem('userInfo', JSON.stringify(info));
  };

  const logout = () => {
    userInfo.value = null;
    accessToken.value = '';
    localStorage.removeItem('accessToken');
    localStorage.removeItem('userInfo');
  };

  const initFromStorage = () => {
    const token = localStorage.getItem('accessToken');
    const userStr = localStorage.getItem('userInfo');

    if (token) {
      accessToken.value = token;
    }

    if (userStr) {
      try {
        userInfo.value = JSON.parse(userStr);
      } catch (e) {
        console.error('Failed to parse user info from localStorage');
      }
    }
  };

  return {
    userInfo,
    accessToken,
    isLoggedIn,
    setToken,
    setUserInfo,
    logout,
    initFromStorage
  };
});
