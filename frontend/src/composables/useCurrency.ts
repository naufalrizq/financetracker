import { useAuthStore } from '@/stores/auth'

/**
 * Shared currency formatting composable.
 * Defaults to IDR (Indonesian Rupiah) when no user currency is set.
 */
export function useCurrency() {
  const authStore = useAuthStore()

  const formatCurrency = (amount: number): string => {
    const currency = authStore.user?.currency || 'IDR'
    const locale = currency === 'IDR' ? 'id-ID' : 'en-US'
    return new Intl.NumberFormat(locale, {
      style: 'currency',
      currency,
      minimumFractionDigits: currency === 'IDR' ? 0 : 2,
    }).format(amount)
  }

  return { formatCurrency }
}
