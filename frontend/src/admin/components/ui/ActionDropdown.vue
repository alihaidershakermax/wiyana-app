<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { MoreHorizontal } from '@lucide/vue'

export interface MenuItem {
  id: string
  label: string
  icon?: any
  danger?: boolean
  disabled?: boolean
  onClick: () => void
}

const props = defineProps<{
  items: MenuItem[]
}>()

const isOpen = ref(false)
const buttonRef = ref<HTMLButtonElement | null>(null)
const menuRef = ref<HTMLDivElement | null>(null)

const menuStyle = ref<{ top: string; left: string }>({ top: '0px', left: '0px' })

function toggleMenu(e: MouseEvent) {
  e.stopPropagation()
  if (isOpen.value) {
    isOpen.value = false
    return
  }
  isOpen.value = true

  nextTick(() => {
    updatePosition()
    // Second tick ensures rendered dimensions are exact
    requestAnimationFrame(updatePosition)
  })
}

function updatePosition() {
  if (!buttonRef.value) return
  const rect = buttonRef.value.getBoundingClientRect()
  const menuWidth = menuRef.value?.offsetWidth || 185
  const menuHeight = menuRef.value?.offsetHeight || 160
  const viewportHeight = window.innerHeight
  const windowWidth = window.innerWidth

  // Align menu right edge with button right edge (RTL standard)
  let left = rect.right - menuWidth
  
  // Strict screen boundary clamping
  if (left < 12) {
    left = Math.max(12, rect.left)
  }
  if (left + menuWidth > windowWidth - 12) {
    left = Math.max(12, windowWidth - menuWidth - 12)
  }

  // Vertical placement: open downwards if space permits, else upwards
  const spaceBelow = viewportHeight - rect.bottom
  const openUpwards = spaceBelow < menuHeight + 10 && rect.top > menuHeight + 10

  let top = openUpwards ? rect.top - menuHeight - 6 : rect.bottom + 6
  if (top < 12) top = 12
  if (top + menuHeight > viewportHeight - 12) {
    top = Math.max(12, viewportHeight - menuHeight - 12)
  }

  menuStyle.value = {
    top: `${Math.round(top)}px`,
    left: `${Math.round(left)}px`,
  }
}

function handleGlobalClick(e: MouseEvent) {
  if (
    buttonRef.value?.contains(e.target as Node) ||
    menuRef.value?.contains(e.target as Node)
  ) {
    return
  }
  isOpen.value = false
}

function handleScrollOrResize() {
  if (isOpen.value) {
    updatePosition()
  }
}

onMounted(() => {
  window.addEventListener('click', handleGlobalClick)
  window.addEventListener('scroll', handleScrollOrResize, true)
  window.addEventListener('resize', handleScrollOrResize)
})

onUnmounted(() => {
  window.removeEventListener('click', handleGlobalClick)
  window.removeEventListener('scroll', handleScrollOrResize, true)
  window.removeEventListener('resize', handleScrollOrResize)
})
</script>

<template>
  <div class="relative inline-block text-right">
    <button
      ref="buttonRef"
      type="button"
      class="p-2 rounded-xl text-[#A1A1AA] hover:text-white hover:bg-[#222222] transition-colors focus:outline-none focus:ring-1 focus:ring-[#F97316] cursor-pointer"
      title="الإجراءات"
      @click="toggleMenu"
    >
      <MoreHorizontal class="h-5 w-5" />
    </button>

    <!-- Teleported Floating Menu clamped within viewport -->
    <Teleport to="body">
      <div
        v-if="isOpen"
        ref="menuRef"
        dir="rtl"
        :style="menuStyle"
        class="fixed z-[9999] min-w-[180px] bg-[#141414] border border-[#2B2B2B] rounded-2xl shadow-2xl py-1.5 animate-dropdown select-none text-right backdrop-blur-md"
        @click.stop
      >
        <button
          v-for="item in items"
          :key="item.id"
          type="button"
          :disabled="item.disabled"
          :class="[
            'w-full px-3.5 py-2.5 text-xs font-bold flex items-center gap-2.5 transition-colors text-right cursor-pointer rounded-xl mx-auto',
            item.danger
              ? 'text-rose-400 hover:bg-rose-950/40 hover:text-rose-300'
              : 'text-[#E4E4E7] hover:bg-[#222222] hover:text-white',
            item.disabled ? 'opacity-40 cursor-not-allowed' : '',
          ]"
          @click="item.onClick(); isOpen = false"
        >
          <component
            :is="item.icon"
            v-if="item.icon"
            class="h-4 w-4 shrink-0"
            :class="item.danger ? 'text-rose-400' : 'text-[#F97316]'"
          />
          <span class="truncate">{{ item.label }}</span>
        </button>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
@keyframes dropdownFade {
  from {
    opacity: 0;
    transform: scale(0.95) translateY(-4px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.animate-dropdown {
  animation: dropdownFade 0.15s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  box-shadow: 0 10px 30px -5px rgba(0, 0, 0, 0.6), 0 0 0 1px rgba(255, 255, 255, 0.05);
}
</style>
