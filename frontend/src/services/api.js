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
  // Get available employees for registration
  getAvailableEmployees: async () => {
    try {
      const response = await api.get("/available-employees");
      if (response.data && response.data.success) {
        return response.data;
      } else {
        throw new Error(response.data?.message || "Failed to fetch employees");
      }
    } catch (error) {
      throw (
        error.response?.data || {
          success: false,
          message: error.message || "Network error",
        }
      );
    }
  },

  // Register user
  register: async (userData) => {
    try {
      const response = await api.post("/register", userData);
      if (response.data.success && response.data.data.token) {
        localStorage.setItem("token", response.data.data.token);
        localStorage.setItem(
          "user",
          JSON.stringify({
            nik: response.data.data.nik,
            name: response.data.data.name,
            email: response.data.data.email,
            level: response.data.data.level,
            profile_picture: response.data.data.profile_picture,
          })
        );
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
      if (response.data.success && response.data.data.token) {
        localStorage.setItem("token", response.data.data.token);
        localStorage.setItem(
          "user",
          JSON.stringify({
            nik: response.data.data.nik,
            name: response.data.data.name,
            email: response.data.data.email,
            level: response.data.data.level,
            profile_picture: response.data.data.profile_picture,
          })
        );
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
  getAllUsers: async (page = 1, pageSize = 10, search = "") => {
    try {
      const response = await api.get(
        `/users?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(
          search
        )}`
      );
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  // Get user by NIK
  getUserByNIK: async (nik) => {
    try {
      const response = await api.get(`/users/${nik}`);
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
      const response = await api.post("/users/create", {
        nik: userData.nik,
        name: userData.name,
        email: userData.email,
        password: userData.password,
        level: userData.level || 1,
      });
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  // Update user
  updateUser: async (nik, userData) => {
    try {
      const response = await api.put(`/users/${nik}`, userData);

      if (!response.data.success) {
        throw new Error(response.data.message || "Failed to update user");
      }

      return response.data;
    } catch (error) {
      if (error.response?.data) {
        throw error.response.data;
      }
      throw { success: false, message: error.message || "Network error" };
    }
  },

  // Update user with profile picture
  updateUserProfile: async (nik, formData) => {
    try {
      const response = await api.put(`/users/${nik}/profile`, formData, {
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

  // Delete user
  deleteUser: async (nik) => {
    try {
      const response = await api.delete(`/users/${nik}`);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  // Change password
  changePassword: async (nik, passwordData) => {
    try {
      const response = await api.put(`/users/${nik}/password`, passwordData);
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
  getAllDepartments: async (page = 1, pageSize = 10, search = "") => {
    try {
      const response = await api.get(
        `/departments?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(
          search
        )}`
      );
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

// Employee Service
export const employeeService = {
  getAllEmployees: async (page = 1, pageSize = 10, search = "") => {
    try {
      const response = await api.get(
        `/employees?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(
          search
        )}`
      );
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  createEmployee: async (employeeData) => {
    try {
      const response = await api.post("/employees", employeeData);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  updateEmployee: async (nik, employeeData) => {
    try {
      const response = await api.put(`/employees/${nik}`, employeeData);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  deleteEmployee: async (nik) => {
    try {
      const response = await api.delete(`/employees/${nik}`);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  hardDeleteEmployee: async (nik) => {
    try {
      const response = await api.delete(`/employees/hard-delete/${nik}`);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  activateEmployee: async (nik) => {
    try {
      const response = await api.put(`/employees/activate/${nik}`);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },
};

// Buyer Service
export const buyerService = {
  getAllBuyers: async (page = 1, pageSize = 10, search = "") => {
    try {
      const response = await api.get(
        `/buyers?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(
          search
        )}`
      );
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  createBuyer: async (buyerData) => {
    try {
      const response = await api.post("/buyers", buyerData);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  updateBuyer: async (shortName, buyerData) => {
    try {
      const response = await api.put(`/buyers/${shortName}`, buyerData);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  deleteBuyer: async (shortName) => {
    try {
      const response = await api.delete(`/buyers/${shortName}`);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  hardDeleteBuyer: async (shortName) => {
    try {
      const response = await api.delete(`/buyers/hard-delete/${shortName}`);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  activateBuyer: async (shortName) => {
    try {
      const response = await api.put(`/buyers/activate/${shortName}`);
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

// Style Service
export const styleService = {
  getAllStyles: async (page = 1, pageSize = 10, search = "") => {
    try {
      const response = await api.get(
        `/styles?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(
          search
        )}`
      );
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  createStyle: async (styleData) => {
    try {
      const response = await api.post("/styles", styleData);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  updateStyle: async (styleNo, styleData) => {
    try {
      const response = await api.put(`/styles/${styleNo}`, styleData);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  deleteStyle: async (styleNo) => {
    try {
      const response = await api.delete(`/styles/${styleNo}`);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  hardDeleteStyle: async (styleNo) => {
    try {
      const response = await api.delete(`/styles/hard-delete/${styleNo}`);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },

  activateStyle: async (styleNo) => {
    try {
      const response = await api.put(`/styles/activate/${styleNo}`);
      return response.data;
    } catch (error) {
      throw (
        error.response?.data || { success: false, message: "Network error" }
      );
    }
  },
};

export default api;
