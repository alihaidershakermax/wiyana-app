<script setup lang="ts">
import { ref, onMounted } from "vue";
import { teacherService } from "@/services/teacherService";
import { StudentQuestion } from "@/types/teacher";
import {
  MessageSquareQuote,
  Send,
  CheckCircle2,
  Clock,
  Play,
  User,
  Sparkles,
} from "@lucide/vue";

const questions = ref<StudentQuestion[]>([]);
const activeReplyId = ref<string | null>(null);
const replyInput = ref("");
const isLoading = ref(true);

async function loadQuestions() {
  isLoading.value = true;
  try {
    const list = await teacherService.getTeacherQuestions();
    questions.value = list;
  } finally {
    isLoading.value = false;
  }
}

onMounted(() => {
  loadQuestions();
});

async function submitReply(q: StudentQuestion) {
  if (!replyInput.value.trim()) return;

  const previousStatus = q.status;
  const previousReply = q.replyText;
  const previousRepliedAt = q.repliedAt;
  const reply = replyInput.value;

  q.replyText = reply;
  q.status = "answered";
  q.repliedAt = new Date().toLocaleString("ar");
  replyInput.value = "";
  activeReplyId.value = null;

  try {
    await teacherService.replyToQuestion(q.id, reply);
  } catch {
    q.status = previousStatus;
    q.replyText = previousReply;
    q.repliedAt = previousRepliedAt;
  }
}
</script>

<template>
  <div class="space-y-6 max-w-4xl text-start">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-black text-[#111111] dark:text-[#FFFFFF] tracking-tight flex items-center gap-2.5">
          <MessageSquareQuote class="w-6 h-6 text-orange-500" />
          <span>أسئلة واستفسارات الطلاب على المحاضرات</span>
        </h2>
        <p class="text-xs text-[#6B7280] dark:text-[#A1A1AA] font-medium">الرد المباشر على استفسارات الطلاب المرتبطة بدقائق وثواني الفيديو</p>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="space-y-4">
      <div v-for="i in 2" :key="i" class="h-40 bg-[#F8F8F8] dark:bg-[#222222] rounded-3xl animate-pulse" />
    </div>

    <!-- Empty state -->
    <div v-else-if="!questions.length" class="bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-3xl p-10 text-center space-y-3">
      <MessageSquareQuote class="w-10 h-10 text-[#9CA3AF] dark:text-[#71717A] mx-auto" />
      <p class="text-sm font-bold text-[#6B7280] dark:text-[#A1A1AA]">لا توجد أسئلة من الطلاب حالياً</p>
      <p class="text-xs text-[#9CA3AF]">ستظهر هنا أسئلة واستفسارات الطلاب على محاضراتك.</p>
    </div>

    <!-- Questions Feed -->
    <div v-else class="space-y-4">
      <div
        v-for="q in questions"
        :key="q.id"
        class="bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-3xl p-5 sm:p-6 shadow-xs space-y-4"
      >
        <!-- Top Meta: Student & Lesson Timestamp -->
        <div class="flex items-start justify-between gap-3">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center font-bold">
              {{ q.studentName.charAt(0) }}
            </div>
            <div>
              <div class="font-black text-sm text-[#111111] dark:text-[#FFFFFF]">
                {{ q.studentName }}
              </div>
              <div class="text-xs text-[#6B7280] dark:text-[#A1A1AA] font-medium flex items-center gap-2 mt-0.5">
                <span>{{ q.lessonTitle }}</span>
                <span v-if="q.timestampInVideo" class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md bg-orange-500/10 text-orange-600 dark:text-orange-400 font-mono font-bold text-[10px]">
                  <Play class="w-2.5 h-2.5" />
                  {{ q.timestampInVideo }}
                </span>
              </div>
            </div>
          </div>

          <span
            :class="[
              'px-2.5 py-1 rounded-xl text-[10px] font-bold',
              q.status === 'answered'
                ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
                : 'bg-amber-500/10 text-amber-600 dark:text-amber-400'
            ]"
          >
            {{ q.status === 'answered' ? 'تم الرد' : 'بانتظار الإجابة' }}
          </span>
        </div>

        <!-- Question Text Bubble -->
        <div class="p-4 rounded-2xl bg-[#F8F8F8] dark:bg-[#222222] border border-[#E5E7EB] dark:border-[#27272A] text-xs font-medium text-[#111111] dark:text-[#FFFFFF] leading-relaxed">
          {{ q.questionText }}
        </div>

        <!-- Teacher Existing Reply -->
        <div
          v-if="q.replyText"
          class="p-4 rounded-2xl bg-orange-500/10 border border-orange-500/20 text-xs font-medium text-[#111111] dark:text-[#FFFFFF] space-y-1.5"
        >
          <div class="flex items-center justify-between text-[11px] font-bold text-orange-500">
            <span>إجابتك للأستاذ:</span>
            <span class="font-mono">{{ q.repliedAt }}</span>
          </div>
          <p class="leading-relaxed">{{ q.replyText }}</p>
        </div>

        <!-- Reply Input Trigger & Box -->
        <div v-if="!q.replyText">
          <div v-if="activeReplyId !== q.id">
            <button
              type="button"
              @click="activeReplyId = q.id"
              class="px-4 py-2 rounded-xl bg-orange-500 hover:bg-[#EA580C] text-white font-bold text-xs flex items-center gap-2 cursor-pointer transition-all shadow-sm shadow-orange-500/20"
            >
              <Send class="w-3.5 h-3.5" />
              <span>كتابة الإجابة للطالب</span>
            </button>
          </div>

          <div v-else class="space-y-3 pt-2">
            <textarea
              v-model="replyInput"
              rows="3"
              placeholder="اكتب توضيحك وإجابتك العلمية للطالب هنا..."
              class="w-full px-4 py-3 rounded-2xl bg-[#F8F8F8] dark:bg-[#0F0F0F] border border-[#E5E7EB] dark:border-[#27272A] text-xs font-medium text-[#111111] dark:text-[#FFFFFF] focus:outline-none focus:ring-2 focus:ring-orange-500"
            ></textarea>

            <div class="flex gap-2">
              <button
                type="button"
                @click="submitReply(q)"
                class="px-5 py-2.5 rounded-xl bg-orange-500 hover:bg-[#EA580C] text-white font-bold text-xs flex items-center gap-2 cursor-pointer shadow-sm shadow-orange-500/20"
              >
                <Send class="w-3.5 h-3.5" />
                <span>إرسال الرد للطالب فوراً</span>
              </button>
              <button
                type="button"
                @click="activeReplyId = null"
                class="px-4 py-2.5 rounded-xl bg-[#F8F8F8] dark:bg-[#222222] text-[#6B7280] dark:text-[#A1A1AA] font-bold text-xs cursor-pointer"
              >
                إلغاء
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
