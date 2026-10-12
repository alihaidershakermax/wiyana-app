<script setup lang="ts">
import { ref, computed } from "vue";
import { useQuery, useQueryClient } from "@tanstack/vue-query";
import { CheckCheck, Tag, BookOpen, Info } from "@lucide/vue";
import { fetchMyNotifications, markNotificationRead, markAllNotificationsRead } from "@/api/notifications";
import { useUserStore } from "@/stores/useUserStore";
import SkeletonCard from "@/components/ui/SkeletonCard.vue";
import EmptyState from "@/components/ui/EmptyState.vue";
import ErrorState from "@/components/ui/ErrorState.vue";

const queryClient = useQueryClient();
const userStore = useUserStore();
const activeFilter = ref<"all" | "unread" | "read">("all");

const { data: rawNotifications, isLoading, isError, refetch } = useQuery({
  queryKey: ["myNotifications"],
  queryFn: fetchMyNotifications,
  staleTime: 1000 * 60 * 2,
});

const notifications = computed(() => {
  const list = Array.isArray(rawNotifications.value) ? rawNotifications.value : [];
  if (activeFilter.value === "unread") {
    return list.filter((n: any) => !n.is_read && !n.isRead);
  }
  if (activeFilter.value === "read") {
    return list.filter((n: any) => n.is_read || n.isRead);
  }
  return list;
});

const unreadCount = computed(() => {
  const list = Array.isArray(rawNotifications.value) ? rawNotifications.value : [];
  return list.filter((n: any) => !n.is_read && !n.isRead).length;
});

async function handleMarkRead(notifId: string) {
  try {
    await markNotificationRead(notifId);
    queryClient.invalidateQueries({ queryKey: ["myNotifications"] });
  } catch {}
}

async function handleMarkAllRead() {
  try {
    await markAllNotificationsRead();
    userStore.showToast("تم تحديد جميع الإشعارات كمقروءة");
    queryClient.invalidateQueries({ queryKey: ["myNotifications"] });
    queryClient.invalidateQueries({ queryKey: ["homeDashboardData"] });
  } catch (err) {
    userStore.showToast("تعذر تحديث الإشعارات", "error");
  }
}
</script>

<template>
  <div class="min-h-screen bg-white dark:bg-[#0F0F0F] text-[#111111] dark:text-[#FFFFFF] flex flex-col font-sans transition-colors duration-200 pb-28 text-start">
    <!-- Header -->
    <header class="sticky top-0 z-30 w-full backdrop-blur-md bg-white/95 dark:bg-[#0F0F0F]/95 border-b border-[#E5E7EB] dark:border-[#27272A]">
      <div class="max-w-7xl mx-auto px-4 py-3.5 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <h1 class="text-xl sm:text-2xl font-black text-[#111111] dark:text-[#FFFFFF]">
            الإشعارات
          </h1>
          <span
            v-if="unreadCount > 0"
            class="px-2.5 py-0.5 rounded-full text-xs font-black bg-orange-500 text-white"
          >
            {{ unreadCount }} جديد
          </span>
        </div>

        <button
          @click="handleMarkAllRead"
          type="button"
          class="inline-flex items-center gap-1.5 text-xs font-bold text-orange-500 hover:text-[#EA580C] px-3 py-1.5 rounded-xl hover:bg-orange-500/10 active:scale-95 transition-all cursor-pointer"
        >
          <CheckCheck class="w-4 h-4" />
          <span>تحديد الكل كمقروء</span>
        </button>
      </div>

      <!-- Filter Tabs -->
      <div class="max-w-7xl mx-auto px-4 flex items-center gap-2 pb-2.5">
        <button
          v-for="tab in [
            { id: 'all', label: 'الكل' },
            { id: 'unread', label: 'غير المقروءة' },
            { id: 'read', label: 'المقروءة' },
          ]"
          :key="tab.id"
          @click="activeFilter = tab.id as any"
          :class="[
            'px-3 py-1 rounded-xl text-xs font-bold transition-all cursor-pointer',
            activeFilter === tab.id
              ? 'bg-orange-500 text-white shadow-xs'
              : 'bg-[#F8F8F8] dark:bg-[#181818] text-[#6B7280] dark:text-[#A1A1AA] hover:text-[#111111] dark:hover:text-white border border-[#E5E7EB] dark:border-[#27272A]'
          ]"
        >
          {{ tab.label }}
        </button>
      </div>
    </header>

    <!-- Loading State -->
    <div v-if="isLoading" class="max-w-7xl mx-auto px-4 py-4 w-full flex flex-col gap-3">
      <SkeletonCard v-for="i in 4" :key="i" />
    </div>

    <!-- Error State -->
    <div v-else-if="isError" class="max-w-7xl mx-auto px-4 py-8 w-full">
      <ErrorState title="تعذر تحميل الإشعارات" message="حدث خطأ أثناء الاتصال بالخادم." @retry="refetch()" />
    </div>

    <!-- Empty State -->
    <div v-else-if="notifications.length === 0" class="max-w-7xl mx-auto px-4 py-12 w-full">
      <EmptyState
        title="لا توجد إشعارات حالياً"
        description="ستصلك هنا تنبيهات المحاضرات الجديدة والتحديثات فور نشرها."
      />
    </div>

    <!-- Main Notifications List -->
    <main v-else class="flex-1 max-w-7xl mx-auto px-4 py-4 w-full flex flex-col gap-3">
      <div
        v-for="item in notifications"
        :key="item.id || item._id"
        @click="handleMarkRead(String(item.id || item._id))"
        :class="[
          'p-4 rounded-3xl border transition-all duration-200 flex items-start gap-3.5 cursor-pointer active:scale-[0.99] shadow-xs',
          (item.is_read || item.isRead)
            ? 'bg-white dark:bg-[#181818] border-[#E5E7EB] dark:border-[#27272A] opacity-80'
            : 'bg-white dark:bg-[#181818] border-orange-500/40 ring-1 ring-orange-500/20'
        ]"
      >
        <div
          class="w-10 h-10 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center shrink-0 mt-0.5"
        >
          <BookOpen v-if="item.type === 'course'" class="w-5 h-5" />
          <Tag v-else-if="item.type === 'promo'" class="w-5 h-5" />
          <Info v-else class="w-5 h-5" />
        </div>

        <div class="flex-1 flex flex-col gap-1 min-w-0">
          <div class="flex items-center justify-between gap-2">
            <h3 class="text-xs sm:text-sm font-bold text-[#111111] dark:text-[#FFFFFF] truncate">
              {{ item.title }}
            </h3>
            <span class="text-[10px] text-[#9CA3AF] dark:text-[#71717A] shrink-0">
              {{ item.created_at ? new Date(item.created_at).toLocaleDateString('ar-IQ') : 'اليوم' }}
            </span>
          </div>

          <p class="text-xs text-[#6B7280] dark:text-[#A1A1AA] leading-relaxed">
            {{ item.message }}
          </p>
        </div>
      </div>
    </main>
  </div>
</template>
