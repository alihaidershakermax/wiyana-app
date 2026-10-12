<script setup lang="ts">
import { useRouter } from "vue-router";
import { ShieldAlert, LogIn } from "@lucide/vue";

const router = useRouter();

function goToLogin() {
  router.push("/admin/login");
}
</script>

<template>
  <div
    dir="rtl"
    class="min-h-screen bg-[#080808] text-[#F5F5F5] flex items-center justify-center p-6"
  >
    <div class="card max-w-md w-full text-center p-8 space-y-6 border-[#242424] bg-[#111111] shadow-2xl">
      <div
        class="mx-auto h-20 w-20 rounded-3xl bg-[#1F0E03] border border-orange-900/60 flex items-center justify-center"
      >
        <ShieldAlert class="h-10 w-10 text-[#F97316]" />
      </div>
      <div class="space-y-2">
        <h1 class="text-xl font-extrabold tracking-tight text-white">تعذر الوصول</h1>
        <p class="text-xs text-[#A1A1AA] leading-relaxed">
          يرجى تسجيل الدخول بحساب يملك الصلاحيات اللازمة للوصول إلى لوحة التحكم.
        </p>
      </div>
      <button class="btn-primary text-xs w-full py-3" @click="goToLogin">
        <LogIn class="h-4 w-4" />
        تسجيل الدخول إلى لوحة التحكم
      </button>
    </div>
  </div>
</template>
