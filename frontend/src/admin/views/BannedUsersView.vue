<script setup lang="ts">
import { ref, onMounted } from "vue";
import { fetchBannedUsers, unbanUserApi, banUserApi, type BannedUserInfo } from "@/api/protectionAdmin";

const bannedUsers = ref<BannedUserInfo[]>([]);
const loading = ref(true);
const actionLoading = ref<string | null>(null);
const manualUserId = ref("");
const alertMessage = ref<{ text: string; type: "success" | "error" } | null>(null);

async function loadData() {
  loading.value = true;
  alertMessage.value = null;
  try {
    const res = await fetchBannedUsers();
    bannedUsers.value = res?.banned_users || [];
  } catch (err: any) {
    bannedUsers.value = [];
  } finally {
    loading.value = false;
  }
}

async function handleUnban(userId: string) {
  if (!confirm(`هل أنت تأكد من رغبتك في رفع الحظر عن المستخدم (${userId})؟`)) return;
  actionLoading.value = userId;
  alertMessage.value = null;
  try {
    const res = await unbanUserApi(userId);
    alertMessage.value = { text: res.message || "تم رفع الحظر بنجاح", type: "success" };
    await loadData();
  } catch (err: any) {
    alertMessage.value = { text: "فشل رفع الحظر عن المستخدم", type: "error" };
  } finally {
    actionLoading.value = null;
  }
}

async function handleManualBan() {
  if (!manualUserId.value.trim()) return;
  actionLoading.value = "manual";
  alertMessage.value = null;
  try {
    const res = await banUserApi(manualUserId.value.trim());
    alertMessage.value = { text: res.message || "تم حظر المستخدم بنجاح", type: "success" };
    manualUserId.value = "";
    await loadData();
  } catch (err: any) {
    alertMessage.value = { text: "فشل حظر المستخدم", type: "error" };
  } finally {
    actionLoading.value = null;
  }
}

function formatUntilTime(untilTs?: number): string {
  if (!untilTs) return "غير محدد";
  const diffSecs = Math.max(0, Math.ceil(untilTs - Date.now() / 1000));
  const mins = Math.floor(diffSecs / 60);
  const secs = diffSecs % 60;
  return `${String(mins).padStart(2, "0")}:${String(secs).padStart(2, "0")}`;
}

onMounted(() => {
  loadData();
});
</script>

<template>
  <div class="p-6 max-w-7xl mx-auto space-y-6" dir="rtl">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-gray-200 dark:border-gray-800 pb-5">
      <div>
        <h1 class="text-2xl font-black text-gray-900 dark:text-white flex items-center gap-2">
          <span>🛡️ إدارة المحظورين وسجل الحماية</span>
        </h1>
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
          مراقبة الطلاب المحظورين بسبب محاولات أخذ لقطات شاشة أو تسجيل فيديو، مع إمكانية رفع الحظر فوراً.
        </p>
      </div>

      <button
        @click="loadData"
        type="button"
        class="px-4 py-2 bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 text-gray-700 dark:text-gray-300 font-bold text-xs rounded-xl transition-all flex items-center gap-2 cursor-pointer self-start sm:self-auto"
      >
        <span>🔄 تحديث القائمة</span>
      </button>
    </div>

    <!-- Alert Box -->
    <div
      v-if="alertMessage"
      :class="[
        'p-4 rounded-xl text-xs font-bold border transition-all',
        alertMessage.type === 'success' ? 'bg-emerald-500/10 border-emerald-500/30 text-emerald-600 dark:text-emerald-400' : 'bg-red-500/10 border-red-500/30 text-red-600 dark:text-red-400',
      ]"
    >
      {{ alertMessage.text }}
    </div>

    <!-- Manual Ban Form Box -->
    <div class="bg-gray-50 dark:bg-[#1A1A1A] border border-gray-200 dark:border-gray-800 rounded-2xl p-5 shadow-sm">
      <h3 class="text-sm font-black text-gray-900 dark:text-white mb-3">حظر يدوي لمستخدم (بالـ Telegram ID):</h3>
      <div class="flex flex-col sm:flex-row gap-3">
        <input
          v-model="manualUserId"
          type="text"
          placeholder="أدخل Telegram ID أو User ID للفيسبوك/الموقع..."
          class="flex-1 px-4 py-2.5 bg-white dark:bg-[#0F0F0F] border border-gray-300 dark:border-gray-700 rounded-xl text-xs font-mono text-gray-900 dark:text-white focus:outline-none focus:border-orange-500"
        />
        <button
          @click="handleManualBan"
          :disabled="actionLoading === 'manual' || !manualUserId.trim()"
          type="button"
          class="px-5 py-2.5 bg-red-500 hover:bg-red-600 active:scale-95 disabled:opacity-50 text-white font-black text-xs rounded-xl transition-all cursor-pointer shadow-lg shadow-red-500/20"
        >
          {{ actionLoading === 'manual' ? 'جاري الحظر...' : 'حظر لمدة 5 دقائق' }}
        </button>
      </div>
    </div>

    <!-- Banned Users Table -->
    <div class="bg-white dark:bg-[#121212] border border-gray-200 dark:border-gray-800 rounded-2xl overflow-hidden shadow-sm">
      <div v-if="loading" class="p-12 text-center text-gray-500 dark:text-gray-400 text-xs font-bold">
        جاري تحميل سجل المحظورين...
      </div>

      <div v-else-if="bannedUsers.length === 0" class="p-12 text-center text-gray-500 dark:text-gray-400">
        <div class="text-4xl mb-3">✅</div>
        <div class="text-sm font-black text-gray-900 dark:text-white">لا يوجد مستخدمين محظورين حالياً</div>
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">جميع الطلاب يعملون بشكل طبيعي وبدون مخالفات حماية.</p>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-right text-xs">
          <thead class="bg-gray-50 dark:bg-[#1A1A1A] border-b border-gray-200 dark:border-gray-800 text-gray-600 dark:text-gray-400 font-black">
            <tr>
              <th class="px-5 py-3.5">مُعرّف الطالب (User ID)</th>
              <th class="px-5 py-3.5">نقاط المخالفات</th>
              <th class="px-5 py-3.5">نوع الانتهاك الأخير</th>
              <th class="px-5 py-3.5">الوقت المتبقي لرفع الحظر</th>
              <th class="px-5 py-3.5 text-center">الإجراءات</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100 dark:divide-gray-800/60 font-medium">
            <tr v-for="user in bannedUsers" :key="user.user_id" class="hover:bg-gray-50/50 dark:hover:bg-gray-800/30 transition-colors">
              <td class="px-5 py-4 font-mono font-bold text-gray-900 dark:text-white">
                {{ user.user_id }}
              </td>
              <td class="px-5 py-4 font-bold text-red-500">
                {{ user.points }} نقطة
              </td>
              <td class="px-5 py-4 text-gray-700 dark:text-gray-300">
                <span class="px-2.5 py-1 rounded-md bg-red-500/10 text-red-500 border border-red-500/20 font-bold font-mono">
                  {{ user.ban?.last_violation || 'محاولة تصوير / سكرين شوت' }}
                </span>
              </td>
              <td class="px-5 py-4 font-mono font-bold text-amber-500">
                ⏱️ {{ formatUntilTime(user.ban?.until_ts) }}
              </td>
              <td class="px-5 py-4 text-center">
                <button
                  @click="handleUnban(user.user_id)"
                  :disabled="actionLoading === user.user_id"
                  type="button"
                  class="px-4 py-2 bg-emerald-500 hover:bg-emerald-600 active:scale-95 text-white font-black text-xs rounded-xl transition-all cursor-pointer shadow-md shadow-emerald-500/20"
                >
                  {{ actionLoading === user.user_id ? 'جاري الفك...' : '🔓 رفع الحظر فوراً' }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
