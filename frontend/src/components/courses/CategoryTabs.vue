<script setup lang="ts">
import { computed } from "vue";
import { CourseTabCategory, CategoryOption } from "@/types/course";
import { Layers, BookMarked, Sparkles, Gift, Tag, Flame } from "@lucide/vue";

const props = defineProps<{
  activeTab: CourseTabCategory;
}>();

const emit = defineEmits<{
  (e: "selectTab", tab: CourseTabCategory): void;
}>();

const tabs = computed<Array<CategoryOption & { icon: any }>>(() => [
  { id: "all", label: "جميع المساقات", icon: Layers },
  { id: "my_courses", label: "اشتراكاتي", icon: BookMarked },
  { id: "newest", label: "المضاف حديثاً", icon: Sparkles },
  { id: "popular", label: "الأكثر مشاهدة", icon: Flame },
  { id: "free", label: "المجاني", icon: Gift },
  { id: "paid", label: "المدفوع", icon: Tag },
  { id: "favorites", label: "المفضلة", icon: BookMarked },
]);
</script>

<template>
  <div class="w-full overflow-x-auto hide-scrollbar py-1 px-0.5 touch-pan-x">
    <div class="flex items-center gap-2 min-w-max">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        @click="emit('selectTab', tab.id)"
        type="button"
        :class="[
          'inline-flex items-center gap-2 px-4 h-11 rounded-2xl text-xs font-bold transition-all duration-200 active:scale-95 shrink-0 focus:outline-none select-none cursor-pointer',
          props.activeTab === tab.id
            ? 'bg-orange-500 text-white shadow-sm shadow-orange-500/20'
            : 'bg-white dark:bg-[#181818] text-[#111111] dark:text-[#FFFFFF] hover:bg-[#F8F8F8] dark:hover:bg-[#222222] border border-[#E5E7EB] dark:border-[#27272A]'
        ]"
      >
        <component
          :is="tab.icon"
          :class="[
            'w-4 h-4 transition-transform duration-200',
            props.activeTab === tab.id ? 'scale-110' : 'text-[#6B7280] dark:text-[#A1A1AA]'
          ]"
        />
        <span>{{ tab.label }}</span>
      </button>
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
