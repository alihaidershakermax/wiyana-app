<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useSubject, useUpdateSubject } from '@/admin/queries/useSubjects'
import { useToast } from '@/admin/composables/useToast'
import ToastContainer from '@/admin/components/ui/ToastContainer.vue'
import { ArrowRight, FolderOpen, Loader2 } from '@lucide/vue'
import type { SubjectUpdate } from '@/admin/types/subject'

const router = useRouter()
const route = useRoute()
const toast = useToast()

const subjectId = ref(route.params.id as string)
const { data: subject, isLoading, isError } = useSubject(subjectId)
const updateSubject = useUpdateSubject()

const form = reactive<Required<SubjectUpdate>>({
  name: '',
  code: '',
  description: '',
  department: '',
  stage: '',
  status: 'active',
})

const errors = reactive<Partial<Record<keyof SubjectUpdate, string>>>({})
const submitting = ref(false)

// Populate form when data loads
watch(subject, (val) => {
  if (val) {
    form.name = val.name ?? ''
    form.code = val.code ?? ''
    form.description = val.description ?? ''
    form.department = val.department ?? ''
    form.stage = val.stage ?? ''
    form.status = val.status ?? 'active'
  }
}, { immediate: true })

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
    await updateSubject.mutateAsync({
      id: subjectId.value,
      payload: {
        name: form.name.trim(),
        code: form.code.trim().toUpperCase(),
        description: form.description?.trim() || undefined,
        department: form.department?.trim() || undefined,
        stage: form.stage || undefined,
        status: form.status,
      },
    })
    toast.success('تم التحديث', 'تم تحديث المادة الدراسية بنجاح')
    router.push({ name: 'admin.subjects' })
  } catch (err: any) {
    toast.error('خطأ', err?.response?.data?.detail ?? 'فشل تحديث المادة')
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
      <button @click="router.push({ name: 'admin.subjects' })" class="btn-ghost border border-slate-200 dark:border-slate-700 p-2">
        <ArrowRight class="h-5 w-5" />
      </button>
      <div>
        <h1 class="text-2xl font-extrabold text-slate-900 dark:text-white">تعديل المادة</h1>
        <p class="text-sm text-slate-500 dark:text-slate-400 mt-0.5">تحديث بيانات المادة الدراسية</p>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="card p-6 space-y-5 animate-pulse">
      <div v-for="i in 5" :key="i" class="space-y-1.5">
        <div class="h-4 bg-slate-200 dark:bg-slate-700 rounded w-1/4" />
        <div class="h-12 bg-slate-200 dark:bg-slate-700 rounded-xl" />
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="isError" class="card p-8 text-center border-red-200 bg-red-50 dark:bg-red-950/20">
      <p class="text-red-600 dark:text-red-400 font-medium mb-3">لم يتم العثور على المادة</p>
      <button @click="router.push({ name: 'admin.subjects' })" class="btn-primary">العودة للقائمة</button>
    </div>

    <!-- Form -->
    <form v-else @submit.prevent="handleSubmit" class="card p-6 space-y-5">
      <div class="space-y-1.5">
        <label class="text-sm font-bold text-slate-700 dark:text-slate-200">اسم المادة <span class="text-red-500">*</span></label>
        <input
          v-model="form.name"
          type="text"
          :class="['w-full px-4 py-3 rounded-xl border text-sm transition-colors focus:outline-none focus:ring-2 focus:ring-orange-400/50', errors.name ? 'border-red-300' : 'border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 focus:border-orange-400']"
        />
        <p v-if="errors.name" class="text-xs text-red-500">{{ errors.name }}</p>
      </div>

      <div class="space-y-1.5">
        <label class="text-sm font-bold text-slate-700 dark:text-slate-200">الكود <span class="text-red-500">*</span></label>
        <input
          v-model="form.code"
          type="text"
          dir="ltr"
          :class="['w-full px-4 py-3 rounded-xl border text-sm font-mono transition-colors focus:outline-none focus:ring-2 focus:ring-orange-400/50', errors.code ? 'border-red-300' : 'border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 focus:border-orange-400']"
        />
        <p v-if="errors.code" class="text-xs text-red-500">{{ errors.code }}</p>
      </div>

      <div class="space-y-1.5">
        <label class="text-sm font-bold text-slate-700 dark:text-slate-200">الوصف</label>
        <textarea
          v-model="form.description"
          rows="3"
          class="w-full px-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-orange-400/50 focus:border-orange-400"
        />
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div class="space-y-1.5">
          <label class="text-sm font-bold text-slate-700 dark:text-slate-200">القسم</label>
          <input
            v-model="form.department"
            type="text"
            class="w-full px-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-sm focus:outline-none focus:ring-2 focus:ring-orange-400/50 focus:border-orange-400"
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

      <div class="space-y-1.5">
        <label class="text-sm font-bold text-slate-700 dark:text-slate-200">الحالة</label>
        <div class="flex gap-3">
          <label
            v-for="s in [{ value: 'active', label: 'نشط' }, { value: 'inactive', label: 'غير نشط' }, { value: 'archived', label: 'مؤرشف' }]"
            :key="s.value"
            :class="['flex-1 flex items-center gap-2 p-3 rounded-xl border cursor-pointer transition-all text-sm font-medium', form.status === s.value ? 'border-orange-400 bg-orange-50 dark:bg-orange-950/20' : 'border-slate-200 dark:border-slate-700']"
          >
            <input type="radio" :value="s.value" v-model="form.status" class="accent-orange-500" />
            {{ s.label }}
          </label>
        </div>
      </div>

      <div class="flex gap-3 pt-2">
        <button type="submit" :disabled="submitting" class="flex-1 btn-primary gap-2 py-3">
          <Loader2 v-if="submitting" class="h-4 w-4 animate-spin" />
          <FolderOpen v-else class="h-4 w-4" />
          {{ submitting ? 'جارٍ الحفظ...' : 'حفظ التغييرات' }}
        </button>
        <button type="button" @click="router.push({ name: 'admin.subjects' })" class="btn-ghost border border-slate-200 dark:border-slate-700 px-6">
          إلغاء
        </button>
      </div>
    </form>
  </div>
</template>
