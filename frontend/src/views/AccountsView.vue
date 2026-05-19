<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Accounts</h1>
        <p class="text-gray-600 dark:text-gray-400">Manage your bank accounts and wallets</p>
      </div>
      <button
        @click="showAddModal = true"
        class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 flex items-center gap-2"
      >
        <PlusIcon class="w-5 h-5" />
        Add Account
      </button>
    </div>

    <!-- Account Summary Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-gradient-to-r from-blue-500 to-blue-600 text-white p-6 rounded-lg">
        <h3 class="text-lg font-semibold mb-2">Total Balance</h3>
        <p class="text-3xl font-bold">{{ formatCurrency(totalBalance) }}</p>
      </div>
      <div class="bg-gradient-to-r from-green-500 to-green-600 text-white p-6 rounded-lg">
        <h3 class="text-lg font-semibold mb-2">Active Accounts</h3>
        <p class="text-3xl font-bold">{{ activeAccounts }}</p>
      </div>
      <div class="bg-gradient-to-r from-purple-500 to-purple-600 text-white p-6 rounded-lg">
        <h3 class="text-lg font-semibold mb-2">Account Types</h3>
        <p class="text-3xl font-bold">{{ uniqueAccountTypes }}</p>
      </div>
    </div>

    <!-- Accounts List -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700">
      <div class="p-4 border-b border-gray-200 dark:border-gray-700">
        <h2 class="text-lg font-semibold dark:text-white">Your Accounts</h2>
      </div>
      
      <div v-if="loading" class="p-8 text-center">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 dark:border-blue-400 mx-auto"></div>
        <p class="mt-2 text-gray-600 dark:text-gray-400">Loading accounts...</p>
      </div>

      <div v-else-if="accounts.length === 0" class="p-8 text-center text-gray-500 dark:text-gray-400">
        <BanknotesIcon class="w-12 h-12 mx-auto mb-4 text-gray-400 dark:text-gray-500" />
        <p class="text-lg font-medium mb-2 dark:text-gray-300">No accounts yet</p>
        <p class="mb-4">Add your first account to start tracking your finances</p>
        <button
          @click="showAddModal = true"
          class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
        >
          Add Account
        </button>
      </div>

      <div v-else class="divide-y divide-gray-200 dark:divide-gray-700">
        <div
          v-for="account in accounts"
          :key="account.id"
          class="p-4 hover:bg-gray-50 dark:hover:bg-gray-700 flex items-center justify-between"
        >
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-full flex items-center justify-center"
                 :class="getAccountTypeColor(account.type)">
              <component :is="getAccountIcon(account.type)" class="w-6 h-6" />
            </div>
            <div>
              <h3 class="font-semibold text-gray-900 dark:text-white">{{ account.name }}</h3>
              <p class="text-sm text-gray-600 dark:text-gray-400 capitalize">{{ account.type.replace('_', ' ') }}</p>
              <p class="text-xs text-gray-500 dark:text-gray-500">Created {{ formatDate(account.created_at) }}</p>
            </div>
          </div>
          <div class="text-right">
            <p class="text-xl font-bold text-gray-900 dark:text-white">{{ formatCurrency(account.balance) }}</p>
            <div class="flex gap-2 mt-2">
              <button
                @click="editAccount(account)"
                class="text-blue-600 hover:text-blue-800 text-sm px-2 py-1 rounded"
              >
                Edit
              </button>
              <button
                @click="deleteAccount(account.id)"
                class="text-red-600 hover:text-red-800 text-sm px-2 py-1 rounded"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Account Modal -->
    <div v-if="showAddModal || showEditModal" class="fixed inset-0 bg-black bg-opacity-50 dark:bg-opacity-70 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md mx-4">
        <h3 class="text-lg font-semibold mb-4 dark:text-white">
          {{ showEditModal ? 'Edit Account' : 'Add Account' }}
        </h3>
        
        <form @submit.prevent="submitAccount" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Account Name</label>
            <input
              v-model="accountForm.name"
              type="text"
              required
              class="w-full border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md px-3 py-2"
              placeholder="e.g., Main Checking, Savings, Cash Wallet"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Account Type</label>
            <select v-model="accountForm.type" required class="w-full border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md px-3 py-2">
              <option value="">Select Type</option>
              <option value="checking">Checking Account</option>
              <option value="savings">Savings Account</option>
              <option value="credit_card">Credit Card</option>
              <option value="cash">Cash</option>
              <option value="investment">Investment Account</option>
              <option value="loan">Loan Account</option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Initial Balance</label>
            <input
              v-model.number="accountForm.balance"
              type="number"
              step="0.01"
              required
              class="w-full border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md px-3 py-2"
              placeholder="0.00"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Description (Optional)</label>
            <textarea
              v-model="accountForm.description"
              rows="3"
              class="w-full border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md px-3 py-2"
              placeholder="Additional notes about this account"
            ></textarea>
          </div>
          
          <div class="flex gap-3 pt-4">
            <button
              type="submit"
              :disabled="submitting"
              class="flex-1 bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              {{ submitting ? 'Saving...' : (showEditModal ? 'Update' : 'Add') }}
            </button>
            <button
              type="button"
              @click="closeModal"
              class="flex-1 bg-gray-300 dark:bg-gray-600 text-gray-700 dark:text-white py-2 px-4 rounded-md hover:bg-gray-400 dark:hover:bg-gray-500"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { 
  PlusIcon, 
  BanknotesIcon,
  CreditCardIcon,
  BuildingLibraryIcon,
  WalletIcon,
  ChartBarIcon,
  DocumentTextIcon
} from '@heroicons/vue/24/outline'
import { api } from '@/services/api'
import { useCurrency } from '@/composables/useCurrency'
import type { Account } from '@/types'

const { formatCurrency } = useCurrency()

// Reactive data
const loading = ref(false)
const submitting = ref(false)
const showAddModal = ref(false)
const showEditModal = ref(false)
const accounts = ref<Account[]>([])

// Account form
const accountForm = reactive({
  id: null as number | null,
  name: '',
  type: '',
  balance: 0,
  description: ''
})

// Computed
const totalBalance = computed(() => {
  return accounts.value.reduce((sum, account) => sum + account.balance, 0)
})

const activeAccounts = computed(() => {
  return accounts.value.length
})

const uniqueAccountTypes = computed(() => {
  const types = new Set(accounts.value.map(account => account.type))
  return types.size
})

// Methods
const loadAccounts = async () => {
  try {
    loading.value = true
    const response = await api.get('/accounts')
    accounts.value = response.data
  } catch (error) {
    console.error('Failed to load accounts:', error)
  } finally {
    loading.value = false
  }
}

const submitAccount = async () => {
  try {
    submitting.value = true
    
    const data = {
      name: accountForm.name,
      type: accountForm.type,
      balance: accountForm.balance,
      description: accountForm.description
    }
    
    if (showEditModal.value && accountForm.id) {
      await api.put(`/accounts/${accountForm.id}`, data)
    } else {
      await api.post('/accounts', data)
    }
    
    await loadAccounts()
    closeModal()
  } catch (error) {
    console.error('Failed to save account:', error)
  } finally {
    submitting.value = false
  }
}

const editAccount = (account: Account) => {
  accountForm.id = account.id
  accountForm.name = account.name
  accountForm.type = account.type
  accountForm.balance = account.balance
  accountForm.description = account.description || ''
  showEditModal.value = true
}

const deleteAccount = async (id: number) => {
  if (!confirm('Are you sure you want to delete this account? This will also delete all associated transactions.')) return
  
  try {
    await api.delete(`/accounts/${id}`)
    await loadAccounts()
  } catch (error) {
    console.error('Failed to delete account:', error)
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  accountForm.id = null
  accountForm.name = ''
  accountForm.type = ''
  accountForm.balance = 0
  accountForm.description = ''
}

const getAccountIcon = (type: string) => {
  const icons = {
    checking: BuildingLibraryIcon,
    savings: BuildingLibraryIcon,
    credit_card: CreditCardIcon,
    cash: WalletIcon,
    investment: ChartBarIcon,
    loan: DocumentTextIcon
  }
  return icons[type as keyof typeof icons] || BanknotesIcon
}

const getAccountTypeColor = (type: string) => {
  const colors = {
    checking: 'bg-blue-100 text-blue-600',
    savings: 'bg-green-100 text-green-600',
    credit_card: 'bg-purple-100 text-purple-600',
    cash: 'bg-yellow-100 text-yellow-600',
    investment: 'bg-indigo-100 text-indigo-600',
    loan: 'bg-red-100 text-red-600'
  }
  return colors[type as keyof typeof colors] || 'bg-gray-100 text-gray-600'
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

// Lifecycle
onMounted(() => {
  loadAccounts()
})
</script>