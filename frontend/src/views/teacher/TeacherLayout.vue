<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { apiClient } from "@/api/client";
import {
  LayoutDashboard,
  Upload,
  MessageSquareQuote,
  Coins,
  ArrowRight,
  Menu,
  GraduationCap,
} from "@lucide/vue";

const isSidebarOpen = ref(false);
const route = useRoute();
const router = useRouter();

const teacherName = ref("");
const openQuestionsCount = ref(0);

onMounted(async () => {
  try {
    const [meRes, questionsRes] = await Promise.all([
      apiClient.get<any>("/users/me"),
      apiClient.get<any[]>("/questions/teacher"),
    ]);
    const u = meRes.data || {};
    teacherName.value = u.name || `${u.first_name || ""} ${u.last_name || ""}`.trim() || u.username || "";
    const questions = Array.isArray(questionsRes.data) ? questionsRes.data : [];
    openQuestionsCount.value = questions.filter((q) => q.status !== "answered").length;
  } catch {
    // keep defaults empty
  }
});

const navItems = computed(() => [
  {
    name: "نظرة عامة على المقررات",
    to: "/teacher",
    icon: LayoutDashboard,
    badge: null,
  },
  {
    name: "رفع المحاضرات والملازم",
    to: "/teacher/uploads",
    icon: Upload,
    badge: "رفع",
    badgeColor: "orange",
  },
  {
    name: "أسئلة واستفسارات الطلاب",
    to: "/teacher/questions",
    icon: MessageSquareQuote,
    badge: openQuestionsCount.value > 0 ? `${openQuestionsCount.value} أسئلة` : null,
    badgeColor: "orange",
  },
  {
    name: "أرباحي ومستحقاتي المالية",
    to: "/teacher/earnings",
    icon: Coins,
    badge: null,
  },
]);

const currentPath = computed(() => route.path);
</script>

<template>
  <div dir="rtl" class="min-h-screen bg-white dark:bg-[#0F0F0F] text-[#111111] dark:text-[#FFFFFF] flex transition-colors duration-200">
    <!-- Teacher Sidebar -->
    <aside
      :class="[
        'w-72 bg-[#111111] dark:bg-[#181818] text-[#FFFFFF] flex flex-col border-l border-[#27272A] shrink-0 h-screen sticky top-0 transition-all duration-300 z-40',
        'fixed inset-y-0 right-0 lg:static lg:block',
        isSidebarOpen ? 'translate-x-0' : 'translate-x-full lg:translate-x-0'
      ]"
    >
      <!-- Header / Teacher Profile Pill -->
      <div class="p-5 border-b border-[#27272A] flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-2xl bg-orange-500 flex items-center justify-center text-white font-black shadow-md shadow-orange-500/20">
            <GraduationCap class="w-5 h-5" />
          </div>
          <div>
            <h1 class="font-black text-sm text-white tracking-tight flex items-center gap-1.5">
              بوابة الأستاذ
              <span class="text-[9px] font-bold px-1.5 py-0.5 rounded-md bg-orange-500/20 text-orange-400 border border-orange-500/30 font-mono">
                TEACHER
              </span>
            </h1>
            <p class="text-xs text-[#A1A1AA] font-medium truncate">{{ teacherName || 'المحاضر' }}</p>
          </div>
        </div>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 overflow-y-auto px-4 py-6 space-y-1.5">
        <div class="text-[11px] font-bold text-[#71717A] uppercase tracking-wider px-3 mb-2">
          لوحة تحكم المساقات التعليمية
        </div>

        <button
          v-for="item in navItems"
          :key="item.to"
          type="button"
          @click="router.push(item.to); isSidebarOpen = false"
          :class="[
            'w-full flex items-center justify-between px-3.5 py-3 rounded-2xl text-xs font-bold transition-all group cursor-pointer text-right',
            currentPath === item.to
              ? 'bg-orange-500 text-white shadow-md shadow-orange-500/25'
              : 'text-[#A1A1AA] hover:text-white hover:bg-[#222222]'
          ]"
        >
          <div class="flex items-center gap-3">
            <component
              :is="item.icon"
              :class="[
                'w-4 h-4 transition-transform duration-200 group-hover:scale-110',
                currentPath === item.to ? 'text-white' : 'text-[#71717A] group-hover:text-orange-400'
              ]"
            />
            <span>{{ item.name }}</span>
          </div>

          <span
            v-if="item.badge"
            :class="[
              'text-[10px] font-bold px-2 py-0.5 rounded-md',
              currentPath === item.to ? 'bg-white/20 text-white' : 'bg-orange-500/10 text-orange-400 border border-orange-500/20'
            ]"
          >
            {{ item.badge }}
          </span>
        </button>
      </nav>

      <!-- Footer Action -->
      <div class="p-4 border-t border-[#27272A] space-y-2">
        <button
          type="button"
          @click="router.push('/'); isSidebarOpen = false"
          class="w-full flex items-center justify-center gap-2 px-4 py-2.5 rounded-2xl bg-[#222222] hover:bg-[#2A2A2A] text-slate-200 text-xs font-bold transition-all border border-[#27272A] cursor-pointer"
        >
          <ArrowRight class="w-4 h-4" />
          <span>العودة لتطبيق الطلاب</span>
        </button>
      </div>
    </aside>

    <!-- Mobile Backdrop -->
    <div
      v-if="isSidebarOpen"
      @click="isSidebarOpen = false"
      class="fixed inset-0 bg-black/60 backdrop-blur-sm z-30 lg:hidden"
    ></div>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <!-- Top Bar -->
      <header class="h-16 bg-white/90 dark:bg-[#181818]/90 backdrop-blur-md border-b border-[#E5E7EB] dark:border-[#27272A] px-6 sm:px-8 flex items-center justify-between sticky top-0 z-30">
        <div class="flex items-center gap-3">
          <button
            type="button"
            @click="isSidebarOpen = !isSidebarOpen"
            class="lg:hidden p-2 rounded-xl bg-[#F8F8F8] dark:bg-[#222222] text-[#111111] dark:text-[#FFFFFF] cursor-pointer"
          >
            <Menu class="w-5 h-5" />
          </button>
          <div class="text-xs font-bold text-[#6B7280] dark:text-[#A1A1AA]">
            لوحة الأستاذ المحاضر / {{ teacherName || 'المحاضر' }}
          </div>
        </div>

        <div class="flex items-center gap-3">
        </div>
      </header>

      <main class="flex-1 p-5 sm:p-8 overflow-y-auto max-w-7xl w-full mx-auto">
        <RouterView v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </RouterView>
      </main>
    </div>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
