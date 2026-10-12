<script setup lang="ts">
import AdminSidebar from './AdminSidebar.vue'
import AdminHeader from './AdminHeader.vue'
import ToastContainer from '@/admin/components/ui/ToastContainer.vue'
</script>

<template>
  <div dir="rtl" class="admin-glass-shell min-h-screen text-[#F5F5F5] flex antialiased selection:bg-[#F97316] selection:text-white">
    <!-- Sidebar Component -->
    <AdminSidebar />

    <!-- Main Wrapper -->
    <div class="flex-1 flex flex-col min-w-0 min-h-screen overflow-x-hidden">
      <!-- Topbar Header Component -->
      <AdminHeader />

      <!-- View Content Area -->
      <main class="admin-glass-shell flex-1 p-4 sm:p-6 lg:p-8">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" class="animate-fade-in" />
          </transition>
        </router-view>
      </main>
    </div>

    <!-- Global Toast Container -->
    <ToastContainer />
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>