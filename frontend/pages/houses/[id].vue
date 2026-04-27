<template>
  <div class="min-h-screen bg-gray-50" v-if="house">
    <div class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <NuxtLink to="/" class="text-blue-600 hover:text-blue-700">
          ← 返回列表
        </NuxtLink>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="lg:col-span-2 space-y-6">
          <div class="bg-white rounded-lg shadow-sm overflow-hidden">
            <div class="relative h-96 bg-gray-200">
              <div v-if="house.images && house.images.length > 0" class="w-full h-full">
                <img
                  :src="getImageUrl(house.images[currentImageIndex])"
                  :alt="house.title"
                  class="w-full h-full object-contain"
                />
              </div>

              <button
                v-if="house.images && house.images.length > 1"
                @click="prevImage"
                class="absolute left-4 top-1/2 -translate-y-1/2 bg-black/50 text-white w-10 h-10 rounded-full flex items-center justify-center hover:bg-black/70"
              >
                ‹
              </button>
              <button
                v-if="house.images && house.images.length > 1"
                @click="nextImage"
                class="absolute right-4 top-1/2 -translate-y-1/2 bg-black/50 text-white w-10 h-10 rounded-full flex items-center justify-center hover:bg-black/70"
              >
                ›
              </button>

              <div v-if="house.images && house.images.length > 0" class="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2">
                <button
                  v-for="(_, index) in house.images"
                  :key="index"
                  @click="currentImageIndex = index"
                  :class="[
                    'w-2 h-2 rounded-full',
                    currentImageIndex === index ? 'bg-white' : 'bg-white/50'
                  ]"
                />
              </div>
            </div>

            <div v-if="house.images && house.images.length > 0" class="p-4 flex gap-2 overflow-x-auto">
              <button
                v-for="(img, index) in house.images"
                :key="index"
                @click="currentImageIndex = index"
                :class="[
                  'w-20 h-20 flex-shrink-0 rounded overflow-hidden border-2',
                  currentImageIndex === index ? 'border-blue-500' : 'border-transparent'
                ]"
              >
                <img :src="getImageUrl(img)" :alt="`图片 ${index + 1}`" class="w-full h-full object-cover" />
              </button>
            </div>
          </div>

          <div class="bg-white rounded-lg shadow-sm p-6">
            <h1 class="text-2xl font-bold text-gray-900">{{ house.title }}</h1>
            <div class="mt-2 flex items-center gap-4">
              <span class="text-3xl font-bold text-red-500">¥{{ house.price }}</span>
              <span class="text-gray-500">/月</span>
              <span class="text-gray-500">|</span>
              <span class="text-gray-600">{{ house.bedrooms }}室{{ house.living_rooms }}厅{{ house.bathrooms }}卫</span>
              <span class="text-gray-500">|</span>
              <span class="text-gray-600">{{ house.area }}㎡</span>
            </div>

            <div class="mt-4 grid grid-cols-2 md:grid-cols-4 gap-4">
              <div class="border rounded-lg p-3">
                <p class="text-gray-500 text-sm">小区</p>
                <p class="font-medium">{{ house.community_name }}</p>
              </div>
              <div class="border rounded-lg p-3">
                <p class="text-gray-500 text-sm">楼层</p>
                <p class="font-medium">{{ house.floor }}/{{ house.total_floors }}层</p>
              </div>
              <div class="border rounded-lg p-3">
                <p class="text-gray-500 text-sm">朝向</p>
                <p class="font-medium">{{ getOrientationLabel(house.orientation) }}</p>
              </div>
              <div class="border rounded-lg p-3">
                <p class="text-gray-500 text-sm">装修</p>
                <p class="font-medium">{{ getDecorationLabel(house.decoration) }}</p>
              </div>
            </div>
          </div>

          <div class="bg-white rounded-lg shadow-sm p-6">
            <h2 class="text-lg font-semibold mb-4">配套设施</h2>
            <div class="flex flex-wrap gap-3">
              <template v-if="house.facilities && house.facilities.length > 0">
                <span
                  v-for="facility in house.facilities"
                  :key="facility"
                  class="px-4 py-2 bg-blue-50 text-blue-700 rounded-full text-sm"
                >
                  {{ getFacilityLabel(facility) }}
                </span>
              </template>
              <span v-else class="text-gray-500">暂无配套设施信息</span>
            </div>
          </div>

          <div class="bg-white rounded-lg shadow-sm p-6">
            <h2 class="text-lg font-semibold mb-4">房源描述</h2>
            <p class="text-gray-600 whitespace-pre-line">{{ house.description || '暂无描述' }}</p>
          </div>

          <div v-if="house.latitude && house.longitude" class="bg-white rounded-lg shadow-sm p-6">
            <h2 class="text-lg font-semibold mb-4">位置信息</h2>
            <p class="text-gray-600 mb-4">{{ house.address }}</p>
            <div ref="mapContainer" class="h-64 rounded-lg overflow-hidden"></div>
          </div>

          <div v-if="similarHouses.length > 0" class="bg-white rounded-lg shadow-sm p-6">
            <h2 class="text-lg font-semibold mb-4">相似房源推荐</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <NuxtLink
                v-for="h in similarHouses"
                :key="h.id"
                :to="`/houses/${h.id}`"
                class="flex gap-4 p-3 border rounded-lg hover:shadow-md transition-shadow"
              >
                <div class="w-32 h-24 flex-shrink-0 bg-gray-200 rounded overflow-hidden">
                  <img
                    v-if="h.images && h.images.length > 0"
                    :src="getImageUrl(h.images[0])"
                    :alt="h.title"
                    class="w-full h-full object-cover"
                  />
                </div>
                <div class="flex-1 min-w-0">
                  <h3 class="font-medium text-gray-900 truncate">{{ h.title }}</h3>
                  <p class="text-red-500 font-bold mt-1">¥{{ h.price }}/月</p>
                  <p class="text-gray-500 text-sm">{{ h.bedrooms }}室{{ h.living_rooms }}厅 | {{ h.area }}㎡</p>
                </div>
              </NuxtLink>
            </div>
          </div>
        </div>

        <div class="space-y-6">
          <div class="bg-white rounded-lg shadow-sm p-6 sticky top-6">
            <h2 class="text-lg font-semibold mb-4">房东信息</h2>
            <div class="flex items-center gap-4">
              <div class="w-16 h-16 bg-gray-200 rounded-full flex items-center justify-center">
                <span class="text-2xl text-gray-500">
                  {{ house.landlord?.nickname?.charAt(0) || '房' }}
                </span>
              </div>
              <div>
                <p class="font-medium">{{ house.landlord?.nickname || '房东' }}</p>
                <p class="text-sm text-gray-500">{{ house.landlord?.phone }}</p>
              </div>
            </div>

            <div class="mt-6 space-y-3">
              <button
                v-if="authStore.isTenant"
                @click="toggleFavorite"
                :class="[
                  'w-full py-3 rounded-lg font-medium',
                  isFavorited
                    ? 'bg-red-50 text-red-600 border border-red-200'
                    : 'bg-white border border-gray-300 hover:bg-gray-50'
                ]"
              >
                {{ isFavorited ? '♥ 已收藏' : '♡ 收藏房源' }}
              </button>

              <button
                v-if="authStore.isTenant"
                @click="showMessageModal = true"
                class="w-full py-3 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700"
              >
                在线咨询
              </button>

              <NuxtLink
                v-if="!authStore.isAuthenticated"
                to="/login"
                class="block w-full py-3 text-center bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700"
              >
                登录后咨询
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showMessageModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg w-full max-w-lg mx-4">
        <div class="p-4 border-b flex justify-between items-center">
          <h3 class="font-semibold">给房东发消息</h3>
          <button @click="showMessageModal = false" class="text-gray-400 hover:text-gray-600">✕</button>
        </div>
        <div class="p-4">
          <textarea
            v-model="messageContent"
            placeholder="请输入您想咨询的问题..."
            class="w-full h-32 border rounded-lg p-3 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div class="p-4 border-t flex justify-end gap-3">
          <button
            @click="showMessageModal = false"
            class="px-4 py-2 border rounded-lg hover:bg-gray-50"
          >
            取消
          </button>
          <button
            @click="sendMessage"
            :disabled="sendingMessage || !messageContent.trim()"
            class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
          >
            {{ sendingMessage ? '发送中...' : '发送' }}
          </button>
        </div>
      </div>
    </div>
  </div>

  <div v-if="loading" class="flex items-center justify-center min-h-screen">
    <p class="text-gray-500">加载中...</p>
  </div>
</template>

<script setup lang="ts">
import type L from 'leaflet'

interface House {
  id: string
  title: string
  community_name: string
  address: string
  longitude: number
  latitude: number
  price: number
  area: number
  bedrooms: number
  living_rooms: number
  bathrooms: number
  floor: number
  total_floors: number
  orientation: string
  decoration: string
  facilities: string[]
  description: string
  images: string[]
  landlord: {
    id: string
    nickname: string
    phone: string
  }
}

const route = useRoute()
const authStore = useAuthStore()
const config = useRuntimeConfig()
const api = useApi()

const house = ref<House | null>(null)
const similarHouses = ref<House[]>([])
const loading = ref(true)
const currentImageIndex = ref(0)
const isFavorited = ref(false)
const showMessageModal = ref(false)
const messageContent = ref('')
const sendingMessage = ref(false)

const mapContainer = ref<HTMLElement>()
let map: L.Map | null = null

const getImageUrl = (path: string) => {
  if (path.startsWith('http')) return path
  return `${config.public.apiUrl}${path}`
}

const getOrientationLabel = (orientation: string): string => {
  const labels: Record<string, string> = {
    north: '北',
    south: '南',
    east: '东',
    west: '西',
    northeast: '东北',
    northwest: '西北',
    southeast: '东南',
    southwest: '西南'
  }
  return labels[orientation] || orientation
}

const getDecorationLabel = (decoration: string): string => {
  const labels: Record<string, string> = {
    rough: '毛坯',
    simple: '简装',
    standard: '精装',
    luxury: '豪装'
  }
  return labels[decoration] || decoration
}

const getFacilityLabel = (facility: string): string => {
  const labels: Record<string, string> = {
    wifi: 'WiFi',
    air_conditioner: '空调',
    washing_machine: '洗衣机',
    refrigerator: '冰箱',
    water_heater: '热水器',
    private_bathroom: '独卫',
    balcony: '阳台',
    tv: '电视',
    wardrobe: '衣柜',
    sofa: '沙发',
    bed: '床',
    desk: '书桌',
    chair: '椅子'
  }
  return labels[facility] || facility
}

const prevImage = () => {
  if (house.value && house.value.images) {
    currentImageIndex.value = (currentImageIndex.value - 1 + house.value.images.length) % house.value.images.length
  }
}

const nextImage = () => {
  if (house.value && house.value.images) {
    currentImageIndex.value = (currentImageIndex.value + 1) % house.value.images.length
  }
}

const fetchHouse = async () => {
  loading.value = true
  try {
    const data = await api.get<House>(`/houses/${route.params.id}`)
    house.value = data
  } catch (error) {
    console.error('Failed to fetch house:', error)
  } finally {
    loading.value = false
  }
}

const fetchSimilarHouses = async () => {
  try {
    const response = await api.get<{ data: House[] }>(`/houses/${route.params.id}/similar?limit=4`)
    similarHouses.value = response.data
  } catch (error) {
    console.error('Failed to fetch similar houses:', error)
  }
}

const checkFavorite = async () => {
  if (!authStore.isAuthenticated) return
  try {
    const response = await api.get<{ data: House[] }>('/favorites')
    isFavorited.value = response.data.some(h => h.id === route.params.id)
  } catch (error) {
    console.error('Failed to check favorite:', error)
  }
}

const toggleFavorite = async () => {
  if (!authStore.isAuthenticated) {
    navigateTo('/login')
    return
  }

  try {
    if (isFavorited.value) {
      await api.del(`/favorites/${route.params.id}`)
      isFavorited.value = false
    } else {
      await api.post(`/favorites/${route.params.id}`)
      isFavorited.value = true
    }
  } catch (error) {
    console.error('Failed to toggle favorite:', error)
  }
}

const sendMessage = async () => {
  if (!house.value || !messageContent.value.trim()) return

  sendingMessage.value = true
  try {
    await api.post('/messages', {
      receiver_id: house.value.landlord.id,
      house_id: house.value.id,
      content: messageContent.value.trim()
    })
    showMessageModal.value = false
    messageContent.value = ''
    alert('消息发送成功！')
  } catch (error: any) {
    alert(error.message || '发送失败，请重试')
  } finally {
    sendingMessage.value = false
  }
}

const initMap = async () => {
  if (!mapContainer.value || !house.value || !house.value.latitude || !house.value.longitude) return

  const L = await import('leaflet')
  
  if (map) {
    map.remove()
  }

  map = L.map(mapContainer.value).setView([house.value.latitude, house.value.longitude], 15)

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© OpenStreetMap contributors'
  }).addTo(map)

  L.marker([house.value.latitude, house.value.longitude])
    .addTo(map)
    .bindPopup(`<b>${house.value.community_name}</b><br>${house.value.address}`)
}

onMounted(() => {
  fetchHouse()
  fetchSimilarHouses()
  checkFavorite()
})

watch(house, () => {
  if (house.value && house.value.latitude && house.value.longitude) {
    nextTick(() => initMap())
  }
}, { deep: true })
</script>
