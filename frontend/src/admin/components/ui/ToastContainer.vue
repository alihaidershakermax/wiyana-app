<script setup lang="ts">
import { useToast } from '@/admin/composables/useToast'
import { CheckCircle, XCircle, AlertTriangle, Info, X } from '@lucide/vue'
import { TransitionGroup } from 'vue'

const { toasts, remove } = useToast()

const icons = {
  success: CheckCircle,
  error: XCircle,
  warning: AlertTriangle,
  info: Info,
}

const styles = {
  success: 'border-emerald-200 bg-emerald-50 dark:border-emerald-900/60 dark:bg-emerald-950/40 text-emerald-800 dark:text-emerald-300',
  error: 'border-red-200 bg-red-50 dark:border-red-900/60 dark:bg-red-950/40 text-red-800 dark:text-red-300',
  warning: 'border-amber-200 bg-amber-50 dark:border-amber-900/60 dark:bg-amber-950/40 text-amber-800 dark:text-amber-300',
  info: 'border-blue-200 bg-blue-50 dark:border-blue-900/60 dark:bg-blue-950/40 text-blue-800 dark:text-blue-300',
}

const iconStyles = {
  success: 'text-emerald-500',
  error: 'text-red-500',
  warning: 'text-amber-500',
  info: 'text-blue-500',
}
</script>

<template>
  <Teleport to="body">
    <div
      class="fixed bottom-6 left-6 z-[9999] flex flex-col gap-3 w-80 max-w-[calc(100vw-2rem)]"
      dir="rtl"
      aria-live="polite"
    >
      <TransitionGroup
        enter-active-class="transition-all duration-300 ease-out"
        enter-from-class="opacity-0 translate-y-4 scale-95"
        enter-to-class="opacity-100 translate-y-0 scale-100"
        leave-active-class="transition-all duration-200 ease-in absolute w-full"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0 translate-x-full"
        move-class="transition-all duration-300"
      >
        <div
          v-for="toast in toasts"
          :key="toast.id"
          :class="['card border flex items-start gap-3 p-4 shadow-lg', styles[toast.type]]"
        >
          <component
            :is="icons[toast.type]"
            :class="['h-5 w-5 mt-0.5 shrink-0', iconStyles[toast.type]]"
          />
          <div class="flex-1 min-w-0">
            <p class="text-sm font-bold leading-snug">{{ toast.title }}</p>
            <p v-if="toast.message" class="text-xs mt-0.5 opacity-80">{{ toast.message }}</p>
          </div>
          <button
            @click="remove(toast.id)"
            class="shrink-0 p-0.5 rounded opacity-60 hover:opacity-100 transition-opacity"
          >
            <X class="h-4 w-4" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>
