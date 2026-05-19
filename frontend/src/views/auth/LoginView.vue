<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-950 py-12 px-4">
    <div class="max-w-md w-full space-y-6">
      <!-- Header -->
      <div class="text-center">
        <div class="mx-auto h-14 w-14 flex items-center justify-center rounded-2xl bg-indigo-900 mb-4">
          <span class="text-3xl">💰</span>
        </div>
        <h1 class="text-3xl font-bold text-white tracking-tight">
          FinanceTracker
        </h1>
        <p class="mt-2 text-gray-400 text-sm">
          Sign in to manage your finances
        </p>
      </div>

      <!-- Login Form -->
      <form
        class="mt-6 space-y-4 bg-gray-800 border border-gray-700 rounded-2xl p-8 shadow-2xl"
        @submit.prevent="handleSubmit"
      >
        <div class="space-y-4">
          <div>
            <label for="email" class="block text-sm font-medium text-gray-300 mb-1">
              Email address
            </label>
            <input
              id="email"
              v-model="form.email"
              name="email"
              type="email"
              autocomplete="email"
              required
              class="block w-full px-4 py-3 bg-gray-900 border border-gray-600 text-white placeholder-gray-500 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition sm:text-sm"
              placeholder="you@example.com"
              :disabled="isLoading"
            />
          </div>
          <div>
            <label for="password" class="block text-sm font-medium text-gray-300 mb-1">
              Password
            </label>
            <input
              id="password"
              v-model="form.password"
              name="password"
              type="password"
              autocomplete="current-password"
              required
              class="block w-full px-4 py-3 bg-gray-900 border border-gray-600 text-white placeholder-gray-500 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition sm:text-sm"
              placeholder="••••••••"
              :disabled="isLoading"
            />
          </div>
        </div>

        <button
          type="submit"
          class="w-full flex justify-center py-3 px-4 border border-transparent text-sm font-semibold rounded-xl text-white bg-indigo-600 hover:bg-indigo-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          :disabled="isLoading"
        >
          <span v-if="isLoading" class="h-5 w-5 border-2 border-white border-t-transparent rounded-full animate-spin mr-2 inline-block"></span>
          {{ isLoading ? 'Signing in...' : 'Sign in' }}
        </button>

        <!-- Demo Account -->
        <div class="mt-4 p-4 bg-indigo-950 rounded-xl border border-indigo-800">
          <div class="flex items-start justify-between">
            <div>
              <h3 class="text-sm font-semibold text-indigo-300 mb-1">🚀 Demo Account</h3>
              <p class="text-xs text-indigo-400">demo@financetracker.com / demo123456</p>
            </div>
            <button
              type="button"
              @click="fillDemoCredentials"
              class="ml-4 text-xs font-medium text-indigo-300 bg-indigo-800 hover:bg-indigo-700 px-3 py-1.5 rounded-lg transition-colors whitespace-nowrap"
            >
              Use Demo
            </button>
          </div>
        </div>

        <p class="text-center text-sm text-gray-400">
          Don't have an account?
          <router-link to="/register" class="font-medium text-indigo-400 hover:text-indigo-300 ml-1">
            Create account
          </router-link>
        </p>
      </form>
    </div>
  </div>
</template>


<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useAuthStore } from '@/stores/auth'
import type { LoginForm } from '@/types'

const authStore = useAuthStore()

const form = reactive<LoginForm>({
  email: '',
  password: ''
})

const rememberMe = ref(false)
const isLoading = ref(false)

const handleSubmit = async () => {
  if (isLoading.value) return

  try {
    isLoading.value = true
    await authStore.login(form)
  } catch (error) {
    console.error('Login failed:', error)
  } finally {
    isLoading.value = false
  }
}

const fillDemoCredentials = () => {
  form.email = 'demo@financetracker.com'
  form.password = 'demo123456'
}
</script>