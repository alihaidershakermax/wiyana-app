<script setup lang="ts">
import { Search, X } from "@lucide/vue";

defineProps<{
  modelValue: string;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
  (e: "clear"): void;
}>();

function onInput(e: Event) {
  const target = e.target as HTMLInputElement;
  emit("update:modelValue", target.value);
}
</script>

<template>
  <div class="w-full bg-white dark:bg-[#181818] border-b border-[#E5E7EB] dark:border-[#27272A] px-4 py-3 shadow-xs">
    <div class="max-w-7xl mx-auto relative flex items-center">
      <Search class="w-5 h-5 text-orange-500 absolute right-3 pointer-events-none shrink-0" />
      
      <input
        :value="modelValue"
        @input="onInput"
        type="text"
        placeholder="ابحث باسم المساق، المرحلة، أو الأستاذ..."
        class="w-full bg-[#F8F8F8] dark:bg-[#0F0F0F] text-[#111111] dark:text-[#FFFFFF] pr-10 pl-10 py-2.5 rounded-xl text-xs sm:text-sm font-semibold border border-transparent focus:border-orange-500 focus:bg-white dark:focus:bg-[#0F0F0F] focus:outline-none transition-all placeholder:text-[#9CA3AF] dark:placeholder:text-[#71717A]"
      />

      <button
        v-if="modelValue.length > 0"
        @click="emit('clear')"
        type="button"
        aria-label="مسح البحث"
        class="absolute left-3 text-[#6B7280] hover:text-[#111111] dark:text-[#A1A1AA] dark:hover:text-white p-1 rounded-lg transition-colors cursor-pointer"
      >
        <X class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>
