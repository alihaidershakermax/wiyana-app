<script setup lang="ts">
import { reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from '@/admin/composables/useToast'
import ToastContainer from '@/admin/components/ui/ToastContainer.vue'
import { ArrowRight, GraduationCap, Loader2 } from '@lucide/vue'
import { useCreateTeacher, useAssignTeacherToCourse } from '@/admin/queries/useTeachers'
import { useCourses } from '@/admin/queries/useCourses'

const router = useRouter()
const toast = useToast()
const createTeacher = useCreateTeacher()
const assignTeacher = useAssignTeacherToCourse()

const { data: coursesData } = useCourses({ limit: 100 })
const allCourses = computed(() => {
  const c = coursesData.value
  if (!c) return []
  if (Array.isArray(c)) return c
  if (Array.isArray(c.data)) return c.data
  return []
})

const form = reactive<{
  name: string
  username: string
  telegram_id: string
  course_id: string
}>({
  name: '',
  username: '',
  telegram_id: '',
  course_id: '',
})

const errors = reactive<Partial<Record<keyof typeof form, string>>>({})
const submitting = ref(false)

function validate(): boolean {
  Object.keys(errors).forEach(k => delete (errors as any)[k])
  const hasUsername = !!form.username?.trim()
  const hasTelegramId = !!form.telegram_id?.trim()
  if (!hasUsername && !hasTelegramId) {
    errors.username = 'أدخل اسم المستخدم أو معرف تيليجرام'
    errors.telegram_id = 'أدخل اسم المستخدم أو معرف تيليجرام'
  }
  if (!form.course_id) errors.course_id = 'اختر المساق'
  return Object.keys(errors).length === 0
}

async function handleSubmit() {
  if (!validate()) return
  submitting.value = true
  try {
    const teacher = await createTeacher.mutateAsync({
      name: form.name?.trim() || undefined,
      username: form.username?.trim() || undefined,
      telegram_id: form.telegram_id?.trim() ? Number(form.telegram_id.trim()) : undefined,
    })
    if (form.course_id) {
      await assignTeacher.mutateAsync({ courseId: form.course_id, teacherId: teacher.id })
    }
    toast.success('تمت الإضافة', 'تم إضافة الأستاذ وتخصيص المساق بنجاح')
    router.push({ name: 'admin.teachers' })
  } catch (err: any) {
    toast.error('خطأ', err?.response?.data?.detail ?? 'حدث خطأ أثناء إضافة الأستاذ')
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
        @click="router.push({ name: 'admin.teachers' })"
        class="btn-ghost border border-slate-200 dark:border-slate-700 p-2"
      >
        <ArrowRight class="h-5 w-5" />
      </button>
      <div>
        <h1 class="text-2xl font-extrabold text-slate-900 dark:text-white">إضافة أستاذ</h1>
        <p class="text-sm text-slate-500 dark:text-slate-400 mt-0.5">إضافة أستاذ جديد وتخصيص مساق له</p>
      </div>
    </div>

    <!-- Form -->
    <form @submit.prevent="handleSubmit" class="card p-6 space-y-5">
      <!-- Teacher name -->
      <div class="space-y-1.5">
        <label class="text-sm font-bold text-slate-700 dark:text-slate-200">اسم الأستاذ</label>
        <input
          v-model="form.name"
          type="text"
          placeholder="مثال: الأستاذ أحمد"
          class="w-full px-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-sm focus:outline-none focus:ring-2 focus:ring-orange-400/50 focus:border-orange-400 transition-colors"
        />
      </div>

      <!-- Telegram username / ID -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div class="space-y-1.5">
          <label class="text-sm font-bold text-slate-700 dark:text-slate-200">
            يوزر نيم تيليجرام <span class="text-red-500">*</span>
          </label>
          <input
            v-model="form.username"
            type="text"
            placeholder="مثال: @ahmed_teacher"
            dir="ltr"
            :class="[
              'w-full px-4 py-3 rounded-xl border text-sm font-mono transition-colors focus:outline-none focus:ring-2 focus:ring-orange-400/50',
              errors.username
                ? 'border-red-300 dark:border-red-700 bg-red-50 dark:bg-red-950/20'
                : 'border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 focus:border-orange-400',
            ]"
          />
          <p v-if="errors.username" class="text-xs text-red-500">{{ errors.username }}</p>
        </div>
        <div class="space-y-1.5">
          <label class="text-sm font-bold text-slate-700 dark:text-slate-200">
            معرف تيليجرام <span class="text-red-500">*</span>
          </label>
          <input
            v-model="form.telegram_id"
            type="number"
            placeholder="مثال: 123456789"
            dir="ltr"
            :class="[
              'w-full px-4 py-3 rounded-xl border text-sm font-mono transition-colors focus:outline-none focus:ring-2 focus:ring-orange-400/50',
              errors.telegram_id
                ? 'border-red-300 dark:border-red-700 bg-red-50 dark:bg-red-950/20'
                : 'border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 focus:border-orange-400',
            ]"
          />
          <p v-if="errors.telegram_id" class="text-xs text-red-500">{{ errors.telegram_id }}</p>
        </div>
      </div>
      <p class="text-xs text-slate-400 -mt-2">أدخل يوزر نيم تيليجرام أو معرف تيليجرام (الرقمي) على الأقل</p>

      <!-- Course -->
      <div class="space-y-1.5">
        <label class="text-sm font-bold text-slate-700 dark:text-slate-200">
          المساق التعليمي <span class="text-red-500">*</span>
        </label>
        <select
          v-model="form.course_id"
          :class="[
            'w-full px-4 py-3 rounded-xl border text-sm transition-colors focus:outline-none focus:ring-2 focus:ring-orange-400/50',
            errors.course_id
              ? 'border-red-300 dark:border-red-700 bg-red-50 dark:bg-red-950/20'
              : 'border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-200 focus:border-orange-400',
          ]"
        >
          <option value="">— اختر المساق —</option>
          <option v-for="course in allCourses" :key="course.id" :value="course.id">
            {{ course.title }} ({{ course.stage || 'عام' }})
          </option>
        </select>
        <p v-if="errors.course_id" class="text-xs text-red-500">{{ errors.course_id }}</p>
      </div>

      <!-- Submit -->
      <div class="flex gap-3 pt-2">
        <button
          type="submit"
          :disabled="submitting"
          class="flex-1 btn-primary gap-2 py-3"
        >
          <Loader2 v-if="submitting" class="h-4 w-4 animate-spin" />
          <GraduationCap v-else class="h-4 w-4" />
          {{ submitting ? 'جارٍ الإضافة...' : 'إضافة الأستاذ' }}
        </button>
        <button
          type="button"
          @click="router.push({ name: 'admin.teachers' })"
          class="btn-ghost border border-slate-200 dark:border-slate-700 px-6"
        >
          إلغاء
        </button>
      </div>
    </form>
  </div>
</template>
