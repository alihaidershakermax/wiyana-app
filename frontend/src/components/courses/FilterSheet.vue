<script setup lang="ts">
import { X, RotateCcw, Check, ArrowUpDown } from "@lucide/vue";
import { CoursesFilterMetadata } from "@/services/courseService";
import { CourseSortOption } from "@/types/course";

defineProps<{
  isOpen: boolean;
  metadata?: CoursesFilterMetadata;
  selectedDepartment: string;
  selectedGrade: string;
  selectedSubject?: string;
  selectedInstructor: string;
  isFreeOnly: boolean;
  isPaidOnly: boolean;
  isSubscribedOnly: boolean;
  sortBy: CourseSortOption;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "update:department", value: string): void;
  (e: "update:grade", value: string): void;
  (e: "update:subject", value: string): void;
  (e: "update:instructor", value: string): void;
  (e: "update:freeOnly", value: boolean): void;
  (e: "update:paidOnly", value: boolean): void;
  (e: "update:subscribedOnly", value: boolean): void;
  (e: "update:sortBy", value: CourseSortOption): void;
  (e: "reset"): void;
}>();

const sortOptions: Array<{ id: CourseSortOption; label: string }> = [
  { id: "newest", label: "الأحدث" },
  { id: "oldest", label: "الأقدم" },
  { id: "popular", label: "الأكثر مشاهدة" },
  { id: "top_rated", label: "الأعلى تقييماً" },
  { id: "alphabetical", label: "الأبجدية (أ-ي)" },
];
</script>

<template>
  <Transition
    enter-active-class="transition duration-300 ease-out"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition duration-200 ease-in"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-end sm:items-center justify-center p-0 sm:p-4"
      @click.self="emit('close')"
    >
      <div
        class="w-full sm:max-w-lg bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-t-3xl sm:rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[85vh] animate-slide-up"
      >
        <!-- Sheet Header -->
        <div class="p-4 border-b border-[#E5E7EB] dark:border-[#27272A] flex items-center justify-between">
          <div class="flex items-center gap-2">
            <h2 class="text-base sm:text-lg font-bold text-[#111111] dark:text-[#FFFFFF]">
              فلترة وتصفية المساقات
            </h2>
          </div>
          
          <button
            @click="emit('close')"
            type="button"
            class="p-2 rounded-xl text-[#6B7280] hover:text-[#111111] dark:text-[#A1A1AA] dark:hover:text-[#FFFFFF] hover:bg-[#F8F8F8] dark:hover:bg-[#222222] transition-colors cursor-pointer"
          >
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Sheet Body / Filter Form -->
        <div class="p-5 overflow-y-auto flex flex-col gap-5 text-right">
          <!-- Sorting Select -->
          <div class="flex flex-col gap-2">
            <label class="text-xs font-bold text-[#111111] dark:text-[#FFFFFF] flex items-center gap-1.5">
              <ArrowUpDown class="w-3.5 h-3.5 text-orange-500" />
              <span>ترتيب النتائج حسب</span>
            </label>
            <div class="grid grid-cols-2 gap-2">
              <button
                v-for="opt in sortOptions"
                :key="opt.id"
                @click="emit('update:sortBy', opt.id)"
                type="button"
                :class="[
                  'px-3 py-2 rounded-xl text-xs font-semibold border text-right transition-all flex items-center justify-between cursor-pointer',
                  sortBy === opt.id
                    ? 'bg-orange-500/10 border-orange-500 text-orange-500'
                    : 'bg-[#F8F8F8] dark:bg-[#0F0F0F] border-[#E5E7EB] dark:border-[#27272A] text-[#6B7280] dark:text-[#A1A1AA]'
                ]"
              >
                <span>{{ opt.label }}</span>
                <Check v-if="sortBy === opt.id" class="w-3.5 h-3.5 text-orange-500" />
              </button>
            </div>
          </div>

          <!-- Department Select -->
          <div v-if="metadata?.departments" class="flex flex-col gap-1.5">
            <label class="text-xs font-bold text-[#111111] dark:text-[#FFFFFF]">القسم / الفرع الدراسي</label>
            <select
              :value="selectedDepartment"
              @change="emit('update:department', ($event.target as HTMLSelectElement).value)"
              class="w-full bg-[#F8F8F8] dark:bg-[#0F0F0F] text-[#111111] dark:text-[#FFFFFF] border border-[#E5E7EB] dark:border-[#27272A] rounded-xl px-3.5 py-2.5 text-xs sm:text-sm font-semibold focus:outline-none focus:border-orange-500"
            >
              <option value="">جميع الأقسام</option>
              <option v-for="dept in metadata.departments" :key="dept.id" :value="dept.id">
                {{ dept.name }}
              </option>
            </select>
          </div>

          <!-- Grade / Level Select -->
          <div v-if="metadata?.grades" class="flex flex-col gap-1.5">
            <label class="text-xs font-bold text-[#111111] dark:text-[#FFFFFF]">المرحلة الدراسية</label>
            <select
              :value="selectedGrade"
              @change="emit('update:grade', ($event.target as HTMLSelectElement).value)"
              class="w-full bg-[#F8F8F8] dark:bg-[#0F0F0F] text-[#111111] dark:text-[#FFFFFF] border border-[#E5E7EB] dark:border-[#27272A] rounded-xl px-3.5 py-2.5 text-xs sm:text-sm font-semibold focus:outline-none focus:border-orange-500"
            >
              <option value="">جميع المراحل</option>
              <option v-for="grade in metadata.grades" :key="grade.id" :value="grade.id">
                {{ grade.name }}
              </option>
            </select>
          </div>



          <!-- Instructor Select -->
          <div v-if="metadata?.instructors" class="flex flex-col gap-1.5">
            <label class="text-xs font-bold text-[#111111] dark:text-[#FFFFFF]">المدرس</label>
            <select
              :value="selectedInstructor"
              @change="emit('update:instructor', ($event.target as HTMLSelectElement).value)"
              class="w-full bg-[#F8F8F8] dark:bg-[#0F0F0F] text-[#111111] dark:text-[#FFFFFF] border border-[#E5E7EB] dark:border-[#27272A] rounded-xl px-3.5 py-2.5 text-xs sm:text-sm font-semibold focus:outline-none focus:border-orange-500"
            >
              <option value="">جميع المدرسين</option>
              <option v-for="inst in metadata.instructors" :key="inst.id" :value="inst.id">
                {{ inst.name }}
              </option>
            </select>
          </div>
        </div>

        <!-- Sheet Footer -->
        <div class="p-4 border-t border-[#E5E7EB] dark:border-[#27272A] flex items-center justify-between gap-3 bg-[#F8F8F8] dark:bg-[#0F0F0F]">
          <button
            @click="emit('reset')"
            type="button"
            class="inline-flex items-center gap-1.5 px-4 py-2.5 rounded-xl border border-[#E5E7EB] dark:border-[#27272A] text-[#6B7280] dark:text-[#A1A1AA] font-bold text-xs hover:bg-white dark:hover:bg-[#181818] transition-colors cursor-pointer"
          >
            <RotateCcw class="w-3.5 h-3.5" />
            <span>إعادة تعيين</span>
          </button>

          <button
            @click="emit('close')"
            type="button"
            class="flex-1 px-5 py-2.5 rounded-xl bg-orange-500 hover:bg-[#EA580C] text-white font-bold text-xs sm:text-sm shadow-sm shadow-orange-500/20 transition-all text-center cursor-pointer"
          >
            عرض النتائج
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
@keyframes slideUp {
  from {
    transform: translateY(100%);
  }
  to {
    transform: translateY(0);
  }
}
.animate-slide-up {
  animation: slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}
</style>
