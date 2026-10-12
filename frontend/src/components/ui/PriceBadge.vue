<script setup lang="ts">
import { computed } from "vue";
import { CheckCircle2 } from "@lucide/vue";
import { CourseEnrollmentStatus } from "@/types/course";

const props = withDefaults(
  defineProps<{
    price: number;
    currency?: string;
    isFree?: boolean;
    enrollmentStatus?: CourseEnrollmentStatus;
  }>(),
  {
    currency: "د.ع",
    isFree: false,
    enrollmentStatus: "not_enrolled",
  }
);

const formattedPrice = computed(() => {
  if (props.isFree || props.price === 0) {
    return "مجاني";
  }
  return `${props.price.toLocaleString("ar-IQ")} ${props.currency}`;
});
</script>

<template>
  <div>
    <!-- Completed Badge -->
    <span
      v-if="enrollmentStatus === 'completed'"
      class="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg text-xs font-bold bg-orange-500/10 border border-orange-500/30 text-orange-500 dark:text-orange-400"
    >
      <CheckCircle2 class="w-3.5 h-3.5" />
      <span>مكتمل</span>
    </span>

    <!-- Subscribed / In Progress Badge -->
    <span
      v-else-if="enrollmentStatus === 'enrolled'"
      class="inline-flex items-center px-2.5 py-1 rounded-lg text-xs font-bold bg-orange-500/10 border border-orange-500/30 text-orange-500 dark:text-orange-400"
    >
      مشترك
    </span>

    <!-- Price / Free Badge -->
    <span
      v-else
      :class="[
        'inline-flex items-center px-2.5 py-1 rounded-lg text-xs font-bold shadow-sm',
        isFree || price === 0
          ? 'bg-orange-500 text-white'
          : 'bg-orange-500 text-white'
      ]"
    >
      {{ formattedPrice }}
    </span>
  </div>
</template>
