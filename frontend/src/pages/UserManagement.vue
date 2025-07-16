<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { useToast } from "vue-toastification";
import { useAuthStore } from "../stores/auth";
import { userService, authService } from "../services/api";
import AuthenticatedLayout from "../layouts/AuthenticatedLayout.vue";

const authStore = useAuthStore();
const toast = useToast();

// Validation error refs
const nameError = ref("");
const emailError = ref("");
const nikError = ref("");
const passwordError = ref("");
const levelError = ref("");

// Reactive data
const users = ref([]);
const loading = ref(false);
const searchQuery = ref("");

// Add these after the existing refs
const employees = ref([]);
const loadingEmployees = ref(false);

// Pagination and page size options
const pageSizeOptions = [10, 20, 50, 100, -1]; // -1 represents "All"
const pageSize = ref(10);
const currentPage = ref(1);
const totalItems = ref(0);
const totalPages = ref(0);

// Add computed property for displayed pages
const displayedPages = computed(() => {
  const delta = 2;
  const range = [];
  const rangeWithDots = [];
  let l;

  for (let i = 1; i <= totalPages.value; i++) {
    if (
      i === 1 ||
      i === totalPages.value ||
      (i >= currentPage.value - delta && i <= currentPage.value + delta)
    ) {
      range.push(i);
    }
  }

  for (let i of range) {
    if (l) {
      if (i - l === 2) {
        rangeWithDots.push(l + 1);
      } else if (i - l !== 1) {
        rangeWithDots.push("...");
      }
    }
    rangeWithDots.push(i);
    l = i;
  }

  return rangeWithDots;
});

// Modal states
const showCreateModal = ref(false);
const showEditModal = ref(false);
const showDeleteModal = ref(false);
const showPassword = ref(false);
const submitting = ref(false);

// Form data
const userForm = ref({
  nik: "",
  name: "",
  email: "",
  password: "",
  level: 1,
});

const userToEdit = ref(null);
const userToDelete = ref(null);

// Computed
const filteredUsers = computed(() => {
  if (!searchQuery.value) return users.value;

  const query = searchQuery.value.toLowerCase();
  return users.value.filter(
    (user) =>
      user.nik.toLowerCase().includes(query) ||
      user.name.toLowerCase().includes(query) ||
      user.email.toLowerCase().includes(query)
  );
});

// Methods for pagination
const handlePageSizeChange = async () => {
  currentPage.value = 1; // Reset to first page when changing page size
  if (pageSize.value === -1) {
    // If "All" is selected, get total count first
    try {
      const response = await userService.getAllUsers(1, 1);
      if (response.success) {
        pageSize.value = response.data.total_items;
      }
    } catch (error) {
      console.error("Error getting total count:", error);
      pageSize.value = 100; // Fallback to 100 if error
    }
  }
  fetchUsers();
};

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--;
    fetchUsers();
  }
};

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
    fetchUsers();
  }
};

const goToPage = (page) => {
  if (page !== "..." && page !== currentPage.value) {
    currentPage.value = page;
    fetchUsers();
  }
};

// Add debounced search method
const debouncedSearch = ref(null);

const handleSearch = () => {
  if (debouncedSearch.value) {
    clearTimeout(debouncedSearch.value);
  }
  debouncedSearch.value = setTimeout(() => {
    currentPage.value = 1; // Reset to first page when searching
    fetchUsers();
  }, 300);
};

// Watch for search query changes
watch(searchQuery, () => {
  handleSearch();
});

// Update fetchUsers method
const fetchUsers = async () => {
  loading.value = true;
  try {
    const data = await userService.getAllUsers(
      currentPage.value,
      pageSize.value === -1 ? 999999 : pageSize.value,
      searchQuery.value
    );
    if (data.success) {
      // Update pagination data
      totalItems.value = data.data.total_items;
      totalPages.value = data.data.total_pages;
      currentPage.value = data.data.current_page;
      pageSize.value = data.data.page_size;
      users.value = data.data.users || [];
    } else {
      console.error("Failed to fetch users:", data.message);
      toast.error("Gagal memuat data pengguna: " + data.message);
    }
  } catch (error) {
    console.error("Error fetching users:", error);
    toast.error("Terjadi kesalahan saat memuat data pengguna");
  } finally {
    loading.value = false;
  }
};

const createUser = async () => {
  submitting.value = true;
  try {
    // Generate a random password
    const randomPassword = Math.random().toString(36).slice(-8);

    const data = await userService.createUser({
      ...userForm.value,
      password: randomPassword,
    });
    if (data.success) {
      await fetchUsers();
      closeModal();
      toast.success("Pengguna berhasil dibuat dengan password default");
    } else {
      toast.error("Gagal membuat pengguna: " + data.message);
    }
  } catch (error) {
    console.error("Error creating user:", error);
    toast.error("Terjadi kesalahan saat membuat pengguna");
  } finally {
    submitting.value = false;
  }
};

const editUser = (user) => {
  userToEdit.value = user;
  userForm.value = {
    nik: user.nik,
    name: user.name,
    email: user.email,
    level: user.level,
    profile_picture: user.profile_picture,
    password: "", // Empty for edit mode
  };
  showEditModal.value = true;
};

const updateUser = async () => {
  submitting.value = true;
  try {
    console.log("Updating user:", userToEdit.value.nik);

    // Create update data object
    const updateData = {
      name: userForm.value.name,
      email: userForm.value.email,
      level: parseInt(userForm.value.level),
    };

    // Only include password if it's provided
    if (userForm.value.password && userForm.value.password.trim() !== "") {
      updateData.password = userForm.value.password;
    }

    console.log("Update data:", {
      ...updateData,
      password: updateData.password ? "[REDACTED]" : undefined,
    });

    const data = await userService.updateUser(userToEdit.value.nik, updateData);
    console.log("Update response:", data);

    if (data.success) {
      await fetchUsers();
      closeModal();
      toast.success("Pengguna berhasil diperbarui");
    } else {
      console.error("Update failed:", data.message);
      toast.error(data.message || "Gagal memperbarui pengguna");
    }
  } catch (error) {
    console.error("Error updating user:", error);
    let errorMessage = "Terjadi kesalahan saat memperbarui pengguna";

    if (error.message) {
      if (error.message.toLowerCase().includes("email")) {
        errorMessage = "Email sudah digunakan oleh pengguna lain";
      } else if (error.message.toLowerCase().includes("password")) {
        errorMessage = "Format password tidak valid";
      } else {
        errorMessage = error.message;
      }
    }

    toast.error(errorMessage);
  } finally {
    submitting.value = false;
  }
};

const deleteUser = (user) => {
  userToDelete.value = user;
  showDeleteModal.value = true;
};

const confirmDelete = async () => {
  if (!userToDelete.value) return;

  submitting.value = true;
  try {
    const response = await userService.deleteUser(userToDelete.value.nik);
    if (response.success) {
      toast.success("User berhasil dihapus");
      showDeleteModal.value = false;
      await fetchUsers();
    } else {
      toast.error(response.message || "Gagal menghapus user");
    }
  } catch (error) {
    console.error("Error deleting user:", error);
    toast.error(error.message || "Terjadi kesalahan saat menghapus user");
  } finally {
    submitting.value = false;
  }
};

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
      toast.error("Gagal memuat data karyawan: " + response.message);
    }
  } catch (err) {
    console.error("Error loading employees:", err);
    toast.error(
      "Gagal memuat data karyawan: " + (err.message || "Unknown error")
    );
  } finally {
    loadingEmployees.value = false;
  }
};

const handleEmployeeSelect = (event) => {
  const selectedNIK = event.target.value;
  console.log("Selected NIK:", selectedNIK);
  const selectedEmployee = employees.value.find(
    (emp) => emp.nik === selectedNIK
  );
  console.log("Selected employee:", selectedEmployee);
  if (selectedEmployee) {
    userForm.value.nik = selectedEmployee.nik;
    userForm.value.name = selectedEmployee.name;
    validateNIK();
  }
};

const openCreateModal = () => {
  userForm.value = {
    nik: "",
    name: "",
    email: "",
    password: "",
    level: 1,
  };
  loadEmployees(); // Load employees when opening create modal
  showCreateModal.value = true;
};

const closeModal = () => {
  showCreateModal.value = false;
  showEditModal.value = false;
  showPassword.value = false;
  userForm.value = {
    nik: "",
    name: "",
    email: "",
    password: "",
    level: 1,
  };
};

const getUserInitials = (name) => {
  return name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
    .substring(0, 2);
};

const formatDate = (dateString) => {
  if (!dateString) return "N/A";

  try {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
      return "Invalid Date";
    }
    return date.toLocaleDateString("id-ID", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
    });
  } catch (error) {
    return "Invalid Date";
  }
};

// Lifecycle
onMounted(() => {
  fetchUsers();
});

// Add these computed properties and methods
const passwordStrength = computed(() => {
  const password = userForm.value.password;
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
  if (!userForm.value.password) return "";
  if (strength <= 2) return "Lemah";
  if (strength <= 3) return "Sedang";
  if (strength <= 4) return "Kuat";
  return "Sangat Kuat";
});

const passwordStrengthColor = computed(() => {
  const strength = passwordStrength.value;
  if (!userForm.value.password) return "";
  if (strength <= 2) return "bg-red-500";
  if (strength <= 3) return "bg-yellow-500";
  if (strength <= 4) return "bg-blue-500";
  return "bg-green-500";
});

const passwordStrengthWidth = computed(() => {
  return `${(passwordStrength.value / 5) * 100}%`;
});

const passwordStrengthTextColor = computed(() => {
  const strength = passwordStrength.value;
  if (!userForm.value.password) return "text-gray-500";
  if (strength <= 2) return "text-red-500";
  if (strength <= 3) return "text-yellow-500";
  if (strength <= 4) return "text-blue-500";
  return "text-green-500";
});

const validatePassword = (password, isEdit = false) => {
  if (!password && !isEdit) {
    passwordError.value = "Password wajib diisi";
    return false;
  }
  if (password && password.length < 8) {
    passwordError.value = "Password minimal 8 karakter";
    return false;
  }
  if (password && password.length > 128) {
    passwordError.value = "Password maksimal 128 karakter";
    return false;
  }
  passwordError.value = "";
  return true;
};

const isFormValid = computed(() => {
  // Basic validation
  if (!userForm.value.name || !userForm.value.email) return false;

  // Check for validation errors
  if (nameError.value || emailError.value) return false;

  // If password is provided, check password validation
  if (userForm.value.password && passwordError.value) return false;

  return true;
});

// Validation functions
const validateName = (name) => {
  if (!name) {
    nameError.value = "Nama wajib diisi";
    return false;
  }
  if (name.length < 3) {
    nameError.value = "Nama minimal 3 karakter";
    return false;
  }
  if (name.length > 25) {
    nameError.value = "Nama maksimal 25 karakter";
    return false;
  }
  nameError.value = "";
  return true;
};

const validateEmail = (email) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!email) {
    emailError.value = "Email wajib diisi";
    return false;
  }
  if (!emailRegex.test(email)) {
    emailError.value = "Format email tidak valid";
    return false;
  }
  if (email.length > 25) {
    emailError.value = "Email maksimal 25 karakter";
    return false;
  }
  emailError.value = "";
  return true;
};

const validateNIK = (nik) => {
  if (!nik) {
    nikError.value = "NIK wajib diisi";
    return false;
  }
  if (!nik.startsWith("HS")) {
    nikError.value = "NIK harus diawali dengan 'HS'";
    return false;
  }
  if (nik.length < 3 || nik.length > 10) {
    nikError.value = "NIK harus antara 3-10 karakter";
    return false;
  }
  if (!/^HS\d+$/.test(nik)) {
    nikError.value = "NIK harus berformat HS diikuti angka";
    return false;
  }
  nikError.value = "";
  return true;
};

const validateLevel = (level) => {
  if (!level) {
    levelError.value = "Level wajib diisi";
    return false;
  }
  if (![1, 2, 3].includes(Number(level))) {
    levelError.value = "Level harus 1, 2, atau 3";
    return false;
  }
  levelError.value = "";
  return true;
};

const validateForm = (isEdit = false) => {
  const isNameValid = validateName(userForm.value.name);
  const isEmailValid = validateEmail(userForm.value.email);
  const isNIKValid = validateNIK(userForm.value.nik);
  const isPasswordValid = validatePassword(userForm.value.password, isEdit);
  const isLevelValid = validateLevel(userForm.value.level);

  return (
    isNameValid && isEmailValid && isNIKValid && isPasswordValid && isLevelValid
  );
};

// Update the handleSubmit function to use validateForm
const handleSubmit = async () => {
  loading.value = true;
  error.value = "";

  try {
    if (!validateForm(isEditMode.value)) {
      error.value = "Mohon perbaiki kesalahan pada form";
      return;
    }

    // ... rest of the submit logic ...
  } catch (err) {
    console.error("Error submitting form:", err);
    error.value = err.message || "Terjadi kesalahan";
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <AuthenticatedLayout :user="authStore.user">
    <div class="space-y-6">
      <!-- Header -->
      <div
        class="bg-white rounded-lg shadow-sm border border-gray-200 p-4 sm:p-6"
      >
        <div
          class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4"
        >
          <div>
            <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">
              User Management
            </h1>
            <p class="mt-1 text-sm text-gray-600">
              Kelola pengguna dalam sistem
            </p>
          </div>
          <button
            @click="openCreateModal"
            class="btn btn-primary flex items-center justify-center gap-2 w-full sm:w-auto"
          >
            <svg
              class="w-5 h-5 flex-shrink-0"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 6v6m0 0v6m0-6h6m-6 0H6"
              ></path>
            </svg>
            <span class="hidden sm:inline">Tambah User</span>
            <span class="sm:hidden">Tambah</span>
          </button>
        </div>
      </div>

      <!-- Search and Filter -->
      <div
        class="bg-white rounded-lg shadow-sm border border-gray-200 p-3 sm:p-4"
      >
        <div class="flex flex-col sm:flex-row gap-3 sm:gap-4">
          <div class="flex-1">
            <div class="relative">
              <svg
                class="w-4 h-4 sm:w-5 sm:h-5 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                ></path>
              </svg>
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Cari berdasarkan NIK, nama, atau email..."
                class="w-full pl-9 sm:pl-10 pr-4 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
          <button
            @click="fetchUsers"
            class="bg-gray-100 hover:bg-gray-200 text-gray-700 px-3 sm:px-4 py-2 rounded-lg flex items-center justify-center gap-2 transition-colors w-full sm:w-auto"
          >
            <svg
              class="w-4 h-4 sm:w-5 sm:h-5 flex-shrink-0"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
              ></path>
            </svg>
            <span class="text-sm sm:text-base">Refresh</span>
          </button>
        </div>
      </div>

      <!-- After the search input -->
      <div
        class="mt-4 flex items-center justify-between px-4 py-3 bg-white border border-gray-200 rounded-lg shadow-sm"
      >
        <div class="flex items-center">
          <span class="text-sm text-gray-700 mr-2">Tampilkan:</span>
          <select
            v-model="pageSize"
            @change="handlePageSizeChange"
            class="border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
            <option v-for="size in pageSizeOptions" :key="size" :value="size">
              {{ size === -1 ? "Semua" : size }}
            </option>
          </select>
          <span class="text-sm text-gray-700 ml-2">per halaman</span>
        </div>
      </div>

      <!-- Users Table -->
      <div
        class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden"
      >
        <div class="overflow-x-auto">
          <table class="w-full min-w-[800px]">
            <thead class="bg-gray-50 border-b border-gray-200">
              <tr>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  NIK
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Nama
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Email
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Level
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Actions
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-if="loading" class="animate-pulse">
                <td
                  colspan="5"
                  class="px-3 sm:px-6 py-3 sm:py-4 text-center text-gray-500 text-sm"
                >
                  Loading users...
                </td>
              </tr>
              <tr v-else-if="filteredUsers.length === 0">
                <td
                  colspan="5"
                  class="px-3 sm:px-6 py-3 sm:py-4 text-center text-gray-500 text-sm"
                >
                  No users found
                </td>
              </tr>
              <tr
                v-else
                v-for="user in filteredUsers"
                :key="user.nik"
                class="hover:bg-gray-50"
              >
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-xs sm:text-sm text-gray-900">
                    {{ user.nik }}
                  </div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 h-8 w-8 sm:h-10 sm:w-10">
                      <div
                        v-if="user.profile_picture"
                        class="h-8 w-8 sm:h-10 sm:w-10 rounded-full bg-cover bg-center"
                        :style="{
                          backgroundImage: `url(${user.profile_picture})`,
                        }"
                      ></div>
                      <div
                        v-else
                        class="h-8 w-8 sm:h-10 sm:w-10 rounded-full bg-blue-100 flex items-center justify-center"
                      >
                        <span
                          class="text-blue-600 font-medium text-xs sm:text-sm"
                          >{{ getUserInitials(user.name) }}</span
                        >
                      </div>
                    </div>
                    <div class="ml-2 sm:ml-4">
                      <div class="text-xs sm:text-sm font-medium text-gray-900">
                        {{ user.name }}
                      </div>
                    </div>
                  </div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-xs sm:text-sm text-gray-900 break-all">
                    {{ user.email }}
                  </div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-xs sm:text-sm text-gray-900">
                    {{ user.level }}
                  </div>
                </td>
                <td
                  class="px-3 sm:px-6 py-3 sm:py-4 text-right text-sm font-medium"
                >
                  <div class="flex justify-end gap-1 sm:gap-2">
                    <button
                      @click="editUser(user)"
                      class="text-blue-600 hover:text-blue-900 p-1 sm:p-1.5 rounded transition-colors"
                      title="Edit User"
                    >
                      <svg
                        class="w-3 h-3 sm:w-4 sm:h-4"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                        ></path>
                      </svg>
                    </button>
                    <button
                      @click="deleteUser(user)"
                      class="text-red-600 hover:text-red-900 p-1 sm:p-1.5 rounded transition-colors"
                      title="Delete User"
                    >
                      <svg
                        class="w-3 h-3 sm:w-4 sm:h-4"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                        ></path>
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- After the users table -->
      <div
        class="mt-4 flex items-center justify-between px-4 py-3 bg-white border-t border-gray-200 sm:px-6"
      >
        <div class="flex justify-between flex-1 sm:hidden">
          <button
            @click="prevPage"
            :disabled="currentPage === 1"
            class="relative inline-flex items-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
            :class="{ 'opacity-50 cursor-not-allowed': currentPage === 1 }"
          >
            Previous
          </button>
          <button
            @click="nextPage"
            :disabled="currentPage === totalPages"
            class="relative ml-3 inline-flex items-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
            :class="{
              'opacity-50 cursor-not-allowed': currentPage === totalPages,
            }"
          >
            Next
          </button>
        </div>
        <div
          class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between"
        >
          <div>
            <p class="text-sm text-gray-700">
              Showing
              <span class="font-medium">{{
                (currentPage - 1) * pageSize + 1
              }}</span>
              to
              <span class="font-medium">{{
                Math.min(currentPage * pageSize, totalItems)
              }}</span>
              of
              <span class="font-medium">{{ totalItems }}</span>
              results
            </p>
          </div>
          <div>
            <nav
              class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px"
              aria-label="Pagination"
            >
              <button
                @click="prevPage"
                :disabled="currentPage === 1"
                class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
                :class="{ 'opacity-50 cursor-not-allowed': currentPage === 1 }"
              >
                <span class="sr-only">Previous</span>
                <svg
                  class="h-5 w-5"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                  aria-hidden="true"
                >
                  <path
                    fill-rule="evenodd"
                    d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z"
                    clip-rule="evenodd"
                  />
                </svg>
              </button>
              <button
                v-for="page in displayedPages"
                :key="page"
                @click="goToPage(page)"
                :class="[
                  page === currentPage
                    ? 'z-10 bg-blue-50 border-blue-500 text-blue-600'
                    : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50',
                  'relative inline-flex items-center px-4 py-2 border text-sm font-medium',
                ]"
              >
                {{ page }}
              </button>
              <button
                @click="nextPage"
                :disabled="currentPage === totalPages"
                class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
                :class="{
                  'opacity-50 cursor-not-allowed': currentPage === totalPages,
                }"
              >
                <span class="sr-only">Next</span>
                <svg
                  class="h-5 w-5"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                  aria-hidden="true"
                >
                  <path
                    fill-rule="evenodd"
                    d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
                    clip-rule="evenodd"
                  />
                </svg>
              </button>
            </nav>
          </div>
        </div>
      </div>

      <!-- Create/Edit User Modal -->
      <div
        v-if="showCreateModal || showEditModal"
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      >
        <div class="bg-white rounded-lg p-4 sm:p-6 w-full max-w-xl mx-auto">
          <h3 class="text-lg sm:text-xl font-semibold mb-4">
            {{ showCreateModal ? "Tambah User Baru" : "Edit Data User" }}
          </h3>

          <form @submit.prevent="showCreateModal ? createUser() : updateUser()">
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <!-- NIK Field - Only show in create mode -->
              <div v-if="showCreateModal">
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >NIK</label
                >
                <select
                  v-model="userForm.nik"
                  @change="handleEmployeeSelect"
                  :disabled="loadingEmployees"
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
                >
                  <option value="">Pilih Karyawan</option>
                  <option
                    v-for="emp in employees"
                    :key="emp.nik"
                    :value="emp.nik"
                  >
                    {{ emp.nik }} - {{ emp.name }}
                  </option>
                </select>
                <p v-if="loadingEmployees" class="mt-1 text-xs text-gray-500">
                  Memuat karyawan...
                </p>
                <p v-else-if="userForm.nik" class="mt-1 text-xs text-gray-500">
                  Pilih karyawan untuk mengisi NIK dan nama.
                </p>
                <p v-else class="mt-1 text-xs text-gray-500">
                  Pilih karyawan dari daftar untuk menambahkan user.
                </p>
              </div>

              <!-- Name Field -->
              <div :class="{ 'sm:col-span-2': showEditModal }">
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Nama</label
                >
                <input
                  v-model="userForm.name"
                  type="text"
                  required
                  maxlength="25"
                  :disabled="showCreateModal || loading"
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
                  :placeholder="
                    showCreateModal
                      ? 'Nama akan terisi otomatis'
                      : 'Masukkan nama'
                  "
                />
                <p v-if="showCreateModal" class="mt-1 text-xs text-gray-500">
                  Nama akan terisi otomatis saat memilih karyawan
                </p>
              </div>

              <!-- Email Field -->
              <div :class="{ 'sm:col-span-2': showEditModal }">
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Email</label
                >
                <input
                  v-model="userForm.email"
                  type="email"
                  required
                  maxlength="25"
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="nama@email.com"
                />
              </div>

              <!-- Level Field -->
              <div :class="{ 'sm:col-span-2': showEditModal }">
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Level</label
                >
                <select
                  v-model="userForm.level"
                  required
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  <option value="">Pilih Level</option>
                  <option value="1">Level 1</option>
                  <option value="2">Level 2</option>
                  <option value="3">Level 3</option>
                </select>
              </div>

              <!-- Password Field -->
              <div :class="{ 'sm:col-span-2': showEditModal }">
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Password</label
                >
                <div class="relative">
                  <input
                    v-model="userForm.password"
                    :type="showPassword ? 'text' : 'password'"
                    :required="showCreateModal"
                    minlength="8"
                    maxlength="128"
                    class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent pr-10"
                    :placeholder="
                      showEditModal
                        ? 'Kosongkan jika tidak ingin mengubah'
                        : 'Minimal 8 karakter'
                    "
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
                <div v-if="userForm.password" class="mt-2">
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
                <p v-if="showEditModal" class="mt-1 text-xs text-gray-500">
                  Kosongkan field ini jika tidak ingin mengubah password. Isi
                  hanya jika ingin mengganti password.
                </p>
                <p v-else class="mt-1 text-xs text-gray-500">
                  Password harus minimal 8 karakter
                </p>
                <p v-if="passwordError" class="mt-1 text-xs text-red-500">
                  {{ passwordError }}
                </p>
              </div>
            </div>

            <!-- Form Buttons -->
            <div class="flex flex-col sm:flex-row justify-end gap-3 mt-6">
              <button
                type="button"
                @click="closeModal"
                class="px-4 py-2 text-sm sm:text-base text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors order-2 sm:order-1"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-4 py-2 text-sm sm:text-base bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors disabled:opacity-50 order-1 sm:order-2"
              >
                {{
                  submitting
                    ? "Menyimpan..."
                    : showCreateModal
                    ? "Tambah User"
                    : "Update User"
                }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Delete Confirmation Modal -->
      <div
        v-if="showDeleteModal"
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      >
        <div class="bg-white rounded-lg p-4 sm:p-6 w-full max-w-md mx-auto">
          <h3 class="text-lg sm:text-xl font-semibold mb-4 text-red-600">
            Hapus User
          </h3>
          <p class="text-sm sm:text-base text-gray-700 mb-6">
            Apakah Anda yakin ingin menghapus user
            <strong>{{ userToDelete?.name }}</strong
            >? Tindakan ini tidak dapat dibatalkan.
          </p>

          <div class="flex flex-col sm:flex-row justify-end gap-3">
            <button
              @click="showDeleteModal = false"
              class="px-4 py-2 text-sm sm:text-base text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors order-2 sm:order-1"
            >
              Batal
            </button>
            <button
              @click="confirmDelete"
              :disabled="submitting"
              class="px-4 py-2 text-sm sm:text-base bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors disabled:opacity-50 order-1 sm:order-2"
            >
              {{ submitting ? "Menghapus..." : "Hapus User" }}
            </button>
          </div>
        </div>
      </div>

      <!-- Loading Overlay -->
      <div
        v-if="loading"
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      >
        <div class="bg-white rounded-lg p-4 sm:p-6 flex items-center gap-3">
          <div
            class="animate-spin rounded-full h-5 w-5 sm:h-6 sm:w-6 border-b-2 border-primary-600"
          ></div>
          <span class="text-sm sm:text-base text-gray-700">Loading...</span>
        </div>
      </div>
    </div>
  </AuthenticatedLayout>
</template>
