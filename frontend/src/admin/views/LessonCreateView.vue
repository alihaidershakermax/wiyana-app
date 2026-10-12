<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  UploadCloud,
  BookOpen,
  FolderOpen,
  Film,
  Image,
  Paperclip,
  CheckCircle2,
  FileText,
  Trash2,
  Sparkles,
  Unlock,
  Layers,
} from '@lucide/vue'
import { useCourses } from '@/admin/queries/useCourses'
import { useSections } from '@/admin/queries/useSections'
import { useCreateLesson } from '@/admin/queries/useLessons'
import { useUploadLessonVideo } from '@/admin/queries/useVideoUpload'
import { useToast } from '@/admin/composables/useToast'
import ToastContainer from '@/admin/components/ui/ToastContainer.vue'
import UploadProgress from '@/admin/components/UploadProgress.vue'

const router = useRouter()
const route = useRoute()
const toast = useToast()

// Data Fetching
const { data: rawCourses } = useCourses()
const coursesList = computed(() => {
  const c = rawCourses.value
  if (!c) return []
  if (Array.isArray(c)) return c
  if (Array.isArray(c.data)) return c.data
  return []
})

const selectedCourseId = ref<string>('')
const selectedSectionId = ref<string>('')

// Initialize from route query or params if provided
onMounted(() => {
  const qCourse = (route.query.courseId || route.params.courseId) as string
  const qSection = (route.query.sectionId || route.params.sectionId) as string
  if (qCourse) selectedCourseId.value = qCourse
  if (qSection) selectedSectionId.value = qSection
})

watch(coursesList, (list) => {
  if (list.length > 0 && !selectedCourseId.value) {
    selectedCourseId.value = list[0].id || list[0]._id || ''
  }
}, { immediate: true })

const assignedCourse = computed(() =>
  coursesList.value.find((c: any) => c.id === selectedCourseId.value || c._id === selectedCourseId.value)
)

const { data: sectionsList } = useSections(selectedCourseId)

// Form State
const title = ref('')
const description = ref('')
const isFree = ref(false)
const status = ref<'published' | 'draft'>('published')

// Video File State
const videoFileInput = ref<HTMLInputElement | null>(null)
const selectedVideo = ref<File | null>(null)

// Thumbnail State
const thumbnailFileInput = ref<HTMLInputElement | null>(null)
const selectedThumbnail = ref<File | null>(null)
const thumbnailPreview = ref<string | null>(null)

// Attachments State
const attachmentFileInput = ref<HTMLInputElement | null>(null)
const attachments = ref<File[]>([])

// Auto-select section when sections load
watch(
  sectionsList,
  (sec) => {
    if (sec && sec.length > 0) {
      if (!selectedSectionId.value || !sec.some(s => s.id === selectedSectionId.value)) {
        selectedSectionId.value = sec[0].id
      }
    }
  },
  { immediate: true }
)

function onVideoChange(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) selectedVideo.value = file
}

function onThumbnailChange(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    selectedThumbnail.value = file
    const reader = new FileReader()
    reader.onload = (ev) => {
      thumbnailPreview.value = ev.target?.result as string
    }
    reader.readAsDataURL(file)
  }
}

function onAttachmentChange(e: Event) {
  const target = e.target as HTMLInputElement
  const files = target.files
  if (files) {
    for (let i = 0; i < files.length; i++) {
      attachments.value.push(files[i])
    }
  }
}

function removeAttachment(index: number) {
  attachments.value.splice(index, 1)
}

// Upload & Save Logic
const createLessonMutation = useCreateLesson()
const videoUploadMutation = useUploadLessonVideo()
const isSubmitting = ref(false)
const isSuccess = ref(false)

async function handlePublishAll() {
  if (!selectedCourseId.value) {
    toast.error('حقل مطلوب', 'يرجى اختيار المساق أولاً')
    return
  }
  if (!title.value.trim()) {
    toast.error('حقل مطلوب', 'يرجى كتابة عنوان المحاضرة')
    return
  }

  const targetSectionId = selectedSectionId.value || undefined

  isSubmitting.value = true
  try {
    // 1. Create Lesson Metadata
    const lessonPayload = {
      title: title.value.trim(),
      description: description.value.trim() || undefined,
      course_id: selectedCourseId.value,
      section_id: targetSectionId,
      thumbnail: thumbnailPreview.value || undefined,
      status: status.value,
      is_free: isFree.value,
    }

    const lesson = await createLessonMutation.mutateAsync({
      sectionId: targetSectionId,
      courseId: selectedCourseId.value,
      input: lessonPayload,
    })

    // 2. Upload Video if selected
    if (selectedVideo.value && lesson?.id) {
      await videoUploadMutation.mutateAsync({
        lessonId: lesson.id,
        file: selectedVideo.value,
        title: title.value,
        privacyStatus: 'unlisted',
      })
    }

    // 3. Upload Attachments if any
    if (attachments.value.length > 0 && lesson?.id) {
      const { adminApiClient } = await import('@/admin/api/adminApi')
      for (const attFile of attachments.value) {
        try {
          const form = new FormData()
          form.append('title', attFile.name)
          form.append('file', attFile)
          await adminApiClient.post(`/panel/lessons/${lesson.id}/attachments`, form, {
            headers: { 'Content-Type': 'multipart/form-data' }
          })
        } catch {
          // continue
        }
      }
    }

    isSuccess.value = true
    toast.success('تمت العملية بنجاح', 'تم رفع ونشر المحاضرة والمرفقات بنجاح')
  } catch (err: any) {
    toast.error('حدث خطأ', err?.message || 'فشل نشر المحاضرة')
  } finally {
    isSubmitting.value = false
  }
}

function resetForm() {
  title.value = ''
  description.value = ''
  selectedVideo.value = null
  selectedThumbnail.value = null
  thumbnailPreview.value = null
  attachments.value = []
  isSuccess.value = false
}
</script>

<template>
  <div dir="rtl" class="max-w-4xl mx-auto space-y-6 animate-fade-in pb-12 text-[#F5F5F5]">
    <ToastContainer />

    <!-- Header -->
    <div class="flex items-center justify-between border-b border-[#242424] pb-4">
      <div>
        <h1 class="text-2xl font-black text-white flex items-center gap-2.5">
          <Sparkles class="h-6 w-6 text-[#F97316]" />
          مركز رفع ونشر المحاضرات للمساق
        </h1>
        <p class="text-xs text-[#A1A1AA] mt-1">
          ارفع الفيديو والغلاف والمرفقات وانشر المحاضرة مباشرة داخل قسم المساق.
        </p>
      </div>
    </div>

    <!-- Success Screen -->
    <div v-if="isSuccess" class="card p-8 text-center space-y-4 border-emerald-500/40 bg-emerald-950/20">
      <div class="h-16 w-16 bg-emerald-500/20 rounded-full flex items-center justify-center mx-auto text-emerald-400">
        <CheckCircle2 class="h-10 w-10" />
      </div>
      <h2 class="text-xl font-black text-white">تم رفع ونشر المحاضرة بنجاح!</h2>
      <p class="text-xs text-slate-300">أصبحت المحاضرة متاحة الآن للطلاب في القسم المختار من المساق.</p>
      <div class="flex justify-center gap-3 pt-2">
        <button class="btn-primary" @click="resetForm">
          رفع محاضرة أخرى
        </button>
        <button class="btn-secondary" @click="router.push(selectedCourseId ? `/admin/courses/${selectedCourseId}` : '/admin/courses')">
          العودة للمساق
        </button>
      </div>
    </div>

    <!-- Unified Form Wizard -->
    <div v-else class="space-y-6">
      <!-- Course & Section Selection -->
      <div class="card p-6 space-y-4">
        <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2">
          <BookOpen class="h-4 w-4 text-[#F97316]" />
          1. تحديد المساق والقسم
        </h3>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <!-- Select Course -->
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-300">المساق *</label>
            <select
              v-model="selectedCourseId"
              class="input w-full text-sm font-semibold"
            >
              <option value="" disabled>-- اختر المساق --</option>
              <option v-for="c in coursesList" :key="c.id || c._id" :value="c.id || c._id">
                {{ c.title }} ({{ c.stage || 'عام' }})
              </option>
            </select>
          </div>

          <!-- Select Section -->
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-300">القسم داخل المساق *</label>
            <select
              v-model="selectedSectionId"
              class="input w-full text-sm font-semibold"
            >
              <option v-if="!sectionsList || sectionsList.length === 0" value="">-- لا توجد أقسام (سيتم إنشاء قسم افتراضي) --</option>
              <option v-for="s in sectionsList" :key="s.id" :value="s.id">
                {{ s.title }}
              </option>
            </select>
          </div>
        </div>
      </div>

      <!-- Step 2: Lesson Basic Details -->
      <div class="card p-6 space-y-4">
        <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2">
          <FileText class="h-4 w-4 text-[#F97316]" />
          2. تفاصيل المحاضرة
        </h3>

        <div class="space-y-1.5">
          <label class="text-xs font-bold text-slate-300">عنوان المحاضرة *</label>
          <input
            v-model="title"
            type="text"
            placeholder="مثال: المحاضرة الأولى - المتسعات وتطبيقاتها"
            class="input w-full text-sm"
          />
        </div>

        <div class="space-y-1.5">
          <label class="text-xs font-bold text-slate-300">وصف أو ملخص المحاضرة</label>
          <textarea
            v-model="description"
            rows="3"
            placeholder="اكتب ملخصاً أو ملاحظات للطلاب حول محتوى هذه المحاضرة..."
            class="input w-full text-sm resize-none"
          />
        </div>

        <div class="flex items-center justify-between p-3 rounded-xl border border-slate-800 bg-slate-900/50">
          <span class="text-xs font-bold text-slate-200 flex items-center gap-2">
            <Unlock class="h-4 w-4 text-emerald-400" />
            إتاحة كـ محاضرة مجانية تجريبية للجميع
          </span>
          <input v-model="isFree" type="checkbox" class="h-4 w-4 accent-orange-500 cursor-pointer" />
        </div>
      </div>

      <!-- Step 3: Media Upload (Video & Thumbnail) -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Video Upload Block -->
        <div class="card p-6 space-y-4">
          <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2">
            <Film class="h-4 w-4 text-[#F97316]" />
            3. ملف الفيديو
          </h3>

          <input ref="videoFileInput" type="file" accept="video/*" class="hidden" @change="onVideoChange" />

          <div
            v-if="!selectedVideo"
            class="border-2 border-dashed border-slate-800 hover:border-orange-500 rounded-2xl p-6 text-center cursor-pointer transition-colors"
            @click="videoFileInput?.click()"
          >
            <UploadCloud class="h-8 w-8 text-slate-500 mx-auto mb-2" />
            <p class="text-xs font-bold text-slate-300">اختر ملف فيديو من جهازك</p>
            <p class="text-[10px] text-slate-500 mt-1">MP4, MOV, MKV, WebM</p>
          </div>

          <div v-else class="space-y-2 p-3 bg-slate-900/60 rounded-xl border border-slate-800 flex items-center justify-between">
            <div class="min-w-0">
              <p class="text-xs font-bold text-white truncate">{{ selectedVideo.name }}</p>
              <p class="text-[10px] text-slate-400">{{ (selectedVideo.size / 1024 / 1024).toFixed(1) }} MB</p>
            </div>
            <button class="text-xs text-red-400 hover:text-red-300 cursor-pointer" @click="selectedVideo = null">إزالة</button>
          </div>
        </div>

        <!-- Thumbnail Upload Block -->
        <div class="card p-6 space-y-4">
          <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2">
            <Image class="h-4 w-4 text-[#F97316]" />
            4. صورة الغلاف
          </h3>

          <input ref="thumbnailFileInput" type="file" accept="image/*" class="hidden" @change="onThumbnailChange" />

          <div
            v-if="!thumbnailPreview"
            class="border-2 border-dashed border-slate-800 hover:border-orange-500 rounded-2xl p-6 text-center cursor-pointer transition-colors"
            @click="thumbnailFileInput?.click()"
          >
            <Image class="h-8 w-8 text-slate-500 mx-auto mb-2" />
            <p class="text-xs font-bold text-slate-300">رفع صورة الغلاف</p>
            <p class="text-[10px] text-slate-500 mt-1">PNG, JPG, WebP</p>
          </div>

          <div v-else class="relative rounded-xl overflow-hidden border border-slate-800 h-32 group">
            <img :src="thumbnailPreview" class="w-full h-full object-cover" />
            <button
              class="absolute top-2 right-2 p-1 bg-red-600/80 text-white rounded-lg opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer"
              @click="thumbnailPreview = null"
            >
              <Trash2 class="h-4 w-4" />
            </button>
          </div>
        </div>
      </div>

      <!-- Step 4: Attachments Upload -->
      <div class="card p-6 space-y-4">
        <div class="flex items-center justify-between">
          <h3 class="text-sm font-bold text-slate-200 flex items-center gap-2">
            <Paperclip class="h-4 w-4 text-[#F97316]" />
            5. المرفقات والملزمات (اختياري)
          </h3>
          <button class="btn-secondary text-xs cursor-pointer" @click="attachmentFileInput?.click()">
            + إضافة ملفات
          </button>
        </div>

        <input ref="attachmentFileInput" type="file" multiple class="hidden" @change="onAttachmentChange" />

        <div v-if="attachments.length" class="space-y-2">
          <div
            v-for="(file, idx) in attachments"
            :key="idx"
            class="flex items-center justify-between p-3 rounded-xl border border-slate-800 bg-slate-900/40 text-xs"
          >
            <span class="font-bold text-slate-300 truncate">{{ file.name }}</span>
            <button class="text-red-400 hover:text-red-300 cursor-pointer" @click="removeAttachment(idx)">حذف</button>
          </div>
        </div>

        <p v-else class="text-xs text-slate-500 text-center py-2">
          يمكنك رفع ملزمة PDF، أسئلة امتحانات، أو ملفات مساعدة للطلاب.
        </p>
      </div>

      <!-- Video Upload Progress Bar -->
      <UploadProgress
        :percent="videoUploadMutation.progress.value"
        :state="videoUploadMutation.uploadState.value"
        :error="videoUploadMutation.errorMessage.value"
      />

      <!-- Submit Action Button -->
      <button
        class="btn-primary w-full py-4 text-base font-black tracking-wide shadow-lg shadow-orange-500/20 active:scale-[0.99] transition-transform cursor-pointer"
        :disabled="isSubmitting || !title.trim()"
        @click="handlePublishAll"
      >
        {{ isSubmitting ? 'جارٍ نشر المحاضرة...' : 'نشر المحاضرة في المساق' }}
      </button>
    </div>
  </div>
</template>
