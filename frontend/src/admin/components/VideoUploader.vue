<script setup lang="ts">
import { ref, computed } from 'vue'
import { UploadCloud, Video, AlertTriangle, Film } from '@lucide/vue'
import { useUploadLessonVideo, useYouTubeStatus } from '@/admin/queries/useVideoUpload'
import UploadProgress from './UploadProgress.vue'

const props = defineProps<{
  lessonId: string
  existingVideoId?: string | null
}>()

const emit = defineEmits<{ uploaded: [data: { lesson: unknown; youtube: Record<string, unknown> }] }>()

const fileInput = ref<HTMLInputElement | null>(null)
const selectedFile = ref<File | null>(null)
const privacy = ref<'public' | 'unlisted' | 'private'>('unlisted')

const { data: ytStatus, isLoading: statusLoading } = useYouTubeStatus()
const upload = useUploadLessonVideo()

const isSimulated = computed(() => ytStatus.value?.simulation_mode)

function onFileChange(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) selectedFile.value = file
}

async function startUpload() {
  if (!selectedFile.value) return
  upload.mutate(
    {
      lessonId: props.lessonId,
      file: selectedFile.value,
      privacyStatus: privacy.value,
    },
    {
      onSuccess: (data) => emit('uploaded', data),
    }
  )
}

function reset() {
  selectedFile.value = null
  if (fileInput.value) fileInput.value.value = ''
}

const hasExistingVideo = computed(() => !!props.existingVideoId && !props.existingVideoId.startsWith('yt_sim_'))
</script>

<template>
  <div class="space-y-4">
    <!-- YouTube status banner -->
    <div v-if="statusLoading" class="skeleton h-10 rounded-xl" />
    <div
      v-else-if="isSimulated"
      class="flex items-start gap-2 rounded-xl border border-amber-200 dark:border-amber-900/60 bg-amber-50 dark:bg-amber-950/30 p-3"
    >
      <AlertTriangle class="h-4 w-4 text-amber-500 mt-0.5 shrink-0" />
      <p class="text-xs text-amber-700 dark:text-amber-400 leading-relaxed">
        لم يتم إعداد حساب YouTube بعد (لا يوجد توكن OAuth). سيتم حفظ بيانات المحاضرة بدون فيديو حقيقي.
        أضف ملف التوكن في الخادم لتفعيل الرفع الفعلي.
      </p>
    </div>

    <!-- Existing video -->
    <div v-if="hasExistingVideo" class="flex items-center gap-3 rounded-xl border border-emerald-200 dark:border-emerald-900/60 bg-emerald-50 dark:bg-emerald-950/30 p-3">
      <Video class="h-5 w-5 text-emerald-500 shrink-0" />
      <div class="min-w-0">
        <p class="text-xs font-bold text-emerald-700 dark:text-emerald-400">فيديو مرفق على YouTube</p>
        <p class="text-[11px] text-emerald-600/80 dark:text-emerald-500/80 truncate" dir="ltr">
          https://www.youtube.com/watch?v={{ existingVideoId }}
        </p>
      </div>
    </div>

    <!-- Upload form -->
    <div class="rounded-2xl border border-slate-200 dark:border-slate-800 p-4 space-y-4">
      <div class="flex items-center gap-2">
        <Film class="h-4 w-4 text-orange-500" />
        <h4 class="text-sm font-bold">رفع فيديو المحاضرة</h4>
      </div>

      <input ref="fileInput" type="file" accept="video/*" class="hidden" @change="onFileChange" />

      <button
        v-if="!selectedFile"
        class="w-full rounded-xl border-2 border-dashed border-slate-300 dark:border-slate-700 py-8 flex flex-col items-center gap-2 hover:border-orange-400 dark:hover:border-orange-700 transition-colors"
        @click="fileInput?.click()"
      >
        <UploadCloud class="h-8 w-8 text-slate-300 dark:text-slate-600" />
        <span class="text-sm font-bold text-slate-500 dark:text-slate-400">اختر ملف فيديو</span>
        <span class="text-[11px] text-slate-400">MP4, MOV, WebM</span>
      </button>

      <div v-else class="space-y-3">
        <div class="flex items-center gap-3 rounded-xl bg-slate-50 dark:bg-slate-800/60 p-3">
          <Film class="h-5 w-5 text-orange-500 shrink-0" />
          <div class="min-w-0 flex-1">
            <p class="text-xs font-bold truncate">{{ selectedFile.name }}</p>
            <p class="text-[11px] text-slate-400">{{ (selectedFile.size / 1024 / 1024).toFixed(1) }} م.ب</p>
          </div>
          <button class="text-[11px] font-bold text-slate-400 hover:text-red-500" @click="reset">
            إزالة
          </button>
        </div>

        <div class="flex items-center gap-3">
          <label class="text-xs font-bold text-slate-500 w-24 shrink-0">الخصوصية</label>
          <div class="flex gap-2">
            <button
              v-for="opt in [{ v: 'public', l: 'عام' }, { v: 'unlisted', l: 'غير مدرج' }, { v: 'private', l: 'خاص' }]"
              :key="opt.v"
              class="text-xs font-bold px-3 py-1.5 rounded-lg border transition-colors"
              :class="privacy === opt.v ? 'border-orange-400 bg-orange-50 dark:bg-orange-950/40 text-orange-600 dark:text-orange-400' : 'border-slate-200 dark:border-slate-700 text-slate-400 hover:border-slate-300'"
              @click="privacy = opt.v as typeof privacy"
            >
              {{ opt.l }}
            </button>
          </div>
        </div>

        <button
          class="btn-primary w-full py-2.5"
          :disabled="upload.isPending.value"
          @click="startUpload"
        >
          <UploadCloud v-if="!upload.isPending.value" class="h-4 w-4" />
          {{ upload.isPending.value ? 'جارٍ الرفع...' : 'رفع إلى YouTube' }}
        </button>
      </div>

      <UploadProgress
        :percent="upload.progress.value"
        :state="upload.uploadState.value"
        :error="upload.errorMessage.value"
        :uploaded-size="selectedFile ? Math.round((selectedFile.size * upload.progress.value) / 100) : undefined"
        :total-size="selectedFile?.size"
        @cancel="reset"
      />
    </div>
  </div>
</template>

<style scoped>
.skeleton {
  background: linear-gradient(90deg, rgb(241 245 249) 25%, rgb(226 232 240) 50%, rgb(241 245 249) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
}
.dark .skeleton {
  background: linear-gradient(90deg, rgb(30 41 59) 25%, rgb(51 65 85) 50%, rgb(30 41 59) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
}
@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>
