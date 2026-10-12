<template>
  <div
    class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-mono border"
    :class="statusClasses"
  >
    <span class="w-2 h-2 rounded-full animate-pulse" :class="dotClasses"></span>
    <span class="font-medium uppercase">{{ level }} ({{ score }})</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { RiskLevel } from '@/security/riskEngine'

const props = defineProps<{
  score: number
  level: RiskLevel
}>()

const statusClasses = computed(() => {
  if (props.level === 'CRITICAL') {
    return 'bg-red-600/20 text-red-300 border-red-600/40'
  }
  if (props.level === 'HIGH_RISK') {
    return 'bg-red-500/10 text-red-400 border-red-500/30'
  }
  if (props.level === 'SUSPICIOUS') {
    return 'bg-amber-500/10 text-amber-400 border-amber-500/30'
  }
  return 'bg-emerald-500/10 text-emerald-400 border-emerald-500/30'
})

const dotClasses = computed(() => {
  if (props.level === 'CRITICAL') return 'bg-red-500'
  if (props.level === 'HIGH_RISK') return 'bg-red-500'
  if (props.level === 'SUSPICIOUS') return 'bg-amber-500'
  return 'bg-emerald-500'
})
</script>