<script setup lang="ts">
// SidebarGroupLabel Component
</script>

<template>
  <div
    class="duration-200 flex h-8 shrink-0 items-center rounded-md px-2 text-[11px] font-black text-slate-500 uppercase tracking-wider select-none outline-none group-data-[state=collapsed]/sidebar:opacity-0 group-data-[state=collapsed]/sidebar:h-0 group-data-[state=collapsed]/sidebar:my-0 group-data-[state=collapsed]/sidebar:py-0 overflow-hidden transition-all"
  >
    <slot />
  </div>
</template>

