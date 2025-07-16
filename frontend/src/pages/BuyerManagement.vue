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
              Manajemen Buyer
            </h1>
            <p class="mt-1 text-sm text-gray-600">Kelola buyer dalam sistem</p>
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
            <span class="hidden sm:inline">Tambah Buyer</span>
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
                placeholder="Cari buyer..."
                class="w-full pl-9 sm:pl-10 pr-4 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
          <button
            @click="fetchBuyers"
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

      <!-- Add page size selector -->
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

      <!-- Buyers Table -->
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
                  Nama Buyer
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[100px]"
                >
                  Status
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
                  colspan="5"
                  class="px-3 sm:px-6 py-3 sm:py-4 text-center text-gray-500 text-sm"
                >
                  Memuat buyer...
                </td>
              </tr>
              <tr v-else-if="buyers.length === 0">
                <td
                  colspan="5"
                  class="px-3 sm:px-6 py-3 sm:py-4 text-center text-gray-500 text-sm"
                >
                  Tidak ada buyer ditemukan
                </td>
              </tr>
              <tr
                v-else
                v-for="buyer in buyers"
                :key="buyer.short_name"
                class="hover:bg-gray-50"
              >
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-sm font-medium text-gray-900">
                    {{ buyer.short_name }}
                  </div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-sm text-gray-900">{{ buyer.long_name }}</div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <span
                    :class="[
                      'px-2 py-1 text-xs rounded-full',
                      buyer.delete_status
                        ? 'bg-red-100 text-red-800'
                        : 'bg-green-100 text-green-800',
                    ]"
                  >
                    {{ buyer.delete_status ? "Non-Aktif" : "Aktif" }}
                  </span>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-xs sm:text-sm text-gray-900">
                    {{ formatDate(buyer.created_at) }}
                  </div>
                </td>
                <td
                  class="px-3 sm:px-6 py-3 sm:py-4 text-right text-sm font-medium"
                >
                  <div class="flex justify-end gap-1 sm:gap-2">
                    <button
                      v-if="!buyer.delete_status"
                      @click="editBuyer(buyer)"
                      class="text-blue-600 hover:text-blue-900 p-1 sm:p-1.5 rounded transition-colors"
                      title="Edit Buyer"
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
                      v-if="!buyer.delete_status"
                      @click="softDeleteBuyer(buyer)"
                      class="text-yellow-600 hover:text-yellow-900 p-1 sm:p-1.5 rounded transition-colors"
                      title="Non-aktifkan Buyer"
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
                          d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"
                        ></path>
                      </svg>
                    </button>
                    <button
                      v-if="buyer.delete_status"
                      @click="activateBuyer(buyer)"
                      class="text-green-600 hover:text-green-900 p-1 sm:p-1.5 rounded transition-colors"
                      title="Aktifkan Buyer"
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
                          d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                        ></path>
                      </svg>
                    </button>
                    <button
                      @click="showHardDeleteConfirm(buyer)"
                      class="text-red-600 hover:text-red-900 p-1 sm:p-1.5 rounded transition-colors"
                      title="Hapus Permanen"
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

      <!-- Pagination -->
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

      <!-- Create/Edit Buyer Modal -->
      <div
        v-if="showCreateModal || showEditModal"
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      >
        <div class="bg-white rounded-lg p-4 sm:p-6 w-full max-w-md mx-auto">
          <h3 class="text-lg sm:text-xl font-semibold mb-4">
            {{ showCreateModal ? "Buat Buyer Baru" : "Edit Buyer" }}
          </h3>

          <form
            @submit.prevent="showCreateModal ? createBuyer() : updateBuyer()"
          >
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Kode Buyer</label
                >
                <input
                  v-model="buyerForm.short_name"
                  type="text"
                  required
                  maxlength="5"
                  :disabled="!showCreateModal"
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
                  placeholder="Masukkan kode buyer"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Nama Buyer</label
                >
                <input
                  v-model="buyerForm.long_name"
                  type="text"
                  required
                  maxlength="50"
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Masukkan nama buyer"
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
                    ? "Buat Buyer"
                    : "Update Buyer"
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
            Hapus Buyer
          </h3>
          <p class="text-sm sm:text-base text-gray-700 mb-6">
            Apakah Anda yakin ingin menghapus buyer
            <strong>{{ buyerToDelete?.long_name }}</strong
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
              {{ submitting ? "Menghapus..." : "Hapus Buyer" }}
            </button>
          </div>
        </div>
      </div>

      <!-- Soft Delete Confirmation Modal -->
      <div
        v-if="showSoftDeleteModal"
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      >
        <div class="bg-white rounded-lg p-4 sm:p-6 w-full max-w-md mx-auto">
          <h3 class="text-lg sm:text-xl font-semibold mb-4 text-yellow-600">
            Non-aktifkan Buyer
          </h3>
          <p class="text-sm sm:text-base text-gray-700 mb-6">
            Apakah Anda yakin ingin menonaktifkan buyer
            <strong>{{ buyerToDelete?.long_name }}</strong
            >? Buyer yang dinonaktifkan masih dapat dilihat dalam sistem.
          </p>

          <div class="flex flex-col sm:flex-row justify-end gap-3">
            <button
              @click="showSoftDeleteModal = false"
              class="px-4 py-2 text-sm sm:text-base text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors order-2 sm:order-1"
            >
              Batal
            </button>
            <button
              @click="confirmSoftDelete"
              :disabled="submitting"
              class="px-4 py-2 text-sm sm:text-base bg-yellow-600 hover:bg-yellow-700 text-white rounded-lg transition-colors disabled:opacity-50 order-1 sm:order-2"
            >
              {{ submitting ? "Menonaktifkan..." : "Non-aktifkan" }}
            </button>
          </div>
        </div>
      </div>

      <!-- Hard Delete Confirmation Modal -->
      <div
        v-if="showHardDeleteModal"
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      >
        <div class="bg-white rounded-lg p-4 sm:p-6 w-full max-w-md mx-auto">
          <h3 class="text-lg sm:text-xl font-semibold mb-4 text-red-600">
            Hapus Permanen Buyer
          </h3>
          <p class="text-sm sm:text-base text-gray-700 mb-6">
            Apakah Anda yakin ingin menghapus buyer
            <strong>{{ buyerToHardDelete?.long_name }}</strong>
            secara permanen? Tindakan ini tidak dapat dibatalkan dan data akan
            dihapus dari sistem.
          </p>

          <div class="flex flex-col sm:flex-row justify-end gap-3">
            <button
              @click="showHardDeleteModal = false"
              class="px-4 py-2 text-sm sm:text-base text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors order-2 sm:order-1"
            >
              Batal
            </button>
            <button
              @click="confirmHardDelete"
              :disabled="submitting"
              class="px-4 py-2 text-sm sm:text-base bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors disabled:opacity-50 order-1 sm:order-2"
            >
              {{ submitting ? "Menghapus..." : "Hapus Permanen" }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </AuthenticatedLayout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { useToast } from "vue-toastification";
import { useAuthStore } from "../stores/auth";
import { buyerService } from "../services/api";
import AuthenticatedLayout from "../layouts/AuthenticatedLayout.vue";

const authStore = useAuthStore();
const toast = useToast();

// Reactive data
const buyers = ref([]);
const loading = ref(false);
const submitting = ref(false);
const searchQuery = ref("");

// Pagination and page size options
const pageSizeOptions = [10, 20, 50, 100, -1]; // -1 represents "All"
const pageSize = ref(10);
const currentPage = ref(1);
const totalItems = ref(0);
const totalPages = ref(0);

// Modal states
const showCreateModal = ref(false);
const showEditModal = ref(false);
const showDeleteModal = ref(false);
const showSoftDeleteModal = ref(false);
const showHardDeleteModal = ref(false);

// Form data
const buyerForm = ref({
  short_name: "",
  long_name: "",
});

const buyerToDelete = ref(null);
const buyerToHardDelete = ref(null);

// Computed
const filteredBuyers = computed(() => buyers.value);

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

// Add debounced search method
const debouncedSearch = ref(null);

const handleSearch = () => {
  if (debouncedSearch.value) {
    clearTimeout(debouncedSearch.value);
  }
  debouncedSearch.value = setTimeout(() => {
    currentPage.value = 1; // Reset to first page when searching
    fetchBuyers();
  }, 300);
};

// Methods
const handlePageSizeChange = async () => {
  currentPage.value = 1; // Reset to first page when changing page size
  if (pageSize.value === -1) {
    // If "All" is selected, get total count first
    try {
      const response = await buyerService.getAllBuyers(1, 1);
      if (response.success) {
        pageSize.value = response.data.total_items;
      }
    } catch (error) {
      console.error("Error getting total count:", error);
      pageSize.value = 100; // Fallback to 100 if error
    }
  }
  fetchBuyers();
};

// Watch for search query changes
watch(searchQuery, () => {
  handleSearch();
});

const fetchBuyers = async () => {
  loading.value = true;
  try {
    const data = await buyerService.getAllBuyers(
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

      buyers.value = data.data.buyers || [];
    } else {
      console.error("Failed to fetch buyers:", data.message);
      toast.error("Gagal memuat data buyer: " + data.message);
    }
  } catch (error) {
    console.error("Error fetching buyers:", error);
    toast.error("Terjadi kesalahan saat memuat data buyer");
  } finally {
    loading.value = false;
  }
};

// Add pagination methods
const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--;
    fetchBuyers();
  }
};

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
    fetchBuyers();
  }
};

const goToPage = (page) => {
  if (page !== "..." && page !== currentPage.value) {
    currentPage.value = page;
    fetchBuyers();
  }
};

const createBuyer = async () => {
  submitting.value = true;
  try {
    const data = await buyerService.createBuyer(buyerForm.value);
    if (data.success) {
      await fetchBuyers();
      closeModal();
      toast.success("Buyer berhasil dibuat");
    } else {
      toast.error("Gagal membuat buyer: " + data.message);
    }
  } catch (error) {
    console.error("Error creating buyer:", error);
    toast.error("Terjadi kesalahan saat membuat buyer");
  } finally {
    submitting.value = false;
  }
};

const editBuyer = (buyer) => {
  buyerForm.value = {
    short_name: buyer.short_name,
    long_name: buyer.long_name,
  };
  showEditModal.value = true;
};

const updateBuyer = async () => {
  submitting.value = true;
  try {
    const data = await buyerService.updateBuyer(
      buyerForm.value.short_name,
      buyerForm.value
    );
    if (data.success) {
      await fetchBuyers();
      closeModal();
      toast.success("Buyer berhasil diperbarui");
    } else {
      toast.error("Gagal memperbarui buyer: " + data.message);
    }
  } catch (error) {
    console.error("Error updating buyer:", error);
    toast.error("Terjadi kesalahan saat memperbarui buyer");
  } finally {
    submitting.value = false;
  }
};

const deleteBuyer = (buyer) => {
  buyerToDelete.value = buyer;
  showDeleteModal.value = true;
};

const confirmDelete = async () => {
  submitting.value = true;
  try {
    const data = await buyerService.deleteBuyer(buyerToDelete.value.short_name);
    if (data.success) {
      await fetchBuyers();
      showDeleteModal.value = false;
      toast.success("Buyer berhasil dihapus");
    } else {
      toast.error("Gagal menghapus buyer: " + data.message);
    }
  } catch (error) {
    console.error("Error deleting buyer:", error);
    toast.error("Terjadi kesalahan saat menghapus buyer");
  } finally {
    submitting.value = false;
  }
};

const openCreateModal = () => {
  buyerForm.value = {
    short_name: "",
    long_name: "",
  };
  showCreateModal.value = true;
};

const closeModal = () => {
  showCreateModal.value = false;
  showEditModal.value = false;
  buyerForm.value = {
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

// Methods for soft delete
const softDeleteBuyer = (buyer) => {
  buyerToDelete.value = buyer;
  showSoftDeleteModal.value = true;
};

const confirmSoftDelete = async () => {
  submitting.value = true;
  try {
    const data = await buyerService.deleteBuyer(buyerToDelete.value.short_name);
    if (data.success) {
      await fetchBuyers();
      showSoftDeleteModal.value = false;
      toast.success("Buyer berhasil dinonaktifkan");
    } else {
      toast.error("Gagal menonaktifkan buyer: " + data.message);
    }
  } catch (error) {
    console.error("Error soft deleting buyer:", error);
    toast.error("Terjadi kesalahan saat menonaktifkan buyer");
  } finally {
    submitting.value = false;
  }
};

// Methods for hard delete
const showHardDeleteConfirm = (buyer) => {
  buyerToHardDelete.value = buyer;
  showHardDeleteModal.value = true;
};

const confirmHardDelete = async () => {
  submitting.value = true;
  try {
    const data = await buyerService.hardDeleteBuyer(
      buyerToHardDelete.value.short_name
    );
    if (data.success) {
      await fetchBuyers();
      showHardDeleteModal.value = false;
      toast.success("Buyer berhasil dihapus secara permanen");
    } else {
      toast.error("Gagal menghapus buyer: " + data.message);
    }
  } catch (error) {
    console.error("Error hard deleting buyer:", error);
    toast.error("Terjadi kesalahan saat menghapus buyer");
  } finally {
    submitting.value = false;
  }
};

// Methods for activation
const activateBuyer = async (buyer) => {
  submitting.value = true;
  try {
    const data = await buyerService.activateBuyer(buyer.short_name);
    if (data.success) {
      await fetchBuyers();
      toast.success("Buyer berhasil diaktifkan");
    } else {
      toast.error("Gagal mengaktifkan buyer: " + data.message);
    }
  } catch (error) {
    console.error("Error activating buyer:", error);
    toast.error("Terjadi kesalahan saat mengaktifkan buyer");
  } finally {
    submitting.value = false;
  }
};

// Lifecycle
onMounted(() => {
  fetchBuyers();
});
</script>
