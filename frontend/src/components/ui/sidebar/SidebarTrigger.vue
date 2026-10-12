<script setup lang="ts">
import { useSidebar } from "./context";
import { PanelRightClose, PanelRightOpen, Menu } from "@lucide/vue";

const { toggleSidebar, isMobile, state } = useSidebar();
</script>

<template>
  <button
    type="button"
    @click="toggleSidebar"
    class="inline-flex h-10 px-3.5 items-center justify-center gap-2 rounded-2xl bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-200 hover:bg-orange-500 hover:text-white dark:hover:bg-orange-500 transition-all cursor-pointer shadow-sm active:scale-95 select-none shrink-0 border border-slate-200/80 dark:border-slate-700"
    :title="isMobile ? 'فتح القائمة الجانبية' : state === 'collapsed' ? 'توسيع القائمة الجانبية (Ctrl+B)' : 'طي القائمة الجانبية (Ctrl+B)'"
  >
    <Menu v-if="isMobile" class="h-4 w-4 text-orange-500 dark:text-orange-400" />
    <template v-else>
      <PanelRightOpen v-if="state === 'collapsed'" class="h-4 w-4 text-orange-500 dark:text-orange-400" />
      <PanelRightClose v-else class="h-4 w-4 text-slate-500 dark:text-slate-400" />
    </template>
    <span class="hidden sm:inline text-xs font-bold font-sans">
      {{ isMobile ? 'القائمة' : state === 'collapsed' ? 'توسيع القائمة' : 'طي القائمة' }}
    </span>
  </button>
</template>
