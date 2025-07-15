import axios from "axios";

// Function to get API base URL from environment
const getApiBaseUrl = () => {
  const baseUrl = import.meta.env.VITE_API_BASE_URL || "http://localhost:8081";
  return `${baseUrl}/api`;
};

// Base URL untuk API backend
const API_BASE_URL = getApiBaseUrl();

// Helper function to get auth header
const getAuthHeader = () => {
  const token = localStorage.getItem("token");
  return {
    Authorization: `Bearer ${token}`,
    Accept: "application/json",
  };
};

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor untuk menambahkan token ke header
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor untuk handle error
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response?.status === 401) {
      // Token expired atau invalid, hapus token dan redirect ke login
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      // Bisa tambahkan redirect ke login page di sini
    }
    return Promise.reject(error);
  }
);

// Auth service functions
export const authService = {
  // Register user
  register: async (userData) => {
    try {
      const response = await api.post("/register", userData);
      if (response.data.success && response.data.token) {
        localStorage.setItem("token", response.data.token);
        localStorage.setItem("user", JSON.stringify(response.data.data));
      }
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  // Login user
  login: async (credentials) => {
    try {
      const response = await api.post("/login", credentials);
      if (response.data.success && response.data.token) {
        localStorage.setItem("token", response.data.token);
        localStorage.setItem("user", JSON.stringify(response.data.user));
      }
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  // Logout user
  logout: () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
  },

  // Get current user from localStorage
  getCurrentUser: () => {
    const user = localStorage.getItem("user");
    return user ? JSON.parse(user) : null;
  },

  // Check if user is authenticated
  isAuthenticated: () => {
    return !!localStorage.getItem("token");
  },
};

// User service functions
export const userService = {
  // Get all users
  getAllUsers: async () => {
    try {
      const response = await api.get("/users");
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  // Get user by ID
  getUserById: async (id) => {
    try {
      const response = await api.get(`/users/${id}`);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  // Create new user
  createUser: async (userData) => {
    try {
      const response = await api.post("/users/create", userData);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  // Update user
  updateUser: async (id, userData) => {
    try {
      const response = await api.put(`/users/${id}`, userData);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  // Update user with image
  updateUserWithImage: async (id, formData) => {
    try {
      const response = await api.put(`/users/${id}/profile`, formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  // Update profile
  updateProfile: async (id, data) => {
    try {
      let requestData = data;
      let contentType = "multipart/form-data";

      // If data is not FormData, convert it to FormData
      if (!(data instanceof FormData)) {
        requestData = new FormData();
        Object.keys(data).forEach((key) => {
          requestData.append(key, data[key]);
        });
      }

      const response = await api.put(`/users/${id}/profile`, requestData, {
        headers: {
          "Content-Type": contentType,
        },
      });
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  // Delete user
  deleteUser: async (id) => {
    try {
      const response = await api.delete(`/users/${id}`);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  // Change password
  changePassword: async (id, passwordData) => {
    try {
      const response = await api.put(`/users/${id}/password`, passwordData);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },
};

// Department Service
export const departmentService = {
  getAllDepartments: async () => {
    try {
      const response = await api.get("/departments");
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  createDepartment: async (departmentData) => {
    try {
      const response = await api.post("/departments", departmentData);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  updateDepartment: async (shortName, departmentData) => {
    try {
      const response = await api.put(
        `/departments/${shortName}`,
        departmentData
      );
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  deleteDepartment: async (shortName) => {
    try {
      const response = await api.delete(`/departments/${shortName}`);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  hardDeleteDepartment: async (shortName) => {
    try {
      const response = await api.delete(
        `/departments/hard-delete/${shortName}`
      );
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  activateDepartment: async (shortName) => {
    try {
      const response = await api.put(`/departments/activate/${shortName}`);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },
};

// Test connection
export const testConnection = async () => {
  try {
    const response = await api.get("/hello");
    return response.data;
  } catch (error) {
    throw error.response?.data || { message: "Network error" };
  }
};

export default api;
