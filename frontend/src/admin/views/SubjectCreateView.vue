<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useCreateSubject } from '@/admin/queries/useSubjects'
import { useToast } from '@/admin/composables/useToast'
import ToastContainer from '@/admin/components/ui/ToastContainer.vue'
import { ArrowRight, FolderOpen, Loader2 } from '@lucide/vue'
import type { SubjectCreate } from '@/admin/types/subject'

const router = useRouter()
const toast = useToast()
const createSubject = useCreateSubject()

const form = reactive<SubjectCreate>({
  name: '',
  code: '',
  description: '',
  department: '',
  stage: '',
  status: 'active',
})

const errors = reactive<Partial<Record<keyof SubjectCreate, string>>>({})
const submitting = ref(false)

const stages = ['الأول', 'الثاني', 'الثالث', 'الرابع', 'الخامس', 'السادس']

function validate(): boolean {
  Object.keys(errors).forEach(k => delete (errors as any)[k])
  if (!form.name?.trim()) errors.name = 'اسم المادة مطلوب'
  if (!form.code?.trim()) errors.code = 'كود المادة مطلوب'
  return Object.keys(errors).length === 0
}

async function handleSubmit() {
  if (!validate()) return
  submitting.value = true
  try {
    await createSubject.mutateAsync({
      name: form.name.trim(),
      code: form.code.trim().toUpperCase(),
      description: form.description?.trim() || undefined,
      department: form.department?.trim() || undefined,
      stage: form.stage || undefined,
      status: form.status,
    })
    toast.success('تم الإنشاء', 'تم إنشاء المادة الدراسية بنجاح')
    router.push({ name: 'admin.subjects' })
  } catch (err: any) {
    const msg = err?.response?.data?.detail ?? 'حدث خطأ أثناء إنشاء المادة'
    toast.error('خطأ', msg)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div dir="rtl" class="max-w-2xl mx-auto space-y-6 animate-fade-in">
    <ToastContainer />

    <!-- Header -->
    <div class="flex items-center gap-4">
      <button
        @click="router.push({ name: 'admin.subjects' })"
        class="btn-ghost border border-slate-200 dark:border-slate-700 p-2"
      >
        <ArrowRight class="h-5 w-5" />
      </button>
      <div>
        <h1 class="text-2xl font-extrabold text-slate-900 dark:text-white">مادة جديدة</h1>
        <p class="text-sm text-slate-500 dark:text-slate-400 mt-0.5">إضافة مادة دراسية جديدة للمنصة</p>
      </div>
    </div>

    <!-- Form -->
    <form @submit.prevent="handleSubmit" class="card p-6 space-y-5">
      <!-- Name -->
      <div class="space-y-1.5">
        <label class="text-sm font-bold text-slate-700 dark:text-slate-200">
          اسم المادة <span class="text-red-500">*</span>
        </label>
        <input
          v-model="form.name"
          type="text"
          placeholder="مثال: الرياضيات التطبيقية"
          :class="[
            'w-full px-4 py-3 rounded-xl border text-sm transition-colors focus:outline-none focus:ring-2 focus:ring-orange-400/50',
            errors.name
              ? 'border-red-300 dark:border-red-700 bg-red-50 dark:bg-red-950/20'
              : 'border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 focus:border-orange-400',
          ]"
        />
        <p v-if="errors.name" class="text-xs text-red-500">{{ errors.name }}</p>
      </div>

      <!-- Code -->
      <div class="space-y-1.5">
        <label class="text-sm font-bold text-slate-700 dark:text-slate-200">
          الكود <span class="text-red-500">*</span>
        </label>
        <input
          v-model="form.code"
          type="text"
          placeholder="مثال: MATH101"
          dir="ltr"
          :class="[
            'w-full px-4 py-3 rounded-xl border text-sm font-mono transition-colors focus:outline-none focus:ring-2 focus:ring-orange-400/50',
            errors.code
              ? 'border-red-300 dark:border-red-700 bg-red-50 dark:bg-red-950/20'
              : 'border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 focus:border-orange-400',
          ]"
        />
        <p v-if="errors.code" class="text-xs text-red-500">{{ errors.code }}</p>
        <p class="text-xs text-slate-400">سيتم تحويله تلقائياً إلى أحرف كبيرة</p>
      </div>

      <!-- Description -->
      <div class="space-y-1.5">
        <label class="text-sm font-bold text-slate-700 dark:text-slate-200">الوصف</label>
        <textarea
          v-model="form.description"
          rows="3"
          placeholder="وصف مختصر للمادة..."
          class="w-full px-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-orange-400/50 focus:border-orange-400 transition-colors"
        />
      </div>

      <!-- Department + Stage -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div class="space-y-1.5">
          <label class="text-sm font-bold text-slate-700 dark:text-slate-200">القسم</label>
          <input
            v-model="form.department"
            type="text"
            placeholder="مثال: علوم الحاسوب"
            class="w-full px-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-sm focus:outline-none focus:ring-2 focus:ring-orange-400/50 focus:border-orange-400 transition-colors"
          />
        </div>
        <div class="space-y-1.5">
          <label class="text-sm font-bold text-slate-700 dark:text-slate-200">المرحلة</label>
          <select
            v-model="form.stage"
            class="w-full px-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-sm focus:outline-none focus:ring-2 focus:ring-orange-400/50 text-slate-700 dark:text-slate-200"
          >
            <option value="">— اختر المرحلة —</option>
            <option v-for="stage in stages" :key="stage" :value="stage">{{ stage }}</option>
          </select>
        </div>
      </div>

      <!-- Status -->
      <div class="space-y-1.5">
        <label class="text-sm font-bold text-slate-700 dark:text-slate-200">الحالة</label>
        <div class="flex gap-3">
          <label
            v-for="s in [{ value: 'active', label: 'نشط' }, { value: 'inactive', label: 'غير نشط' }]"
            :key="s.value"
            :class="[
              'flex-1 flex items-center gap-3 p-3 rounded-xl border cursor-pointer transition-all',
              form.status === s.value
                ? 'border-orange-400 bg-orange-50 dark:bg-orange-950/20'
                : 'border-slate-200 dark:border-slate-700 hover:border-slate-300',
            ]"
          >
            <input type="radio" :value="s.value" v-model="form.status" class="accent-orange-500" />
            <span class="text-sm font-medium">{{ s.label }}</span>
          </label>
        </div>
      </div>

      <!-- Submit -->
      <div class="flex gap-3 pt-2">
        <button
          type="submit"
          :disabled="submitting"
          class="flex-1 btn-primary gap-2 py-3"
        >
          <Loader2 v-if="submitting" class="h-4 w-4 animate-spin" />
          <FolderOpen v-else class="h-4 w-4" />
          {{ submitting ? 'جارٍ الإنشاء...' : 'إنشاء المادة' }}
        </button>
        <button
          type="button"
          @click="router.push({ name: 'admin.subjects' })"
          class="btn-ghost border border-slate-200 dark:border-slate-700 px-6"
        >
          إلغاء
        </button>
      </div>
    </form>
  </div>
</template>
