<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Budgets</h1>
        <p class="text-gray-600">Set spending limits and track your progress</p>
      </div>
      <button
        @click="showAddModal = true"
        class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 flex items-center gap-2"
      >
        <PlusIcon class="w-5 h-5" />
        Create Budget
      </button>
    </div>

    <!-- Budget Overview -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-gradient-to-r from-green-500 to-green-600 text-white p-6 rounded-lg">
        <h3 class="text-lg font-semibold mb-2">Total Budget</h3>
        <p class="text-3xl font-bold">{{ formatCurrency(totalBudget) }}</p>
      </div>
      <div class="bg-gradient-to-r from-blue-500 to-blue-600 text-white p-6 rounded-lg">
        <h3 class="text-lg font-semibold mb-2">Total Spent</h3>
        <p class="text-3xl font-bold">{{ formatCurrency(totalSpent) }}</p>
      </div>
      <div class="bg-gradient-to-r from-purple-500 to-purple-600 text-white p-6 rounded-lg">
        <h3 class="text-lg font-semibold mb-2">Remaining</h3>
        <p class="text-3xl font-bold">{{ formatCurrency(totalRemaining) }}</p>
      </div>
    </div>

    <!-- Budget Period Filter -->
    <div class="bg-white p-4 rounded-lg shadow-sm border">
      <div class="flex items-center gap-4">
        <label class="text-sm font-medium text-gray-700">Period:</label>
        <select v-model="selectedPeriod" class="border border-gray-300 rounded-md px-3 py-2">
          <option value="monthly">Monthly</option>
          <option value="weekly">Weekly</option>
          <option value="yearly">Yearly</option>
        </select>
        <div class="flex items-center gap-2 ml-auto">
          <button
            @click="previousPeriod"
            class="p-2 text-gray-600 hover:text-gray-800"
          >
            <ChevronLeftIcon class="w-5 h-5" />
          </button>
          <span class="text-sm font-medium text-gray-900 min-w-[120px] text-center">
            {{ currentPeriodLabel }}
          </span>
          <button
            @click="nextPeriod"
            class="p-2 text-gray-600 hover:text-gray-800"
          >
            <ChevronRightIcon class="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>

    <!-- Budgets List -->
    <div class="bg-white rounded-lg shadow-sm border">
      <div class="p-4 border-b">
        <h2 class="text-lg font-semibold">Your Budgets</h2>
      </div>
      
      <div v-if="loading" class="p-8 text-center">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
        <p class="mt-2 text-gray-600">Loading budgets...</p>
      </div>

      <div v-else-if="budgets.length === 0" class="p-8 text-center text-gray-500">
        <ChartBarIcon class="w-12 h-12 mx-auto mb-4 text-gray-400" />
        <p class="text-lg font-medium mb-2">No budgets yet</p>
        <p class="mb-4">Create your first budget to start tracking spending</p>
        <button
          @click="showAddModal = true"
          class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
        >
          Create Budget
        </button>
      </div>

      <div v-else class="divide-y">
        <div
          v-for="budget in budgets"
          :key="budget.id"
          class="p-4"
        >
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-full flex items-center justify-center"
                   :class="getBudgetStatusColor(budget)">
                <component :is="getCategoryIcon(budget.category?.icon)" class="w-5 h-5" />
              </div>
              <div>
                <h3 class="font-semibold text-gray-900">{{ budget.category?.name || 'All Categories' }}</h3>
                <p class="text-sm text-gray-600 capitalize">{{ budget.period }} Budget</p>
              </div>
            </div>
            <div class="text-right">
              <p class="text-lg font-bold text-gray-900">
                {{ formatCurrency(getBudgetSpent(budget)) }} / {{ formatCurrency(budget.amount) }}
              </p>
              <div class="flex gap-2 mt-1">
                <button
                  @click="editBudget(budget)"
                  class="text-blue-600 hover:text-blue-800 text-sm"
                >
                  Edit
                </button>
                <button
                  @click="deleteBudget(budget.id)"
                  class="text-red-600 hover:text-red-800 text-sm"
                >
                  Delete
                </button>
              </div>
            </div>
          </div>
          
          <!-- Progress Bar -->
          <div class="mb-3">
            <div class="flex justify-between text-sm text-gray-600 mb-1">
              <span>Progress</span>
              <span>{{ getBudgetProgress(budget).toFixed(1) }}%</span>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2">
              <div
                class="h-2 rounded-full transition-all duration-300"
                :class="getBudgetProgressColor(budget)"
                :style="{ width: Math.min(getBudgetProgress(budget), 100) + '%' }"
              ></div>
            </div>
          </div>
          
          <!-- Budget Status -->
          <div class="flex justify-between items-center text-sm">
            <span class="text-gray-600">
            Remaining: {{ formatCurrency(Math.max(budget.amount - getBudgetSpent(budget), 0)) }}
            </span>
            <span
              :class="getBudgetProgress(budget) > 100 ? 'text-red-600 font-medium' : 
                     getBudgetProgress(budget) > 80 ? 'text-yellow-600 font-medium' : 'text-green-600'"
            >
              {{ getBudgetStatusText(budget) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Budget Modal -->
    <div v-if="showAddModal || showEditModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4">
        <h3 class="text-lg font-semibold mb-4">
          {{ showEditModal ? 'Edit Budget' : 'Create Budget' }}
        </h3>
        
        <form @submit.prevent="submitBudget" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Category</label>
            <select v-model="budgetForm.category_id" class="w-full border border-gray-300 rounded-md px-3 py-2">
              <option value="">All Categories</option>
              <option v-for="category in expenseCategories" :key="category.id" :value="category.id">
                {{ category.name }}
              </option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Budget Amount</label>
            <input
              v-model.number="budgetForm.amount"
              type="number"
              step="0.01"
              min="0"
              required
              class="w-full border border-gray-300 rounded-md px-3 py-2"
              placeholder="0.00"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Period</label>
            <select v-model="budgetForm.period" required class="w-full border border-gray-300 rounded-md px-3 py-2">
              <option value="weekly">Weekly</option>
              <option value="monthly">Monthly</option>
              <option value="yearly">Yearly</option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Start Date</label>
            <input
              v-model="budgetForm.start_date"
              type="date"
              required
              class="w-full border border-gray-300 rounded-md px-3 py-2"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">End Date</label>
            <input
              v-model="budgetForm.end_date"
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
              {{ submitting ? 'Saving...' : (showEditModal ? 'Update' : 'Create') }}
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
import { ref, reactive, onMounted, computed } from 'vue'
import { 
  PlusIcon, 
  ChartBarIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  ShoppingCartIcon
} from '@heroicons/vue/24/outline'
import { api } from '@/services/api'
import { useCurrency } from '@/composables/useCurrency'
import type { Budget, Category, Transaction } from '@/types'

const { formatCurrency } = useCurrency()

// Reactive data
const loading = ref(false)
const submitting = ref(false)
const showAddModal = ref(false)
const showEditModal = ref(false)
const selectedPeriod = ref('monthly')
const currentDate = ref(new Date())
const budgets = ref<Budget[]>([])
const categories = ref<Category[]>([])
const transactions = ref<Transaction[]>([])

// Budget form
const budgetForm = reactive({
  id: null as number | null,
  category_id: '',
  amount: 0,
  period: 'monthly',
  start_date: '',
  end_date: ''
})

// Computed
const expenseCategories = computed(() => {
  return categories.value.filter(category => category.type === 'expense')
})

const totalBudget = computed(() => {
  return budgets.value.reduce((sum, budget) => sum + budget.amount, 0)
})

const totalSpent = computed(() => {
  return budgets.value.reduce((sum, budget) => sum + getBudgetSpent(budget), 0)
})

const totalRemaining = computed(() => {
  return totalBudget.value - totalSpent.value
})

const currentPeriodLabel = computed(() => {
  const date = currentDate.value
  switch (selectedPeriod.value) {
    case 'weekly':
      return `Week of ${date.toLocaleDateString()}`
    case 'monthly':
      return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
    case 'yearly':
      return date.getFullYear().toString()
    default:
      return ''
  }
})

// Methods
const loadBudgets = async () => {
  try {
    loading.value = true
    const response = await api.get('/budgets')
    budgets.value = response.data
  } catch (error) {
    console.error('Failed to load budgets:', error)
  } finally {
    loading.value = false
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

const loadTransactions = async () => {
  try {
    const response = await api.get('/transactions')
    transactions.value = response.data
  } catch (error) {
    console.error('Failed to load transactions:', error)
  }
}

const submitBudget = async () => {
  try {
    submitting.value = true
    
    const data = {
      category_id: budgetForm.category_id ? parseInt(budgetForm.category_id) : null,
      amount: budgetForm.amount,
      period: budgetForm.period,
      start_date: budgetForm.start_date,
      end_date: budgetForm.end_date
    }
    
    if (showEditModal.value && budgetForm.id) {
      await api.put(`/budgets/${budgetForm.id}`, data)
    } else {
      await api.post('/budgets', data)
    }
    
    await loadBudgets()
    closeModal()
  } catch (error) {
    console.error('Failed to save budget:', error)
  } finally {
    submitting.value = false
  }
}

const editBudget = (budget: Budget) => {
  budgetForm.id = budget.id
  budgetForm.category_id = budget.category_id?.toString() || ''
  budgetForm.amount = budget.amount
  budgetForm.period = budget.period
  budgetForm.start_date = budget.start_date.split('T')[0]
  budgetForm.end_date = budget.end_date.split('T')[0]
  showEditModal.value = true
}

const deleteBudget = async (id: number) => {
  if (!confirm('Are you sure you want to delete this budget?')) return
  
  try {
    await api.delete(`/budgets/${id}`)
    await loadBudgets()
  } catch (error) {
    console.error('Failed to delete budget:', error)
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  budgetForm.id = null
  budgetForm.category_id = ''
  budgetForm.amount = 0
  budgetForm.period = 'monthly'
  budgetForm.start_date = ''
  budgetForm.end_date = ''
}

const getBudgetSpent = (budget: Budget) => {
  const startDate = new Date(budget.start_date)
  const endDate = new Date(budget.end_date)
  
  return transactions.value
    .filter(transaction => {
      const transactionDate = new Date(transaction.date)
      const isInPeriod = transactionDate >= startDate && transactionDate <= endDate
      const isExpense = transaction.type === 'expense'
      const matchesCategory = !budget.category_id || transaction.category_id === budget.category_id
      
      return isInPeriod && isExpense && matchesCategory
    })
    .reduce((sum, transaction) => sum + transaction.amount, 0)
}

const getBudgetProgress = (budget: Budget) => {
  const spent = getBudgetSpent(budget)
  return budget.amount > 0 ? (spent / budget.amount) * 100 : 0
}

const getBudgetProgressColor = (budget: Budget) => {
  const progress = getBudgetProgress(budget)
  if (progress > 100) return 'bg-red-500'
  if (progress > 80) return 'bg-yellow-500'
  return 'bg-green-500'
}

const getBudgetStatusColor = (budget: Budget) => {
  const progress = getBudgetProgress(budget)
  if (progress > 100) return 'bg-red-100 text-red-600'
  if (progress > 80) return 'bg-yellow-100 text-yellow-600'
  return 'bg-green-100 text-green-600'
}

const getBudgetStatusText = (budget: Budget) => {
  const progress = getBudgetProgress(budget)
  if (progress > 100) return 'Over Budget'
  if (progress > 80) return 'Near Limit'
  return 'On Track'
}

const getCategoryIcon = (iconName?: string) => {
  // This would need to be implemented based on your icon system
  return ShoppingCartIcon
}

const previousPeriod = () => {
  const date = new Date(currentDate.value)
  switch (selectedPeriod.value) {
    case 'weekly':
      date.setDate(date.getDate() - 7)
      break
    case 'monthly':
      date.setMonth(date.getMonth() - 1)
      break
    case 'yearly':
      date.setFullYear(date.getFullYear() - 1)
      break
  }
  currentDate.value = date
}

const nextPeriod = () => {
  const date = new Date(currentDate.value)
  switch (selectedPeriod.value) {
    case 'weekly':
      date.setDate(date.getDate() + 7)
      break
    case 'monthly':
      date.setMonth(date.getMonth() + 1)
      break
    case 'yearly':
      date.setFullYear(date.getFullYear() + 1)
      break
  }
  currentDate.value = date
}

// Lifecycle
onMounted(() => {
  loadBudgets()
  loadCategories()
  loadTransactions()
  
  // Set default dates for new budget
  const today = new Date()
  const nextMonth = new Date(today.getFullYear(), today.getMonth() + 1, 0)
  budgetForm.start_date = today.toISOString().split('T')[0]
  budgetForm.end_date = nextMonth.toISOString().split('T')[0]
})
</script>