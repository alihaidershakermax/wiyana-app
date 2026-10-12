<script setup lang="ts">
import { reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useCreateCourse } from '@/admin/queries/useCourses'
import { useSubjects } from '@/admin/queries/useSubjects'
import { useToast } from '@/admin/composables/useToast'
import ToastContainer from '@/admin/components/ui/ToastContainer.vue'
import { ArrowRight, BookOpen, Loader2, Upload, X } from '@lucide/vue'
import type { CourseCreate } from '@/admin/types/course'
import { adminApiClient } from '@/admin/api/adminApi'

const router = useRouter()
const toast = useToast()
const createCourse = useCreateCourse()

// Load subjects for the select
const { data: subjectsData } = useSubjects({ limit: 100, status: 'active' })
const subjects = computed(() => subjectsData.value?.data ?? [])

const form = reactive<CourseCreate>({
  title: '',
  description: '',
  subject_id: '',
  department: '',
  stage: '',
  price: 0,
  status: 'draft',
  thumbnail_url: '',
})

const errors = reactive<Partial<Record<keyof CourseCreate, string>>>({})
const submitting = ref(false)

// Cover image upload
const coverFile = ref<File | null>(null)
const coverPreview = ref<string>('')
const uploadingCover = ref(false)

function onCoverSelected(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  if (file.size > 5 * 1024 * 1024) {
    toast.error('خطأ', 'حجم الصورة يتجاوز 5MB')
    return
  }
  coverFile.value = file
  const reader = new FileReader()
  reader.onload = () => { coverPreview.value = reader.result as string }
  reader.readAsDataURL(file)
}

function removeCover() {
  coverFile.value = null
  coverPreview.value = ''
  form.thumbnail_url = ''
}

async function uploadCoverToImgbb(file: File): Promise<string | null> {
  const formData = new FormData()
  formData.append('image', file)
  try {
    const resp = await adminApiClient.post('/images/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    return resp.data?.url || null
  } catch (err: any) {
    const detail = err?.response?.data?.detail || err?.message || 'فشل رفع الصورة'
    toast.error('خطأ في رفع الصورة', detail)
    return null
  }
}

function validate() {
  Object.keys(errors).forEach(k => delete (errors as any)[k])
  if (!form.title?.trim()) errors.title = 'عنوان المساق مطلوب'
  return Object.keys(errors).length === 0
}

async function handleSubmit() {
  if (!validate()) return
  submitting.value = true
  try {
    // Upload cover image if selected
    if (coverFile.value) {
      uploadingCover.value = true
      const url = await uploadCoverToImgbb(coverFile.value)
      uploadingCover.value = false
      if (url) form.thumbnail_url = url
    }

    const course = await createCourse.mutateAsync({
      title: form.title.trim(),
      description: form.description?.trim() || undefined,
      thumbnail_url: form.thumbnail_url || undefined,
      subject_id: form.subject_id || undefined,
      department: form.department?.trim() || undefined,
      stage: form.stage || undefined,
      price: form.price ?? 0,
      status: form.status ?? 'draft',
    })
    toast.success('تم الإنشاء', 'تم إنشاء المساق بنجاح')
    router.push({ name: 'admin.course-detail', params: { id: course.id } })
  } catch (err: any) {
    toast.error('خطأ', err?.response?.data?.detail ?? 'فشل إنشاء المساق')
  } finally {
    submitting.value = false
    uploadingCover.value = false
  }
}
</script>

<template>
  <div dir="rtl" class="max-w-2xl mx-auto space-y-6 animate-fade-in">
    <ToastContainer />

    <div class="flex items-center gap-4">
      <button @click="router.push({ name: 'admin.courses' })" class="btn-ghost border border-slate-200 dark:border-slate-700 p-2">
        <ArrowRight class="h-5 w-5" />
      </button>
      <div>
        <h1 class="text-2xl font-extrabold text-slate-900 dark:text-white">مساق جديد</h1>
        <p class="text-sm text-slate-500 dark:text-slate-400 mt-0.5">إنشاء مساق تعليمي جديد</p>
      </div>
    </div>

    <form @submit.prevent="handleSubmit" class="card p-6 space-y-5">
      <!-- Title -->
      <div class="space-y-1.5">
        <label class="text-sm font-bold text-slate-700 dark:text-slate-200">عنوان المساق <span class="text-red-500">*</span></label>
        <input
          v-model="form.title"
          type="text"
          placeholder="مثال: مقدمة في الرياضيات"
          :class="['w-full px-4 py-3 rounded-xl border text-sm transition-colors focus:outline-none focus:ring-2 focus:ring-orange-400/50', errors.title ? 'border-red-300' : 'border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 focus:border-orange-400']"
        />
        <p v-if="errors.title" class="text-xs text-red-500">{{ errors.title }}</p>
      </div>

      <!-- Description -->
      <div class="space-y-1.5">
        <label class="text-sm font-bold text-slate-700 dark:text-slate-200">الوصف</label>
        <textarea
          v-model="form.description"
          rows="3"
          placeholder="وصف مختصر للمساق..."
          class="w-full px-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-orange-400/50 focus:border-orange-400"
        />
      </div>

      <!-- Cover Image Upload -->
      <div class="space-y-1.5">
        <label class="text-sm font-bold text-slate-700 dark:text-slate-200">صورة الغلاف</label>
        <div v-if="coverPreview" class="relative">
          <img :src="coverPreview" class="w-full h-48 object-cover rounded-xl border border-slate-200 dark:border-slate-700" />
          <button type="button" @click="removeCover" class="absolute top-2 left-2 p-1.5 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors">
            <X class="w-4 h-4" />
          </button>
        </div>
        <label v-else class="flex flex-col items-center justify-center w-full h-32 border-2 border-dashed border-slate-300 dark:border-slate-600 rounded-xl cursor-pointer hover:border-orange-400 transition-colors bg-slate-50 dark:bg-slate-800/50">
          <Upload class="w-8 h-8 text-slate-400 mb-2" />
          <span class="text-sm text-slate-500">اضغط لرفع صورة الغلاف</span>
          <span class="text-xs text-slate-400 mt-1">PNG, JPEG, WEBP — حد أقصى 5MB</span>
          <input type="file" accept="image/*" class="hidden" @change="onCoverSelected" />
        </label>
      </div>



      <!-- Price -->
      <div class="space-y-1.5">
        <label class="text-sm font-bold text-slate-700 dark:text-slate-200">السعر ($)</label>
        <input
          v-model.number="form.price"
          type="number"
          min="0"
          step="0.01"
          dir="ltr"
          class="w-full px-4 py-3 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-sm focus:outline-none focus:ring-2 focus:ring-orange-400/50"
        />
      </div>

      <!-- Status -->
      <div class="space-y-1.5">
        <label class="text-sm font-bold text-slate-700 dark:text-slate-200">الحالة الأولية</label>
        <div class="flex gap-3">
          <label
            v-for="s in [{ value: 'draft', label: 'مسودة' }, { value: 'published', label: 'منشور مباشرة' }]"
            :key="s.value"
            :class="['flex-1 flex items-center gap-3 p-3 rounded-xl border cursor-pointer transition-all text-sm font-medium', form.status === s.value ? 'border-orange-400 bg-orange-50 dark:bg-orange-950/20' : 'border-slate-200 dark:border-slate-700']"
          >
            <input type="radio" :value="s.value" v-model="form.status" class="accent-orange-500" />
            {{ s.label }}
          </label>
        </div>
      </div>

      <div class="flex gap-3 pt-2">
        <button type="submit" :disabled="submitting" class="flex-1 btn-primary gap-2 py-3">
          <Loader2 v-if="submitting" class="h-4 w-4 animate-spin" />
          <BookOpen v-else class="h-4 w-4" />
          {{ submitting ? 'جارٍ الإنشاء...' : 'إنشاء المساق' }}
        </button>
        <button type="button" @click="router.push({ name: 'admin.courses' })" class="btn-ghost border border-slate-200 dark:border-slate-700 px-6">إلغاء</button>
      </div>
    </form>
  </div>
</template>
