<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ChevronDown, LogOut, UserRound, X } from '@lucide/vue'
import BrandLogo from '@/components/BrandLogo.vue'
import { useAdminStore } from '@/admin/stores/useAdminStore'
import { useSidebarState } from '@/admin/composables/useSidebarState'
import { getNavigationSections, type NavSection, type NavChildItem } from '@/admin/config/navigation'
import ConfirmLogoutDialog from '@/admin/components/ConfirmLogoutDialog.vue'

const route = useRoute()
const router = useRouter()
const adminStore = useAdminStore()
const { isCollapsed, isMobileOpen, closeMobileSidebar } = useSidebarState()

const navigationSections = computed<NavSection[]>(() => {
  const role = adminStore.userRole || 'OWNER'
  return getNavigationSections(role)
})

// Track open state of collapsible sections
const openSections = ref<Record<string, boolean>>({
  content: true,
  users: true,
  subscriptions: true,
  communication: true,
  analytics: true,
  system: true,
  'teacher-content': true,
  'teacher-students': true,
  'teacher-system': true,
})

function toggleSection(sectionId: string) {
  if (isCollapsed.value) return
  openSections.value[sectionId] = !openSections.value[sectionId]
}

function isChildActive(child: NavChildItem): boolean {
  return route.path === child.route
}

function isSectionActive(sec: NavSection): boolean {
  if (sec.route && route.path === sec.route) return true
  if (sec.children) {
    return sec.children.some((child) => route.path === child.route)
  }
  return false
}

function autoExpandActiveSection() {
  for (const sec of navigationSections.value) {
    if (sec.children && sec.children.some((child) => route.path === child.route)) {
      openSections.value[sec.id] = true
    }
  }
}

watch(() => route.path, () => {
  autoExpandActiveSection()
})

onMounted(() => {
  autoExpandActiveSection()
})

function navigateTo(targetRoute: string) {
  router.push(targetRoute)
  closeMobileSidebar()
}

const hoveredSection = ref<string | null>(null)
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
  <div>
    <!-- Desktop Sidebar -->
    <aside
      :class="[
        'glass-liquid-nav hidden lg:flex flex-col h-screen sticky top-0 bg-[#0B0B0B]/90 border-l border-[#242424] transition-all duration-300 ease-in-out z-30 select-none',
        isCollapsed ? 'w-[72px]' : 'w-[260px]',
      ]"
    >
      <!-- Header / Logo -->
      <div class="h-16 flex items-center px-4 border-b border-[#242424] shrink-0 justify-between">
        <BrandLogo :size="isCollapsed ? 'sm' : 'md'" :show-label="!isCollapsed" subtitle="لوحة التحكم" />
      </div>

      <!-- Navigation List -->
      <div class="flex-1 overflow-y-auto overflow-x-hidden py-4 px-3 space-y-1.5 custom-scrollbar">
        <div v-for="sec in navigationSections" :key="sec.id" class="relative group">
          <!-- Case A: Single Item (e.g. الرئيسية) -->
          <div
            v-if="!sec.children"
            :class="[
              'flex items-center gap-3 px-3 py-2.5 rounded-xl cursor-pointer text-sm font-semibold transition-all duration-150 relative',
              isSectionActive(sec)
                ? 'bg-[#1F0E03] text-[#F97316] font-extrabold shadow-sm'
                : 'text-[#A1A1AA] hover:bg-[#161616] hover:text-white',
              isCollapsed ? 'justify-center px-0' : '',
            ]"
            @mouseenter="hoveredSection = sec.id"
            @mouseleave="hoveredSection = null"
            @click="navigateTo(sec.route!)"
          >
            <div v-if="isSectionActive(sec)" class="absolute right-0 top-1.5 bottom-1.5 w-1 bg-[#F97316] rounded-l-full" />
            <component
              :is="sec.icon"
              :class="[
                'h-5 w-5 shrink-0 transition-colors',
                isSectionActive(sec) ? 'text-[#F97316]' : 'text-[#71717A] group-hover:text-white',
              ]"
            />
            <span v-if="!isCollapsed" class="truncate flex-1 text-right text-xs font-bold">{{ sec.label }}</span>
            <!-- Tooltip -->
            <div v-if="isCollapsed && hoveredSection === sec.id" class="absolute right-full top-1/2 -translate-y-1/2 mr-3 px-3 py-1.5 bg-[#111111] text-white border border-[#242424] text-xs font-bold rounded-lg shadow-xl whitespace-nowrap z-50 pointer-events-none animate-fade-in">
              {{ sec.label }}
            </div>
          </div>

          <!-- Case B: Collapsible Section with Children -->
          <div v-else class="space-y-1">
            <!-- Section Header (20px icon) -->
            <div
              :class="[
                'flex items-center gap-3 px-3 py-2.5 rounded-xl cursor-pointer text-sm font-semibold transition-all duration-150 relative select-none',
                isSectionActive(sec)
                  ? 'text-[#F97316] font-extrabold'
                  : 'text-[#A1A1AA] hover:bg-[#161616] hover:text-white',
                isCollapsed ? 'justify-center px-0' : '',
              ]"
              @mouseenter="hoveredSection = sec.id"
              @mouseleave="hoveredSection = null"
              @click="toggleSection(sec.id)"
            >
              <component
                :is="sec.icon"
                :class="[
                  'h-5 w-5 shrink-0 transition-colors',
                  isSectionActive(sec) ? 'text-[#F97316]' : 'text-[#71717A] group-hover:text-white',
                ]"
              />

              <span v-if="!isCollapsed" class="truncate flex-1 text-right text-xs font-bold">
                {{ sec.label }}
              </span>

              <!-- Chevron Arrow -->
              <ChevronDown
                v-if="!isCollapsed"
                :class="['h-4 w-4 shrink-0 transition-transform duration-200 text-[#71717A]', openSections[sec.id] ? 'rotate-180 text-white' : '']"
              />

              <!-- Floating Tooltip for Collapsed Mode -->
              <div v-if="isCollapsed && hoveredSection === sec.id" class="absolute right-full top-1/2 -translate-y-1/2 mr-3 px-3 py-1.5 bg-[#111111] text-white border border-[#242424] text-xs font-bold rounded-lg shadow-xl whitespace-nowrap z-50 pointer-events-none animate-fade-in">
                {{ sec.label }}
              </div>
            </div>

            <!-- Submenu Items (18px icon) -->
            <div
              v-if="!isCollapsed && openSections[sec.id]"
              class="pr-6 pt-0.5 pb-1 space-y-1 transition-all duration-200 animate-fade-in"
            >
              <div
                v-for="child in sec.children"
                :key="child.id"
                :class="[
                  'flex items-center gap-2.5 px-3 py-2 rounded-lg text-xs font-semibold cursor-pointer transition-all duration-150 text-right relative',
                  isChildActive(child)
                    ? 'bg-[#1F0E03] text-[#F97316] font-bold border-r-2 border-[#F97316]'
                    : 'text-[#71717A] hover:bg-[#161616] hover:text-white',
                ]"
                @click="navigateTo(child.route)"
              >
                <component
                  :is="child.icon"
                  :class="[
                    'h-[18px] w-[18px] shrink-0 transition-colors',
                    isChildActive(child) ? 'text-[#F97316]' : 'text-[#71717A]',
                  ]"
                />
                <span class="truncate flex-1">{{ child.label }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer / User Info -->
      <div class="p-3 border-t border-[#242424] shrink-0 bg-[#0B0B0B]">
        <div v-if="!isCollapsed" class="flex items-center gap-3 p-2 rounded-xl">
          <div class="h-9 w-9 rounded-full bg-[#1F0E03] text-[#F97316] flex items-center justify-center font-bold text-xs shrink-0 border border-orange-950">
            <UserRound class="h-5 w-5" />
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-xs font-extrabold text-white truncate">{{ adminStore.userName }}</p>
            <p class="text-[10px] text-[#71717A] capitalize truncate">{{ adminStore.userRole }}</p>
          </div>
          <button
            class="p-1.5 text-[#71717A] hover:text-red-400 hover:bg-[#181818] rounded-lg transition-colors"
            title="تسجيل الخروج"
            @click="showLogoutDialog = true"
          >
            <LogOut class="h-[18px] w-[18px]" />
          </button>
        </div>

        <div v-else class="flex justify-center py-1">
          <button
            class="p-2 text-[#71717A] hover:text-red-400 hover:bg-[#181818] rounded-xl transition-colors"
            title="تسجيل الخروج"
            @click="showLogoutDialog = true"
          >
            <LogOut class="h-5 w-5" />
          </button>
        </div>
      </div>
    </aside>

    <!-- Mobile Drawer Overlay -->
    <div v-if="isMobileOpen" class="fixed inset-0 z-50 lg:hidden">
      <div class="fixed inset-0 bg-black/80 backdrop-blur-sm transition-opacity" @click="closeMobileSidebar" />
      <aside class="glass-liquid-nav fixed inset-y-0 right-0 w-[280px] max-w-[85vw] bg-[#0B0B0B]/95 border-l border-[#242424] flex flex-col z-50 shadow-2xl animate-slide-down">
        <div class="h-16 flex items-center justify-between px-4 border-b border-[#242424]">
          <BrandLogo size="sm" subtitle="لوحة التحكم" />
          <button class="p-2 text-[#71717A] hover:text-white rounded-lg" @click="closeMobileSidebar">
            <X class="h-5 w-5" />
          </button>
        </div>

        <div class="flex-1 overflow-y-auto py-4 px-3 space-y-2">
          <div v-for="sec in navigationSections" :key="sec.id" class="space-y-1">
            <div v-if="!sec.children" class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-semibold cursor-pointer" :class="isSectionActive(sec) ? 'bg-[#1F0E03] text-[#F97316]' : 'text-[#A1A1AA]'" @click="navigateTo(sec.route!)">
              <component :is="sec.icon" class="h-5 w-5 shrink-0 text-[#F97316]" />
              <span class="text-xs font-bold">{{ sec.label }}</span>
            </div>
            <div v-else class="space-y-1">
              <div class="flex items-center justify-between px-3 py-2 text-xs font-extrabold text-[#71717A] uppercase tracking-wider cursor-pointer" @click="toggleSection(sec.id)">
                <span class="flex items-center gap-2">
                  <component :is="sec.icon" class="h-5 w-5 text-[#F97316]" />
                  {{ sec.label }}
                </span>
                <ChevronDown class="h-4 w-4" :class="openSections[sec.id] ? 'rotate-180' : ''" />
              </div>
              <div v-if="openSections[sec.id]" class="pr-4 space-y-1">
                <div
                  v-for="child in sec.children"
                  :key="child.id"
                  :class="[
                    'flex items-center gap-2.5 px-3 py-2 rounded-lg text-xs font-semibold cursor-pointer',
                    isChildActive(child) ? 'bg-[#1F0E03] text-[#F97316] font-bold' : 'text-[#A1A1AA]',
                  ]"
                  @click="navigateTo(child.route)"
                >
                  <component :is="child.icon" class="h-[18px] w-[18px] shrink-0 text-[#F97316]" />
                  <span>{{ child.label }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="p-4 border-t border-[#242424] bg-[#0B0B0B]">
          <button class="btn-secondary w-full text-xs text-red-400 border-red-900/40" @click="showLogoutDialog = true">
            <LogOut class="h-[18px] w-[18px]" />
            تسجيل الخروج
          </button>
        </div>
      </aside>
    </div>

    <!-- Confirm Logout Dialog -->
    <ConfirmLogoutDialog
      :open="showLogoutDialog"
      :loading="loggingOut"
      @confirm="handleLogout"
      @cancel="showLogoutDialog = false"
    />
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #242424;
  border-radius: 4px;
}
</style>