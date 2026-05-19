<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Transactions</h1>
        <p class="text-gray-600">Manage your income and expenses</p>
      </div>
      <button
        @click="showAddModal = true"
        class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 flex items-center gap-2"
      >
        <PlusIcon class="w-5 h-5" />
        Add Transaction
      </button>
    </div>

    <!-- Filters -->
    <div class="bg-white p-4 rounded-lg shadow-sm border">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
          <select v-model="filters.type" class="w-full border border-gray-300 rounded-md px-3 py-2">
            <option value="">All Types</option>
            <option value="income">Income</option>
            <option value="expense">Expense</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Category</label>
          <select v-model="filters.categoryId" class="w-full border border-gray-300 rounded-md px-3 py-2">
            <option value="">All Categories</option>
            <option v-for="category in categories" :key="category.id" :value="category.id">
              {{ category.name }}
            </option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Account</label>
          <select v-model="filters.accountId" class="w-full border border-gray-300 rounded-md px-3 py-2">
            <option value="">All Accounts</option>
            <option v-for="account in accounts" :key="account.id" :value="account.id">
              {{ account.name }}
            </option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Date Range</label>
          <select v-model="filters.dateRange" class="w-full border border-gray-300 rounded-md px-3 py-2">
            <option value="7">Last 7 days</option>
            <option value="30">Last 30 days</option>
            <option value="90">Last 3 months</option>
            <option value="365">Last year</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Transactions List -->
    <div class="bg-white rounded-lg shadow-sm border">
      <div class="p-4 border-b">
        <h2 class="text-lg font-semibold">Recent Transactions</h2>
      </div>
      
      <div v-if="loading" class="p-8 text-center">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
        <p class="mt-2 text-gray-600">Loading transactions...</p>
      </div>

      <div v-else-if="transactions.length === 0" class="p-8 text-center text-gray-500">
        <p>No transactions found</p>
      </div>

      <div v-else class="divide-y">
        <div
          v-for="transaction in transactions"
          :key="transaction.id"
          class="p-4 hover:bg-gray-50 flex items-center justify-between"
        >
          <div class="flex items-center gap-4">
            <div class="w-10 h-10 rounded-full flex items-center justify-center"
                 :class="transaction.type === 'income' ? 'bg-green-100 text-green-600' : 'bg-red-100 text-red-600'">
              <ArrowUpIcon v-if="transaction.type === 'income'" class="w-5 h-5" />
              <ArrowDownIcon v-else class="w-5 h-5" />
            </div>
            <div>
              <p class="font-medium text-gray-900">{{ transaction.description }}</p>
              <p class="text-sm text-gray-600">
                {{ getCategoryName(transaction.category_id) }} • {{ getAccountName(transaction.account_id) }}
              </p>
              <p class="text-xs text-gray-500">{{ formatDate(transaction.date) }}</p>
            </div>
          </div>
          <div class="text-right">
            <p class="font-semibold"
               :class="transaction.type === 'income' ? 'text-green-600' : 'text-red-600'">
              {{ transaction.type === 'income' ? '+' : '-' }}{{ formatCurrency(transaction.amount) }}
            </p>
            <div class="flex gap-2 mt-1">
              <button
                @click="editTransaction(transaction)"
                class="text-blue-600 hover:text-blue-800 text-sm"
              >
                Edit
              </button>
              <button
                @click="deleteTransaction(transaction.id)"
                class="text-red-600 hover:text-red-800 text-sm"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Transaction Modal -->
    <div v-if="showAddModal || showEditModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4">
        <h3 class="text-lg font-semibold mb-4">
          {{ showEditModal ? 'Edit Transaction' : 'Add Transaction' }}
        </h3>
        
        <form @submit.prevent="submitTransaction" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
            <select v-model="transactionForm.type" required class="w-full border border-gray-300 rounded-md px-3 py-2">
              <option value="income">Income</option>
              <option value="expense">Expense</option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Amount</label>
            <input
              v-model.number="transactionForm.amount"
              type="number"
              step="0.01"
              min="0"
              required
              class="w-full border border-gray-300 rounded-md px-3 py-2"
              placeholder="0.00"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <input
              v-model="transactionForm.description"
              type="text"
              required
              class="w-full border border-gray-300 rounded-md px-3 py-2"
              placeholder="Transaction description"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Category</label>
            <select v-model="transactionForm.category_id" required class="w-full border border-gray-300 rounded-md px-3 py-2">
              <option value="">Select Category</option>
              <option v-for="category in categories" :key="category.id" :value="category.id">
                {{ category.name }}
              </option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Account</label>
            <select v-model="transactionForm.account_id" required class="w-full border border-gray-300 rounded-md px-3 py-2">
              <option value="">Select Account</option>
              <option v-for="account in accounts" :key="account.id" :value="account.id">
                {{ account.name }}
              </option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Date</label>
            <input
              v-model="transactionForm.date"
              type="date"
              required
              class="w-full border border-gray-300 rounded-md px-3 py-2"
            />
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
              class="flex-1 bg-gray-300 text-gray-700 py-2 px-4 rounded-md hover:bg-gray-400"
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
import { ref, reactive, onMounted, watch } from 'vue'
import { PlusIcon, ArrowUpIcon, ArrowDownIcon } from '@heroicons/vue/24/outline'
import { api } from '@/services/api'
import { useCurrency } from '@/composables/useCurrency'
import type { Transaction, Account, Category } from '@/types'

const { formatCurrency } = useCurrency()

const loading = ref(false)
const submitting = ref(false)
const showAddModal = ref(false)
const showEditModal = ref(false)
const transactions = ref<Transaction[]>([])
const accounts = ref<Account[]>([])
const categories = ref<Category[]>([])

const filters = reactive({
  type: '',
  categoryId: '',
  accountId: '',
  dateRange: '30'
})

const transactionForm = reactive({
  id: null as string | null,
  type: 'expense' as 'income' | 'expense',
  amount: 0,
  description: '',
  category_id: '',
  account_id: '',
  date: new Date().toISOString().split('T')[0]
})

const loadTransactions = async () => {
  try {
    loading.value = true
    const response = await api.get('/transactions')
    transactions.value = response.data
  } catch (error) {
    console.error('Failed to load transactions:', error)
  } finally {
    loading.value = false
  }
}

const loadAccounts = async () => {
  try {
    const response = await api.get('/accounts')
    accounts.value = response.data
  } catch (error) {
    console.error('Failed to load accounts:', error)
  }
}

const loadCategories = async () => {
  try {
    const response = await api.get('/categories')
    categories.value = response.data
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

const submitTransaction = async () => {
  try {
    submitting.value = true
    const data = {
      type: transactionForm.type,
      amount: transactionForm.amount,
      description: transactionForm.description,
      category_id: transactionForm.category_id || undefined,
      account_id: transactionForm.account_id,
      date: transactionForm.date
    }
    if (showEditModal.value && transactionForm.id) {
      await api.put(`/transactions/${transactionForm.id}`, data)
    } else {
      await api.post('/transactions', data)
    }
    await loadTransactions()
    closeModal()
  } catch (error) {
    console.error('Failed to save transaction:', error)
  } finally {
    submitting.value = false
  }
}

const editTransaction = (transaction: Transaction) => {
  transactionForm.id = transaction.id
  transactionForm.type = transaction.type === 'transfer' ? 'expense' : transaction.type
  transactionForm.amount = transaction.amount
  transactionForm.description = transaction.description
  transactionForm.category_id = transaction.category_id ?? ''
  transactionForm.account_id = transaction.account_id
  transactionForm.date = transaction.date.split('T')[0]
  showEditModal.value = true
}

const deleteTransaction = async (id: string) => {
  if (!confirm('Are you sure you want to delete this transaction?')) return
  try {
    await api.delete(`/transactions/${id}`)
    await loadTransactions()
  } catch (error) {
    console.error('Failed to delete transaction:', error)
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  transactionForm.id = null
  transactionForm.type = 'expense'
  transactionForm.amount = 0
  transactionForm.description = ''
  transactionForm.category_id = ''
  transactionForm.account_id = ''
  transactionForm.date = new Date().toISOString().split('T')[0]
}

const getCategoryName = (categoryId: string | null) => {
  if (!categoryId) return 'Uncategorized'
  const category = categories.value.find(c => c.id === categoryId)
  return category?.name || 'Unknown'
}

const getAccountName = (accountId: string) => {
  const account = accounts.value.find(a => a.id === accountId)
  return account?.name || 'Unknown'
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

watch(filters, () => {}, { deep: true })

onMounted(() => {
  loadTransactions()
  loadAccounts()
  loadCategories()
})
</script>