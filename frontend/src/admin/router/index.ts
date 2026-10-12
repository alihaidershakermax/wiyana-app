import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import { useAdminStore } from "@/admin/stores/useAdminStore";
import AdminLayout from "@/admin/components/layout/AdminLayout.vue";

declare module "vue-router" {
  interface RouteMeta {
    requiresAuth?: boolean;
    roles?: string[]; // panel roles allowed on this route
    permissions?: string[]; // all required permissions
  }
}

export const adminRoutes: RouteRecordRaw = {
  path: "/admin",
  component: AdminLayout,
  meta: { requiresAuth: true },
  children: [
    {
      path: "",
      name: "admin-root",
      redirect: () => {
        const store = useAdminStore();
        return store.isOwner ? "/admin/dashboard" : "/admin/teacher/dashboard";
      },
    },
    {
      path: "dashboard",
      name: "admin.dashboard",
      component: () => import("@/admin/views/OwnerDashboardView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"] },
    },
    {
      path: "analytics",
      name: "admin.analytics",
      component: () => import("@/admin/views/AnalyticsView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"], permissions: ["analytics.read"] },
    },
    {
      path: "courses",
      name: "admin.courses",
      component: () => import("@/admin/views/CoursesView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"], permissions: ["courses.read"] },
    },
    {
      path: "courses/create",
      name: "admin.course-create",
      component: () => import("@/admin/views/CourseCreateView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"], permissions: ["courses.create"] },
    },
    {
      path: "courses/:id",
      name: "admin.course-detail",
      component: () => import("@/admin/views/CourseDetailView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER", "TEACHER"], permissions: ["courses.read"] },
    },
    {
      path: "courses/:id/edit",
      name: "admin.course-edit",
      component: () => import("@/admin/views/CourseEditView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER", "TEACHER"], permissions: ["courses.update"] },
    },
    {
      path: "permissions",
      name: "admin.permissions",
      component: () => import("@/admin/views/PlaceholderView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"], permissions: ["permissions.manage"] },
    },
    {
      path: "subjects",
      name: "admin.subjects",
      redirect: "/admin/courses",
    },
    {
      path: "subjects/create",
      name: "admin.subject-create",
      redirect: "/admin/courses/create",
    },
    {
      path: "subjects/:id/edit",
      name: "admin.subject-edit",
      redirect: "/admin/courses",
    },
    {
      path: "courses/:courseId/lessons",
      name: "admin.course-lessons",
      component: () => import("@/admin/views/LessonsView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER", "TEACHER"], permissions: ["lessons.read"] },
    },
    {
      path: "courses/:courseId/lessons/create",
      name: "admin.lesson-create",
      component: () => import("@/admin/views/LessonCreateView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER", "TEACHER"], permissions: ["lessons.create"] },
    },
    {
      path: "lessons/:id/edit",
      name: "admin.lesson-edit",
      component: () => import("@/admin/views/LessonEditView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER", "TEACHER"], permissions: ["lessons.update"] },
    },
    {
      path: "lessons/:id",
      name: "admin.lesson-detail",
      component: () => import("@/admin/views/LessonDetailsView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER", "TEACHER"], permissions: ["lessons.read"] },
    },
    {
      path: "lessons",
      name: "admin.lessons",
      component: () => import("@/admin/views/LessonsView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"], permissions: ["lessons.read"] },
    },
    {
      path: "files",
      name: "admin.files",
      component: () => import("@/admin/views/PlaceholderView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"] },
    },
    {
      path: "reports",
      name: "admin.reports",
      component: () => import("@/admin/views/PlaceholderView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"], permissions: ["analytics.read"] },
    },
    {
      path: "students",
      name: "admin.students",
      component: () => import("@/admin/views/StudentsView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"], permissions: ["students.read"] },
    },
    {
      path: "students/:id",
      name: "admin.student-detail",
      component: () => import("@/admin/views/StudentDetailsView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER", "TEACHER"], permissions: ["students.read"] },
    },
    {
      path: "courses/:courseId/students",
      name: "admin.course-students",
      component: () => import("@/admin/views/CourseStudentsView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER", "TEACHER"], permissions: ["students.read"] },
    },
    {
      path: "teachers",
      name: "admin.teachers",
      component: () => import("@/admin/views/TeachersView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"], permissions: ["teachers.read"] },
    },
    {
      path: "teachers/create",
      name: "admin.teacher-create",
      component: () => import("@/admin/views/TeacherCreateView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"], permissions: ["teachers.manage"] },
    },
    {
      path: "subscriptions",
      name: "admin.subscriptions",
      component: () => import("@/admin/views/SubscriptionsView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"], permissions: ["subscriptions.read"] },
    },
    {
      path: "subscription-codes",
      name: "admin.subscription-codes",
      component: () => import("@/admin/views/SubscriptionCodesView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"], permissions: ["subscriptions.manage"] },
    },
    {
      path: "banned-users",
      name: "admin.banned-users",
      component: () => import("@/admin/views/BannedUsersView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"] },
    },
    {
      path: "notifications",
      name: "admin.notifications",
      component: () => import("@/admin/views/AdminNotificationsView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"] },
    },
    {
      path: "settings",
      name: "admin.settings",
      component: () => import("@/admin/views/SettingsView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"], permissions: ["settings.manage"] },
    },
    {
      path: "audit-logs",
      name: "admin.audit-logs",
      component: () => import("@/admin/views/AuditLogsView.vue"),
      meta: { requiresAuth: true, roles: ["OWNER"] },
    },
    // ---- TEACHER panel ----
    {
      path: "teacher",
      name: "admin.teacher-root",
      redirect: { name: "admin.teacher.dashboard" },
    },
    {
      path: "teacher/dashboard",
      name: "admin.teacher.dashboard",
      component: () => import("@/admin/views/TeacherDashboardView.vue"),
      meta: { requiresAuth: true, roles: ["TEACHER", "OWNER"] },
    },
    {
      path: "teacher/courses",
      name: "admin.teacher.courses",
      component: () => import("@/views/teacher/TeacherUploadsView.vue"),
      meta: { requiresAuth: true, roles: ["TEACHER", "OWNER"] },
    },
    {
      path: "teacher/lessons",
      name: "admin.teacher.lessons",
      component: () => import("@/views/teacher/TeacherUploadsView.vue"),
      meta: { requiresAuth: true, roles: ["TEACHER", "OWNER"] },
    },
    {
      path: "teacher/upload",
      name: "admin.teacher.upload",
      component: () => import("@/views/teacher/TeacherUploadsView.vue"),
      meta: { requiresAuth: true, roles: ["TEACHER", "OWNER"] },
    },
    {
      path: "teacher/students",
      name: "admin.teacher.students",
      component: () => import("@/admin/views/StudentsView.vue"),
      meta: { requiresAuth: true, roles: ["TEACHER", "OWNER"] },
    },
    {
      path: "teacher/attachments",
      name: "admin.teacher.attachments",
      component: () => import("@/views/teacher/TeacherUploadsView.vue"),
      meta: { requiresAuth: true, roles: ["TEACHER", "OWNER"] },
    },
    {
      path: "teacher/analytics",
      name: "admin.teacher.analytics",
      component: () => import("@/admin/views/AnalyticsView.vue"),
      meta: { requiresAuth: true, roles: ["TEACHER", "OWNER"] },
    },
    {
      path: "teacher/notifications",
      name: "admin.teacher.notifications",
      component: () => import("@/admin/views/AdminNotificationsView.vue"),
      meta: { requiresAuth: true, roles: ["TEACHER", "OWNER"] },
    },
    {
      path: "teacher/profile",
      name: "admin.teacher.profile",
      component: () => import("@/admin/views/PlaceholderView.vue"),
      meta: { requiresAuth: true, roles: ["TEACHER"] },
    },
  ],
};

export const adminAuthRoutes: RouteRecordRaw[] = [
  {
    path: "/admin/login",
    name: "admin.login",
    component: () => import("@/admin/views/LoginView.vue"),
    meta: { requiresAuth: false },
  },
  {
    path: "/admin/forbidden",
    name: "admin.forbidden",
    component: () => import("@/admin/views/ForbiddenView.vue"),
    meta: { requiresAuth: false },
  },
  {
    path: "/admin/unauthorized",
    name: "admin.unauthorized",
    component: () => import("@/admin/views/UnauthorizedView.vue"),
    meta: { requiresAuth: false },
  },
];

// Guard bootstrap: make sure the session is restored once before guarding.
let authInitialized = false;

async function ensureAuth() {
  const store = useAdminStore();
  if (!authInitialized) {
    authInitialized = true;
    await store.initializeAuth();
  }
}

function roleHomePath(store: ReturnType<typeof useAdminStore>): string {
  return store.isOwner ? "/admin/dashboard" : "/admin/teacher/dashboard";
}

export function registerAdminGuards(router: ReturnType<typeof createRouter>): void {
  router.beforeEach(async (to) => {
    if (to.path.startsWith("/admin")) {
      const store = useAdminStore();

      // The Control Panel is always dark-themed; child views use `dark:` variants,
      // so force the `dark` class on <html> regardless of the student app's theme.
      document.documentElement.classList.add("dark");

      // Restore session on first admin navigation.
      if (!authInitialized) {
        await ensureAuth();
      }

      // Session expired / network failure → back to login.
      if (
        ["session_expired", "network_error", "server_error"].includes(
          store.authState
        ) &&
        to.name !== "admin.login"
      ) {
        return { name: "admin.login", query: { redirect: to.fullPath } };
      }

      if (to.meta.requiresAuth) {
        if (!store.isAuthenticated) {
          return { name: "admin.login", query: { redirect: to.fullPath } };
        }

        // Role guard (enforced by backend too — this is UX only).
        const allowedRoles = to.meta.roles;
        if (allowedRoles && allowedRoles.length && !store.hasRole(allowedRoles)) {
          return { name: "admin.forbidden" };
        }

        // Permission guard (backend enforces as well).
        const requiredPerms = to.meta.permissions;
        if (requiredPerms && requiredPerms.length) {
          if (!store.hasAllPermissions(requiredPerms)) {
            return { name: "admin.forbidden" };
          }
        }
      }

      // Authenticated user visiting /admin/login → role home.
      if (to.name === "admin.login" && store.isAuthenticated) {
        const redirect = typeof to.query.redirect === "string" ? to.query.redirect : null;
        return redirect || roleHomePath(store);
      }

      if (store.isAuthenticated && to.name === "admin-root" ) {
        return roleHomePath(store);
      }
    } else {
      // Leaving the Control Panel: restore the student app's own theme preference.
      const darkMode = localStorage.getItem("wy_dark_mode") === "true";
      document.documentElement.classList.toggle("dark", darkMode);
    }
    return true;
  });
}

// Session-expired event fired by the API interceptor.
if (typeof window !== "undefined") {
  window.addEventListener("admin:session-expired", () => {
    const store = useAdminStore();
    store.$patch({ authState: "session_expired" });
    // The guard handles redirection on next navigation; also hard-redirect here
    // so the user lands on the login page immediately.
    const adminRouter = adminAuthRoutes; // no-op to keep ref
    void adminRouter;
    window.location.href = "/admin/login";
  });
}

export function createAdminRouter() {
  const router = createRouter({
    history: createWebHistory(),
    routes: [...adminAuthRoutes, adminRoutes],
  });
  registerAdminGuards(router);
  return router;
}
