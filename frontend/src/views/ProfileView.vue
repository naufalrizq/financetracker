<template>
  <div class="space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold text-gray-900">Profile Settings</h1>
      <p class="text-gray-600">Manage your account information and preferences</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Profile Information -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Personal Information -->
        <div class="bg-white rounded-lg shadow-sm border">
          <div class="p-4 border-b">
            <h3 class="text-lg font-semibold">Personal Information</h3>
          </div>
          <div class="p-4">
            <form @submit.prevent="updateProfile" class="space-y-4">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">First Name</label>
                  <input
                    v-model="profileForm.firstName"
                    type="text"
                    required
                    class="w-full border border-gray-300 rounded-md px-3 py-2"
                    placeholder="Enter your first name"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Last Name</label>
                  <input
                    v-model="profileForm.lastName"
                    type="text"
                    required
                    class="w-full border border-gray-300 rounded-md px-3 py-2"
                    placeholder="Enter your last name"
                  />
                </div>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Email Address</label>
                <input
                  v-model="profileForm.email"
                  type="email"
                  required
                  class="w-full border border-gray-300 rounded-md px-3 py-2"
                  placeholder="Enter your email"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Phone Number</label>
                <input
                  v-model="profileForm.phone"
                  type="tel"
                  class="w-full border border-gray-300 rounded-md px-3 py-2"
                  placeholder="Enter your phone number"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Date of Birth</label>
                <input
                  v-model="profileForm.dateOfBirth"
                  type="date"
                  class="w-full border border-gray-300 rounded-md px-3 py-2"
                />
              </div>
              
              <div class="flex justify-end">
                <button
                  type="submit"
                  :disabled="updating"
                  class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 disabled:opacity-50"
                >
                  {{ updating ? 'Updating...' : 'Update Profile' }}
                </button>
              </div>
            </form>
          </div>
        </div>

        <!-- Change Password -->
        <div class="bg-white rounded-lg shadow-sm border">
          <div class="p-4 border-b">
            <h3 class="text-lg font-semibold">Change Password</h3>
          </div>
          <div class="p-4">
            <form @submit.prevent="changePassword" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Current Password</label>
                <input
                  v-model="passwordForm.currentPassword"
                  type="password"
                  required
                  class="w-full border border-gray-300 rounded-md px-3 py-2"
                  placeholder="Enter current password"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">New Password</label>
                <input
                  v-model="passwordForm.newPassword"
                  type="password"
                  required
                  minlength="6"
                  class="w-full border border-gray-300 rounded-md px-3 py-2"
                  placeholder="Enter new password"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Confirm New Password</label>
                <input
                  v-model="passwordForm.confirmPassword"
                  type="password"
                  required
                  minlength="6"
                  class="w-full border border-gray-300 rounded-md px-3 py-2"
                  placeholder="Confirm new password"
                />
              </div>
              
              <div v-if="passwordError" class="text-red-600 text-sm">
                {{ passwordError }}
              </div>
              
              <div class="flex justify-end">
                <button
                  type="submit"
                  :disabled="changingPassword"
                  class="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 disabled:opacity-50"
                >
                  {{ changingPassword ? 'Changing...' : 'Change Password' }}
                </button>
              </div>
            </form>
          </div>
        </div>

        <!-- Preferences -->
        <div class="bg-white rounded-lg shadow-sm border">
          <div class="p-4 border-b">
            <h3 class="text-lg font-semibold">Preferences</h3>
          </div>
          <div class="p-4 space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Currency</label>
              <select v-model="preferencesForm.currency" class="w-full border border-gray-300 rounded-md px-3 py-2">
                <option value="USD">USD - US Dollar</option>
                <option value="EUR">EUR - Euro</option>
                <option value="GBP">GBP - British Pound</option>
                <option value="JPY">JPY - Japanese Yen</option>
                <option value="CAD">CAD - Canadian Dollar</option>
                <option value="AUD">AUD - Australian Dollar</option>
              </select>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Date Format</label>
              <select v-model="preferencesForm.dateFormat" class="w-full border border-gray-300 rounded-md px-3 py-2">
                <option value="MM/DD/YYYY">MM/DD/YYYY</option>
                <option value="DD/MM/YYYY">DD/MM/YYYY</option>
                <option value="YYYY-MM-DD">YYYY-MM-DD</option>
              </select>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Theme</label>
              <select v-model="preferencesForm.theme" class="w-full border border-gray-300 rounded-md px-3 py-2">
                <option value="light">Light</option>
                <option value="dark">Dark</option>
                <option value="auto">Auto (System)</option>
              </select>
            </div>
            
            <div class="space-y-3">
              <h4 class="text-sm font-medium text-gray-700">Notifications</h4>
              <div class="space-y-2">
                <label class="flex items-center">
                  <input
                    v-model="preferencesForm.emailNotifications"
                    type="checkbox"
                    class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                  <span class="ml-2 text-sm text-gray-700">Email notifications</span>
                </label>
                <label class="flex items-center">
                  <input
                    v-model="preferencesForm.budgetAlerts"
                    type="checkbox"
                    class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                  <span class="ml-2 text-sm text-gray-700">Budget alerts</span>
                </label>
                <label class="flex items-center">
                  <input
                    v-model="preferencesForm.goalReminders"
                    type="checkbox"
                    class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                  <span class="ml-2 text-sm text-gray-700">Goal reminders</span>
                </label>
              </div>
            </div>
            
            <div class="flex justify-end">
              <button
                @click="updatePreferences"
                :disabled="updatingPreferences"
                class="bg-green-600 text-white px-4 py-2 rounded-md hover:bg-green-700 disabled:opacity-50"
              >
                {{ updatingPreferences ? 'Saving...' : 'Save Preferences' }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Sidebar -->
      <div class="space-y-6">
        <!-- Profile Picture -->
        <div class="bg-white rounded-lg shadow-sm border">
          <div class="p-4 border-b">
            <h3 class="text-lg font-semibold">Profile Picture</h3>
          </div>
          <div class="p-4 text-center">
            <div class="w-24 h-24 rounded-full bg-gray-200 mx-auto mb-4 flex items-center justify-center">
              <UserIcon class="w-12 h-12 text-gray-400" />
            </div>
            <button class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 text-sm">
              Upload Photo
            </button>
            <p class="text-xs text-gray-500 mt-2">JPG, PNG up to 2MB</p>
          </div>
        </div>

        <!-- Account Stats -->
        <div class="bg-white rounded-lg shadow-sm border">
          <div class="p-4 border-b">
            <h3 class="text-lg font-semibold">Account Statistics</h3>
          </div>
          <div class="p-4 space-y-3">
            <div class="flex justify-between">
              <span class="text-sm text-gray-600">Member since</span>
              <span class="text-sm font-medium">{{ formatDate(user?.created_at) }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-sm text-gray-600">Total transactions</span>
              <span class="text-sm font-medium">{{ accountStats.totalTransactions }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-sm text-gray-600">Active accounts</span>
              <span class="text-sm font-medium">{{ accountStats.totalAccounts }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-sm text-gray-600">Active budgets</span>
              <span class="text-sm font-medium">{{ accountStats.totalBudgets }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-sm text-gray-600">Financial goals</span>
              <span class="text-sm font-medium">{{ accountStats.totalGoals }}</span>
            </div>
          </div>
        </div>

        <!-- Danger Zone -->
        <div class="bg-white rounded-lg shadow-sm border border-red-200">
          <div class="p-4 border-b border-red-200">
            <h3 class="text-lg font-semibold text-red-600">Danger Zone</h3>
          </div>
          <div class="p-4 space-y-3">
            <div>
              <h4 class="text-sm font-medium text-gray-900 mb-1">Export Data</h4>
              <p class="text-xs text-gray-600 mb-2">Download all your financial data</p>
              <button class="w-full bg-gray-600 text-white px-3 py-2 rounded-md hover:bg-gray-700 text-sm">
                Export All Data
              </button>
            </div>
            <div>
              <h4 class="text-sm font-medium text-gray-900 mb-1">Delete Account</h4>
              <p class="text-xs text-gray-600 mb-2">Permanently delete your account and all data</p>
              <button
                @click="confirmDeleteAccount"
                class="w-full bg-red-600 text-white px-3 py-2 rounded-md hover:bg-red-700 text-sm"
              >
                Delete Account
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { UserIcon } from '@heroicons/vue/24/outline'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/services/api'

// Store
const authStore = useAuthStore()

// Reactive data
const updating = ref(false)
const changingPassword = ref(false)
const updatingPreferences = ref(false)
const passwordError = ref('')
const user = ref(authStore.user)
const accountStats = ref({
  totalTransactions: 0,
  totalAccounts: 0,
  totalBudgets: 0,
  totalGoals: 0
})

// Forms
const profileForm = reactive({
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
  dateOfBirth: ''
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const preferencesForm = reactive({
  currency: 'IDR',
  dateFormat: 'DD/MM/YYYY',
  theme: 'light',
  emailNotifications: true,
  budgetAlerts: true,
  goalReminders: true
})

// Methods
const loadUserData = async () => {
  try {
    const response = await api.get('/user/profile')
    const userData = response.data
    
    // Populate profile form
    profileForm.firstName = userData.first_name || ''
    profileForm.lastName = userData.last_name || ''
    profileForm.email = userData.email || ''
    profileForm.phone = userData.phone || ''
    profileForm.dateOfBirth = userData.date_of_birth ? userData.date_of_birth.split('T')[0] : ''
    
    // Load preferences
    if (userData.preferences) {
      Object.assign(preferencesForm, userData.preferences)
    }
  } catch (error) {
    console.error('Failed to load user data:', error)
  }
}

const loadAccountStats = async () => {
  try {
    const [transactions, accounts, budgets, goals] = await Promise.all([
      api.get('/transactions'),
      api.get('/accounts'),
      api.get('/budgets'),
      api.get('/goals')
    ])
    
    accountStats.value = {
      totalTransactions: transactions.data.length,
      totalAccounts: accounts.data.length,
      totalBudgets: budgets.data.length,
      totalGoals: goals.data.length
    }
  } catch (error) {
    console.error('Failed to load account stats:', error)
  }
}

const updateProfile = async () => {
  try {
    updating.value = true
    
    const data = {
      first_name: profileForm.firstName,
      last_name: profileForm.lastName,
      email: profileForm.email,
      phone: profileForm.phone,
      date_of_birth: profileForm.dateOfBirth
    }
    
    await api.put('/user/profile', data)
    
    // Update auth store
    authStore.updateUser({
      ...authStore.user,
      ...data
    })
    
    alert('Profile updated successfully!')
  } catch (error) {
    console.error('Failed to update profile:', error)
    alert('Failed to update profile. Please try again.')
  } finally {
    updating.value = false
  }
}

const changePassword = async () => {
  passwordError.value = ''
  
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    passwordError.value = 'New passwords do not match'
    return
  }
  
  if (passwordForm.newPassword.length < 6) {
    passwordError.value = 'Password must be at least 6 characters long'
    return
  }
  
  try {
    changingPassword.value = true
    
    await api.put('/user/change-password', {
      current_password: passwordForm.currentPassword,
      new_password: passwordForm.newPassword
    })
    
    // Clear form
    passwordForm.currentPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
    
    alert('Password changed successfully!')
  } catch (error: any) {
    console.error('Failed to change password:', error)
    passwordError.value = error.response?.data?.message || 'Failed to change password'
  } finally {
    changingPassword.value = false
  }
}

const updatePreferences = async () => {
  try {
    updatingPreferences.value = true
    
    await api.put('/user/preferences', preferencesForm)
    
    alert('Preferences saved successfully!')
  } catch (error) {
    console.error('Failed to update preferences:', error)
    alert('Failed to save preferences. Please try again.')
  } finally {
    updatingPreferences.value = false
  }
}

const confirmDeleteAccount = () => {
  const confirmation = prompt(
    'This action cannot be undone. Type "DELETE" to confirm account deletion:'
  )
  
  if (confirmation === 'DELETE') {
    deleteAccount()
  }
}

const deleteAccount = async () => {
  try {
    await api.delete('/user/account')
    
    // Logout and redirect
    authStore.logout()
    alert('Account deleted successfully.')
  } catch (error) {
    console.error('Failed to delete account:', error)
    alert('Failed to delete account. Please try again.')
  }
}

const formatDate = (dateString?: string) => {
  if (!dateString) return 'N/A'
  return new Date(dateString).toLocaleDateString()
}

// Lifecycle
onMounted(() => {
  loadUserData()
  loadAccountStats()
})
</script>