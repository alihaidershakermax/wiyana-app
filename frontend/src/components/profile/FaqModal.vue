<script setup lang="ts">
import { ref } from "vue";
import { X, HelpCircle, ChevronDown, Sparkles } from "@lucide/vue";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const activeIndex = ref<number | null>(0);

const faqs = [
  {
    q: "كيف أقوم بتفعيل كود شحن الاشتراك؟",
    a: "يمكنك الضغط على زر 'تفعيل كود' في بطاقة الاشتراك بأعلى الصفحة، ثم كتابة الرمز الموجود على كارت الشحن (مثل WY-9081-VIP) والضغط على 'تفعيل الكود'."
  },
  {
    q: "هل يمكنني مشاهدة المحاضرات بدون اتصال بالإنترنت؟",
    a: "نعم بكل تأكيد! يمكنك الضغط على أيقونة التنزيل داخل صفحة أي محاضرة لتنزيلها بجودة عالية ومشاهدتها أوفلاين في أي وقت بدون استهلاك بيانات."
  },
  {
    q: "كيف أتحكم في سرعة وجودة الفيديو؟",
    a: "من خلال قسم 'إعدادات مشغّل الفيديو' في صفحة حسابي، يمكنك اختيار جودة التشغيل الافتراضية (1080p, 720p...) وسرعة التشغيل (1.25x, 1.5x...) وسيتم تطبيقها على كافة المحاضرات تلقائياً."
  },
  {
    q: "ما هي ميزة ربط حساب تيليجرام؟",
    a: "ربط حسابك يتيح لك الدخول بضغطة زر وتلقي إشعارات الدروس الجديدة، مواعيد البث المباشر، وتنبيهات الواجبات فور نشرها."
  },
  {
    q: "هل يتم حفظ نقطة توقف المشاهدة تلقائياً؟",
    a: "نعم! بفضل ميزة Auto-Save Progress، يستمر تشغيل المحاضرة من آخر ثانية توقفت عندها بدقة عبر جميع أجهزتك."
  },
];

function toggleFaq(index: number) {
  activeIndex.value = activeIndex.value === index ? null : index;
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
              <HelpCircle class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-lg font-black text-slate-900 dark:text-slate-100">
                الأسئلة الشائعة والمعلومات
              </h3>
              <p class="text-xs text-slate-500 dark:text-slate-400">
                إجابات سريعة على أكثر الأسئلة تكراراً
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

        <!-- FAQ Items -->
        <div class="flex flex-col gap-2.5 max-h-96 overflow-y-auto pr-1">
          <div
            v-for="(item, idx) in faqs"
            :key="idx"
            class="rounded-2xl border border-slate-100 dark:border-slate-800/60 overflow-hidden bg-slate-50 dark:bg-slate-950 transition-colors"
          >
            <button
              @click="toggleFaq(idx)"
              type="button"
              class="w-full p-4 flex items-center justify-between gap-3 text-right focus:outline-none cursor-pointer"
            >
              <span class="text-xs sm:text-sm font-bold text-slate-900 dark:text-slate-100">
                {{ item.q }}
              </span>
              <ChevronDown
                :class="[
                  'w-4 h-4 text-orange-500 transition-transform duration-200 shrink-0',
                  activeIndex === idx ? 'rotate-180' : 'rotate-0'
                ]"
              />
            </button>
            <div
              v-show="activeIndex === idx"
              class="px-4 pb-4 pt-1 text-xs text-slate-600 dark:text-slate-400 leading-relaxed border-t border-slate-100 dark:border-slate-800/40"
            >
              {{ item.a }}
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="pt-2 border-t border-slate-100 dark:border-slate-800/80">
          <button
            type="button"
            @click="emit('close')"
            class="w-full py-3 rounded-2xl border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-bold text-xs sm:text-sm hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
          >
            إغلاق
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>
