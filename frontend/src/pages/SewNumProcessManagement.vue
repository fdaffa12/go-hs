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
              Manajemen Sew Numbering Process
            </h1>
            <p class="mt-1 text-sm text-gray-600">
              Kelola proses penomoran jahitan
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
              <span class="hidden sm:inline">Tambah Proses</span>
              <span class="sm:hidden">Tambah</span>
            </button>

            <!-- Add Save All button -->
            <button
              v-if="newRows.length > 0"
              @click="showSaveAllConfirm"
              class="btn btn-success flex items-center justify-center gap-2 w-full sm:w-auto"
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
                  d="M5 13l4 4L19 7"
                ></path>
              </svg>
              <span>Save All ({{ newRows.length }})</span>
            </button>

            <!-- Add bulk action buttons -->
            <div v-if="hasSelection" class="flex gap-2 w-full sm:w-auto">
              <button
                @click="showBulkSoftDeleteConfirm"
                class="btn btn-warning flex items-center justify-center gap-2 flex-1 sm:flex-none"
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
                  ></path>
                </svg>
                <span>Non-aktifkan ({{ selectedCount }})</span>
              </button>
              <button
                @click="showBulkActivateConfirm"
                class="btn btn-success flex items-center justify-center gap-2 flex-1 sm:flex-none"
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
                  ></path>
                </svg>
                <span>Aktifkan ({{ selectedCount }})</span>
              </button>
              <button
                @click="showBulkHardDeleteConfirm"
                class="btn btn-danger flex items-center justify-center gap-2 flex-1 sm:flex-none"
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
                  ></path>
                </svg>
                <span>Hapus ({{ selectedCount }})</span>
              </button>
            </div>

            <!-- Import Button -->
            <button
              @click="$refs.importFile.click()"
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
                  d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
                ></path>
              </svg>
              <span>Import</span>
            </button>
            <input
              ref="importFile"
              type="file"
              accept=".xlsx,.xls"
              class="hidden"
              @change="handleImport"
            />

            <!-- Export Button -->
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
          <div class="form-group">
            <label class="form-label">Buyer</label>
            <select
              v-model="filters.buyer"
              class="form-select h-[38px]"
              @change="handleBuyerChange"
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
          </div>
          <div class="form-group">
            <label class="form-label">Style</label>
            <select
              v-model="filters.style"
              class="form-select h-[38px]"
              :disabled="!filters.buyer"
              @change="handleStyleChange"
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
          </div>
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
                <th class="w-4 px-6 py-3">
                  <input
                    type="checkbox"
                    :checked="selectAll"
                    @change="toggleSelectAll"
                    class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  No Process
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Category
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Sub Category
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  SMV GSD
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  SMV EST
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Process Name (ENG)
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Process Name (IND)
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Machine Code
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
                <td class="w-4 px-6 py-4">
                  <span class="text-blue-500 text-xs">New</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <input
                    v-model.number="row.no_process"
                    type="number"
                    class="form-input"
                    placeholder="No Process"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <select v-model="row.category" class="form-select">
                    <option value="">Select Category</option>
                    <option value="SUPPORT PART">SUPPORT PART</option>
                    <option value="ASSEMBLY PART">ASSEMBLY PART</option>
                  </select>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <input
                    v-model="row.sub_category"
                    class="form-input"
                    placeholder="Sub Category"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <input
                    v-model.number="row.smv_proc_gsd"
                    type="number"
                    step="0.0001"
                    class="form-input"
                    placeholder="SMV GSD"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <input
                    v-model.number="row.smv_proc_est"
                    type="number"
                    step="0.0001"
                    class="form-input"
                    placeholder="SMV EST"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <input
                    v-model="row.process_name_eng"
                    class="form-input"
                    placeholder="Process Name (ENG)"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <input
                    v-model="row.process_name_ind"
                    class="form-input"
                    placeholder="Process Name (IND)"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <input
                    v-model="row.machine_code"
                    class="form-input"
                    placeholder="Machine Code"
                    maxlength="2"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    class="px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-800"
                  >
                    New
                  </span>
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
              <tr v-if="!filters.buyer">
                <td colspan="11" class="px-6 py-4 text-center text-gray-500">
                  Silakan pilih buyer terlebih dahulu
                </td>
              </tr>
              <tr v-else-if="!filters.style">
                <td colspan="11" class="px-6 py-4 text-center text-gray-500">
                  Silakan pilih style terlebih dahulu
                </td>
              </tr>
              <tr v-else-if="loading">
                <td colspan="11" class="px-6 py-4 text-center text-gray-500">
                  Loading...
                </td>
              </tr>
              <tr v-else-if="!processes || processes.length === 0">
                <td colspan="11" class="px-6 py-4 text-center text-gray-500">
                  Tidak ada data proses untuk style ini
                </td>
              </tr>
              <tr
                v-for="process in processes"
                :key="process.id"
                class="hover:bg-gray-50"
              >
                <td class="w-4 px-6 py-4">
                  <input
                    type="checkbox"
                    :checked="selectedRows.has(process.id)"
                    @change="toggleRowSelection(process.id)"
                    class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div v-if="process.isEditing">
                    <input
                      v-model.number="process.editedNoProcess"
                      type="number"
                      class="form-input"
                      placeholder="No Process"
                    />
                  </div>
                  <div v-else class="text-sm text-gray-900">
                    {{ process.no_process }}
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div v-if="process.isEditing">
                    <select
                      v-model="process.editedCategory"
                      class="form-select"
                    >
                      <option value="">Select Category</option>
                      <option value="SUPPORT PART">SUPPORT PART</option>
                      <option value="ASSEMBLY PART">ASSEMBLY PART</option>
                    </select>
                  </div>
                  <div v-else class="text-sm text-gray-900">
                    {{ process.category }}
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div v-if="process.isEditing">
                    <input
                      v-model="process.editedSubCategory"
                      class="form-input"
                      placeholder="Sub Category"
                    />
                  </div>
                  <div v-else class="text-sm text-gray-900">
                    {{ process.sub_category }}
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div v-if="process.isEditing">
                    <input
                      v-model.number="process.editedSmvProcGsd"
                      type="number"
                      step="0.0001"
                      class="form-input"
                      placeholder="SMV GSD"
                    />
                  </div>
                  <div v-else class="text-sm text-gray-900">
                    {{ process.smv_proc_gsd }}
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div v-if="process.isEditing">
                    <input
                      v-model.number="process.editedSmvProcEst"
                      type="number"
                      step="0.0001"
                      class="form-input"
                      placeholder="SMV EST"
                    />
                  </div>
                  <div v-else class="text-sm text-gray-900">
                    {{ process.smv_proc_est }}
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div v-if="process.isEditing">
                    <input
                      v-model="process.editedProcessNameEng"
                      class="form-input"
                      placeholder="Process Name (ENG)"
                    />
                  </div>
                  <div v-else class="text-sm text-gray-900">
                    {{ process.process_name_eng }}
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div v-if="process.isEditing">
                    <input
                      v-model="process.editedProcessNameInd"
                      class="form-input"
                      placeholder="Process Name (IND)"
                    />
                  </div>
                  <div v-else class="text-sm text-gray-900">
                    {{ process.process_name_ind }}
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div v-if="process.isEditing">
                    <input
                      v-model="process.editedMachineCode"
                      class="form-input"
                      placeholder="Machine Code"
                      maxlength="2"
                    />
                  </div>
                  <div v-else class="text-sm text-gray-900">
                    {{ process.machine_code }}
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    :class="[
                      'px-2 py-1 text-xs rounded-full',
                      process.delete_status
                        ? 'bg-red-100 text-red-800'
                        : 'bg-green-100 text-green-800',
                    ]"
                  >
                    {{ process.delete_status ? "Non-Aktif" : "Aktif" }}
                  </span>
                </td>
                <td
                  class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                >
                  <div class="flex justify-end space-x-2">
                    <template v-if="process.isEditing">
                      <button
                        @click="saveInlineEdit(process)"
                        class="text-green-600 hover:text-green-900"
                        title="Simpan"
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
                        @click="cancelInlineEdit(process)"
                        class="text-red-600 hover:text-red-900"
                        title="Batal"
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
                    </template>
                    <template v-else>
                      <button
                        v-if="!process.delete_status"
                        @click="startInlineEdit(process)"
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
                      <!-- Action buttons in table -->
                      <button
                        v-if="!process.delete_status"
                        @click="showSoftDeleteConfirm(process)"
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
                      <!-- Activate button in the table -->
                      <button
                        v-if="process.delete_status"
                        @click="showActivateConfirm(process)"
                        class="text-green-600 hover:text-green-900"
                        title="Aktifkan"
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
                        @click="showHardDeleteConfirm(process)"
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
                    </template>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Modals -->
    <div
      v-if="showSoftDeleteModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
      <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
        <h3 class="text-lg font-semibold mb-4">Konfirmasi Non-aktifkan</h3>
        <p class="text-gray-600 mb-6">
          Apakah Anda yakin ingin menonaktifkan proses ini?
        </p>
        <div class="flex justify-end gap-2">
          <button
            @click="showSoftDeleteModal = false"
            class="btn btn-secondary"
          >
            Batal
          </button>
          <button @click="confirmSoftDelete" class="btn btn-warning">
            Non-aktifkan
          </button>
        </div>
      </div>
    </div>

    <div
      v-if="showBulkSoftDeleteModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
      <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
        <h3 class="text-lg font-semibold mb-4">
          Konfirmasi Non-aktifkan Massal
        </h3>
        <p class="text-gray-600 mb-6">
          Apakah Anda yakin ingin menonaktifkan {{ selectedCount }} proses yang
          dipilih?
        </p>
        <div class="flex justify-end gap-2">
          <button
            @click="showBulkSoftDeleteModal = false"
            class="btn btn-secondary"
          >
            Batal
          </button>
          <button @click="confirmBulkSoftDelete" class="btn btn-warning">
            Non-aktifkan
          </button>
        </div>
      </div>
    </div>

    <!-- Activate Modal -->
    <div
      v-if="showActivateModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
      <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
        <h3 class="text-lg font-semibold mb-4">Konfirmasi Aktivasi</h3>
        <p class="text-gray-600 mb-6">
          Apakah Anda yakin ingin mengaktifkan kembali proses ini?
        </p>
        <div class="flex justify-end gap-2">
          <button @click="showActivateModal = false" class="btn btn-secondary">
            Batal
          </button>
          <button @click="confirmActivate" class="btn btn-success">
            Aktifkan
          </button>
        </div>
      </div>
    </div>

    <div
      v-if="showBulkActivateModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
      <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
        <h3 class="text-lg font-semibold mb-4">Konfirmasi Aktivasi Massal</h3>
        <p class="text-gray-600 mb-6">
          Apakah Anda yakin ingin mengaktifkan {{ selectedCount }} proses yang
          dipilih?
        </p>
        <div class="flex justify-end gap-2">
          <button
            @click="showBulkActivateModal = false"
            class="btn btn-secondary"
          >
            Batal
          </button>
          <button @click="confirmBulkActivate" class="btn btn-success">
            Aktifkan
          </button>
        </div>
      </div>
    </div>

    <div
      v-if="showHardDeleteModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
      <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
        <h3 class="text-lg font-semibold mb-4">Konfirmasi Hapus Permanen</h3>
        <p class="text-gray-600 mb-6">
          Apakah Anda yakin ingin menghapus permanen proses ini? Tindakan ini
          tidak dapat dibatalkan.
        </p>
        <div class="flex justify-end gap-2">
          <button
            @click="showHardDeleteModal = false"
            class="btn btn-secondary"
          >
            Batal
          </button>
          <button @click="confirmHardDelete" class="btn btn-danger">
            Hapus Permanen
          </button>
        </div>
      </div>
    </div>

    <div
      v-if="showBulkHardDeleteModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
      <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
        <h3 class="text-lg font-semibold mb-4">
          Konfirmasi Hapus Permanen Massal
        </h3>
        <p class="text-gray-600 mb-6">
          Apakah Anda yakin ingin menghapus permanen {{ selectedCount }} proses
          yang dipilih? Tindakan ini tidak dapat dibatalkan.
        </p>
        <div class="flex justify-end gap-2">
          <button
            @click="showBulkHardDeleteModal = false"
            class="btn btn-secondary"
          >
            Batal
          </button>
          <button @click="confirmBulkHardDelete" class="btn btn-danger">
            Hapus Permanen
          </button>
        </div>
      </div>
    </div>

    <div
      v-if="showSaveAllModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
      <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
        <h3 class="text-lg font-semibold mb-4">Konfirmasi Simpan Semua</h3>
        <p class="text-gray-600 mb-6">
          Apakah Anda yakin ingin menyimpan {{ newRows.length }} proses baru?
        </p>
        <div class="flex justify-end gap-2">
          <button @click="showSaveAllModal = false" class="btn btn-secondary">
            Batal
          </button>
          <button @click="confirmSaveAll" class="btn btn-success">
            Simpan Semua
          </button>
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
  sewNumProcessService,
  buyerService,
  styleService,
} from "../services/api";
import * as XLSX from "xlsx";

const authStore = useAuthStore();
const toast = useToast();

// Reactive state
const loading = ref(false);
const processes = ref([]);
const buyers = ref([]);
const styles = ref([]);
const newRows = ref([]);
const selectedRows = ref(new Set());
const selectAll = ref(false);

// Modal states
const showSoftDeleteModal = ref(false);
const showBulkSoftDeleteModal = ref(false);
const showActivateModal = ref(false);
const showBulkActivateModal = ref(false);
const showHardDeleteModal = ref(false);
const showBulkHardDeleteModal = ref(false);
const showSaveAllModal = ref(false);

// Store process to be deleted/activated
const selectedProcess = ref(null);

// Filters
const filters = ref({
  buyer: "",
  style: "",
});

// Computed properties
const hasSelection = computed(() => selectedRows.value.size > 0);
const selectedCount = computed(() => selectedRows.value.size);

const filteredStyles = computed(() => {
  if (!filters.value.buyer) return [];
  return styles.value.filter(
    (style) => style.buyer_short_name === filters.value.buyer
  );
});

// Methods
const handleBuyerChange = async () => {
  filters.value.style = "";
  processes.value = [];

  if (!filters.value.buyer) return;

  try {
    const response = await styleService.getAllStyles(
      1,
      1000,
      filters.value.buyer
    );
    if (response.success) {
      styles.value = response.data.styles;
    }
  } catch (error) {
    toast.error("Failed to fetch styles");
  }
};

const handleStyleChange = async () => {
  if (!filters.value.style) {
    processes.value = [];
    return;
  }

  loading.value = true;
  try {
    const response = await sewNumProcessService.getAllProcesses(
      filters.value.style
    );
    if (response.success) {
      processes.value = response.data;
    } else {
      toast.error(response.message || "Failed to fetch processes");
    }
  } catch (error) {
    toast.error("An error occurred while fetching processes");
  } finally {
    loading.value = false;
  }
};

// Lifecycle hooks
onMounted(async () => {
  try {
    const response = await buyerService.getAllBuyers(1, 100);
    if (response.success) {
      buyers.value = response.data.buyers;
    }
  } catch (error) {
    toast.error("Failed to fetch buyers");
  }
});

// Add these methods to the script section
const addNewRow = () => {
  if (!filters.value.buyer || !filters.value.style) {
    toast.error("Please select buyer and style first");
    return;
  }

  newRows.value.push({
    buyer_short_name: filters.value.buyer,
    style_no: filters.value.style,
    no_process: null,
    category: "",
    sub_category: "",
    smv_proc_gsd: null,
    smv_proc_est: 0,
    process_name_eng: null,
    process_name_ind: null,
    machine_code: "",
  });
};

const saveNewRow = async (row, index) => {
  if (!validateNewRow(row)) {
    return;
  }

  try {
    const response = await sewNumProcessService.bulkSave([row]);
    if (response.success) {
      toast.success("Process saved successfully");
      newRows.value.splice(index, 1);
      handleStyleChange();
    } else {
      toast.error(response.message || "Failed to save process");
    }
  } catch (error) {
    toast.error("An error occurred while saving process");
  }
};

const cancelNewRow = (index) => {
  newRows.value.splice(index, 1);
};

const validateNewRow = (row) => {
  if (
    !row.no_process ||
    !row.category ||
    !row.sub_category ||
    !row.machine_code
  ) {
    toast.error("Please fill in all required fields");
    return false;
  }

  if (row.smv_proc_gsd && isNaN(row.smv_proc_gsd)) {
    toast.error("SMV GSD must be a number");
    return false;
  }

  if (isNaN(row.smv_proc_est)) {
    toast.error("SMV EST must be a number");
    return false;
  }

  return true;
};

const startInlineEdit = (process) => {
  process.isEditing = true;
  process.editedNoProcess = process.no_process;
  process.editedCategory = process.category;
  process.editedSubCategory = process.sub_category;
  process.editedSmvProcGsd = process.smv_proc_gsd;
  process.editedSmvProcEst = process.smv_proc_est;
  process.editedProcessNameEng = process.process_name_eng;
  process.editedProcessNameInd = process.process_name_ind;
  process.editedMachineCode = process.machine_code;
};

const saveInlineEdit = async (process) => {
  const updatedProcess = {
    id: process.id,
    buyer_short_name: process.buyer_short_name,
    style_no: process.style_no,
    no_process: process.editedNoProcess,
    category: process.editedCategory,
    sub_category: process.editedSubCategory,
    smv_proc_gsd: process.editedSmvProcGsd,
    smv_proc_est: process.editedSmvProcEst,
    process_name_eng: process.editedProcessNameEng,
    process_name_ind: process.editedProcessNameInd,
    machine_code: process.editedMachineCode,
  };

  if (!validateNewRow(updatedProcess)) {
    return;
  }

  try {
    const response = await sewNumProcessService.bulkSave([updatedProcess]);
    if (response.success) {
      toast.success("Process updated successfully");
      process.isEditing = false;
      handleStyleChange();
    } else {
      toast.error(response.message || "Failed to update process");
    }
  } catch (error) {
    toast.error("An error occurred while updating process");
  }
};

const cancelInlineEdit = (process) => {
  process.isEditing = false;
  delete process.editedNoProcess;
  delete process.editedCategory;
  delete process.editedSubCategory;
  delete process.editedSmvProcGsd;
  delete process.editedSmvProcEst;
  delete process.editedProcessNameEng;
  delete process.editedProcessNameInd;
  delete process.editedMachineCode;
};

const showSoftDeleteConfirm = (process) => {
  selectedProcess.value = process;
  showSoftDeleteModal.value = true;
};

const confirmSoftDelete = async () => {
  try {
    const response = await sewNumProcessService.deleteProcess(
      selectedProcess.value.id
    );
    if (response.success) {
      toast.success("Proses berhasil dinonaktifkan");
      handleStyleChange();
    } else {
      toast.error(response.message || "Gagal menonaktifkan proses");
    }
  } catch (error) {
    toast.error("Terjadi kesalahan saat menonaktifkan proses");
  } finally {
    showSoftDeleteModal.value = false;
    selectedProcess.value = null;
  }
};

const showActivateConfirm = (process) => {
  selectedProcess.value = process;
  showActivateModal.value = true;
};

const confirmActivate = async () => {
  try {
    const response = await sewNumProcessService.activateProcess(
      selectedProcess.value.id
    );
    if (response.success) {
      toast.success("Proses berhasil diaktifkan");
      handleStyleChange();
    } else {
      toast.error(response.message || "Gagal mengaktifkan proses");
    }
  } catch (error) {
    toast.error("Terjadi kesalahan saat mengaktifkan proses");
  } finally {
    showActivateModal.value = false;
    selectedProcess.value = null;
  }
};

const showHardDeleteConfirm = (process) => {
  selectedProcess.value = process;
  showHardDeleteModal.value = true;
};

const confirmHardDelete = async () => {
  try {
    const response = await sewNumProcessService.hardDeleteProcess(
      selectedProcess.value.id
    );
    if (response.success) {
      toast.success("Proses berhasil dihapus permanen");
      handleStyleChange();
    } else {
      toast.error(response.message || "Gagal menghapus proses");
    }
  } catch (error) {
    toast.error("Terjadi kesalahan saat menghapus proses");
  } finally {
    showHardDeleteModal.value = false;
    selectedProcess.value = null;
  }
};

const toggleSelectAll = () => {
  selectAll.value = !selectAll.value;
  if (selectAll.value) {
    processes.value.forEach((process) => {
      selectedRows.value.add(process.id);
    });
  } else {
    selectedRows.value.clear();
  }
};

const toggleRowSelection = (id) => {
  if (selectedRows.value.has(id)) {
    selectedRows.value.delete(id);
    selectAll.value = false;
  } else {
    selectedRows.value.add(id);
    selectAll.value = processes.value.every((process) =>
      selectedRows.value.has(process.id)
    );
  }
};

const exportToExcel = () => {
  if (!processes.value.length) {
    toast.error("No data to export");
    return;
  }

  // Create workbook and worksheet
  const wb = XLSX.utils.book_new();
  const ws_data = [];

  // Add header row
  const header = [
    "NO_PROCESS",
    "CATEGORY",
    "SUB_CATEGORY",
    "SMV_PROC_GSD",
    "SMV_PROC_EST",
    "PROCESS_NAME_ENG",
    "PROCESS_NAME_IND",
    "MACHINE_CODE",
  ];
  ws_data.push(header);

  // Add data rows
  processes.value.forEach((process) => {
    const row = [
      process.no_process,
      process.category,
      process.sub_category,
      process.smv_proc_gsd,
      process.smv_proc_est,
      process.process_name_eng,
      process.process_name_ind,
      process.machine_code,
    ];
    ws_data.push(row);
  });

  // Create worksheet and append to workbook
  const ws = XLSX.utils.aoa_to_sheet(ws_data);
  XLSX.utils.book_append_sheet(wb, ws, "Sew Num Process");

  // Save file
  XLSX.writeFile(wb, `SewNumProcess_${filters.value.style}.xlsx`);
  toast.success("Data exported to Excel");
};

const handleImport = async (event) => {
  const file = event.target.files[0];
  if (!file) return;

  if (!filters.value.buyer || !filters.value.style) {
    toast.error("Please select buyer and style first");
    event.target.value = "";
    return;
  }

  try {
    const reader = new FileReader();
    reader.onload = async (e) => {
      const data = new Uint8Array(e.target.result);
      const workbook = XLSX.read(data, { type: "array" });
      const firstSheet = workbook.Sheets[workbook.SheetNames[0]];
      const jsonData = XLSX.utils.sheet_to_json(firstSheet, { header: 1 });

      const importedData = [];
      const invalidRows = [];

      // Skip header row and process each data row
      for (let i = 1; i < jsonData.length; i++) {
        const row = jsonData[i];
        const importData = {
          buyer_short_name: filters.value.buyer,
          style_no: filters.value.style,
          no_process: row[0],
          category: row[1],
          sub_category: row[2],
          smv_proc_gsd: row[3],
          smv_proc_est: row[4],
          process_name_eng: row[5],
          process_name_ind: row[6],
          machine_code: row[7],
        };

        // Basic validation
        if (
          !importData.no_process ||
          !importData.category ||
          !importData.sub_category ||
          !importData.machine_code
        ) {
          invalidRows.push({
            row: i + 1,
            reason: "Required fields missing",
          });
          continue;
        }

        importedData.push(importData);
      }

      // Show validation errors if any
      if (invalidRows.length > 0) {
        const errorMessages = invalidRows
          .map((error) => `Row ${error.row}: ${error.reason}`)
          .join("\n");
        toast.error(`Some data is invalid:\n${errorMessages}`);
        return;
      }

      // Save valid data
      try {
        const response = await sewNumProcessService.bulkSave(importedData);
        if (response.success) {
          toast.success(
            `${importedData.length} processes imported successfully`
          );
          handleStyleChange();
        } else {
          toast.error(response.message || "Failed to import processes");
        }
      } catch (error) {
        toast.error("An error occurred while importing processes");
      }
    };

    reader.readAsArrayBuffer(file);
  } catch (error) {
    toast.error("Error importing data: " + error.message);
  } finally {
    event.target.value = "";
  }
};

const saveAllNewRows = async () => {
  if (newRows.value.length === 0) {
    toast.error("No new rows to save.");
    return;
  }

  const validRows = newRows.value.filter(validateNewRow);
  if (validRows.length === 0) {
    toast.error("No valid new rows to save.");
    return;
  }

  try {
    const response = await sewNumProcessService.bulkSave(validRows);
    if (response.success) {
      toast.success(`${validRows.length} new processes saved successfully`);
      newRows.value = [];
      handleStyleChange();
    } else {
      toast.error(response.message || "Failed to save new processes");
    }
  } catch (error) {
    toast.error("An error occurred while saving new processes");
  }
};

const showBulkSoftDeleteConfirm = () => {
  if (selectedRows.value.size === 0) {
    toast.error("Please select at least one process to deactivate.");
    return;
  }
  showBulkSoftDeleteModal.value = true;
};

const confirmBulkSoftDelete = async () => {
  try {
    const response = await sewNumProcessService.bulkSoftDelete(
      Array.from(selectedRows.value)
    );
    if (response.success) {
      toast.success(`${selectedRows.value.size} proses berhasil dinonaktifkan`);
      selectedRows.value.clear();
      selectAll.value = false;
      handleStyleChange();
    } else {
      toast.error(response.message || "Gagal menonaktifkan proses");
    }
  } catch (error) {
    toast.error("Terjadi kesalahan saat menonaktifkan proses");
  } finally {
    showBulkSoftDeleteModal.value = false;
  }
};

const showBulkActivateConfirm = () => {
  if (selectedRows.value.size === 0) {
    toast.error("Please select at least one process to activate.");
    return;
  }
  showBulkActivateModal.value = true;
};

const confirmBulkActivate = async () => {
  try {
    const response = await sewNumProcessService.bulkActivate(
      Array.from(selectedRows.value)
    );
    if (response.success) {
      toast.success(`${selectedRows.value.size} proses berhasil diaktifkan`);
      selectedRows.value.clear();
      selectAll.value = false;
      handleStyleChange();
    } else {
      toast.error(response.message || "Gagal mengaktifkan proses");
    }
  } catch (error) {
    toast.error("Terjadi kesalahan saat mengaktifkan proses");
  } finally {
    showBulkActivateModal.value = false;
  }
};

const showBulkHardDeleteConfirm = () => {
  if (selectedRows.value.size === 0) {
    toast.error("Please select at least one process to delete permanently.");
    return;
  }
  showBulkHardDeleteModal.value = true;
};

const confirmBulkHardDelete = async () => {
  try {
    const response = await sewNumProcessService.bulkHardDelete(
      Array.from(selectedRows.value)
    );
    if (response.success) {
      toast.success(
        `${selectedRows.value.size} proses berhasil dihapus permanen`
      );
      selectedRows.value.clear();
      selectAll.value = false;
      handleStyleChange();
    } else {
      toast.error(response.message || "Gagal menghapus proses");
    }
  } catch (error) {
    toast.error("Terjadi kesalahan saat menghapus proses");
  } finally {
    showBulkHardDeleteModal.value = false;
  }
};

const showSaveAllConfirm = () => {
  if (newRows.value.length === 0) {
    toast.error("Tidak ada proses baru untuk disimpan.");
    return;
  }
  showSaveAllModal.value = true;
};

const confirmSaveAll = async () => {
  const validRows = newRows.value.filter(validateNewRow);
  if (validRows.length === 0) {
    toast.error("Tidak ada proses baru yang valid untuk disimpan.");
    showSaveAllModal.value = false;
    return;
  }

  try {
    const response = await sewNumProcessService.bulkSave(validRows);
    if (response.success) {
      toast.success(`${validRows.length} proses baru berhasil disimpan`);
      newRows.value = [];
      handleStyleChange();
    } else {
      toast.error(response.message || "Gagal menyimpan proses baru");
    }
  } catch (error) {
    toast.error("Terjadi kesalahan saat menyimpan proses baru");
  } finally {
    showSaveAllModal.value = false;
  }
};
</script>

<style scoped>
/* Add any component-specific styles here */
</style>
