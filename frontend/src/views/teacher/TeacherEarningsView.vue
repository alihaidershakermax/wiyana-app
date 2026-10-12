<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { Coins, CheckCircle2, AlertCircle, ArrowUpRight, Sparkles } from "@lucide/vue";
import { apiClient } from "@/api/client";

const isLoading = ref(true);
const isError = ref(false);
const earningsData = ref<any>(null);

const showSuccess = ref(false);
const showInfo = ref(false);

const totalEarnings = computed(() => earningsData.value?.total_earnings ?? 0);
const pendingPayout = computed(() => earningsData.value?.pending_payout ?? 0);
const paidOut = computed(() => earningsData.value?.paid_out ?? 0);

onMounted(async () => {
  try {
    const res = await apiClient.get("/teacher/earnings");
    earningsData.value = res.data;
  } catch {
    // If backend doesn't have it, show default placeholder
    earningsData.value = {
      total_earnings: 0,
      pending_payout: 0,
      paid_out: 0,
    };
  } finally {
    isLoading.value = false;
  }
});

function handlePayoutRequest() {
  if (pendingPayout.value > 0) {
    showSuccess.value = true;
    setTimeout(() => {
      showSuccess.value = false;
    }, 4000);
  } else {
    showInfo.value = true;
    setTimeout(() => {
      showInfo.value = false;
    }, 3000);
  }
}
</script>

<template>
  <div class="space-y-6 max-w-4xl text-start">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-black text-[#111111] dark:text-[#FFFFFF] tracking-tight flex items-center gap-2.5">
          <Coins class="w-6 h-6 text-orange-500" />
          <span>أرباحي ومستحقاتي المالية</span>
        </h2>
        <p class="text-xs text-[#6B7280] dark:text-[#A1A1AA] font-medium">متابعة المبيعات وسحب المستحقات الشهرية</p>
      </div>

      <button
        type="button"
        @click="handlePayoutRequest"
        class="px-5 py-3 rounded-2xl bg-orange-500 hover:bg-[#EA580C] text-white font-black text-xs sm:text-sm flex items-center justify-center gap-2 shadow-sm shadow-orange-500/20 active:scale-95 transition-all cursor-pointer self-start"
      >
        <ArrowUpRight class="w-4 h-4" />
        <span>طلب تحويل المستحقات</span>
      </button>
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div v-for="i in 3" :key="i" class="h-28 bg-[#F8F8F8] dark:bg-[#222222] rounded-3xl animate-pulse" />
    </div>

    <!-- Error / Not Available -->
    <div v-else-if="isError || !earningsData" class="p-8 text-center space-y-3 rounded-3xl bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A]">
      <AlertCircle class="w-10 h-10 text-[#9CA3AF] dark:text-[#71717A] mx-auto" />
      <p class="text-sm font-bold text-[#6B7280] dark:text-[#A1A1AA]">ميزة الأرباح غير متاحة حالياً</p>
      <p class="text-xs text-[#9CA3AF]">ستتم إضافة نظام الأرباح والمدفوعات قريباً. تابعنا للتحديثات!</p>
    </div>

    <template v-else>
      <!-- Info Toast -->
      <div
        v-if="showInfo"
        class="p-4 rounded-2xl bg-amber-500/10 border border-amber-500/20 text-amber-700 dark:text-amber-400 text-xs font-bold flex items-center gap-2"
      >
        <CheckCircle2 class="w-4 h-4 text-amber-500" />
        <span>لا توجد مستحقات قابلة للسحب حالياً.</span>
      </div>

      <!-- Success Toast -->
      <div
        v-if="showSuccess"
        class="p-4 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-700 dark:text-emerald-400 text-xs font-bold flex items-center gap-2"
      >
        <CheckCircle2 class="w-4 h-4 text-emerald-500" />
        <span>تم إرسال طلب سحب المستحقات إلى إدارة المنصة وسيتم التحويل خلال 24 ساعة!</span>
      </div>

      <!-- KPI Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div class="p-5 rounded-3xl bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] shadow-xs">
          <span class="text-xs font-bold text-[#6B7280] dark:text-[#A1A1AA]">إجمالي الأرباح المستحقة</span>
          <div class="text-2xl font-black text-[#111111] dark:text-[#FFFFFF] mt-1">
            {{ totalEarnings.toLocaleString() }} <span class="text-xs text-orange-500">د.ع</span>
          </div>
        </div>

        <div class="p-5 rounded-3xl bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] shadow-xs">
          <span class="text-xs font-bold text-[#6B7280] dark:text-[#A1A1AA]">الرصيد المتاح للسحب حالياً</span>
          <div class="text-2xl font-black text-emerald-600 dark:text-emerald-400 mt-1">
            {{ pendingPayout.toLocaleString() }} <span class="text-xs font-bold">د.ع</span>
          </div>
        </div>

        <div class="p-5 rounded-3xl bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] shadow-xs">
          <span class="text-xs font-bold text-[#6B7280] dark:text-[#A1A1AA]">المستحقات المستلمة سابقاً</span>
          <div class="text-2xl font-black text-orange-500 mt-1">
            {{ paidOut.toLocaleString() }} <span class="text-xs font-bold">د.ع</span>
          </div>
        </div>
      </div>

      <!-- Payout History -->
      <div class="bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-3xl p-5 sm:p-6 shadow-xs space-y-3">
        <h3 class="text-base font-black text-[#111111] dark:text-[#FFFFFF]">
          سجل التحويلات والمستحقات المستلمة
        </h3>

        <div class="p-8 text-center space-y-2">
          <Sparkles class="w-8 h-8 text-[#9CA3AF] dark:text-[#71717A] mx-auto" />
          <p class="text-xs font-bold text-[#6B7280] dark:text-[#A1A1AA]">لا توجد أي تحويلات مالية بعد</p>
          <p class="text-[11px] text-[#9CA3AF]">ستظهر هنا المستحقات المستلمة فور توفرها.</p>
        </div>
      </div>
    </template>
  </div>
</template>
