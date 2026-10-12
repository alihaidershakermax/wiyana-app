<script setup lang="ts">
import { useRouter } from "vue-router";
import { ShieldX, ArrowRight, LayoutDashboard } from "@lucide/vue";
import { useAdminStore } from "@/admin/stores/useAdminStore";

const router = useRouter();
const adminStore = useAdminStore();

function goToDashboard() {
  if (adminStore.userRole === "TEACHER") {
    router.push("/admin/teacher/dashboard");
  } else {
    router.push("/admin/dashboard");
  }
}
</script>

<template>
  <div
    dir="rtl"
    class="min-h-screen bg-[#080808] text-[#F5F5F5] flex items-center justify-center p-6"
  >
    <div class="card max-w-md w-full text-center p-8 space-y-6 border-[#242424] bg-[#111111] shadow-2xl">
      <div
        class="mx-auto h-20 w-20 rounded-3xl bg-red-950/40 border border-red-900/60 flex items-center justify-center"
      >
        <ShieldX class="h-10 w-10 text-red-500" />
      </div>
      <div class="space-y-2">
        <h1 class="text-4xl font-black tracking-tight text-white">403</h1>
        <h2 class="text-lg font-bold text-white">لا تملك صلاحية الوصول إلى هذه الصفحة</h2>
        <p class="text-xs text-[#A1A1AA] leading-relaxed">
          حسابك لا يملك الصلاحيات اللازمة لعرض هذا المحتوى. إذا كنت تعتقد أن هذا خطأ،
          يرجى التواصل مع إدارة المنصة.
        </p>
      </div>
      <div class="flex items-center justify-center gap-3 pt-2">
        <button class="btn-secondary text-xs" @click="router.back()">
          <ArrowRight class="h-4 w-4" />
          العودة
        </button>
        <button class="btn-primary text-xs" @click="goToDashboard">
          <LayoutDashboard class="h-4 w-4" />
          لوحة التحكم
        </button>
      </div>
    </div>
  </div>
</template>
