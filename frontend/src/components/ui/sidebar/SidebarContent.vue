<script setup lang="ts">
// SidebarContent Component
</script>

<template>
  <div class="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto px-3 py-4 custom-scrollbar">
    <slot />
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 4px;
}
</style>
