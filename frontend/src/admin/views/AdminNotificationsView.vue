<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Bell, Send, RefreshCw, Trash2, ShieldCheck, Megaphone } from '@lucide/vue'
import { adminApiClient } from '@/admin/api/adminApi'
import { useCourses } from '@/admin/queries/useCourses'
import { useAdminStore } from '@/admin/stores/useAdminStore'
import StatePanel from '@/admin/components/StatePanel.vue'
import { useToast } from '@/admin/composables/useToast'
import { formatDate } from '@/utils/format'

interface NotificationItem {
  id?: string
  _id?: string
  title: string
  message: string
  type: string
  created_at: string
}

const toast = useToast()
const adminStore = useAdminStore()
const isOwner = computed(() => adminStore.userRole === 'OWNER')

const { data: rawCourses } = useCourses()
const coursesList = computed(() => {
  const c = rawCourses.value
  if (!c) return []
  if (Array.isArray(c)) return c
  if (Array.isArray(c.data)) return c.data
  return []
})

const notifications = ref<NotificationItem[]>([])
const loading = ref(false)
const errorMsg = ref<string | null>(null)

const showBroadcastModal = ref(false)
const broadcasting = ref(false)

const form = ref({
  title: '',
  message: '',
  target_type: 'ALL_STUDENTS',
  target_id: '',
})

async function fetchNotifications() {
  loading.value = true
  errorMsg.value = null
  try {
    const { data } = await adminApiClient.get('/notifications/')
    notifications.value = data.data ?? data
  } catch (err: unknown) {
    errorMsg.value = (err as Error)?.message || 'فشل تحميل الإشعارات'
  } finally {
    loading.value = false
  }
}

async function handleBroadcast() {
  if (!isOwner.value) {
    toast.error('صلاحية غير متاحة', 'هذه الميزة مخصصة لمالك المنصة فقط')
    return
  }

  if (!form.value.title.trim() || !form.value.message.trim()) {
    toast.error('تنبيه', 'يرجى كتابة عنوان ورسالة الإشعار')
    return
  }
  broadcasting.value = true
  try {
    const { data } = await adminApiClient.post('/notifications/broadcast', {
      title: form.value.title.trim(),
      message: form.value.message.trim(),
      target_type: form.value.target_type,
      target_id: form.value.target_id || undefined,
    })
    toast.success(data.message || 'تم إرسال التنويه بنجاح')
    showBroadcastModal.value = false
    form.value.title = ''
    form.value.message = ''
    fetchNotifications()
  } catch (err: unknown) {
    toast.error('فشل إرسال التنويه', (err as Error)?.message)
  } finally {
    broadcasting.value = false
  }
}

async function deleteNotification(id: string) {
  if (!isOwner.value) return
  try {
    await adminApiClient.delete(`/notifications/${id}`)
    toast.success('تم حذف الإشعار')
    fetchNotifications()
  } catch (err: unknown) {
    toast.error('فشل الحذف', (err as Error)?.message)
  }
}

const showAppReleaseModal = ref(false)
const publishingRelease = ref(false)

const releaseForm = ref({
  version: '1.2.0',
  title: 'تحديث التطبيق',
  featuresText: 'تحسينات على مشغل الفيديو\nتحديثات على حماية حسابات الطلاب\nتحسينات في واجهة Telegram Mini App',
})

async function handlePublishAppRelease() {
  if (!releaseForm.value.version.trim() || !releaseForm.value.title.trim()) {
    toast.error('تنبيه', 'يرجى إدخال رقم الإصدار وعنوان التحديث')
    return
  }

  publishingRelease.value = true
  try {
    const featuresList = releaseForm.value.featuresText
      .split('\n')
      .map((s) => s.trim())
      .filter(Boolean)

    const payload = {
      version: releaseForm.value.version.trim(),
      title: releaseForm.value.title.trim(),
      features: featuresList,
      publishedAt: new Date().toISOString(),
    }

    localStorage.setItem('wiyana_latest_app_release', JSON.stringify(payload))
    window.dispatchEvent(new CustomEvent('wiyana:app-update-published', { detail: payload }))

    try {
      await adminApiClient.post('/notifications/app-release', payload)
    } catch {}

    toast.success('تم نشر تحديث التطبيق')
    showAppReleaseModal.value = false
  } catch (err: unknown) {
    toast.error('فشل نشر التحديث', (err as Error)?.message)
  } finally {
    publishingRelease.value = false
  }
}

onMounted(() => {
  fetchNotifications()
})
</script>

<template>
  <div class="space-y-6 text-[#F5F5F5]">
    <!-- Header -->
    <div class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <div class="flex items-center gap-3">
          <h1 class="text-xl font-extrabold flex items-center gap-2">
            <Bell class="h-5 w-5 text-[#F97316]" />
            نظام الإشعارات وتحديثات التطبيق
          </h1>
          <span class="px-2.5 py-1 text-[11px] font-black text-[#F97316] bg-[#1F0E03] border border-orange-900/50 rounded-full flex items-center gap-1">
            <ShieldCheck class="h-3.5 w-3.5" />
            خاص بمالك المنصة
          </span>
        </div>
        <p class="text-xs text-[#A1A1AA] mt-1">
          إدارة التنويهات العامة وتحديثات التطبيق التي تظهر للطلاب.
        </p>
      </div>

      <div class="flex items-center gap-2">
        <button
          v-if="isOwner"
          class="px-4 py-2 rounded-xl bg-amber-500 hover:bg-amber-600 active:scale-95 text-slate-950 font-extrabold text-xs transition-all flex items-center gap-2 shadow-md shadow-amber-500/20"
          @click="showAppReleaseModal = true"
        >
          <Megaphone class="h-4 w-4" />
          نشر تحديث التطبيق
        </button>
        <button
          v-if="isOwner"
          class="btn-primary text-xs font-bold"
          @click="showBroadcastModal = true"
        >
          <Send class="h-4 w-4" />
          إرسال تنويه جديد
        </button>
      </div>
    </div>

    <!-- Notification List Container -->
    <div class="card overflow-hidden">
      <StatePanel
        :loading="loading"
        :error="errorMsg"
        :empty="!notifications.length"
        empty-title="لا توجد إشعارات مسجلة"
        empty-message="يمكنك بث تنويه جديد ليصل فوراً إلى المستهدفين"
        :skeleton-count="4"
        @retry="fetchNotifications"
      >
        <div class="divide-y divide-[#242424]">
          <div
            v-for="item in notifications"
            :key="item.id || item._id"
            class="p-4 flex items-start justify-between gap-4 hover:bg-[#161616] transition-colors"
          >
            <div class="space-y-1">
              <div class="flex items-center gap-2">
                <span class="text-[10px] font-extrabold text-[#F97316] bg-[#1F0E03] border border-orange-900/40 px-2 py-0.5 rounded-md flex items-center gap-1">
                  <Megaphone class="h-3 w-3" />
                  {{ item.type || 'تنويه مالك المنصة' }}
                </span>
                <h3 class="text-sm font-extrabold text-white">{{ item.title }}</h3>
              </div>
              <p class="text-xs text-[#A1A1AA] leading-relaxed">
                {{ item.message }}
              </p>
              <span class="text-[10px] text-[#71717A] block pt-1">
                تاريخ البث: {{ formatDate(item.created_at) }}
              </span>
            </div>

            <button
              v-if="isOwner"
              class="btn-ghost text-xs text-[#71717A] hover:text-red-400 p-1.5 rounded-lg"
              title="حذف الإشعار"
              @click="deleteNotification(String(item.id || item._id))"
            >
              <Trash2 class="h-4 w-4" />
            </button>
          </div>
        </div>
      </StatePanel>
    </div>

    <!-- Broadcast Modal (Owner Only) -->
    <div v-if="showBroadcastModal && isOwner" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm animate-fade-in">
      <div class="card max-w-lg w-full p-6 space-y-4 border-[#242424] bg-[#111111]">
        <div class="flex items-center justify-between border-b border-[#242424] pb-3">
          <h3 class="text-base font-extrabold text-white flex items-center gap-2">
            <Send class="h-5 w-5 text-[#F97316]" />
            بث تنويه جديد (مالك المنصة)
          </h3>
          <span class="text-[10px] font-bold text-[#F97316] bg-[#1F0E03] px-2 py-0.5 rounded-md">OWNER</span>
        </div>

        <form class="space-y-4" @submit.prevent="handleBroadcast">
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-[#A1A1AA]">عنوان الإشعار *</label>
            <input
              v-model="form.title"
              type="text"
              placeholder="مثال: تم إتاحة المحاضرة الجديدة في الفيزياء"
              class="input w-full text-sm"
              required
            />
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-[#A1A1AA]">نص الرسالة *</label>
            <textarea
              v-model="form.message"
              rows="4"
              placeholder="تفاصيل التنويه والتعليمات الموجهة للطلاب أو الأساتذة..."
              class="input w-full text-sm resize-none"
              required
            />
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-[#A1A1AA]">الجمهور المستهدف</label>
            <select v-model="form.target_type" class="input w-full text-sm font-semibold">
              <option value="ALL_STUDENTS">جميع الطلاب في المنصة</option>
              <option value="TEACHERS">جميع الأساتذة</option>
              <option value="COURSE">طلاب مساق دراسي معين</option>
            </select>
          </div>

          <div v-if="form.target_type === 'COURSE'" class="space-y-1.5">
            <label class="text-xs font-bold text-[#A1A1AA]">اختر المساق الدراسي</label>
            <select v-model="form.target_id" class="input w-full text-sm font-semibold" required>
              <option value="">-- حدد المساق --</option>
              <option v-for="c in coursesList" :key="c.id" :value="c.id">{{ c.title }}</option>
            </select>
          </div>

          <div class="flex items-center justify-end gap-2 pt-2 border-t border-[#242424]">
            <button type="button" class="btn-secondary text-xs" @click="showBroadcastModal = false">
              إلغاء
            </button>
            <button type="submit" class="btn-primary text-xs" :disabled="broadcasting">
              {{ broadcasting ? 'جارٍ البث...' : 'إرسال التنويه الآن' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- App Release Features Announcement Modal (Owner Only) -->
    <div v-if="showAppReleaseModal && isOwner" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm animate-fade-in text-right">
      <div class="card max-w-lg w-full p-6 space-y-4 border-amber-500/30 bg-[#111111] shadow-2xl">
        <div class="flex items-center justify-between border-b border-[#242424] pb-3">
          <h3 class="text-base font-extrabold text-white flex items-center gap-2">
            <Megaphone class="h-5 w-5 text-amber-500" />
            نشر تحديث تطبيق جديد كـ نافذة منبثقة للجميع
          </h3>
          <span class="text-[10px] font-bold text-amber-400 bg-amber-500/10 border border-amber-500/30 px-2 py-0.5 rounded-md">RELEASE POPUP</span>
        </div>

        <form class="space-y-4" @submit.prevent="handlePublishAppRelease">
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-[#A1A1AA]">رقم الإصدار (e.g. 1.2.0) *</label>
            <input
              v-model="releaseForm.version"
              type="text"
              placeholder="1.2.0"
              class="input w-full text-sm font-mono"
              required
            />
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-[#A1A1AA]">عنوان النافذة المنبثقة *</label>
            <input
              v-model="releaseForm.title"
              type="text"
              placeholder="تحديث التطبيق"
              class="input w-full text-sm font-bold"
              required
            />
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-[#A1A1AA]">قائمة التغييرات (كل تغيير في سطر) *</label>
            <textarea
              v-model="releaseForm.featuresText"
              rows="4"
              placeholder="تحسينات على مشغل الفيديو&#10;تحديثات على حماية الحسابات&#10;تحسينات في واجهة Telegram Mini App"
              class="input w-full text-sm resize-none"
              required
            />
          </div>

          <div class="flex items-center justify-end gap-2 pt-2 border-t border-[#242424]">
            <button type="button" class="btn-secondary text-xs" @click="showAppReleaseModal = false">
              إلغاء
            </button>
            <button type="submit" class="px-5 py-2 rounded-xl bg-amber-500 hover:bg-amber-600 active:scale-95 text-slate-950 font-extrabold text-xs transition-all" :disabled="publishingRelease">
              {{ publishingRelease ? 'جارٍ النشر...' : 'نشر التحديث كنافذة منبثقة' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
