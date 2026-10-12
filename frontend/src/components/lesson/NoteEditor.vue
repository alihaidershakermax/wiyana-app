<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { Check, Loader2, Trash2, StickyNote } from "@lucide/vue";
import { notesService } from "@/services/notesService";
import { formatTimeOfDay } from "@/utils/format";

const props = defineProps<{
  lessonId: string;
}>();

type SaveStatus = "idle" | "saving" | "saved";

const text = ref("");
const status = ref<SaveStatus>("idle");
const lastSavedAt = ref("");
let saveTimer: ReturnType<typeof setTimeout> | null = null;

const wordCount = computed(() => {
  const words = text.value.trim().split(/\s+/).filter(Boolean);
  return words.length;
});

const statusLabel = computed(() => {
  if (status.value === "saving") return "جارٍ الحفظ...";
  if (status.value === "saved") return "تم الحفظ";
  return "ملاحظاتك تُحفظ تلقائياً";
});

const statusTime = computed(() => formatTimeOfDay(lastSavedAt.value));

function saveNow() {
  const note = notesService.saveNote(props.lessonId, text.value);
  lastSavedAt.value = note.updatedAt;
  status.value = "saved";
}

function scheduleSave() {
  if (saveTimer) clearTimeout(saveTimer);
  status.value = "saving";
  saveTimer = setTimeout(() => {
    saveNow();
    saveTimer = null;
  }, 600);
}

function clearNotes() {
  if (saveTimer) {
    clearTimeout(saveTimer);
    saveTimer = null;
  }
  text.value = "";
  notesService.clearNote(props.lessonId);
  lastSavedAt.value = "";
  status.value = "saved";
}

watch(text, () => {
  scheduleSave();
});

onMounted(() => {
  const note = notesService.getNote(props.lessonId);
  text.value = note.text;
  if (note.updatedAt) lastSavedAt.value = note.updatedAt;
});

onBeforeUnmount(() => {
  if (saveTimer) {
    clearTimeout(saveTimer);
  }
  saveNow();
});
</script>

<template>
  <div class="bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-2xl overflow-hidden text-start">
    <div class="flex items-center justify-between px-4 pt-3.5">
      <span class="inline-flex items-center gap-1.5 text-xs font-bold text-[#111111] dark:text-[#FFFFFF]">
        <StickyNote class="w-4 h-4 text-orange-500" />
        ملاحظاتك على المحاضرة
      </span>

      <button
        @click="clearNotes"
        type="button"
        :disabled="!text"
        class="inline-flex items-center gap-1 px-2.5 py-1.5 rounded-lg text-[10px] font-bold text-rose-500 hover:bg-rose-500/10 transition-colors disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
      >
        <Trash2 class="w-3.5 h-3.5" />
        مسح
      </button>
    </div>

    <div class="p-4">
      <textarea
        v-model="text"
        rows="6"
        placeholder="اكتب ملاحظاتك أو أفكارك حول هذه المحاضرة هنا... سيتم الحفظ تلقائياً."
        class="w-full bg-[#F8F8F8] dark:bg-[#222222] border border-[#E5E7EB] dark:border-[#27272A] rounded-xl p-3.5 text-xs sm:text-sm text-[#111111] dark:text-[#FFFFFF] placeholder:text-[#9CA3AF] dark:placeholder:text-[#71717A] focus:border-orange-500 focus:outline-none focus:ring-2 focus:ring-orange-500/20 transition-all resize-none"
      ></textarea>

      <div class="flex items-center justify-between mt-2.5">
        <span class="inline-flex items-center gap-1.5 text-[10px] font-bold text-[#6B7280] dark:text-[#A1A1AA]">
          <Check v-if="status === 'saved'" class="w-3.5 h-3.5 text-emerald-500" />
          <Loader2 v-else-if="status === 'saving'" class="w-3.5 h-3.5 text-orange-500 animate-spin" />
          <span>{{ statusLabel }}</span>
          <span v-if="status === 'saved' && statusTime">· {{ statusTime }}</span>
        </span>
        <span class="text-[10px] font-bold text-[#6B7280] dark:text-[#A1A1AA] tabular-nums">{{ wordCount }} كلمة</span>
      </div>
    </div>
  </div>
</template>
