<script setup lang="ts">
import { ref } from "vue";
import { KeyRound, X, CheckCircle2, AlertCircle } from "@lucide/vue";
import { redeemSubscriptionCode } from "@/api/subscription";
import { useQueryClient } from "@tanstack/vue-query";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "activated"): void;
}>();

const queryClient = useQueryClient();
const code = ref("");
const isLoading = ref(false);
const errorMessage = ref("");
const successMessage = ref("");

async function handleSubmit() {
  errorMessage.value = "";
  successMessage.value = "";
  if (!code.value.trim()) {
    errorMessage.value = "يرجى إدخال رمز كود الاشتراك أولاً.";
    return;
  }

  isLoading.value = true;
  try {
    const res = await redeemSubscriptionCode(code.value.trim());
    successMessage.value = res?.message || "تم تفعيل الاشتراك بنجاح!";
    
    // Invalidate queries so Profile, Home, Subscriptions, Courses all refresh with real active status
    queryClient.invalidateQueries({ queryKey: ["userProfileData"] });
    queryClient.invalidateQueries({ queryKey: ["userSubscriptionInfo"] });
    queryClient.invalidateQueries({ queryKey: ["homeDashboardData"] });
    queryClient.invalidateQueries({ queryKey: ["mySubscriptions"] });
    queryClient.invalidateQueries({ queryKey: ["coursesList"] });

    setTimeout(() => {
      emit("activated");
      emit("close");
      code.value = "";
      successMessage.value = "";
    }, 1200);
  } catch (err: any) {
    if (err?.response?.data?.detail) {
      errorMessage.value = err.response.data.detail;
    } else if (err?.message) {
      errorMessage.value = err.message;
    } else {
      errorMessage.value = "رمز الاشتراك غير صالح أو منتهي الصلاحية.";
    }
  } finally {
    isLoading.value = false;
  }
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
        class="w-full max-w-sm sm:max-w-md bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-3xl p-6 shadow-2xl flex flex-col gap-4 my-8 text-start"
      >
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-[#E5E7EB] dark:border-[#27272A] pb-3.5">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center">
              <KeyRound class="w-5 h-5 stroke-[2.2]" />
            </div>
            <div>
              <h3 class="text-base font-black text-[#111111] dark:text-[#FFFFFF]">
                تفعيل كود الاشتراك
              </h3>
              <p class="text-[11px] text-[#6B7280] dark:text-[#A1A1AA]">
                أدخل الكود للوصول إلى المحاضرات
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

        <p class="text-xs text-[#6B7280] dark:text-[#A1A1AA] leading-relaxed">
          أدخل رمز الاشتراك المكون من أرقام أو حروف المطبوع على كارت الاشتراك الخاص بك.
        </p>

        <!-- Input Form -->
        <form @submit.prevent="handleSubmit" class="flex flex-col gap-3">
          <input
            v-model="code"
            type="text"
            placeholder="مثال: WY-9082-X72"
            class="w-full bg-[#F8F8F8] dark:bg-[#222222] text-[#111111] dark:text-[#FFFFFF] px-4 py-3 rounded-2xl text-sm font-mono tracking-wider border border-[#E5E7EB] dark:border-[#27272A] focus:border-orange-500 focus:outline-none uppercase"
          />

          <p v-if="errorMessage" class="text-xs font-bold text-rose-500 flex items-center gap-1.5">
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{{ errorMessage }}</span>
          </p>

          <p v-if="successMessage" class="text-xs font-bold text-emerald-600 dark:text-emerald-400 flex items-center gap-1.5">
            <CheckCircle2 class="w-4 h-4 shrink-0" />
            <span>{{ successMessage }}</span>
          </p>

          <button
            type="submit"
            :disabled="isLoading"
            class="w-full py-3 rounded-2xl bg-orange-500 hover:bg-[#EA580C] text-white font-bold text-xs sm:text-sm shadow-sm shadow-orange-500/20 transition-all text-center mt-1 disabled:opacity-50 active:scale-95 cursor-pointer"
          >
            {{ isLoading ? "جاري التحقق والتفعيل..." : "تفعيل الاشتراك" }}
          </button>
        </form>
      </div>
    </div>
  </Transition>
</template>
