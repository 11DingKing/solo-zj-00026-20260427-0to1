<template>
  <div class="min-h-screen bg-gray-50">
    <div class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex flex-wrap gap-4 items-center">
          <input
            v-model="filters.keyword"
            type="text"
            placeholder="搜索小区、地址、描述..."
            class="flex-1 min-w-64 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            @keyup.enter="handleSearch"
          />
          
          <div class="flex items-center gap-2">
            <label class="text-sm text-gray-600">价格:</label>
            <input
              v-model.number="filters.minPrice"
              type="number"
              placeholder="最低"
              class="w-24 px-3 py-2 border border-gray-300 rounded-md text-sm"
            />
            <span class="text-gray-400">-</span>
            <input
              v-model.number="filters.maxPrice"
              type="number"
              placeholder="最高"
              class="w-24 px-3 py-2 border border-gray-300 rounded-md text-sm"
            />
          </div>

          <select
            v-model="filters.bedrooms"
            class="px-3 py-2 border border-gray-300 rounded-md text-sm"
          >
            <option value="">户型不限</option>
            <option value="1">1室</option>
            <option value="2">2室</option>
            <option value="3">3室</option>
            <option value="4">4室</option>
            <option value="5">5室及以上</option>
          </select>

          <div class="flex items-center gap-2">
            <label class="text-sm text-gray-600">面积:</label>
            <input
              v-model.number="filters.minArea"
              type="number"
              placeholder="最小"
              class="w-20 px-3 py-2 border border-gray-300 rounded-md text-sm"
            />
            <span class="text-gray-400">-</span>
            <input
              v-model.number="filters.maxArea"
              type="number"
              placeholder="最大"
              class="w-20 px-3 py-2 border border-gray-300 rounded-md text-sm"
            />
            <span class="text-gray-400">㎡</span>
          </div>

          <button
            @click="showFacilitiesFilter = !showFacilitiesFilter"
            class="px-4 py-2 border border-gray-300 rounded-md text-sm hover:bg-gray-50"
          >
            配套设施 ▼
          </button>

          <button
            @click="handleSearch"
            class="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            搜索
          </button>
        </div>

        <div v-if="showFacilitiesFilter" class="mt-4 pt-4 border-t">
          <div class="flex flex-wrap gap-3">
            <label
              v-for="facility in allFacilities"
              :key="facility"
              class="flex items-center gap-1 cursor-pointer"
            >
              <input
                type="checkbox"
                :value="facility"
                v-model="filters.facilities"
                class="rounded"
              />
              <span class="text-sm">{{ getFacilityLabel(facility) }}</span>
            </label>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center gap-4 py-2">
          <button
            @click="viewMode = 'list'"
            :class="[
              'px-4 py-2 text-sm rounded-md',
              viewMode === 'list' ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100'
            ]"
          >
            列表模式
          </button>
          <button
            @click="viewMode = 'map'"
            :class="[
              'px-4 py-2 text-sm rounded-md',
              viewMode === 'map' ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100'
            ]"
          >
            地图模式
          </button>
          <span class="text-sm text-gray-500 ml-4">
            共 {{ total }} 套房源
          </span>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div v-if="viewMode === 'list'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <NuxtLink
          v-for="house in houses"
          :key="house.id"
          :to="`/houses/${house.id}`"
          class="bg-white rounded-lg shadow-sm overflow-hidden hover:shadow-md transition-shadow"
        >
          <div class="relative h-48 bg-gray-200">
            <img
              v-if="house.images && house.images.length > 0"
              :src="getImageUrl(house.images[0])"
              :alt="house.title"
              class="w-full h-full object-cover"
            />
            <div class="absolute top-3 left-3">
              <span class="bg-red-500 text-white px-2 py-1 rounded text-sm font-bold">
                ¥{{ house.price }}/月
              </span>
            </div>
          </div>
          <div class="p-4">
            <h3 class="font-semibold text-gray-900 truncate">{{ house.title }}</h3>
            <p class="text-gray-500 text-sm mt-1 truncate">{{ house.community_name }}</p>
            <div class="flex items-center gap-2 mt-2 text-sm text-gray-600">
              <span>{{ house.bedrooms }}室{{ house.living_rooms }}厅{{ house.bathrooms }}卫</span>
              <span>|</span>
              <span>{{ house.area }}㎡</span>
              <span>|</span>
              <span>{{ getDecorationLabel(house.decoration) }}</span>
            </div>
            <div class="flex flex-wrap gap-1 mt-2">
              <span
                v-for="facility in house.facilities?.slice(0, 3)"
                :key="facility"
                class="text-xs bg-gray-100 text-gray-600 px-2 py-1 rounded"
              >
                {{ getFacilityLabel(facility) }}
              </span>
            </div>
          </div>
        </NuxtLink>
      </div>

      <div v-else class="bg-white rounded-lg shadow-sm overflow-hidden" style="height: 600px">
        <div ref="mapContainer" class="w-full h-full"></div>
      </div>

      <div v-if="viewMode === 'list' && houses.length > 0" class="mt-8 flex justify-center">
        <nav class="flex items-center gap-2">
          <button
            @click="currentPage > 1 && (currentPage--)"
            :disabled="currentPage <= 1"
            class="px-3 py-1 border rounded disabled:opacity-50 disabled:cursor-not-allowed"
          >
            上一页
          </button>
          <span class="text-gray-600">
            第 {{ currentPage }} 页
          </span>
          <button
            @click="(currentPage < totalPages) && (currentPage++)"
            :disabled="currentPage >= totalPages"
            class="px-3 py-1 border rounded disabled:opacity-50 disabled:cursor-not-allowed"
          >
            下一页
          </button>
        </nav>
      </div>

      <div v-if="houses.length === 0 && !loading" class="text-center py-12">
        <p class="text-gray-500">暂无符合条件的房源</p>
      </div>

      <div v-if="loading" class="text-center py-12">
        <p class="text-gray-500">加载中...</p>
      </div>
    </div>
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
  images: string[]
  landlord: {
    id: string
    nickname: string
    phone: string
  }
}

interface Filters {
  keyword: string
  minPrice: number | null
  maxPrice: number | null
  minArea: number | null
  maxArea: number | null
  bedrooms: string
  orientation: string
  decoration: string
  facilities: string[]
}

const config = useRuntimeConfig()
const api = useApi()

const houses = ref<House[]>([])
const total = ref(0)
const loading = ref(false)
const viewMode = ref<'list' | 'map'>('list')
const currentPage = ref(1)
const pageSize = 20
const showFacilitiesFilter = ref(false)

const mapContainer = ref<HTMLElement>()
let map: L.Map | null = null
let markers: L.Marker[] = []

const allFacilities = [
  'wifi', 'air_conditioner', 'washing_machine', 'refrigerator',
  'water_heater', 'private_bathroom', 'balcony', 'tv',
  'wardrobe', 'sofa', 'bed', 'desk', 'chair'
]

const filters = ref<Filters>({
  keyword: '',
  minPrice: null,
  maxPrice: null,
  minArea: null,
  maxArea: null,
  bedrooms: '',
  orientation: '',
  decoration: '',
  facilities: []
})

const totalPages = computed(() => Math.ceil(total.value / pageSize))

const getImageUrl = (path: string) => {
  if (path.startsWith('http')) return path
  return `${config.public.apiUrl}${path}`
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

const getDecorationLabel = (decoration: string): string => {
  const labels: Record<string, string> = {
    rough: '毛坯',
    simple: '简装',
    standard: '精装',
    luxury: '豪装'
  }
  return labels[decoration] || decoration
}

const fetchHouses = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams()
    params.append('page', currentPage.value.toString())
    params.append('limit', pageSize.toString())

    if (filters.value.keyword) {
      params.append('keyword', filters.value.keyword)
    }
    if (filters.value.minPrice !== null) {
      params.append('min_price', filters.value.minPrice.toString())
    }
    if (filters.value.maxPrice !== null) {
      params.append('max_price', filters.value.maxPrice.toString())
    }
    if (filters.value.minArea !== null) {
      params.append('min_area', filters.value.minArea.toString())
    }
    if (filters.value.maxArea !== null) {
      params.append('max_area', filters.value.maxArea.toString())
    }
    if (filters.value.bedrooms) {
      params.append('bedrooms', filters.value.bedrooms)
    }
    if (filters.value.orientation) {
      params.append('orientation', filters.value.orientation)
    }
    if (filters.value.decoration) {
      params.append('decoration', filters.value.decoration)
    }
    if (filters.value.facilities.length > 0) {
      params.append('facilities', filters.value.facilities.join(','))
    }

    const response = await api.get<{ data: House[]; meta: { total: number } }>(`/houses?${params.toString()}`)
    houses.value = response.data
    total.value = response.meta.total
  } catch (error) {
    console.error('Failed to fetch houses:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchHouses()
}

const initMap = async () => {
  if (!mapContainer.value || !houses.value.length) return

  const L = await import('leaflet')
  
  if (map) {
    map.remove()
  }

  map = L.map(mapContainer.value).setView([39.9042, 116.4074], 12)

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© OpenStreetMap contributors'
  }).addTo(map)

  markers.forEach(m => m.remove())
  markers = []

  houses.value.forEach(house => {
    if (house.latitude && house.longitude) {
      const customIcon = L.divIcon({
        className: 'custom-marker',
        html: `<div class="bg-red-500 text-white px-2 py-1 rounded text-xs font-bold shadow-lg">
          ¥${house.price}
        </div>`,
        iconSize: [80, 30],
        iconAnchor: [40, 30]
      })

      const marker = L.marker([house.latitude, house.longitude], { icon: customIcon })
        .bindPopup(`
          <div class="p-2 min-w-48">
            <img src="${getImageUrl(house.images?.[0] || '')}" alt="${house.title}" class="w-full h-24 object-cover rounded mb-2" onerror="this.style.display='none'" />
            <h4 class="font-semibold text-sm">${house.title}</h4>
            <p class="text-red-500 font-bold">¥${house.price}/月</p>
            <p class="text-gray-500 text-xs">${house.bedrooms}室${house.living_rooms}厅 | ${house.area}㎡</p>
            <a href="/houses/${house.id}" class="text-blue-500 text-sm">查看详情</a>
          </div>
        `)
        .addTo(map!)

      markers.push(marker)
    }
  })

  if (markers.length > 0) {
    const group = L.featureGroup(markers)
    map.fitBounds(group.getBounds().pad(0.1))
  }
}

watch([viewMode, houses], () => {
  if (viewMode.value === 'map' && houses.value.length > 0) {
    nextTick(() => initMap())
  }
})

watch(currentPage, () => {
  fetchHouses()
})

watch(filters, () => {
  currentPage.value = 1
  fetchHouses()
}, { deep: true })

onMounted(() => {
  fetchHouses()
})
</script>
