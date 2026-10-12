<template>
  <div
    ref="watermarkContainer"
    class="pointer-events-none absolute inset-0 z-40 overflow-hidden select-none"
    aria-hidden="true"
  >
    <!-- Primary moving watermark layer -->
    <div
      ref="primaryLayer"
      data-wiyana-wm="primary"
      class="absolute transition-all duration-1000 ease-in-out px-3 py-1.5 rounded-lg bg-black/40 backdrop-blur-xs border border-amber-500/30 text-amber-200 text-xs font-mono shadow-lg flex flex-col gap-0.5 items-center justify-center text-center"
      :style="{
        top: `${pos.topPercent}%`,
        left: `${pos.leftPercent}%`,
        opacity: pos.opacity,
        transform: `rotate(${pos.rotationDeg}deg)`,
      }"
    >
      <div class="font-bold text-[11px] tracking-wider text-amber-400">
        {{ config.displayName || 'Weyana Student' }}
      </div>
      <div class="text-[10px] text-amber-200/90 font-mono">
        {{ config.privacyId || 'ID: ****' }} | {{ sessionLabel }}
      </div>
      <div class="text-[9px] text-amber-300/70 font-mono">
        {{ timestampStr }} • Weyana
      </div>
    </div>

    <!-- Secondary background low-opacity repeating watermark layer -->
    <div
      ref="secondaryLayer"
      class="absolute inset-0 grid grid-cols-2 sm:grid-cols-3 gap-8 p-6 opacity-[0.12] pointer-events-none select-none text-amber-500/50 font-mono text-[10px]"
    >
      <div v-for="i in 6" :key="i" class="flex flex-col items-center justify-center rotate-[-15deg]">
        <span>WEYANA PROTECTED</span>
        <span>{{ config.displayName }}</span>
        <span>{{ config.privacyId }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import type { WatermarkConfig, WatermarkPosition } from '@/security/watermark'

const props = defineProps<{
  config: WatermarkConfig
}>()

const emit = defineEmits<{
  (e: 'tampered', reason?: string): void
}>()

const watermarkContainer = ref<HTMLDivElement | null>(null)
const primaryLayer = ref<HTMLDivElement | null>(null)
const secondaryLayer = ref<HTMLDivElement | null>(null)

/** Deterministic PRNG seeded from session + name so placement is stable per
 *  session yet randomized across sessions. */
function seededRandom(seed: string): () => number {
  let h = 1779033703 ^ seed.length
  for (let i = 0; i < seed.length; i++) {
    h = Math.imul(h ^ seed.charCodeAt(i), 3432918353)
    h = (h << 13) | (h >>> 19)
  }
  let a = h >>> 0
  return () => {
    a |= 0
    a = (a + 0x6d2b79f5) | 0
    let t = Math.imul(a ^ (a >>> 15), 1 | a)
    t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296
  }
}

const rng = seededRandom(`${props.config.sessionId}|${props.config.displayName}`)

/** Deterministic-but-random placement within the safe viewport area. */
function nextPosition(prev?: WatermarkPosition): WatermarkPosition {
  const r = rng()
  const top = Math.floor(r * 65) + 10
  const left = Math.floor(rng() * 60) + 10
  const opacity = Math.round((rng() * 0.3 + 0.4) * 100) / 100 // 0.40 - 0.70
  const rotationDeg = Math.floor(rng() * 12) - 6 // -6deg to +6deg
  return { topPercent: top, leftPercent: left, opacity, rotationDeg }
}

const pos = ref<WatermarkPosition>(nextPosition())
const sessionLabel = (props.config.sessionId || 'SESS').slice(0, 8).toUpperCase()
const timestampStr = ref(new Date().toLocaleTimeString())

let moveTimer: ReturnType<typeof setInterval> | null = null
let clockTimer: ReturnType<typeof setInterval> | null = null
let integrityTimer: ReturnType<typeof setInterval> | null = null
let observer: MutationObserver | null = null
let tamperFired = false

const fireTamper = (reason: string) => {
  if (tamperFired) return
  tamperFired = true
  emit('tampered', reason)
}

const checkLayerHidden = (layer: HTMLElement | null): boolean => {
  if (!layer) return false
  const style = window.getComputedStyle(layer)
  return (
    style.display === 'none' ||
    style.visibility === 'hidden' ||
    parseFloat(style.opacity || '1') <= 0.04
  )
}

const startMovement = () => {
  moveTimer = setInterval(() => {
    const next = nextPosition(pos.value)
    pos.value = {
      ...next,
      // Keep the previous opacity stable between border sweeps so the
      // watermark stays visible across reprobes.
      opacity: Math.max(0.4, next.opacity),
    }
  }, 3500)

  clockTimer = setInterval(() => {
    timestampStr.value = new Date().toLocaleTimeString()
  }, 1000)

  // Integrity probe — catches style/attribute tampering that MutationObserver
  // misses (e.g. CSS overriding opacity/display without DOM attribute change).
  integrityTimer = setInterval(() => {
    if (checkLayerHidden(primaryLayer.value) || checkLayerHidden(secondaryLayer.value)) {
      fireTamper('watermark-hidden')
    }
    if (!watermarkContainer.value?.isConnected) {
      fireTamper('watermark-removed')
    }
  }, 2500)
}

const setupAntiTamperingObserver = () => {
  if (!watermarkContainer.value) return
  observer = new MutationObserver((mutations) => {
    for (const mutation of mutations) {
      if (mutation.type === 'childList') {
        const removed = Array.from(mutation.removedNodes)
        if (
          removed.includes(primaryLayer.value as Node) ||
          removed.includes(secondaryLayer.value as Node)
        ) {
          fireTamper('watermark-node-removed')
          return
        }
      }
      if (mutation.type === 'attributes') {
        if (checkLayerHidden(mutation.target as HTMLElement)) {
          fireTamper('watermark-attr-tampered')
          return
        }
      }
    }
  })

  observer.observe(watermarkContainer.value, {
    childList: true,
    subtree: true,
    attributes: true,
    attributeFilter: ['style', 'class', 'hidden'],
  })
}

onMounted(() => {
  startMovement()
  setupAntiTamperingObserver()
})

onBeforeUnmount(() => {
  if (moveTimer) clearInterval(moveTimer)
  if (clockTimer) clearInterval(clockTimer)
  if (integrityTimer) clearInterval(integrityTimer)
  if (observer) observer.disconnect()
})
</script>