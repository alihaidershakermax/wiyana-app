<script setup lang="ts">
import { useNetwork } from "@vueuse/core";
import { WifiOff, RotateCw } from "@lucide/vue";

const { isOnline } = useNetwork();

const emit = defineEmits<{
  (e: "refresh"): void;
}>();
</script>

<template>
  <div v-if="!isOnline" class="w-full bg-amber-500 text-slate-950 px-4 py-2.5 shadow-md flex items-center justify-between text-xs sm:text-sm font-bold sticky top-0 z-50 animate-slide-down">
    <div class="flex items-center gap-2">
      <WifiOff class="w-4 h-4 shrink-0" />
      <span>أنت غير متصل بالإنترنت حالياً. يتم عرض البيانات من التخزين المؤقت.</span>
    </div>
    <button
      @click="emit('refresh')"
      type="button"
      class="inline-flex items-center gap-1 bg-slate-950/10 hover:bg-slate-950/20 px-2.5 py-1 rounded-lg transition-colors text-slate-950"
    >
      <RotateCw class="w-3.5 h-3.5" />
      <span>تحديث</span>
    </button>
  </div>
</template>
