<script setup lang="ts">
import { useRoute } from "vue-router";
import { Construction, ArrowRight } from "@lucide/vue";

const route = useRoute();
const title = computedTitle();

function computedTitle(): string {
  const name = String(route.name ?? "");
  const map: Record<string, string> = {
    "admin.analytics": "التحليلات",
    "admin.courses": "المساقات",
    "admin.students": "الطلاب",
    "admin.teachers": "الأساتذة",
    "admin.subscriptions": "الاشتراكات",
    "admin.notifications": "الإشعارات",
    "admin.settings": "الإعدادات",
    "admin.teacher.courses": "مساقاتي",
    "admin.teacher.lessons": "محاضراتي",
    "admin.teacher.upload": "رفع محاضرة",
    "admin.teacher.students": "طلابي",
    "admin.teacher.profile": "الملف الشخصي",
  };
  return map[name] ?? "قيد الإنشاء";
}
</script>

<template>
  <div dir="rtl" class="flex flex-col items-center justify-center min-h-[50vh] text-center space-y-4">
    <div
      class="h-16 w-16 rounded-2xl bg-orange-50 dark:bg-orange-950/40 border border-orange-200 dark:border-orange-900/60 flex items-center justify-center"
    >
      <Construction class="h-8 w-8 text-orange-500" />
    </div>
    <div class="space-y-1">
      <h1 class="text-lg font-extrabold">{{ title }}</h1>
      <p class="text-sm text-slate-500 dark:text-slate-400">
        هذه الصفحة قيد الإنشاء وستتوفر في المراحل القادمة
      </p>
    </div>
    <button class="btn-secondary" @click="$router.back()">
      <ArrowRight class="h-4 w-4" />
      العودة
    </button>
  </div>
</template>
