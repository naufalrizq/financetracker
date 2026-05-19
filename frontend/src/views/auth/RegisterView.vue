<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <!-- Header -->
      <div>
        <div class="mx-auto h-12 w-12 flex items-center justify-center rounded-full bg-primary-100">
          <span class="text-2xl">💰</span>
        </div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Create your account
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Or
          <router-link to="/login" class="font-medium text-primary-600 hover:text-primary-500">
            sign in to your existing account
          </router-link>
        </p>
      </div>

      <!-- Register Form -->
      <form class="mt-8 space-y-6" @submit.prevent="handleSubmit">
        <div class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label for="first_name" class="form-label">First Name</label>
              <input
                id="first_name"
                v-model="form.first_name"
                name="first_name"
                type="text"
                required
                class="form-input"
                placeholder="John"
                :disabled="isLoading"
              />
            </div>
            <div>
              <label for="last_name" class="form-label">Last Name</label>
              <input
                id="last_name"
                v-model="form.last_name"
                name="last_name"
                type="text"
                required
                class="form-input"
                placeholder="Doe"
                :disabled="isLoading"
              />
            </div>
          </div>

          <div>
            <label for="email" class="form-label">Email Address</label>
            <input
              id="email"
              v-model="form.email"
              name="email"
              type="email"
              autocomplete="email"
              required
              class="form-input"
              placeholder="john@example.com"
              :disabled="isLoading"
            />
          </div>

          <div>
            <label for="password" class="form-label">Password</label>
            <input
              id="password"
              v-model="form.password"
              name="password"
              type="password"
              autocomplete="new-password"
              required
              class="form-input"
              placeholder="Enter your password"
              :disabled="isLoading"
            />
            <p class="mt-1 text-sm text-gray-500">
              Password must be at least 8 characters long
            </p>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label for="currency" class="form-label">Currency</label>
              <select
                id="currency"
                v-model="form.currency"
                name="currency"
                class="form-select"
                :disabled="isLoading"
              >
                <option value="USD">USD - US Dollar</option>
                <option value="EUR">EUR - Euro</option>
                <option value="GBP">GBP - British Pound</option>
                <option value="JPY">JPY - Japanese Yen</option>
                <option value="CAD">CAD - Canadian Dollar</option>
                <option value="AUD">AUD - Australian Dollar</option>
                <option value="CHF">CHF - Swiss Franc</option>
                <option value="CNY">CNY - Chinese Yuan</option>
                <option value="INR">INR - Indian Rupee</option>
              </select>
            </div>
            <div>
              <label for="timezone" class="form-label">Timezone</label>
              <select
                id="timezone"
                v-model="form.timezone"
                name="timezone"
                class="form-select"
                :disabled="isLoading"
              >
                <option value="UTC">UTC</option>
                <option value="America/New_York">Eastern Time</option>
                <option value="America/Chicago">Central Time</option>
                <option value="America/Denver">Mountain Time</option>
                <option value="America/Los_Angeles">Pacific Time</option>
                <option value="Europe/London">London</option>
                <option value="Europe/Paris">Paris</option>
                <option value="Asia/Tokyo">Tokyo</option>
                <option value="Asia/Shanghai">Shanghai</option>
                <option value="Asia/Kolkata">Mumbai</option>
              </select>
            </div>
          </div>
        </div>

        <div class="flex items-center">
          <input
            id="agree-terms"
            v-model="agreeTerms"
            name="agree-terms"
            type="checkbox"
            class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
            required
          />
          <label for="agree-terms" class="ml-2 block text-sm text-gray-900">
            I agree to the
            <a href="#" class="text-primary-600 hover:text-primary-500">Terms of Service</a>
            and
            <a href="#" class="text-primary-600 hover:text-primary-500">Privacy Policy</a>
          </label>
        </div>

        <div>
          <button
            type="submit"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
            :disabled="isLoading || !agreeTerms"
          >
            <span class="absolute left-0 inset-y-0 flex items-center pl-3">
              <svg
                v-if="!isLoading"
                class="h-5 w-5 text-primary-500 group-hover:text-primary-400"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M10 1a4.5 4.5 0 00-4.5 4.5V9H5a2 2 0 00-2 2v6a2 2 0 002 2h10a2 2 0 002-2v-6a2 2 0 00-2-2h-.5V5.5A4.5 4.5 0 0010 1zm3 8V5.5a3 3 0 10-6 0V9h6z"
                  clip-rule="evenodd"
                />
              </svg>
              <div
                v-else
                class="h-5 w-5 border-2 border-white border-t-transparent rounded-full animate-spin"
              ></div>
            </span>
            {{ isLoading ? 'Creating account...' : 'Create account' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useAuthStore } from '@/stores/auth'
import type { RegisterForm } from '@/types'

const authStore = useAuthStore()

const form = reactive<RegisterForm>({
  email: '',
  password: '',
  first_name: '',
  last_name: '',
  currency: 'IDR',
  timezone: 'Asia/Jakarta'
})

const agreeTerms = ref(false)
const isLoading = ref(false)

const handleSubmit = async () => {
  if (isLoading.value || !agreeTerms.value) return

  try {
    isLoading.value = true
    await authStore.register(form)
  } catch (error) {
    console.error('Registration failed:', error)
  } finally {
    isLoading.value = false
  }
}
</script>