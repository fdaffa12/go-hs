<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { useToast } from "vue-toastification";
import { useAuthStore } from "../stores/auth";
import { styleService, buyerService } from "../services/api";
import AuthenticatedLayout from "../layouts/AuthenticatedLayout.vue";

// Add click-outside directive
const vClickOutside = {
  mounted(el, binding) {
    el._clickOutside = (event) => {
      if (!(el === event.target || el.contains(event.target))) {
        binding.value(event);
      }
    };
    document.addEventListener("click", el._clickOutside);
  },
  unmounted(el) {
    document.removeEventListener("click", el._clickOutside);
  },
};

const authStore = useAuthStore();
const toast = useToast();

// Reactive data
const styles = ref([]);
const buyers = ref([]); // Add buyers list
const buyerSearchQuery = ref(""); // Add buyer search query
const showBuyerDropdown = ref(false); // Control dropdown visibility
const loading = ref(false);
const submitting = ref(false);
const searchQuery = ref("");

// Modal states
const showCreateModal = ref(false);
const showEditModal = ref(false);
const showDeleteModal = ref(false);
const showSoftDeleteModal = ref(false);
const showHardDeleteModal = ref(false);

// Form data - Moved up before any methods that use it
const styleForm = ref({
  style_no: "",
  buyer_short_name: "",
  unit: "",
  t_b: "",
  sub_category: "",
  fabric: "",
  smv_accum: null,
});

const styleToDelete = ref(null);
const styleToHardDelete = ref(null);

// Pagination and page size options
const pageSizeOptions = [10, 20, 50, 100, -1];
const pageSize = ref(10);
const currentPage = ref(1);
const totalItems = ref(0);
const totalPages = ref(0);

// Add debounced search method
const debouncedSearch = ref(null);

const handleSearch = () => {
  if (debouncedSearch.value) {
    clearTimeout(debouncedSearch.value);
  }
  debouncedSearch.value = setTimeout(() => {
    currentPage.value = 1; // Reset to first page when searching
    fetchStyles();
  }, 300);
};

// Watch for search query changes
watch(searchQuery, () => {
  handleSearch();
});

// Methods
const handlePageSizeChange = async () => {
  currentPage.value = 1; // Reset to first page when changing page size
  if (pageSize.value === -1) {
    try {
      const response = await styleService.getAllStyles(1, 1);
      if (response.success) {
        pageSize.value = response.data.total_items;
      }
    } catch (error) {
      console.error("Error getting total count:", error);
      pageSize.value = 100; // Fallback to 100 if error
    }
  }
  fetchStyles();
};

const fetchStyles = async () => {
  loading.value = true;
  try {
    const data = await styleService.getAllStyles(
      currentPage.value,
      pageSize.value === -1 ? 999999 : pageSize.value,
      searchQuery.value
    );
    if (data.success) {
      totalItems.value = data.data.total_items;
      totalPages.value = data.data.total_pages;
      currentPage.value = data.data.current_page;
      pageSize.value = data.data.page_size;
      styles.value = data.data.styles || [];
    } else {
      console.error("Failed to fetch styles:", data.message);
      toast.error("Gagal memuat data style: " + data.message);
    }
  } catch (error) {
    console.error("Error fetching styles:", error);
    toast.error("Terjadi kesalahan saat memuat data style");
  } finally {
    loading.value = false;
  }
};

// Add method to fetch buyers
const fetchBuyers = async (search = "") => {
  try {
    const data = await buyerService.getAllBuyers(1, 100, search);
    if (data.success) {
      buyers.value = data.data.buyers || [];
    } else {
      console.error("Failed to fetch buyers:", data.message);
    }
  } catch (error) {
    console.error("Error fetching buyers:", error);
  }
};

// Add method to handle buyer search
const handleBuyerSearch = () => {
  fetchBuyers(buyerSearchQuery.value);
};

// Modify the selectBuyer method
const selectBuyer = (buyer) => {
  styleForm.value.buyer_short_name = buyer.short_name;
  buyerSearchQuery.value = `${buyer.short_name} - ${buyer.long_name}`; // Update search input to show selected buyer
  showBuyerDropdown.value = false;
};

// Add method to clear buyer selection
const clearBuyerSelection = () => {
  styleForm.value.buyer_short_name = "";
  buyerSearchQuery.value = "";
};

// Add watcher for buyer search
watch(buyerSearchQuery, () => {
  handleBuyerSearch();
});

// Add method to format style number
const formatStyleNo = (styleNo, tb) => {
  if (tb) {
    return `${styleNo}(${tb})`;
  }
  return styleNo;
};

// Add watcher for T/B changes
watch(
  () => styleForm.value.t_b,
  (newTB) => {
    if (styleForm.value.style_no && styleForm.value.unit === "SET") {
      // Remove any existing T/B suffix first
      let baseStyleNo = styleForm.value.style_no.replace(/\(T1\)|\(B1\)$/, "");
      styleForm.value.style_no = formatStyleNo(baseStyleNo, newTB);
    }
  }
);

// Add watcher for unit changes
watch(
  () => styleForm.value.unit,
  (newUnit) => {
    if (newUnit !== "SET") {
      // Remove T/B value and any T/B suffix from style number
      styleForm.value.t_b = "";
      if (styleForm.value.style_no) {
        styleForm.value.style_no = styleForm.value.style_no.replace(
          /\(T1\)|\(B1\)$/,
          ""
        );
      }
    }
  }
);

// Lifecycle
onMounted(() => {
  fetchStyles();
});

// Add pagination methods
const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--;
    fetchStyles();
  }
};

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
    fetchStyles();
  }
};

const goToPage = (page) => {
  if (page !== "..." && page !== currentPage.value) {
    currentPage.value = page;
    fetchStyles();
  }
};

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

// Methods for style operations
const openCreateModal = () => {
  styleForm.value = {
    style_no: "",
    buyer_short_name: "",
    unit: "",
    t_b: "",
    sub_category: "",
    fabric: "",
    smv_accum: null,
  };
  buyerSearchQuery.value = "";
  fetchBuyers();
  showCreateModal.value = true;
};

const editStyle = (style) => {
  // Remove T/B suffix from style number for form
  const baseStyleNo = style.style_no.replace(/\(T1\)|\(B1\)$/, "");
  styleForm.value = {
    ...style,
    style_no: baseStyleNo,
  };

  // Find the buyer to get the long name
  fetchBuyers().then(() => {
    const selectedBuyer = buyers.value.find(
      (b) => b.short_name === style.buyer_short_name
    );
    if (selectedBuyer) {
      buyerSearchQuery.value = `${selectedBuyer.short_name} - ${selectedBuyer.long_name}`;
    } else {
      buyerSearchQuery.value = style.buyer_short_name;
    }
  });
  showEditModal.value = true;
};

const closeModal = () => {
  showCreateModal.value = false;
  showEditModal.value = false;
  styleForm.value = {
    style_no: "",
    buyer_short_name: "",
    unit: "",
    t_b: "",
    sub_category: "",
    fabric: "",
    smv_accum: null,
  };
};

const createStyle = async () => {
  submitting.value = true;
  try {
    // Format the style number if unit is SET
    if (styleForm.value.unit === "SET" && styleForm.value.t_b) {
      let baseStyleNo = styleForm.value.style_no.replace(/\(T1\)|\(B1\)$/, "");
      styleForm.value.style_no = formatStyleNo(
        baseStyleNo,
        styleForm.value.t_b
      );
    }

    const data = await styleService.createStyle(styleForm.value);
    if (data.success) {
      await fetchStyles();
      closeModal();
      toast.success("Style berhasil dibuat");
    } else {
      toast.error("Gagal membuat style: " + data.message);
    }
  } catch (error) {
    console.error("Error creating style:", error);
    toast.error("Terjadi kesalahan saat membuat style");
  } finally {
    submitting.value = false;
  }
};

const updateStyle = async () => {
  submitting.value = true;
  try {
    // Format the style number if unit is SET
    if (styleForm.value.unit === "SET" && styleForm.value.t_b) {
      let baseStyleNo = styleForm.value.style_no.replace(/\(T1\)|\(B1\)$/, "");
      styleForm.value.style_no = formatStyleNo(
        baseStyleNo,
        styleForm.value.t_b
      );
    }

    const data = await styleService.updateStyle(
      styleForm.value.style_no,
      styleForm.value
    );
    if (data.success) {
      await fetchStyles();
      closeModal();
      toast.success("Style berhasil diperbarui");
    } else {
      toast.error("Gagal memperbarui style: " + data.message);
    }
  } catch (error) {
    console.error("Error updating style:", error);
    toast.error("Terjadi kesalahan saat memperbarui style");
  } finally {
    submitting.value = false;
  }
};

const softDeleteStyle = (style) => {
  styleToDelete.value = style;
  showSoftDeleteModal.value = true;
};

const confirmSoftDelete = async () => {
  submitting.value = true;
  try {
    const data = await styleService.deleteStyle(styleToDelete.value.style_no);
    if (data.success) {
      await fetchStyles();
      showSoftDeleteModal.value = false;
      toast.success("Style berhasil dinonaktifkan");
    } else {
      toast.error("Gagal menonaktifkan style: " + data.message);
    }
  } catch (error) {
    console.error("Error soft deleting style:", error);
    toast.error("Terjadi kesalahan saat menonaktifkan style");
  } finally {
    submitting.value = false;
  }
};

const showHardDeleteConfirm = (style) => {
  styleToHardDelete.value = style;
  showHardDeleteModal.value = true;
};

const confirmHardDelete = async () => {
  submitting.value = true;
  try {
    const data = await styleService.hardDeleteStyle(
      styleToHardDelete.value.style_no
    );
    if (data.success) {
      await fetchStyles();
      showHardDeleteModal.value = false;
      toast.success("Style berhasil dihapus secara permanen");
    } else {
      toast.error("Gagal menghapus style: " + data.message);
    }
  } catch (error) {
    console.error("Error hard deleting style:", error);
    toast.error("Terjadi kesalahan saat menghapus style");
  } finally {
    submitting.value = false;
  }
};

const activateStyle = async (style) => {
  submitting.value = true;
  try {
    const data = await styleService.activateStyle(style.style_no);
    if (data.success) {
      await fetchStyles();
      toast.success("Style berhasil diaktifkan");
    } else {
      toast.error("Gagal mengaktifkan style: " + data.message);
    }
  } catch (error) {
    console.error("Error activating style:", error);
    toast.error("Terjadi kesalahan saat mengaktifkan style");
  } finally {
    submitting.value = false;
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
              Manajemen Style
            </h1>
            <p class="mt-1 text-sm text-gray-600">Kelola style dalam sistem</p>
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
            <span class="hidden sm:inline">Tambah Style</span>
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
                placeholder="Cari style..."
                class="w-full pl-9 sm:pl-10 pr-4 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
          <button
            @click="fetchStyles"
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

      <!-- Page Size Selector -->
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

      <!-- Styles Table -->
      <div
        class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden"
      >
        <div class="overflow-x-auto">
          <table class="w-full min-w-[800px]">
            <thead class="bg-gray-50 border-b border-gray-200">
              <tr>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[120px]"
                >
                  Style No
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[120px]"
                >
                  Buyer
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[80px]"
                >
                  Unit
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[80px]"
                >
                  T/B
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[120px]"
                >
                  Sub Category
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[120px]"
                >
                  Fabric
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[100px]"
                >
                  SMV
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[120px]"
                >
                  E-Style No
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[100px]"
                >
                  Status
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
                  colspan="10"
                  class="px-3 sm:px-6 py-3 sm:py-4 text-center text-gray-500 text-sm"
                >
                  Memuat style...
                </td>
              </tr>
              <tr v-else-if="styles.length === 0">
                <td
                  colspan="10"
                  class="px-3 sm:px-6 py-3 sm:py-4 text-center text-gray-500 text-sm"
                >
                  Tidak ada style ditemukan
                </td>
              </tr>
              <tr
                v-else
                v-for="style in styles"
                :key="style.style_no"
                class="hover:bg-gray-50"
              >
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-sm font-medium text-gray-900">
                    {{ style.style_no }}
                  </div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-sm text-gray-900">
                    {{ style.buyer_short_name }}
                  </div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-sm text-gray-900">{{ style.unit }}</div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-sm text-gray-900">
                    {{ style.t_b || "-" }}
                  </div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-sm text-gray-900">
                    {{ style.sub_category || "-" }}
                  </div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-sm text-gray-900">
                    {{ style.fabric || "-" }}
                  </div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-sm text-gray-900">
                    {{ style.smv_accum || "-" }}
                  </div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-sm text-gray-900">
                    {{ style.e_style_no || "-" }}
                  </div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <span
                    :class="[
                      'px-2 py-1 text-xs rounded-full',
                      style.delete_status
                        ? 'bg-red-100 text-red-800'
                        : 'bg-green-100 text-green-800',
                    ]"
                  >
                    {{ style.delete_status ? "Non-Aktif" : "Aktif" }}
                  </span>
                </td>
                <td
                  class="px-3 sm:px-6 py-3 sm:py-4 text-right text-sm font-medium"
                >
                  <div class="flex justify-end gap-1 sm:gap-2">
                    <button
                      v-if="!style.delete_status"
                      @click="editStyle(style)"
                      class="text-blue-600 hover:text-blue-900 p-1 sm:p-1.5 rounded transition-colors"
                      title="Edit Style"
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
                      v-if="!style.delete_status"
                      @click="softDeleteStyle(style)"
                      class="text-yellow-600 hover:text-yellow-900 p-1 sm:p-1.5 rounded transition-colors"
                      title="Non-aktifkan Style"
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
                      v-if="style.delete_status"
                      @click="activateStyle(style)"
                      class="text-green-600 hover:text-green-900 p-1 sm:p-1.5 rounded transition-colors"
                      title="Aktifkan Style"
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
                      @click="showHardDeleteConfirm(style)"
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

      <!-- Modals will be added in the next part -->

      <!-- Create/Edit Style Modal -->
      <div
        v-if="showCreateModal || showEditModal"
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      >
        <div class="bg-white rounded-lg p-4 sm:p-6 w-full max-w-md mx-auto">
          <h3 class="text-lg sm:text-xl font-semibold mb-4">
            {{ showCreateModal ? "Buat Style Baru" : "Edit Style" }}
          </h3>

          <form
            @submit.prevent="showCreateModal ? createStyle() : updateStyle()"
          >
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Style No</label
                >
                <input
                  v-model="styleForm.style_no"
                  type="text"
                  required
                  maxlength="10"
                  :disabled="!showCreateModal"
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
                  placeholder="Masukkan nomor style"
                />
              </div>

              <div class="relative">
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Buyer</label
                >
                <div
                  class="relative"
                  v-click-outside="() => (showBuyerDropdown = false)"
                >
                  <div class="relative">
                    <input
                      v-model="buyerSearchQuery"
                      type="text"
                      required
                      @focus="showBuyerDropdown = true"
                      @input="styleForm.buyer_short_name = ''"
                      class="w-full px-3 py-2 pr-8 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      placeholder="Cari buyer..."
                    />
                    <button
                      v-if="buyerSearchQuery"
                      @click="clearBuyerSelection"
                      class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                      type="button"
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
                          d="M6 18L18 6M6 6l12 12"
                        />
                      </svg>
                    </button>
                  </div>
                  <div
                    v-if="showBuyerDropdown"
                    class="absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-60 overflow-y-auto"
                  >
                    <div
                      v-if="buyers.length === 0"
                      class="px-4 py-2 text-sm text-gray-500"
                    >
                      Tidak ada buyer ditemukan
                    </div>
                    <div
                      v-else
                      v-for="buyer in buyers"
                      :key="buyer.short_name"
                      @click="selectBuyer(buyer)"
                      class="px-4 py-2 text-sm hover:bg-gray-100 cursor-pointer"
                      :class="{
                        'bg-blue-50':
                          buyer.short_name === styleForm.buyer_short_name,
                      }"
                    >
                      {{ buyer.short_name }} - {{ buyer.long_name }}
                    </div>
                  </div>
                </div>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Unit</label
                >
                <select
                  v-model="styleForm.unit"
                  required
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  <option value="">Pilih unit</option>
                  <option value="SET">SET</option>
                  <option value="PCS">PCS</option>
                </select>
              </div>

              <div v-if="styleForm.unit === 'SET'">
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >T/B</label
                >
                <select
                  v-model="styleForm.t_b"
                  required
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  <option value="">Pilih T/B</option>
                  <option value="T1">T1</option>
                  <option value="B1">B1</option>
                </select>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Sub Category</label
                >
                <input
                  v-model="styleForm.sub_category"
                  type="text"
                  maxlength="50"
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Masukkan sub kategori"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >Fabric</label
                >
                <input
                  v-model="styleForm.fabric"
                  type="text"
                  maxlength="50"
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Masukkan jenis fabric"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1"
                  >SMV</label
                >
                <input
                  v-model="styleForm.smv_accum"
                  type="number"
                  step="0.01"
                  min="0"
                  class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Masukkan nilai SMV"
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
                    ? "Buat Style"
                    : "Update Style"
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
            Hapus Style
          </h3>
          <p class="text-sm sm:text-base text-gray-700 mb-6">
            Apakah Anda yakin ingin menghapus style
            <strong>{{ styleToDelete?.style_no }}</strong
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
              {{ submitting ? "Menghapus..." : "Hapus Style" }}
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
            Non-aktifkan Style
          </h3>
          <p class="text-sm sm:text-base text-gray-700 mb-6">
            Apakah Anda yakin ingin menonaktifkan style
            <strong>{{ styleToDelete?.style_no }}</strong
            >? Style yang dinonaktifkan masih dapat dilihat dalam sistem.
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
            Hapus Permanen Style
          </h3>
          <p class="text-sm sm:text-base text-gray-700 mb-6">
            Apakah Anda yakin ingin menghapus style
            <strong>{{ styleToHardDelete?.style_no }}</strong>
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
