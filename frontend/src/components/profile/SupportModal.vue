<script setup lang="ts">
import { ref } from "vue";
import { X, Headphones, Send, MessageCircle, Check, Sparkles } from "@lucide/vue";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "sent"): void;
}>();

const ticketSubject = ref("");
const ticketMessage = ref("");
const isSubmitting = ref(false);
const isSuccess = ref(false);

function submitTicket() {
  if (!ticketMessage.value.trim()) return;
  isSubmitting.value = true;
  setTimeout(() => {
    isSubmitting.value = false;
    isSuccess.value = true;
    setTimeout(() => {
      isSuccess.value = false;
      ticketSubject.value = "";
      ticketMessage.value = "";
      emit("sent");
      emit("close");
    }, 1500);
  }, 900);
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
      <div class="w-full max-w-lg bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-3xl p-6 shadow-2xl flex flex-col gap-5 my-8">
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-800/80 pb-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center">
              <Headphones class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-lg font-black text-slate-900 dark:text-slate-100">
                مركز الدعم الفني والمساعدة
              </h3>
              <p class="text-xs text-slate-500 dark:text-slate-400">
                فريقنا متاح للرد على جميع استفساراتك 24/7
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

        <!-- Direct Channels -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <a
            href="https://t.me/wiyana_support"
            target="_blank"
            rel="noopener noreferrer"
            class="p-4 rounded-2xl bg-orange-500/10 hover:bg-orange-500/20 border border-orange-500/20 flex items-center gap-3 transition-colors group text-right cursor-pointer"
          >
            <div class="w-10 h-10 rounded-xl bg-orange-500 text-white flex items-center justify-center shrink-0">
              <Send class="w-5 h-5" />
            </div>
            <div>
              <h4 class="text-xs font-bold text-slate-900 dark:text-slate-100 group-hover:text-orange-500 transition-colors">
                دعم تيليجرام الفوري
              </h4>
              <p class="text-[11px] text-slate-500 dark:text-slate-400">@wiyana_support</p>
            </div>
          </a>

          <a
            href="https://wa.me/9647700000000"
            target="_blank"
            rel="noopener noreferrer"
            class="p-4 rounded-2xl bg-orange-500/10 hover:bg-orange-500/20 border border-orange-500/20 flex items-center gap-3 transition-colors group text-right cursor-pointer"
          >
              <div class="w-10 h-10 rounded-xl bg-orange-500 text-white flex items-center justify-center shrink-0">
              <MessageCircle class="w-5 h-5" />
            </div>
            <div>
                <h4 class="text-xs font-bold text-slate-900 dark:text-slate-100 group-hover:text-orange-500 transition-colors">
                خدمة واتساب
              </h4>
              <p class="text-[11px] text-slate-500 dark:text-slate-400">محادثة مباشرة</p>
            </div>
          </a>
        </div>

        <!-- Send Direct Ticket Form -->
        <form @submit.prevent="submitTicket" class="flex flex-col gap-3 pt-2 border-t border-slate-100 dark:border-slate-800/80">
          <span class="text-xs font-bold text-slate-700 dark:text-slate-300">أو أرسل تذكرة دعم مباشرة داخل التطبيق:</span>

          <input
            v-model="ticketSubject"
            type="text"
            required
            placeholder="موضوع المشكلة (مثال: كود الاشتراك، تشغيل الفيديو...)"
            class="w-full bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 px-4 py-2.5 rounded-2xl text-xs font-semibold border border-slate-200 dark:border-slate-800 focus:border-orange-500 focus:outline-none"
          />

          <textarea
            v-model="ticketMessage"
            required
            rows="3"
            placeholder="اشرح المشكلة بالتفصيل لنتمكن من مساعدتك فوراً..."
            class="w-full bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 px-4 py-2.5 rounded-2xl text-xs font-semibold border border-slate-200 dark:border-slate-800 focus:border-orange-500 focus:outline-none resize-none"
          ></textarea>

          <p v-if="isSuccess" class="text-xs font-bold text-orange-500 flex items-center gap-1.5">
            <Check class="w-4 h-4" />
            تم إرسال تذكرتك بنجاح! سيتم التواصل معك عبر تيليجرام خلال دقائق.
          </p>

          <div class="flex items-center gap-3 mt-1">
            <button
              type="button"
              @click="emit('close')"
              class="flex-1 py-3 rounded-2xl border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-bold text-xs sm:text-sm hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
            >
              إلغاء
            </button>
            <button
              type="submit"
              :disabled="isSubmitting"
              class="flex-1 py-3 rounded-2xl bg-orange-500 hover:bg-orange-600 text-white font-bold text-xs sm:text-sm shadow-md shadow-orange-500/20 transition-all text-center flex items-center justify-center gap-1.5 active:scale-95 disabled:opacity-50 cursor-pointer"
            >
              <Send class="w-4 h-4" />
              {{ isSubmitting ? 'جاري الإرسال...' : 'إرسال التذكرة' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </Transition>
</template>
