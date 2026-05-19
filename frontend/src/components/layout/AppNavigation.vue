<template>
  <div class="fixed inset-y-0 left-0 z-50 w-64 bg-white dark:bg-gray-800 dark:border-r dark:border-gray-700 shadow-lg">
    <!-- Sidebar Header -->
    <div class="flex items-center justify-between h-16 px-6 border-b border-gray-200 dark:border-gray-700">
      <div class="flex items-center">
        <span class="text-2xl">💰</span>
        <span class="ml-2 text-xl font-bold text-gray-900 dark:text-white">FinanceTracker</span>
      </div>
    </div>

    <!-- Navigation Menu -->
    <nav class="flex-1 px-4 py-6 space-y-2">
      <router-link
        v-for="item in navigationItems"
        :key="item.name"
        :to="item.to"
        :class="[
          'nav-link',
          $route.name === item.name ? 'nav-link-active' : 'nav-link-inactive'
        ]"
      >
        <span class="text-xl mr-3">{{ item.icon }}</span>
        {{ item.label }}
      </router-link>
    </nav>

    <!-- User Menu -->
    <div class="border-t border-gray-200 dark:border-gray-700 p-4">
      <div class="flex items-center mb-4">
        <div class="w-10 h-10 bg-primary-100 dark:bg-primary-900 rounded-full flex items-center justify-center">
          <span class="text-primary-600 dark:text-primary-300 font-medium">
            {{ userInitials }}
          </span>
        </div>
        <div class="ml-3">
          <p class="text-sm font-medium text-gray-900 dark:text-white">{{ authStore.userFullName }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400">{{ authStore.user?.email }}</p>
        </div>
      </div>

      <div class="space-y-1">
        <router-link
          to="/profile"
          class="nav-link nav-link-inactive text-sm"
        >
          <span class="text-lg mr-3">⚙️</span>
          Settings
        </router-link>
        
        <button
          @click="handleLogout"
          class="w-full nav-link nav-link-inactive text-sm text-left"
        >
          <span class="text-lg mr-3">🚪</span>
          Sign out
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const navigationItems = [
  { name: 'dashboard', to: '/dashboard', label: 'Dashboard', icon: '📊' },
  { name: 'transactions', to: '/transactions', label: 'Transactions', icon: '💳' },
  { name: 'accounts', to: '/accounts', label: 'Accounts', icon: '🏦' },
  { name: 'categories', to: '/categories', label: 'Categories', icon: '🏷️' },
  { name: 'budgets', to: '/budgets', label: 'Budgets', icon: '📈' },
  { name: 'goals', to: '/goals', label: 'Goals', icon: '🎯' },
  { name: 'reports', to: '/reports', label: 'Reports', icon: '📋' },
]

const userInitials = computed(() => {
  const user = authStore.user
  if (!user) return 'U'
  return (user.first_name[0] + user.last_name[0]).toUpperCase()
})

const handleLogout = async () => {
  await authStore.logout()
}
</script>