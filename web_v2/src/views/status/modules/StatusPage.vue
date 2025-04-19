<template>
  <div class="flex justify-center items-center h-screen">
    <div class="h-1/2">
      <h1 class="text-5xl italic font-bold text-center text-gray-300 mb-10">Ech0s Status</h1>
      <div class="">
        <!-- TheStatus -->
        <TheStatus />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import TheStatus from './TheStatus.vue'
</script>
