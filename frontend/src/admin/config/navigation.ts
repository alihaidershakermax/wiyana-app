import type { Component } from 'vue'
import {
  LayoutDashboard,
  FolderOpen,
  BookOpen,
  Library,
  Layers3,
  PlayCircle,
  Users,
  UserRound,
  GraduationCap,
  ShieldCheck,
  CreditCard,
  Ticket,
  WalletCards,
  Bell,
  Megaphone,
  ChartNoAxesCombined,
  Settings,
  Send,
  Video,
  History,
} from '@lucide/vue'

export interface NavChildItem {
  id: string
  label: string
  route: string
  icon: Component
  roles?: string[]
  permissions?: string[]
}

export interface NavSection {
  id: string
  label: string
  icon: Component
  roles?: string[]
  permissions?: string[]
  children?: NavChildItem[]
  route?: string
}

// ---------------------------------------------------------------------------
// OWNER iOS-Style Icon Navigation Configuration
// ---------------------------------------------------------------------------
export const ownerNavSections: NavSection[] = [
  {
    id: 'dashboard',
    label: 'الرئيسية',
    icon: LayoutDashboard,
    route: '/admin/dashboard',
    roles: ['OWNER'],
  },
  {
    id: 'content',
    label: 'إدارة المحتوى',
    icon: FolderOpen,
    roles: ['OWNER'],
    children: [
      { id: 'courses', label: 'المساقات (الدورات)', route: '/admin/courses', icon: BookOpen },
      { id: 'sections', label: 'أقسام المساقات', route: '/admin/courses', icon: Layers3 },
      { id: 'lessons', label: 'المحاضرات ومقاطع الفيديو', route: '/admin/teacher/lessons', icon: PlayCircle },
    ],
  },
  {
    id: 'users',
    label: 'المستخدمون والحماية',
    icon: Users,
    roles: ['OWNER'],
    children: [
      { id: 'students', label: 'الطلاب', route: '/admin/students', icon: UserRound },
      { id: 'teachers', label: 'الأساتذة', route: '/admin/teachers', icon: GraduationCap },
      { id: 'banned-users', label: 'المحظورون وسجل الحماية', route: '/admin/banned-users', icon: ShieldCheck },
      { id: 'permissions', label: 'الصلاحيات', route: '/admin/settings', icon: ShieldCheck },
    ],
  },
  {
    id: 'subscriptions',
    label: 'الاشتراكات والأكواد',
    icon: CreditCard,
    roles: ['OWNER'],
    children: [
      { id: 'subs-codes', label: 'توليد وأكواد الاشتراك', route: '/admin/subscription-codes', icon: Ticket },
      { id: 'subs-list', label: 'سجل اشتراكات الطلاب', route: '/admin/subscriptions', icon: CreditCard },
      { id: 'subs-payments', label: 'المدفوعات والحسابات', route: '/admin/subscriptions', icon: WalletCards },
    ],
  },
  {
    id: 'communication',
    label: 'التواصل',
    icon: Bell,
    roles: ['OWNER'],
    children: [
      { id: 'notifications', label: 'الإشعارات', route: '/admin/notifications', icon: Bell },
      { id: 'announcements', label: 'الإعلانات', route: '/admin/notifications', icon: Megaphone },
    ],
  },
  {
    id: 'analytics',
    label: 'التحليلات',
    icon: ChartNoAxesCombined,
    roles: ['OWNER'],
    children: [
      { id: 'analytics-overview', label: 'نظرة عامة', route: '/admin/analytics', icon: ChartNoAxesCombined },
      { id: 'analytics-students', label: 'الطلاب', route: '/admin/analytics', icon: UserRound },
      { id: 'analytics-courses', label: 'المساقات', route: '/admin/analytics', icon: BookOpen },
    ],
  },
  {
    id: 'system',
    label: 'النظام',
    icon: Settings,
    roles: ['OWNER'],
    children: [
      { id: 'settings', label: 'إعدادات المنصة', route: '/admin/settings', icon: Settings },
      { id: 'telegram', label: 'Telegram', route: '/admin/settings', icon: Send },
      { id: 'youtube', label: 'YouTube', route: '/admin/settings', icon: Video },
      { id: 'audit', label: 'سجل النشاط', route: '/admin/audit-logs', icon: History },
    ],
  },
]

// ---------------------------------------------------------------------------
// TEACHER iOS-Style Icon Navigation Configuration
// ---------------------------------------------------------------------------
export const teacherNavSections: NavSection[] = [
  {
    id: 'teacher-dashboard',
    label: 'الرئيسية',
    icon: LayoutDashboard,
    route: '/admin/teacher/dashboard',
    roles: ['TEACHER'],
  },
  {
    id: 'teacher-content',
    label: 'المحتوى والمحاضرات',
    icon: FolderOpen,
    roles: ['TEACHER'],
    children: [
      { id: 't-courses', label: 'مساقاتي', route: '/admin/teacher/courses', icon: BookOpen },
      { id: 't-lessons', label: 'محاضراتي', route: '/admin/teacher/lessons', icon: PlayCircle },
      { id: 't-upload', label: 'رفع محاضرة جديدة', route: '/admin/teacher/upload', icon: PlayCircle },
    ],
  },
  {
    id: 'teacher-students',
    label: 'الطلاب والتحليلات',
    icon: Users,
    roles: ['TEACHER'],
    children: [
      { id: 't-students', label: 'طلابي', route: '/admin/teacher/students', icon: UserRound },
      { id: 't-analytics', label: 'تحليلات الأداء', route: '/admin/teacher/analytics', icon: ChartNoAxesCombined },
    ],
  },
  {
    id: 'teacher-system',
    label: 'الحساب والتواصل',
    icon: Settings,
    roles: ['TEACHER'],
    children: [
      { id: 't-notifications', label: 'الإشعارات', route: '/admin/teacher/notifications', icon: Bell },
    ],
  },
]

export function getNavigationSections(role: string): NavSection[] {
  return role === 'OWNER' ? ownerNavSections : teacherNavSections
}
