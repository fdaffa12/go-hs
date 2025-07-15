<script setup>
import { useAuthStore } from "../stores/auth";
import AuthenticatedLayout from "../layouts/AuthenticatedLayout.vue";

// Initialize auth store
const authStore = useAuthStore();
</script>

<template>
  <AuthenticatedLayout :user="authStore.user">
    <div class="space-y-8">
      <!-- Settings Header -->
      <section class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h2 class="text-2xl font-bold text-gray-900 mb-2">Settings</h2>
        <p class="text-gray-600">
          Kelola pengaturan aplikasi dan preferensi Anda
        </p>
      </section>

      <!-- Application Settings -->
      <section class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-6">
          Pengaturan Aplikasi
        </h3>
        <div class="space-y-6">
          <!-- Theme Setting -->
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-gray-900">Theme</h4>
              <p class="text-sm text-gray-500">Pilih tema tampilan aplikasi</p>
            </div>
            <select
              class="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            >
              <option value="light">Light</option>
              <option value="dark">Dark</option>
              <option value="auto">Auto</option>
            </select>
          </div>

          <!-- Language Setting -->
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-gray-900">Bahasa</h4>
              <p class="text-sm text-gray-500">Pilih bahasa interface</p>
            </div>
            <select
              class="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            >
              <option value="id">Bahasa Indonesia</option>
              <option value="en">English</option>
            </select>
          </div>

          <!-- Notifications -->
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-gray-900">Notifikasi</h4>
              <p class="text-sm text-gray-500">
                Terima notifikasi dari aplikasi
              </p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" class="sr-only peer" checked />
              <div
                class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-600"
              ></div>
            </label>
          </div>
        </div>
      </section>

      <!-- Security Settings -->
      <section class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-6">Keamanan</h3>
        <div class="space-y-6">
          <!-- Two Factor Authentication -->
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-gray-900">
                Two-Factor Authentication
              </h4>
              <p class="text-sm text-gray-500">
                Tambahkan lapisan keamanan ekstra
              </p>
            </div>
            <button class="btn btn-secondary btn-sm">Enable 2FA</button>
          </div>

          <!-- Session Management -->
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-gray-900">Active Sessions</h4>
              <p class="text-sm text-gray-500">Kelola sesi login aktif</p>
            </div>
            <button class="btn btn-secondary btn-sm">View Sessions</button>
          </div>

          <!-- Password Change -->
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-gray-900">Password</h4>
              <p class="text-sm text-gray-500">
                Terakhir diubah 30 hari yang lalu
              </p>
            </div>
            <button class="btn btn-primary btn-sm">Change Password</button>
          </div>
        </div>
      </section>

      <!-- API Settings -->
      <section class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-6">
          API & Integrasi
        </h3>
        <div class="space-y-6">
          <!-- API Key -->
          <div>
            <h4 class="text-sm font-medium text-gray-900 mb-2">API Key</h4>
            <div class="flex space-x-2">
              <input
                type="password"
                value="sk-1234567890abcdef"
                readonly
                class="flex-1 px-3 py-2 border border-gray-300 rounded-lg bg-gray-50 text-gray-500"
              />
              <button class="btn btn-secondary btn-sm">
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
                    d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                  />
                </svg>
              </button>
              <button class="btn btn-primary btn-sm">Regenerate</button>
            </div>
          </div>

          <!-- Webhook URL -->
          <div>
            <h4 class="text-sm font-medium text-gray-900 mb-2">Webhook URL</h4>
            <input
              type="url"
              placeholder="https://your-app.com/webhook"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            />
          </div>
        </div>
      </section>

      <!-- Danger Zone -->
      <section class="bg-white rounded-lg shadow-sm border border-red-200 p-6">
        <h3 class="text-lg font-semibold text-red-900 mb-6">Danger Zone</h3>
        <div class="space-y-4">
          <div class="p-4 bg-red-50 rounded-lg border border-red-200">
            <h4 class="text-sm font-medium text-red-900 mb-2">
              Delete Account
            </h4>
            <p class="text-sm text-red-700 mb-4">
              Permanently delete your account and all associated data. This
              action cannot be undone.
            </p>
            <button class="btn btn-danger btn-sm">Delete Account</button>
          </div>
        </div>
      </section>

      <!-- Save Button -->
      <div class="flex justify-end">
        <button class="btn btn-primary">
          <svg
            class="w-4 h-4 mr-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M5 13l4 4L19 7"
            />
          </svg>
          Save Changes
        </button>
      </div>
    </div>
  </AuthenticatedLayout>
</template>

<style scoped>
/* Styles are now handled by Tailwind CSS and global.css */
</style>
