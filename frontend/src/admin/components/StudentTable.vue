<script setup lang="ts">
import { computed } from 'vue'
import { Eye, Pencil, Power, Trash2, BookOpen, PlayCircle } from '@lucide/vue'
import type { Student } from '@/admin/types/panel'
import { formatDate } from '@/utils/format'
import ActionDropdown, { type MenuItem } from '@/admin/components/ui/ActionDropdown.vue'

const props = defineProps<{
  students: Student[]
  loading?: boolean
}>()

const emit = defineEmits<{
  view: [student: Student]
  toggleStatus: [student: Student]
  delete: [student: Student]
}>()

function fullName(s: Student): string {
  return [s.first_name, s.last_name].filter(Boolean).join(' ') || (s as any).name || s.username || `طالب ${s.telegram_id}`
}

function getActionItems(s: Student): MenuItem[] {
  const isActive = s.is_active || (s as any).status === 'ACTIVE'
  return [
    {
      id: 'view',
      label: 'عرض تفاصيل الطالب',
      icon: Eye,
      onClick: () => emit('view', s),
    },
    {
      id: 'edit',
      label: 'تعديل البيانات',
      icon: Pencil,
      onClick: () => emit('view', s),
    },
    {
      id: 'status',
      label: isActive ? 'إيقاف الحساب' : 'تفعيل الحساب',
      icon: Power,
      danger: isActive,
      onClick: () => emit('toggleStatus', s),
    },
    {
      id: 'delete',
      label: 'حذف الطالب',
      icon: Trash2,
      danger: true,
      onClick: () => emit('delete' as any, s),
    },
  ]
}
</script>

<template>
  <div class="overflow-x-auto">
    <table class="w-full text-sm text-[#F5F5F5]">
      <thead>
        <tr class="bg-[#0F0F0F] border-b border-[#242424] text-right">
          <th class="py-3.5 px-4 text-xs font-bold text-[#A1A1AA]">الطالب</th>
          <th class="py-3.5 px-4 text-xs font-bold text-[#A1A1AA]">الحالة</th>
          <th class="py-3.5 px-4 text-xs font-bold text-[#A1A1AA] hidden md:table-cell">المساقات</th>
          <th class="py-3.5 px-4 text-xs font-bold text-[#A1A1AA] hidden md:table-cell">اشتراكات نشطة</th>
          <th class="py-3.5 px-4 text-xs font-bold text-[#A1A1AA] hidden lg:table-cell">آخر نشاط</th>
          <th class="py-3.5 px-4 text-xs font-bold text-[#A1A1AA] hidden lg:table-cell">تاريخ التسجيل</th>
          <th class="py-3.5 px-4 text-xs font-bold text-[#A1A1AA] w-16 text-center">إجراءات</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-[#242424] bg-[#111111]">
        <tr
          v-for="s in students"
          :key="s.id"
          class="hover:bg-[#161616] transition-colors"
        >
          <td class="py-3.5 px-4">
            <div class="flex items-center gap-3">
              <div class="h-9 w-9 rounded-full bg-[#1F0E03] border border-orange-900/40 overflow-hidden shrink-0 flex items-center justify-center font-bold text-xs text-[#F97316]">
                <img v-if="s.photo_url" :src="s.photo_url" class="h-full w-full object-cover" alt="" />
                <span v-else>{{ (s.first_name || (s as any).name || s.username || '؟')[0] }}</span>
              </div>
              <div class="min-w-0">
                <p class="font-bold text-white truncate">{{ fullName(s) }}</p>
                <p class="text-[11px] text-[#A1A1AA] truncate" dir="ltr">
                  @{{ s.username || '—' }} · {{ s.telegram_id }}
                </p>
              </div>
            </div>
          </td>
          <td class="py-3.5 px-4">
            <span
              :class="[
                'text-[10px] font-extrabold px-2.5 py-1 rounded-full border inline-block',
                (s.is_active || (s as any).status === 'ACTIVE')
                  ? 'bg-[#052e16] text-[#22C55E] border-[#14532d]'
                  : 'bg-[#450a0a] text-[#EF4444] border-[#7f1d1d]',
              ]"
            >
              {{ (s.is_active || (s as any).status === 'ACTIVE') ? 'نشط' : 'موقوف' }}
            </span>
          </td>
          <td class="py-3.5 px-4 hidden md:table-cell">
            <span class="flex items-center gap-1 text-xs text-[#A1A1AA]">
              <BookOpen class="h-3.5 w-3.5 text-[#71717A]" />
              {{ s.courses_count ?? (s as any).coursesCount ?? 0 }}
            </span>
          </td>
          <td class="py-3.5 px-4 hidden md:table-cell">
            <span class="flex items-center gap-1 text-xs text-[#A1A1AA]">
              <PlayCircle class="h-3.5 w-3.5 text-[#71717A]" />
              {{ s.active_subscriptions ?? (s as any).coursesCount ?? 0 }}
            </span>
          </td>
          <td class="py-3.5 px-4 hidden lg:table-cell text-xs text-[#71717A]">{{ formatDate(s.last_active_at || (s as any).lastActive) ?? '—' }}</td>
          <td class="py-3.5 px-4 hidden lg:table-cell text-xs text-[#71717A]">{{ formatDate(s.created_at || (s as any).createdAt) ?? '—' }}</td>
          <td class="py-3.5 px-4 text-center">
            <ActionDropdown :items="getActionItems(s)" />
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
