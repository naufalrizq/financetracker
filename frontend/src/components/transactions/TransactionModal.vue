<template>
  <div class="modal-overlay" @click="$emit('close')">
    <div class="modal-container">
      <div class="modal-content">
        <div class="modal-panel" @click.stop>
          <!-- Header -->
          <div class="px-6 py-4 border-b border-gray-200">
            <h3 class="text-lg font-medium text-gray-900">
              {{ isEditing ? 'Edit Transaction' : 'Add Transaction' }}
            </h3>
          </div>

          <!-- Form -->
          <form @submit.prevent="handleSubmit" class="px-6 py-4 space-y-4">
            <!-- Transaction Type -->
            <div>
              <label class="form-label">Type *</label>
              <div class="grid grid-cols-3 gap-2">
                <button
                  v-for="type in transactionTypes"
                  :key="type.value"
                  type="button"
                  @click="form.type = type.value"
                  :class="[
                    'flex flex-col items-center p-3 rounded-lg border-2 transition-colors',
                    form.type === type.value
                      ? 'border-primary-500 bg-primary-50 text-primary-700'
                      : 'border-gray-300 hover:border-gray-400'
                  ]"
                >
                  <span class="text-2xl mb-1">{{ type.icon }}</span>
                  <span class="text-sm font-medium">{{ type.label }}</span>
                </button>
              </div>
            </div>

            <!-- Account -->
            <div>
              <label class="form-label">Account *</label>
              <select v-model="form.account_id" class="form-select" required>
                <option value="">Select an account</option>
                <option v-for="account in accounts" :key="account.id" :value="account.id">
                  {{ account.icon }} {{ account.name }} ({{ formatCurrency(account.balance) }})
                </option>
              </select>
            </div>

            <!-- Transfer To Account (only for transfers) -->
            <div v-if="form.type === 'transfer'">
              <label class="form-label">Transfer To *</label>
              <select v-model="form.to_account_id" class="form-select" required>
                <option value="">Select destination account</option>
                <option 
                  v-for="account in accounts" 
                  :key="account.id" 
                  :value="account.id"
                  :disabled="account.id === form.account_id"
                >
                  {{ account.icon }} {{ account.name }} ({{ formatCurrency(account.balance) }})
                </option>
              </select>
            </div>

            <!-- Category (not for transfers) -->
            <div v-if="form.type !== 'transfer'">
              <label class="form-label">Category</label>
              <select v-model="form.category_id" class="form-select">
                <option value="">Select a category</option>
                <option 
                  v-for="category in filteredCategories" 
                  :key="category.id" 
                  :value="category.id"
                >
                  {{ category.icon }} {{ category.name }}
                </option>
              </select>
            </div>

            <!-- Amount -->
            <div>
              <label class="form-label">Amount *</label>
              <div class="relative">
                <span class="absolute left-3 top-2 text-gray-500">
                  {{ currencySymbol }}
                </span>
                <input
                  v-model.number="form.amount"
                  type="number"
                  step="0.01"
                  min="0"
                  class="form-input pl-8"
                  placeholder="0.00"
                  required
                />
              </div>
            </div>

            <!-- Description -->
            <div>
              <label class="form-label">Description *</label>
              <input
                v-model="form.description"
                type="text"
                class="form-input"
                placeholder="Enter transaction description"
                required
              />
            </div>

            <!-- Notes -->
            <div>
              <label class="form-label">Notes</label>
              <textarea
                v-model="form.notes"
                class="form-textarea"
                rows="3"
                placeholder="Additional notes (optional)"
              ></textarea>
            </div>

            <!-- Date -->
            <div>
              <label class="form-label">Date *</label>
              <input
                v-model="form.date"
                type="date"
                class="form-input"
                required
              />
            </div>

            <!-- Actions -->
            <div class="flex justify-end space-x-3 pt-4 border-t border-gray-200">
              <button
                type="button"
                @click="$emit('close')"
                class="btn-outline"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="isSubmitting"
                class="btn-primary"
              >
                <div v-if="isSubmitting" class="loading-spinner w-4 h-4 mr-2"></div>
                {{ isSubmitting ? 'Saving...' : (isEditing ? 'Update' : 'Create') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { useToast } from 'vue-toastification'
import { transactionsApi, accountsApi, categoriesApi } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import type { Transaction, CreateTransactionForm } from '@/types'

interface Props {
  transaction?: Transaction | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
  saved: []
}>()

const authStore = useAuthStore()
const toast = useToast()

const isSubmitting = ref(false)
const isEditing = computed(() => !!props.transaction)

const form = reactive<CreateTransactionForm>({
  account_id: '',
  category_id: '',
  type: 'expense' as CreateTransactionForm['type'],
  amount: 0,
  description: '',
  notes: '',
  date: new Date().toISOString().split('T')[0],
  to_account_id: '',
})

const transactionTypes = [
  { value: 'income', label: 'Income', icon: '💰' },
  { value: 'expense', label: 'Expense', icon: '💸' },
  { value: 'transfer', label: 'Transfer', icon: '🔄' },
]

// Queries
const { data: accountsData } = useQuery({
  queryKey: ['accounts'],
  queryFn: () => accountsApi.getAccounts(),
})

const { data: categoriesData } = useQuery({
  queryKey: ['categories'],
  queryFn: () => categoriesApi.getCategories(),
})

const accounts = computed(() => accountsData.value?.accounts || [])
const categories = computed(() => categoriesData.value?.categories || [])

const filteredCategories = computed(() => {
  return categories.value.filter(category => category.type === form.type)
})

const currencySymbol = computed(() => {
  const currency = authStore.user?.currency || 'IDR'
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: currency,
  }).formatToParts(0).find(part => part.type === 'currency')?.value || 'Rp'
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

const handleSubmit = async () => {
  if (isSubmitting.value) return

  try {
    isSubmitting.value = true

    // Prepare form data
    const transactionData = {
      ...form,
      category_id: form.category_id || undefined,
      to_account_id: form.to_account_id || undefined,
    }

    if (isEditing.value && props.transaction) {
      await transactionsApi.updateTransaction(props.transaction.id, transactionData)
      toast.success('Transaction updated successfully')
    } else {
      await transactionsApi.createTransaction(transactionData)
      toast.success('Transaction created successfully')
    }

    emit('saved')
  } catch (error: any) {
    const message = error.response?.data?.error || 'Failed to save transaction'
    toast.error(message)
  } finally {
    isSubmitting.value = false
  }
}

// Initialize form with transaction data if editing
onMounted(() => {
  if (props.transaction) {
    Object.assign(form, {
      account_id: props.transaction.account_id,
      category_id: props.transaction.category_id || '',
      type: props.transaction.type,
      amount: props.transaction.amount,
      description: props.transaction.description,
      notes: props.transaction.notes || '',
      date: props.transaction.date.split('T')[0],
      to_account_id: props.transaction.to_account_id || '',
    })
  }
})
</script>