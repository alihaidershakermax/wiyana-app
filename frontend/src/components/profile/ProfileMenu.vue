<script setup lang="ts">
import {
  Settings,
  Headphones,
  Info,
  ShieldCheck,
} from "@lucide/vue";
import MenuItem from "./MenuItem.vue";

const emit = defineEmits<{
  (e: "navigate", routeKey: string): void;
}>();

const menuItems = [
  { key: "settings", label: "إعدادات المنصة والتطبيق", icon: Settings },
  { key: "support", label: "الدعم الفني المباشر", icon: Headphones },
  { key: "about", label: "حول تطبيق ويانا", icon: Info },
  { key: "privacy", label: "سياسة الخصوصية", icon: ShieldCheck },
];
</script>

<template>
  <div class="flex flex-col w-full">
    <div class="flex flex-col bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800/80 rounded-3xl overflow-hidden shadow-xs divide-y divide-slate-100 dark:divide-slate-800/60">
      <MenuItem
        v-for="item in menuItems"
        :key="item.key"
        :label="item.label"
        :icon="item.icon"
        @click="emit('navigate', item.key)"
      />
    </div>
  </div>
</template>

