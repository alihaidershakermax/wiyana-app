<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Menu,
  Bell,
  User,
  ChevronDown,
  LogOut,
  Settings,
  ChevronLeft,
} from '@lucide/vue'
import { useAdminStore } from '@/admin/stores/useAdminStore'
import { useSidebarState } from '@/admin/composables/useSidebarState'
import ConfirmLogoutDialog from '@/admin/components/ConfirmLogoutDialog.vue'

const route = useRoute()
const router = useRouter()
const adminStore = useAdminStore()
const { toggleSidebar, openMobileSidebar } = useSidebarState()

// Breadcrumb Generation
const breadcrumbs = computed(() => {
  const pathSegments = route.path.split('/').filter(Boolean)
  const items: Array<{ label: string; route?: string }> = [{ label: 'الرئيسية', route: '/admin/dashboard' }]

  if (pathSegments.includes('students')) {
    items.push({ label: 'الطلاب', route: '/admin/students' })
    if (pathSegments.length > 2) items.push({ label: 'تفاصيل الطالب' })
  } else if (pathSegments.includes('courses')) {
    items.push({ label: 'المساقات', route: '/admin/courses' })
    if (pathSegments.includes('create')) items.push({ label: 'إنشاء مساق' })
    else if (pathSegments.length > 2) items.push({ label: 'تفاصيل المساق' })
  } else if (pathSegments.includes('subscription-codes')) {
    items.push({ label: 'أكواد الاشتراك' })
  } else if (pathSegments.includes('subscriptions')) {
    items.push({ label: 'الاشتراكات' })
  } else if (pathSegments.includes('analytics')) {
    items.push({ label: 'التحليلات' })
  } else if (pathSegments.includes('notifications')) {
    items.push({ label: 'الإشعارات' })
  } else if (pathSegments.includes('settings')) {
    items.push({ label: 'إعدادات المنصة' })
  } else if (pathSegments.includes('audit-logs')) {
    items.push({ label: 'سجل الأمان' })
  }

  return items
})

// User Menu Dropdown
const showUserMenu = ref(false)

function handleUserMenuClick(targetRoute: string) {
  showUserMenu.value = false
  router.push(targetRoute)
}

// Logout Dialog
const showLogoutDialog = ref(false)
const loggingOut = ref(false)

async function handleLogout() {
  if (loggingOut.value) return
  loggingOut.value = true
  try {
    await adminStore.logout()
    router.replace('/admin/login')
  } finally {
    loggingOut.value = false
    showLogoutDialog.value = false
  }
}
</script>

<template>
  <header class="glass-liquid-nav h-16 bg-[#0A0A0A]/80 border-b border-[#242424] px-4 lg:px-6 flex items-center justify-between sticky top-0 z-20 shrink-0 select-none">
    <!-- Left Group: Sidebar Toggle & Breadcrumbs -->
    <div class="flex items-center gap-3">
      <!-- Desktop Toggle Button -->
      <button
        class="hidden lg:flex p-2 rounded-xl text-[#A1A1AA] hover:text-white hover:bg-[#181818] transition-colors"
        title="طي/توسيع القائمة"
        @click="toggleSidebar"
      >
        <Menu class="h-5 w-5" />
      </button>

      <!-- Mobile Toggle Button -->
      <button
        class="lg:hidden p-2 rounded-xl text-[#A1A1AA] hover:text-white hover:bg-[#181818] transition-colors"
        title="فتح القائمة"
        @click="openMobileSidebar"
      >
        <Menu class="h-5 w-5" />
      </button>

      <!-- Breadcrumbs Navigation -->
      <nav class="hidden sm:flex items-center gap-1.5 text-xs text-[#71717A] font-medium">
        <template v-for="(b, idx) in breadcrumbs" :key="idx">
          <ChevronLeft v-if="idx > 0" class="h-3.5 w-3.5 text-[#333333]" />
          <button
            v-if="b.route && idx < breadcrumbs.length - 1"
            type="button"
            @click="router.push(b.route)"
            class="hover:text-[#F97316] transition-colors cursor-pointer"
          >
            {{ b.label }}
          </button>
          <span v-else :class="idx === breadcrumbs.length - 1 ? 'font-bold text-white' : ''">
            {{ b.label }}
          </span>
        </template>
      </nav>
    </div>

    <!-- Right Group: Actions (Notifications, User Menu) -->
    <div class="flex items-center gap-2">
      <!-- Notifications Icon Button -->
      <button
        type="button"
        @click="router.push('/admin/notifications')"
        class="p-2 rounded-xl text-[#A1A1AA] hover:text-white hover:bg-[#181818] transition-colors relative cursor-pointer"
        title="الإشعارات"
      >
        <Bell class="h-5 w-5" />
        <span class="absolute top-1.5 right-1.5 h-2 w-2 rounded-full bg-[#F97316] ring-2 ring-[#0A0A0A]" />
      </button>

      <!-- User Profile Dropdown -->
      <div class="relative">
        <button
          class="flex items-center gap-2 p-1.5 rounded-xl hover:bg-[#181818] transition-colors"
          @click="showUserMenu = !showUserMenu"
        >
          <div class="h-8 w-8 rounded-full bg-[#F97316] text-white flex items-center justify-center font-bold text-xs">
            <User class="h-4 w-4" />
          </div>
          <span class="hidden md:inline-block text-xs font-extrabold text-[#F5F5F5]">
            {{ adminStore.userName }}
          </span>
          <ChevronDown class="h-3.5 w-3.5 text-[#71717A]" />
        </button>

        <!-- Dropdown Menu -->
        <div
          v-if="showUserMenu"
          class="absolute left-0 mt-2 w-48 bg-[#111111] border border-[#242424] rounded-2xl shadow-2xl py-2 z-50 animate-slide-down"
        >
          <div class="px-4 py-2 border-b border-[#242424]">
            <p class="text-xs font-extrabold text-white truncate">{{ adminStore.userName }}</p>
            <p class="text-[10px] text-[#A1A1AA] capitalize">{{ adminStore.userRole }}</p>
          </div>

          <button
            class="w-full px-4 py-2 text-right text-xs font-semibold hover:bg-[#181818] flex items-center gap-2 text-[#F5F5F5]"
            @click="handleUserMenuClick('/admin/settings')"
          >
            <Settings class="h-4 w-4 text-[#71717A]" />
            الإعدادات
          </button>

          <button
            class="w-full px-4 py-2 text-right text-xs font-semibold hover:bg-red-950/40 flex items-center gap-2 text-red-400"
            @click="showLogoutDialog = true; showUserMenu = false"
          >
            <LogOut class="h-4 w-4" />
            تسجيل الخروج
          </button>
        </div>
      </div>
    </div>

    <!-- Confirm Logout Dialog -->
    <ConfirmLogoutDialog
      :open="showLogoutDialog"
      :loading="loggingOut"
      @confirm="handleLogout"
      @cancel="showLogoutDialog = false"
    />
  </header>
</template>
