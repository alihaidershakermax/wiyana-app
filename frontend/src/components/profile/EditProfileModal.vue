<script setup lang="ts">
import { ref, watch } from "vue";
import { X, User, Sparkles, Check, Camera, Send, Hash } from "@lucide/vue";
import { UserProfile } from "@/types/user";

const props = defineProps<{
  isOpen: boolean;
  user?: UserProfile | null;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "save", updated: Partial<UserProfile>): void;
}>();

const form = ref({
  name: props.user?.name || "",
  telegramUsername: props.user?.telegramUsername || "",
  avatarUrl: props.user?.avatarUrl || "https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&q=80&w=250",
});

watch(
  () => props.user,
  (newUser) => {
    if (newUser) {
      form.value = {
        name: newUser.name || "",
        telegramUsername: newUser.telegramUsername || "",
        avatarUrl: newUser.avatarUrl || form.value.avatarUrl,
      };
    }
  },
  { deep: true, immediate: true }
);

const avatarPresets = [
  "https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&q=80&w=250",
  "https://images.unsplash.com/photo-1539571696357-5a69c17a67c6?auto=format&fit=crop&q=80&w=250",
  "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&q=80&w=250",
  "https://images.unsplash.com/photo-1517841905240-472988babdf9?auto=format&fit=crop&q=80&w=250",
  "https://images.unsplash.com/photo-1524504388940-b1c1722653e1?auto=format&fit=crop&q=80&w=250",
];

function handleSave() {
  if (!form.value.name.trim()) return;
  emit("save", {
    name: form.value.name.trim(),
    telegramUsername: form.value.telegramUsername.replace("@", "").trim(),
    avatarUrl: form.value.avatarUrl,
  });
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
      class="fixed inset-0 z-50 bg-slate-950/80 backdrop-blur-md flex items-center justify-center p-4 overflow-y-auto"
      @click.self="emit('close')"
    >
      <div class="w-full max-w-lg bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 shadow-2xl flex flex-col gap-5 my-8">
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-800/80 pb-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center">
              <Send class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-lg font-black text-slate-900 dark:text-slate-100">
                بيانات حساب تيليجرام
              </h3>
              <p class="text-xs text-slate-500 dark:text-slate-400">
                تعديل الاسم والصورة المعروضة ومعرف حساب تيليجرام
              </p>
            </div>
          </div>

          <button
            @click="emit('close')"
            type="button"
            class="p-2 rounded-xl text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
          >
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleSave" class="flex flex-col gap-4">
          <!-- Avatar Preview & Preset Selection -->
          <div class="flex flex-col items-center gap-3 py-2 bg-slate-50 dark:bg-slate-950/60 rounded-2xl p-4 border border-slate-100 dark:border-slate-800/60">
            <div class="relative group cursor-pointer">
              <img
                :src="form.avatarUrl"
                alt="Avatar"
                class="w-20 h-20 rounded-full object-cover ring-4 ring-orange-500/30 shadow-md"
              />
              <div class="absolute inset-0 rounded-full bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                <Camera class="w-6 h-6 text-white" />
              </div>
            </div>

            <span class="text-xs font-bold text-slate-600 dark:text-slate-300">اختر صورة الحساب:</span>
            <div class="flex items-center gap-2">
              <button
                v-for="(preset, idx) in avatarPresets"
                :key="idx"
                type="button"
                @click="form.avatarUrl = preset"
                :class="[
                  'w-9 h-9 rounded-full overflow-hidden transition-all duration-200 border-2 cursor-pointer',
                  form.avatarUrl === preset ? 'border-orange-500 scale-110 shadow-sm shadow-orange-500/30' : 'border-transparent opacity-70 hover:opacity-100'
                ]"
              >
                <img :src="preset" class="w-full h-full object-cover" />
              </button>
            </div>
          </div>

          <!-- Full Name -->
          <div class="flex flex-col gap-1.5">
            <label class="text-xs font-bold text-slate-700 dark:text-slate-300 flex items-center gap-1.5">
              <User class="w-3.5 h-3.5 text-orange-500" />
              الاسم الكامل
            </label>
            <input
              v-model="form.name"
              type="text"
              required
              placeholder="اكتب اسمك..."
              class="w-full bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 px-4 py-3 rounded-2xl text-xs sm:text-sm font-semibold border border-slate-200 dark:border-slate-800 focus:border-orange-500 focus:outline-none transition-colors"
            />
          </div>

          <!-- Telegram Username -->
          <div class="flex flex-col gap-1.5">
            <label class="text-xs font-bold text-slate-700 dark:text-slate-300 flex items-center gap-1.5">
              <Send class="w-3.5 h-3.5 text-orange-500" />
              معرف حساب تيليجرام (Telegram Username)
            </label>
            <div class="relative">
              <span class="absolute right-3.5 top-1/2 -translate-y-1/2 text-xs font-bold text-slate-400">@</span>
              <input
                v-model="form.telegramUsername"
                type="text"
                placeholder="username"
                class="w-full bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 pr-8 pl-3 py-3 rounded-2xl text-xs font-semibold border border-slate-200 dark:border-slate-800 focus:border-orange-500 focus:outline-none transition-colors"
              />
            </div>
          </div>

          <!-- Telegram ID info box -->
          <div class="p-3.5 rounded-2xl bg-orange-500/10 border border-orange-500/20 flex items-center justify-between">
            <div class="flex items-center gap-2 text-xs font-bold text-slate-800 dark:text-slate-200">
              <Hash class="w-4 h-4 text-orange-500" />
              <span>معرف تيليجرام ID:</span>
              <span class="font-mono text-orange-600 dark:text-orange-400">{{ user?.telegramId || '—' }}</span>
            </div>
            <span class="text-[10px] font-black text-orange-600 dark:text-orange-400 bg-orange-500/10 px-2 py-0.5 rounded-full">
              متصل تلقائياً
            </span>
          </div>

          <!-- Buttons -->
          <div class="flex items-center gap-3 mt-2 pt-3 border-t border-slate-100 dark:border-slate-800/80">
            <button
              type="button"
              @click="emit('close')"
              class="flex-1 py-3 rounded-2xl border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-bold text-xs sm:text-sm hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
            >
              إلغاء
            </button>
            <button
              type="submit"
              class="flex-1 py-3 rounded-2xl bg-orange-500 hover:bg-orange-600 text-white font-bold text-xs sm:text-sm shadow-md shadow-orange-500/20 transition-all flex items-center justify-center gap-1.5 active:scale-95 cursor-pointer"
            >
              <Check class="w-4 h-4" />
              حفظ التعديلات
            </button>
          </div>
        </form>
      </div>
    </div>
  </Transition>
</template>
