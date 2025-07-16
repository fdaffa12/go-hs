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
              Manajemen Karyawan
            </h1>
            <p class="mt-1 text-sm text-gray-600">
              Kelola data karyawan dalam sistem
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
            <span class="hidden sm:inline">Tambah Karyawan</span>
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
                placeholder="Cari karyawan..."
                class="w-full pl-9 sm:pl-10 pr-4 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
          <div class="flex gap-2">
            <select
              v-model="filterDepartment"
              class="w-full sm:w-48 px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="">Semua Departemen</option>
              <option
                v-for="dept in departments"
                :key="dept.short_name"
                :value="dept.short_name"
              >
                {{ dept.long_name }}
              </option>
            </select>
            <button
              @click="fetchEmployees"
              class="bg-gray-100 hover:bg-gray-200 text-gray-700 px-3 sm:px-4 py-2 rounded-lg flex items-center justify-center gap-2 transition-colors"
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

      <!-- Employees Table -->
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
                  Departemen
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Jabatan
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Status
                </th>
                <th
                  class="px-3 sm:px-6 py-2 sm:py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Aksi
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-if="loading" class="animate-pulse">
                <td
                  colspan="6"
                  class="px-3 sm:px-6 py-3 sm:py-4 text-center text-gray-500 text-sm"
                >
                  Memuat data karyawan...
                </td>
              </tr>
              <tr v-else-if="filteredEmployees.length === 0">
                <td
                  colspan="6"
                  class="px-3 sm:px-6 py-3 sm:py-4 text-center text-gray-500 text-sm"
                >
                  Tidak ada karyawan ditemukan
                </td>
              </tr>
              <tr
                v-else
                v-for="employee in filteredEmployees"
                :key="employee.NIK"
                class="hover:bg-gray-50"
              >
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-sm font-medium text-gray-900">
                    {{ employee.NIK }}
                  </div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-sm text-gray-900">{{ employee.NAME }}</div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-sm text-gray-900">
                    {{ getDepartmentName(employee.DEPT_SHORT_NAME) }}
                  </div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <div class="text-sm text-gray-900">{{ employee.TITLE }}</div>
                </td>
                <td class="px-3 sm:px-6 py-3 sm:py-4">
                  <span
                    :class="[
                      'px-2 py-1 text-xs rounded-full',
                      employee.DELETE_STATUS
                        ? 'bg-red-100 text-red-800'
                        : 'bg-green-100 text-green-800',
                    ]"
                  >
                    {{ employee.DELETE_STATUS ? "Non-Aktif" : "Aktif" }}
                  </span>
                </td>
                <td
                  class="px-3 sm:px-6 py-3 sm:py-4 text-right text-sm font-medium"
                >
                  <div class="flex justify-end gap-1 sm:gap-2">
                    <button
                      v-if="!employee.DELETE_STATUS"
                      @click="editEmployee(employee)"
                      class="text-blue-600 hover:text-blue-900 p-1 sm:p-1.5 rounded transition-colors"
                      title="Edit Karyawan"
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
                      v-if="!employee.DELETE_STATUS"
                      @click="softDeleteEmployee(employee)"
                      class="text-yellow-600 hover:text-yellow-900 p-1 sm:p-1.5 rounded transition-colors"
                      title="Non-aktifkan Karyawan"
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
                      v-if="employee.DELETE_STATUS"
                      @click="activateEmployee(employee)"
                      class="text-green-600 hover:text-green-900 p-1 sm:p-1.5 rounded transition-colors"
                      title="Aktifkan Karyawan"
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
                      @click="showHardDeleteConfirm(employee)"
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
    </div>

    <!-- After the table -->
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
      <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
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
  </AuthenticatedLayout>

  <!-- Create/Edit Employee Modal -->
  <div
    v-if="showCreateModal || showEditModal"
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
  >
    <div class="bg-white rounded-lg p-4 sm:p-6 w-full max-w-xl mx-auto">
      <h3 class="text-lg sm:text-xl font-semibold mb-4">
        {{ showCreateModal ? "Tambah Karyawan Baru" : "Edit Data Karyawan" }}
      </h3>

      <form
        @submit.prevent="showCreateModal ? createEmployee() : updateEmployee()"
      >
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >NIK</label
            >
            <input
              v-model="employeeForm.NIK"
              type="text"
              required
              maxlength="10"
              :disabled="!showCreateModal"
              class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
              placeholder="HS12345678"
            />
            <p class="mt-1 text-xs text-gray-500">
              Format: HS diikuti 1-8 angka (contoh: HS12345678)
            </p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Nama</label
            >
            <input
              v-model="employeeForm.NAME"
              type="text"
              required
              maxlength="22"
              class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Masukkan nama"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Departemen</label
            >
            <select
              v-model="employeeForm.DEPT_SHORT_NAME"
              required
              class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="">Pilih Departemen</option>
              <option
                v-for="dept in departments"
                :key="dept.short_name"
                :value="dept.short_name"
              >
                {{ dept.long_name }}
              </option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Tanggal Masuk</label
            >
            <input
              v-model="employeeForm.ENTERANCE_DATE"
              type="date"
              required
              class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Jabatan</label
            >
            <input
              v-model="employeeForm.TITLE"
              type="text"
              required
              maxlength="22"
              class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Masukkan jabatan"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Jenis Kelamin</label
            >
            <select
              v-model="employeeForm.GENDER"
              required
              class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="">Pilih Jenis Kelamin</option>
              <option :value="1">Laki-laki</option>
              <option :value="0">Perempuan</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >RFID ID</label
            >
            <input
              v-model="employeeForm.RFID_ID"
              type="text"
              maxlength="15"
              class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Masukkan RFID ID (opsional)"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >Area Kerja</label
            >
            <select
              v-model="employeeForm.WORKING_AREA"
              required
              class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="">Pilih Area Kerja</option>
              <option :value="1">ALL AREA</option>
              <option :value="2">F1</option>
              <option :value="3">F2</option>
              <option :value="4">F3</option>
            </select>
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
                ? "Tambah Karyawan"
                : "Update Karyawan"
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
        Hapus Karyawan
      </h3>
      <p class="text-sm sm:text-base text-gray-700 mb-6">
        Apakah Anda yakin ingin menghapus karyawan
        <strong>{{ employeeToDelete?.NAME }}</strong
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
          {{ submitting ? "Menghapus..." : "Hapus Karyawan" }}
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
        Non-aktifkan Karyawan
      </h3>
      <p class="text-sm sm:text-base text-gray-700 mb-6">
        Apakah Anda yakin ingin menonaktifkan karyawan
        <strong>{{ employeeToDelete?.NAME }}</strong
        >? Karyawan yang dinonaktifkan masih dapat dilihat dalam sistem.
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
        Hapus Permanen Karyawan
      </h3>
      <p class="text-sm sm:text-base text-gray-700 mb-6">
        Apakah Anda yakin ingin menghapus karyawan
        <strong>{{ employeeToHardDelete?.NAME }}</strong>
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
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { useToast } from "vue-toastification";
import { useAuthStore } from "../stores/auth";
import { employeeService, departmentService } from "../services/api";
import AuthenticatedLayout from "../layouts/AuthenticatedLayout.vue";

const authStore = useAuthStore();
const toast = useToast();

// Reactive data
const employees = ref([]);
const departments = ref([]);
const loading = ref(false);
const submitting = ref(false);
const searchQuery = ref("");
const filterDepartment = ref("");

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
const employeeForm = ref({
  NIK: "",
  NAME: "",
  DEPT_SHORT_NAME: "",
  ENTERANCE_DATE: "",
  TITLE: "",
  GENDER: "",
  RFID_ID: "",
  WORKING_AREA: "",
});

const employeeToDelete = ref(null);
const employeeToHardDelete = ref(null);

// Computed
const filteredEmployees = computed(() => {
  let result = employees.value;

  // Filter by department
  if (filterDepartment.value) {
    result = result.filter(
      (emp) => emp.DEPT_SHORT_NAME === filterDepartment.value
    );
  }

  return result;
});

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

// Methods
const handlePageSizeChange = async () => {
  currentPage.value = 1; // Reset to first page when changing page size
  if (pageSize.value === -1) {
    // If "All" is selected, get total count first
    try {
      const response = await employeeService.getAllEmployees(1, 1);
      if (response.success) {
        pageSize.value = response.data.total_items;
      }
    } catch (error) {
      console.error("Error getting total count:", error);
      pageSize.value = 100; // Fallback to 100 if error
    }
  }
  fetchEmployees();
};

const fetchEmployees = async () => {
  loading.value = true;
  try {
    const data = await employeeService.getAllEmployees(
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

      // Transform data to match frontend property names
      employees.value = (data.data.employees || []).map((emp) => ({
        NIK: emp.nik,
        NAME: emp.name,
        DEPT_SHORT_NAME: emp.dept_short_name,
        ENTERANCE_DATE: emp.enterance_date,
        TITLE: emp.title,
        GENDER: emp.gender,
        RFID_ID: emp.rfid_id,
        WORKING_AREA: emp.working_area,
        DELETE_STATUS: emp.delete_status,
        DEPT_LONG_NAME: emp.dept_long_name,
      }));
    } else {
      console.error("Failed to fetch employees:", data.message);
      toast.error("Gagal memuat data karyawan: " + data.message);
    }
  } catch (error) {
    console.error("Error fetching employees:", error);
    toast.error("Terjadi kesalahan saat memuat data karyawan");
  } finally {
    loading.value = false;
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
    fetchEmployees();
  }, 300);
};

// Watch for search query changes
watch(searchQuery, () => {
  handleSearch();
});

const fetchDepartments = async () => {
  try {
    const data = await departmentService.getAllDepartments(1, 999999); // Get all departments
    if (data.success) {
      departments.value = data.data.departments || [];
    } else {
      console.error("Failed to fetch departments:", data.message);
      toast.error("Gagal memuat data departemen: " + data.message);
    }
  } catch (error) {
    console.error("Error fetching departments:", error);
    toast.error("Terjadi kesalahan saat memuat data departemen");
  }
};

const getDepartmentName = (shortName) => {
  const dept = departments.value.find((d) => d.short_name === shortName);
  return dept ? dept.long_name : shortName;
};

const openCreateModal = () => {
  employeeForm.value = {
    NIK: "HS",
    NAME: "",
    DEPT_SHORT_NAME: "",
    ENTERANCE_DATE: "",
    TITLE: "",
    GENDER: "",
    RFID_ID: "",
    WORKING_AREA: "",
  };
  showCreateModal.value = true;
};

// Tambahkan fungsi untuk memvalidasi dan memformat NIK
const validateAndFormatNIK = (nik) => {
  // Hapus spasi
  let formattedNIK = nik.trim().toUpperCase();

  // Jika NIK kosong, kembalikan HS
  if (!formattedNIK) {
    return "HS";
  }

  // Jika NIK tidak dimulai dengan HS, tambahkan
  if (!formattedNIK.startsWith("HS")) {
    formattedNIK = "HS" + formattedNIK;
  }

  // Jika ada karakter non-alphanumeric selain HS, hapus
  formattedNIK = "HS" + formattedNIK.substring(2).replace(/[^0-9]/g, "");

  // Batasi panjang total menjadi 10 karakter (HS + 8 angka)
  return formattedNIK.substring(0, 10);
};

// Watch perubahan NIK untuk auto format
watch(
  () => employeeForm.value.NIK,
  (newValue) => {
    if (newValue !== undefined && showCreateModal.value) {
      employeeForm.value.NIK = validateAndFormatNIK(newValue);
    }
  }
);

const editEmployee = (employee) => {
  employeeForm.value = {
    NIK: employee.NIK,
    NAME: employee.NAME,
    DEPT_SHORT_NAME: employee.DEPT_SHORT_NAME,
    ENTERANCE_DATE: formatDateForInput(employee.ENTERANCE_DATE),
    TITLE: employee.TITLE,
    GENDER: employee.GENDER,
    RFID_ID: employee.RFID_ID || "",
    WORKING_AREA: employee.WORKING_AREA,
  };
  showEditModal.value = true;
};

const closeModal = () => {
  showCreateModal.value = false;
  showEditModal.value = false;
  employeeForm.value = {
    NIK: "",
    NAME: "",
    DEPT_SHORT_NAME: "",
    ENTERANCE_DATE: "",
    TITLE: "",
    GENDER: "",
    RFID_ID: "",
    WORKING_AREA: "",
  };
};

const createEmployee = async () => {
  // Validasi format NIK
  if (!employeeForm.value.NIK.match(/^HS\d{1,8}$/)) {
    toast.error("Format NIK tidak valid. Harus berformat HS diikuti 1-8 angka");
    return;
  }

  submitting.value = true;
  try {
    const data = await employeeService.createEmployee(employeeForm.value);
    if (data.success) {
      await fetchEmployees();
      closeModal();
      toast.success("Karyawan berhasil ditambahkan");
    } else {
      toast.error("Gagal menambahkan karyawan: " + data.message);
    }
  } catch (error) {
    console.error("Error creating employee:", error);
    toast.error("Terjadi kesalahan saat menambahkan karyawan");
  } finally {
    submitting.value = false;
  }
};

const updateEmployee = async () => {
  submitting.value = true;
  try {
    const data = await employeeService.updateEmployee(
      employeeForm.value.NIK,
      employeeForm.value
    );
    if (data.success) {
      await fetchEmployees();
      closeModal();
      toast.success("Data karyawan berhasil diperbarui");
    } else {
      toast.error("Gagal memperbarui data karyawan: " + data.message);
    }
  } catch (error) {
    console.error("Error updating employee:", error);
    toast.error("Terjadi kesalahan saat memperbarui data karyawan");
  } finally {
    submitting.value = false;
  }
};

const softDeleteEmployee = (employee) => {
  if (!employee || !employee.NIK) {
    console.error("Invalid employee data:", employee);
    toast.error("Data karyawan tidak valid");
    return;
  }
  employeeToDelete.value = employee;
  showSoftDeleteModal.value = true;
};

const confirmSoftDelete = async () => {
  if (!employeeToDelete.value || !employeeToDelete.value.NIK) {
    toast.error("Data karyawan tidak valid");
    return;
  }

  submitting.value = true;
  try {
    const data = await employeeService.deleteEmployee(
      employeeToDelete.value.NIK
    );
    if (data.success) {
      await fetchEmployees();
      showSoftDeleteModal.value = false;
      toast.success("Karyawan berhasil dinonaktifkan");
    } else {
      toast.error("Gagal menonaktifkan karyawan: " + data.message);
    }
  } catch (error) {
    console.error("Error soft deleting employee:", error);
    toast.error("Terjadi kesalahan saat menonaktifkan karyawan");
  } finally {
    submitting.value = false;
  }
};

const showHardDeleteConfirm = (employee) => {
  employeeToHardDelete.value = employee;
  showHardDeleteModal.value = true;
};

const confirmHardDelete = async () => {
  submitting.value = true;
  try {
    const data = await employeeService.hardDeleteEmployee(
      employeeToHardDelete.value.NIK
    );
    if (data.success) {
      await fetchEmployees();
      showHardDeleteModal.value = false;
      toast.success("Karyawan berhasil dihapus secara permanen");
    } else {
      toast.error("Gagal menghapus karyawan: " + data.message);
    }
  } catch (error) {
    console.error("Error hard deleting employee:", error);
    toast.error("Terjadi kesalahan saat menghapus karyawan");
  } finally {
    submitting.value = false;
  }
};

const activateEmployee = async (employee) => {
  submitting.value = true;
  try {
    const data = await employeeService.activateEmployee(employee.NIK);
    if (data.success) {
      await fetchEmployees();
      toast.success("Karyawan berhasil diaktifkan");
    } else {
      toast.error("Gagal mengaktifkan karyawan: " + data.message);
    }
  } catch (error) {
    console.error("Error activating employee:", error);
    toast.error("Terjadi kesalahan saat mengaktifkan karyawan");
  } finally {
    submitting.value = false;
  }
};

const formatDateForInput = (dateString) => {
  if (!dateString) return "";
  try {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
      return "";
    }
    return date.toISOString().split("T")[0];
  } catch (error) {
    return "";
  }
};

// Add pagination methods
const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--;
    fetchEmployees();
  }
};

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
    fetchEmployees();
  }
};

const goToPage = (page) => {
  if (page !== "..." && page !== currentPage.value) {
    currentPage.value = page;
    fetchEmployees();
  }
};

// Lifecycle
onMounted(() => {
  fetchDepartments();
  fetchEmployees();
});
</script>
