<template>
  <AuthenticatedLayout :user="authStore.user">
    <div class="space-y-8">
      <!-- Welcome Section -->
      <section class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div class="text-center">
          <h2
            class="text-3xl font-bold text-gray-900 mb-2 flex items-center justify-center gap-2"
          >
            <span class="animate-bounce">👋</span>
            Selamat Datang, {{ authStore.user?.username }}!
          </h2>
          <p class="text-lg text-gray-600 mb-6">
            Anda berhasil masuk ke aplikasi CORS dengan Vue.js dan Go Backend
          </p>
          <div class="flex justify-center gap-8 flex-wrap">
            <div class="flex flex-col gap-1">
              <span class="text-sm text-gray-500 font-medium">Email:</span>
              <span class="text-base text-gray-900 font-semibold">{{
                authStore.user?.email
              }}</span>
            </div>
            <div class="flex flex-col gap-1">
              <span class="text-sm text-gray-500 font-medium">User ID:</span>
              <span class="text-base text-gray-900 font-semibold"
                >#{{ authStore.user?.id }}</span
              >
            </div>
          </div>
        </div>
      </section>

      <!-- Features Section -->
      <section class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h3 class="text-2xl font-bold text-gray-900 mb-6 text-center">
          Fitur Aplikasi
        </h3>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div class="card p-4 text-center hover:shadow-md transition-shadow">
            <div class="text-4xl mb-3">🔐</div>
            <h4 class="text-lg font-semibold text-gray-900 mb-2">
              Autentikasi JWT
            </h4>
            <p class="text-sm text-gray-600">
              Sistem keamanan dengan JSON Web Token untuk melindungi data
              pengguna
            </p>
          </div>

          <div class="card p-4 text-center hover:shadow-md transition-shadow">
            <div class="text-4xl mb-3">🌐</div>
            <h4 class="text-lg font-semibold text-gray-900 mb-2">
              CORS Support
            </h4>
            <p class="text-sm text-gray-600">
              Cross-Origin Resource Sharing yang memungkinkan komunikasi antar
              domain
            </p>
          </div>

          <div class="card p-4 text-center hover:shadow-md transition-shadow">
            <div class="text-4xl mb-3">⚡</div>
            <h4 class="text-lg font-semibold text-gray-900 mb-2">
              Vue.js Frontend
            </h4>
            <p class="text-sm text-gray-600">
              Interface modern dan responsif dengan Vue.js 3 dan Composition API
            </p>
          </div>

          <div class="card p-4 text-center hover:shadow-md transition-shadow">
            <div class="text-4xl mb-3">🚀</div>
            <h4 class="text-lg font-semibold text-gray-900 mb-2">Go Backend</h4>
            <p class="text-sm text-gray-600">
              Backend yang cepat dan efisien menggunakan Go dengan arsitektur
              MVC
            </p>
          </div>
        </div>
      </section>

      <!-- Quick Actions -->
      <section class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h3 class="text-2xl font-bold text-gray-900 mb-6 text-center">
          Aksi Cepat
        </h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <button
            @click="testConnection"
            :disabled="testing"
            class="btn btn-primary flex items-center justify-center gap-2"
          >
            <span class="text-lg">🔗</span>
            <span>{{ testing ? "Testing..." : "Test Koneksi Backend" }}</span>
          </button>

          <button
            @click="refreshUserData"
            :disabled="refreshing"
            class="btn btn-secondary flex items-center justify-center gap-2"
          >
            <span class="text-lg">🔄</span>
            <span>{{
              refreshing ? "Refreshing..." : "Refresh Data User"
            }}</span>
          </button>
        </div>
      </section>

      <!-- Status Messages -->
      <div
        v-if="message"
        class="alert"
        :class="{
          'alert-success': messageType === 'success',
          'alert-error': messageType === 'error',
          'alert-info': messageType === 'info',
        }"
      >
        {{ message }}
      </div>
    </div>
  </AuthenticatedLayout>
</template>

<script setup>
import { ref } from "vue";
import { useAuthStore } from "../stores/auth";
import AuthenticatedLayout from "../layouts/AuthenticatedLayout.vue";

// Initialize auth store
const authStore = useAuthStore();

// Reactive data
const testing = ref(false);
const refreshing = ref(false);
const message = ref("");
const messageType = ref("success");

// Methods
const testConnection = async () => {
  testing.value = true;
  message.value = "";

  try {
    const baseUrl =
      import.meta.env.VITE_API_BASE_URL || "http://localhost:8081";
    const response = await fetch(`${baseUrl}/api/hello`, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${authStore.token}`,
        "Content-Type": "application/json",
      },
    });

    if (response.ok) {
      message.value = "Koneksi ke backend berhasil!";
      messageType.value = "success";
    } else {
      message.value = "Test koneksi gagal";
      messageType.value = "error";
    }
  } catch (error) {
    message.value = error.message || "Terjadi kesalahan saat test koneksi";
    messageType.value = "error";
  } finally {
    testing.value = false;
    // Clear message after 5 seconds
    setTimeout(() => {
      message.value = "";
    }, 5000);
  }
};

const refreshUserData = async () => {
  refreshing.value = true;
  message.value = "";

  try {
    await authStore.getCurrentUser();
    message.value = "Data user berhasil diperbarui";
    messageType.value = "success";
  } catch (error) {
    message.value = error.message || "Terjadi kesalahan saat memperbarui data";
    messageType.value = "error";
  } finally {
    refreshing.value = false;
    // Clear message after 5 seconds
    setTimeout(() => {
      message.value = "";
    }, 5000);
  }
};
</script>

<style scoped>
/* Styles are now handled by Tailwind CSS and global.css */
</style>
