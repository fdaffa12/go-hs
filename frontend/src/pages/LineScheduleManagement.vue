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
              Manajemen Line Schedule
            </h1>
            <p class="mt-1 text-sm text-gray-600">
              Kelola jadwal line produksi
            </p>
          </div>
          <div class="flex flex-wrap gap-2">
            <button
              @click="addNewRow"
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
              <span class="hidden sm:inline">Tambah Schedule</span>
              <span class="sm:hidden">Tambah</span>
            </button>
            <button
              @click="exportToExcel"
              class="btn btn-secondary flex items-center justify-center gap-2 w-full sm:w-auto"
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
                  d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                ></path>
              </svg>
              <span>Export</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Filters -->
      <div
        class="bg-white rounded-lg shadow-sm border border-gray-200 p-4 sm:p-6"
      >
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Factory</label
            >
            <select
              v-model="filters.factory"
              class="form-select w-full border border-gray-300"
              @change="handleFilterChange"
            >
              <option value="">Pilih Factory</option>
              <option value="F1">F1</option>
              <option value="F2">F2</option>
              <option value="F3">F3</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Tanggal</label
            >
            <input
              type="date"
              v-model="filters.date"
              class="form-input w-full"
              @change="handleFilterChange"
            />
          </div>
        </div>
        <div class="mt-4 flex justify-end">
          <button
            @click="applyFilters"
            class="btn btn-primary"
            :disabled="!filters.date"
          >
            Filter Data
          </button>
        </div>
      </div>

      <!-- Table -->
      <div
        class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden"
      >
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  ID Registrasi
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Type
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Line
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Buyer
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Style No
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  MP
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Start Date
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Working Days
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Status
                </th>
                <th
                  class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Actions
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <!-- New Rows Forms -->
              <tr
                v-for="(row, index) in newRows"
                :key="index"
                class="bg-blue-50"
              >
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  <em>Auto-generated</em>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <select
                    v-model="row.type"
                    class="form-select w-full"
                    @change="handleNewRowTypeChange(row)"
                  >
                    <option value="assembly">Assembly</option>
                    <option value="component">Component</option>
                  </select>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <select
                    v-if="row.type === 'assembly'"
                    v-model="row.line"
                    class="form-select w-full"
                  >
                    <option value="">Pilih Line</option>
                    <option v-for="line in lines" :key="line" :value="line">
                      {{ line }}
                    </option>
                  </select>
                  <input
                    v-else
                    v-model="row.line"
                    class="form-input w-full"
                    disabled
                    value="AREA"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <select
                    v-model="row.buyer_short_name"
                    class="form-select w-full"
                    @change="handleBuyerChange(row)"
                  >
                    <option value="">Pilih Buyer</option>
                    <option
                      v-for="buyer in buyers"
                      :key="buyer.short_name"
                      :value="buyer.short_name"
                    >
                      {{ buyer.long_name }}
                    </option>
                  </select>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <select
                    v-model="row.style_no"
                    class="form-select w-full"
                    :disabled="!row.buyer_short_name"
                  >
                    <option value="">Pilih Style</option>
                    <option
                      v-for="style in filteredStyles"
                      :key="style.style_no"
                      :value="style.style_no"
                    >
                      {{ style.style_no }}
                    </option>
                  </select>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <input
                    v-model.number="row.number_of_mp"
                    type="number"
                    class="form-input w-full"
                    placeholder="MP"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <input
                    v-model="row.start_date"
                    type="date"
                    class="form-input w-full"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <input
                    v-model.number="row.working_day"
                    type="number"
                    class="form-input w-full"
                    placeholder="Days"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    class="px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-800"
                    >New</span
                  >
                </td>
                <td
                  class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                >
                  <div class="flex justify-end space-x-2">
                    <button
                      @click="saveNewRow(row, index)"
                      class="text-green-600 hover:text-green-900"
                      title="Save"
                    >
                      <svg
                        class="w-5 h-5"
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
                    </button>
                    <button
                      @click="cancelNewRow(index)"
                      class="text-red-600 hover:text-red-900"
                      title="Cancel"
                    >
                      <svg
                        class="w-5 h-5"
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
                </td>
              </tr>

              <!-- Existing Rows -->
              <tr v-if="!isFilterApplied">
                <td colspan="10" class="px-6 py-4 text-center text-gray-500">
                  Silakan pilih filter terlebih dahulu
                </td>
              </tr>
              <tr v-else-if="loading" class="animate-pulse">
                <td colspan="10" class="px-6 py-4 text-center text-gray-500">
                  Loading...
                </td>
              </tr>
              <tr v-else-if="schedules.length === 0">
                <td colspan="10" class="px-6 py-4 text-center text-gray-500">
                  Tidak ada data yang sesuai dengan filter yang dipilih
                </td>
              </tr>
              <tr
                v-for="schedule in schedules"
                :key="schedule.row_id"
                class="hover:bg-gray-50"
              >
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ schedule.id_registrasi }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ schedule.type }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ schedule.line }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ schedule.buyer_short_name }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ schedule.style_no }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ schedule.number_of_mp }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ formatDate(schedule.start_date) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ schedule.working_day }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    :class="[
                      'px-2 py-1 text-xs rounded-full',
                      schedule.delete_status
                        ? 'bg-red-100 text-red-800'
                        : 'bg-green-100 text-green-800',
                    ]"
                  >
                    {{ schedule.delete_status ? "Non-Aktif" : "Aktif" }}
                  </span>
                </td>
                <td
                  class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                >
                  <div class="flex justify-end space-x-2">
                    <button
                      v-if="!schedule.delete_status"
                      @click="editSchedule(schedule)"
                      class="text-blue-600 hover:text-blue-900"
                      title="Edit"
                    >
                      <svg
                        class="w-5 h-5"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                        />
                      </svg>
                    </button>
                    <button
                      v-if="!schedule.delete_status"
                      @click="showSoftDeleteConfirm(schedule)"
                      class="text-yellow-600 hover:text-yellow-900"
                      title="Non-aktifkan"
                    >
                      <svg
                        class="w-5 h-5"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"
                        />
                      </svg>
                    </button>
                    <button
                      v-if="schedule.delete_status"
                      @click="activateSchedule(schedule)"
                      class="text-green-600 hover:text-green-900"
                      title="Activate"
                    >
                      <svg
                        class="w-5 h-5"
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
                    </button>
                    <button
                      @click="showHardDeleteConfirm(schedule)"
                      class="text-red-600 hover:text-red-900"
                      title="Hapus Permanen"
                    >
                      <svg
                        class="w-5 h-5"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                        />
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Soft Delete Confirmation Modal -->
      <div
        v-if="showSoftDeleteModal"
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      >
        <div class="bg-white rounded-lg p-4 sm:p-6 w-full max-w-md mx-auto">
          <h3 class="text-lg sm:text-xl font-semibold mb-4 text-yellow-600">
            Non-aktifkan Line Schedule
          </h3>
          <p class="text-sm sm:text-base text-gray-700 mb-6">
            Apakah Anda yakin ingin menonaktifkan line schedule
            <strong>{{ scheduleToDelete?.id_registrasi }}</strong
            >? Line schedule yang dinonaktifkan masih dapat dilihat dalam
            sistem.
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
            Hapus Permanen Line Schedule
          </h3>
          <p class="text-sm sm:text-base text-gray-700 mb-6">
            Apakah Anda yakin ingin menghapus line schedule
            <strong>{{ scheduleToHardDelete?.id_registrasi }}</strong>
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

      <!-- Activation Confirmation Modal -->
      <div
        v-if="showActivateModal"
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      >
        <div class="bg-white rounded-lg p-4 sm:p-6 w-full max-w-md mx-auto">
          <h3 class="text-lg sm:text-xl font-semibold mb-4 text-green-600">
            Aktifkan Line Schedule
          </h3>
          <p class="text-sm sm:text-base text-gray-700 mb-6">
            Apakah Anda yakin ingin mengaktifkan kembali line schedule
            <strong>{{ scheduleToActivate?.id_registrasi }}</strong
            >?
          </p>

          <div class="flex flex-col sm:flex-row justify-end gap-3">
            <button
              @click="showActivateModal = false"
              class="px-4 py-2 text-sm sm:text-base text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors order-2 sm:order-1"
            >
              Batal
            </button>
            <button
              @click="confirmActivate"
              :disabled="submitting"
              class="px-4 py-2 text-sm sm:text-base bg-green-600 hover:bg-green-700 text-white rounded-lg transition-colors disabled:opacity-50 order-1 sm:order-2"
            >
              {{ submitting ? "Mengaktifkan..." : "Aktifkan" }}
            </button>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div
        class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6"
      >
        <div class="flex-1 flex justify-between sm:hidden">
          <button
            @click="prevPage"
            :disabled="currentPage === 1"
            class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
          >
            Previous
          </button>
          <button
            @click="nextPage"
            :disabled="currentPage === totalPages"
            class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
          >
            Next
          </button>
        </div>
        <div
          class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between"
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
            >
              <button
                @click="prevPage"
                :disabled="currentPage === 1"
                class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
              >
                <span class="sr-only">Previous</span>
                <svg
                  class="h-5 w-5"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 20 20"
                  fill="currentColor"
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
              >
                <span class="sr-only">Next</span>
                <svg
                  class="h-5 w-5"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 20 20"
                  fill="currentColor"
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
    </div>
  </AuthenticatedLayout>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useToast } from "vue-toastification";
import { useAuthStore } from "../stores/auth";
import AuthenticatedLayout from "../layouts/AuthenticatedLayout.vue";
import {
  lineScheduleService,
  buyerService,
  styleService,
} from "../services/api";

const authStore = useAuthStore();
const toast = useToast();

// Reactive state
const loading = ref(false);
const schedules = ref([]);
const buyers = ref([]);
const styles = ref([]);
const currentPage = ref(1);
const pageSize = ref(10);
const totalItems = ref(0);
const totalPages = ref(0);
const newRows = ref([]); // Change to array for multiple rows

// Available lines for assembly
const lines = ref(["LINE 1", "LINE 2", "LINE 3", "LINE 4", "LINE 5"]);

// Filters
const filters = ref({
  factory: "",
  date: new Date().toISOString().split("T")[0],
});

const canFetchData = computed(() => {
  return filters.value.date; // Only require date for fetching data
});

// Computed
const filteredStyles = computed(() => {
  if (!newRows.value[0]?.buyer_short_name) return [];
  return styles.value.filter(
    (style) => style.buyer_short_name === newRows.value[0].buyer_short_name
  );
});

// Computed property for pagination
const displayedPages = computed(() => {
  const delta = 2;
  const range = [];
  const rangeWithDots = [];
  let l;

  for (
    let i = Math.max(2, currentPage.value - delta);
    i <= Math.min(totalPages.value - 1, currentPage.value + delta);
    i++
  ) {
    range.push(i);
  }

  if (currentPage.value - delta > 2) {
    range.unshift("...");
  }
  if (currentPage.value + delta < totalPages.value - 1) {
    range.push("...");
  }

  range.unshift(1);
  if (totalPages.value > 1) {
    range.push(totalPages.value);
  }

  return range;
});

// In LineScheduleManagement.vue, add these helper functions:
const mapTypeToFactory = (type) => {
  switch (type) {
    case "component":
      return "F2"; // Use F2 for component area
    case "assembly":
    default:
      return "F1"; // Default to F1 for assembly
  }
};

const mapFactoryToType = (factory) => {
  switch (factory) {
    case "F2": // Assuming F2 is used for component area
      return "component";
    default:
      return "assembly";
  }
};

// Add generateRegistrationId function
const generateRegistrationId = (type, line, styleNo, startDate) => {
  const factory = mapTypeToFactory(type);
  const prefix = type === "component" ? "A0" : "L";
  const lineValue =
    type === "component" ? "" : line?.replace("LINE ", "") || "";
  const formattedDate = startDate ? startDate.replace(/-/g, "") : "";
  const formattedStyleNo = styleNo || "XXXX";

  return `${factory}-${prefix}${lineValue}-${formattedStyleNo}-${formattedDate}`;
};

// Modify the saveNewRow function:
const saveNewRow = async (row, index) => {
  row.id_registrasi = generateRegistrationId(
    row.type,
    row.line,
    row.style_no,
    row.start_date
  );

  if (!validateNewRow(row)) {
    return;
  }

  // Format the dates properly
  const startDate = row.start_date ? new Date(row.start_date) : null;
  const currentDate = new Date();

  const payload = {
    id_registrasi: row.id_registrasi,
    date: currentDate,
    factory: mapTypeToFactory(row.type),
    line: row.line,
    buyer_short_name: row.buyer_short_name,
    style_no: row.style_no,
    start_date: startDate,
    number_of_mp: parseInt(row.number_of_mp),
    working_day: parseInt(row.working_day),
  };

  try {
    const response = await lineScheduleService.createLineSchedule(payload);
    if (response.success) {
      toast.success("Schedule created successfully");
      newRows.value.splice(index, 1); // Remove only the saved row
      fetchSchedules();
    } else {
      toast.error(response.message || "Failed to create schedule");
    }
  } catch (error) {
    console.error("Error creating schedule:", error);
    toast.error("An error occurred while creating the schedule");
  }
};

const cancelNewRow = (index) => {
  newRows.value.splice(index, 1);
};

const validateNewRow = (row) => {
  const required = [
    "type",
    "line",
    "buyer_short_name",
    "style_no",
    "number_of_mp",
    "start_date",
    "working_day",
  ];
  const missing = required.filter((field) => !row[field]);

  if (missing.length > 0) {
    toast.error(`Please fill in all required fields: ${missing.join(", ")}`);
    return false;
  }
  return true;
};

const editSchedule = (schedule) => {
  newRow.value = { ...schedule };
};

const deactivateSchedule = (schedule) => {
  showSoftDeleteConfirm(schedule);
};

const activateSchedule = (schedule) => {
  showActivateConfirm(schedule);
};

const deleteSchedule = (schedule) => {
  showHardDeleteConfirm(schedule);
};

const exportToExcel = () => {
  const url = `${window.location.origin}/api/line-schedules/export?type=${filters.value.type}&date=${filters.value.date}`;
  window.open(url, "_blank");
};

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--;
    fetchSchedules();
  }
};

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
    fetchSchedules();
  }
};

const goToPage = (page) => {
  if (page !== "..." && page !== currentPage.value) {
    currentPage.value = page;
    fetchSchedules();
  }
};

const formatDate = (date) => {
  if (!date) return "";
  return new Date(date).toLocaleDateString("id-ID", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  });
};

// When displaying data, map factory back to type:
const displayedSchedules = computed(() => {
  return schedules.value.map((schedule) => ({
    ...schedule,
    type: mapFactoryToType(schedule.factory),
  }));
});

// Add the addNewRow function
const addNewRow = () => {
  newRows.value.push({
    id_registrasi: "", // Will be auto-generated before saving
    type: "assembly",
    line: "",
    buyer_short_name: "",
    style_no: "",
    number_of_mp: null,
    start_date: new Date().toISOString().split("T")[0],
    working_day: null,
  });
};

// Add new ref for tracking if filter has been applied
const isFilterApplied = ref(false);

// Add new refs for modals and schedule to delete
const showSoftDeleteModal = ref(false);
const showHardDeleteModal = ref(false);
const scheduleToDelete = ref(null);
const scheduleToHardDelete = ref(null);
const submitting = ref(false);

// Add new refs for activation modal
const showActivateModal = ref(false);
const scheduleToActivate = ref(null);

// Methods
const fetchSchedules = async () => {
  if (!canFetchData.value) {
    schedules.value = [];
    totalItems.value = 0;
    totalPages.value = 0;
    return;
  }

  loading.value = true;
  try {
    const response = await lineScheduleService.getAllLineSchedules(
      currentPage.value,
      pageSize.value,
      {
        date: filters.value.date,
        factory: filters.value.factory,
      }
    );
    if (response.success) {
      // Safely handle null or undefined line_schedules
      if (response.data?.line_schedules) {
        schedules.value = response.data.line_schedules.map((schedule) => ({
          ...schedule,
          type: mapFactoryToType(schedule.factory),
        }));
        totalItems.value = response.data.total_items;
        totalPages.value = response.data.total_pages;
        currentPage.value = response.data.current_page;
      } else {
        schedules.value = [];
        totalItems.value = 0;
        totalPages.value = 0;
      }
    } else {
      schedules.value = [];
      totalItems.value = 0;
      totalPages.value = 0;
      toast.error(response.message || "Failed to fetch schedules");
    }
  } catch (error) {
    console.error("Error fetching schedules:", error);
    schedules.value = [];
    totalItems.value = 0;
    totalPages.value = 0;
    toast.error("An error occurred while fetching schedules");
  } finally {
    loading.value = false;
  }
};

const fetchBuyers = async () => {
  try {
    const response = await buyerService.getAllBuyers(1, 100);
    if (response.success) {
      buyers.value = response.data.buyers;
    }
  } catch (error) {
    console.error("Error fetching buyers:", error);
    toast.error("Failed to fetch buyers");
  }
};

const fetchStyles = async () => {
  try {
    const response = await styleService.getAllStyles(1, 1000);
    if (response.success) {
      styles.value = response.data.styles;
    }
  } catch (error) {
    console.error("Error fetching styles:", error);
    toast.error("Failed to fetch styles");
  }
};

const handleFilterChange = () => {
  currentPage.value = 1; // Reset to first page when filter changes
};

const applyFilters = () => {
  isFilterApplied.value = true;
  fetchSchedules();
};

// Add new methods for soft delete
const showSoftDeleteConfirm = (schedule) => {
  scheduleToDelete.value = schedule;
  showSoftDeleteModal.value = true;
};

const confirmSoftDelete = async () => {
  submitting.value = true;
  try {
    const response = await lineScheduleService.deleteLineSchedule(
      scheduleToDelete.value.row_id
    );
    if (response.success) {
      await fetchSchedules();
      showSoftDeleteModal.value = false;
      toast.success("Line schedule berhasil dinonaktifkan");
    } else {
      toast.error(response.message || "Gagal menonaktifkan line schedule");
    }
  } catch (error) {
    console.error("Error deactivating schedule:", error);
    toast.error("Terjadi kesalahan saat menonaktifkan line schedule");
  } finally {
    submitting.value = false;
  }
};

// Add new methods for hard delete
const showHardDeleteConfirm = (schedule) => {
  scheduleToHardDelete.value = schedule;
  showHardDeleteModal.value = true;
};

const confirmHardDelete = async () => {
  submitting.value = true;
  try {
    const response = await lineScheduleService.hardDeleteLineSchedule(
      scheduleToHardDelete.value.row_id
    );
    if (response.success) {
      await fetchSchedules();
      showHardDeleteModal.value = false;
      toast.success("Line schedule berhasil dihapus secara permanen");
    } else {
      toast.error(response.message || "Gagal menghapus line schedule");
    }
  } catch (error) {
    console.error("Error hard deleting schedule:", error);
    toast.error("Terjadi kesalahan saat menghapus line schedule");
  } finally {
    submitting.value = false;
  }
};

// Add new methods for activation
const showActivateConfirm = (schedule) => {
  scheduleToActivate.value = schedule;
  showActivateModal.value = true;
};

const confirmActivate = async () => {
  submitting.value = true;
  try {
    const response = await lineScheduleService.activateLineSchedule(
      scheduleToActivate.value.row_id
    );
    if (response.success) {
      await fetchSchedules();
      showActivateModal.value = false;
      toast.success("Line schedule berhasil diaktifkan");
    } else {
      toast.error(response.message || "Gagal mengaktifkan line schedule");
    }
  } catch (error) {
    console.error("Error activating schedule:", error);
    toast.error("Terjadi kesalahan saat mengaktifkan line schedule");
  } finally {
    submitting.value = false;
  }
};

// Lifecycle hooks
onMounted(() => {
  filters.value.date = new Date().toISOString().split("T")[0];
  fetchBuyers();
  fetchStyles();
  // Remove initial fetchSchedules call
});

const handleBuyerChange = (row) => {
  // Reset style when buyer changes
  row.style_no = "";
};

const handleNewRowTypeChange = (row) => {
  if (row.type === "component") {
    row.line = "AREA";
  } else {
    row.line = "";
  }
};
</script>

<style scoped>
.btn {
  @apply px-4 py-2 rounded-lg text-sm font-medium transition-colors duration-200;
}

.btn-primary {
  @apply bg-blue-600 text-white hover:bg-blue-700;
}

.btn-secondary {
  @apply bg-gray-100 text-gray-700 hover:bg-gray-200;
}

.form-select {
  @apply mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-lg;
}

.form-input {
  @apply mt-1 block w-full px-3 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-lg;
}
</style>
