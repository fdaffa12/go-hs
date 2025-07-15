<script setup>
import { onMounted, ref, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useAuthStore } from "./stores/auth";

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

// Global loading state
const isLoading = ref(false);
const loadingProgress = ref(0);

// Watch for auth store loading state
watch(
  () => authStore.isLoading,
  (newValue) => {
    isLoading.value = newValue;
    if (newValue) {
      // Simulate progress animation
      loadingProgress.value = 0;
      const interval = setInterval(() => {
        if (loadingProgress.value < 90) {
          loadingProgress.value += Math.random() * 30;
        }
        if (!authStore.isLoading) {
          loadingProgress.value = 100;
          setTimeout(() => {
            isLoading.value = false;
            loadingProgress.value = 0;
          }, 200);
          clearInterval(interval);
        }
      }, 100);
    }
  }
);

// Router navigation guards for loading
router.beforeEach((to, from, next) => {
  // Show loading bar for all route changes (except initial load)
  if (from.name !== undefined) {
    isLoading.value = true;
    loadingProgress.value = 0;

    // Simulate loading progress
    const interval = setInterval(() => {
      if (loadingProgress.value < 70) {
        loadingProgress.value += Math.random() * 20;
      }
    }, 50);

    // Store interval for cleanup
    router.loadingInterval = interval;
  }
  next();
});

router.afterEach(() => {
  // Complete loading after route change
  if (router.loadingInterval) {
    clearInterval(router.loadingInterval);
  }

  loadingProgress.value = 100;
  setTimeout(() => {
    isLoading.value = false;
    loadingProgress.value = 0;
  }, 200);
});

// Initialize auth store on app mount
onMounted(async () => {
  isLoading.value = true;
  await authStore.init();

  // Redirect logic
  if (
    authStore.isAuthenticated &&
    (route.path === "/login" || route.path === "/register")
  ) {
    router.push("/dashboard");
  } else if (!authStore.isAuthenticated && route.meta.requiresAuth) {
    router.push("/login");
  }

  isLoading.value = false;
});
</script>

<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <!-- Global Loading Bar -->
    <div
      v-if="isLoading"
      class="fixed top-0 left-0 right-0 z-[9999] h-1 bg-gray-200"
    >
      <div
        class="h-full bg-gradient-to-r from-primary-500 to-primary-600 transition-all duration-300 ease-out"
        :style="{ width: loadingProgress + '%' }"
      >
        <div class="h-full bg-white opacity-30 animate-pulse"></div>
      </div>
    </div>

    <router-view />
  </div>
</template>

<style scoped>
/* Styles are now handled by Tailwind CSS and global.css */
</style>
