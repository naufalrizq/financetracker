<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Reports & Analytics</h1>
        <p class="text-gray-600">Analyze your financial data and spending patterns</p>
      </div>
      <div class="flex gap-3">
        <select v-model="selectedPeriod" class="form-select bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-700 text-gray-900 dark:text-gray-100 rounded-md py-2 pl-3 pr-10 focus:ring-primary-500 focus:border-primary-500 shadow-sm appearance-none cursor-pointer">
          <option value="7">Last 7 days</option>
          <option value="30">Last 30 days</option>
          <option value="90">Last 3 months</option>
          <option value="365">Last year</option>
        </select>
        <button
          @click="exportReport"
          class="bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 flex items-center gap-2"
        >
          <DocumentArrowDownIcon class="w-5 h-5" />
          Export
        </button>
      </div>
    </div>

    <!-- Key Metrics -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <div class="bg-gradient-to-r from-green-500 to-green-600 text-white p-6 rounded-lg">
        <h3 class="text-lg font-semibold mb-2">Total Income</h3>
        <p class="text-3xl font-bold">{{ formatCurrency(totalIncome) }}</p>
        <p class="text-sm opacity-90">{{ selectedPeriod }} days</p>
      </div>
      <div class="bg-gradient-to-r from-red-500 to-red-600 text-white p-6 rounded-lg">
        <h3 class="text-lg font-semibold mb-2">Total Expenses</h3>
        <p class="text-3xl font-bold">{{ formatCurrency(totalExpenses) }}</p>
        <p class="text-sm opacity-90">{{ selectedPeriod }} days</p>
      </div>
      <div class="bg-gradient-to-r from-blue-500 to-blue-600 text-white p-6 rounded-lg">
        <h3 class="text-lg font-semibold mb-2">Net Income</h3>
        <p class="text-3xl font-bold">{{ formatCurrency(netIncome) }}</p>
        <p class="text-sm opacity-90" :class="netIncome >= 0 ? 'text-green-200' : 'text-red-200'">
          {{ netIncome >= 0 ? 'Surplus' : 'Deficit' }}
        </p>
      </div>
      <div class="bg-gradient-to-r from-purple-500 to-purple-600 text-white p-6 rounded-lg">
        <h3 class="text-lg font-semibold mb-2">Avg Daily Spending</h3>
        <p class="text-3xl font-bold">{{ formatCurrency(averageDailySpending) }}</p>
        <p class="text-sm opacity-90">per day</p>
      </div>
    </div>

    <!-- Charts Row -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Income vs Expenses Chart -->
      <div class="bg-white p-6 rounded-lg shadow-sm border">
        <h3 class="text-lg font-semibold mb-4">Income vs Expenses</h3>
        <div class="h-64 flex items-center justify-center">
          <div class="text-center">
            <ChartBarIcon class="w-16 h-16 text-gray-400 mx-auto mb-4" />
            <p class="text-gray-500">Chart visualization would go here</p>
            <p class="text-sm text-gray-400">Integration with Chart.js or similar library</p>
          </div>
        </div>
      </div>

      <!-- Expense Categories Pie Chart -->
      <div class="bg-white p-6 rounded-lg shadow-sm border">
        <h3 class="text-lg font-semibold mb-4">Expenses by Category</h3>
        <div class="h-64 flex items-center justify-center">
          <div class="text-center">
            <ChartPieIcon class="w-16 h-16 text-gray-400 mx-auto mb-4" />
            <p class="text-gray-500">Pie chart visualization would go here</p>
            <p class="text-sm text-gray-400">Shows spending breakdown by category</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Detailed Reports -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Top Categories -->
      <div class="bg-white rounded-lg shadow-sm border">
        <div class="p-4 border-b">
          <h3 class="text-lg font-semibold">Top Spending Categories</h3>
        </div>
        <div class="p-4">
          <div v-if="loading" class="text-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
          </div>
          <div v-else-if="topCategories.length === 0" class="text-center py-8 text-gray-500">
            <p>No spending data available</p>
          </div>
          <div v-else class="space-y-4">
            <div
              v-for="(category, index) in topCategories"
              :key="category.id"
              class="flex items-center justify-between"
            >
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold"
                     :class="getCategoryRankColor(index)">
                  {{ index + 1 }}
                </div>
                <div>
                  <p class="font-medium text-gray-900">{{ category.name }}</p>
                  <p class="text-sm text-gray-600">{{ category.transactionCount }} transactions</p>
                </div>
              </div>
              <div class="text-right">
                <p class="font-semibold text-gray-900">{{ formatCurrency(category.totalAmount) }}</p>
                <p class="text-sm text-gray-600">{{ category.percentage.toFixed(1) }}%</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Monthly Trends -->
      <div class="bg-white rounded-lg shadow-sm border">
        <div class="p-4 border-b">
          <h3 class="text-lg font-semibold">Monthly Trends</h3>
        </div>
        <div class="p-4">
          <div v-if="loading" class="text-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
          </div>
          <div v-else-if="monthlyTrends.length === 0" class="text-center py-8 text-gray-500">
            <p>No trend data available</p>
          </div>
          <div v-else class="space-y-4">
            <div
              v-for="trend in monthlyTrends"
              :key="trend.month"
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div>
                <p class="font-medium text-gray-900">{{ trend.monthName }}</p>
                <p class="text-sm text-gray-600">{{ trend.transactionCount }} transactions</p>
              </div>
              <div class="text-right">
                <p class="font-semibold text-gray-900">{{ formatCurrency(trend.totalExpenses) }}</p>
                <p class="text-sm" :class="trend.change >= 0 ? 'text-red-600' : 'text-green-600'">
                  {{ trend.change >= 0 ? '+' : '' }}{{ trend.change.toFixed(1) }}%
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Account Balances -->
    <div class="bg-white rounded-lg shadow-sm border">
      <div class="p-4 border-b">
        <h3 class="text-lg font-semibold">Account Balances</h3>
      </div>
      <div class="p-4">
        <div v-if="loading" class="text-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
        </div>
        <div v-else-if="accounts.length === 0" class="text-center py-8 text-gray-500">
          <p>No accounts found</p>
        </div>
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div
            v-for="account in accounts"
            :key="account.id"
            class="p-4 border rounded-lg"
          >
            <div class="flex items-center justify-between mb-2">
              <h4 class="font-medium text-gray-900">{{ account.name }}</h4>
              <span class="text-sm text-gray-600 capitalize">{{ account.type.replace('_', ' ') }}</span>
            </div>
            <p class="text-2xl font-bold text-gray-900">{{ formatCurrency(account.balance) }}</p>
            <div class="mt-2 flex justify-between text-sm text-gray-600">
              <span>Last updated</span>
              <span>{{ formatDate(account.updated_at) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Budget Performance -->
    <div class="bg-white rounded-lg shadow-sm border">
      <div class="p-4 border-b">
        <h3 class="text-lg font-semibold">Budget Performance</h3>
      </div>
      <div class="p-4">
        <div v-if="loading" class="text-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
        </div>
        <div v-else-if="budgets.length === 0" class="text-center py-8 text-gray-500">
          <p>No budgets found</p>
          <button
            @click="$router.push('/budgets')"
            class="mt-2 text-blue-600 hover:text-blue-800"
          >
            Create your first budget
          </button>
        </div>
        <div v-else class="space-y-4">
          <div
            v-for="budget in budgets"
            :key="budget.id"
            class="p-4 border rounded-lg"
          >
            <div class="flex items-center justify-between mb-3">
              <div>
                <h4 class="font-medium text-gray-900">
                  {{ budget.category?.name || 'All Categories' }}
                </h4>
                <p class="text-sm text-gray-600 capitalize">{{ budget.period }} Budget</p>
              </div>
              <div class="text-right">
                <p class="font-semibold text-gray-900">
                  {{ formatCurrency(getBudgetSpent(budget)) }} / {{ formatCurrency(budget.amount) }}
                </p>
                <p class="text-sm" :class="getBudgetProgress(budget) > 100 ? 'text-red-600' : 'text-green-600'">
                  {{ getBudgetProgress(budget).toFixed(1) }}% used
                </p>
              </div>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2">
              <div
                class="h-2 rounded-full transition-all duration-300"
                :class="getBudgetProgress(budget) > 100 ? 'bg-red-500' : 
                       getBudgetProgress(budget) > 80 ? 'bg-yellow-500' : 'bg-green-500'"
                :style="{ width: Math.min(getBudgetProgress(budget), 100) + '%' }"
              ></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { 
  DocumentArrowDownIcon,
  ChartBarIcon,
  ChartPieIcon
} from '@heroicons/vue/24/outline'
import { api } from '@/services/api'
import { useCurrency } from '@/composables/useCurrency'
import type { Transaction, Account, Category, Budget } from '@/types'

const { formatCurrency } = useCurrency()

// Reactive data
const loading = ref(false)
const selectedPeriod = ref('30')
const transactions = ref<Transaction[]>([])
const accounts = ref<Account[]>([])
const categories = ref<Category[]>([])
const budgets = ref<Budget[]>([])

// Computed
const filteredTransactions = computed(() => {
  const now = new Date()
  const daysAgo = new Date(now.getTime() - parseInt(selectedPeriod.value) * 24 * 60 * 60 * 1000)
  return transactions.value.filter(t => new Date(t.date) >= daysAgo)
})

const totalIncome = computed(() => {
  return filteredTransactions.value
    .filter(t => t.type === 'income')
    .reduce((sum, t) => sum + t.amount, 0)
})

const totalExpenses = computed(() => {
  return filteredTransactions.value
    .filter(t => t.type === 'expense')
    .reduce((sum, t) => sum + t.amount, 0)
})

const netIncome = computed(() => {
  return totalIncome.value - totalExpenses.value
})

const averageDailySpending = computed(() => {
  const days = parseInt(selectedPeriod.value)
  return days > 0 ? totalExpenses.value / days : 0
})

const topCategories = computed(() => {
  const categoryTotals = new Map<string, {
    id: string
    name: string
    totalAmount: number
    transactionCount: number
    percentage: number
  }>()

  // Calculate totals for each category
  filteredTransactions.value
    .filter(t => t.type === 'expense')
    .forEach(transaction => {
      const category = categories.value.find(c => c.id === transaction.category_id)
      if (category) {
        const existing = categoryTotals.get(category.id) || {
          id: category.id,
          name: category.name,
          totalAmount: 0,
          transactionCount: 0,
          percentage: 0
        }
        existing.totalAmount += transaction.amount
        existing.transactionCount += 1
        categoryTotals.set(category.id, existing)
      }
    })

  // Calculate percentages and sort
  const total = totalExpenses.value
  const result = Array.from(categoryTotals.values())
    .map(cat => ({
      ...cat,
      percentage: total > 0 ? (cat.totalAmount / total) * 100 : 0
    }))
    .sort((a, b) => b.totalAmount - a.totalAmount)
    .slice(0, 5)

  return result
})

const monthlyTrends = computed(() => {
  const monthlyData = new Map<string, {
    month: string
    monthName: string
    totalExpenses: number
    transactionCount: number
    change: number
  }>()

  // Group transactions by month
  filteredTransactions.value
    .filter(t => t.type === 'expense')
    .forEach(transaction => {
      const date = new Date(transaction.date)
      const monthKey = `${date.getFullYear()}-${date.getMonth()}`
      const monthName = date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
      
      const existing = monthlyData.get(monthKey) || {
        month: monthKey,
        monthName,
        totalExpenses: 0,
        transactionCount: 0,
        change: 0
      }
      existing.totalExpenses += transaction.amount
      existing.transactionCount += 1
      monthlyData.set(monthKey, existing)
    })

  // Calculate month-over-month changes
  const sortedMonths = Array.from(monthlyData.values())
    .sort((a, b) => a.month.localeCompare(b.month))

  for (let i = 1; i < sortedMonths.length; i++) {
    const current = sortedMonths[i]
    const previous = sortedMonths[i - 1]
    if (previous.totalExpenses > 0) {
      current.change = ((current.totalExpenses - previous.totalExpenses) / previous.totalExpenses) * 100
    }
  }

  return sortedMonths.slice(-6) // Last 6 months
})

// Methods
const loadData = async () => {
  try {
    loading.value = true
    
    const [transactionsRes, accountsRes, categoriesRes, budgetsRes] = await Promise.all([
      api.get('/transactions'),
      api.get('/accounts'),
      api.get('/categories'),
      api.get('/budgets')
    ])
    
    transactions.value = transactionsRes.data
    accounts.value = accountsRes.data
    categories.value = categoriesRes.data
    budgets.value = budgetsRes.data
  } catch (error) {
    console.error('Failed to load report data:', error)
  } finally {
    loading.value = false
  }
}

const getBudgetSpent = (budget: Budget) => {
  const startDate = new Date(budget.start_date)
  const endDate = budget.end_date ? new Date(budget.end_date) : new Date()
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

const getCategoryRankColor = (index: number) => {
  const colors = [
    'bg-yellow-100 text-yellow-800', // 1st place
    'bg-gray-100 text-gray-800',     // 2nd place
    'bg-orange-100 text-orange-800', // 3rd place
    'bg-blue-100 text-blue-800',     // 4th place
    'bg-purple-100 text-purple-800'  // 5th place
  ]
  return colors[index] || 'bg-gray-100 text-gray-800'
}

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) return 'N/A'
  return new Date(dateString).toLocaleDateString()
}

const exportReport = () => {
  // This would generate and download a CSV or PDF report
  console.log('Exporting report for period:', selectedPeriod.value)
  
  // Example CSV export
  const csvData = [
    ['Date', 'Type', 'Category', 'Description', 'Amount'],
    ...filteredTransactions.value.map(t => [
      t.date.split('T')[0],
      t.type,
      categories.value.find(c => c.id === t.category_id)?.name || 'Unknown',
      t.description,
      t.amount.toFixed(2)
    ])
  ]
  
  const csvContent = csvData.map(row => row.join(',')).join('\n')
  const blob = new Blob([csvContent], { type: 'text/csv' })
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `financial-report-${selectedPeriod.value}-days.csv`
  a.click()
  window.URL.revokeObjectURL(url)
}

// Watch for period changes
watch(selectedPeriod, () => {
  // Data will be automatically recalculated due to computed properties
})

// Lifecycle
onMounted(() => {
  loadData()
})
</script>