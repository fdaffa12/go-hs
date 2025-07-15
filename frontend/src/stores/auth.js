import { ref, computed } from "vue";
import { defineStore } from "pinia";
import { authService } from "../services/api";

export const useAuthStore = defineStore("auth", () => {
  // State
  const user = ref(null);
  const token = ref(localStorage.getItem("token") || null);
  const isLoading = ref(false);
  const error = ref(null);

  // Getters
  const isAuthenticated = computed(() => {
    return !!token.value && !!user.value;
  });

  const currentUser = computed(() => {
    return user.value;
  });

  // Actions
  const login = async (credentials) => {
    try {
      isLoading.value = true;
      error.value = null;

      const result = await authService.login(credentials);

      if (result.success) {
        token.value = result.data.token;
        user.value = result.data.user;
        // Ensure data is saved to localStorage
        localStorage.setItem("token", result.data.token);
        localStorage.setItem("user", JSON.stringify(result.data.user));
        return { success: true, data: result.data };
      } else {
        error.value = result.message || "Login failed";
        return { success: false, message: result.message };
      }
    } catch (err) {
      error.value = err.message || "Network error occurred";
      return {
        success: false,
        message: err.message || "Network error occurred",
      };
    } finally {
      isLoading.value = false;
    }
  };

  const register = async (userData) => {
    try {
      isLoading.value = true;
      error.value = null;

      const result = await authService.register(userData);

      if (result.success) {
        return { success: true, data: result };
      } else {
        error.value = result.message || "Registration failed";
        return { success: false, message: result.message };
      }
    } catch (err) {
      error.value = err.message || "Network error occurred";
      return {
        success: false,
        message: err.message || "Network error occurred",
      };
    } finally {
      isLoading.value = false;
    }
  };

  const logout = () => {
    user.value = null;
    token.value = null;
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    error.value = null;
  };

  const getCurrentUser = async () => {
    if (!token.value) return;

    try {
      isLoading.value = true;
      // Use direct fetch since we need /user endpoint, not /users
      const baseUrl =
        import.meta.env.VITE_API_BASE_URL || "http://localhost:8081";
      const response = await fetch(`${baseUrl}/api/user`, {
        method: "GET",
        headers: {
          Authorization: `Bearer ${token.value}`,
          "Content-Type": "application/json",
        },
      });

      const data = await response.json();

      if (response.ok) {
        user.value = data.data;
      } else {
        // Token might be invalid, logout
        logout();
      }
    } catch (err) {
      console.error("Failed to get current user:", err);
      logout();
    } finally {
      isLoading.value = false;
    }
  };

  const clearError = () => {
    error.value = null;
  };

  // Initialize store
  const init = async () => {
    const storedToken = localStorage.getItem("token");
    const storedUser = localStorage.getItem("user");
    
    if (storedToken) {
      token.value = storedToken;
      
      // Try to get user from localStorage first
      if (storedUser) {
        try {
          user.value = JSON.parse(storedUser);
        } catch (err) {
          console.error("Failed to parse stored user:", err);
        }
      }
      
      // Verify token is still valid by getting current user
      try {
        await getCurrentUser();
      } catch (err) {
        // Token might be invalid, clear it
        logout();
      }
    }
  };

  // Alias for getCurrentUser to match Profile.vue expectation
  const fetchUser = getCurrentUser;

  return {
    // State
    user,
    token,
    isLoading,
    error,
    // Getters
    isAuthenticated,
    currentUser,
    // Actions
    init,
    login,
    register,
    logout,
    getCurrentUser,
    fetchUser,
    clearError,
  };
});
