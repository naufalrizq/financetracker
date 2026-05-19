// User types
export interface User {
  id: string
  email: string
  first_name: string
  last_name: string
  currency: string
  timezone: string
  is_active: boolean
  email_verified: boolean
  last_login_at: string | null
  created_at: string
  updated_at: string
}

export interface AuthResponse {
  user: User
  access_token: string
  refresh_token: string
  expires_in: number
}

// Category types
export type CategoryType = 'income' | 'expense'

export interface Category {
  id: string
  user_id: string
  name: string
  type: CategoryType
  color: string
  icon: string
  description: string
  is_default: boolean
  is_active: boolean
  created_at: string
  updated_at: string
}

// Account types
export type AccountType = 'checking' | 'savings' | 'credit' | 'cash' | 'investment'

export interface Account {
  id: string
  user_id: string
  name: string
  type: AccountType
  balance: number
  currency: string
  color: string
  icon: string
  description: string
  is_active: boolean
  include_in_total: boolean
  created_at: string
  updated_at: string
}

export interface AccountSummary extends Account {
  total_income: number
  total_expense: number
  transaction_count: number
  last_transaction_at: string | null
}

// Transaction types
export type TransactionType = 'income' | 'expense' | 'transfer'

export interface Transaction {
  id: string
  user_id: string
  account_id: string
  category_id: string | null
  type: TransactionType
  amount: number
  currency: string
  description: string
  notes: string
  date: string
  to_account_id: string | null
  is_recurring: boolean
  recurring_type: string
  recurring_until: string | null
  tags: string
  created_at: string
  updated_at: string
  
  // Relationships
  category?: Category
  account?: Account
  to_account?: Account
}

export interface TransactionFilter {
  account_id?: string
  category_id?: string
  type?: TransactionType
  date_from?: string
  date_to?: string
  amount_min?: number
  amount_max?: number
  search?: string
  page?: number
  limit?: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

export interface PaginationInfo {
  page: number
  limit: number
  total: number
  total_pages: number
}

export interface TransactionResponse {
  transactions: Transaction[]
  pagination: PaginationInfo
}

// Budget types
export type BudgetPeriod = 'weekly' | 'monthly' | 'yearly'

export interface Budget {
  id: string
  user_id: string
  category_id: string | null
  name: string
  amount: number
  currency: string
  period: BudgetPeriod
  start_date: string
  end_date: string | null
  is_active: boolean
  alert_threshold: number
  alert_enabled: boolean
  last_alert_sent: string | null
  created_at: string
  updated_at: string
  
  // Relationships
  category?: Category
}

export interface BudgetStatus extends Budget {
  spent: number
  remaining: number
  percentage_used: number
  is_over_budget: boolean
  days_remaining: number
  category_name?: string
  period_start: string
  period_end: string
}

// Goal types
export type GoalType = 'saving' | 'debt_payoff' | 'investment' | 'other'

export interface Goal {
  id: string
  user_id: string
  name: string
  description: string
  type: GoalType
  target_amount: number
  current_amount: number
  currency: string
  target_date: string | null
  is_completed: boolean
  completed_at: string | null
  color: string
  icon: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface GoalStatus extends Goal {
  progress: number
  remaining_amount: number
  days_remaining: number | null
  required_monthly_amount: number | null
  is_on_track: boolean
  status: 'completed' | 'on_track' | 'behind' | 'overdue' | 'in_progress'
}

// Report types
export interface FinancialSummary {
  total_income: number
  total_expense: number
  net_income: number
  total_accounts: number
  total_transactions: number
  period_start: string
  period_end: string
  currency: string
}

export interface ExpenseAnalysis {
  total_expense: number
  categories: Array<{
    category_id: string
    category_name: string
    amount: number
    percentage: number
    color: string
  }>
  period_start: string
  period_end: string
}

export interface IncomeAnalysis {
  total_income: number
  categories: Array<{
    category_id: string
    category_name: string
    amount: number
    percentage: number
    color: string
  }>
  period_start: string
  period_end: string
}

export interface SpendingTrend {
  date: string
  income: number
  expense: number
  net: number
}

export interface SpendingTrends {
  trends: SpendingTrend[]
  period: 'daily' | 'weekly' | 'monthly'
  period_start: string
  period_end: string
}

// Form types
export interface LoginForm {
  email: string
  password: string
}

export interface RegisterForm {
  email: string
  password: string
  first_name: string
  last_name: string
  currency?: string
  timezone?: string
}

export interface CreateTransactionForm {
  account_id: string
  category_id?: string
  type: TransactionType
  amount: number
  currency?: string
  description: string
  notes?: string
  date: string
  to_account_id?: string
  tags?: string[]
}

export interface CreateAccountForm {
  name: string
  type: AccountType
  balance?: number
  currency?: string
  color?: string
  icon?: string
  description?: string
  include_in_total?: boolean
}

export interface CreateCategoryForm {
  name: string
  type: CategoryType
  color?: string
  icon?: string
  description?: string
}

export interface CreateBudgetForm {
  category_id?: string
  name: string
  amount: number
  currency?: string
  period: BudgetPeriod
  start_date: string
  end_date?: string
  alert_threshold?: number
  alert_enabled?: boolean
}

export interface CreateGoalForm {
  name: string
  description?: string
  type: GoalType
  target_amount: number
  current_amount?: number
  currency?: string
  target_date?: string
  color?: string
  icon?: string
}

// API Error types
export interface ApiError {
  error: string
  details?: string | string[]
  status?: number
}

// Theme types
export type Theme = 'light' | 'dark' | 'system'

// Chart data types
export interface ChartData {
  labels: string[]
  datasets: Array<{
    label: string
    data: number[]
    backgroundColor?: string | string[]
    borderColor?: string | string[]
    borderWidth?: number
  }>
}

// Utility types
export type LoadingState = 'idle' | 'loading' | 'success' | 'error'

export interface SelectOption {
  value: string
  label: string
  icon?: string
  color?: string
}