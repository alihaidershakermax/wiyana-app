<script setup lang="ts">
import { computed } from "vue";
import { useSidebar } from "./context";
import { X } from "@lucide/vue";

const props = withDefaults(
  defineProps<{
    side?: "left" | "right";
    variant?: "sidebar" | "floating" | "inset";
    collapsible?: "offcanvas" | "icon" | "none";
    className?: string;
  }>(),
  {
    side: "right",
    variant: "sidebar",
    collapsible: "icon",
  }
);

const { state, openMobile, setOpenMobile, isMobile } = useSidebar();
</script>

<template>
  <!-- 1. Mobile Drawer Mode (Sheet Overlay) -->
  <div v-if="isMobile">
    <!-- Backdrop overlay -->
    <transition name="fade">
      <div
        v-if="openMobile"
        @click="setOpenMobile(false)"
        class="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-sm cursor-pointer transition-opacity"
      ></div>
    </transition>

    <!-- Mobile Sheet Panel -->
    <aside
      :class="[
        'fixed inset-y-0 z-50 flex h-full w-[var(--sidebar-width-mobile)] max-w-[85vw] flex-col bg-slate-900 border-slate-800 text-slate-100 shadow-2xl transition-all duration-300 ease-in-out',
        side === 'right' ? 'right-0 border-l' : 'left-0 border-r',
        openMobile
          ? 'translate-x-0 opacity-100 pointer-events-auto'
          : side === 'right'
          ? 'translate-x-full opacity-0 pointer-events-none'
          : '-translate-x-full opacity-0 pointer-events-none'
      ]"
    >
      <!-- Mobile Close Trigger -->
      <div class="absolute top-4 left-4 z-10">
        <button
          type="button"
          @click="setOpenMobile(false)"
          class="p-2 rounded-xl bg-slate-800/90 text-slate-400 hover:text-white transition-colors cursor-pointer"
          title="إغلاق القائمة"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <slot />
    </aside>
  </div>

  <!-- 2. Desktop Mode (Collapsible Sidebar) -->
  <aside
    v-else-if="collapsible !== 'offcanvas' || state === 'expanded'"
    :data-state="state"
    :data-collapsible="collapsible"
    :data-side="side"
    :data-variant="variant"
    :class="[
      'group/sidebar relative z-30 hidden lg:flex flex-col bg-slate-900 text-slate-100 border-slate-800 transition-all duration-300 ease-in-out select-none shrink-0 overflow-visible',
      side === 'right' ? 'border-l' : 'border-r',
      variant === 'floating'
        ? 'm-3 rounded-3xl border shadow-xl h-[calc(100vh-1.5rem)]'
        : 'h-screen sticky top-0',
      state === 'collapsed' && collapsible === 'icon'
        ? 'w-[var(--sidebar-width-icon)]'
        : 'w-[var(--sidebar-width)]',
      className
    ]"
  >
    <slot />
  </aside>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>

