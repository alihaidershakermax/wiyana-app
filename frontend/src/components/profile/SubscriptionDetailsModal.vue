<script setup lang="ts">
import { X, Crown, CheckCircle2, ShieldCheck, Zap, DownloadCloud, Sparkles, KeyRound } from "@lucide/vue";
import { SubscriptionInfo } from "@/types/subscription";

defineProps<{
  isOpen: boolean;
  subscription?: SubscriptionInfo | null;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "activateCode"): void;
}>();
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
              <Crown class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-lg font-black text-slate-900 dark:text-slate-100">
                تفاصيل الاشتراك
              </h3>
              <p class="text-xs text-slate-500 dark:text-slate-400">
                {{ subscription ? 'تفاصيل خطتك الحالية والوصول المفتوح' : 'لم يتم تفعيل اشتراك بعد' }}
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

        <!-- Active Subscription -->
        <template v-if="subscription && subscription.status === 'active'">
          <!-- Plan Status Banner -->
          <div class="p-4 rounded-2xl bg-gradient-to-r from-orange-500/15 via-orange-500/10 to-orange-500/5 border border-orange-500/20 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <Sparkles class="w-6 h-6 text-orange-500 shrink-0" />
              <div>
                <h4 class="text-xs sm:text-sm font-black text-slate-900 dark:text-slate-100">
                  {{ subscription.planName }}
                </h4>
                <p class="text-[11px] text-orange-600 dark:text-orange-400 font-bold">
                  نشط حتى {{ subscription.endDate }} (متبقي {{ subscription.remainingDays }} يوم)
                </p>
              </div>
            </div>
            <span class="px-2.5 py-1 rounded-full bg-orange-500/10 text-orange-500 text-[10px] font-black border border-orange-500/20">
              مفعّل
            </span>
          </div>

          <!-- Features List -->
          <div class="flex flex-col gap-2.5">
            <span class="text-xs font-black text-slate-700 dark:text-slate-300">ما يشمله اشتراكك:</span>
            <div
              v-for="(feat, idx) in [
                { title: 'مشاهدة جميع الدروس والمحاضرات', desc: 'وصول كامل لجميع المحتوى التعليمي' },
                { title: 'تنزيل الدروس للمشاهدة بدون إنترنت', desc: 'احفظ المحاضرات على جهازك واستمتع بالمشاهدة في أي مكان' },
                { title: 'متابعة التقدم والإحصائيات', desc: 'تتبع نسب إنجازك وساعات مشاهدتك' },
                { title: 'الوصول لملخصات واختبارات', desc: 'بنك أسئلة شامل مع تصحيح ذكي' },
                { title: 'إشعارات تذكير بالمحاضرات', desc: 'تنبيهات فورية عن الدروس الجديدة' },
              ]"
              :key="idx"
              class="flex items-start gap-3 p-3 rounded-2xl bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800/60"
            >
              <CheckCircle2 class="w-4 h-4 text-orange-500 shrink-0 mt-0.5" />
              <div class="flex flex-col gap-0.5">
                <span class="text-xs font-bold text-slate-900 dark:text-slate-100">{{ feat.title }}</span>
                <span class="text-[11px] text-slate-500 dark:text-slate-400">{{ feat.desc }}</span>
              </div>
            </div>
          </div>
        </template>

        <!-- No Subscription -->
        <template v-else>
          <div class="flex flex-col items-center justify-center py-8 gap-4">
            <div class="w-16 h-16 rounded-3xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center">
              <KeyRound class="w-8 h-8 text-slate-400" />
            </div>
            <div class="text-center">
              <h4 class="text-base font-bold text-slate-900 dark:text-slate-100 mb-1">لا يوجد اشتراك نشط</h4>
              <p class="text-sm text-slate-500 dark:text-slate-400">قم بتفعيل كود الاشتراك للوصول إلى كافة الدروس والمواد الدراسية.</p>
            </div>
          </div>
        </template>

        <!-- Actions -->
        <div class="flex items-center gap-3 pt-2">
          <button
            type="button"
            @click="emit('close')"
            class="flex-1 py-3 rounded-2xl border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-bold text-xs sm:text-sm hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
          >
            إغلاق
          </button>
          <button
            type="button"
            @click="emit('close'); emit('activateCode')"
            class="flex-1 py-3 rounded-2xl bg-orange-500 hover:bg-orange-600 text-white font-bold text-xs sm:text-sm shadow-md shadow-orange-500/20 transition-all text-center active:scale-95 cursor-pointer"
          >
            تفعيل كود اشتراك
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>
