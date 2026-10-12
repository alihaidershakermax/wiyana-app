<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowRight, User, Power, Activity, BookOpen, GraduationCap } from '@lucide/vue'
import {
  useStudent,
  useStudentCourses,
  useStudentActivity,
  useSuspendStudent,
  useActivateStudent,
} from '@/admin/queries/useStudents'
import StudentCard from '@/admin/components/StudentCard.vue'
import StudentActivity from '@/admin/components/StudentActivity.vue'
import StatePanel from '@/admin/components/StatePanel.vue'
import ConfirmDialog from '@/admin/components/ui/ConfirmDialog.vue'
import { useToast } from '@/admin/composables/useToast'
import { formatDate } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const studentId = computed(() => String(route.params.id ?? ''))
const { data: student, isLoading, isError, error, refetch } = useStudent(studentId)
const { data: courses, isLoading: coursesLoading, refetch: refetchCourses } = useStudentCourses(studentId)
const { data: activities, isLoading: activityLoading, refetch: refetchActivity } = useStudentActivity(studentId)

const suspend = useSuspendStudent()
const activate = useActivateStudent()
const showToggleDialog = ref(false)
const actionLoading = ref(false)

const isActive = computed(() => student.value?.is_active ?? true)

function openToggle() {
  showToggleDialog.value = true
}

function confirmToggle() {
  if (!student.value) return
  actionLoading.value = true
  const id = student.value.id
  const mutation = isActive.value ? suspend : activate
  mutation.mutate(id, {
    onSuccess: () => {
      toast.success(isActive.value ? 'تم إيقاف حساب الطالب' : 'تم تفعيل حساب الطالب')
      showToggleDialog.value = false
      actionLoading.value = false
      refetch()
    },
    onError: (err: unknown) => {
      toast.error('فشلت العملية', (err as Error)?.message)
      showToggleDialog.value = false
      actionLoading.value = false
    },
  })
}

function fullName(): string {
  const s = student.value
  if (!s) return ''
  return [s.first_name, s.last_name].filter(Boolean).join(' ') || s.username || `طالب ${s.telegram_id}`
}
</script>

<template>
  <div class="space-y-6">
    <button class="btn-ghost text-sm -mb-2" @click="router.push({ name: 'admin.students' })">
      <ArrowRight class="h-4 w-4" />
      العودة إلى الطلاب
    </button>

    <StatePanel
      :loading="isLoading"
      :error="isError ? (error as Error)?.message : null"
      :empty="!student"
      empty-title="الطالب غير موجود"
      :skeleton-count="4"
      @retry="() => refetch()"
    >
      <template v-if="student">
        <!-- Profile header -->
        <div class="card p-5">
          <div class="flex flex-wrap items-center gap-4">
            <div class="h-16 w-16 rounded-2xl bg-orange-50 dark:bg-orange-950/40 border border-orange-100 dark:border-orange-900/60 overflow-hidden shrink-0">
              <img v-if="student.photo_url" :src="student.photo_url" class="h-full w-full object-cover" alt="" />
              <div v-else class="h-full w-full flex items-center justify-center text-xl font-extrabold text-orange-500">
                {{ (student.first_name || student.username || '؟')[0] }}
              </div>
            </div>
            <div class="flex-1 min-w-0">
              <h1 class="text-lg font-extrabold">{{ fullName() }}</h1>
              <div class="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-slate-400 mt-1">
                <span dir="ltr">@{{ student.username || '—' }}</span>
                <span dir="ltr">ID: {{ student.telegram_id }}</span>
                <span>مسجل منذ {{ formatDate(student.created_at) ?? '—' }}</span>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <span
                :class="['text-[10px] font-bold px-2.5 py-1 rounded-full', isActive ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-400' : 'bg-red-100 text-red-700 dark:bg-red-950/40 dark:text-red-400']"
              >
                {{ isActive ? 'نشط' : 'موقوف' }}
              </span>
              <button
                class="btn-secondary text-sm"
                :class="isActive ? 'text-red-500 hover:bg-red-50 dark:hover:bg-red-950/40' : 'text-emerald-600'"
                @click="openToggle"
              >
                <Power class="h-4 w-4" />
                {{ isActive ? 'إيقاف الحساب' : 'تفعيل الحساب' }}
              </button>
            </div>
          </div>
        </div>

        <!-- Courses -->
        <div class="space-y-2">
          <div class="flex items-center gap-2 px-1">
            <BookOpen class="h-4 w-4 text-orange-500" />
            <h2 class="text-sm font-extrabold">المساقات المسجل بها</h2>
            <span class="text-[10px] font-bold bg-orange-50 dark:bg-orange-950/40 text-orange-600 dark:text-orange-400 px-2 py-0.5 rounded-full">
              {{ courses?.length ?? 0 }}
            </span>
          </div>
          <StatePanel
            :loading="coursesLoading"
            :error="null"
            :empty="!courses?.length"
            empty-title="لا توجد مساقات"
            empty-message="الطالب غير مسجل في أي مساق"
            :skeleton-count="2"
            @retry="refetchCourses"
          >
            <StudentCard
              :student="student"
              :courses="courses ?? []"
              @view-course="(courseId: string) => router.push({ name: 'admin.course-details', params: { id: courseId } })"
            />
          </StatePanel>
        </div>

        <!-- Activity -->
        <div class="space-y-2">
          <div class="flex items-center gap-2 px-1">
            <Activity class="h-4 w-4 text-orange-500" />
            <h2 class="text-sm font-extrabold">آخر النشاط</h2>
          </div>
          <StatePanel
            :loading="activityLoading"
            :error="null"
            :empty="!activities?.length"
            empty-title="لا يوجد نشاط"
            empty-message="لم يسجل الطالب أي نشاط بعد"
            :skeleton-count="3"
            @retry="refetchActivity"
          >
            <StudentActivity :activities="activities ?? []" />
          </StatePanel>
        </div>
      </template>
    </StatePanel>

    <ConfirmDialog
      :open="showToggleDialog"
      :title="isActive ? 'إيقاف حساب الطالب' : 'تفعيل حساب الطالب'"
      :message="isActive ? 'سيتم منع الطالب من الوصول إلى المنصة. هل أنت متأكد؟' : 'سيتمكن الطالب من الوصول إلى المنصة. هل أنت متأكد؟'"
      confirm-label="تأكيد"
      cancel-label="إلغاء"
      variant="danger"
      :loading="actionLoading"
      @confirm="confirmToggle"
      @cancel="showToggleDialog = false"
    />
  </div>
</template>
