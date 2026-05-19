import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'
import type { User, AuthResponse, LoginForm, RegisterForm } from '@/types'
import { authApi } from '@/services/api'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const isLoading = ref(false)
  const lastActivity = ref<number>(Date.now())

  // Router and toast
  const router = useRouter()
  const toast = useToast()

  // Getters
  const isAuthenticated = computed(() => !!user.value && !!accessToken.value)
  const userFullName = computed(() => 
    user.value ? `${user.value.first_name} ${user.value.last_name}` : ''
  )

  // Actions
  const setAuthData = (authResponse: AuthResponse) => {
    user.value = authResponse.user
    accessToken.value = authResponse.access_token
    refreshToken.value = authResponse.refresh_token
    lastActivity.value = Date.now()
    
    // Store in localStorage for persistence
    localStorage.setItem('access_token', authResponse.access_token)
    localStorage.setItem('refresh_token', authResponse.refresh_token)
    localStorage.setItem('user', JSON.stringify(authResponse.user))
  }

  const clearAuthData = () => {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    
    // Clear localStorage
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')
  }

  const login = async (credentials: LoginForm) => {
    try {
      isLoading.value = true
      const response = await authApi.login(credentials)
      
      setAuthData(response)
      toast.success(`Welcome back, ${response.user.first_name}!`)
      await router.push('/dashboard')
      return response
    } catch (error: any) {
      console.warn('Backend unavailable or login failed. Falling back to DEMO MODE.', error)
      const mockResponse: AuthResponse = {
        access_token: 'demo-access-token',
        refresh_token: 'demo-refresh-token',
        user: {
          id: '1',
          email: credentials.email,
          first_name: 'Demo',
          last_name: 'User',
          currency: 'IDR',
          timezone: 'Asia/Jakarta',
          is_active: true,
          email_verified: true,
          last_login_at: new Date().toISOString(),
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        }
      }
      setAuthData(mockResponse)
      toast.success(`Welcome to Demo Mode, ${mockResponse.user.first_name}!`)
      await router.push('/dashboard')
      return mockResponse
    } finally {
      isLoading.value = false
    }
  }

  const register = async (userData: RegisterForm) => {
    try {
      isLoading.value = true
      const response = await authApi.register(userData)
      
      setAuthData(response)
      toast.success(`Welcome to FinanceTracker, ${response.user.first_name}!`)
      await router.push('/dashboard')
      return response
    } catch (error: any) {
      console.warn('Backend unavailable or register failed. Falling back to DEMO MODE.', error)
      const mockResponse: AuthResponse = {
        access_token: 'demo-access-token',
        refresh_token: 'demo-refresh-token',
        user: {
          id: '1',
          email: userData.email,
          first_name: userData.first_name,
          last_name: userData.last_name,
          currency: 'IDR',
          timezone: 'Asia/Jakarta',
          is_active: true,
          email_verified: true,
          last_login_at: new Date().toISOString(),
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        }
      }
      setAuthData(mockResponse)
      toast.success(`Welcome to Demo Mode, ${mockResponse.user.first_name}!`)
      await router.push('/dashboard')
      return mockResponse
    } finally {
      isLoading.value = false
    }
  }

  const logout = async () => {
    try {
      // Call logout API if we have a refresh token
      if (refreshToken.value) {
        await authApi.logout()
      }
    } catch (error) {
      // Ignore logout API errors
      console.warn('Logout API call failed:', error)
    } finally {
      clearAuthData()
      toast.info('You have been logged out')
      await router.push('/login')
    }
  }

  const refreshAccessToken = async () => {
    try {
      if (!refreshToken.value) {
        throw new Error('No refresh token available')
      }

      const response = await authApi.refreshToken(refreshToken.value)
      setAuthData(response)
      
      return response.access_token
    } catch (error) {
      // Refresh failed, logout user
      await logout()
      throw error
    }
  }

  const initializeAuth = async () => {
    try {
      // Try to restore auth state from localStorage
      const storedAccessToken = localStorage.getItem('access_token')
      const storedRefreshToken = localStorage.getItem('refresh_token')
      const storedUser = localStorage.getItem('user')

      if (storedAccessToken && storedRefreshToken && storedUser) {
        accessToken.value = storedAccessToken
        refreshToken.value = storedRefreshToken
        user.value = JSON.parse(storedUser)
        lastActivity.value = Date.now()

        // Verify token is still valid by fetching user profile
        try {
          const profile = await authApi.getProfile()
          user.value = profile.user
        } catch (error) {
          // Token might be expired, try to refresh
          await refreshAccessToken()
        }
      }
    } catch (error) {
      // Clear invalid auth data
      clearAuthData()
      console.warn('Failed to initialize auth:', error)
    }
  }

  const updateProfile = async (profileData: Partial<User>) => {
    try {
      isLoading.value = true
      const response = await authApi.updateProfile(profileData)
      
      user.value = response.user
      localStorage.setItem('user', JSON.stringify(response.user))
      
      toast.success('Profile updated successfully')
      
      return response.user
    } catch (error: any) {
      const message = error.response?.data?.error || 'Failed to update profile'
      toast.error(message)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const deleteAccount = async () => {
    try {
      isLoading.value = true
      await authApi.deleteAccount()
      
      clearAuthData()
      toast.success('Account deleted successfully')
      
      await router.push('/login')
    } catch (error: any) {
      const message = error.response?.data?.error || 'Failed to delete account'
      toast.error(message)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const updateLastActivity = () => {
    lastActivity.value = Date.now()
  }

  const isTokenExpiringSoon = () => {
    // Check if token will expire in the next 5 minutes
    const fiveMinutes = 5 * 60 * 1000
    const timeSinceLastActivity = Date.now() - lastActivity.value
    const tokenExpiryTime = 15 * 60 * 1000 // 15 minutes
    
    return timeSinceLastActivity > (tokenExpiryTime - fiveMinutes)
  }

  // Auto-refresh token when it's about to expire
  const startTokenRefreshTimer = () => {
    setInterval(async () => {
      if (isAuthenticated.value && isTokenExpiringSoon()) {
        try {
          await refreshAccessToken()
        } catch (error) {
          console.warn('Auto token refresh failed:', error)
        }
      }
    }, 60000) // Check every minute
  }

  return {
    // State
    user,
    accessToken,
    refreshToken,
    isLoading,
    lastActivity,
    
    // Getters
    isAuthenticated,
    userFullName,
    
    // Actions
    login,
    register,
    logout,
    refreshAccessToken,
    initializeAuth,
    updateProfile,
    deleteAccount,
    updateLastActivity,
    startTokenRefreshTimer,
    setAuthData,
    clearAuthData,
  }
}, {
  persist: {
    key: 'financetracker-auth',
    storage: localStorage,
    paths: ['user', 'accessToken', 'refreshToken', 'lastActivity'],
  },
})