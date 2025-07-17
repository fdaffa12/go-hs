<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Top Navigation Bar -->
    <header
      class="bg-white border-b border-gray-200 shadow-sm sticky top-0 z-50"
    >
      <div class="px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <!-- Mobile menu button -->
          <button
            @click="sidebarOpen = true"
            class="lg:hidden p-2 rounded-md text-gray-400 hover:text-gray-600 hover:bg-gray-100"
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
                d="M4 6h16M4 12h16M4 18h16"
              />
            </svg>
          </button>

          <!-- Logo and App Name (Mobile) -->
          <div class="flex items-center space-x-3 lg:hidden">
            <div
              class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center"
            >
              <svg
                class="w-5 h-5 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M13 10V3L4 14h7v7l9-11h-7z"
                />
              </svg>
            </div>
            <h1 class="text-lg font-bold text-gray-900">CORS App</h1>
          </div>

          <!-- Desktop Hamburger and Page Title -->
          <div class="hidden lg:flex items-center space-x-4">
            <button
              @click="sidebarCollapsed = !sidebarCollapsed"
              class="p-2 rounded-md text-gray-400 hover:text-gray-600 hover:bg-gray-100"
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
                  d="M4 6h16M4 12h16M4 18h16"
                />
              </svg>
            </button>
            <h2 class="text-lg font-semibold text-gray-900">
              {{ pageTitle }}
            </h2>
          </div>

          <!-- Right side actions -->
          <div class="flex items-center space-x-4">
            <!-- Notifications -->
            <button
              class="p-2 rounded-md text-gray-400 hover:text-gray-600 hover:bg-gray-100 relative"
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
                  d="M15 17h5l-5 5v-5zM10.5 3.5a6 6 0 0 1 6 6v2l1.5 1.5v1h-13v-1L6.5 11.5v-2a6 6 0 0 1 6-6z"
                />
              </svg>
              <!-- Notification badge -->
              <span
                class="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"
              ></span>
            </button>

            <!-- User menu (Desktop) -->
            <div class="hidden sm:flex items-center space-x-3">
              <div class="text-sm text-gray-600">
                Selamat datang,
                <span class="font-medium text-gray-900">{{
                  user?.name || "User"
                }}</span>
              </div>
              <!-- Profile Picture -->
              <div
                class="w-8 h-8 rounded-full overflow-hidden border-2 border-gray-200"
              >
                <img
                  v-if="user?.profile_picture"
                  :src="getProfileImageUrl(user.profile_picture)"
                  :alt="user?.name || 'User'"
                  class="w-full h-full object-cover"
                />
                <div
                  v-else
                  class="w-full h-full bg-primary-100 flex items-center justify-center"
                >
                  <span class="text-xs font-medium text-primary-600">
                    {{ user?.name?.charAt(0).toUpperCase() || "U" }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- Sidebar Overlay (Mobile) -->
    <div
      v-if="sidebarOpen"
      @click="sidebarOpen = false"
      class="fixed inset-0 z-40 bg-black bg-opacity-50 lg:hidden"
    ></div>

    <!-- Sidebar -->
    <aside
      :class="[
        'fixed top-16 bottom-0 left-0 z-50 bg-white border-r border-gray-200 transform transition-all duration-300 ease-in-out lg:top-0 lg:z-30',
        'lg:translate-x-0',
        sidebarOpen ? 'translate-x-0' : '-translate-x-full',
        sidebarCollapsed ? 'lg:w-16' : 'lg:w-64',
        'w-64',
      ]"
    >
      <!-- Sidebar Header (Desktop only) -->
      <div
        class="hidden lg:flex items-center justify-center h-16 px-6 border-b border-gray-200"
      >
        <div class="flex items-center space-x-3">
          <div
            class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center"
          >
            <svg
              class="w-5 h-5 text-white"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M13 10V3L4 14h7v7l9-11h-7z"
              />
            </svg>
          </div>
          <h1
            v-if="!sidebarCollapsed"
            class="text-lg font-bold text-gray-900 transition-opacity duration-300"
          >
            CORS App
          </h1>
        </div>
      </div>

      <!-- Mobile Header -->
      <div
        class="lg:hidden flex items-center justify-between p-4 border-b border-gray-200"
      >
        <div class="flex items-center space-x-3">
          <div
            class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center"
          >
            <svg
              class="w-5 h-5 text-white"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M13 10V3L4 14h7v7l9-11h-7z"
              />
            </svg>
          </div>
          <h1 class="text-lg font-bold text-gray-900">CORS App</h1>
        </div>
        <button
          @click="sidebarOpen = false"
          class="p-2 rounded-md text-gray-400 hover:text-gray-600 hover:bg-gray-100"
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

      <!-- Navigation Menu -->
      <nav class="flex-1 px-4 py-4 lg:py-6 space-y-2 overflow-y-auto">
        <router-link
          to="/dashboard"
          :class="[
            'nav-item',
            { 'nav-item-active': $route.path === '/dashboard' },
            sidebarCollapsed ? 'lg:justify-center lg:px-2' : '',
          ]"
          :title="sidebarCollapsed ? 'Dashboard' : ''"
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
              d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"
            />
          </svg>
          <span v-if="!sidebarCollapsed" class="transition-opacity duration-300"
            >Dashboard</span
          >
        </router-link>

        <router-link
          to="/profile"
          :class="[
            'nav-item',
            { 'nav-item-active': $route.path === '/profile' },
            sidebarCollapsed ? 'lg:justify-center lg:px-2' : '',
          ]"
          :title="sidebarCollapsed ? 'Profile' : ''"
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
              d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
            />
          </svg>
          <span v-if="!sidebarCollapsed" class="transition-opacity duration-300"
            >Profile</span
          >
        </router-link>

        <router-link
          to="/settings"
          :class="[
            'nav-item',
            { 'nav-item-active': $route.path === '/settings' },
            sidebarCollapsed ? 'lg:justify-center lg:px-2' : '',
          ]"
          :title="sidebarCollapsed ? 'Settings' : ''"
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
              d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
            />
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
            />
          </svg>
          <span v-if="!sidebarCollapsed" class="transition-opacity duration-300"
            >Settings</span
          >
        </router-link>

        <router-link
          to="/users"
          :class="[
            'nav-item',
            { 'nav-item-active': $route.path === '/users' },
            sidebarCollapsed ? 'lg:justify-center lg:px-2' : '',
          ]"
          :title="sidebarCollapsed ? 'User Management' : ''"
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
              d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"
            />
          </svg>
          <span v-if="!sidebarCollapsed" class="transition-opacity duration-300"
            >User Management</span
          >
        </router-link>

        <!-- Employee Management Dropdown -->
        <div class="space-y-1">
          <button
            @click="toggleEmployeeDropdown"
            :class="[
              'nav-item w-full',
              {
                'nav-item-active':
                  $route.path === '/departments' ||
                  $route.path === '/employees',
              },
              sidebarCollapsed ? 'lg:justify-center lg:px-2' : '',
            ]"
            :title="sidebarCollapsed ? 'Employee Management' : ''"
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
                d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
              />
            </svg>
            <span
              v-if="!sidebarCollapsed"
              class="flex-1 transition-opacity duration-300"
              >Employee</span
            >
            <svg
              v-if="!sidebarCollapsed"
              :class="[
                'w-4 h-4 transition-transform duration-200',
                { 'rotate-180': employeeDropdownOpen },
              ]"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </button>

          <!-- Dropdown Items -->
          <div
            v-show="employeeDropdownOpen || sidebarCollapsed"
            :class="[
              'space-y-1',
              sidebarCollapsed
                ? 'lg:absolute lg:left-full lg:top-0 lg:w-48 lg:mt-0 lg:bg-white lg:rounded-lg lg:shadow-lg lg:border lg:border-gray-200'
                : 'pl-4',
            ]"
          >
            <router-link
              to="/departments"
              :class="[
                'nav-item',
                { 'nav-item-active': $route.path === '/departments' },
                sidebarCollapsed ? 'lg:rounded-none' : '',
              ]"
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
                  d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
                />
              </svg>
              <span class="transition-opacity duration-300"
                >Manajemen Departemen</span
              >
            </router-link>

            <router-link
              to="/employees"
              :class="[
                'nav-item',
                { 'nav-item-active': $route.path === '/employees' },
                sidebarCollapsed ? 'lg:rounded-none' : '',
              ]"
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
                  d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                />
              </svg>
              <span class="transition-opacity duration-300"
                >Manajemen Karyawan</span
              >
            </router-link>
          </div>
        </div>

        <!-- Order Management Dropdown -->
        <div class="space-y-1">
          <button
            @click="toggleOrderDropdown"
            :class="[
              'nav-item w-full',
              {
                'nav-item-active':
                  $route.path === '/buyers' ||
                  $route.path === '/styles' ||
                  $route.path === '/line-schedules',
              },
              sidebarCollapsed ? 'lg:justify-center lg:px-2' : '',
            ]"
            :title="sidebarCollapsed ? 'Order Management' : ''"
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
                d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
              />
            </svg>
            <span
              v-if="!sidebarCollapsed"
              class="flex-1 transition-opacity duration-300"
              >Order</span
            >
            <svg
              v-if="!sidebarCollapsed"
              :class="[
                'w-4 h-4 transition-transform duration-200',
                { 'rotate-180': orderDropdownOpen },
              ]"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </button>

          <!-- Order Dropdown Items -->
          <div
            v-show="orderDropdownOpen || sidebarCollapsed"
            :class="[
              'space-y-1',
              sidebarCollapsed
                ? 'lg:absolute lg:left-full lg:top-0 lg:w-48 lg:mt-0 lg:bg-white lg:rounded-lg lg:shadow-lg lg:border lg:border-gray-200'
                : 'pl-4',
            ]"
          >
            <router-link
              to="/buyers"
              :class="[
                'nav-item',
                { 'nav-item-active': $route.path === '/buyers' },
                sidebarCollapsed ? 'lg:rounded-none' : '',
              ]"
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
                  d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                />
              </svg>
              <span class="transition-opacity duration-300"
                >Manajemen Buyer</span
              >
            </router-link>

            <router-link
              to="/styles"
              :class="[
                'nav-item',
                { 'nav-item-active': $route.path === '/styles' },
                sidebarCollapsed ? 'lg:rounded-none' : '',
              ]"
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
                  d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01"
                />
              </svg>
              <span class="transition-opacity duration-300"
                >Manajemen Style</span
              >
            </router-link>

            <router-link
              to="/line-schedules"
              :class="[
                'nav-item',
                { 'nav-item-active': $route.path === '/line-schedules' },
                sidebarCollapsed ? 'lg:rounded-none' : '',
              ]"
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
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                />
              </svg>
              <span class="transition-opacity duration-300">Line Schedule</span>
            </router-link>
          </div>
        </div>

        <!-- Schedule Management Dropdown -->
        <div class="space-y-1">
          <button
            @click="toggleScheduleDropdown"
            :class="[
              'nav-item w-full',
              {
                'nav-item-active':
                  $route.path === '/line-schedules' ||
                  $route.path === '/holidays',
              },
              sidebarCollapsed ? 'lg:justify-center lg:px-2' : '',
            ]"
            :title="sidebarCollapsed ? 'Schedule Management' : ''"
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
                d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
              />
            </svg>
            <span
              v-if="!sidebarCollapsed"
              class="flex-1 transition-opacity duration-300"
              >Schedule</span
            >
            <svg
              v-if="!sidebarCollapsed"
              :class="[
                'w-4 h-4 transition-transform duration-200',
                { 'rotate-180': scheduleDropdownOpen },
              ]"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </button>

          <!-- Schedule Dropdown Items -->
          <div
            v-show="scheduleDropdownOpen || sidebarCollapsed"
            :class="[
              'space-y-1',
              sidebarCollapsed
                ? 'lg:absolute lg:left-full lg:top-0 lg:w-48 lg:mt-0 lg:bg-white lg:rounded-lg lg:shadow-lg lg:border lg:border-gray-200'
                : 'pl-4',
            ]"
          >
            <router-link
              to="/line-schedules"
              :class="[
                'nav-item',
                { 'nav-item-active': $route.path === '/line-schedules' },
                sidebarCollapsed ? 'lg:rounded-none' : '',
              ]"
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
                  d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
                />
              </svg>
              <span class="transition-opacity duration-300">Line Schedule</span>
            </router-link>

            <router-link
              to="/holidays"
              :class="[
                'nav-item',
                { 'nav-item-active': $route.path === '/holidays' },
                sidebarCollapsed ? 'lg:rounded-none' : '',
              ]"
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
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                />
              </svg>
              <span class="transition-opacity duration-300"
                >Manajemen Hari Libur</span
              >
            </router-link>
          </div>
        </div>

        <!-- Divider -->
        <div class="border-t border-gray-200 my-4"></div>

        <!-- User Section -->
        <div v-if="!sidebarCollapsed || !isLargeScreen" class="px-3 py-2">
          <div
            class="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2"
          >
            User
          </div>
          <div class="flex items-center space-x-3 p-3 rounded-lg bg-gray-50">
            <div
              class="w-10 h-10 rounded-full overflow-hidden border-2 border-gray-200"
            >
              <img
                v-if="user?.profile_picture"
                :src="getProfileImageUrl(user.profile_picture)"
                :alt="user?.name || 'User'"
                class="w-full h-full object-cover"
              />
              <div
                v-else
                class="w-full h-full bg-primary-100 flex items-center justify-center"
              >
                <span class="text-sm font-medium text-primary-600">
                  {{ user?.name?.charAt(0).toUpperCase() || "U" }}
                </span>
              </div>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-900 truncate">
                {{ user?.name || "User" }}
              </p>
              <p class="text-xs text-gray-500 truncate">
                {{ user?.nik || "NIK" }} -
                {{ user?.email || "user@example.com" }}
              </p>
            </div>
          </div>
        </div>

        <!-- Collapsed User Avatar (Desktop only) -->
        <div
          v-if="sidebarCollapsed && isLargeScreen"
          class="hidden lg:flex px-3 py-2 justify-center"
        >
          <div
            class="w-8 h-8 rounded-full overflow-hidden border-2 border-gray-200"
            :title="user?.name || 'User'"
          >
            <img
              v-if="user?.profile_picture"
              :src="getProfileImageUrl(user.profile_picture)"
              :alt="user?.name || 'User'"
              class="w-full h-full object-cover"
            />
            <div
              v-else
              class="w-full h-full bg-primary-100 flex items-center justify-center"
            >
              <span class="text-sm font-medium text-primary-600">
                {{ user?.name?.charAt(0).toUpperCase() || "U" }}
              </span>
            </div>
          </div>
        </div>

        <!-- Logout Button -->
        <button
          @click="handleLogout"
          :class="[
            'w-full nav-item text-red-600 hover:bg-red-50 hover:text-red-700 mt-4',
            sidebarCollapsed ? 'lg:justify-center lg:px-2' : '',
          ]"
          :title="sidebarCollapsed ? 'Logout' : ''"
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
              d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
            />
          </svg>
          <span
            v-if="!sidebarCollapsed || !isLargeScreen"
            class="transition-opacity duration-300"
            >Logout</span
          >
        </button>
      </nav>
    </aside>

    <!-- Main Content Area -->
    <div
      :class="[
        'pt-16 lg:pt-0 transition-all duration-300 min-h-screen flex flex-col',
        sidebarCollapsed ? 'lg:pl-16' : 'lg:pl-64',
      ]"
    >
      <!-- Main Content -->
      <main class="flex-1">
        <div class="px-4 sm:px-6 lg:px-8 py-6 lg:py-8">
          <slot></slot>
        </div>
      </main>

      <!-- Footer -->
      <footer class="bg-white border-t border-gray-200 mt-auto">
        <div class="px-4 sm:px-6 lg:px-8 py-4">
          <p class="text-center text-sm text-gray-500">
            &copy; 2025 CORS Application - Vue.js & Go Backend
          </p>
        </div>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";

// Router and store
const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();

// Props
const props = defineProps({
  user: {
    type: Object,
    default: null,
  },
});

// Emits
const emit = defineEmits(["logout"]);

// Reactive data
const sidebarOpen = ref(false);
const sidebarCollapsed = ref(false);
const isLargeScreen = ref(window?.innerWidth >= 1024);
const employeeDropdownOpen = ref(false);
const orderDropdownOpen = ref(false);
const scheduleDropdownOpen = ref(false);

// Watch window resize
const handleResize = () => {
  isLargeScreen.value = window.innerWidth >= 1024;
};

onMounted(() => {
  window.addEventListener("resize", handleResize);
});

onUnmounted(() => {
  window.removeEventListener("resize", handleResize);
});

// Computed properties
const pageTitle = computed(() => {
  const routeNames = {
    "/dashboard": "Dashboard",
    "/profile": "Profile",
    "/settings": "Settings",
    "/holidays": "Manajemen Hari Libur",
  };
  return routeNames[route.path] || "Dashboard";
});

// Methods
const getProfileImageUrl = (profilePicture) => {
  if (!profilePicture) return null;
  if (profilePicture.startsWith("http")) {
    return profilePicture;
  }
  const baseUrl = import.meta.env.VITE_API_BASE_URL || "http://localhost:8081";
  return `${baseUrl}${profilePicture}`;
};

const toggleEmployeeDropdown = () => {
  employeeDropdownOpen.value = !employeeDropdownOpen.value;
};

const toggleOrderDropdown = () => {
  orderDropdownOpen.value = !orderDropdownOpen.value;
};

const toggleScheduleDropdown = () => {
  scheduleDropdownOpen.value = !scheduleDropdownOpen.value;
};

const handleLogout = () => {
  authStore.logout();
  router.push("/login");
  emit("logout");
};
</script>

<style scoped>
.nav-item {
  @apply flex items-center space-x-3 px-3 py-2 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors duration-200;
}

.nav-item-active {
  @apply bg-primary-100 text-primary-700 hover:bg-primary-200;
}

.nav-item svg {
  @apply flex-shrink-0;
}

/* Sidebar animation improvements */
@media (max-width: 1023px) {
  .sidebar-enter-active,
  .sidebar-leave-active {
    transition: transform 0.3s ease-in-out;
  }

  .sidebar-enter-from,
  .sidebar-leave-to {
    transform: translateX(-100%);
  }
}
</style>
