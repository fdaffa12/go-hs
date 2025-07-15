<template>
  <div class="min-h-screen flex flex-col bg-gray-50">
    <!-- Header -->
    <header
      class="bg-white border-b border-gray-200 shadow-sm sticky top-0 z-50"
    >
      <div class="container">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center space-x-3">
            <img
              alt="Vue logo"
              class="w-8 h-8 rounded-lg"
              src="../assets/logo.svg"
            />
            <div>
              <h1 class="text-xl font-bold text-gray-900">CORS App</h1>
              <p class="text-xs text-gray-500">Vue.js & Go</p>
            </div>
          </div>

          <nav class="flex space-x-1">
            <button
              @click="$emit('navigate', 'login')"
              :class="['nav-link', currentPage === 'login' ? 'active' : '']"
            >
              <svg
                class="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"
                />
              </svg>
              Masuk
            </button>
            <button
              @click="$emit('navigate', 'register')"
              :class="['nav-link', currentPage === 'register' ? 'active' : '']"
            >
              <svg
                class="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"
                />
              </svg>
              Daftar
            </button>
          </nav>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1">
      <slot></slot>
    </main>

    <!-- Footer -->
    <footer class="bg-white border-t border-gray-200 mt-auto">
      <div class="container py-12">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div>
            <h3 class="text-lg font-semibold text-gray-900 mb-4">
              CORS Application
            </h3>
            <p class="text-gray-600 text-sm leading-relaxed">
              Aplikasi modern dengan Vue.js frontend dan Go backend yang
              mendemonstrasikan autentikasi JWT dan penanganan CORS yang aman.
            </p>
          </div>

          <div>
            <h4 class="text-sm font-semibold text-gray-900 mb-4">Teknologi</h4>
            <ul class="space-y-2 text-sm text-gray-600">
              <li class="flex items-center">
                <span
                  class="w-1.5 h-1.5 bg-primary-500 rounded-full mr-2"
                ></span>
                Vue.js 3 + Composition API
              </li>
              <li class="flex items-center">
                <span
                  class="w-1.5 h-1.5 bg-primary-500 rounded-full mr-2"
                ></span>
                Go dengan Gin Framework
              </li>
              <li class="flex items-center">
                <span
                  class="w-1.5 h-1.5 bg-primary-500 rounded-full mr-2"
                ></span>
                MySQL Database
              </li>
              <li class="flex items-center">
                <span
                  class="w-1.5 h-1.5 bg-primary-500 rounded-full mr-2"
                ></span>
                JWT Authentication
              </li>
            </ul>
          </div>

          <div>
            <h4 class="text-sm font-semibold text-gray-900 mb-4">Fitur</h4>
            <ul class="space-y-2 text-sm text-gray-600">
              <li class="flex items-center">
                <span
                  class="w-1.5 h-1.5 bg-primary-500 rounded-full mr-2"
                ></span>
                Registrasi & Login User
              </li>
              <li class="flex items-center">
                <span
                  class="w-1.5 h-1.5 bg-primary-500 rounded-full mr-2"
                ></span>
                Autentikasi JWT
              </li>
              <li class="flex items-center">
                <span
                  class="w-1.5 h-1.5 bg-primary-500 rounded-full mr-2"
                ></span>
                CORS Configuration
              </li>
              <li class="flex items-center">
                <span
                  class="w-1.5 h-1.5 bg-primary-500 rounded-full mr-2"
                ></span>
                Responsive Design
              </li>
            </ul>
          </div>
        </div>

        <div class="border-t border-gray-200 pt-8 mt-8">
          <p class="text-center text-sm text-gray-500">
            &copy; 2025 CORS Application. Dibuat dengan ❤️ menggunakan Vue.js &
            Go
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from "vue";

// Define props
const props = defineProps({
  currentPage: {
    type: String,
    default: "",
  },
});

// Define emits
const emit = defineEmits(["navigate"]);
</script>

<style scoped>
/* Styles are now handled by Tailwind CSS and global.css */
</style>
