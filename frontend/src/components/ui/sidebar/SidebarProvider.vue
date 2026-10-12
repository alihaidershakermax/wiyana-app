<script setup lang="ts">
import { ref, computed, provide, onMounted, onUnmounted, watch } from "vue";
import { SIDEBAR_INJECTION_KEY, type SidebarContext } from "./context";

const SIDEBAR_COOKIE_NAME = "sidebar_state";
const SIDEBAR_KEYBOARD_SHORTCUT = "b";

const props = withDefaults(
  defineProps<{
    defaultOpen?: boolean;
    open?: boolean;
    storageKey?: string;
  }>(),
  {
    defaultOpen: true,
    storageKey: SIDEBAR_COOKIE_NAME,
  }
);

const emit = defineEmits<{
  (e: "update:open", value: boolean): void;
}>();

// Get initial state from localStorage if available
function getInitialState(): boolean {
  if (typeof window === "undefined") return props.defaultOpen;
  try {
    const stored = localStorage.getItem(props.storageKey);
    if (stored !== null) {
      return stored === "true";
    }
  } catch (e) {
    // Ignore storage errors
  }
  return props.defaultOpen;
}

const isOpenInternal = ref(getInitialState());
const isMobile = ref(typeof window !== "undefined" ? window.innerWidth < 1024 : false);
const openMobile = ref(false);

const open = computed({
  get: () => (props.open !== undefined ? props.open : isOpenInternal.value),
  set: (val: boolean) => {
    isOpenInternal.value = val;
    emit("update:open", val);
    if (typeof window !== "undefined") {
      try {
        localStorage.setItem(props.storageKey, String(val));
      } catch (e) {
        // Ignore storage errors
      }
    }
  },
});

const state = computed(() => (open.value ? "expanded" : "collapsed"));

function toggleSidebar() {
  if (isMobile.value) {
    openMobile.value = !openMobile.value;
  } else {
    open.value = !open.value;
  }
}

function handleResize() {
  const mobile = window.innerWidth < 1024;
  isMobile.value = mobile;
  if (!mobile) {
    openMobile.value = false;
  }
}

function handleKeyDown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === SIDEBAR_KEYBOARD_SHORTCUT) {
    e.preventDefault();
    toggleSidebar();
  }
}

onMounted(() => {
  handleResize();
  window.addEventListener("resize", handleResize);
  window.addEventListener("keydown", handleKeyDown);
});

onUnmounted(() => {
  window.removeEventListener("resize", handleResize);
  window.removeEventListener("keydown", handleKeyDown);
});

const context: SidebarContext = {
  open,
  setOpen: (val: boolean) => {
    open.value = val;
  },
  openMobile,
  setOpenMobile: (val: boolean) => {
    openMobile.value = val;
  },
  isMobile,
  toggleSidebar,
  state,
};

provide(SIDEBAR_INJECTION_KEY, context);
</script>

<template>
  <div
    class="group/sidebar-wrapper flex min-h-screen w-full bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 transition-colors duration-200"
    :data-state="state"
    :style="{
      '--sidebar-width': '17.5rem',
      '--sidebar-width-icon': '4.5rem',
      '--sidebar-width-mobile': '18rem',
    }"
  >
    <slot />
  </div>
</template>

