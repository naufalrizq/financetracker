<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Financial Goals</h1>
        <p class="text-gray-600">Set and track your savings and financial objectives</p>
      </div>
      <button
        @click="showAddModal = true"
        class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 flex items-center gap-2"
      >
        <PlusIcon class="w-5 h-5" />
        Add Goal
      </button>
    </div>

    <!-- Goals Overview -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <div class="bg-gradient-to-r from-blue-500 to-blue-600 text-white p-6 rounded-lg">
        <h3 class="text-lg font-semibold mb-2">Total Goals</h3>
        <p class="text-3xl font-bold">{{ goals.length }}</p>
      </div>
      <div class="bg-gradient-to-r from-green-500 to-green-600 text-white p-6 rounded-lg">
        <h3 class="text-lg font-semibold mb-2">Completed</h3>
        <p class="text-3xl font-bold">{{ completedGoals }}</p>
      </div>
      <div class="bg-gradient-to-r from-yellow-500 to-yellow-600 text-white p-6 rounded-lg">
        <h3 class="text-lg font-semibold mb-2">In Progress</h3>
        <p class="text-3xl font-bold">{{ inProgressGoals }}</p>
      </div>
      <div class="bg-gradient-to-r from-purple-500 to-purple-600 text-white p-6 rounded-lg">
        <h3 class="text-lg font-semibold mb-2">Total Target</h3>
        <p class="text-3xl font-bold">{{ formatCurrency(totalTargetAmount) }}</p>
      </div>
    </div>

    <!-- Goal Status Filter -->
    <div class="bg-white rounded-lg shadow-sm border">
      <div class="border-b">
        <nav class="flex space-x-8 px-4">
          <button
            @click="activeTab = 'all'"
            :class="[
              'py-4 px-1 border-b-2 font-medium text-sm',
              activeTab === 'all'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            ]"
          >
            All Goals ({{ goals.length }})
          </button>
          <button
            @click="activeTab = 'active'"
            :class="[
              'py-4 px-1 border-b-2 font-medium text-sm',
              activeTab === 'active'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            ]"
          >
            Active ({{ inProgressGoals }})
          </button>
          <button
            @click="activeTab = 'completed'"
            :class="[
              'py-4 px-1 border-b-2 font-medium text-sm',
              activeTab === 'completed'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            ]"
          >
            Completed ({{ completedGoals }})
          </button>
        </nav>
      </div>

      <!-- Goals List -->
      <div v-if="loading" class="p-8 text-center">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
        <p class="mt-2 text-gray-600">Loading goals...</p>
      </div>

      <div v-else-if="filteredGoals.length === 0" class="p-8 text-center text-gray-500">
        <FlagIcon class="w-12 h-12 mx-auto mb-4 text-gray-400" />
        <p class="text-lg font-medium mb-2">No goals yet</p>
        <p class="mb-4">Set your first financial goal to start saving</p>
        <button
          @click="showAddModal = true"
          class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
        >
          Add Goal
        </button>
      </div>

      <div v-else class="p-4">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div
            v-for="goal in filteredGoals"
            :key="goal.id"
            class="border rounded-lg p-6 hover:shadow-md transition-shadow"
          >
            <!-- Goal Header -->
            <div class="flex items-start justify-between mb-4">
              <div class="flex items-center gap-3">
                <div class="w-12 h-12 rounded-full flex items-center justify-center"
                     :class="getGoalStatusColor(goal)">
                  <component :is="getGoalIcon(goal.type)" class="w-6 h-6" />
                </div>
                <div>
                  <h3 class="font-semibold text-gray-900">{{ goal.name }}</h3>
                  <p class="text-sm text-gray-600 capitalize">{{ goal.type.replace('_', ' ') }}</p>
                </div>
              </div>
              <div class="flex gap-1">
                <button
                  @click="editGoal(goal)"
                  class="text-gray-400 hover:text-blue-600 p-1"
                >
                  <PencilIcon class="w-4 h-4" />
                </button>
                <button
                  @click="deleteGoal(goal.id)"
                  class="text-gray-400 hover:text-red-600 p-1"
                >
                  <TrashIcon class="w-4 h-4" />
                </button>
              </div>
            </div>

            <!-- Goal Progress -->
            <div class="mb-4">
              <div class="flex justify-between text-sm text-gray-600 mb-2">
                <span>Progress</span>
                <span>{{ formatCurrency(goal.current_amount) }} / {{ formatCurrency(goal.target_amount) }}</span>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-3">
                <div
                  class="h-3 rounded-full transition-all duration-300"
                  :class="getGoalProgressColor(goal)"
                  :style="{ width: Math.min(getGoalProgress(goal), 100) + '%' }"
                ></div>
              </div>
              <div class="flex justify-between text-sm mt-2">
                <span class="text-gray-600">{{ getGoalProgress(goal).toFixed(1) }}% complete</span>
                <span
                  :class="isGoalCompleted(goal) ? 'text-green-600 font-medium' : 'text-gray-600'"
                >
                  {{ isGoalCompleted(goal) ? 'Completed!' : `${formatCurrency(goal.target_amount - goal.current_amount)} to go` }}
                </span>
              </div>
            </div>

            <!-- Goal Details -->
            <div class="space-y-2 text-sm text-gray-600">
              <div class="flex justify-between">
                <span>Target Date:</span>
                <span>{{ formatDate(goal.target_date) }}</span>
              </div>
              <div class="flex justify-between">
                <span>Days Remaining:</span>
                <span :class="getDaysRemainingColor(goal)">
                  {{ getDaysRemaining(goal) }}
                </span>
              </div>
              <div v-if="goal.description" class="pt-2 border-t">
                <p class="text-gray-700">{{ goal.description }}</p>
              </div>
            </div>

            <!-- Quick Actions -->
            <div class="flex gap-2 mt-4 pt-4 border-t">
              <button
                @click="addContribution(goal)"
                class="flex-1 bg-green-600 text-white py-2 px-3 rounded-md hover:bg-green-700 text-sm"
                :disabled="isGoalCompleted(goal)"
              >
                Add Money
              </button>
              <button
                @click="viewGoalHistory(goal)"
                class="flex-1 bg-gray-600 text-white py-2 px-3 rounded-md hover:bg-gray-700 text-sm"
              >
                View History
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Goal Modal -->
    <div v-if="showAddModal || showEditModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4">
        <h3 class="text-lg font-semibold mb-4">
          {{ showEditModal ? 'Edit Goal' : 'Add Goal' }}
        </h3>
        
        <form @submit.prevent="submitGoal" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Goal Name</label>
            <input
              v-model="goalForm.name"
              type="text"
              required
              class="w-full border border-gray-300 rounded-md px-3 py-2"
              placeholder="e.g., Emergency Fund, Vacation, New Car"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Goal Type</label>
            <select v-model="goalForm.type" required class="w-full border border-gray-300 rounded-md px-3 py-2">
              <option value="">Select Type</option>
              <option value="emergency_fund">Emergency Fund</option>
              <option value="vacation">Vacation</option>
              <option value="house_down_payment">House Down Payment</option>
              <option value="car_purchase">Car Purchase</option>
              <option value="education">Education</option>
              <option value="retirement">Retirement</option>
              <option value="debt_payoff">Debt Payoff</option>
              <option value="other">Other</option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Target Amount</label>
            <input
              v-model.number="goalForm.target_amount"
              type="number"
              step="0.01"
              min="0"
              required
              class="w-full border border-gray-300 rounded-md px-3 py-2"
              placeholder="0.00"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Current Amount</label>
            <input
              v-model.number="goalForm.current_amount"
              type="number"
              step="0.01"
              min="0"
              class="w-full border border-gray-300 rounded-md px-3 py-2"
              placeholder="0.00"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Target Date</label>
            <input
              v-model="goalForm.target_date"
              type="date"
              required
              class="w-full border border-gray-300 rounded-md px-3 py-2"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Description (Optional)</label>
            <textarea
              v-model="goalForm.description"
              rows="3"
              class="w-full border border-gray-300 rounded-md px-3 py-2"
              placeholder="Additional details about this goal"
            ></textarea>
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

    <!-- Add Contribution Modal -->
    <div v-if="showContributionModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4">
        <h3 class="text-lg font-semibold mb-4">Add Contribution</h3>
        <p class="text-gray-600 mb-4">Add money to: {{ selectedGoal?.name }}</p>
        
        <form @submit.prevent="submitContribution" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Amount</label>
            <input
              v-model.number="contributionAmount"
              type="number"
              step="0.01"
              min="0"
              required
              class="w-full border border-gray-300 rounded-md px-3 py-2"
              placeholder="0.00"
            />
          </div>
          
          <div class="flex gap-3 pt-4">
            <button
              type="submit"
              :disabled="submitting"
              class="flex-1 bg-green-600 text-white py-2 px-4 rounded-md hover:bg-green-700 disabled:opacity-50"
            >
              {{ submitting ? 'Adding...' : 'Add Contribution' }}
            </button>
            <button
              type="button"
              @click="showContributionModal = false"
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
  FlagIcon,
  PencilIcon,
  TrashIcon,
  HomeIcon,
  CarIcon,
  AcademicCapIcon,
  HeartIcon,
  BanknotesIcon,
  ShieldCheckIcon,
  DocumentTextIcon
} from '@heroicons/vue/24/outline'
import { api } from '@/services/api'
import { useCurrency } from '@/composables/useCurrency'
import type { Goal } from '@/types'

const { formatCurrency } = useCurrency()

// Reactive data
const loading = ref(false)
const submitting = ref(false)
const showAddModal = ref(false)
const showEditModal = ref(false)
const showContributionModal = ref(false)
const activeTab = ref('all')
const goals = ref<Goal[]>([])
const selectedGoal = ref<Goal | null>(null)
const contributionAmount = ref(0)

// Goal form
const goalForm = reactive({
  id: null as number | null,
  name: '',
  type: '',
  target_amount: 0,
  current_amount: 0,
  target_date: '',
  description: ''
})

// Computed
const filteredGoals = computed(() => {
  if (activeTab.value === 'all') return goals.value
  if (activeTab.value === 'completed') return goals.value.filter(goal => isGoalCompleted(goal))
  if (activeTab.value === 'active') return goals.value.filter(goal => !isGoalCompleted(goal))
  return goals.value
})

const completedGoals = computed(() => {
  return goals.value.filter(goal => isGoalCompleted(goal)).length
})

const inProgressGoals = computed(() => {
  return goals.value.filter(goal => !isGoalCompleted(goal)).length
})

const totalTargetAmount = computed(() => {
  return goals.value.reduce((sum, goal) => sum + goal.target_amount, 0)
})

// Methods
const loadGoals = async () => {
  try {
    loading.value = true
    const response = await api.get('/goals')
    goals.value = response.data
  } catch (error) {
    console.error('Failed to load goals:', error)
  } finally {
    loading.value = false
  }
}

const submitGoal = async () => {
  try {
    submitting.value = true
    
    const data = {
      name: goalForm.name,
      type: goalForm.type,
      target_amount: goalForm.target_amount,
      current_amount: goalForm.current_amount,
      target_date: goalForm.target_date,
      description: goalForm.description
    }
    
    if (showEditModal.value && goalForm.id) {
      await api.put(`/goals/${goalForm.id}`, data)
    } else {
      await api.post('/goals', data)
    }
    
    await loadGoals()
    closeModal()
  } catch (error) {
    console.error('Failed to save goal:', error)
  } finally {
    submitting.value = false
  }
}

const editGoal = (goal: Goal) => {
  goalForm.id = goal.id
  goalForm.name = goal.name
  goalForm.type = goal.type
  goalForm.target_amount = goal.target_amount
  goalForm.current_amount = goal.current_amount
  goalForm.target_date = goal.target_date.split('T')[0]
  goalForm.description = goal.description || ''
  showEditModal.value = true
}

const deleteGoal = async (id: number) => {
  if (!confirm('Are you sure you want to delete this goal?')) return
  
  try {
    await api.delete(`/goals/${id}`)
    await loadGoals()
  } catch (error) {
    console.error('Failed to delete goal:', error)
  }
}

const addContribution = (goal: Goal) => {
  selectedGoal.value = goal
  contributionAmount.value = 0
  showContributionModal.value = true
}

const submitContribution = async () => {
  if (!selectedGoal.value) return
  
  try {
    submitting.value = true
    
    const newAmount = selectedGoal.value.current_amount + contributionAmount.value
    await api.put(`/goals/${selectedGoal.value.id}`, {
      ...selectedGoal.value,
      current_amount: Math.min(newAmount, selectedGoal.value.target_amount)
    })
    
    await loadGoals()
    showContributionModal.value = false
    selectedGoal.value = null
    contributionAmount.value = 0
  } catch (error) {
    console.error('Failed to add contribution:', error)
  } finally {
    submitting.value = false
  }
}

const viewGoalHistory = (goal: Goal) => {
  // This would open a modal or navigate to a detailed view
  console.log('View history for goal:', goal.name)
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  goalForm.id = null
  goalForm.name = ''
  goalForm.type = ''
  goalForm.target_amount = 0
  goalForm.current_amount = 0
  goalForm.target_date = ''
  goalForm.description = ''
}

const isGoalCompleted = (goal: Goal) => {
  return goal.current_amount >= goal.target_amount
}

const getGoalProgress = (goal: Goal) => {
  return goal.target_amount > 0 ? (goal.current_amount / goal.target_amount) * 100 : 0
}

const getGoalProgressColor = (goal: Goal) => {
  if (isGoalCompleted(goal)) return 'bg-green-500'
  const progress = getGoalProgress(goal)
  if (progress > 75) return 'bg-blue-500'
  if (progress > 50) return 'bg-yellow-500'
  return 'bg-gray-400'
}

const getGoalStatusColor = (goal: Goal) => {
  if (isGoalCompleted(goal)) return 'bg-green-100 text-green-600'
  const progress = getGoalProgress(goal)
  if (progress > 75) return 'bg-blue-100 text-blue-600'
  if (progress > 50) return 'bg-yellow-100 text-yellow-600'
  return 'bg-gray-100 text-gray-600'
}

const getGoalIcon = (type: string) => {
  const icons = {
    emergency_fund: ShieldCheckIcon,
    vacation: HeartIcon,
    house_down_payment: HomeIcon,
    car_purchase: CarIcon,
    education: AcademicCapIcon,
    retirement: BanknotesIcon,
    debt_payoff: DocumentTextIcon,
    other: FlagIcon
  }
  return icons[type as keyof typeof icons] || FlagIcon
}

const getDaysRemaining = (goal: Goal) => {
  const today = new Date()
  const targetDate = new Date(goal.target_date)
  const diffTime = targetDate.getTime() - today.getTime()
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays < 0) return 'Overdue'
  if (diffDays === 0) return 'Today'
  if (diffDays === 1) return '1 day'
  return `${diffDays} days`
}

const getDaysRemainingColor = (goal: Goal) => {
  const today = new Date()
  const targetDate = new Date(goal.target_date)
  const diffTime = targetDate.getTime() - today.getTime()
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays < 0) return 'text-red-600 font-medium'
  if (diffDays <= 30) return 'text-yellow-600 font-medium'
  return 'text-gray-600'
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

// Lifecycle
onMounted(() => {
  loadGoals()
  
  // Set default target date to 1 year from now
  const nextYear = new Date()
  nextYear.setFullYear(nextYear.getFullYear() + 1)
  goalForm.target_date = nextYear.toISOString().split('T')[0]
})
</script>