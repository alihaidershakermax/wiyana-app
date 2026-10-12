<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Key, Plus, Copy, Download, Search, RefreshCw, Tag, Power, Trash2, X, CheckCircle, AlertCircle } from '@lucide/vue'
import { adminApiClient } from '@/admin/api/adminApi'
import { useCourses } from '@/admin/queries/useCourses'
import StatePanel from '@/admin/components/StatePanel.vue'
import ActionDropdown, { type MenuItem } from '@/admin/components/ui/ActionDropdown.vue'
import ToastContainer from '@/admin/components/ui/ToastContainer.vue'
import { useToast } from '@/admin/composables/useToast'
import { formatDate } from '@/utils/format'

interface SubscriptionCode {
  id?: string
  _id?: string
  code: string
  course_id?: string
  course_ids?: string[]
  max_uses: number
  used_count: number
  status?: string
  is_active: boolean
  duration_days?: number
  created_at: string
}

const toast = useToast()
const { data: rawCourses } = useCourses()
const coursesList = computed(() => {
  const c = rawCourses.value
  if (!c) return []
  if (Array.isArray(c)) return c
  if (Array.isArray(c.data)) return c.data
  return []
})

const codes = ref<SubscriptionCode[]>([])
const loading = ref(false)
const errorMsg = ref<string | null>(null)
const search = ref('')
const statusFilter = ref('')
const courseFilter = ref('')

const showGenerateModal = ref(false)
const generating = ref(false)

const form = ref({
  course_id: '',
  count: 10,
  duration_days: 30,
  prefix: 'WYN',
})

function getCourseTitle(cid?: string, course_ids?: string[]): string {
  if (course_ids && course_ids.length > 1) {
    return `متعدد (${course_ids.length} مساقات)`
  }
  const targetId = cid || (course_ids && course_ids[0])
  if (!targetId || targetId === 'ALL') return 'كافة المساقات (الاشتراك الشامل)'
  const course = coursesList.value.find((c: any) => c.id === targetId || c._id === targetId)
  return course?.title || 'مساق محدد'
}

async function fetchCodes() {
  loading.value = true
  errorMsg.value = null
  try {
    const { data } = await adminApiClient.get('/subscription-codes/', {
      params: {
        search: search.value.trim() || undefined,
        status: statusFilter.value || undefined,
        course_id: courseFilter.value || undefined,
      },
    })
    const res = data?.data ?? data
    codes.value = Array.isArray(res) ? res : (Array.isArray(res?.items) ? res.items : [])
  } catch (err: unknown) {
    errorMsg.value = (err as Error)?.message || 'فشل تحميل أكواد الاشتراك'
  } finally {
    loading.value = false
  }
}

async function handleGenerate() {
  generating.value = true
  try {
    const selectedCourseIds = form.value.course_id && form.value.course_id !== 'ALL'
      ? [form.value.course_id]
      : []

    const { data } = await adminApiClient.post('/subscription-codes/batch-generate', {
      course_ids: selectedCourseIds,
      count: Math.max(1, Number(form.value.count) || 1),
      duration_days: Math.max(1, Number(form.value.duration_days) || 30),
      prefix: form.value.prefix.trim().toUpperCase() || 'WYN',
    })

    const generatedCount = data?.count || form.value.count || 1
    toast.success('تم التوليد بنجاح', `تم توليد ${generatedCount} كود اشتراك بنجاح`)
    showGenerateModal.value = false
    await fetchCodes()
  } catch (err: any) {
    const detail = err?.response?.data?.detail || err?.message || 'حدث خطأ أثناء توليد الأكواد'
    toast.error('فشل توليد الأكواد', detail)
  } finally {
    generating.value = false
  }
}

function copyCode(code: string) {
  navigator.clipboard.writeText(code)
  toast.success('تم النسخ', 'تم نسخ كود الاشتراك إلى الحافظة')
}

function getActionItems(c: SubscriptionCode): MenuItem[] {
  return [
    {
      id: 'copy',
      label: 'نسخ كود الاشتراك',
      icon: Copy,
      onClick: () => copyCode(c.code),
    },
    {
      id: 'status',
      label: c.is_active ? 'تعطيل الكود' : 'تفعيل الكود',
      icon: Power,
      danger: c.is_active,
      onClick: async () => {
        try {
          const codeId = c.id || c._id
          await adminApiClient.patch(`/subscription-codes/${codeId}`, { is_active: !c.is_active })
          c.is_active = !c.is_active
          toast.success('تم التحديث', c.is_active ? 'تم تفعيل الكود بنجاح' : 'تم تعطيل الكود بنجاح')
        } catch {
          toast.error('خطأ', 'تعذر تحديث حالة الكود')
        }
      },
    },
    {
      id: 'delete',
      label: 'حذف الكود',
      icon: Trash2,
      danger: true,
      onClick: async () => {
        try {
          const codeId = c.id || c._id
          await adminApiClient.delete(`/subscription-codes/${codeId}`)
          codes.value = codes.value.filter(i => (i.id || i._id) !== codeId)
          toast.success('تم الحذف', 'تم حذف كود الاشتراك بنجاح')
        } catch {
          toast.error('خطأ', 'تعذر حذف الكود')
        }
      },
    },
  ]
}

function exportCodes() {
  if (!codes.value.length) return
  const text = codes.value.map((c) => c.code).join('\n')
  const blob = new Blob([text], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `subscription_codes_${new Date().toISOString().slice(0, 10)}.txt`
  a.click()
  URL.revokeObjectURL(url)
  toast.success('تم التصدير', 'تم تصدير ملف الأكواد بنجاح')
}

onMounted(() => {
  fetchCodes()
})
</script>

<template>
  <div class="space-y-6 text-[#F5F5F5]">
    <ToastContainer />

    <div class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="text-xl font-extrabold flex items-center gap-2">
          <Key class="h-5 w-5 text-[#F97316]" />
          أكواد الاشتراكات
        </h1>
        <p class="text-xs text-[#A1A1AA] mt-1">
          توليد وتصدير وإدارة أكواد تفعيل المساقات للطلاب
        </p>
      </div>

      <div class="flex items-center gap-2">
        <button class="btn-secondary text-xs" :disabled="!codes.length" @click="exportCodes">
          <Download class="h-4 w-4" />
          تصدير النص
        </button>
        <button class="btn-primary text-xs cursor-pointer" @click="showGenerateModal = true">
          <Plus class="h-4 w-4" />
          توليد أكواد جديدة
        </button>
      </div>
    </div>

    <!-- Filters -->
    <div class="card p-4">
      <div class="flex flex-wrap items-center gap-3">
        <div class="relative flex-1 min-w-[200px]">
          <Search class="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-[#71717A]" />
          <input
            v-model="search"
            type="text"
            placeholder="ابحث بالكود..."
            class="input w-full pr-9 text-sm"
            @keyup.enter="fetchCodes"
          />
        </div>

        <select v-model="courseFilter" class="input text-sm w-44 font-semibold" @change="fetchCodes">
          <option value="">كل المساقات</option>
          <option v-for="c in coursesList" :key="c.id" :value="c.id">{{ c.title }}</option>
        </select>

        <select v-model="statusFilter" class="input text-sm w-36 font-semibold" @change="fetchCodes">
          <option value="">كل الحالات</option>
          <option value="AVAILABLE">متاح</option>
          <option value="EXPIRED">معطل / مستخدم</option>
        </select>

        <button class="btn-secondary text-xs cursor-pointer" @click="fetchCodes">
          <RefreshCw class="h-4 w-4" />
          تحديث
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="card overflow-hidden">
      <StatePanel
        :loading="loading"
        :error="errorMsg"
        :empty="!codes.length"
        empty-title="لا توجد أكواد اشتراك"
        empty-message="قم بتوليد دفعة جديدة من الأكواد للمساقات"
        :skeleton-count="5"
        @retry="fetchCodes"
      >
        <div class="overflow-x-auto">
          <table class="w-full text-right border-collapse text-sm text-[#F5F5F5]">
            <thead>
              <tr class="bg-[#0F0F0F] text-[#A1A1AA] border-b border-[#242424] text-xs font-extrabold">
                <th class="py-3.5 px-4">الكود</th>
                <th class="py-3.5 px-4">المساق</th>
                <th class="py-3.5 px-4">الاستخدام</th>
                <th class="py-3.5 px-4">المدة</th>
                <th class="py-3.5 px-4">الحالة</th>
                <th class="py-3.5 px-4">تاريخ الإنشاء</th>
                <th class="py-3.5 px-4 text-center w-16">إجراءات</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-[#242424] bg-[#111111]">
              <tr v-for="c in codes" :key="c.id || c._id" class="hover:bg-[#161616] transition-colors">
                <td class="py-3.5 px-4 font-mono font-bold text-[#F97316]">
                  {{ c.code }}
                </td>
                <td class="py-3.5 px-4 text-xs font-bold text-slate-300">
                  {{ getCourseTitle(c.course_id, c.course_ids) }}
                </td>
                <td class="py-3.5 px-4 text-xs text-[#A1A1AA]">
                  {{ c.used_count || 0 }} / {{ c.max_uses || 1 }}
                </td>
                <td class="py-3.5 px-4 text-xs text-[#71717A]">
                  {{ c.duration_days ?? 30 }} يوم
                </td>
                <td class="py-3.5 px-4">
                  <span
                    :class="[
                      'text-[10px] font-extrabold px-2.5 py-1 rounded-full border inline-block',
                      (c.used_count && c.used_count >= c.max_uses)
                        ? 'bg-[#18181b] text-[#A1A1AA] border-[#27272a]'
                        : c.is_active !== false
                        ? 'bg-[#052e16] text-[#22C55E] border-[#14532d]'
                        : 'bg-[#450a0a] text-[#EF4444] border-[#7f1d1d]'
                    ]"
                  >
                    {{ (c.used_count && c.used_count >= c.max_uses) ? 'مستعمل' : (c.is_active !== false ? 'متاح' : 'معطل') }}
                  </span>
                </td>
                <td class="py-3.5 px-4 text-xs text-[#71717A]">
                  {{ formatDate(c.created_at) }}
                </td>
                <td class="py-3.5 px-4 text-center">
                  <ActionDropdown :items="getActionItems(c)" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </StatePanel>
    </div>

    <!-- Generate Modal -->
    <div
      v-if="showGenerateModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm animate-fade-in"
      @click.self="showGenerateModal = false"
    >
      <div class="card max-w-md w-full p-6 space-y-4 border-[#242424] bg-[#111111] shadow-2xl relative">
        <div class="flex items-center justify-between border-b border-[#242424] pb-3">
          <h3 class="text-base font-extrabold flex items-center gap-2 text-white">
            <Tag class="h-5 w-5 text-[#F97316]" />
            توليد دفعة أكواد اشتراك
          </h3>
          <button
            type="button"
            class="p-1 rounded-lg text-[#71717A] hover:text-white hover:bg-[#1C1C1C] transition-colors cursor-pointer"
            @click="showGenerateModal = false"
          >
            <X class="h-5 w-5" />
          </button>
        </div>

        <form class="space-y-4" @submit.prevent="handleGenerate">
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-[#A1A1AA]">تحديد المساق</label>
            <select v-model="form.course_id" class="input w-full text-sm font-semibold">
              <option value="ALL">كافة المساقات (الاشتراك الشامل)</option>
              <option v-for="c in coursesList" :key="c.id || c._id" :value="c.id || c._id">
                {{ c.title }}
              </option>
            </select>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1.5">
              <label class="text-xs font-bold text-[#A1A1AA]">عدد الأكواد المراد توليدها</label>
              <input v-model.number="form.count" type="number" min="1" max="500" class="input w-full text-sm font-bold" required />
            </div>
            <div class="space-y-1.5">
              <label class="text-xs font-bold text-[#A1A1AA]">مدة الصلاحية (بالأيام)</label>
              <input v-model.number="form.duration_days" type="number" min="1" max="365" class="input w-full text-sm font-bold" required />
            </div>
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-[#A1A1AA]">بادئة الكود (Prefix)</label>
            <input v-model="form.prefix" type="text" placeholder="مثال: WYN أو MATH" class="input w-full text-sm uppercase font-bold" />
          </div>

          <div class="flex items-center justify-end gap-2 pt-2 border-t border-[#242424]">
            <button type="button" class="btn-secondary text-xs cursor-pointer" @click="showGenerateModal = false">
              إلغاء
            </button>
            <button type="submit" class="btn-primary text-xs cursor-pointer" :disabled="generating">
              {{ generating ? 'جارٍ التوليد...' : 'توليد الأكواد' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
