<template>
  <div class="p-6 space-y-6 bg-neutral-950 text-neutral-100 min-h-screen" dir="rtl">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-neutral-800 pb-5">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-white flex items-center gap-2">
          <ShieldAlert class="w-6 h-6 text-amber-500" />
          لوحة الحماية والأمان
        </h1>
        <p class="text-xs text-neutral-400 mt-1">
          مراقبة الجلسات وتنبيهات الأمان وتقييم المخاطر المباشر لطلاب Weyana
        </p>
      </div>

      <button
        @click="loadAll"
        :disabled="loading"
        class="px-4 py-2 bg-neutral-800 hover:bg-neutral-700 text-neutral-200 rounded-xl text-xs font-semibold flex items-center gap-2 transition cursor-pointer"
      >
        <RotateCcw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
        تحديث البيانات
      </button>
    </div>

    <!-- Active Sessions summary -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="bg-neutral-900 border border-neutral-800 rounded-2xl p-4">
        <p class="text-[11px] text-neutral-400">الجلسات النشطة</p>
        <p class="text-2xl font-black text-emerald-400">{{ sessions.filter((s) => s.status === 'active').length }}</p>
      </div>
      <div class="bg-neutral-900 border border-neutral-800 rounded-2xl p-4">
        <p class="text-[11px] text-neutral-400">أحداث أمنية</p>
        <p class="text-2xl font-black text-amber-400">{{ events.length }}</p>
      </div>
      <div class="bg-neutral-900 border border-neutral-800 rounded-2xl p-4">
        <p class="text-[11px] text-neutral-400">جلسات عالية الخطورة</p>
        <p
          class="text-2xl font-black"
          :class="highRiskCount > 0 ? 'text-red-400' : 'text-neutral-600'"
        >
          {{ highRiskCount }}
        </p>
      </div>
    </div>

    <!-- Security Events Table -->
    <div class="bg-neutral-900 border border-neutral-800 rounded-2xl overflow-hidden shadow-xl">
      <div class="p-4 border-b border-neutral-800 flex items-center justify-between">
        <h2 class="text-sm font-bold text-neutral-200">سجل أحداث الأمان الأخيرة</h2>
        <span class="text-xs text-neutral-400 font-mono">{{ events.length }} حدث</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-right text-xs">
          <thead class="bg-neutral-950/80 text-neutral-400 border-b border-neutral-800 font-mono">
            <tr>
              <th class="p-3.5">المستخدم</th>
              <th class="p-3.5">المنصة</th>
              <th class="p-3.5">الجلسة</th>
              <th class="p-3.5">الحدث</th>
              <th class="p-3.5">الخطر</th>
              <th class="p-3.5">الثقة</th>
              <th class="p-3.5">حالة الجلسة</th>
              <th class="p-3.5">التاريخ والوقت</th>
              <th class="p-3.5 text-center">إجراءات</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-neutral-800/60 font-mono">
            <tr v-if="events.length === 0">
              <td colspan="9" class="p-8 text-center text-neutral-500 text-sm">
                لا توجد أحداث أمان مسجلة حالياً
              </td>
            </tr>
            <tr v-for="ev in events" :key="ev.eventId" class="hover:bg-neutral-800/40 transition">
              <td class="p-3.5 font-bold text-amber-400">{{ (ev.userId || '').slice(0, 12) }}</td>
              <td class="p-3.5 text-neutral-300">{{ platformLabel(ev.platform || ev.metadata?.platform) }}</td>
              <td class="p-3.5 text-neutral-400">{{ (ev.sessionId || '').slice(0, 8) }}...</td>
              <td class="p-3.5 font-semibold text-neutral-200">{{ ev.eventType }}</td>
              <td class="p-3.5 font-bold" :class="riskClass(ev.riskScore || 0)">
                {{ ev.riskScore || 0 }}
              </td>
              <td class="p-3.5 text-neutral-300">{{ confidenceLabel(ev.confidence) }}</td>
              <td class="p-3.5">
                <span
                  class="px-2 py-0.5 rounded-md text-[10px] font-bold"
                  :class="statusClass(sessionStatus(ev.sessionId))"
                >
                  {{ sessionStatus(ev.sessionId) }}
                </span>
              </td>
              <td class="p-3.5 text-neutral-400">
                {{ new Date((ev.timestamp || 0) * 1000).toLocaleString('ar-IQ') }}
              </td>
              <td class="p-3.5 text-center whitespace-nowrap">
                <button
                  @click="openUserHistory(ev.userId)"
                  class="px-2 py-1 bg-sky-500/10 hover:bg-sky-500/20 text-sky-400 border border-sky-500/30 rounded-lg text-[11px] font-semibold transition cursor-pointer ml-1"
                >
                  السجل
                </button>
                <button
                  @click="revoke(ev.sessionId)"
                  class="px-2 py-1 bg-red-500/10 hover:bg-red-500/20 text-red-400 border border-red-500/30 rounded-lg text-[11px] font-semibold transition cursor-pointer"
                >
                  إلغاء
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Active Sessions -->
    <div class="bg-neutral-900 border border-neutral-800 rounded-2xl overflow-hidden shadow-xl">
      <div class="p-4 border-b border-neutral-800 flex items-center justify-between">
        <h2 class="text-sm font-bold text-neutral-200">الجلسات النشطة</h2>
        <span class="text-xs text-neutral-400 font-mono">{{ sessions.length }} جلسة</span>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-right text-xs">
          <thead class="bg-neutral-950/80 text-neutral-400 border-b border-neutral-800 font-mono">
            <tr>
              <th class="p-3.5">الجلسة</th>
              <th class="p-3.5">المستخدم</th>
              <th class="p-3.5">الفيديو</th>
              <th class="p-3.5">المنصة</th>
              <th class="p-3.5">الخطر</th>
              <th class="p-3.5">الحالة</th>
              <th class="p-3.5 text-center">إجراء</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-neutral-800/60 font-mono">
            <tr v-if="sessions.length === 0">
              <td colspan="7" class="p-8 text-center text-neutral-500 text-sm">لا توجد جلسات</td>
            </tr>
            <tr v-for="s in sessions" :key="s.sessionId" class="hover:bg-neutral-800/40 transition">
              <td class="p-3.5 text-neutral-400">{{ (s.sessionId || '').slice(0, 8) }}...</td>
              <td class="p-3.5 font-bold text-amber-400">{{ (s.userId || '').slice(0, 12) }}</td>
              <td class="p-3.5 text-neutral-300">{{ (s.videoId || '').slice(0, 12) }}</td>
              <td class="p-3.5 text-neutral-300">{{ platformLabel(s.platform) }}</td>
              <td class="p-3.5 font-bold" :class="riskClass(s.riskScore || 0)">{{ s.riskScore || 0 }}</td>
              <td class="p-3.5">
                <span class="px-2 py-0.5 rounded-md text-[10px] font-bold" :class="statusClass(s.status)">
                  {{ s.status }}
                </span>
              </td>
              <td class="p-3.5 text-center">
                <button
                  @click="revokeAllUser(s.userId)"
                  class="px-2.5 py-1 bg-red-500/10 hover:bg-red-500/20 text-red-400 border border-red-500/30 rounded-lg text-[11px] font-semibold transition cursor-pointer"
                >
                  إلغاء جميع جلسات المستخدم
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- User Security History Modal -->
    <div
      v-if="historyOpen"
      class="fixed inset-0 z-[9998] flex items-center justify-center bg-black/85 backdrop-blur-sm p-4"
      dir="rtl"
      @click.self="historyOpen = false"
    >
      <div class="w-full max-w-lg max-h-[80vh] flex flex-col bg-neutral-900 border border-neutral-800 rounded-2xl shadow-2xl overflow-hidden">
        <div class="p-4 border-b border-neutral-800 flex items-center justify-between">
          <h3 class="text-sm font-bold text-neutral-100">سجل أمان المستخدم {{ historyUserId }}</h3>
          <button @click="historyOpen = false" class="text-neutral-400 hover:text-white text-lg cursor-pointer">×</button>
        </div>
        <div class="overflow-y-auto divide-y divide-neutral-800/60 font-mono">
          <div v-if="historyLoading" class="p-8 text-center text-neutral-500">جارٍ التحميل…</div>
          <div v-else-if="userHistory.length === 0" class="p-8 text-center text-neutral-500">لا يوجد سجل</div>
          <div v-for="(e, i) in userHistory" :key="i" class="p-3.5 text-xs flex items-center justify-between gap-3">
            <div class="min-w-0">
              <div class="font-semibold text-neutral-200">{{ e.eventType }}</div>
              <div class="text-neutral-500 mt-0.5">
                {{ new Date((e.timestamp || 0) * 1000).toLocaleString('ar-IQ') }} • خطر {{ e.riskScore || 0 }}
              </div>
            </div>
            <div class="text-neutral-400 shrink-0">{{ e.sessionId?.slice(0, 8) }}…</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ShieldAlert, RotateCcw } from '@lucide/vue'
import {
  fetchAdminSecurityEvents,
  fetchAdminSessions,
  revokeSessionAdmin,
  revokeAllUserSessionsAdmin,
  fetchUserSecurityHistory,
} from '@/api/securityApi'

const events = ref<any[]>([])
const sessions = ref<any[]>([])
const loading = ref(false)
const historyOpen = ref(false)
const historyLoading = ref(false)
const historyUserId = ref('')
const userHistory = ref<any[]>([])

const sessionStatusMap = computed(() => {
  const map: Record<string, string> = {}
  for (const s of sessions.value) map[s.sessionId] = s.status || 'unknown'
  return map
})

const highRiskCount = computed(
  () => events.value.filter((e) => (e.riskScore || 0) >= 60).length
)

const sessionStatus = (sessionId: string) => sessionStatusMap.value[sessionId] || '—'

const loadEvents = async () => {
  events.value = await fetchAdminSecurityEvents(50)
}

const loadSessions = async () => {
  try {
    sessions.value = await fetchAdminSessions()
  } catch {
    sessions.value = []
  }
}

const loadAll = async () => {
  loading.value = true
  try {
    await Promise.all([loadEvents(), loadSessions()])
  } catch (err) {
    console.error('Failed to load security data:', err)
  } finally {
    loading.value = false
  }
}

const revoke = async (sessionId: string) => {
  if (!sessionId) return
  if (confirm('هل أنت متأكد من إلغاء جلسة البث هذه فوراً؟')) {
    try {
      await revokeSessionAdmin(sessionId, 'admin_manual_revoke')
      await loadAll()
    } catch (err) {
      alert('تعذر إلغاء الجلسة')
    }
  }
}

const revokeAllUser = async (userId: string) => {
  if (!userId) return
  if (confirm('هل أنت متأكد من إلغاء جميع جلسات هذا المستخدم؟')) {
    try {
      await revokeAllUserSessionsAdmin(userId)
      await loadAll()
    } catch (err) {
      alert('تعذر إلغاء الجلسات')
    }
  }
}

const openUserHistory = async (userId: string) => {
  historyUserId.value = userId
  historyOpen.value = true
  historyLoading.value = true
  userHistory.value = []
  try {
    userHistory.value = await fetchUserSecurityHistory(userId, 50)
  } catch (err) {
    console.error('Failed to load user history:', err)
  } finally {
    historyLoading.value = false
  }
}

const platformLabel = (p?: string) => {
  const labels: Record<string, string> = {
    ANDROID_TELEGRAM_WEBVIEW: 'Android Telegram',
    IOS_TELEGRAM_WEBVIEW: 'iOS Telegram',
    TELEGRAM_MINI_APP: 'Telegram Mini App',
    ANDROID_WEB: 'Android Web',
    IOS_WEB: 'iOS Web',
    DESKTOP_WEB: 'Desktop',
    ANDROID_NATIVE: 'Android Native',
  }
  return labels[p || ''] || p || '—'
}

const confidenceLabel = (c?: number) => {
  const v = Number(c ?? 0)
  if (v >= 0.8) return 'HIGH'
  if (v >= 0.5) return 'MEDIUM'
  if (v > 0) return 'LOW'
  return 'NONE'
}

const riskClass = (score: number) =>
  score >= 80 ? 'text-red-300' : score >= 60 ? 'text-red-400' : score >= 30 ? 'text-amber-400' : 'text-emerald-400'

const statusClass = (status: string) => {
  switch ((status || '').toLowerCase()) {
    case 'active': return 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/30'
    case 'revoked': return 'bg-red-500/10 text-red-400 border border-red-500/30'
    case 'superseded': return 'bg-amber-500/10 text-amber-400 border border-amber-500/30'
    default: return 'bg-neutral-700/30 text-neutral-400 border border-neutral-700'
  }
}

onMounted(() => {
  loadAll()
})
</script>