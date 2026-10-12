<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowRight, Users } from '@lucide/vue'
import { useCourseStudents } from '@/admin/queries/useStudents'
import StudentTable from '@/admin/components/StudentTable.vue'
import StatePanel from '@/admin/components/StatePanel.vue'
import { useToast } from '@/admin/composables/useToast'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const courseId = ref(String(route.params.courseId ?? route.params.id ?? ''))
const { data, isLoading, isError, error, refetch } = useCourseStudents(courseId)

function viewStudent(s: { id: string }) {
  router.push({ name: 'admin.student-detail', params: { id: s.id } })
}
</script>

<template>
  <div class="space-y-6">
    <button class="btn-ghost text-sm -mb-2" @click="router.back()">
      <ArrowRight class="h-4 w-4" />
      العودة
    </button>

    <div class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="text-xl font-extrabold flex items-center gap-2">
          <Users class="h-5 w-5 text-orange-500" />
          طلاب المساق
        </h1>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
          الطلاب المسجلون في هذا المساق
        </p>
      </div>
    </div>

    <div class="card overflow-hidden">
      <StatePanel
        :loading="isLoading"
        :error="isError ? (error as Error)?.message : null"
        :empty="!data?.length"
        empty-title="لا يوجد طلاب"
        empty-message="لا يوجد طلاب مسجلون في هذا المساق"
        :skeleton-count="5"
        @retry="() => refetch()"
      >
        <StudentTable :students="data ?? []" @view="viewStudent" />
      </StatePanel>
    </div>
  </div>
</template>
