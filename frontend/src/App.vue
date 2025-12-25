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
  box-sizing: border-box;
}
</style>
