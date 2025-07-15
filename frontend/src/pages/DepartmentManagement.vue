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
              Manajemen Departemen
            </h1>
            <p class="mt-1 text-sm text-gray-600">
              Kelola departemen dalam sistem
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
            <span class="hidden sm:inline">Tambah Departemen</span>
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
                placeholder="Cari departemen..."
                class="w-full pl-9 sm:pl-10 pr-4 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
          <button
            @click="fetchDepartments"
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

      <!-- Departments Table -->
      <div
        class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden"
      >
        <div class="overflow-x-auto">
          <table class="w-full min-w-[600px]">
            <thead class="bg-gray-50 border-b border-gray-200">
              <tr>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[120px]"
                >
                  Kode
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[200px]"
                >
                  Nama Departemen
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[120px]"
                >
                  Dibuat
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[100px]"
                >
                  Aksi
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-if="loading" class="animate-pulse">
                <td
                  colspan="4"
                  class="px-3 sm:px-6 py-3 sm:py-4 text-center text-gray-500 text-sm"
                >
                  Memuat departemen...
                </td>
              </tr>
              <tr v-else-if="filteredDepartments.length === 0">
                <td
                  colspan="4"
                  class="px-3 sm:px-6 py-3 sm:py-4 text-center text-gray-500 text-sm"
                >
                  Tidak ada departemen ditemukan
                </td>
              </tr>
              <tr
                v-else
                v-for="dept in filteredDepartments"
                :key="dept.short_name"
                class="hover:bg-gray-50"
              >
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-sm font-medium text-gray-900">
                    {{ dept.short_name }}
                  </div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-sm text-gray-900">{{ dept.long_name }}</div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-xs sm:text-sm text-gray-900">
                    {{ formatDate(dept.created_at) }}
                  </div>
                </td>
                <td
                  class="px-3 sm:px-6 py-3 sm:py-4 text-right text-sm font-medium"
                >
                  <div class="flex justify-end gap-1 sm:gap-2">
                    <button
                      @click="editDepartment(dept)"
                      class="text-blue-600 hover:text-blue-900 p-1 sm:p-1.5 rounded transition-colors"
                      title="Edit Departemen"
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
                      @click="deleteDepartment(dept)"
                      class="text-red-600 hover:text-red-900 p-1 sm:p-1.5 rounded transition-colors"
                      title="Hapus Departemen"
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

      <!-- Create/Edit Department Modal -->
      <div
        v-if="showCreateModal || showEditModal"
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      >
        <div class="bg-white rounded-lg p-4 sm:p-6 w-full max-w-md mx-auto">
          <h3 class="text-lg sm:text-xl font-semibold mb-4">
            {{ showCreateModal ? "Buat Departemen Baru" : "Edit Departemen" }}
          </h3>

          <form
            @submit.prevent="
              showCreateModal ? createDepartment() : updateDepartment()
            "
          >
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Kode Departemen</label
                >
                <input
                  v-model="departmentForm.short_name"
                  type="text"
                  required
                  maxlength="5"
                  :disabled="!showCreateModal"
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
                  placeholder="Masukkan kode departemen"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Nama Departemen</label
                >
                <input
                  v-model="departmentForm.long_name"
                  type="text"
                  required
                  maxlength="50"
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Masukkan nama departemen"
                />
              </div>
            </div>

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
                    ? "Buat Departemen"
                    : "Update Departemen"
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
            Hapus Departemen
          </h3>
          <p class="text-sm sm:text-base text-gray-700 mb-6">
            Apakah Anda yakin ingin menghapus departemen
            <strong>{{ departmentToDelete?.long_name }}</strong
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
              {{ submitting ? "Menghapus..." : "Hapus Departemen" }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </AuthenticatedLayout>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useToast } from "vue-toastification";
import { useAuthStore } from "../stores/auth";
import { departmentService } from "../services/api";
import AuthenticatedLayout from "../layouts/AuthenticatedLayout.vue";

const authStore = useAuthStore();
const toast = useToast();

// Reactive data
const departments = ref([]);
const loading = ref(false);
const submitting = ref(false);
const searchQuery = ref("");

// Modal states
const showCreateModal = ref(false);
const showEditModal = ref(false);
const showDeleteModal = ref(false);

// Form data
const departmentForm = ref({
  short_name: "",
  long_name: "",
});

const departmentToDelete = ref(null);

// Computed
const filteredDepartments = computed(() => {
  if (!searchQuery.value) return departments.value;

  const query = searchQuery.value.toLowerCase();
  return departments.value.filter(
    (dept) =>
      dept.short_name.toLowerCase().includes(query) ||
      dept.long_name.toLowerCase().includes(query)
  );
});

// Methods
const fetchDepartments = async () => {
  loading.value = true;
  try {
    const data = await departmentService.getAllDepartments();
    if (data.success) {
      departments.value = data.data || [];
    } else {
      console.error("Failed to fetch departments:", data.message);
      toast.error("Gagal memuat data departemen: " + data.message);
    }
  } catch (error) {
    console.error("Error fetching departments:", error);
    toast.error("Terjadi kesalahan saat memuat data departemen");
  } finally {
    loading.value = false;
  }
};

const createDepartment = async () => {
  submitting.value = true;
  try {
    const data = await departmentService.createDepartment(departmentForm.value);
    if (data.success) {
      await fetchDepartments();
      closeModal();
      toast.success("Departemen berhasil dibuat");
    } else {
      toast.error("Gagal membuat departemen: " + data.message);
    }
  } catch (error) {
    console.error("Error creating department:", error);
    toast.error("Terjadi kesalahan saat membuat departemen");
  } finally {
    submitting.value = false;
  }
};

const editDepartment = (dept) => {
  departmentForm.value = {
    short_name: dept.short_name,
    long_name: dept.long_name,
  };
  showEditModal.value = true;
};

const updateDepartment = async () => {
  submitting.value = true;
  try {
    const data = await departmentService.updateDepartment(
      departmentForm.value.short_name,
      departmentForm.value
    );
    if (data.success) {
      await fetchDepartments();
      closeModal();
      toast.success("Departemen berhasil diperbarui");
    } else {
      toast.error("Gagal memperbarui departemen: " + data.message);
    }
  } catch (error) {
    console.error("Error updating department:", error);
    toast.error("Terjadi kesalahan saat memperbarui departemen");
  } finally {
    submitting.value = false;
  }
};

const deleteDepartment = (dept) => {
  departmentToDelete.value = dept;
  showDeleteModal.value = true;
};

const confirmDelete = async () => {
  submitting.value = true;
  try {
    const data = await departmentService.deleteDepartment(
      departmentToDelete.value.short_name
    );
    if (data.success) {
      await fetchDepartments();
      showDeleteModal.value = false;
      toast.success("Departemen berhasil dihapus");
    } else {
      toast.error("Gagal menghapus departemen: " + data.message);
    }
  } catch (error) {
    console.error("Error deleting department:", error);
    toast.error("Terjadi kesalahan saat menghapus departemen");
  } finally {
    submitting.value = false;
  }
};

const openCreateModal = () => {
  departmentForm.value = {
    short_name: "",
    long_name: "",
  };
  showCreateModal.value = true;
};

const closeModal = () => {
  showCreateModal.value = false;
  showEditModal.value = false;
  departmentForm.value = {
    short_name: "",
    long_name: "",
  };
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
    });
  } catch (error) {
    return "Invalid Date";
  }
};

// Lifecycle
onMounted(() => {
  fetchDepartments();
});
</script>
