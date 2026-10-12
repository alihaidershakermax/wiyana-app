<script setup lang="ts">
import { ref, computed } from "vue";
import { useRouter, useRoute } from "vue-router";
import { Loader2, ShieldAlert, Send, Lock, UserCheck, KeyRound } from "@lucide/vue";
import BrandLogo from "@/components/BrandLogo.vue";
import { useAdminStore } from "@/admin/stores/useAdminStore";

const router = useRouter();
const route = useRoute();
const adminStore = useAdminStore();

const identifier = ref("");
const isLoading = ref(false);
const localError = ref<string | null>(null);

const errorMessage = computed(() => localError.value || adminStore.error);

// Redirect target (preserved from ?redirect= or role home).
function homePath(): string {
  const redirect = typeof route.query.redirect === "string" ? route.query.redirect : null;
  if (redirect && redirect.startsWith("/admin")) return redirect;
  return adminStore.isOwner ? "/admin/dashboard" : "/admin/teacher/dashboard";
}

async function handleDirectLogin() {
  localError.value = null;
  const clean = identifier.value.trim();
  if (!clean) {
    localError.value = "يرجى كتابة معرف تليجرام (ID) أو اسم المستخدم (Username)";
    return;
  }

  isLoading.value = true;
  try {
    const ok = await adminStore.loginIdentifier(clean);
    if (ok) {
      router.replace(homePath());
    }
  } catch (err: any) {
    localError.value = err?.message || "حدث خطأ أثناء تسجيل الدخول";
  } finally {
    isLoading.value = false;
  }
}

function goBack() {
  window.location.href = "/";
}
</script>

<template>
  <div
    dir="rtl"
    class="min-h-screen bg-slate-950 text-slate-100 flex flex-col items-center justify-center p-6 transition-colors duration-200 selection:bg-orange-500 selection:text-white"
  >
    <!-- Background decoration -->
    <div class="pointer-events-none fixed inset-0 overflow-hidden" aria-hidden="true">
      <div class="absolute -top-32 -left-32 h-96 w-96 rounded-full bg-orange-500/10 blur-3xl" />
      <div class="absolute -bottom-32 -right-32 h-96 w-96 rounded-full bg-orange-500/10 blur-3xl" />
    </div>

    <div class="relative w-full max-w-md">
      <!-- Card -->
      <div
        class="glass-liquid-strong card rounded-3xl shadow-2xl p-8 sm:p-10 space-y-6"
      >
        <!-- Logo & Title -->
        <div class="flex flex-col items-center text-center">
          <BrandLogo size="lg" subtitle="لوحة التحكم" />
          <p class="mt-5 text-sm font-bold text-orange-400">
            تسجيل الدخول للوحة التحكم
          </p>
          <p class="mt-2 text-xs text-slate-400 leading-relaxed">
            أدخل معرف تليجرام الخاص بك (ID) أو اسم المستخدم للوصول إلى اللوحة
          </p>
        </div>

        <!-- Direct ID / Username Login Form -->
        <form @submit.prevent="handleDirectLogin" class="space-y-4 pt-2">
          <div>
            <label class="block text-xs font-bold text-slate-300 mb-2">
              معرف تليجرام (Telegram ID) أو اليوزر نيم:
            </label>
            <div class="relative">
              <input
                v-model="identifier"
                type="text"
                placeholder="مثال: 2018954602 أو @dextermorgenk"
                class="glass-liquid-input w-full px-4 py-3.5 rounded-xl text-white placeholder-slate-500 focus:outline-none text-sm transition-all dir-ltr text-right"
                :disabled="isLoading"
                required
              />
              <UserCheck class="absolute left-3.5 top-3.5 h-5 w-5 text-slate-500 pointer-events-none" />
            </div>
          </div>

          <!-- Error Alert -->
          <div
            v-if="errorMessage"
            class="p-3.5 rounded-xl bg-red-950/60 border border-red-800/80 text-xs text-red-300 text-center flex items-center justify-center gap-2"
          >
            <ShieldAlert class="h-4 w-4 shrink-0 text-red-400" />
            <span>{{ errorMessage }}</span>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="isLoading"
            class="w-full py-3.5 px-6 rounded-xl bg-gradient-to-r from-orange-500 to-amber-600 hover:from-orange-400 hover:to-amber-500 text-white font-bold text-sm shadow-lg shadow-orange-500/25 hover:shadow-orange-500/40 flex items-center justify-center gap-2 transition-all duration-200 hover:scale-[1.01] active:scale-[0.99] disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <Loader2 v-if="isLoading" class="h-4 w-4 animate-spin" />
            <KeyRound v-else class="h-4 w-4" />
            <span>{{ isLoading ? 'جارٍ التحقق...' : 'دخول لوحة التحكم' }}</span>
          </button>
        </form>

        <!-- Footer links -->
        <div class="pt-4 border-t border-slate-800/80 flex items-center justify-center gap-1 text-xs text-slate-500">
          <Lock class="h-3 w-3" />
          <span>لوحة الإدارة المشفرة الآمنة</span>
        </div>
      </div>

      <!-- Back to platform -->
      <div class="mt-6 text-center">
        <button
          class="inline-flex items-center gap-1.5 text-xs font-medium text-slate-400 hover:text-orange-400 transition-colors"
          @click="goBack"
        >
          <Send class="h-3.5 w-3.5 rotate-180" />
          العودة إلى منصة ويانا
        </button>
      </div>
    </div>
  </div>
</template>
