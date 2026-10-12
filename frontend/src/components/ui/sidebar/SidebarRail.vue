<script setup lang="ts">
import { useSidebar } from "./context";

const { toggleSidebar } = useSidebar();
</script>

<template>
  <button
    type="button"
    @click="toggleSidebar"
    aria-label="Toggle Sidebar"
    title="تبديل القائمة الجانبية (انقر أو اسحب)"
    tabindex="-1"
    class="absolute inset-y-0 z-40 hidden w-3 transition-all ease-linear after:absolute after:inset-y-0 after:left-1/2 after:w-[3px] hover:after:bg-orange-500/60 group-data-[side=left]/sidebar:right-0 group-data-[side=left]/sidebar:translate-x-1/2 group-data-[side=right]/sidebar:left-0 group-data-[side=right]/sidebar:-translate-x-1/2 sm:flex cursor-col-resize border-none outline-none bg-transparent"
  />
</template>
