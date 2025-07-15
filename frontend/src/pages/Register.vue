<script setup>
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";
import { authService } from "../services/api";

// Initialize router and auth store
const router = useRouter();
const authStore = useAuthStore();

// Add employees data
const employees = ref([]);
const loading = ref(false);
const loadingEmployees = ref(false);

const form = ref({
  nik: "",
  name: "",
  email: "",
  password: "",
  confirmPassword: "",
  level: 1,
});

// Load employees on component mount
const loadEmployees = async () => {
  try {
    loadingEmployees.value = true;
    console.log("Fetching employees...");
    const response = await authService.getAvailableEmployees();
    console.log("API Response:", response);
    if (response.success) {
      employees.value = response.data;
      console.log("Loaded employees:", employees.value);
    } else {
      console.error("Failed to load employees:", response.message);
      error.value = "Gagal memuat data karyawan: " + response.message;
    }
  } catch (err) {
    console.error("Error loading employees:", err);
    error.value =
      "Gagal memuat data karyawan: " + (err.message || "Unknown error");
  } finally {
    loadingEmployees.value = false;
  }
};

// Handle employee selection
const handleEmployeeSelect = (event) => {
  const selectedNIK = event.target.value;
  console.log("Selected NIK:", selectedNIK);
  const selectedEmployee = employees.value.find(
    (emp) => emp.nik === selectedNIK
  );
  console.log("Selected employee:", selectedEmployee);
  if (selectedEmployee) {
    form.value.nik = selectedEmployee.nik;
    form.value.name = selectedEmployee.name;
    validateNIK();
  }
};

onMounted(() => {
  loadEmployees();
});

const error = ref("");
const success = ref("");
const showPassword = ref(false);
const showConfirmPassword = ref(false);
const acceptTerms = ref(false);

// Validation errors
const nikError = ref("");
const nameError = ref("");
const emailError = ref("");
const passwordError = ref("");
const confirmPasswordError = ref("");

// Computed properties
const isFormValid = computed(() => {
  return (
    form.value.nik &&
    form.value.name &&
    form.value.email &&
    form.value.password &&
    form.value.confirmPassword &&
    !nikError.value &&
    !nameError.value &&
    !emailError.value &&
    !passwordError.value &&
    !confirmPasswordError.value &&
    form.value.password === form.value.confirmPassword
  );
});

const passwordStrength = computed(() => {
  const password = form.value.password;
  if (!password) return 0;

  let strength = 0;
  if (password.length >= 8) strength++;
  if (/[a-z]/.test(password)) strength++;
  if (/[A-Z]/.test(password)) strength++;
  if (/[0-9]/.test(password)) strength++;
  if (/[^A-Za-z0-9]/.test(password)) strength++;

  return strength;
});

const passwordStrengthText = computed(() => {
  const strength = passwordStrength.value;
  if (strength === 0) return "";
  if (strength <= 2) return "Lemah";
  if (strength <= 3) return "Sedang";
  if (strength <= 4) return "Kuat";
  return "Sangat Kuat";
});

const passwordStrengthColor = computed(() => {
  const strength = passwordStrength.value;
  if (strength <= 2) return "bg-red-500";
  if (strength <= 3) return "bg-yellow-500";
  if (strength <= 4) return "bg-blue-500";
  return "bg-green-500";
});

const passwordStrengthTextColor = computed(() => {
  const strength = passwordStrength.value;
  if (strength <= 2) return "text-red-600";
  if (strength <= 3) return "text-yellow-600";
  if (strength <= 4) return "text-blue-600";
  return "text-green-600";
});

const passwordStrengthWidth = computed(() => {
  return `${(passwordStrength.value / 5) * 100}%`;
});

// Validation methods
const isValidEmail = (email) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

const validateNIK = () => {
  const nik = form.value.nik.trim();
  if (!nik) {
    nikError.value = "NIK wajib diisi";
  } else if (!nik.startsWith("HS")) {
    nikError.value = "NIK harus diawali dengan 'HS'";
  } else if (nik.length < 3) {
    nikError.value = "NIK minimal 3 karakter";
  } else if (nik.length > 10) {
    // Changed from 7 to 10 to match employee table
    nikError.value = "NIK maksimal 10 karakter";
  } else if (!/^HS\d+$/.test(nik)) {
    nikError.value = "NIK harus berformat HS diikuti angka";
  } else {
    nikError.value = "";
  }
};

const validateName = () => {
  const name = form.value.name.trim();
  if (!name) {
    nameError.value = "Nama wajib diisi";
  } else if (name.length < 3) {
    nameError.value = "Nama minimal 3 karakter";
  } else if (name.length > 25) {
    nameError.value = "Nama maksimal 25 karakter";
  } else {
    nameError.value = "";
  }
};

const validatePassword = () => {
  const password = form.value.password;
  if (!password) {
    passwordError.value = "Password wajib diisi";
  } else if (password.length < 8) {
    passwordError.value = "Password minimal 8 karakter";
  } else if (password.length > 128) {
    passwordError.value = "Password maksimal 128 karakter";
  } else {
    passwordError.value = "";
  }

  // Re-validate confirm password if it exists
  if (form.value.confirmPassword) {
    validateConfirmPassword();
  }
};

const validateConfirmPassword = () => {
  const confirmPassword = form.value.confirmPassword;
  if (!confirmPassword) {
    confirmPasswordError.value = "Konfirmasi password wajib diisi";
  } else if (confirmPassword !== form.value.password) {
    confirmPasswordError.value = "Password tidak cocok";
  } else {
    confirmPasswordError.value = "";
  }
};

const validateEmail = () => {
  const email = form.value.email.trim();
  if (!email) {
    emailError.value = "Email wajib diisi";
  } else if (!isValidEmail(email)) {
    emailError.value = "Format email tidak valid";
  } else if (email.length > 25) {
    emailError.value = "Email maksimal 25 karakter";
  } else {
    emailError.value = "";
  }
};

const handleRegister = async () => {
  try {
    loading.value = true;
    error.value = "";
    success.value = "";

    // Final validation
    validateNIK();
    validateName();
    validateEmail();
    validatePassword();
    validateConfirmPassword();

    if (!isFormValid.value) {
      error.value = "Mohon perbaiki kesalahan pada form";
      return;
    }

    if (!acceptTerms.value) {
      error.value = "Anda harus menyetujui syarat dan ketentuan";
      return;
    }

    await authStore.register({
      nik: form.value.nik.trim(),
      name: form.value.name.trim(),
      email: form.value.email.trim().toLowerCase(),
      password: form.value.password,
      level: form.value.level,
    });

    success.value = "Registrasi berhasil! Selamat datang.";

    // Reset form
    form.value = {
      nik: "HS",
      name: "",
      email: "",
      password: "",
      confirmPassword: "",
      level: 1,
    };
    acceptTerms.value = false;

    // Clear validation errors
    nikError.value = "";
    nameError.value = "";
    emailError.value = "";
    passwordError.value = "";
    confirmPasswordError.value = "";

    // Redirect to dashboard after successful registration
    setTimeout(() => {
      router.push("/dashboard");
    }, 1000);
  } catch (err) {
    console.error("Registration error:", err);
    if (err.message?.toLowerCase().includes("nik not found")) {
      nikError.value = "NIK tidak terdaftar sebagai karyawan";
    } else if (err.message?.toLowerCase().includes("nik already")) {
      nikError.value = "NIK sudah terdaftar";
    } else if (err.message?.toLowerCase().includes("email")) {
      emailError.value = "Email sudah terdaftar";
    } else {
      error.value = err.message || "Registrasi gagal";
    }
  } finally {
    loading.value = false;
  }
};
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
              d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"
            />
          </svg>
        </div>
        <h2 class="text-3xl font-bold text-gray-900 mb-2">
          Bergabung dengan Kami
        </h2>
        <p class="text-gray-600">
          Buat akun baru untuk memulai perjalanan Anda
        </p>
      </div>

      <!-- Register Form -->
      <div class="bg-white shadow-xl rounded-2xl p-8 border border-gray-100">
        <form @submit.prevent="handleRegister" class="space-y-6">
          <!-- NIK Field -->
          <div class="space-y-2">
            <label for="nik" class="block text-sm font-semibold text-gray-700">
              NIK
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
                    d="M10 6H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V8a2 2 0 00-2-2h-5m-4 0V5a2 2 0 114 0v1m-4 0a2 2 0 104 0m-5 8a2 2 0 100-4 2 2 0 000 4zm0 0c1.306 0 2.417.835 2.83 2M9 14a3.001 3.001 0 00-2.83 2M15 11h3m-3 4h2"
                  />
                </svg>
              </div>
              <select
                id="nik"
                v-model="form.nik"
                @change="handleEmployeeSelect"
                :disabled="loadingEmployees"
                class="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition-colors duration-200 disabled:bg-gray-50 disabled:text-gray-500"
              >
                <option value="">Pilih Karyawan</option>
                <option
                  v-for="employee in employees"
                  :key="employee.nik"
                  :value="employee.nik"
                >
                  {{ employee.nik }} - {{ employee.name }}
                </option>
              </select>
            </div>
            <p v-if="nikError" class="text-sm text-red-600">
              {{ nikError }}
            </p>
            <p
              v-else-if="form.nik && !nikError"
              class="text-sm text-green-600 flex items-center"
            >
              <svg
                class="h-4 w-4 mr-1"
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
              Format NIK valid
            </p>
          </div>

          <div class="space-y-2">
            <label for="name" class="block text-sm font-semibold text-gray-700">
              Nama Lengkap
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
                    d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                  />
                </svg>
              </div>
              <input
                id="name"
                v-model="form.name"
                type="text"
                required
                :disabled="loading"
                maxlength="25"
                placeholder="Masukkan nama lengkap"
                class="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition-colors duration-200 disabled:bg-gray-50 disabled:text-gray-500"
                :class="{
                  'border-red-300 focus:ring-red-500 focus:border-red-500':
                    nameError,
                }"
                @blur="validateName"
                @input="validateName"
              />
            </div>
            <p v-if="nameError" class="text-sm text-red-600">
              {{ nameError }}
            </p>
            <p
              v-else-if="form.name && !nameError"
              class="text-sm text-green-600 flex items-center"
            >
              <svg
                class="h-4 w-4 mr-1"
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
              Format nama valid
            </p>
          </div>

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
                @input="validateEmail"
              />
            </div>
            <p v-if="emailError" class="text-sm text-red-600">
              {{ emailError }}
            </p>
            <p
              v-else-if="form.email && !emailError"
              class="text-sm text-green-600 flex items-center"
            >
              <svg
                class="h-4 w-4 mr-1"
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
              Format email valid
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
                placeholder="Minimal 8 karakter"
                class="block w-full pl-10 pr-10 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition-colors duration-200 disabled:bg-gray-50 disabled:text-gray-500"
                :class="{
                  'border-red-300 focus:ring-red-500 focus:border-red-500':
                    passwordError,
                }"
                @blur="validatePassword"
                @input="validatePassword"
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
            <div v-if="form.password" class="space-y-1">
              <div class="flex items-center space-x-2">
                <div class="flex-1 bg-gray-200 rounded-full h-2">
                  <div
                    class="h-2 rounded-full transition-all duration-300"
                    :class="passwordStrengthColor"
                    :style="{ width: passwordStrengthWidth }"
                  ></div>
                </div>
                <span
                  class="text-xs font-medium"
                  :class="passwordStrengthTextColor"
                >
                  {{ passwordStrengthText }}
                </span>
              </div>
            </div>
            <p v-if="passwordError" class="text-sm text-red-600">
              {{ passwordError }}
            </p>
          </div>

          <!-- Confirm Password Field -->
          <div class="space-y-2">
            <label
              for="confirmPassword"
              class="block text-sm font-semibold text-gray-700"
            >
              Konfirmasi Password
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
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
              </div>
              <input
                id="confirmPassword"
                v-model="form.confirmPassword"
                :type="showConfirmPassword ? 'text' : 'password'"
                required
                :disabled="loading"
                placeholder="Ulangi password Anda"
                class="block w-full pl-10 pr-10 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition-colors duration-200 disabled:bg-gray-50 disabled:text-gray-500"
                :class="{
                  'border-red-300 focus:ring-red-500 focus:border-red-500':
                    confirmPasswordError,
                }"
                @blur="validateConfirmPassword"
                @input="validateConfirmPassword"
              />
              <button
                type="button"
                @click="showConfirmPassword = !showConfirmPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center"
              >
                <svg
                  v-if="showConfirmPassword"
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
            <p v-if="confirmPasswordError" class="text-sm text-red-600">
              {{ confirmPasswordError }}
            </p>
            <p
              v-else-if="
                form.confirmPassword && form.password === form.confirmPassword
              "
              class="text-sm text-green-600 flex items-center"
            >
              <svg
                class="h-4 w-4 mr-1"
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
              Password cocok
            </p>
          </div>

          <!-- Terms and Conditions -->
          <div class="flex items-start">
            <input
              id="terms"
              v-model="acceptTerms"
              type="checkbox"
              class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded mt-1"
            />
            <label for="terms" class="ml-2 block text-sm text-gray-700">
              Saya menyetujui
              <a
                href="#"
                class="font-medium text-primary-600 hover:text-primary-500 transition-colors"
              >
                Syarat dan Ketentuan
              </a>
              serta
              <a
                href="#"
                class="font-medium text-primary-600 hover:text-primary-500 transition-colors"
              >
                Kebijakan Privasi
              </a>
            </label>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="loading || !isFormValid || !acceptTerms"
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
            {{ loading ? "Memproses..." : "Buat Akun" }}
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

        <!-- Login Link -->
        <div class="mt-6 text-center">
          <p class="text-sm text-gray-600">
            Sudah memiliki akun?
            <button
              type="button"
              @click="router.push('/login')"
              class="ml-1 font-semibold text-primary-600 hover:text-primary-500 transition-colors duration-200 underline-offset-4 hover:underline"
            >
              Masuk di sini
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
