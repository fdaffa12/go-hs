# API Configuration Guide

## Overview

Aplikasi ini telah dikonfigurasi untuk menggunakan environment variables untuk URL backend API, menggantikan hardcoded URLs yang tersebar di berbagai file.

## Environment Variables

### Frontend Environment Files

#### `.env` (Production/Default)

```env
# Default Environment Variables
# Backend API URL
VITE_API_BASE_URL=http://172.30.177.46:8081

# Development settings
VITE_DEV_MODE=false
```

#### `.env.local` (Local Development)

```env
# Local Development Environment Variables
# This file is used for local development only

# Backend API URL for local development
VITE_API_BASE_URL=http://localhost:8081

# Development settings
VITE_DEV_MODE=true
```

## Centralized API Configuration

### API Service (`src/services/api.js`)

Semua komunikasi dengan backend sekarang menggunakan konfigurasi terpusat:

```javascript
// Function to get API base URL from environment
const getApiBaseUrl = () => {
  const baseUrl = import.meta.env.VITE_API_BASE_URL || "http://localhost:8081";
  return `${baseUrl}/api`;
};
```

### Services Available

- `authService`: Login, Register, Logout
- `userService`: CRUD operations untuk user management
- `testConnection`: Test koneksi ke backend

## Files Updated

Berikut adalah file-file yang telah diperbarui untuk menggunakan konfigurasi terpusat:

### 1. `src/services/api.js`

- ✅ Menambahkan fungsi `getApiBaseUrl()`
- ✅ Menggunakan environment variable untuk base URL
- ✅ Memperbaiki endpoint `/users/create`

### 2. `src/stores/auth.js`

- ✅ Menggunakan `authService` dari `api.js`
- ✅ Menghapus hardcoded URLs
- ✅ Menggunakan environment variable untuk `getCurrentUser`

### 3. `src/pages/UserManagement.vue`

- ✅ Menggunakan `userService` dari `api.js`
- ✅ Menghapus semua hardcoded URLs
- ✅ Menggunakan service methods untuk CRUD operations

### 4. `src/pages/Profile.vue`

- ✅ Menggunakan environment variable untuk profile image URLs
- ✅ Menggunakan `userService` untuk update operations

### 5. `src/pages/Dashboard.vue`

- ✅ Menggunakan environment variable untuk test connection

### 6. `src/layouts/AuthenticatedLayout.vue`

- ✅ Menggunakan environment variable untuk profile image URLs

## Benefits

### 1. **Maintainability**

- URL backend hanya perlu diubah di satu tempat (environment variables)
- Tidak ada hardcoded URLs yang tersebar di berbagai file

### 2. **Flexibility**

- Mudah beralih antara development dan production environments
- Support untuk multiple environments (local, staging, production)

### 3. **Consistency**

- Semua API calls menggunakan konfigurasi yang sama
- Centralized error handling dan interceptors

### 4. **Security**

- Environment variables tidak ter-commit ke repository
- Sensitive URLs dapat disimpan secara terpisah

## Usage

### Development Mode

```bash
# Menggunakan localhost untuk development
cp .env.local .env
npm run dev
```

### Production Mode

```bash
# Menggunakan IP remote untuk production
# .env sudah dikonfigurasi dengan IP production
npm run build
npm run preview
```

### Switch Mode Script

Gunakan script `switch-mode.bat` untuk beralih antara mode:

```bash
# Windows
.\switch-mode.bat

# Pilih:
# 1. Local Mode (localhost:8081)
# 2. Remote Mode (172.30.177.46:8081)
```

## API Endpoints

### Authentication

- `POST /api/register` - Register user baru
- `POST /api/login` - Login user
- `GET /api/user` - Get current user info

### User Management

- `GET /api/users` - Get all users
- `POST /api/users/create` - Create new user
- `GET /api/users/{id}` - Get user by ID
- `PUT /api/users/{id}` - Update user
- `PUT /api/users/{id}/profile` - Update user profile with image
- `PUT /api/users/{id}/password` - Change user password
- `DELETE /api/users/{id}` - Delete user

### Testing

- `GET /api/hello` - Test backend connection (public)
- `GET /api/test` - Test backend connection (authenticated)

## Troubleshooting

### 1. API Connection Issues

- Periksa `VITE_API_BASE_URL` di file `.env`
- Pastikan backend server berjalan di URL yang benar
- Gunakan browser console untuk melihat network errors

### 2. Environment Variables Not Working

- Restart development server setelah mengubah `.env`
- Pastikan variable name dimulai dengan `VITE_`
- Periksa file `.env` ada di root directory frontend

### 3. CORS Issues

- Pastikan backend CORS middleware mengizinkan origin frontend
- Periksa konfigurasi CORS di `backend/middleware/auth.go`

## Best Practices

1. **Selalu gunakan environment variables** untuk URLs dan konfigurasi
2. **Jangan commit file `.env.local`** ke repository
3. **Gunakan centralized API service** untuk semua HTTP requests
4. **Test di berbagai environments** sebelum deployment
5. **Dokumentasikan perubahan environment variables** untuk tim
