<script setup>
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";

// Initialize router and auth store
const router = useRouter();
const authStore = useAuthStore();

// Reactive data
const form = ref({
  email: "",
  password: "",
});

const loading = ref(false);
const error = ref("");
const success = ref("");
const showPassword = ref(false);
const rememberMe = ref(false);
const emailError = ref("");
const passwordError = ref("");

// Computed properties
const isFormValid = computed(() => {
  return (
    form.value.email.trim() !== "" &&
    form.value.password.trim() !== "" &&
    !emailError.value &&
    !passwordError.value &&
    isValidEmail(form.value.email)
  );
});

// Validation methods
const isValidEmail = (email) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

const validateEmail = () => {
  if (!form.value.email.trim()) {
    emailError.value = "Email harus diisi";
  } else if (!isValidEmail(form.value.email)) {
    emailError.value = "Format email tidak valid";
  } else {
    emailError.value = "";
  }
};

const validatePassword = () => {
  if (!form.value.password.trim()) {
    passwordError.value = "Password harus diisi";
  } else if (form.value.password.length < 6) {
    passwordError.value = "Password minimal 6 karakter";
  } else {
    passwordError.value = "";
  }
};

// Methods
const handleLogin = async () => {
  // Validate form before submission
  validateEmail();
  validatePassword();

  if (!isFormValid.value) {
    error.value = "Mohon periksa kembali data yang Anda masukkan";
    return;
  }

  loading.value = true;
  error.value = "";
  success.value = "";

  try {
    const result = await authStore.login({
      email: form.value.email,
      password: form.value.password,
    });

    if (result.success) {
      success.value = "Login berhasil! Selamat datang kembali.";

      // Save remember me preference
      if (rememberMe.value) {
        localStorage.setItem("rememberMe", "true");
        localStorage.setItem("userEmail", form.value.email);
      } else {
        localStorage.removeItem("rememberMe");
        localStorage.removeItem("userEmail");
      }

      // Reset form
      form.value = {
        email: "",
        password: "",
      };

      // Redirect to dashboard after successful login
      setTimeout(() => {
        router.push("/dashboard");
      }, 1000);
    } else {
      error.value = result.message || "Email atau password tidak valid";
    }
  } catch (err) {
    console.error("Login error:", err);
    error.value = err.message || "Email atau password tidak valid";
  } finally {
    loading.value = false;
  }
};

// Load remembered email on component mount
if (localStorage.getItem("rememberMe") === "true") {
  const savedEmail = localStorage.getItem("userEmail");
  if (savedEmail) {
    form.value.email = savedEmail;
    rememberMe.value = true;
  }
}
</script>

<template>
  <div
    class="min-h-screen bg-gradient-to-br from-primary-50 to-primary-100 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8"
  >
    <div class="max-w-md w-full">
      <!-- Header Section -->
      <div class="text-center mb-8">
        <div
          class="mx-auto h-16 w-16 bg-primary-600 rounded-full flex items-center justify-center mb-4"
        >
          <svg
            class="h-8 w-8 text-white"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
            />
          </svg>
        </div>
        <h2 class="text-3xl font-bold text-gray-900 mb-2">
          Selamat Datang Kembali
        </h2>
        <p class="text-gray-600">Masuk ke akun Anda untuk melanjutkan</p>
      </div>

      <!-- Login Form -->
      <div class="bg-white shadow-xl rounded-2xl p-8 border border-gray-100">
        <form @submit.prevent="handleLogin" class="space-y-6">
          <!-- Email Field -->
          <div class="space-y-2">
            <label
              for="email"
              class="block text-sm font-semibold text-gray-700"
            >
              Alamat Email
            </label>
            <div class="relative">
              <div
                class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
              >
                <svg
                  class="h-5 w-5 text-gray-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207"
                  />
                </svg>
              </div>
              <input
                id="email"
                v-model="form.email"
                type="email"
                required
                :disabled="loading"
                placeholder="nama@email.com"
                class="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition-colors duration-200 disabled:bg-gray-50 disabled:text-gray-500"
                :class="{
                  'border-red-300 focus:ring-red-500 focus:border-red-500':
                    emailError,
                }"
                @blur="validateEmail"
              />
            </div>
            <p v-if="emailError" class="text-sm text-red-600">
              {{ emailError }}
            </p>
          </div>

          <!-- Password Field -->
          <div class="space-y-2">
            <label
              for="password"
              class="block text-sm font-semibold text-gray-700"
            >
              Password
            </label>
            <div class="relative">
              <div
                class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
              >
                <svg
                  class="h-5 w-5 text-gray-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                  />
                </svg>
              </div>
              <input
                id="password"
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                required
                :disabled="loading"
                placeholder="Masukkan password Anda"
                class="block w-full pl-10 pr-10 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition-colors duration-200 disabled:bg-gray-50 disabled:text-gray-500"
                :class="{
                  'border-red-300 focus:ring-red-500 focus:border-red-500':
                    passwordError,
                }"
                @blur="validatePassword"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center"
              >
                <svg
                  v-if="showPassword"
                  class="h-5 w-5 text-gray-400 hover:text-gray-600"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"
                  />
                </svg>
                <svg
                  v-else
                  class="h-5 w-5 text-gray-400 hover:text-gray-600"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                  />
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                  />
                </svg>
              </button>
            </div>
            <p v-if="passwordError" class="text-sm text-red-600">
              {{ passwordError }}
            </p>
          </div>

          <!-- Remember Me & Forgot Password -->
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <input
                id="remember-me"
                v-model="rememberMe"
                type="checkbox"
                class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
              />
              <label for="remember-me" class="ml-2 block text-sm text-gray-700">
                Ingat saya
              </label>
            </div>
            <a
              href="#"
              class="text-sm font-medium text-primary-600 hover:text-primary-500 transition-colors"
            >
              Lupa password?
            </a>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="loading || !isFormValid"
            class="w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 transform hover:scale-[1.02] active:scale-[0.98]"
          >
            <svg
              v-if="loading"
              class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
            {{ loading ? "Memproses..." : "Masuk ke Akun" }}
          </button>
        </form>

        <!-- Error & Success Messages -->
        <div
          v-if="error"
          class="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg"
        >
          <div class="flex items-center">
            <svg
              class="h-5 w-5 text-red-400 mr-2"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <p class="text-sm text-red-800">{{ error }}</p>
          </div>
        </div>

        <div
          v-if="success"
          class="mt-4 p-4 bg-green-50 border border-green-200 rounded-lg"
        >
          <div class="flex items-center">
            <svg
              class="h-5 w-5 text-green-400 mr-2"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <p class="text-sm text-green-800">{{ success }}</p>
          </div>
        </div>

        <!-- Divider -->
        <div class="mt-6">
          <div class="relative">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t border-gray-300"></div>
            </div>
            <div class="relative flex justify-center text-sm">
              <span class="px-2 bg-white text-gray-500">atau</span>
            </div>
          </div>
        </div>

        <!-- Register Link -->
        <div class="mt-6 text-center">
          <p class="text-sm text-gray-600">
            Belum memiliki akun?
            <button
              type="button"
              @click="router.push('/register')"
              class="ml-1 font-semibold text-primary-600 hover:text-primary-500 transition-colors duration-200 underline-offset-4 hover:underline"
            >
              Daftar sekarang
            </button>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Styles are now handled by Tailwind CSS and global.css */
</style>
