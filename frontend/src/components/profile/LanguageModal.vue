<script setup lang="ts">
import { X, Globe, Check } from "@lucide/vue";
import { useUserStore } from "@/stores/useUserStore";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const userStore = useUserStore();

function selectLanguage(lang: "ar" | "en") {
  userStore.setLanguage(lang);
  emit("close");
}
</script>

<template>
  <Transition
    enter-active-class="transition duration-200 ease-out"
    enter-from-class="opacity-0 scale-95"
    enter-to-class="opacity-100 scale-100"
    leave-active-class="transition duration-150 ease-in"
    leave-from-class="opacity-100 scale-100"
    leave-to-class="opacity-0 scale-95"
  >
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 bg-black/70 backdrop-blur-md flex items-center justify-center p-4 overflow-y-auto"
      @click.self="emit('close')"
    >
      <div
        class="w-full max-w-sm bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-3xl p-6 shadow-2xl flex flex-col gap-4 text-start"
      >
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-[#E5E7EB] dark:border-[#27272A] pb-3.5">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center">
              <Globe class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-base font-black text-[#111111] dark:text-[#FFFFFF]">
                لغة التطبيق
              </h3>
              <p class="text-[11px] text-[#6B7280] dark:text-[#A1A1AA]">
                اختر لغة عرض الواجهة
              </p>
            </div>
          </div>

          <button
            @click="emit('close')"
            type="button"
            class="p-2 rounded-xl text-[#9CA3AF] hover:text-[#111111] dark:hover:text-white hover:bg-[#F8F8F8] dark:hover:bg-[#222222] transition-colors cursor-pointer"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Languages List -->
        <div class="flex flex-col gap-2">
          <!-- Arabic -->
          <button
            type="button"
            @click="selectLanguage('ar')"
            class="w-full flex items-center justify-between p-3.5 rounded-2xl border transition-all cursor-pointer text-start"
            :class="[
              userStore.settings.language === 'ar'
                ? 'border-orange-500 bg-orange-500/5 text-[#111111] dark:text-white font-bold'
                : 'border-[#E5E7EB] dark:border-[#27272A] hover:bg-[#F8F8F8] dark:hover:bg-[#222222] text-[#111111] dark:text-white'
            ]"
          >
            <div class="flex flex-col">
              <span class="text-sm font-bold">العربية</span>
              <span class="text-[11px] text-[#6B7280] dark:text-[#A1A1AA] font-normal">اللغة الافتراضية للواجهة</span>
            </div>
            <div
              v-if="userStore.settings.language === 'ar'"
              class="w-6 h-6 rounded-full bg-orange-500 text-white flex items-center justify-center shadow-xs"
            >
              <Check class="w-3.5 h-3.5 stroke-[2.5]" />
            </div>
          </button>

          <!-- English -->
          <button
            type="button"
            @click="selectLanguage('en')"
            class="w-full flex items-center justify-between p-3.5 rounded-2xl border transition-all cursor-pointer text-start"
            :class="[
              userStore.settings.language === 'en'
                ? 'border-orange-500 bg-orange-500/5 text-[#111111] dark:text-white font-bold'
                : 'border-[#E5E7EB] dark:border-[#27272A] hover:bg-[#F8F8F8] dark:hover:bg-[#222222] text-[#111111] dark:text-white'
            ]"
          >
            <div class="flex flex-col">
              <span class="text-sm font-bold">English</span>
              <span class="text-[11px] text-[#6B7280] dark:text-[#A1A1AA] font-normal">English Language Interface</span>
            </div>
            <div
              v-if="userStore.settings.language === 'en'"
              class="w-6 h-6 rounded-full bg-orange-500 text-white flex items-center justify-center shadow-xs"
            >
              <Check class="w-3.5 h-3.5 stroke-[2.5]" />
            </div>
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>
