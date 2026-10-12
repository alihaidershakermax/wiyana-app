<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Sparkles, CheckCircle2, X, ArrowRight } from '@lucide/vue'
import { appApiClient } from '@/api/client'

export interface AppReleaseNotice {
  version: string
  title: string
  features: string[]
  publishedAt?: string
}

const isVisible = ref(false)
const releaseInfo = ref<AppReleaseNotice | null>(null)

const STORAGE_KEY = 'wiyana_seen_app_release_version'

async function checkReleaseNotice() {
  try {
    const { data } = await appApiClient.get('/api/v1/notifications/latest-release')
    const release = data?.data || data
    if (release && release.version) {
      const seenVersion = localStorage.getItem(STORAGE_KEY)
      if (release.version !== seenVersion) {
        releaseInfo.value = release
        isVisible.value = true
        return
      }
    }
  } catch (e) {
    console.warn('Notice fetch skipped:', e)
  }

  // Fallback to local storage
  try {
    const raw = localStorage.getItem('wiyana_latest_app_release')
    if (!raw) return
    const parsed: AppReleaseNotice = JSON.parse(raw)
    const seenVersion = localStorage.getItem(STORAGE_KEY)
    
    if (parsed.version && parsed.version !== seenVersion) {
      releaseInfo.value = parsed
      isVisible.value = true
    }
  } catch (e) {
    console.error('Error parsing app release notice:', e)
  }
}

function dismissModal() {
  if (releaseInfo.value?.version) {
    localStorage.setItem(STORAGE_KEY, releaseInfo.value.version)
  }
  isVisible.value = false
}

// Global event listener to open modal dynamically when broadcasted
onMounted(() => {
  checkReleaseNotice()
  window.addEventListener('wiyana:app-update-published', (event: any) => {
    if (event.detail) {
      releaseInfo.value = event.detail
      isVisible.value = true
    }
  })
})
</script>

<template>
  <Transition
    enter-active-class="transition duration-300 ease-out"
    enter-from-class="opacity-0 scale-95"
    enter-to-class="opacity-100 scale-100"
    leave-active-class="transition duration-200 ease-in"
    leave-from-class="opacity-100 scale-100"
    leave-to-class="opacity-0 scale-95"
  >
    <div
      v-if="isVisible && releaseInfo"
      class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/80 backdrop-blur-md"
    >
      <div class="relative w-full max-w-lg overflow-hidden rounded-3xl bg-slate-900 border border-amber-500/30 p-6 sm:p-8 shadow-2xl shadow-amber-500/10 text-white text-right">
        <!-- Background Gradient Glow -->
        <div class="absolute -top-24 -right-24 w-48 h-48 rounded-full bg-amber-500/20 blur-3xl pointer-events-none"></div>
        <div class="absolute -bottom-24 -left-24 w-48 h-48 rounded-full bg-orange-500/20 blur-3xl pointer-events-none"></div>

        <!-- Header Badge -->
        <div class="flex items-center justify-between mb-4">
          <button
            @click="dismissModal"
            class="p-2 rounded-full text-slate-400 hover:text-white hover:bg-white/10 transition-colors"
          >
            <X class="w-5 h-5" />
          </button>
          <div class="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-gradient-to-r from-amber-500/20 to-orange-500/20 border border-amber-500/40 text-amber-400 text-xs font-bold">
            <Sparkles class="w-4 h-4 animate-spin-slow" />
            <span>تحديث جديد الإصدار v{{ releaseInfo.version }}</span>
          </div>
        </div>

        <!-- Title -->
        <h3 class="text-xl sm:text-2xl font-extrabold text-white mb-2 leading-snug">
          {{ releaseInfo.title || 'مميزات وتحديثات جديدة في منصة ويانا 🎉' }}
        </h3>
        <p class="text-xs sm:text-sm text-slate-300 mb-6 leading-relaxed">
          تم إطلاق تحديث جديد للمنصة يتضمن تحسينات ومميزات فايقة لخدمتكم:
        </p>

        <!-- Features List -->
        <div class="space-y-3 mb-8 bg-slate-800/60 rounded-2xl p-4 border border-white/5">
          <div
            v-for="(feature, idx) in releaseInfo.features"
            :key="idx"
            class="flex items-start gap-3 text-sm text-slate-200"
          >
            <CheckCircle2 class="w-5 h-5 text-amber-400 shrink-0 mt-0.5" />
            <span class="font-medium leading-normal">{{ feature }}</span>
          </div>
        </div>

        <!-- Confirm Button -->
        <button
          @click="dismissModal"
          class="w-full py-3.5 px-6 rounded-2xl bg-gradient-to-r from-amber-500 to-orange-600 hover:from-amber-600 hover:to-orange-700 text-white font-extrabold text-sm sm:text-base shadow-lg shadow-amber-500/25 transition-all transform active:scale-95 flex items-center justify-center gap-2"
        >
          <span>تم الفهم ومتابعة الدراسة 🚀</span>
          <ArrowRight class="w-5 h-5 rotate-180" />
        </button>
      </div>
    </div>
  </Transition>
</template>
