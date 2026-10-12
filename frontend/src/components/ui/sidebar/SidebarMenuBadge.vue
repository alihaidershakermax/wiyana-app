<script setup lang="ts">
// SidebarMenuBadge Component
</script>

<template>
  <span
    class="pointer-events-none inline-flex h-5 select-none items-center justify-center rounded-md px-1.5 font-mono text-[10px] font-black shrink-0 transition-all group-data-[state=collapsed]/sidebar:hidden"
  >
    <slot />
  </span>
</template>

