<script setup lang="ts">
import { ref } from "vue";
import { X, Smartphone, Laptop, Tablet, ShieldCheck, LogOut, Check } from "@lucide/vue";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "revoked"): void;
}>();

const devices = ref([
  {
    id: "dev_1",
    name: "هذا الجهاز الحالي (iPhone 15 Pro)",
    type: "mobile",
    ip: "151.236.14.88 (بغداد، العراق)",
    lastActive: "نشط الآن",
    isCurrent: true,
  },
  {
    id: "dev_2",
    name: "Chrome Browser (Windows 11 PC)",
    type: "laptop",
    ip: "37.237.210.12 (العراق)",
    lastActive: "منذ 4 ساعات",
    isCurrent: false,
  },
  {
    id: "dev_3",
    name: "iPad Air (5th Generation)",
    type: "tablet",
    ip: "151.236.14.88 (العراق)",
    lastActive: "منذ يومين",
    isCurrent: false,
  },
]);

function revokeOtherDevices() {
  devices.value = devices.value.filter(d => d.isCurrent);
  emit("revoked");
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
      <div class="w-full max-w-lg bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 shadow-2xl flex flex-col gap-4 my-8">
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-800/80 pb-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center">
              <Smartphone class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-lg font-black text-slate-900 dark:text-slate-100">
                الأجهزة المسجلة والجلسات
              </h3>
              <p class="text-xs text-slate-500 dark:text-slate-400">
                إدارة الأجهزة المسموح لها بتسجيل الدخول إلى حسابك
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

        <!-- Devices List -->
        <div class="flex flex-col gap-3 max-h-80 overflow-y-auto pr-1">
          <div
            v-for="dev in devices"
            :key="dev.id"
            class="p-4 rounded-2xl bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800/60 flex items-center justify-between gap-3"
          >
            <div class="flex items-center gap-3 min-w-0">
              <div class="w-9 h-9 rounded-xl bg-orange-500/10 text-orange-500 flex items-center justify-center shrink-0">
                <Laptop v-if="dev.type === 'laptop'" class="w-4 h-4" />
                <Tablet v-else-if="dev.type === 'tablet'" class="w-4 h-4" />
                <Smartphone v-else class="w-4 h-4" />
              </div>
              <div class="flex flex-col min-w-0">
                <div class="flex items-center gap-2">
                  <h4 class="text-xs sm:text-sm font-bold text-slate-900 dark:text-slate-100 truncate">
                    {{ dev.name }}
                  </h4>
                    <span
                      v-if="dev.isCurrent"
                      class="px-2 py-0.5 rounded-full bg-orange-500/10 text-orange-500 text-[9px] font-black border border-orange-500/20 shrink-0"
                    >
                      هذا الجهاز
                    </span>
                </div>
                <p class="text-[11px] text-slate-500 dark:text-slate-400 truncate">
                  {{ dev.ip }} • {{ dev.lastActive }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-col sm:flex-row items-center gap-3 pt-2 border-t border-slate-100 dark:border-slate-800/80">
          <button
            v-if="devices.length > 1"
            type="button"
            @click="revokeOtherDevices"
            class="w-full sm:w-auto flex-1 py-3 px-4 rounded-2xl border border-rose-200 dark:border-rose-900/50 text-rose-600 dark:text-rose-400 font-bold text-xs hover:bg-rose-50 dark:hover:bg-rose-950/30 transition-colors flex items-center justify-center gap-1.5 cursor-pointer"
          >
            <LogOut class="w-4 h-4" />
            تسجيل الخروج من بقية الأجهزة
          </button>
          <button
            type="button"
            @click="emit('close')"
            class="w-full sm:w-auto flex-1 py-3 rounded-2xl bg-orange-500 hover:bg-orange-600 text-white font-bold text-xs sm:text-sm shadow-md shadow-orange-500/20 transition-all text-center active:scale-95 cursor-pointer"
          >
            تم
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>
