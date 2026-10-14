import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import {
  adminRoutes,
  adminAuthRoutes,
  registerAdminGuards,
} from "@/admin/router";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "home",
    component: () => import("@/views/HomeView.vue"),
  },
  {
    path: "/courses",
    name: "courses",
    component: () => import("@/views/CoursesView.vue"),
  },
  {
    path: "/favorites",
    name: "favorites",
    component: () => import("@/views/FavoritesView.vue"),
  },
  {
    path: "/profile",
    name: "profile",
    component: () => import("@/views/ProfileView.vue"),
  },
  {
    path: "/courses/:id",
    name: "course-detail",
    component: () => import("@/views/CourseDetailView.vue"),
  },
  {
    path: "/subjects/:id",
    name: "subject-detail",
    component: () => import("@/views/SubjectDetailView.vue"),
  },
  {
    path: "/lessons/:id",
    name: "lesson-player",
    component: () => import("@/views/VideoPlayerView.vue"),
  },
  {
    path: "/search",
    name: "search",
    component: () => import("@/views/SearchView.vue"),
  },
  {
    path: "/notifications",
    name: "notifications",
    component: () => import("@/views/NotificationsView.vue"),
  },
  {
    path: "/teacher",
    component: () => import("@/views/teacher/TeacherLayout.vue"),
    children: [
      {
        path: "",
        name: "teacher-dashboard",
        component: () => import("@/views/teacher/TeacherDashboardView.vue"),
      },
      {
        path: "uploads",
        name: "teacher-uploads",
        component: () => import("@/views/teacher/TeacherUploadsView.vue"),
      },
      {
        path: "questions",
        name: "teacher-questions",
        component: () => import("@/views/teacher/TeacherQuestionsView.vue"),
      },
      {
        path: "earnings",
        name: "teacher-earnings",
        component: () => import("@/views/teacher/TeacherEarningsView.vue"),
      },
    ],
  },
    // ---- Control Panel (external) routes ----
    ...adminAuthRoutes,
    adminRoutes,
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

if (import.meta.env.DEV) {
  router.beforeEach((to, from) => {
    console.log(`[Router] Navigating from "${from.path}" to "${to.path}"`);
  });
}

router.onError((err) => {
  if (import.meta.env.DEV) console.error("[Router] Navigation error:", err);
});

// Register the Control Panel auth/role/permission guards.
registerAdminGuards(router);

export default router;
