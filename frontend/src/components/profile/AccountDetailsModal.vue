<script setup lang="ts">
import { ref } from "vue";
import { X, User, Copy, Check } from "@lucide/vue";
import { UserProfile } from "@/types/user";

const props = defineProps<{
  isOpen: boolean;
  user: UserProfile | null;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const isCopied = ref(false);

function copyId() {
  if (!props.user?.telegramId) return;
  navigator.clipboard?.writeText(props.user.telegramId);
  isCopied.value = true;
  setTimeout(() => {
    isCopied.value = false;
  }, 2000);
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
              <User class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-base font-black text-[#111111] dark:text-[#FFFFFF]">
                معلومات الحساب
              </h3>
              <p class="text-[11px] text-[#6B7280] dark:text-[#A1A1AA]">
                بيانات الطالب الموثقة
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

        <!-- Account Fields -->
        <div class="flex flex-col gap-2.5">
          <!-- Name -->
          <div class="p-3 rounded-2xl bg-[#F8F8F8] dark:bg-[#222222] border border-[#E5E7EB] dark:border-[#27272A] flex items-center justify-between">
            <span class="text-xs text-[#6B7280] dark:text-[#A1A1AA] font-semibold">الاسم</span>
            <span class="text-xs font-bold text-[#111111] dark:text-[#FFFFFF]">{{ user?.name || "طالب ويانا" }}</span>
          </div>

          <!-- Telegram Username -->
          <div class="p-3 rounded-2xl bg-[#F8F8F8] dark:bg-[#222222] border border-[#E5E7EB] dark:border-[#27272A] flex items-center justify-between">
            <span class="text-xs text-[#6B7280] dark:text-[#A1A1AA] font-semibold">معرف تيليجرام</span>
            <span class="text-xs font-bold text-orange-500" dir="ltr">
              {{ user?.telegramUsername ? `@${user.telegramUsername}` : "غير متوفر" }}
            </span>
          </div>

          <!-- Telegram ID -->
          <div class="p-3 rounded-2xl bg-[#F8F8F8] dark:bg-[#222222] border border-[#E5E7EB] dark:border-[#27272A] flex items-center justify-between">
            <span class="text-xs text-[#6B7280] dark:text-[#A1A1AA] font-semibold">معرف الحساب (ID)</span>
            <div class="flex items-center gap-2">
              <span class="text-xs font-mono font-bold text-[#111111] dark:text-[#FFFFFF]" dir="ltr">
                {{ user?.telegramId || "—" }}
              </span>
              <button
                v-if="user?.telegramId"
                type="button"
                @click="copyId"
                class="p-1 rounded-lg hover:bg-[#E5E7EB] dark:hover:bg-[#333333] transition-colors cursor-pointer"
                title="نسخ"
              >
                <component :is="isCopied ? Check : Copy" :class="['w-3.5 h-3.5', isCopied ? 'text-emerald-500' : 'text-orange-500']" />
              </button>
            </div>
          </div>

          <!-- Joined Date -->
          <div v-if="user?.joinedDate" class="p-3 rounded-2xl bg-[#F8F8F8] dark:bg-[#222222] border border-[#E5E7EB] dark:border-[#27272A] flex items-center justify-between">
            <span class="text-xs text-[#6B7280] dark:text-[#A1A1AA] font-semibold">تاريخ الانضمام</span>
            <span class="text-xs font-bold text-[#111111] dark:text-[#FFFFFF]">{{ user.joinedDate }}</span>
          </div>
        </div>

        <!-- Close Button -->
        <button
          type="button"
          @click="emit('close')"
          class="w-full py-2.5 rounded-2xl bg-[#F8F8F8] dark:bg-[#222222] hover:bg-[#E5E7EB] dark:hover:bg-[#2A2A2A] text-[#111111] dark:text-[#FFFFFF] font-bold text-xs transition-all cursor-pointer text-center mt-1"
        >
          إغلاق
        </button>
      </div>
    </div>
  </Transition>
</template>
