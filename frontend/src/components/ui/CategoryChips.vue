<script setup lang="ts">
import { computed } from "vue";
import { CategoryFilterId } from "@/types/course";
import { Layers, Sparkles, Flame, Gift, Atom, Beaker, Calculator, Dna, Languages } from "@lucide/vue";

const props = defineProps<{
  activeCategory: CategoryFilterId | string;
}>();

const emit = defineEmits<{
  (e: "selectCategory", category: CategoryFilterId | string): void;
}>();

const categories = computed(() => [
  { id: "all", label: "الكل", icon: Layers },
  { id: "popular", label: "الأكثر شعبية", icon: Flame },
  { id: "phys", label: "الفيزياء", icon: Atom },
  { id: "chem", label: "الكيمياء", icon: Beaker },
  { id: "math", label: "الرياضيات", icon: Calculator },
  { id: "bio", label: "الأحياء", icon: Dna },
  { id: "eng", label: "الإنكليزية", icon: Languages },
  { id: "free", label: "الكورسات المجانية", icon: Gift },
  { id: "new", label: "الحديث والمضاف مؤخراً", icon: Sparkles },
]);
</script>

<template>
  <div class="w-full flex flex-col gap-2">
    <!-- Header Section Label -->
    <div class="flex items-center justify-between px-1">
      <span class="text-xs font-bold text-[#6B7280] dark:text-[#A1A1AA] uppercase tracking-wider">
        تصفح حسب التصنيف
      </span>
      <span class="text-[10px] font-bold text-orange-500">
        {{ categories.length }} تصنيفات
      </span>
    </div>

    <!-- Scrollable Touch Pills Carousel -->
    <div class="w-full overflow-x-auto hide-scrollbar py-1.5 px-0.5 -mx-0.5 touch-pan-x">
      <div class="flex items-center gap-2.5 min-w-max">
        <button
          v-for="cat in categories"
          :key="cat.id"
          @click="emit('selectCategory', cat.id)"
          type="button"
          :class="[
            'relative inline-flex items-center gap-2.5 px-4 h-11 rounded-2xl text-xs font-bold transition-all duration-200 active:scale-95 shrink-0 focus:outline-none select-none cursor-pointer',
            props.activeCategory === cat.id
              ? 'bg-orange-500 text-white shadow-sm shadow-orange-500/20'
              : 'bg-white dark:bg-[#181818] text-[#111111] dark:text-[#FFFFFF] hover:bg-[#F8F8F8] dark:hover:bg-[#222222] border border-[#E5E7EB] dark:border-[#27272A]'
          ]"
        >
          <component
            :is="cat.icon"
            :class="[
              'w-4 h-4 transition-transform duration-200',
              props.activeCategory === cat.id ? 'scale-110' : 'text-[#6B7280] dark:text-[#A1A1AA]'
            ]"
          />
          <span>{{ cat.label }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.hide-scrollbar::-webkit-scrollbar {
  display: none;
}
.hide-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
