<script setup lang="ts">
import { useRouter } from "vue-router";

const props = withDefaults(
  defineProps<{
    to?: string;
    isActive?: boolean;
    size?: "sm" | "md";
  }>(),
  {
    size: "md",
  }
);

const router = useRouter();

function handleClick() {
  if (props.to) {
    router.push(props.to);
  }
}
</script>

<template>
  <button
    type="button"
    @click="handleClick"
    :class="[
      'flex h-7 min-w-0 items-center gap-2 overflow-hidden rounded-xl px-2 text-xs text-slate-400 hover:bg-slate-800 hover:text-white transition-colors cursor-pointer select-none w-full text-right',
      isActive ? 'font-black text-orange-400 bg-orange-500/10' : ''
    ]"
  >
    <slot />
  </button>
</template>
