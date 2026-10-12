<template>
  <Teleport to="body">
    <div
      v-if="open"
      class="fixed inset-0 z-[9999] flex items-center justify-center bg-black/90 backdrop-blur-md p-4 animate-in fade-in duration-200"
      dir="rtl"
    >
      <div
        class="w-full max-w-md rounded-2xl p-6 shadow-2xl text-center space-y-5 border"
        :class="levelClasses.container"
      >
        <div
          class="mx-auto w-14 h-14 rounded-full flex items-center justify-center"
          :class="levelClasses.iconWrap"
        >
          <ShieldAlert class="w-8 h-8" :class="levelClasses.icon" />
        </div>

        <div class="space-y-2">
          <h3 class="text-xl font-bold text-white tracking-tight">
            {{ title || 'تنبيه أمني' }}
          </h3>
          <p class="text-xs font-medium" :class="levelClasses.accent">
            {{ message }}
          </p>
        </div>

        <p
          class="text-sm text-neutral-300 leading-relaxed bg-neutral-950/80 p-3.5 rounded-xl border border-neutral-800"
        >
          {{ detail }}
        </p>

        <div class="pt-2 space-y-2">
          <button
            type="button"
            @click="onPrimary"
            class="w-full py-3 px-5 rounded-xl text-black font-bold text-sm shadow-lg active:scale-98 transition-all cursor-pointer"
            :class="levelClasses.button"
          >
            {{ primaryLabel }}
          </button>

          <p v-if="requireRevalidation" class="text-[11px] text-neutral-400">
            تتطلب المتابعة إعادة التحقق من جلسة العرض.
          </p>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ShieldAlert } from '@lucide/vue'
import type { RiskLevel } from '@/security/riskEngine'

const props = withDefaults(
  defineProps<{
    open: boolean
    level?: RiskLevel
    title?: string
    message?: string
    detail?: string
    requireRevalidation?: boolean
  }>(),
  {
    level: 'HIGH_RISK',
    title: '',
    message: 'تم اكتشاف نشاط مشبوه. يرجى إيقاف أي برامج تسجيل أو مشاركة شاشة للمتابعة.',
    detail: '',
    requireRevalidation: false,
  }
)

const emit = defineEmits<{
  (e: 'dismiss'): void
  (e: 'revalidate'): void
}>()

const isCritical = computed(() => props.level === 'CRITICAL')
const isHigh = computed(() => props.level === 'HIGH_RISK')
const isAdvisory = computed(() => props.level === 'SUSPICIOUS')

const levelClasses = computed(() => {
  if (isCritical.value) {
    return {
      container: 'bg-neutral-900 border-red-500/40',
      iconWrap: 'bg-red-500/10 border border-red-500/30',
      icon: 'text-red-500',
      accent: 'text-red-400',
      button: 'bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 shadow-red-500/20',
    }
  }
  if (isHigh.value) {
    return {
      container: 'bg-neutral-900 border-amber-500/40',
      iconWrap: 'bg-amber-500/10 border border-amber-500/30',
      icon: 'text-amber-400',
      accent: 'text-amber-400',
      button: 'bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-600 hover:to-amber-700 shadow-amber-500/20',
    }
  }
  // Advisory / suspicious
  return {
    container: 'bg-neutral-900 border-amber-400/30',
    iconWrap: 'bg-amber-400/10 border border-amber-400/30',
    icon: 'text-amber-300',
    accent: 'text-amber-300',
    button: 'bg-gradient-to-r from-amber-400 to-amber-500 hover:from-amber-500 hover:to-amber-600 shadow-amber-500/10',
  }
})

const primaryLabel = computed(() => {
  if (isCritical.value) return 'إعادة فتح المحاضرة'
  if (isAdvisory.value) return 'متابعة'
  return 'متابعة المشاهدة'
})

const onPrimary = () => {
  if (isCritical.value) {
    emit('revalidate')
    return
  }
  emit('dismiss')
}
</script>