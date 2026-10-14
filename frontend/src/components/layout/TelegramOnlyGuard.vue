<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRoute } from "vue-router";
import { ShieldAlert, ExternalLink, MessageSquare, Sparkles, Smartphone } from "@lucide/vue";

const route = useRoute();
const isTelegramEnv = ref<boolean>(true);
const isDevBypass = ref<boolean>(false);

// Admin route or dev mode check
const isAdminRoute = computed(() => route.path.startsWith("/admin"));

onMounted(() => {
  if (typeof window === "undefined") return;

  // Auto-bypass in development (localhost)
  const isLocalhost = window.location.hostname === "localhost" ||
    window.location.hostname === "127.0.0.1" ||
    window.location.hostname === "0.0.0.0";
  if (isLocalhost) {
    isDevBypass.value = true;
    return;
  }

  // Check URL parameters for developer bypass e.g. ?dev=true or #dev
  const urlParams = new URLSearchParams(window.location.search);
  if (urlParams.get("dev") === "true" || window.location.hash.includes("dev")) {
    isDevBypass.value = true;
    return;
  }

  // Telegram WebApp detection.
  // NOTE: the Telegram SDK (telegram-web-app.js) defines a stub
  // Telegram.WebApp even outside Telegram (version 6.0, empty initData).
  // Therefore we require real evidence of being inside Telegram:
  // non-empty initData (signed by Telegram), or initDataUnsafe.user
  // (set only in real Telegram WebViews), or a genuine platform value.
  const SESSION_KEY = "wiyana_tg_env";

  // Bypass: direct login links (?jwt= or ?auth=) prove the user came from
  // the bot — no need for Telegram WebApp context.
  const bypassParams = new URLSearchParams(window.location.search);
  if (bypassParams.get("jwt") || bypassParams.get("auth") || bypassParams.get("dev") === "true") {
    isTelegramEnv.value = true;
    sessionStorage.setItem(SESSION_KEY, "1");
    return;
  }

  // Persist the result across in-app navigations so the guard never
  // flickers between screens (was causing "opens briefly then blocks").
  const stored = sessionStorage.getItem(SESSION_KEY);
  if (stored === "1") {
    isTelegramEnv.value = true;
    return;
  }
  if (stored === "0") {
    isTelegramEnv.value = false;
    return;
  }

  const evaluate = (): boolean => {
    const tg = (window as any).Telegram?.WebApp;
    if (!tg) return false;
    const initData = tg.initData || "";
    const platform = tg.platform || "";
    const unsafe = tg.initDataUnsafe || {};
    const user = unsafe.user;
    const receiver = unsafe.receiver;
    // Real signed initData is the strongest evidence (any Telegram signing format).
    if (initData && initData.includes("hash=") && initData.length > 40) return true;
    // initDataUnsafe.user exists only when Telegram actually injects it.
    if (user && typeof user.id === "number") return true;
    // receiver (bot object) is injected only in real Telegram WebViews.
    if (receiver && receiver.id) return true;
    // Genuine platform values never appear in the external-browser stub.
    const realPlatforms = ["android", "ios", "tdesktop", "weba", "webk", "unigram", "macos", "windows"];
    if (realPlatforms.includes(String(platform).toLowerCase())) return true;
    return false;
  };

  // Evaluate eagerly; if positive, lock it in permanently.
  if (evaluate()) {
    isTelegramEnv.value = true;
    sessionStorage.setItem(SESSION_KEY, "1");
    return;
  }

  // During the first seconds some Telegram clients inject the data
  // progressively (user → receiver → signed initData). Poll regularly
  // for up to 6s. Once positive, we never flip back to false.
  const deadline = Date.now() + 6000;
  const timer = setInterval(() => {
    if (evaluate()) {
      isTelegramEnv.value = true;
      sessionStorage.setItem(SESSION_KEY, "1");
      clearInterval(timer);
      return;
    }
    if (Date.now() > deadline) {
      clearInterval(timer);
      isTelegramEnv.value = false;
      sessionStorage.setItem(SESSION_KEY, "0");
    }
  }, 500);
});

const telegramBotUrl =
  import.meta.env.VITE_TELEGRAM_BOT_URL || "https://t.me/xfile3_bot";
</script>

<template>
  <div v-if="!isTelegramEnv && !isAdminRoute && !isDevBypass" class="fixed inset-0 z-[99999] flex items-center justify-center bg-[#0F0F0F] text-white p-4 sm:p-6 overflow-y-auto selection:bg-orange-500 selection:text-white">
    <!-- Main Guard Card -->
    <div class="relative w-full max-w-lg bg-[#181818] border border-[#27272A] rounded-3xl p-6 sm:p-8 shadow-2xl text-center space-y-6 animate-scale-in">
      
      <!-- Top Icon Badge -->
      <div class="mx-auto w-20 h-20 bg-orange-500 rounded-3xl flex items-center justify-center shadow-lg shadow-orange-500/20 ring-4 ring-orange-500/10">
        <Smartphone class="w-10 h-10 text-white" />
      </div>

      <!-- Header Title & Badge -->
      <div class="space-y-2">
        <div class="inline-flex items-center gap-1.5 px-3.5 py-1 rounded-full text-xs font-bold bg-orange-500/10 text-orange-400 border border-orange-500/20">
          <Sparkles class="w-3.5 h-3.5" />
          <span>Telegram Mini App</span>
        </div>
        <h1 class="text-2xl sm:text-3xl font-black tracking-tight text-white">
          منصة ويانا التعليمية
        </h1>
      </div>

      <!-- Description Body -->
      <div class="bg-[#0F0F0F] rounded-2xl p-5 border border-[#27272A] text-xs sm:text-sm text-[#A1A1AA] leading-relaxed text-right space-y-3">
        <div class="flex items-start gap-3">
          <div class="p-2 rounded-xl bg-orange-500/10 text-orange-400 shrink-0 mt-0.5">
            <ShieldAlert class="w-5 h-5" />
          </div>
          <div>
            <p class="font-bold text-white mb-1">تنبيه الوصول:</p>
            <p class="text-[#A1A1AA]">
              هذه المنصة مخصصة للعمل بكامل ميزاتها حصرياً داخل تطبيق تليجرام (Telegram Mini App). لا يمكن فتح الدروس أو تصفح المحتوى عبر المتصفح الخارجي.
            </p>
          </div>
        </div>
        
        <div class="border-t border-[#27272A] pt-3 text-xs text-[#A1A1AA]">
          للدخول إلى حسابك ومساقاتك، يرجى فتح البوت الرسمي على تليجرام والضغط على زر فتح المنصة.
        </div>
      </div>

      <!-- CTA Button to Open Telegram Bot -->
      <div class="space-y-3 pt-2">
        <a
          :href="telegramBotUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="flex items-center justify-center gap-2.5 w-full py-4 px-6 rounded-2xl bg-orange-500 hover:bg-[#EA580C] text-white font-bold text-sm sm:text-base shadow-lg shadow-orange-500/20 hover:shadow-orange-500/35 transition-all duration-200 active:scale-[0.98]"
        >
          <MessageSquare class="w-5 h-5" />
          <span>الانتقال إلى بوت تليجرام</span>
          <ExternalLink class="w-4 h-4 opacity-80" />
        </a>


      </div>

    </div>
  </div>
</template>
