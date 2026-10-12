<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import { useSidebar } from "./context";

const props = withDefaults(
  defineProps<{
    isActive?: boolean;
    to?: string;
    tooltip?: string;
    variant?: "default" | "outline";
    size?: "default" | "sm" | "lg";
    asChild?: boolean;
    className?: string;
  }>(),
  {
    isActive: false,
    variant: "default",
    size: "default",
    asChild: false,
  }
);

const { state, isMobile, setOpenMobile } = useSidebar();
const router = useRouter();

function handleClick() {
  if (isMobile.value) {
    setOpenMobile(false);
  }
  if (props.to) {
    router.push(props.to);
  }
}

const sizeClasses = computed(() => {
  switch (props.size) {
    case "sm":
      return "h-8 text-xs px-2.5";
    case "lg":
      return "h-12 text-sm px-3.5";
    default:
      return "h-10 text-xs px-3";
  }
});
</script>

<template>
  <button
    type="button"
    @click="handleClick"
    :title="state === 'collapsed' ? tooltip : undefined"
    :class="[
      'peer/menu-button flex items-center justify-between overflow-hidden rounded-2xl font-bold outline-none transition-all duration-150 group/button cursor-pointer select-none',
      sizeClasses,
      state === 'collapsed' ? 'w-10 h-10 p-0 mx-auto justify-center' : 'w-full',
      isActive
        ? 'bg-gradient-to-r from-orange-600 to-amber-600 text-white shadow-lg shadow-orange-600/25 ring-1 ring-orange-400'
        : 'text-slate-300 hover:text-white hover:bg-slate-800/80 active:scale-[0.98]',
      variant === 'outline' ? 'border border-slate-700 bg-transparent' : '',
      className
    ]"
  >
    <slot />
  </button>
</template>

