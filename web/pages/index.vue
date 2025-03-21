<template>
  <UContainer class="py-2 pb-4">
    <!-- <HeaderBar /> -->
    <AddForm />
    <MessageList />
  </UContainer>
  <Notification />
</template>

<script setup lang="ts">
// import HeaderBar from '@/components/index/HeaderBar.vue'
import AddForm from '@/components/index/AddForm.vue'
import MessageList from '@/components/index/MessageList.vue'
import Notification from '~/components/widgets/Notification.vue';
</script>

<style>
body {
  background-color: #f4f1ec;
}
</style>