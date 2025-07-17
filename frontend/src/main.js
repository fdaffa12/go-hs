import "./assets/global.css";

import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import { createPinia } from "pinia";
import Toast from "vue-toastification";
import "vue-toastification/dist/index.css";
import App from "./App.vue";
import "./assets/global.css";
import { useAuthStore } from "./stores/auth";

// Import pages
import Login from "./pages/Login.vue";
import Register from "./pages/Register.vue";
import Dashboard from "./pages/Dashboard.vue";
import Profile from "./pages/Profile.vue";
import Settings from "./pages/Settings.vue";
import UserManagement from "./pages/UserManagement.vue";
import DepartmentManagement from "./pages/DepartmentManagement.vue";
import EmployeeManagement from "./pages/EmployeeManagement.vue";
import BuyerManagement from "./pages/BuyerManagement.vue";
import StyleManagement from "./pages/StyleManagement.vue";
import LineScheduleManagement from "./pages/LineScheduleManagement.vue";
import HolidayManagement from "./pages/HolidayManagement.vue";

// Define routes
const routes = [
  { path: "/", redirect: "/login" },
  { path: "/login", component: Login },
  { path: "/register", component: Register },
  { path: "/dashboard", component: Dashboard, meta: { requiresAuth: true } },
  { path: "/profile", component: Profile, meta: { requiresAuth: true } },
  { path: "/settings", component: Settings, meta: { requiresAuth: true } },
  { path: "/users", component: UserManagement, meta: { requiresAuth: true } },
  {
    path: "/departments",
    component: DepartmentManagement,
    meta: { requiresAuth: true },
  },
  {
    path: "/employees",
    component: EmployeeManagement,
    meta: { requiresAuth: true },
  },
  {
    path: "/buyers",
    component: BuyerManagement,
    meta: { requiresAuth: true },
  },
  {
    path: "/styles",
    component: StyleManagement,
    meta: { requiresAuth: true },
  },
  {
    path: "/line-schedules",
    component: LineScheduleManagement,
    meta: { requiresAuth: true },
  },
  {
    path: "/holidays",
    component: HolidayManagement,
    meta: { requiresAuth: true },
  },
];

// Create router instance
const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Create Pinia instance
const pinia = createPinia();

// Toast configuration
const toastOptions = {
  position: "top-right",
  timeout: 5000,
  closeOnClick: true,
  pauseOnFocusLoss: true,
  pauseOnHover: true,
  draggable: true,
  draggablePercent: 0.6,
  showCloseButtonOnHover: false,
  hideProgressBar: false,
  closeButton: "button",
  icon: true,
  rtl: false,
};

// Create and mount app
const app = createApp(App);
app.use(pinia);
app.use(router);
app.use(Toast, toastOptions);

// Router guards (after pinia is installed)
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next("/login");
  } else if (
    (to.path === "/login" || to.path === "/register") &&
    authStore.isAuthenticated
  ) {
    next("/dashboard");
  } else {
    next();
  }
});

app.mount("#app");
