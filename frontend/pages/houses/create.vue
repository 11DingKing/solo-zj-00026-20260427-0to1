<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
      <h1 class="text-2xl font-bold text-gray-900 mb-6">发布房源</h1>

      <form @submit.prevent="handleSubmit" class="bg-white rounded-lg shadow-sm p-6 space-y-6">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">标题 *</label>
          <input
            v-model="form.title"
            type="text"
            required
            placeholder="例如：精装两室一厅 近地铁 拎包入住"
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">小区名称 *</label>
            <input
              v-model="form.community_name"
              type="text"
              required
              placeholder="例如：阳光花园"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">月租金（元） *</label>
            <input
              v-model.number="form.price"
              type="number"
              required
              min="0"
              placeholder="例如：3000"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">详细地址 *</label>
          <input
            v-model="form.address"
            type="text"
            required
            placeholder="例如：北京市朝阳区建国路88号"
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">经度</label>
            <input
              v-model.number="form.longitude"
              type="number"
              step="0.0000001"
              placeholder="例如：116.4074"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">纬度</label>
            <input
              v-model.number="form.latitude"
              type="number"
              step="0.0000001"
              placeholder="例如：39.9042"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">面积（㎡） *</label>
            <input
              v-model.number="form.area"
              type="number"
              required
              min="0"
              step="0.01"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">楼层</label>
            <input
              v-model.number="form.floor"
              type="number"
              min="0"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">总楼层</label>
            <input
              v-model.number="form.total_floors"
              type="number"
              min="0"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">卧室（室） *</label>
            <input
              v-model.number="form.bedrooms"
              type="number"
              required
              min="0"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">客厅（厅） *</label>
            <input
              v-model.number="form.living_rooms"
              type="number"
              required
              min="0"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">卫生间（卫） *</label>
            <input
              v-model.number="form.bathrooms"
              type="number"
              required
              min="0"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">朝向</label>
            <select
              v-model="form.orientation"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">请选择</option>
              <option value="north">北</option>
              <option value="south">南</option>
              <option value="east">东</option>
              <option value="west">西</option>
              <option value="northeast">东北</option>
              <option value="northwest">西北</option>
              <option value="southeast">东南</option>
              <option value="southwest">西南</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">装修程度</label>
            <select
              v-model="form.decoration"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">请选择</option>
              <option value="rough">毛坯</option>
              <option value="simple">简装</option>
              <option value="standard">精装</option>
              <option value="luxury">豪装</option>
            </select>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">入住时间</label>
            <input
              v-model="form.move_in_date"
              type="date"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">最短租期（月）</label>
            <input
              v-model.number="form.min_lease_term"
              type="number"
              min="0"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">配套设施</label>
          <div class="flex flex-wrap gap-3">
            <label
              v-for="facility in allFacilities"
              :key="facility.value"
              class="flex items-center gap-1 cursor-pointer"
            >
              <input
                type="checkbox"
                :value="facility.value"
                v-model="form.facilities"
                class="rounded"
              />
              <span class="text-sm">{{ facility.label }}</span>
            </label>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">房源描述</label>
          <textarea
            v-model="form.description"
            rows="5"
            placeholder="请详细描述房源情况，如装修情况、周边配套、交通情况等..."
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">房源图片（最多9张）</label>
          <div class="flex flex-wrap gap-3">
            <div
              v-for="(img, index) in form.images"
              :key="index"
              class="relative w-24 h-24 rounded-lg overflow-hidden border"
            >
              <img :src="getImageUrl(img)" class="w-full h-full object-cover" />
              <button
                type="button"
                @click="removeImage(index)"
                class="absolute top-1 right-1 w-5 h-5 bg-black/50 text-white rounded-full text-xs hover:bg-black/70"
              >
                ✕
              </button>
            </div>
            <label
              v-if="form.images.length < 9"
              class="w-24 h-24 flex items-center justify-center border-2 border-dashed border-gray-300 rounded-lg cursor-pointer hover:border-blue-500"
            >
              <input
                type="file"
                ref="fileInput"
                accept="image/*"
                multiple
                class="hidden"
                @change="handleImageUpload"
              />
              <span class="text-gray-400 text-2xl">+</span>
            </label>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">发布状态</label>
          <select
            v-model="form.status"
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="published">立即发布</option>
            <option value="draft">保存为草稿</option>
          </select>
        </div>

        <div v-if="error" class="text-red-600 text-sm">
          {{ error }}
        </div>

        <div class="flex justify-end gap-4">
          <NuxtLink
            to="/profile"
            class="px-6 py-2 border border-gray-300 rounded-md hover:bg-gray-50"
          >
            取消
          </NuxtLink>
          <button
            type="submit"
            :disabled="loading"
            class="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
          >
            {{ loading ? '发布中...' : (editMode ? '保存修改' : '发布房源') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const config = useRuntimeConfig()
const api = useApi()

const editMode = computed(() => !!route.query.id)
const fileInput = ref<HTMLInputElement | null>(null)

const allFacilities = [
  { value: 'wifi', label: 'WiFi' },
  { value: 'air_conditioner', label: '空调' },
  { value: 'washing_machine', label: '洗衣机' },
  { value: 'refrigerator', label: '冰箱' },
  { value: 'water_heater', label: '热水器' },
  { value: 'private_bathroom', label: '独卫' },
  { value: 'balcony', label: '阳台' },
  { value: 'tv', label: '电视' },
  { value: 'wardrobe', label: '衣柜' },
  { value: 'sofa', label: '沙发' },
  { value: 'bed', label: '床' },
  { value: 'desk', label: '书桌' },
  { value: 'chair', label: '椅子' },
]

interface FormData {
  title: string
  community_name: string
  address: string
  longitude: number | null
  latitude: number | null
  price: number | null
  area: number | null
  bedrooms: number | null
  living_rooms: number | null
  bathrooms: number | null
  floor: number | null
  total_floors: number | null
  orientation: string
  decoration: string
  facilities: string[]
  move_in_date: string
  min_lease_term: number | null
  description: string
  images: string[]
  status: string
}

const form = ref<FormData>({
  title: '',
  community_name: '',
  address: '',
  longitude: null,
  latitude: null,
  price: null,
  area: null,
  bedrooms: null,
  living_rooms: null,
  bathrooms: null,
  floor: null,
  total_floors: null,
  orientation: '',
  decoration: '',
  facilities: [],
  move_in_date: '',
  min_lease_term: null,
  description: '',
  images: [],
  status: 'published',
})

const loading = ref(false)
const error = ref('')

const getImageUrl = (path: string) => {
  if (path.startsWith('http')) return path
  return `${config.public.apiUrl}${path}`
}

const handleImageUpload = async (event: Event) => {
  const input = event.target as HTMLInputElement
  if (!input.files || input.files.length === 0) return

  const remainingSlots = 9 - form.value.images.length
  const filesToUpload = Array.from(input.files).slice(0, remainingSlots)

  for (const file of filesToUpload) {
    try {
      const formData = new FormData()
      formData.append('images', file)

      const config = useRuntimeConfig()
      const tokenCookie = useCookie('auth_token')

      const uploadUrl = editMode.value
        ? `${config.public.apiUrl}/api/houses/${route.query.id}/images`
        : `${config.public.apiUrl}/api/houses/temp/images`

      const response = await fetch(uploadUrl, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${tokenCookie.value}`,
        },
        body: formData,
      })

      if (response.ok) {
        const data = await response.json()
        form.value.images.push(...data.urls)
      }
    } catch (e) {
      console.error('Failed to upload image:', e)
    }
  }

  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const removeImage = (index: number) => {
  form.value.images.splice(index, 1)
}

const loadHouseData = async () => {
  if (!editMode.value) return

  try {
    const house = await api.get<{
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
      move_in_date: string
      min_lease_term: number
      description: string
      images: string[]
      status: string
    }>(`/houses/${route.query.id}`)

    form.value = {
      title: house.title,
      community_name: house.community_name,
      address: house.address,
      longitude: house.longitude,
      latitude: house.latitude,
      price: house.price,
      area: house.area,
      bedrooms: house.bedrooms,
      living_rooms: house.living_rooms,
      bathrooms: house.bathrooms,
      floor: house.floor,
      total_floors: house.total_floors,
      orientation: house.orientation,
      decoration: house.decoration,
      facilities: house.facilities,
      move_in_date: house.move_in_date ? new Date(house.move_in_date).toISOString().split('T')[0] : '',
      min_lease_term: house.min_lease_term,
      description: house.description,
      images: house.images,
      status: house.status,
    }
  } catch (e) {
    console.error('Failed to load house:', e)
  }
}

const handleSubmit = async () => {
  loading.value = true
  error.value = ''

  try {
    const submitData: Record<string, any> = {
      title: form.value.title,
      community_name: form.value.community_name,
      address: form.value.address,
      price: form.value.price,
      area: form.value.area,
      bedrooms: form.value.bedrooms,
      living_rooms: form.value.living_rooms,
      bathrooms: form.value.bathrooms,
      status: form.value.status,
    }

    if (form.value.longitude !== null) submitData.longitude = form.value.longitude
    if (form.value.latitude !== null) submitData.latitude = form.value.latitude
    if (form.value.floor !== null) submitData.floor = form.value.floor
    if (form.value.total_floors !== null) submitData.total_floors = form.value.total_floors
    if (form.value.orientation) submitData.orientation = form.value.orientation
    if (form.value.decoration) submitData.decoration = form.value.decoration
    if (form.value.facilities.length > 0) submitData.facilities = form.value.facilities
    if (form.value.move_in_date) submitData.move_in_date = form.value.move_in_date
    if (form.value.min_lease_term !== null) submitData.min_lease_term = form.value.min_lease_term
    if (form.value.description) submitData.description = form.value.description
    if (form.value.images.length > 0) submitData.images = form.value.images

    let response
    if (editMode.value) {
      response = await api.put(`/houses/${route.query.id}`, submitData)
    } else {
      response = await api.post('/houses', submitData)
    }

    router.push('/profile')
  } catch (e: any) {
    error.value = e.message || '发布失败，请重试'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadHouseData()
})
</script>
