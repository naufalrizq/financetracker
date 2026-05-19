import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// Lazy load components
const Login = () => import('@/views/auth/LoginView.vue')
const Register = () => import('@/views/auth/RegisterView.vue')
const Dashboard = () => import('@/views/DashboardView.vue')
const Transactions = () => import('@/views/TransactionsView.vue')
const Accounts = () => import('@/views/AccountsView.vue')
const Categories = () => import('@/views/CategoriesView.vue')
const Budgets = () => import('@/views/BudgetsView.vue')
const Goals = () => import('@/views/GoalsView.vue')
const Reports = () => import('@/views/ReportsView.vue')
const Profile = () => import('@/views/ProfileView.vue')

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/login',
      name: 'login',
      component: Login,
      meta: { requiresGuest: true }
    },
    {
      path: '/register',
      name: 'register',
      component: Register,
      meta: { requiresGuest: true }
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: Dashboard,
      meta: { requiresAuth: true }
    },
    {
      path: '/transactions',
      name: 'transactions',
      component: Transactions,
      meta: { requiresAuth: true }
    },
    {
      path: '/accounts',
      name: 'accounts',
      component: Accounts,
      meta: { requiresAuth: true }
    },
    {
      path: '/categories',
      name: 'categories',
      component: Categories,
      meta: { requiresAuth: true }
    },
    {
      path: '/budgets',
      name: 'budgets',
      component: Budgets,
      meta: { requiresAuth: true }
    },
    {
      path: '/goals',
      name: 'goals',
      component: Goals,
      meta: { requiresAuth: true }
    },
    {
      path: '/reports',
      name: 'reports',
      component: Reports,
      meta: { requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: Profile,
      meta: { requiresAuth: true }
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      redirect: '/dashboard'
    }
  ]
})

// Navigation guards
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // Check if route requires authentication
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
    return
  }
  
  // Check if route requires guest (not authenticated)
  if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next('/dashboard')
    return
  }
  
  // Update last activity
  if (authStore.isAuthenticated) {
    authStore.updateLastActivity()
  }
  
  next()
})

export default router