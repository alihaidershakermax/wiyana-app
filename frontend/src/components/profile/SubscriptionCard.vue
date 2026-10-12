<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import { SubscriptionInfo } from "@/types/subscription";
import {
  BookOpen,
  Calendar,
  Clock,
  KeyRound,
  RotateCcw,
  CheckCircle2,
  AlertCircle,
  ClockAlert,
  ArrowLeft,
} from "@lucide/vue";

const props = defineProps<{
  subscription?: SubscriptionInfo | null;
  isLoading?: boolean;
  isError?: boolean;
  progressPercentage?: number;
}>();

const emit = defineEmits<{
  (e: "activateCode"): void;
  (e: "retry"): void;
}>();

const router = useRouter();

// Normalize subscription status from real API
const status = computed<"active" | "expired" | "none">(() => {
  if (!props.subscription) return "none";
  if (props.subscription.status === "active") return "active";
  if (props.subscription.status === "expired") return "expired";
  return "none";
});

const courseTitle = computed(() => {
  return props.subscription?.planName || "المساق الدراسي";
});

const startDate = computed(() => props.subscription?.startDate || "");
const endDate = computed(() => props.subscription?.endDate || "");
const remainingDays = computed(() => props.subscription?.remainingDays ?? 0);

const subscriptionsList = computed(() => {
  return props.subscription?.subscriptionsList || [];
});

function handleContinueLearning() {
  router.push("/courses");
}
</script>

<template>
  <div class="w-full flex flex-col gap-2 text-start">
    <!-- iOS Section Label -->
    <div class="flex items-center justify-between px-3">
      <span class="text-xs font-bold text-[#6B7280] dark:text-[#A1A1AA] select-none">
        حالة الاشتراك
      </span>
      <button
        v-if="!isLoading && !isError"
        type="button"
        @click="emit('activateCode')"
        class="text-[11px] font-bold text-orange-500 hover:text-[#EA580C] transition-colors flex items-center gap-1 cursor-pointer"
      >
        <KeyRound class="w-3 h-3" />
        <span>{{ status === 'active' ? 'تفعيل كود آخر' : 'تفعيل كود' }}</span>
      </button>
    </div>

    <!-- 1. LOADING STATE: Skeleton -->
    <div
      v-if="isLoading"
      class="w-full bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-3xl p-5 shadow-xs animate-pulse flex flex-col gap-4"
    >
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-2xl bg-[#F8F8F8] dark:bg-[#222222]"></div>
        <div class="flex flex-col gap-2 flex-1">
          <div class="h-4 w-32 rounded bg-[#F8F8F8] dark:bg-[#222222]"></div>
          <div class="h-3 w-20 rounded bg-[#F8F8F8] dark:bg-[#222222]"></div>
        </div>
      </div>
      <div class="h-10 rounded-xl bg-[#F8F8F8] dark:bg-[#222222]"></div>
    </div>

    <!-- 2. ERROR STATE: Retry Card -->
    <div
      v-else-if="isError"
      class="w-full bg-white dark:bg-[#181818] border border-rose-200/80 dark:border-rose-900/40 rounded-3xl p-5 shadow-xs flex flex-col sm:flex-row items-center justify-between gap-4 text-center sm:text-start transition-colors"
    >
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-2xl bg-rose-500/10 text-rose-500 flex items-center justify-center shrink-0">
          <AlertCircle class="w-5 h-5" />
        </div>
        <div class="flex flex-col">
          <h4 class="text-xs sm:text-sm font-bold text-[#111111] dark:text-[#FFFFFF]">
            تعذر تحميل حالة الاشتراك
          </h4>
          <p class="text-[11px] text-[#6B7280] dark:text-[#A1A1AA] mt-0.5">
            تحقق من اتصال الإنترنت ثم أعد المحاولة
          </p>
        </div>
      </div>

      <button
        type="button"
        @click="emit('retry')"
        class="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-[#F8F8F8] dark:bg-[#222222] hover:bg-[#E5E7EB] dark:hover:bg-[#2A2A2A] text-[#111111] dark:text-white font-bold text-xs transition-all active:scale-95 cursor-pointer"
      >
        <RotateCcw class="w-3.5 h-3.5" />
        <span>إعادة المحاولة</span>
      </button>
    </div>

    <!-- 3. ACTIVE SUBSCRIPTION STATE -->
    <div
      v-else-if="status === 'active'"
      class="w-full bg-white dark:bg-[#181818] border border-orange-500/30 dark:border-orange-500/25 rounded-3xl p-5 shadow-xs flex flex-col gap-4 transition-colors"
    >
      <!-- Top Row: Course icon & Title & Active Badge -->
      <div class="flex items-start justify-between gap-3">
        <div class="flex items-center gap-3 min-w-0">
          <div class="w-11 h-11 rounded-2xl bg-orange-500 text-white flex items-center justify-center shrink-0 shadow-sm shadow-orange-500/20">
            <BookOpen class="w-5 h-5 stroke-[2.3]" />
          </div>
          <div class="flex flex-col min-w-0">
            <h4 class="text-sm sm:text-base font-black text-[#111111] dark:text-[#FFFFFF] truncate">
              {{ courseTitle }}
            </h4>
            <div class="flex items-center gap-1.5 mt-0.5">
              <span class="inline-flex items-center gap-1 text-[11px] font-bold text-emerald-600 dark:text-emerald-400 bg-emerald-500/10 px-2 py-0.5 rounded-full">
                <CheckCircle2 class="w-3 h-3" />
                <span>الاشتراك فعال</span>
              </span>
              <span v-if="remainingDays > 0" class="text-[11px] font-bold text-orange-600 dark:text-orange-400 bg-orange-500/10 px-2 py-0.5 rounded-full">
                {{ remainingDays }} يوم متبقي
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Multiple Subscriptions Sub-list if user is enrolled in >1 course -->
      <div v-if="subscriptionsList.length > 1" class="flex flex-col gap-1.5 pt-2 border-t border-[#E5E7EB] dark:border-[#27272A]">
        <span class="text-[10px] font-bold text-[#6B7280] dark:text-[#A1A1AA]">المساقات المشترك بها ({{ subscriptionsList.length }}):</span>
        <div class="flex flex-col gap-1">
          <div
            v-for="item in subscriptionsList"
            :key="item.id"
            class="flex items-center justify-between p-2 rounded-xl bg-[#F8F8F8] dark:bg-[#222222] text-xs font-bold"
          >
            <span class="truncate text-[#111111] dark:text-[#FFFFFF]">{{ item.courseTitle }}</span>
            <span class="text-[10px] text-orange-500 shrink-0">{{ item.remainingDays > 0 ? `${item.remainingDays} يوم` : 'فعال' }}</span>
          </div>
        </div>
      </div>

      <!-- Dates Row -->
      <div class="grid grid-cols-2 gap-2 bg-[#F8F8F8] dark:bg-[#222222] p-3 rounded-2xl text-xs border border-[#E5E7EB] dark:border-[#27272A]">
        <div class="flex items-center gap-2 text-[#6B7280] dark:text-[#A1A1AA]">
          <Calendar class="w-3.5 h-3.5 text-orange-500 shrink-0" />
          <div class="flex flex-col">
            <span class="text-[10px] text-[#9CA3AF] dark:text-[#71717A] font-semibold">تاريخ البدء</span>
            <span class="font-bold text-[#111111] dark:text-[#FFFFFF] text-xs">{{ startDate || "غير محدد" }}</span>
          </div>
        </div>

        <div class="flex items-center gap-2 text-[#6B7280] dark:text-[#A1A1AA]">
          <Clock class="w-3.5 h-3.5 text-orange-500 shrink-0" />
          <div class="flex flex-col">
            <span class="text-[10px] text-[#9CA3AF] dark:text-[#71717A] font-semibold">تاريخ الانتهاء</span>
            <span class="font-bold text-[#111111] dark:text-[#FFFFFF] text-xs">{{ endDate || "مستمر" }}</span>
          </div>
        </div>
      </div>

      <!-- Optional Progress Indicator -->
      <div v-if="progressPercentage !== undefined && progressPercentage > 0" class="flex flex-col gap-1.5">
        <div class="flex items-center justify-between text-[11px] font-bold">
          <span class="text-[#6B7280] dark:text-[#A1A1AA]">التقدم الدراسي</span>
          <span class="text-orange-500">{{ progressPercentage }}%</span>
        </div>
        <div class="w-full h-2 rounded-full bg-[#F8F8F8] dark:bg-[#222222] overflow-hidden">
          <div
            class="h-full bg-orange-500 rounded-full transition-all duration-500"
            :style="{ width: `${Math.min(100, Math.max(0, progressPercentage))}%` }"
          ></div>
        </div>
      </div>

      <!-- Action Button: متابعة التعلم -->
      <button
        type="button"
        @click="handleContinueLearning"
        class="w-full py-3 px-4 rounded-2xl bg-orange-500 hover:bg-[#EA580C] text-white font-bold text-xs sm:text-sm shadow-sm shadow-orange-500/20 transition-all flex items-center justify-center gap-2 active:scale-95 cursor-pointer"
      >
        <span>متابعة التعلم</span>
        <ArrowLeft class="w-4 h-4" />
      </button>
    </div>

    <!-- 4. EXPIRED SUBSCRIPTION STATE -->
    <div
      v-else-if="status === 'expired'"
      class="w-full bg-white dark:bg-[#181818] border border-amber-200/80 dark:border-amber-900/40 rounded-3xl p-5 shadow-xs flex flex-col gap-4 transition-colors"
    >
      <div class="flex items-start justify-between gap-3">
        <div class="flex items-center gap-3 min-w-0">
          <div class="w-11 h-11 rounded-2xl bg-amber-500/10 text-amber-500 flex items-center justify-center shrink-0">
            <ClockAlert class="w-5 h-5" />
          </div>
          <div class="flex flex-col min-w-0">
            <h4 class="text-sm sm:text-base font-black text-[#111111] dark:text-[#FFFFFF] truncate">
              {{ courseTitle }}
            </h4>
            <div class="flex items-center gap-1.5 mt-0.5">
              <span class="inline-flex items-center gap-1 text-[11px] font-bold text-amber-600 dark:text-amber-400 bg-amber-500/10 px-2 py-0.5 rounded-full">
                انتهى الاشتراك
              </span>
              <span v-if="endDate" class="text-[11px] text-[#6B7280] dark:text-[#A1A1AA]">
                انتهى في {{ endDate }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <button
        type="button"
        @click="emit('activateCode')"
        class="w-full py-3 px-4 rounded-2xl bg-orange-500 hover:bg-[#EA580C] text-white font-bold text-xs sm:text-sm shadow-sm shadow-orange-500/20 transition-all flex items-center justify-center gap-2 active:scale-95 cursor-pointer"
      >
        <KeyRound class="w-4 h-4" />
        <span>تفعيل اشتراك جديد</span>
      </button>
    </div>

    <!-- 5. NO SUBSCRIPTION (EMPTY) STATE -->
    <div
      v-else
      class="w-full bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-3xl p-5 shadow-xs flex flex-col sm:flex-row items-center justify-between gap-4 text-center sm:text-start transition-colors"
    >
      <div class="flex items-center gap-3.5">
        <div class="w-11 h-11 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center shrink-0">
          <KeyRound class="w-5 h-5 stroke-[2.3]" />
        </div>
        <div class="flex flex-col">
          <h4 class="text-xs sm:text-sm font-bold text-[#111111] dark:text-[#FFFFFF]">
            لا يوجد اشتراك فعال
          </h4>
          <p class="text-[11px] text-[#6B7280] dark:text-[#A1A1AA] mt-0.5">
            اشترك للوصول إلى المساقات التعليمية
          </p>
        </div>
      </div>

      <button
        type="button"
        @click="emit('activateCode')"
        class="w-full sm:w-auto px-5 py-2.5 rounded-xl bg-orange-500 hover:bg-[#EA580C] text-white font-bold text-xs shadow-sm shadow-orange-500/20 transition-all active:scale-95 text-center shrink-0 cursor-pointer flex items-center justify-center gap-1.5"
      >
        <KeyRound class="w-3.5 h-3.5" />
        <span>تفعيل اشتراك</span>
      </button>
    </div>
  </div>
</template>
