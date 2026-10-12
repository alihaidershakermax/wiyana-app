<script setup lang="ts">
withDefaults(
  defineProps<{
    size?: "sm" | "md" | "lg";
    showLabel?: boolean;
    label?: string;
    subtitle?: string;
  }>(),
  {
    size: "md",
    showLabel: true,
    label: "ويانا",
    subtitle: "منصة التعلم",
  },
);
</script>

<template>
  <div class="brand-logo" :class="`brand-logo--${size}`" dir="rtl">
    <span class="brand-logo__mark" aria-hidden="true">
      <svg viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path
          d="M32 5.5 55 14v16.8c0 13.8-9.6 23.4-23 27.7C18.6 54.2 9 44.6 9 30.8V14L32 5.5Z"
          fill="currentColor"
          fill-opacity=".16"
          stroke="currentColor"
          stroke-width="3"
        />
        <path d="m16 23 16-7 16 7-16 7-16-7Z" fill="currentColor" />
        <path d="M23 28v7.2c4.9 3.2 13.1 3.2 18 0V28l-9 4-9-4Z" fill="currentColor" fill-opacity=".88" />
        <path d="M49 24v9" stroke="currentColor" stroke-width="3" stroke-linecap="round" />
      </svg>
    </span>
    <span v-if="showLabel" class="brand-logo__copy">
      <strong>{{ label }}</strong>
      <small>{{ subtitle }}</small>
    </span>
  </div>
</template>

<style scoped>
.brand-logo {
  display: inline-flex;
  align-items: center;
  gap: 0.7rem;
  min-width: 0;
  color: #f97316;
}

.brand-logo__mark {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 auto;
  border: 1px solid rgba(249, 115, 22, 0.36);
  border-radius: 0.95rem;
  background: linear-gradient(145deg, rgba(249, 115, 22, 0.24), rgba(249, 115, 22, 0.07));
  box-shadow: 0 9px 24px rgba(249, 115, 22, 0.16), inset 0 1px rgba(255, 255, 255, 0.16);
  backdrop-filter: blur(18px) saturate(145%);
}

.brand-logo__mark svg {
  width: 62%;
  height: 62%;
}

.brand-logo__copy {
  display: flex;
  min-width: 0;
  flex-direction: column;
  line-height: 1.15;
}

.brand-logo__copy strong {
  overflow: hidden;
  color: #fff;
  font-size: 0.98rem;
  font-weight: 900;
  letter-spacing: -0.02em;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.brand-logo__copy small {
  margin-top: 0.22rem;
  overflow: hidden;
  color: rgba(249, 115, 22, 0.88);
  font-size: 0.65rem;
  font-weight: 800;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.brand-logo--sm .brand-logo__mark { width: 2.25rem; height: 2.25rem; border-radius: 0.7rem; }
.brand-logo--md .brand-logo__mark { width: 2.75rem; height: 2.75rem; }
.brand-logo--lg .brand-logo__mark { width: 4.25rem; height: 4.25rem; border-radius: 1.25rem; }
.brand-logo--lg .brand-logo__copy strong { font-size: 1.35rem; }
.brand-logo--lg .brand-logo__copy small { font-size: 0.72rem; }

@media (prefers-reduced-motion: reduce) {
  .brand-logo__mark { transition: none; }
}
</style>
