<script setup lang="ts">
defineProps<{
  open: boolean
  title?: string
  message?: string
  confirmLabel?: string
  cancelLabel?: string
  variant?: 'default' | 'danger'
  loading?: boolean
}>()

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-all duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-all duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="open" class="fixed inset-0 z-[9998] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" dir="rtl">
        <Transition
          enter-active-class="transition-all duration-200 ease-out"
          enter-from-class="opacity-0 scale-95 translate-y-2"
          enter-to-class="opacity-100 scale-100 translate-y-0"
          appear
        >
          <div class="card w-full max-w-sm p-6 shadow-2xl bg-white dark:bg-slate-900" @click.stop>
            <h3 class="text-base font-bold mb-2">{{ title ?? 'تأكيد العملية' }}</h3>
            <p class="text-sm text-slate-500 dark:text-slate-400 leading-relaxed mb-6">
              {{ message ?? 'هل أنت متأكد من تنفيذ هذه العملية؟' }}
            </p>
            <div class="flex gap-3">
              <button
                :disabled="loading"
                @click="emit('confirm')"
                :class="[
                  'flex-1 btn-primary py-2.5',
                  variant === 'danger' ? 'bg-red-500 hover:bg-red-600' : '',
                  loading ? 'opacity-60 cursor-not-allowed' : '',
                ]"
              >
                <span v-if="loading" class="animate-spin inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full ml-2" />
                {{ confirmLabel ?? 'تأكيد' }}
              </button>
              <button
                :disabled="loading"
                @click="emit('cancel')"
                class="flex-1 btn-ghost border border-slate-200 dark:border-slate-700 py-2.5"
              >
                {{ cancelLabel ?? 'إلغاء' }}
              </button>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>
