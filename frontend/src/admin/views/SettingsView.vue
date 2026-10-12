<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
  Settings,
  Save,
  Lock,
  Video,
  Send,
  Database,
  Shield,
  CheckCircle2,
  AlertCircle,
  RefreshCw,
  ExternalLink,
  Unlink,
  KeyRound,
  Check,
  Radio,
  Tv,
} from '@lucide/vue'
import { adminApiClient } from '@/admin/api/adminApi'
import { youtubeApi } from '@/admin/api/panelApi'
import StatePanel from '@/admin/components/StatePanel.vue'
import { useToast } from '@/admin/composables/useToast'

const toast = useToast()
const loading = ref(false)
const saving = ref(false)
const errorMsg = ref<string | null>(null)

const activeTab = ref<'general' | 'youtube' | 'telegram' | 'storage'>('general')

const settings = ref({
  app_name: 'منصة ويانا التعليمية',
  support_contact: '@WiyanaSupport',
  allow_registration: true,
  youtube_privacy_default: 'unlisted',
  telegram_bot_active: true,
  storage_provider: 'local',
})

// YouTube status & OAuth
const ytLoading = ref(false)
const ytStatus = ref<{
  configured: boolean
  is_authenticated: boolean
  token_file_exists: boolean
  token_in_db: boolean
  simulation_mode: boolean
  channel_id?: string
  channel_name?: string
  channel_avatar?: string
  subscribers?: string
  videos_count?: string
} | null>(null)

const showManualCode = ref(false)
const manualCode = ref('')
const exchangingCode = ref(false)
const disconnectingYt = ref(false)

async function fetchSettings() {
  loading.value = true
  errorMsg.value = null
  try {
    const { data } = await adminApiClient.get('/settings/')
    if (data.data && Object.keys(data.data).length > 0) {
      settings.value = { ...settings.value, ...data.data }
    }
  } catch (err: unknown) {
    errorMsg.value = (err as Error)?.message || 'فشل تحميل الإعدادات'
  } finally {
    loading.value = false
  }
}

async function fetchYtStatus() {
  ytLoading.value = true
  try {
    const res = await youtubeApi.getStatus()
    ytStatus.value = res
  } catch (err) {
    console.warn('Failed to fetch YouTube status', err)
  } finally {
    ytLoading.value = false
  }
}

async function saveSettings() {
  saving.value = true
  try {
    await adminApiClient.post('/settings/', settings.value)
    toast.success('تم حفظ إعدادات النظام بنجاح')
  } catch (err: unknown) {
    toast.error('فشل حفظ الإعدادات', (err as Error)?.message)
  } finally {
    saving.value = false
  }
}

function startGoogleOAuth() {
  const backendBase = import.meta.env.VITE_ADMIN_API_BASE_URL || '/api/v1'
  window.open(`${backendBase}/youtube/oauth/authorize`, '_blank')
}

async function handleExchangeCode() {
  if (!manualCode.value.trim()) {
    toast.error('يرجى إدخال كود التفويض أو الرابط')
    return
  }
  exchangingCode.value = true
  try {
    await youtubeApi.exchangeCode(manualCode.value.trim())
    toast.success('تم ربط قناة YouTube بنجاح!')
    manualCode.value = ''
    showManualCode.value = false
    await fetchYtStatus()
  } catch (err: any) {
    toast.error('فشل استبدال كود التفويض', err?.response?.data?.detail || err?.message)
  } finally {
    exchangingCode.value = false
  }
}

async function handleDisconnectYt() {
  if (!confirm('هل أنت متأكد من رغبتك في إلغاء ربط قناة YouTube؟')) return
  disconnectingYt.value = true
  try {
    await youtubeApi.disconnect()
    toast.success('تم إلغاء ربط YouTube بنجاح')
    await fetchYtStatus()
  } catch (err: any) {
    toast.error('فشل إلغاء الربط', err?.message)
  } finally {
    disconnectingYt.value = false
  }
}

onMounted(() => {
  fetchSettings()
  fetchYtStatus()
})
</script>

<template>
  <div class="max-w-4xl mx-auto space-y-6">
    <div class="flex items-center justify-between gap-3">
      <div>
        <h1 class="text-xl font-extrabold flex items-center gap-2">
          <Settings class="h-5 w-5 text-orange-500" />
          إعدادات النظام
        </h1>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
          إدارة إعدادات المنصة، ربط YouTube، وتهيئات Telegram والتخزين
        </p>
      </div>

      <button class="btn-primary text-sm" :disabled="saving" @click="saveSettings">
        <Save class="h-4 w-4" />
        {{ saving ? 'جارٍ الحفظ...' : 'حفظ الإعدادات' }}
      </button>
    </div>

    <!-- Tabs -->
    <div class="flex border-b border-slate-200 dark:border-slate-800 gap-4 text-xs font-extrabold">
      <button
        v-for="t in [
          { key: 'general', label: 'عامة' },
          { key: 'youtube', label: 'YouTube' },
          { key: 'telegram', label: 'Telegram' },
          { key: 'storage', label: 'التخزين' },
        ]"
        :key="t.key"
        class="pb-3 border-b-2 transition-colors"
        :class="activeTab === t.key ? 'border-orange-500 text-orange-600 dark:text-orange-400' : 'border-transparent text-slate-400 hover:text-slate-600'"
        @click="activeTab = t.key as typeof activeTab"
      >
        {{ t.label }}
      </button>
    </div>

    <StatePanel
      :loading="loading"
      :error="errorMsg"
      :empty="false"
      :skeleton-count="4"
      @retry="() => fetchSettings()"
    >
      <div class="card p-6 space-y-6">
        <!-- General -->
        <template v-if="activeTab === 'general'">
          <div class="space-y-4">
            <div class="space-y-1.5">
              <label class="text-xs font-bold text-slate-600 dark:text-slate-300">اسم المنصة</label>
              <input v-model="settings.app_name" type="text" class="input w-full text-sm" />
            </div>
            <div class="space-y-1.5">
              <label class="text-xs font-bold text-slate-600 dark:text-slate-300">معرف الدعم الفني Telegram</label>
              <input v-model="settings.support_contact" type="text" dir="ltr" class="input w-full text-sm" />
            </div>
            <div class="pt-2">
              <label class="flex items-center gap-3 cursor-pointer">
                <input v-model="settings.allow_registration" type="checkbox" class="h-4 w-4 rounded accent-orange-500" />
                <span class="text-sm font-bold text-slate-700 dark:text-slate-200">السماح بتسجيل الطلاب الجدد تلقائياً</span>
              </label>
            </div>
          </div>
        </template>

        <!-- YouTube -->
        <template v-else-if="activeTab === 'youtube'">
          <div class="space-y-6">
            <!-- Connection Status Card -->
            <div
              class="rounded-2xl p-5 border transition-all"
              :class="ytStatus?.is_authenticated
                ? 'bg-emerald-950/20 border-emerald-500/30'
                : 'bg-slate-900/60 border-slate-800'"
            >
              <div class="flex items-start justify-between gap-4">
                <div class="flex items-center gap-3">
                  <div
                    class="h-12 w-12 rounded-xl flex items-center justify-center shrink-0"
                    :class="ytStatus?.is_authenticated ? 'bg-emerald-500/20 text-emerald-400' : 'bg-red-500/20 text-red-400'"
                  >
                    <Tv class="h-6 w-6" />
                  </div>
                  <div>
                    <div class="flex items-center gap-2">
                      <h3 class="font-extrabold text-base text-white">حالة ربط حساب YouTube</h3>
                      <span
                        class="px-2.5 py-0.5 text-xs font-bold rounded-full"
                        :class="ytStatus?.is_authenticated
                          ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'
                          : 'bg-amber-500/20 text-amber-400 border border-amber-500/30'"
                      >
                        {{ ytStatus?.is_authenticated ? 'متصل بنجاح' : 'غير متصل (وضع المحاكاة)' }}
                      </span>
                    </div>
                    <p class="text-xs text-slate-400 mt-1">
                      {{ ytStatus?.is_authenticated
                        ? 'تم توثيق الاتصال بنجاح. يتم رفع الفيديوهات وتوليد الروابط المشفرة تلقائياً.'
                        : 'لم يتم ربط القناة بعد. قم بتسجيل الدخول بحساب Google لربط القناة ورفع الفيديوهات.' }}
                    </p>
                  </div>
                </div>

                <button
                  type="button"
                  class="btn-secondary text-xs px-3 py-1.5 flex items-center gap-1.5 shrink-0"
                  :disabled="ytLoading"
                  @click="fetchYtStatus"
                >
                  <RefreshCw class="h-3.5 w-3.5" :class="ytLoading ? 'animate-spin' : ''" />
                  <span>تحديث</span>
                </button>
              </div>

              <!-- Channel Details if Authenticated -->
              <div
                v-if="ytStatus?.is_authenticated"
                class="mt-4 pt-4 border-t border-slate-800/80 grid grid-cols-1 sm:grid-cols-3 gap-3 text-xs"
              >
                <div class="p-3 rounded-xl bg-slate-950/60 border border-slate-800">
                  <span class="text-slate-500 block mb-1">اسم القناة</span>
                  <span class="font-bold text-white flex items-center gap-1.5">
                    <img
                      v-if="ytStatus.channel_avatar"
                      :src="ytStatus.channel_avatar"
                      class="w-5 h-5 rounded-full"
                    />
                    {{ ytStatus.channel_name || 'القناة المرتبطة' }}
                  </span>
                </div>
                <div class="p-3 rounded-xl bg-slate-950/60 border border-slate-800">
                  <span class="text-slate-500 block mb-1">معرف القناة (Channel ID)</span>
                  <span class="font-mono text-slate-300 select-all">{{ ytStatus.channel_id || '—' }}</span>
                </div>
                <div class="p-3 rounded-xl bg-slate-950/60 border border-slate-800">
                  <span class="text-slate-500 block mb-1">مكان حفظ التوكن</span>
                  <span class="text-emerald-400 font-bold">
                    {{ ytStatus.token_in_db ? 'سحابياً (Supabase)' : 'ملف محلي' }}
                  </span>
                </div>
              </div>

              <!-- Actions -->
              <div class="mt-5 flex flex-wrap items-center gap-3">
                <button
                  type="button"
                  class="btn-primary text-xs px-4 py-2 flex items-center gap-2 bg-gradient-to-r from-red-600 to-rose-600 hover:from-red-500 hover:to-rose-500 border-none shadow-lg shadow-red-500/20"
                  @click="startGoogleOAuth"
                >
                  <ExternalLink class="h-4 w-4" />
                  <span>{{ ytStatus?.is_authenticated ? 'إعادة ربط القناة عبر Google' : 'تسجيل الدخول وربط YouTube' }}</span>
                </button>

                <button
                  type="button"
                  class="btn-secondary text-xs px-3.5 py-2 flex items-center gap-1.5"
                  @click="showManualCode = !showManualCode"
                >
                  <KeyRound class="h-3.5 w-3.5" />
                  <span>الربط اليدوي بكود التفويض</span>
                </button>

                <button
                  v-if="ytStatus?.is_authenticated"
                  type="button"
                  class="text-xs px-3 py-2 text-red-400 hover:text-red-300 hover:bg-red-950/40 rounded-xl transition-colors flex items-center gap-1.5 mr-auto"
                  :disabled="disconnectingYt"
                  @click="handleDisconnectYt"
                >
                  <Unlink class="h-3.5 w-3.5" />
                  <span>{{ disconnectingYt ? 'جارٍ الإلغاء...' : 'إلغاء الربط' }}</span>
                </button>
              </div>

              <!-- Manual Code Exchange Drawer -->
              <div
                v-if="showManualCode"
                class="mt-4 p-4 rounded-xl bg-slate-950 border border-slate-800 space-y-3"
              >
                <label class="block text-xs font-bold text-slate-300">
                  الصق كود التفويض (Authorization Code) أو رابط إعادة التوجيه الكامل بعد الموافقة:
                </label>
                <div class="flex gap-2">
                  <input
                    v-model="manualCode"
                    type="text"
                    dir="ltr"
                    placeholder="https://.../?code=4/0A... أو 4/0A..."
                    class="input flex-1 text-xs font-mono"
                    :disabled="exchangingCode"
                  />
                  <button
                    type="button"
                    class="btn-primary text-xs px-4"
                    :disabled="exchangingCode || !manualCode.trim()"
                    @click="handleExchangeCode"
                  >
                    {{ exchangingCode ? 'جارٍ التحقق...' : 'تأكيد وحفظ' }}
                  </button>
                </div>
              </div>
            </div>

            <!-- Upload Settings -->
            <div class="space-y-1.5">
              <label class="text-xs font-bold text-slate-600 dark:text-slate-300">الخصوصية الافتراضية للفيديوهات المرفوعة</label>
              <select v-model="settings.youtube_privacy_default" class="input w-full text-sm">
                <option value="unlisted">غير مدرج (Unlisted) — موصى به</option>
                <option value="private">خاص (Private)</option>
                <option value="public">عام (Public)</option>
              </select>
              <p class="text-[11px] text-slate-500">
                وضع "غير مدرج" يسمح بتشغيل الفيديو داخل منصة ويانا فقط للطلاب المشتركين دون ظهوره في بحث يوتيوب العام.
              </p>
            </div>
          </div>
        </template>

        <!-- Telegram -->
        <template v-else-if="activeTab === 'telegram'">
          <div class="space-y-4">
            <label class="flex items-center justify-between cursor-pointer p-4 rounded-xl bg-slate-900/60 border border-slate-800">
              <div>
                <span class="text-sm font-bold text-white block">تفعيل بوت Telegram التلقائي</span>
                <span class="text-xs text-slate-400">إرسال إشعارات الاشتراك والدروس والتسجيل عبر البوت</span>
              </div>
              <input v-model="settings.telegram_bot_active" type="checkbox" class="h-5 w-5 rounded accent-orange-500" />
            </label>
          </div>
        </template>

        <!-- Storage -->
        <template v-else-if="activeTab === 'storage'">
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-600 dark:text-slate-300">نظام تخزين المرفقات والملفات</label>
            <select v-model="settings.storage_provider" class="input w-full text-sm">
              <option value="local">تخزين محلي (Local Disk Storage)</option>
              <option value="s3">تخزين سحابي S3 Compatible</option>
            </select>
          </div>
        </template>
      </div>
    </StatePanel>
  </div>
</template>
