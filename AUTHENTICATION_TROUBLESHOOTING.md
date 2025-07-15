# Authentication Troubleshooting Guide

## Error 401 Unauthorized - "Authorization header required"

### Deskripsi Masalah
Error ini terjadi ketika user mencoba mengakses endpoint yang memerlukan authentication tanpa login terlebih dahulu atau token sudah expired.

### Penyebab Umum
1. **User belum login** - Belum ada token yang tersimpan di localStorage
2. **Token expired** - Token JWT sudah kedaluwarsa
3. **Token invalid** - Token rusak atau tidak valid
4. **Session cleared** - Browser cache/localStorage dibersihkan

### Solusi

#### 1. Login Terlebih Dahulu
```
1. Buka aplikasi di http://localhost:5173
2. Klik "Login" atau navigasi ke /login
3. Masukkan credentials yang valid
4. Setelah login berhasil, coba akses halaman yang memerlukan authentication
```

#### 2. Periksa Status Authentication
Buka Developer Tools (F12) dan jalankan di Console:
```javascript
// Periksa apakah ada token
console.log('Token:', localStorage.getItem('token'));

// Periksa user data
console.log('User:', localStorage.getItem('user'));

// Periksa status authentication dari store
console.log('Is Authenticated:', window.__VUE_DEVTOOLS_GLOBAL_HOOK__);
```

#### 3. Clear dan Login Ulang
Jika token ada tapi masih error 401:
```javascript
// Clear semua data authentication
localStorage.removeItem('token');
localStorage.removeItem('user');

// Refresh halaman dan login ulang
location.reload();
```

### Flow Authentication di Aplikasi

#### 1. Route Protection
```javascript
// Routes yang memerlukan authentication:
- /dashboard
- /profile
- /settings
- /users (UserManagement)

// Routes public:
- /login
- /register
```

#### 2. Authentication Check
```javascript
// Auth store mengecek:
const isAuthenticated = computed(() => {
  return !!token.value && !!user.value;
});
```

#### 3. Token Handling
```javascript
// Token dikirim di header setiap request:
config.headers.Authorization = `Bearer ${token}`;
```

### Debugging Steps

#### 1. Periksa Network Tab
1. Buka Developer Tools (F12)
2. Pilih tab "Network"
3. Coba akses halaman yang error
4. Lihat request ke `/api/users`
5. Periksa:
   - Request Headers: Apakah ada `Authorization: Bearer <token>`?
   - Response: Status code dan error message

#### 2. Periksa Console Errors
1. Buka tab "Console" di Developer Tools
2. Lihat error messages
3. Periksa apakah ada error saat:
   - Mengambil token dari localStorage
   - Mengirim request ke API
   - Parsing response

#### 3. Test API Endpoint Langsung
```bash
# Test dengan curl (ganti <token> dengan token yang valid)
curl -H "Authorization: Bearer <token>" http://localhost:8081/api/users

# Test tanpa token (harus return 401)
curl http://localhost:8081/api/users
```

### Common Scenarios

#### Scenario 1: Fresh Browser/Incognito
**Problem**: User membuka aplikasi di browser baru atau incognito mode
**Solution**: Login dengan credentials yang valid

#### Scenario 2: Token Expired
**Problem**: User sudah login sebelumnya tapi token sudah expired
**Solution**: 
1. Aplikasi akan otomatis redirect ke login
2. Login ulang dengan credentials

#### Scenario 3: Server Restart
**Problem**: Backend server di-restart, session hilang
**Solution**: Login ulang (token di frontend masih ada tapi server tidak recognize)

#### Scenario 4: Development Mode
**Problem**: Sering terjadi saat development karena server restart
**Solution**: 
1. Pastikan backend server running
2. Login ulang jika perlu
3. Gunakan persistent session untuk development

### Prevention Tips

#### 1. Implement Token Refresh
```javascript
// TODO: Implement automatic token refresh
// Ketika token hampir expired, request token baru
```

#### 2. Better Error Handling
```javascript
// Interceptor sudah handle 401, tapi bisa diperbaiki:
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Clear token dan redirect ke login
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      router.push("/login");
    }
    return Promise.reject(error);
  }
);
```

#### 3. Session Persistence
```javascript
// Gunakan sessionStorage untuk development
// atau implement "Remember Me" functionality
```

### Quick Fix Commands

#### Reset Authentication State
```javascript
// Jalankan di browser console
localStorage.clear();
sessionStorage.clear();
location.href = '/login';
```

#### Check Backend Status
```bash
# Test backend connection
curl http://localhost:8081/api/hello

# Should return: {"message": "Hello from backend!"}
```

#### Restart Development Servers
```bash
# Backend
cd d:\PROJECT\cors2\backend
go run main.go

# Frontend (new terminal)
cd d:\PROJECT\cors2\frontend
npm run dev
```

### Contact Information
Jika masalah masih berlanjut:
1. Periksa log backend untuk error details
2. Periksa konfigurasi CORS
3. Pastikan database connection normal
4. Verify JWT secret configuration