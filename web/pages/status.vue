<template>
    <client-only>
    <UContainer class="flex justify-center items-center h-screen">
        <StatusPanel />
    </UContainer>
    <Notification />
</client-only>
</template>

<script setup lang="ts">
import StatusPanel from '@/components/index/StatusPanel.vue'
import Notification from '~/components/widgets/Notification.vue';
</script>

<style scoped>
body {
    background-color: #f4f1ec;
}
</style>