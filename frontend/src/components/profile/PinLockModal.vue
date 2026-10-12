<script setup lang="ts">
import { ref } from "vue";
import { X, KeyRound, Check } from "@lucide/vue";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "saved", pin: string): void;
}>();

const pin = ref("");
const isSaved = ref(false);

function handleSavePin() {
  if (pin.value.length < 4) return;
  isSaved.value = true;
  localStorage.setItem("wy_app_pin", pin.value);
  setTimeout(() => {
    isSaved.value = false;
    emit("saved", pin.value);
    emit("close");
  }, 1200);
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
      <div class="w-full max-w-sm bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 shadow-2xl flex flex-col gap-4 my-8 items-center text-center">
        <!-- Header -->
        <div class="w-12 h-12 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center">
          <KeyRound class="w-6 h-6" />
        </div>

        <div>
          <h3 class="text-lg font-black text-slate-900 dark:text-slate-100">
            قفل التطبيق برمز PIN
          </h3>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
            عيّن رمزاً سرياً مكوناً من 4 أرقام لفتح المنصة
          </p>
        </div>

        <!-- Input -->
        <form @submit.prevent="handleSavePin" class="w-full flex flex-col gap-4">
          <input
            v-model="pin"
            type="password"
            maxlength="4"
            pattern="[0-9]{4}"
            inputmode="numeric"
            placeholder="••••"
            class="w-full bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 text-center tracking-[1em] text-2xl font-mono py-3 rounded-2xl border border-slate-200 dark:border-slate-800 focus:border-orange-500 focus:outline-none"
          />

          <p v-if="isSaved" class="text-xs font-bold text-orange-500 flex items-center justify-center gap-1.5">
            <Check class="w-4 h-4" />
            تم حفظ رمز القفل بنجاح!
          </p>

          <div class="flex items-center gap-3">
            <button
              type="button"
              @click="emit('close')"
              class="flex-1 py-3 rounded-2xl border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-bold text-xs sm:text-sm hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
            >
              إلغاء
            </button>
            <button
              type="submit"
              :disabled="pin.length < 4"
              class="flex-1 py-3 rounded-2xl bg-orange-500 hover:bg-orange-600 text-white font-bold text-xs sm:text-sm shadow-md shadow-orange-500/20 transition-all text-center active:scale-95 disabled:opacity-50 cursor-pointer"
            >
              حفظ الرمز
            </button>
          </div>
        </form>
      </div>
    </div>
  </Transition>
</template>
