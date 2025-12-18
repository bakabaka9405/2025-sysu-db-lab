<script setup lang="ts">
import { NMessageProvider, NDialogProvider, NNotificationProvider } from 'naive-ui'
</script>

<template>
  <NMessageProvider>
    <NDialogProvider>
      <NNotificationProvider>
        <router-view />
      </NNotificationProvider>
    </NDialogProvider>
  </NMessageProvider>
</template>

<style>
html,
body {
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
    'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol',
    'Noto Color Emoji';
}

#app {
  width: 100%;
  height: 100%;
}
</style>
