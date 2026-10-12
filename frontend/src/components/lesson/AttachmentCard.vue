<script setup lang="ts">
import { computed, ref } from "vue";
import { FileText, File, Image as ImageIcon, Link2, Download, Loader2 } from "@lucide/vue";
import type { LessonAttachment } from "@/types/lesson";

const props = defineProps<{
  attachment: LessonAttachment;
}>();

const isLink = computed(() => props.attachment.type === "link");
const isTelegramFile = computed(() => props.attachment.type === "telegram_file");

const sizeText = computed(() => {
  const bytes = props.attachment.sizeBytes;
  if (!bytes || bytes <= 0) return null;
  if (bytes < 1024 * 1024) return `${Math.round(bytes / 1024)} ك.ب`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} م.ب`;
});

const downloading = ref(false);

const TELEGRAM_BOT_USERNAME = "WenaAppBot";

async function openAttachment() {
  // External links are blocked: no link is exposed to students.
  if (isLink.value) return;

  // Telegram-hosted files open through the bot, never a raw URL.
  const telegramFileId = (props.attachment as any).telegram_file_id;
  if (isTelegramFile.value && telegramFileId) {
    window.open(`https://t.me/${TELEGRAM_BOT_USERNAME}?start=file_${telegramFileId}`, "_blank");
    return;
  }

  if (!props.attachment.url || downloading.value) return;
  downloading.value = true;
  try {
    // Fetch the file as a blob so the real storage URL is never shown.
    const res = await fetch(props.attachment.url, { credentials: "same-origin" });
    if (!res.ok) throw new Error("fetch failed");
    const blob = await res.blob();
    const objectUrl = URL.createObjectURL(blob);
    if (props.attachment.type === "image" || props.attachment.type === "pdf") {
      window.open(objectUrl, "_blank");
    } else {
      const a = document.createElement("a");
      a.href = objectUrl;
      a.download = props.attachment.title || "ملف";
      document.body.appendChild(a);
      a.click();
      a.remove();
    }
    setTimeout(() => URL.revokeObjectURL(objectUrl), 60000);
  } catch {
    // Do not fall back to revealing the raw URL.
  } finally {
    downloading.value = false;
  }
}
</script>

<template>
  <div class="flex items-center justify-between gap-3 p-3.5 bg-white dark:bg-[#181818] border border-[#E5E7EB] dark:border-[#27272A] rounded-2xl">
    <div class="flex items-center gap-3 min-w-0">
      <div
        class="w-10 h-10 rounded-xl bg-orange-500/10 text-orange-500 flex items-center justify-center shrink-0"
      >
        <FileText v-if="attachment.type === 'pdf'" class="w-5 h-5" />
        <ImageIcon v-else-if="attachment.type === 'image'" class="w-5 h-5" />
        <Link2 v-else-if="attachment.type === 'link'" class="w-5 h-5" />
        <File v-else class="w-5 h-5" />
      </div>

      <div class="flex flex-col gap-0.5 min-w-0">
        <span class="text-xs font-bold text-[#111111] dark:text-[#FFFFFF] truncate">
          {{ attachment.title }}
        </span>
        <span v-if="sizeText" class="text-[10px] text-[#6B7280] dark:text-[#A1A1AA]">{{ sizeText }}</span>
      </div>
    </div>

    <button
      v-if="!isLink"
      @click="openAttachment"
      type="button"
      :disabled="downloading"
      :aria-label="isTelegramFile ? 'فتح عبر تيلجرام' : 'تحميل الملف'"
      class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl bg-orange-500 hover:bg-[#EA580C] text-white font-bold text-[11px] shadow-sm shadow-orange-500/20 active:scale-95 transition-all shrink-0 cursor-pointer disabled:opacity-60 disabled:cursor-wait"
    >
      <Loader2 v-if="downloading" class="w-3.5 h-3.5 animate-spin" />
      <Download v-else class="w-3.5 h-3.5" />
      <span>{{ isTelegramFile ? "فتح" : "تنزيل" }}</span>
    </button>
  </div>
</template>
