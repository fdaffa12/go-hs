<script setup>
import { ref, reactive } from "vue";
import { useAuthStore } from "../stores/auth";
import { useToast } from "vue-toastification";
import AuthenticatedLayout from "../layouts/AuthenticatedLayout.vue";
import { userService } from "../services/api";

// Initialize stores and utilities
const authStore = useAuthStore();
const toast = useToast();

// Reactive data
const loading = ref(false);
const loadingMessage = ref("");
const showEditModal = ref(false);
const showPasswordModal = ref(false);
const editForm = reactive({
  username: "",
  email: "",
  profileImage: null,
});
const profileImagePreview = ref(null);
const passwordForm = reactive({
  currentPassword: "",
  newPassword: "",
  confirmPassword: "",
});

// Methods
const formatDate = (dateString) => {
  if (!dateString) return "Unknown";
  return new Date(dateString).toLocaleDateString("id-ID", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
};

const getProfileImageUrl = (profilePicture) => {
  if (!profilePicture) return null;
  // If it's already a full URL, return as is
  if (profilePicture.startsWith("http")) {
    return profilePicture;
  }
  // Otherwise, prepend the backend URL from environment
  const baseUrl = import.meta.env.VITE_API_BASE_URL || "http://localhost:8081";
  return `${baseUrl}${profilePicture}`;
};

const openEditModal = () => {
  // Pre-fill form with current user data
  editForm.username = authStore.user?.username || "";
  editForm.email = authStore.user?.email || "";
  showEditModal.value = true;
};

const closeEditModal = () => {
  showEditModal.value = false;
  // Reset form
  editForm.username = "";
  editForm.email = "";
  editForm.profileImage = null;
  profileImagePreview.value = null;
};

const handleImageUpload = (event) => {
  const file = event.target.files[0];
  if (!file) return;

  // Validate file type
  const allowedTypes = ["image/jpeg", "image/jpg", "image/png", "image/gif"];
  if (!allowedTypes.includes(file.type)) {
    toast.error("Format file tidak didukung. Gunakan JPG, PNG, atau GIF.");
    return;
  }

  // Validate file size (2MB)
  const maxSize = 2 * 1024 * 1024; // 2MB in bytes
  if (file.size > maxSize) {
    toast.error("Ukuran file terlalu besar. Maksimal 2MB.");
    return;
  }

  // Store file and create preview
  editForm.profileImage = file;

  // Create preview URL
  const reader = new FileReader();
  reader.onload = (e) => {
    profileImagePreview.value = e.target.result;
  };
  reader.readAsDataURL(file);
};

const removeProfileImage = async () => {
  try {
    // If user has an existing profile picture, delete it from server
    if (authStore.user?.profile_picture && !profileImagePreview.value) {
      loading.value = true;
      loadingMessage.value = "Menghapus foto profil...";

      const response = await userService.updateProfile(authStore.user.id, {
        username: authStore.user.username,
        email: authStore.user.email,
        remove_profile_picture: true,
      });

      if (response.success) {
        // Update auth store with new user data
        await authStore.fetchUser();
        toast.success("Foto profil berhasil dihapus");
      } else {
        toast.error(response.message || "Gagal menghapus foto profil");
        return;
      }
    }

    // Clear form data and preview
    editForm.profileImage = null;
    profileImagePreview.value = null;

    // Clear file input
    const fileInput = document.querySelector('input[type="file"]');
    if (fileInput) {
      fileInput.value = "";
    }
  } catch (error) {
    console.error("Error removing profile image:", error);
    toast.error("Terjadi kesalahan saat menghapus foto profil");
  } finally {
    loading.value = false;
    loadingMessage.value = "";
  }
};

const updateProfile = async () => {
  try {
    loading.value = true;
    loadingMessage.value = "Memperbarui profile...";

    // Validate form
    if (!editForm.username.trim() || !editForm.email.trim()) {
      toast.error("Username dan email harus diisi");
      return;
    }

    // Prepare data for API call
    let response;

    if (editForm.profileImage) {
      // If there's a profile image, use FormData
      const formData = new FormData();
      formData.append("username", editForm.username.trim());
      formData.append("email", editForm.email.trim());
      formData.append("profile_picture", editForm.profileImage);

      response = await userService.updateProfile(authStore.user.id, formData);
    } else {
      // Regular update without image
      const updateData = {
        username: editForm.username.trim(),
        email: editForm.email.trim(),
      };

      response = await userService.updateProfile(authStore.user.id, updateData);
    }

    if (response.success) {
      // Update auth store with new user data
      await authStore.fetchUser();

      toast.success("Profile berhasil diperbarui");
      closeEditModal();
    } else {
      toast.error(response.message || "Gagal memperbarui profile");
    }
  } catch (error) {
    console.error("Error updating profile:", error);
    toast.error("Terjadi kesalahan saat memperbarui profile");
  } finally {
    loading.value = false;
    loadingMessage.value = "";
  }
};

const openPasswordModal = () => {
  // Reset form
  passwordForm.currentPassword = "";
  passwordForm.newPassword = "";
  passwordForm.confirmPassword = "";
  showPasswordModal.value = true;
};

const closePasswordModal = () => {
  showPasswordModal.value = false;
  // Reset form
  passwordForm.currentPassword = "";
  passwordForm.newPassword = "";
  passwordForm.confirmPassword = "";
};

const changePassword = async () => {
  try {
    loading.value = true;
    loadingMessage.value = "Mengubah password...";

    // Validate form
    if (
      !passwordForm.currentPassword ||
      !passwordForm.newPassword ||
      !passwordForm.confirmPassword
    ) {
      toast.error("Semua field password harus diisi");
      return;
    }

    if (passwordForm.newPassword.length < 6) {
      toast.error("Password baru minimal 6 karakter");
      return;
    }

    if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      toast.error("Konfirmasi password tidak cocok");
      return;
    }

    // Call API to change password
    const response = await userService.changePassword(authStore.user.id, {
      currentPassword: passwordForm.currentPassword,
      newPassword: passwordForm.newPassword,
    });

    if (response.success) {
      toast.success("Password berhasil diubah");
      closePasswordModal();
    } else {
      toast.error(response.message || "Gagal mengubah password");
    }
  } catch (error) {
    console.error("Error changing password:", error);
    toast.error("Terjadi kesalahan saat mengubah password");
  } finally {
    loading.value = false;
    loadingMessage.value = "";
  }
};
</script>

<template>
  <AuthenticatedLayout :user="authStore.user">
    <div class="space-y-8">
      <!-- Profile Header -->
      <section
        class="bg-white rounded-xl shadow-lg border border-gray-200 p-4 sm:p-6 lg:p-8 relative overflow-hidden"
      >
        <!-- Background Pattern -->
        <div class="absolute inset-0 bg-gray-50">
          <div
            class="absolute inset-0"
            style="
              background-image: radial-gradient(
                  circle at 25% 25%,
                  rgba(100, 116, 139, 0.05) 0%,
                  transparent 50%
                ),
                radial-gradient(
                  circle at 75% 75%,
                  rgba(100, 116, 139, 0.05) 0%,
                  transparent 50%
                );
            "
          ></div>
        </div>

        <div
          class="relative z-10 flex flex-col sm:flex-row items-center sm:items-center space-y-4 sm:space-y-0 sm:space-x-6 lg:space-x-8"
        >
          <div class="relative flex-shrink-0">
            <img
              v-if="authStore.user?.profile_picture"
              :src="getProfileImageUrl(authStore.user.profile_picture)"
              alt="Profile Picture"
              class="w-20 h-20 sm:w-24 sm:h-24 rounded-full object-cover border-4 border-gray-200 shadow-xl"
            />
            <div
              v-else
              class="w-20 h-20 sm:w-24 sm:h-24 bg-primary-600 rounded-full flex items-center justify-center border-4 border-gray-200 shadow-xl"
            >
              <span class="text-2xl sm:text-3xl font-bold text-white">
                {{ authStore.user?.username?.charAt(0).toUpperCase() || "U" }}
              </span>
            </div>
            <!-- Online indicator -->
            <div
              class="absolute -bottom-1 -right-1 w-5 h-5 sm:w-6 sm:h-6 bg-green-500 rounded-full border-4 border-white shadow-lg"
            >
              <div
                class="w-full h-full bg-green-500 rounded-full animate-ping opacity-75"
              ></div>
            </div>
          </div>

          <div class="flex-1 text-center sm:text-left">
            <div
              class="flex flex-col sm:flex-row sm:items-center space-y-2 sm:space-y-0 sm:space-x-3 mb-2"
            >
              <h2 class="text-2xl sm:text-3xl font-bold text-gray-900">
                {{ authStore.user?.username || "User" }}
              </h2>
              <span
                class="inline-flex items-center justify-center sm:justify-start px-3 py-1 rounded-full text-xs font-medium bg-primary-100 text-primary-800 border border-primary-200"
              >
                <svg
                  class="w-3 h-3 mr-1"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path
                    fill-rule="evenodd"
                    d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                    clip-rule="evenodd"
                  />
                </svg>
                Verified
              </span>
            </div>

            <p class="text-gray-600 text-base sm:text-lg mb-2 break-all">
              {{ authStore.user?.email || "user@example.com" }}
            </p>

            <div
              class="flex flex-col sm:flex-row sm:items-center space-y-2 sm:space-y-0 sm:space-x-4 text-sm text-gray-500"
            >
              <div
                class="flex items-center justify-center sm:justify-start space-x-2"
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
                    d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                  />
                </svg>
                <span
                  >Member since
                  {{ formatDate(authStore.user?.created_at) }}</span
                >
              </div>

              <div
                class="flex items-center justify-center sm:justify-start space-x-2"
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
                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
                <span>Last active: Today</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Profile Information -->
      <section
        class="bg-white rounded-lg shadow-sm border border-gray-200 p-4 sm:p-6"
      >
        <h3 class="text-lg font-semibold text-gray-900 mb-4 sm:mb-6">
          Informasi Profile
        </h3>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-6">
          <!-- Username Card -->
          <div class="bg-gray-50 rounded-lg p-3 sm:p-4 border border-gray-200">
            <div class="flex items-center space-x-3">
              <div class="flex-shrink-0">
                <svg
                  class="w-6 h-6 sm:w-8 sm:h-8 text-primary-600"
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
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-gray-600">Username</p>
                <p
                  class="text-base sm:text-lg font-semibold text-gray-900 truncate"
                >
                  {{ authStore.user?.username || "N/A" }}
                </p>
              </div>
            </div>
          </div>

          <!-- Email Card -->
          <div class="bg-gray-50 rounded-lg p-3 sm:p-4 border border-gray-200">
            <div class="flex items-center space-x-3">
              <div class="flex-shrink-0">
                <svg
                  class="w-6 h-6 sm:w-8 sm:h-8 text-primary-600"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                  />
                </svg>
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-gray-600">Email</p>
                <p
                  class="text-base sm:text-lg font-semibold text-gray-900 truncate"
                >
                  {{ authStore.user?.email || "N/A" }}
                </p>
              </div>
            </div>
          </div>

          <!-- User ID Card -->
          <div class="bg-gray-50 rounded-lg p-3 sm:p-4 border border-gray-200">
            <div class="flex items-center space-x-3">
              <div class="flex-shrink-0">
                <svg
                  class="w-6 h-6 sm:w-8 sm:h-8 text-primary-600"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"
                  />
                </svg>
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-gray-600">User ID</p>
                <p class="text-base sm:text-lg font-semibold text-gray-900">
                  #{{ authStore.user?.id || "N/A" }}
                </p>
              </div>
            </div>
          </div>

          <!-- Status Card -->
          <div class="bg-gray-50 rounded-lg p-3 sm:p-4 border border-gray-200">
            <div class="flex items-center space-x-3">
              <div class="flex-shrink-0">
                <svg
                  class="w-6 h-6 sm:w-8 sm:h-8 text-primary-600"
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
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-gray-600">Status</p>
                <div class="flex items-center mt-1">
                  <span
                    class="inline-flex items-center px-2 sm:px-3 py-1 rounded-full text-xs sm:text-sm font-medium bg-green-100 text-green-800"
                  >
                    <span
                      class="w-2 h-2 bg-green-500 rounded-full mr-2 animate-pulse"
                    ></span>
                    Active
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Account Actions -->
      <section
        class="bg-white rounded-lg shadow-sm border border-gray-200 p-4 sm:p-6"
      >
        <h3 class="text-lg font-semibold text-gray-900 mb-4 sm:mb-6">
          Aksi Akun
        </h3>

        <div class="space-y-3">
          <!-- Edit Profile Action -->
          <button
            @click="openEditModal"
            class="w-full flex items-center justify-between p-3 sm:p-4 text-left border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors duration-200"
          >
            <div class="flex items-center space-x-3">
              <svg
                class="w-5 h-5 text-gray-600 flex-shrink-0"
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
              <div class="min-w-0 flex-1">
                <h4 class="font-medium text-gray-900">Edit Profile</h4>
                <p class="text-sm text-gray-500 hidden sm:block">
                  Update your personal information
                </p>
              </div>
            </div>
            <svg
              class="w-5 h-5 text-gray-400 flex-shrink-0"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 5l7 7-7 7"
              />
            </svg>
          </button>

          <!-- Change Password Action -->
          <button
            @click="openPasswordModal"
            class="w-full flex items-center justify-between p-3 sm:p-4 text-left border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors duration-200"
          >
            <div class="flex items-center space-x-3">
              <svg
                class="w-5 h-5 text-gray-600 flex-shrink-0"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2-2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                />
              </svg>
              <div class="min-w-0 flex-1">
                <h4 class="font-medium text-gray-900">Change Password</h4>
                <p class="text-sm text-gray-500 hidden sm:block">
                  Update your account security
                </p>
              </div>
            </div>
            <svg
              class="w-5 h-5 text-gray-400 flex-shrink-0"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 5l7 7-7 7"
              />
            </svg>
          </button>
        </div>
      </section>
    </div>

    <!-- Edit Profile Modal -->
    <div
      v-if="showEditModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      @click="closeEditModal"
    >
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md" @click.stop>
        <div class="p-6">
          <div class="flex justify-between items-center mb-6">
            <h3 class="text-lg font-semibold text-gray-900">Edit Profile</h3>
            <button
              @click="closeEditModal"
              class="text-gray-400 hover:text-gray-600"
            >
              <svg
                class="w-6 h-6"
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

          <form @submit.prevent="updateProfile">
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Username
                </label>
                <input
                  v-model="editForm.username"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                  placeholder="Masukkan username"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Email
                </label>
                <input
                  v-model="editForm.email"
                  type="email"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                  placeholder="Masukkan email"
                />
              </div>

              <!-- Profile Picture Upload -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Foto Profil
                </label>
                <div class="flex items-center space-x-4">
                  <!-- Current/Preview Image -->
                  <div class="relative">
                    <img
                      v-if="
                        profileImagePreview || authStore.user?.profile_picture
                      "
                      :src="
                        profileImagePreview ||
                        getProfileImageUrl(authStore.user?.profile_picture)
                      "
                      alt="Profile Picture"
                      class="w-20 h-20 rounded-full object-cover border-2 border-gray-300"
                    />
                    <div
                      v-else
                      class="w-20 h-20 rounded-full bg-primary-100 flex items-center justify-center border-2 border-gray-300"
                    >
                      <span class="text-2xl font-semibold text-primary-600">
                        {{
                          authStore.user?.username?.charAt(0).toUpperCase() ||
                          "U"
                        }}
                      </span>
                    </div>
                  </div>

                  <!-- Upload Controls -->
                  <div class="flex-1">
                    <input
                      ref="fileInput"
                      type="file"
                      accept="image/*"
                      @change="handleImageUpload"
                      class="hidden"
                    />
                    <div class="space-y-2">
                      <button
                        type="button"
                        @click="$refs.fileInput.click()"
                        class="px-3 py-2 text-sm font-medium text-primary-600 bg-primary-50 border border-primary-200 rounded-lg hover:bg-primary-100"
                      >
                        {{ profileImagePreview ? "Ganti Foto" : "Upload Foto" }}
                      </button>
                      <button
                        v-if="
                          profileImagePreview || authStore.user?.profile_picture
                        "
                        type="button"
                        @click="removeProfileImage"
                        class="ml-2 px-3 py-2 text-sm font-medium text-red-600 bg-red-50 border border-red-200 rounded-lg hover:bg-red-100"
                      >
                        Hapus Foto
                      </button>
                    </div>
                    <p class="text-xs text-gray-500 mt-1">
                      Format: JPG, PNG, GIF. Maksimal 2MB.
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <div class="flex justify-end space-x-3 mt-6">
              <button
                type="button"
                @click="closeEditModal"
                class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 border border-gray-300 rounded-lg hover:bg-gray-200"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="loading"
                class="px-4 py-2 text-sm font-medium text-white bg-primary-600 border border-transparent rounded-lg hover:bg-primary-700 disabled:opacity-50"
              >
                {{ loading ? "Menyimpan..." : "Simpan" }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Change Password Modal -->
    <div
      v-if="showPasswordModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      @click="closePasswordModal"
    >
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md" @click.stop>
        <div class="p-6">
          <div class="flex justify-between items-center mb-6">
            <h3 class="text-lg font-semibold text-gray-900">Change Password</h3>
            <button
              @click="closePasswordModal"
              class="text-gray-400 hover:text-gray-600"
            >
              <svg
                class="w-6 h-6"
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

          <form @submit.prevent="changePassword">
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Password Lama
                </label>
                <input
                  v-model="passwordForm.currentPassword"
                  type="password"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                  placeholder="Masukkan password lama"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Password Baru
                </label>
                <input
                  v-model="passwordForm.newPassword"
                  type="password"
                  required
                  minlength="6"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                  placeholder="Masukkan password baru (min. 6 karakter)"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Konfirmasi Password Baru
                </label>
                <input
                  v-model="passwordForm.confirmPassword"
                  type="password"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                  placeholder="Konfirmasi password baru"
                />
              </div>
            </div>

            <div class="flex justify-end space-x-3 mt-6">
              <button
                type="button"
                @click="closePasswordModal"
                class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 border border-gray-300 rounded-lg hover:bg-gray-200"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="loading"
                class="px-4 py-2 text-sm font-medium text-white bg-primary-600 border border-transparent rounded-lg hover:bg-primary-700 disabled:opacity-50"
              >
                {{ loading ? "Mengubah..." : "Ubah Password" }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Loading Overlay -->
    <div
      v-if="loading"
      class="fixed inset-0 bg-black bg-opacity-25 flex items-center justify-center z-40"
    >
      <div class="bg-white rounded-lg p-6 shadow-xl">
        <div class="flex items-center space-x-3">
          <div
            class="animate-spin rounded-full h-6 w-6 border-b-2 border-primary-600"
          ></div>
          <span class="text-gray-700">{{ loadingMessage }}</span>
        </div>
      </div>
    </div>
  </AuthenticatedLayout>
</template>

<style scoped>
/* Styles are now handled by Tailwind CSS and global.css */
</style>
