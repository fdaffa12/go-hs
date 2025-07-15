<script setup>
import { ref, computed, onMounted } from "vue";
import { useToast } from "vue-toastification";
import { useAuthStore } from "../stores/auth";
import { userService } from "../services/api";
import AuthenticatedLayout from "../layouts/AuthenticatedLayout.vue";

const authStore = useAuthStore();
const toast = useToast();

// Reactive data
const users = ref([]);
const loading = ref(false);
const submitting = ref(false);
const searchQuery = ref("");

// Modal states
const showCreateModal = ref(false);
const showEditModal = ref(false);
const showDeleteModal = ref(false);

// Form data
const userForm = ref({
  username: "",
  email: "",
  password: "",
});

const userToEdit = ref(null);
const userToDelete = ref(null);

// Computed
const filteredUsers = computed(() => {
  if (!searchQuery.value) return users.value;

  const query = searchQuery.value.toLowerCase();
  return users.value.filter(
    (user) =>
      user.username.toLowerCase().includes(query) ||
      user.email.toLowerCase().includes(query)
  );
});

// Methods
const fetchUsers = async () => {
  loading.value = true;
  try {
    const data = await userService.getAllUsers();
    if (data.success) {
      users.value = data.data || [];
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
    const data = await userService.createUser(userForm.value);
    if (data.success) {
      await fetchUsers();
      closeModal();
      toast.success("Pengguna berhasil dibuat");
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
    username: user.username,
    email: user.email,
    password: "",
  };
  showEditModal.value = true;
};

const updateUser = async () => {
  submitting.value = true;
  try {
    const data = await userService.updateUser(userToEdit.value.id, {
      username: userForm.value.username,
      email: userForm.value.email,
    });
    if (data.success) {
      await fetchUsers();
      closeModal();
      toast.success("Pengguna berhasil diperbarui");
    } else {
      toast.error("Gagal memperbarui pengguna: " + data.message);
    }
  } catch (error) {
    console.error("Error updating user:", error);
    toast.error("Terjadi kesalahan saat memperbarui pengguna");
  } finally {
    submitting.value = false;
  }
};

const deleteUser = (user) => {
  userToDelete.value = user;
  showDeleteModal.value = true;
};

const confirmDelete = async () => {
  submitting.value = true;
  try {
    const data = await userService.deleteUser(userToDelete.value.id);
    if (data.success) {
      await fetchUsers();
      showDeleteModal.value = false;
      toast.success("Pengguna berhasil dihapus");
    } else {
      toast.error("Gagal menghapus pengguna: " + data.message);
    }
  } catch (error) {
    console.error("Error deleting user:", error);
    toast.error("Terjadi kesalahan saat menghapus pengguna");
  } finally {
    submitting.value = false;
  }
};

const openCreateModal = () => {
  userForm.value = {
    username: "",
    email: "",
    password: "",
  };
  showCreateModal.value = true;
};

const closeModal = () => {
  showCreateModal.value = false;
  showEditModal.value = false;
  userForm.value = {
    username: "",
    email: "",
    password: "",
  };
  userToEdit.value = null;
};

const getUserInitials = (username) => {
  return username.substring(0, 2).toUpperCase();
};

const formatDate = (dateString) => {
  if (!dateString) return "N/A";

  // If it's already in YYYY-MM-DD format, just return it
  if (dateString.match(/^\d{4}-\d{2}-\d{2}$/)) {
    return dateString;
  }

  // Try to parse and format the date
  try {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
      return "Invalid Date";
    }
    return date.toLocaleDateString("id-ID", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
    });
  } catch (error) {
    return "Invalid Date";
  }
};

// Lifecycle
onMounted(() => {
  fetchUsers();
});
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
                placeholder="Search users..."
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

      <!-- Users Table -->
      <div
        class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden"
      >
        <div class="overflow-x-auto">
          <table class="w-full min-w-[600px]">
            <thead class="bg-gray-50 border-b border-gray-200">
              <tr>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[180px]"
                >
                  User
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[200px]"
                >
                  Email
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[120px]"
                >
                  Created
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[100px]"
                >
                  Actions
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-if="loading" class="animate-pulse">
                <td
                  colspan="4"
                  class="px-3 sm:px-6 py-3 sm:py-4 text-center text-gray-500 text-sm"
                >
                  Loading users...
                </td>
              </tr>
              <tr v-else-if="filteredUsers.length === 0">
                <td
                  colspan="4"
                  class="px-3 sm:px-6 py-3 sm:py-4 text-center text-gray-500 text-sm"
                >
                  No users found
                </td>
              </tr>
              <tr
                v-else
                v-for="user in filteredUsers"
                :key="user.id"
                class="hover:bg-gray-50"
              >
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 h-8 w-8 sm:h-10 sm:w-10">
                      <div
                        class="h-8 w-8 sm:h-10 sm:w-10 rounded-full bg-blue-100 flex items-center justify-center"
                      >
                        <span
                          class="text-blue-600 font-medium text-xs sm:text-sm"
                          >{{ getUserInitials(user.username) }}</span
                        >
                      </div>
                    </div>
                    <div class="ml-2 sm:ml-4 min-w-0 flex-1">
                      <div
                        class="text-xs sm:text-sm font-medium text-gray-900 truncate"
                      >
                        {{ user.username }}
                      </div>
                      <div class="text-xs text-gray-500">ID: {{ user.id }}</div>
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
                    {{ formatDate(user.created_at) }}
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

      <!-- Create/Edit User Modal -->
      <div
        v-if="showCreateModal || showEditModal"
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      >
        <div class="bg-white rounded-lg p-4 sm:p-6 w-full max-w-md mx-auto">
          <h3 class="text-lg sm:text-xl font-semibold mb-4">
            {{ showCreateModal ? "Create New User" : "Edit User" }}
          </h3>

          <form @submit.prevent="showCreateModal ? createUser() : updateUser()">
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Username</label
                >
                <input
                  v-model="userForm.username"
                  type="text"
                  required
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Enter username"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Email</label
                >
                <input
                  v-model="userForm.email"
                  type="email"
                  required
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Enter email"
                />
              </div>

              <div v-if="showCreateModal">
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Password</label
                >
                <input
                  v-model="userForm.password"
                  type="password"
                  required
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Enter password"
                />
              </div>
            </div>

            <div class="flex flex-col sm:flex-row justify-end gap-3 mt-6">
              <button
                type="button"
                @click="closeModal"
                class="px-4 py-2 text-sm sm:text-base text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors order-2 sm:order-1"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-4 py-2 text-sm sm:text-base bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors disabled:opacity-50 order-1 sm:order-2"
              >
                {{
                  submitting
                    ? "Saving..."
                    : showCreateModal
                    ? "Create User"
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
            Delete User
          </h3>
          <p class="text-sm sm:text-base text-gray-700 mb-6">
            Are you sure you want to delete user
            <strong>{{ userToDelete?.username }}</strong
            >? This action cannot be undone.
          </p>

          <div class="flex flex-col sm:flex-row justify-end gap-3">
            <button
              @click="showDeleteModal = false"
              class="px-4 py-2 text-sm sm:text-base text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors order-2 sm:order-1"
            >
              Cancel
            </button>
            <button
              @click="confirmDelete"
              :disabled="submitting"
              class="px-4 py-2 text-sm sm:text-base bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors disabled:opacity-50 order-1 sm:order-2"
            >
              {{ submitting ? "Deleting..." : "Delete User" }}
            </button>
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
    </div>
  </AuthenticatedLayout>
</template>
