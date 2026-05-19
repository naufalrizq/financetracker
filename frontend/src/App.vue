<template>
  <div id="app" class="min-h-screen bg-gray-50 dark:bg-gray-900 dark:text-gray-100">
    <!-- Loading overlay -->
    <div
      v-if="isInitializing"
      class="fixed inset-0 bg-white z-50 flex items-center justify-center"
    >
      <div class="text-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
        <p class="text-gray-600">Loading FinanceTracker...</p>
      </div>
    </div>

    <!-- Main app content -->
    <div v-else>
      <!-- Navigation -->
      <AppNavigation v-if="isAuthenticated" />
      
      <!-- Main content -->
      <main :class="{ 'ml-64 p-6': isAuthenticated }">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterView } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import AppNavigation from '@/components/layout/AppNavigation.vue'

const authStore = useAuthStore()
const isInitializing = ref(true)

const isAuthenticated = computed(() => authStore.isAuthenticated)

onMounted(async () => {
  try {
    // Try to restore authentication state
    await authStore.initializeAuth()
  } catch (error) {
    console.error('Failed to initialize auth:', error)
  } finally {
    isInitializing.value = false
  }
})
</script>

<style>
/* Global styles */
#app {
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* Custom scrollbar */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f5f9;
}

::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.3s ease;
}

.slide-enter-from {
  transform: translateX(-100%);
}

.slide-leave-to {
  transform: translateX(100%);
}
</style>