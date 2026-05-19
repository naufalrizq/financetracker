import axios, { type AxiosInstance, type AxiosResponse } from 'axios'
import type { 
  AuthResponse, 
  LoginForm, 
  RegisterForm, 
  User,
  Transaction,
  TransactionResponse,
  TransactionFilter,
  CreateTransactionForm,
  Account,
  CreateAccountForm,
  Category,
  CreateCategoryForm,
  Budget,
  CreateBudgetForm,
  Goal,
  CreateGoalForm
} from '@/types'

// Create axios instance
export const api: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8000/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        const refreshToken = localStorage.getItem('refresh_token')
        if (refreshToken) {
          const response = await authApi.refreshToken(refreshToken)
          localStorage.setItem('access_token', response.access_token)
          localStorage.setItem('refresh_token', response.refresh_token)
          
          // Retry original request with new token
          originalRequest.headers.Authorization = `Bearer ${response.access_token}`
          return api(originalRequest)
        }
      } catch (refreshError) {
        // Refresh failed, redirect to login
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        localStorage.removeItem('user')
        window.location.href = '/login'
      }
    }

    return Promise.reject(error)
  }
)

// Auth API
export const authApi = {
  login: async (credentials: LoginForm): Promise<AuthResponse> => {
    const response: AxiosResponse<AuthResponse> = await api.post('/auth/login', credentials)
    return response.data
  },

  register: async (userData: RegisterForm): Promise<AuthResponse> => {
    const response: AxiosResponse<AuthResponse> = await api.post('/auth/register', userData)
    return response.data
  },

  refreshToken: async (refreshToken: string): Promise<AuthResponse> => {
    const response: AxiosResponse<AuthResponse> = await api.post('/auth/refresh', {
      refresh_token: refreshToken
    })
    return response.data
  },

  logout: async (): Promise<void> => {
    await api.post('/auth/logout')
  },

  getProfile: async (): Promise<{ user: User }> => {
    const response: AxiosResponse<{ user: User }> = await api.get('/users/profile')
    return response.data
  },

  updateProfile: async (profileData: Partial<User>): Promise<{ user: User }> => {
    const response: AxiosResponse<{ user: User }> = await api.put('/users/profile', profileData)
    return response.data
  },

  deleteAccount: async (): Promise<void> => {
    await api.delete('/users/account')
  },
}

// Transactions API
export const transactionsApi = {
  getTransactions: async (filters?: TransactionFilter): Promise<TransactionResponse> => {
    const params = new URLSearchParams()
    
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          params.append(key, String(value))
        }
      })
    }

    const response: AxiosResponse<TransactionResponse> = await api.get(`/transactions?${params}`)
    return response.data
  },

  getTransaction: async (id: string): Promise<{ transaction: Transaction }> => {
    const response: AxiosResponse<{ transaction: Transaction }> = await api.get(`/transactions/${id}`)
    return response.data
  },

  createTransaction: async (transactionData: CreateTransactionForm): Promise<{ transaction: Transaction }> => {
    const response: AxiosResponse<{ transaction: Transaction }> = await api.post('/transactions', transactionData)
    return response.data
  },

  updateTransaction: async (id: string, transactionData: Partial<CreateTransactionForm>): Promise<{ transaction: Transaction }> => {
    const response: AxiosResponse<{ transaction: Transaction }> = await api.put(`/transactions/${id}`, transactionData)
    return response.data
  },

  deleteTransaction: async (id: string): Promise<void> => {
    await api.delete(`/transactions/${id}`)
  },

  createBulkTransactions: async (transactions: CreateTransactionForm[]): Promise<{ transactions: Transaction[] }> => {
    const response: AxiosResponse<{ transactions: Transaction[] }> = await api.post('/transactions/bulk', {
      transactions
    })
    return response.data
  },
}

// Accounts API
export const accountsApi = {
  getAccounts: async (): Promise<{ accounts: Account[] }> => {
    const response: AxiosResponse<{ accounts: Account[] }> = await api.get('/accounts')
    return response.data
  },

  getAccount: async (id: string): Promise<{ account: Account }> => {
    const response: AxiosResponse<{ account: Account }> = await api.get(`/accounts/${id}`)
    return response.data
  },

  createAccount: async (accountData: CreateAccountForm): Promise<{ account: Account }> => {
    const response: AxiosResponse<{ account: Account }> = await api.post('/accounts', accountData)
    return response.data
  },

  updateAccount: async (id: string, accountData: Partial<CreateAccountForm>): Promise<{ account: Account }> => {
    const response: AxiosResponse<{ account: Account }> = await api.put(`/accounts/${id}`, accountData)
    return response.data
  },

  deleteAccount: async (id: string): Promise<void> => {
    await api.delete(`/accounts/${id}`)
  },

  getAccountSummary: async (id: string): Promise<{ account_summary: any }> => {
    const response: AxiosResponse<{ account_summary: any }> = await api.get(`/accounts/${id}/summary`)
    return response.data
  },
}

// Categories API
export const categoriesApi = {
  getCategories: async (type?: 'income' | 'expense'): Promise<{ categories: Category[] }> => {
    const params = type ? `?type=${type}` : ''
    const response: AxiosResponse<{ categories: Category[] }> = await api.get(`/categories${params}`)
    return response.data
  },

  getCategory: async (id: string): Promise<{ category: Category }> => {
    const response: AxiosResponse<{ category: Category }> = await api.get(`/categories/${id}`)
    return response.data
  },

  createCategory: async (categoryData: CreateCategoryForm): Promise<{ category: Category }> => {
    const response: AxiosResponse<{ category: Category }> = await api.post('/categories', categoryData)
    return response.data
  },

  updateCategory: async (id: string, categoryData: Partial<CreateCategoryForm>): Promise<{ category: Category }> => {
    const response: AxiosResponse<{ category: Category }> = await api.put(`/categories/${id}`, categoryData)
    return response.data
  },

  deleteCategory: async (id: string): Promise<void> => {
    await api.delete(`/categories/${id}`)
  },
}

// Budgets API
export const budgetsApi = {
  getBudgets: async (): Promise<{ budgets: Budget[] }> => {
    const response: AxiosResponse<{ budgets: Budget[] }> = await api.get('/budgets')
    return response.data
  },

  getBudget: async (id: string): Promise<{ budget: Budget }> => {
    const response: AxiosResponse<{ budget: Budget }> = await api.get(`/budgets/${id}`)
    return response.data
  },

  createBudget: async (budgetData: CreateBudgetForm): Promise<{ budget: Budget }> => {
    const response: AxiosResponse<{ budget: Budget }> = await api.post('/budgets', budgetData)
    return response.data
  },

  updateBudget: async (id: string, budgetData: Partial<CreateBudgetForm>): Promise<{ budget: Budget }> => {
    const response: AxiosResponse<{ budget: Budget }> = await api.put(`/budgets/${id}`, budgetData)
    return response.data
  },

  deleteBudget: async (id: string): Promise<void> => {
    await api.delete(`/budgets/${id}`)
  },
}

// Goals API
export const goalsApi = {
  getGoals: async (): Promise<{ goals: Goal[] }> => {
    const response: AxiosResponse<{ goals: Goal[] }> = await api.get('/goals')
    return response.data
  },

  getGoal: async (id: string): Promise<{ goal: Goal }> => {
    const response: AxiosResponse<{ goal: Goal }> = await api.get(`/goals/${id}`)
    return response.data
  },

  createGoal: async (goalData: CreateGoalForm): Promise<{ goal: Goal }> => {
    const response: AxiosResponse<{ goal: Goal }> = await api.post('/goals', goalData)
    return response.data
  },

  updateGoal: async (id: string, goalData: Partial<CreateGoalForm>): Promise<{ goal: Goal }> => {
    const response: AxiosResponse<{ goal: Goal }> = await api.put(`/goals/${id}`, goalData)
    return response.data
  },

  deleteGoal: async (id: string): Promise<void> => {
    await api.delete(`/goals/${id}`)
  },

  updateGoalProgress: async (id: string, amount: number): Promise<{ goal: Goal }> => {
    const response: AxiosResponse<{ goal: Goal }> = await api.post(`/goals/${id}/progress`, { amount })
    return response.data
  },
}

// Reports API
export const reportsApi = {
  getFinancialSummary: async (dateFrom?: string, dateTo?: string): Promise<{ summary: any }> => {
    const params = new URLSearchParams()
    if (dateFrom) params.append('date_from', dateFrom)
    if (dateTo) params.append('date_to', dateTo)
    
    const response: AxiosResponse<{ summary: any }> = await api.get(`/reports/summary?${params}`)
    return response.data
  },

  getExpenseAnalysis: async (dateFrom?: string, dateTo?: string): Promise<{ expense_analysis: any }> => {
    const params = new URLSearchParams()
    if (dateFrom) params.append('date_from', dateFrom)
    if (dateTo) params.append('date_to', dateTo)
    
    const response: AxiosResponse<{ expense_analysis: any }> = await api.get(`/reports/expenses?${params}`)
    return response.data
  },

  getIncomeAnalysis: async (dateFrom?: string, dateTo?: string): Promise<{ income_analysis: any }> => {
    const params = new URLSearchParams()
    if (dateFrom) params.append('date_from', dateFrom)
    if (dateTo) params.append('date_to', dateTo)
    
    const response: AxiosResponse<{ income_analysis: any }> = await api.get(`/reports/income?${params}`)
    return response.data
  },

  getSpendingTrends: async (period: string = 'daily', dateFrom?: string, dateTo?: string): Promise<{ spending_trends: any }> => {
    const params = new URLSearchParams()
    params.append('period', period)
    if (dateFrom) params.append('date_from', dateFrom)
    if (dateTo) params.append('date_to', dateTo)
    
    const response: AxiosResponse<{ spending_trends: any }> = await api.get(`/reports/trends?${params}`)
    return response.data
  },

  getCategoryAnalysis: async (dateFrom?: string, dateTo?: string): Promise<{ category_analysis: any[] }> => {
    const params = new URLSearchParams()
    if (dateFrom) params.append('date_from', dateFrom)
    if (dateTo) params.append('date_to', dateTo)
    
    const response: AxiosResponse<{ category_analysis: any[] }> = await api.get(`/reports/categories?${params}`)
    return response.data
  },

  getMonthlyReport: async (month?: string, year?: string): Promise<{ monthly_report: any }> => {
    const params = new URLSearchParams()
    if (month) params.append('month', month)
    if (year) params.append('year', year)
    
    const response: AxiosResponse<{ monthly_report: any }> = await api.get(`/reports/monthly?${params}`)
    return response.data
  },
}

export default api
