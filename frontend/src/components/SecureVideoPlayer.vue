<template>
  <div 
    ref="playerContainerRef" 
    class="secure-player-container" 
    @contextmenu.prevent
    @selectstart.prevent
  >
    <!-- شاشة تحذير سوداء تظهر فوراً عند كشف تسجيل الشاشة أو فقدان التركيز -->
    <div v-if="isBlocked" class="security-warning-overlay">
      <div class="warning-content">
        <div class="warning-icon">⚠️</div>
        <h3>تم إيقاف التشغيل لحماية المحتوى</h3>
        <p>تم كشف مغادرة النافذة، سحب الإشعارات، أو محاولة تسجيل الشاشة.</p>
        <p class="sub-text">يرجى العودة للتطبيق وإغلاق أي أدوات خارجية للمتابعة.</p>
        <button class="resume-btn" @click="attemptResume">إعادة التشغيل</button>
      </div>
    </div>

    <!-- العلامة المائية الديناميكية المتحركة عشوائياً -->
    <div
      ref="watermarkRef"
      class="dynamic-watermark"
      :style="watermarkStyle"
    >
      <span class="wm-name">{{ studentName }}</span>
      <span class="wm-separator">|</span>
      <span class="wm-id">ID: {{ telegramId }}</span>
    </div>

    <!-- مشغل الفيديو المحمي -->
    <video
      ref="videoRef"
      class="secure-video-element"
      controls
      disablepictureinpicture
      controlslist="nodownload nodownload noremoteplayback"
      playsinline
      :src="videoStreamUrl"
      @play="onPlay"
      @pause="onPause"
      @ended="onEnded"
    >
      متصفحك لا يدعم مشغل الفيديو المحمي.
    </video>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount, watch } from 'vue'

// التعريف بالخصائص المدخلة للمكون
const props = defineProps<{
  videoStreamUrl: string
  studentName: string
  telegramId: string | number
}>()

// المراجع والمؤشرات
const playerContainerRef = ref<HTMLDivElement | null>(null)
const videoRef = ref<HTMLVideoElement | null>(null)
const watermarkRef = ref<HTMLDivElement | null>(null)

// حالة الأمان والمنع
const isBlocked = ref(false)
let watermarkInterval: ReturnType<typeof setInterval> | null = null
let mutationObserver: MutationObserver | null = null

// موضع العلامة المائية الديناميكية
const watermarkStyle = reactive({
  top: '20%',
  left: '20%',
  opacity: '0.45',
  transform: 'rotate(-5deg)'
})

// ------------------------------------------------------------------
// 1. نظام حماية فقدان التركيز وتسجيل الشاشة (Blur & Visibility)
// ------------------------------------------------------------------
const triggerSecurityBlock = (reason: string) => {
  console.warn(`[Security Alert] Video blocked: ${reason}`)
  isBlocked.value = true
  if (videoRef.value && !videoRef.value.paused) {
    videoRef.value.pause()
  }
}

const handleVisibilityChange = () => {
  if (document.hidden || document.visibilityState === 'hidden') {
    triggerSecurityBlock('Document hidden or switching apps')
  }
}

const handleWindowBlur = () => {
  triggerSecurityBlock('Window lost focus or notifications shade pulled')
}

const handleKeyDown = (e: KeyboardEvent) => {
  // منع اختصارات التقاط الشاشة والأدوات الخارجية (PrintScreen, F12, Ctrl+Shift+I, etc.)
  if (
    e.key === 'PrintScreen' ||
    e.code === 'PrintScreen' ||
    (e.ctrlKey && e.shiftKey && (e.key === 'I' || e.key === 'i' || e.key === 'C' || e.key === 'c')) ||
    (e.ctrlKey && (e.key === 'u' || e.key === 'U')) ||
    e.key === 'F12'
  ) {
    e.preventDefault()
    triggerSecurityBlock('Forbidden key shortcut detected')
  }
}

const attemptResume = () => {
  if (!document.hidden) {
    isBlocked.value = false
  } else {
    alert('يرجى العودة إلى نافذة التطبيق أولاً.')
  }
}

// ------------------------------------------------------------------
// 2. تحريك العلامة المائية عشوائياً (Dynamic Watermark Movement)
// ------------------------------------------------------------------
const moveWatermarkRandomly = () => {
  // توليد مواقع عشوائية بين 10% و 80% لمنع خروجها خارج نطاق الشاشة
  const randomTop = Math.floor(Math.random() * 70 + 10)
  const randomLeft = Math.floor(Math.random() * 65 + 10)
  const randomRotate = Math.floor(Math.random() * 20 - 10) // بين -10 و 10 درجات
  const randomOpacity = (Math.random() * 0.25 + 0.35).toFixed(2) // شفافية ديناميكية بين 0.35 و 0.60

  watermarkStyle.top = `${randomTop}%`
  watermarkStyle.left = `${randomLeft}%`
  watermarkStyle.opacity = randomOpacity
  watermarkStyle.transform = `rotate(${randomRotate}deg)`
}

// ------------------------------------------------------------------
// 3. كاشف التلاعب بـ DOM (DOM Anti-Tamper via MutationObserver)
// ------------------------------------------------------------------
const setupAntiTamperObserver = () => {
  if (!playerContainerRef.value || !watermarkRef.value) return

  mutationObserver = new MutationObserver((mutations) => {
    for (const mutation of mutations) {
      // 1. كشف حذف عنصر العلامة المائية من الصفحة
      if (mutation.type === 'childList') {
        const removed = Array.from(mutation.removedNodes)
        if (removed.includes(watermarkRef.value as Node)) {
          triggerSecurityBlock('Watermark element removed from DOM')
          // إعادة إضافة العلامة المائية فوراً
          playerContainerRef.value?.appendChild(watermarkRef.value!)
        }
      }
      // 2. كشف التلاعب بأنماط CSS إخفاء العلامة المائية (display: none / opacity: 0)
      if (mutation.type === 'attributes' && mutation.target === watermarkRef.value) {
        const el = watermarkRef.value
        const computed = window.getComputedStyle(el)
        if (computed.display === 'none' || computed.visibility === 'hidden' || parseFloat(computed.opacity) < 0.1) {
          triggerSecurityBlock('Watermark element styles tampered')
          el.style.display = 'block'
          el.style.visibility = 'visible'
          el.style.opacity = '0.45'
        }
      }
    }
  })

  mutationObserver.observe(playerContainerRef.value, {
    childList: true,
    subtree: true,
    attributes: true,
    attributeFilter: ['style', 'class', 'hidden']
  })
}

// أحداث المشغل
const onPlay = () => {
  if (isBlocked.value) {
    videoRef.value?.pause()
  }
}
const onPause = () => {}
const onEnded = () => {}

// ------------------------------------------------------------------
// 4. دورة الحياة (Lifecycle Hooks)
// ------------------------------------------------------------------
onMounted(() => {
  // إضافة مستمعي الأحداث للحماية
  window.addEventListener('blur', handleWindowBlur)
  document.addEventListener('visibilitychange', handleVisibilityChange)
  window.addEventListener('keydown', handleKeyDown)

  // بدء تحريك العلامة المائية كل 4 ثوانٍ
  moveWatermarkRandomly()
  watermarkInterval = setInterval(moveWatermarkRandomly, 4000)

  // تفعيل مراقب التلاعب
  setupAntiTamperObserver()
})

onBeforeUnmount(() => {
  // تنظيف مستمعي الأحداث والمؤقتات
  window.removeEventListener('blur', handleWindowBlur)
  document.removeEventListener('visibilitychange', handleVisibilityChange)
  window.removeEventListener('keydown', handleKeyDown)

  if (watermarkInterval) clearInterval(watermarkInterval)
  if (mutationObserver) mutationObserver.disconnect()
})
</script>

<style scoped>
/* حاوية مشغل الفيديو المحمي */
.secure-player-container {
  position: relative;
  width: 100%;
  max-width: 900px;
  margin: 0 auto;
  aspect-ratio: 16 / 9;
  background-color: #000000;
  border-radius: 12px;
  overflow: hidden;
  user-select: none;
  -webkit-user-select: none;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
}

/* عنصر الفيديو */
.secure-video-element {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
}

/* الشاشة التحذيرية السوداء عند حظر المحتوى */
.security-warning-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: #050505;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 24px;
  box-sizing: border-box;
}

.warning-content {
  color: #ffffff;
  max-width: 420px;
}

.warning-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.warning-content h3 {
  font-size: 20px;
  font-weight: 700;
  color: #ff4d4d;
  margin-bottom: 8px;
}

.warning-content p {
  font-size: 14px;
  color: #cccccc;
  line-height: 1.5;
  margin-bottom: 6px;
}

.warning-content .sub-text {
  font-size: 12px;
  color: #888888;
  margin-bottom: 20px;
}

.resume-btn {
  background: linear-gradient(135deg, #0088cc, #005588);
  color: #ffffff;
  border: none;
  padding: 10px 24px;
  font-size: 14px;
  font-weight: 600;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s ease, opacity 0.2s ease;
}

.resume-btn:hover {
  transform: scale(1.03);
  opacity: 0.95;
}

/* العلامة المائية الديناميكية */
.dynamic-watermark {
  position: absolute;
  z-index: 50;
  pointer-events: none; /* تضمن عدم إعاقة النقر على المشغل */
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(0, 0, 0, 0.45);
  color: rgba(255, 255, 255, 0.9);
  padding: 6px 14px;
  border-radius: 20px;
  font-family: monospace, sans-serif;
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.5px;
  border: 1px solid rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(2px);
  transition: top 0.8s ease-in-out, left 0.8s ease-in-out, transform 0.8s ease-in-out, opacity 0.8s ease-in-out;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.wm-name {
  color: #64ffda;
}

.wm-separator {
  color: #888888;
}

.wm-id {
  color: #ffca28;
}
</style>
