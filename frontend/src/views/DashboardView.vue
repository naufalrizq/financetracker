<template>
  <div>
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
      <p class="text-gray-600">Welcome back, {{ authStore.user?.first_name }}!</p>
    </div>

    <!-- Quick Stats -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <div class="stat-card">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-success-100 rounded-lg flex items-center justify-center">
              <span class="text-success-600 text-lg">💰</span>
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">Total Income</p>
            <p class="text-2xl font-bold text-gray-900">
              {{ formatCurrency(summary?.total_income || 0) }}
            </p>
          </div>
        </div>
      </div>

      <div class="stat-card">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-danger-100 rounded-lg flex items-center justify-center">
              <span class="text-danger-600 text-lg">💸</span>
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">Total Expenses</p>
            <p class="text-2xl font-bold text-gray-900">
              {{ formatCurrency(summary?.total_expense || 0) }}
            </p>
          </div>
        </div>
      </div>

      <div class="stat-card">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-primary-100 rounded-lg flex items-center justify-center">
              <span class="text-primary-600 text-lg">📊</span>
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">Net Income</p>
            <p class="text-2xl font-bold" :class="netIncomeClass">
              {{ formatCurrency(summary?.net_income || 0) }}
            </p>
          </div>
        </div>
      </div>

      <div class="stat-card">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-warning-100 rounded-lg flex items-center justify-center">
              <span class="text-warning-600 text-lg">🏦</span>
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">Accounts</p>
            <p class="text-2xl font-bold text-gray-900">
              {{ summary?.total_accounts || 0 }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Charts Row -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
      <!-- Expense Breakdown -->
      <div class="card">
        <div class="card-header">
          <h3 class="text-lg font-medium text-gray-900">Expense Breakdown</h3>
        </div>
        <div class="card-body">
          <div v-if="expenseAnalysis?.categories?.length" class="space-y-4">
            <div
              v-for="category in expenseAnalysis.categories.slice(0, 5)"
              :key="category.category_id"
              class="flex items-center justify-between"
            >
              <div class="flex items-center">
                <div
                  class="w-4 h-4 rounded-full mr-3"
                  :style="{ backgroundColor: category.color }"
                ></div>
                <span class="text-sm font-medium text-gray-900">
                  {{ category.category_name }}
                </span>
              </div>
              <div class="text-right">
                <p class="text-sm font-medium text-gray-900">
                  {{ formatCurrency(category.amount) }}
                </p>
                <p class="text-xs text-gray-500">
                  {{ category.percentage.toFixed(1) }}%
                </p>
              </div>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-500">
            <span class="text-4xl mb-2 block">📊</span>
            <p>No expense data available</p>
          </div>
        </div>
      </div>

      <!-- Recent Transactions -->
      <div class="card">
        <div class="card-header flex justify-between items-center">
          <h3 class="text-lg font-medium text-gray-900">Recent Transactions</h3>
          <router-link
            to="/transactions"
            class="text-sm text-primary-600 hover:text-primary-500"
          >
            View all
          </router-link>
        </div>
        <div class="card-body">
          <div v-if="recentTransactions?.length" class="space-y-4">
            <div
              v-for="transaction in recentTransactions"
              :key="transaction.id"
              class="flex items-center justify-between py-2"
            >
              <div class="flex items-center">
                <div
                  class="w-8 h-8 rounded-lg flex items-center justify-center mr-3"
                  :class="transaction.type === 'income' ? 'bg-success-100' : 'bg-danger-100'"
                >
                  <span :class="transaction.type === 'income' ? 'text-success-600' : 'text-danger-600'">
                    {{ transaction.category?.icon || (transaction.type === 'income' ? '💰' : '💸') }}
                  </span>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-900">
                    {{ transaction.description }}
                  </p>
                  <p class="text-xs text-gray-500">
                    {{ transaction.category?.name || 'Uncategorized' }}
                  </p>
                </div>
              </div>
              <div class="text-right">
                <p
                  class="text-sm font-medium"
                  :class="transaction.type === 'income' ? 'text-success-600' : 'text-danger-600'"
                >
                  {{ transaction.type === 'income' ? '+' : '-' }}{{ formatCurrency(transaction.amount) }}
                </p>
                <p class="text-xs text-gray-500">
                  {{ formatDate(transaction.date) }}
                </p>
              </div>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-500">
            <span class="text-4xl mb-2 block">📝</span>
            <p>No transactions yet</p>
            <router-link
              to="/transactions"
              class="text-sm text-primary-600 hover:text-primary-500 mt-2 inline-block"
            >
              Add your first transaction
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="card">
      <div class="card-header">
        <h3 class="text-lg font-medium text-gray-900">Quick Actions</h3>
      </div>
      <div class="card-body">
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <router-link
            to="/transactions"
            class="flex flex-col items-center p-4 rounded-lg border-2 border-dashed border-gray-300 hover:border-primary-500 hover:bg-primary-50 transition-colors duration-200"
          >
            <span class="text-2xl mb-2">💳</span>
            <span class="text-sm font-medium text-gray-900">Add Transaction</span>
          </router-link>

          <router-link
            to="/accounts"
            class="flex flex-col items-center p-4 rounded-lg border-2 border-dashed border-gray-300 hover:border-primary-500 hover:bg-primary-50 transition-colors duration-200"
          >
            <span class="text-2xl mb-2">🏦</span>
            <span class="text-sm font-medium text-gray-900">Manage Accounts</span>
          </router-link>

          <router-link
            to="/budgets"
            class="flex flex-col items-center p-4 rounded-lg border-2 border-dashed border-gray-300 hover:border-primary-500 hover:bg-primary-50 transition-colors duration-200"
          >
            <span class="text-2xl mb-2">📊</span>
            <span class="text-sm font-medium text-gray-900">Set Budget</span>
          </router-link>

          <router-link
            to="/reports"
            class="flex flex-col items-center p-4 rounded-lg border-2 border-dashed border-gray-300 hover:border-primary-500 hover:bg-primary-50 transition-colors duration-200"
          >
            <span class="text-2xl mb-2">📈</span>
            <span class="text-sm font-medium text-gray-900">View Reports</span>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { reportsApi, transactionsApi } from '@/services/api'
import type { Transaction } from '@/types'

const authStore = useAuthStore()

const summary = ref<any>(null)
const expenseAnalysis = ref<any>(null)
const recentTransactions = ref<Transaction[]>([])
const isLoading = ref(true)

const netIncomeClass = computed(() => {
  const netIncome = summary.value?.net_income || 0
  return netIncome >= 0 ? 'text-success-600' : 'text-danger-600'
})

const formatCurrency = (amount: number): string => {
  const currency = authStore.user?.currency || 'IDR'
  const locale = currency === 'IDR' ? 'id-ID' : 'en-US'
  return new Intl.NumberFormat(locale, {
    style: 'currency',
    currency,
    minimumFractionDigits: currency === 'IDR' ? 0 : 2,
  }).format(amount)
}

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
  })
}

const loadDashboardData = async () => {
  try {
    isLoading.value = true

    // Get current month date range
    const now = new Date()
    const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
    const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0)

    const dateFrom = startOfMonth.toISOString().split('T')[0]
    const dateTo = endOfMonth.toISOString().split('T')[0]

    // Load summary data
    const [summaryResponse, expenseResponse, transactionsResponse] = await Promise.all([
      reportsApi.getFinancialSummary(dateFrom, dateTo),
      reportsApi.getExpenseAnalysis(dateFrom, dateTo),
      transactionsApi.getTransactions({ limit: 5, sort_by: 'date', sort_order: 'desc' })
    ])

    summary.value = summaryResponse.summary
    expenseAnalysis.value = expenseResponse.expense_analysis
    recentTransactions.value = transactionsResponse.transactions
  } catch (error) {
    console.error('Failed to load dashboard data:', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  loadDashboardData()
})
</script>