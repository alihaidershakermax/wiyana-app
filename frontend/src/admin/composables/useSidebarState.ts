import { ref, watch, onMounted } from 'vue'

const STORAGE_KEY = 'wiyana_admin_sidebar_collapsed'

const isCollapsed = ref<boolean>(false)
const isMobileOpen = ref<boolean>(false)

// Read initial collapsed state from localStorage if client side
if (typeof window !== 'undefined') {
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved !== null) {
    isCollapsed.value = saved === 'true'
  }
}

watch(isCollapsed, (newVal) => {
  if (typeof window !== 'undefined') {
    localStorage.setItem(STORAGE_KEY, String(newVal))
  }
})

export function useSidebarState() {
  function toggleSidebar() {
    isCollapsed.value = !isCollapsed.value
  }

  function setCollapsed(val: boolean) {
    isCollapsed.value = val
  }

  function toggleMobileSidebar() {
    isMobileOpen.value = !isMobileOpen.value
  }

  function closeMobileSidebar() {
    isMobileOpen.value = false
  }

  function openMobileSidebar() {
    isMobileOpen.value = true
  }

  return {
    isCollapsed,
    isMobileOpen,
    toggleSidebar,
    setCollapsed,
    toggleMobileSidebar,
    closeMobileSidebar,
    openMobileSidebar,
  }
}
