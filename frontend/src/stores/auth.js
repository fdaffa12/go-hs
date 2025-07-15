import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { authService } from "../services/api";

export const useAuthStore = defineStore("auth", () => {
  // State
  const user = ref(authService.getCurrentUser());
  const token = ref(localStorage.getItem("token"));

  // Getters
  const isAuthenticated = computed(() => !!token.value);

  // Actions
  const init = async () => {
    const savedToken = localStorage.getItem("token");
    const savedUser = authService.getCurrentUser();

    if (savedToken && savedUser) {
      token.value = savedToken;
      user.value = savedUser;
    }
  };

  const register = async (userData) => {
    try {
      const response = await authService.register({
        nik: userData.nik,
        name: userData.name,
        email: userData.email,
        password: userData.password,
        level: userData.level || 1,
      });

      if (response.success) {
        user.value = {
          nik: response.data.nik,
          name: response.data.name,
          email: response.data.email,
          level: response.data.level,
          profile_picture: response.data.profile_picture,
        };
        token.value = response.data.token;
        localStorage.setItem("token", response.data.token);
        localStorage.setItem("user", JSON.stringify(user.value));
      }

      return response;
    } catch (error) {
      console.error("Registration error:", error);
      throw error;
    }
  };

  const login = async (credentials) => {
    try {
      const response = await authService.login(credentials);

      if (response.success) {
        user.value = {
          nik: response.data.nik,
          name: response.data.name,
          email: response.data.email,
          level: response.data.level,
          profile_picture: response.data.profile_picture,
        };
        token.value = response.data.token;
        localStorage.setItem("token", response.data.token);
        localStorage.setItem("user", JSON.stringify(user.value));
      }

      return response;
    } catch (error) {
      console.error("Login error:", error);
      throw error;
    }
  };

  const logout = () => {
    user.value = null;
    token.value = null;
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    authService.logout();
  };

  const updateProfile = async (userData) => {
    try {
      const response = await authService.updateProfile(userData);

      if (response.success) {
        user.value = {
          ...user.value,
          name: response.data.name,
          email: response.data.email,
          profile_picture: response.data.profile_picture,
        };
        localStorage.setItem("user", JSON.stringify(user.value));
      }

      return response;
    } catch (error) {
      console.error("Update profile error:", error);
      throw error;
    }
  };

  return {
    user,
    token,
    isAuthenticated,
    init,
    register,
    login,
    logout,
    updateProfile,
  };
});
