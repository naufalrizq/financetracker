<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Categories</h1>
        <p class="text-gray-600">Organize your transactions with custom categories</p>
      </div>
      <button
        @click="showAddModal = true"
        class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 flex items-center gap-2"
      >
        <PlusIcon class="w-5 h-5" />
        Add Category
      </button>
    </div>

    <!-- Category Type Tabs -->
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
            All Categories ({{ categories.length }})
          </button>
          <button
            @click="activeTab = 'income'"
            :class="[
              'py-4 px-1 border-b-2 font-medium text-sm',
              activeTab === 'income'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            ]"
          >
            Income ({{ incomeCategories.length }})
          </button>
          <button
            @click="activeTab = 'expense'"
            :class="[
              'py-4 px-1 border-b-2 font-medium text-sm',
              activeTab === 'expense'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            ]"
          >
            Expense ({{ expenseCategories.length }})
          </button>
        </nav>
      </div>

      <!-- Categories Grid -->
      <div v-if="loading" class="p-8 text-center">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
        <p class="mt-2 text-gray-600">Loading categories...</p>
      </div>

      <div v-else-if="filteredCategories.length === 0" class="p-8 text-center text-gray-500">
        <TagIcon class="w-12 h-12 mx-auto mb-4 text-gray-400" />
        <p class="text-lg font-medium mb-2">No categories yet</p>
        <p class="mb-4">Create your first category to organize transactions</p>
        <button
          @click="showAddModal = true"
          class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
        >
          Add Category
        </button>
      </div>

      <div v-else class="p-4">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div
            v-for="category in filteredCategories"
            :key="category.id"
            class="border rounded-lg p-4 hover:shadow-md transition-shadow"
          >
            <div class="flex items-start justify-between mb-3">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full flex items-center justify-center"
                     :style="{ backgroundColor: category.color + '20', color: category.color }">
                  <component :is="getCategoryIcon(category.icon)" class="w-5 h-5" />
                </div>
                <div>
                  <h3 class="font-semibold text-gray-900">{{ category.name }}</h3>
                  <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium"
                        :class="category.type === 'income' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'">
                    {{ category.type }}
                  </span>
                </div>
              </div>
              <div class="flex gap-1">
                <button
                  @click="editCategory(category)"
                  class="text-gray-400 hover:text-blue-600 p-1"
                >
                  <PencilIcon class="w-4 h-4" />
                </button>
                <button
                  @click="deleteCategory(category.id)"
                  class="text-gray-400 hover:text-red-600 p-1"
                >
                  <TrashIcon class="w-4 h-4" />
                </button>
              </div>
            </div>
            
            <p v-if="category.description" class="text-sm text-gray-600 mb-3">
              {{ category.description }}
            </p>
            
            <div class="flex justify-between items-center text-sm text-gray-500">
              <span>{{ getCategoryUsage(category.id) }} transactions</span>
              <span>Created {{ formatDate(category.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Category Modal -->
    <div v-if="showAddModal || showEditModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4">
        <h3 class="text-lg font-semibold mb-4">
          {{ showEditModal ? 'Edit Category' : 'Add Category' }}
        </h3>
        
        <form @submit.prevent="submitCategory" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Category Name</label>
            <input
              v-model="categoryForm.name"
              type="text"
              required
              class="w-full border border-gray-300 rounded-md px-3 py-2"
              placeholder="e.g., Groceries, Salary, Entertainment"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
            <select v-model="categoryForm.type" required class="w-full border border-gray-300 rounded-md px-3 py-2">
              <option value="">Select Type</option>
              <option value="income">Income</option>
              <option value="expense">Expense</option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Icon</label>
            <div class="grid grid-cols-6 gap-2 p-3 border border-gray-300 rounded-md max-h-32 overflow-y-auto">
              <button
                v-for="icon in availableIcons"
                :key="icon.name"
                type="button"
                @click="categoryForm.icon = icon.name"
                :class="[
                  'p-2 rounded-md border-2 flex items-center justify-center hover:bg-gray-50',
                  categoryForm.icon === icon.name ? 'border-blue-500 bg-blue-50' : 'border-gray-200'
                ]"
              >
                <component :is="icon.component" class="w-5 h-5" />
              </button>
            </div>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Color</label>
            <div class="flex gap-2 flex-wrap">
              <button
                v-for="color in availableColors"
                :key="color"
                type="button"
                @click="categoryForm.color = color"
                :class="[
                  'w-8 h-8 rounded-full border-2',
                  categoryForm.color === color ? 'border-gray-800' : 'border-gray-300'
                ]"
                :style="{ backgroundColor: color }"
              ></button>
            </div>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Description (Optional)</label>
            <textarea
              v-model="categoryForm.description"
              rows="3"
              class="w-full border border-gray-300 rounded-md px-3 py-2"
              placeholder="Additional notes about this category"
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { 
  PlusIcon, 
  TagIcon,
  PencilIcon,
  TrashIcon,
  ShoppingCartIcon,
  HomeIcon,
  CarIcon,
  HeartIcon,
  AcademicCapIcon,
  CurrencyDollarIcon,
  GiftIcon,
  FilmIcon,
  PhoneIcon,
  LightBulbIcon,
  BeakerIcon,
  CameraIcon
} from '@heroicons/vue/24/outline'
import { api } from '@/services/api'
import type { Category } from '@/types'

// Reactive data
const loading = ref(false)
const submitting = ref(false)
const showAddModal = ref(false)
const showEditModal = ref(false)
const activeTab = ref('all')
const categories = ref<Category[]>([])
const transactionCounts = ref<Record<string, number>>({})

// Category form
const categoryForm = reactive({
  id: null as string | null,
  name: '',
  type: '',
  icon: 'ShoppingCartIcon',
  color: '#3B82F6',
  description: ''
})

// Available icons and colors
const availableIcons = [
  { name: 'ShoppingCartIcon', component: ShoppingCartIcon },
  { name: 'HomeIcon', component: HomeIcon },
  { name: 'CarIcon', component: CarIcon },
  { name: 'HeartIcon', component: HeartIcon },
  { name: 'AcademicCapIcon', component: AcademicCapIcon },
  { name: 'CurrencyDollarIcon', component: CurrencyDollarIcon },
  { name: 'GiftIcon', component: GiftIcon },
  { name: 'FilmIcon', component: FilmIcon },
  { name: 'PhoneIcon', component: PhoneIcon },
  { name: 'LightBulbIcon', component: LightBulbIcon },
  { name: 'BeakerIcon', component: BeakerIcon },
  { name: 'CameraIcon', component: CameraIcon }
]

const availableColors = [
  '#3B82F6', '#EF4444', '#10B981', '#F59E0B', '#8B5CF6',
  '#EC4899', '#06B6D4', '#84CC16', '#F97316', '#6366F1'
]

// Computed
const filteredCategories = computed(() => {
  if (activeTab.value === 'all') return categories.value
  return categories.value.filter(category => category.type === activeTab.value)
})

const incomeCategories = computed(() => {
  return categories.value.filter(category => category.type === 'income')
})

const expenseCategories = computed(() => {
  return categories.value.filter(category => category.type === 'expense')
})

// Methods
const loadCategories = async () => {
  try {
    loading.value = true
    const response = await api.get('/categories')
    categories.value = response.data
    
    // Load transaction counts for each category
    for (const category of categories.value) {
      try {
        const countResponse = await api.get(`/transactions/count-by-category/${category.id}`)
        transactionCounts.value[category.id] = countResponse.data.count || 0
      } catch (error) {
        transactionCounts.value[category.id] = 0
      }
    }
  } catch (error) {
    console.error('Failed to load categories:', error)
  } finally {
    loading.value = false
  }
}

const submitCategory = async () => {
  try {
    submitting.value = true
    
    const data = {
      name: categoryForm.name,
      type: categoryForm.type,
      icon: categoryForm.icon,
      color: categoryForm.color,
      description: categoryForm.description
    }
    
    if (showEditModal.value && categoryForm.id) {
      await api.put(`/categories/${categoryForm.id}`, data)
    } else {
      await api.post('/categories', data)
    }
    
    await loadCategories()
    closeModal()
  } catch (error) {
    console.error('Failed to save category:', error)
  } finally {
    submitting.value = false
  }
}

const editCategory = (category: Category) => {
  categoryForm.id = category.id
  categoryForm.name = category.name
  categoryForm.type = category.type
  categoryForm.icon = category.icon
  categoryForm.color = category.color
  categoryForm.description = category.description || ''
  showEditModal.value = true
}

const deleteCategory = async (id: string) => {
  const count = transactionCounts.value[id] || 0
  const message = count > 0 
    ? `This category has ${count} transactions. Are you sure you want to delete it?`
    : 'Are you sure you want to delete this category?'
    
  if (!confirm(message)) return
  
  try {
    await api.delete(`/categories/${id}`)
    await loadCategories()
  } catch (error) {
    console.error('Failed to delete category:', error)
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  categoryForm.id = null
  categoryForm.name = ''
  categoryForm.type = ''
  categoryForm.icon = 'ShoppingCartIcon'
  categoryForm.color = '#3B82F6'
  categoryForm.description = ''
}

const getCategoryIcon = (iconName: string) => {
  const icon = availableIcons.find(i => i.name === iconName)
  return icon?.component || ShoppingCartIcon
}

const getCategoryUsage = (categoryId: string) => {
  return transactionCounts.value[categoryId] || 0
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

// Lifecycle
onMounted(() => {
  loadCategories()
})
</script>