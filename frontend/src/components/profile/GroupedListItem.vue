<script setup lang="ts">
import type { Component } from "vue";
import { ChevronLeft } from "@lucide/vue";

const props = withDefaults(
  defineProps<{
    title: string;
    description?: string;
    icon?: Component;
    iconColorClass?: string;
    iconBgClass?: string;
    valueText?: string;
    badgeText?: string;
    badgeColorClass?: string;
    showChevron?: boolean;
    isClickable?: boolean;
  }>(),
  {
    showChevron: true,
    isClickable: true,
    iconColorClass: "text-orange-500",
    iconBgClass: "bg-orange-500/10 dark:bg-orange-500/15",
  }
);

const emit = defineEmits<{
  (e: "click"): void;
}>();

function handleClick(event: MouseEvent) {
  if (props.isClickable) {
    emit("click");
  }
}
</script>

<template>
  <div
    @click="handleClick"
    :class="[
      'w-full flex items-center justify-between p-3.5 sm:p-4 text-right transition-colors select-none text-start',
      isClickable
        ? 'hover:bg-[#F8F8F8] dark:hover:bg-[#222222] active:bg-[#E5E7EB]/50 dark:active:bg-[#2A2A2A] cursor-pointer'
        : 'cursor-default'
    ]"
    role="button"
    :tabindex="isClickable ? 0 : undefined"
  >
    <!-- Left/Leading Section: Icon + Title + Description -->
    <div class="flex items-center gap-3.5 min-w-0 flex-1">
      <!-- iOS Rounded Icon Container -->
      <div
        v-if="icon"
        :class="[
          'w-9 h-9 sm:w-10 sm:h-10 rounded-xl sm:rounded-2xl flex items-center justify-center shrink-0 transition-transform duration-200',
          iconBgClass,
          iconColorClass
        ]"
      >
        <component :is="icon" class="w-4.5 h-4.5 sm:w-5 sm:h-5 stroke-[2.2]" />
      </div>

      <div class="flex flex-col min-w-0 flex-1">
        <span class="text-xs sm:text-sm font-bold text-[#111111] dark:text-[#FFFFFF] tracking-tight truncate">
          {{ title }}
        </span>
        <span
          v-if="description"
          class="text-[11px] sm:text-xs text-[#6B7280] dark:text-[#A1A1AA] font-normal leading-normal truncate mt-0.5"
        >
          {{ description }}
        </span>
      </div>
    </div>

    <!-- Right/Trailing Section: Value / Badge / Custom Slot / Chevron -->
    <div class="flex items-center gap-2 shrink-0 mr-2">
      <!-- Slot for custom controls like switch or dropdown -->
      <slot name="trailing">
        <span
          v-if="badgeText"
          :class="[
            'px-2.5 py-0.5 rounded-full text-[10px] font-bold',
            badgeColorClass || 'bg-orange-500/10 text-orange-500'
          ]"
        >
          {{ badgeText }}
        </span>

        <span
          v-if="valueText"
          class="text-xs font-semibold text-[#6B7280] dark:text-[#A1A1AA]"
        >
          {{ valueText }}
        </span>

        <ChevronLeft
          v-if="showChevron"
          class="w-4 h-4 text-[#9CA3AF] dark:text-[#71717A] transition-colors"
        />
      </slot>
    </div>
  </div>
</template>
