<script setup lang="ts">
// SidebarMenuAction Component
</script>

<template>
  <button
    type="button"
    class="absolute left-1 top-1.5 flex aspect-square w-6 items-center justify-center rounded-lg p-0 text-slate-400 outline-none ring-slate-400 transition-transform hover:bg-slate-800 hover:text-white focus-visible:ring-2 cursor-pointer peer-hover/menu-button:text-slate-200"
  >
    <slot />
  </button>
</template>
