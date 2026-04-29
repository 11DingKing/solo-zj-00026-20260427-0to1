<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="bg-white rounded-lg shadow-sm p-6 mb-6">
        <div class="flex items-center gap-4">
          <div class="w-20 h-20 bg-blue-100 rounded-full flex items-center justify-center">
            <span class="text-3xl text-blue-600 font-bold">
              {{ authStore.user?.nickname?.charAt(0) || authStore.user?.phone?.charAt(0) || '用' }}
            </span>
          </div>
          <div>
            <h1 class="text-2xl font-bold text-gray-900">
              {{ authStore.user?.nickname || '用户' }}
            </h1>
            <p class="text-gray-500">{{ authStore.user?.phone }}</p>
            <p class="text-sm text-gray-400">
              {{ authStore.isLandlord ? '房东账号' : '租客账号' }}
            </p>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow-sm overflow-hidden">
        <div class="flex border-b">
          <button
            v-if="authStore.isLandlord"
            @click="activeTab = 'houses'"
            :class="[
              'px-6 py-3 font-medium',
              activeTab === 'houses' ? 'text-blue-600 border-b-2 border-blue-600' : 'text-gray-500 hover:text-gray-700'
            ]"
          >
            我的房源
          </button>
          <button
            v-if="authStore.isLandlord"
            @click="activeTab = 'messages'"
            :class="[
              'px-6 py-3 font-medium',
              activeTab === 'messages' ? 'text-blue-600 border-b-2 border-blue-600' : 'text-gray-500 hover:text-gray-700'
            ]"
          >
            收到的留言
            <span v-if="unreadCount > 0" class="ml-2 bg-red-500 text-white text-xs px-2 py-0.5 rounded-full">
              {{ unreadCount }}
            </span>
          </button>
          <button
            v-if="authStore.isTenant"
            @click="activeTab = 'favorites'"
            :class="[
              'px-6 py-3 font-medium',
              activeTab === 'favorites' ? 'text-blue-600 border-b-2 border-blue-600' : 'text-gray-500 hover:text-gray-700'
            ]"
          >
            我的收藏
          </button>
          <button
            v-if="authStore.isTenant"
            @click="activeTab = 'sent_messages'"
            :class="[
              'px-6 py-3 font-medium',
              activeTab === 'sent_messages' ? 'text-blue-600 border-b-2 border-blue-600' : 'text-gray-500 hover:text-gray-700'
            ]"
          >
            我的留言
          </button>
        </div>

        <div class="p-6">
          <div v-if="activeTab === 'houses'" class="space-y-4">
            <div class="flex justify-between items-center mb-4">
              <select
                v-model="houseStatusFilter"
                class="px-3 py-2 border border-gray-300 rounded-md text-sm"
              >
                <option value="">全部状态</option>
                <option value="draft">草稿</option>
                <option value="published">已发布</option>
                <option value="off_shelves">已下架</option>
                <option value="rented">已出租</option>
              </select>
              <NuxtLink
                to="/houses/create"
                class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
              >
                + 发布新房源
              </NuxtLink>
            </div>

            <div v-if="myHouses.length === 0" class="text-center py-12 text-gray-500">
              暂无房源
            </div>

            <div v-else class="space-y-4">
              <div
                v-for="house in myHouses"
                :key="house.id"
                class="flex gap-4 p-4 border rounded-lg"
              >
                <div class="w-32 h-24 flex-shrink-0 bg-gray-200 rounded overflow-hidden">
                  <img
                    v-if="house.images && house.images.length > 0"
                    :src="getImageUrl(house.images[0])"
                    :alt="house.title"
                    class="w-full h-full object-cover"
                  />
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex justify-between items-start">
                    <div>
                      <h3 class="font-medium text-gray-900 truncate">{{ house.title }}</h3>
                      <p class="text-red-500 font-bold">¥{{ house.price }}/月</p>
                      <p class="text-gray-500 text-sm">
                        {{ house.bedrooms }}室{{ house.living_rooms }}厅 | {{ house.area }}㎡
                      </p>
                    </div>
                    <span
                      :class="[
                        'px-2 py-1 text-xs rounded',
                        house.status === 'published' ? 'bg-green-100 text-green-700' :
                        house.status === 'draft' ? 'bg-gray-100 text-gray-700' :
                        house.status === 'off_shelves' ? 'bg-yellow-100 text-yellow-700' :
                        'bg-blue-100 text-blue-700'
                      ]"
                    >
                      {{ getStatusLabel(house.status) }}
                    </span>
                  </div>
                  <div class="flex gap-2 mt-3">
                    <NuxtLink
                      :to="`/houses/${house.id}`"
                      class="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50"
                    >
                      查看
                    </NuxtLink>
                    <NuxtLink
                      :to="`/houses/create?id=${house.id}`"
                      class="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50"
                    >
                      编辑
                    </NuxtLink>
                    <template v-if="house.status === 'published'">
                      <button
                        @click="updateHouseStatus(house.id, 'off_shelves')"
                        class="px-3 py-1 text-sm border border-yellow-300 text-yellow-700 rounded hover:bg-yellow-50"
                      >
                        下架
                      </button>
                    </template>
                    <template v-else-if="house.status === 'off_shelves' || house.status === 'draft'">
                      <button
                        @click="updateHouseStatus(house.id, 'published')"
                        class="px-3 py-1 text-sm border border-green-300 text-green-700 rounded hover:bg-green-50"
                      >
                        发布
                      </button>
                    </template>
                    <button
                      @click="deleteHouse(house.id)"
                      class="px-3 py-1 text-sm border border-red-300 text-red-700 rounded hover:bg-red-50"
                    >
                      删除
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-else-if="activeTab === 'messages'" class="space-y-4">
            <div v-if="receivedMessages.length === 0" class="text-center py-12 text-gray-500">
              暂无留言
            </div>

            <div v-else class="space-y-4">
              <div
                v-for="msg in receivedMessages"
                :key="msg.id"
                :class="[
                  'p-4 border rounded-lg',
                  !msg.is_read ? 'bg-blue-50' : ''
                ]"
              >
                <div class="flex justify-between items-start mb-2">
                  <div class="flex items-center gap-2">
                    <span class="font-medium">{{ msg.sender?.nickname || msg.sender?.phone }}</span>
                    <span v-if="!msg.is_read" class="w-2 h-2 bg-blue-500 rounded-full"></span>
                  </div>
                  <span class="text-gray-400 text-sm">
                    {{ formatDate(msg.created_at) }}
                  </span>
                </div>
                <p class="text-gray-700">{{ msg.content }}</p>
                <div class="mt-2 text-sm">
                  <span class="text-gray-500">关于房源：</span>
                  <NuxtLink
                    v-if="msg.house"
                    :to="`/houses/${msg.house.id}`"
                    class="text-blue-600 hover:underline"
                  >
                    {{ msg.house.title }}
                  </NuxtLink>
                </div>
                <div v-if="replyingTo === msg.id" class="mt-3">
                  <textarea
                    v-model="replyContent"
                    placeholder="输入回复内容..."
                    class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    rows="3"
                  />
                  <div class="flex gap-2 mt-2">
                    <button
                      @click="sendReply(msg.id)"
                      :disabled="!replyContent.trim()"
                      class="px-4 py-1 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
                    >
                      发送
                    </button>
                    <button
                      @click="replyingTo = null"
                      class="px-4 py-1 border border-gray-300 rounded-md hover:bg-gray-50"
                    >
                      取消
                    </button>
                  </div>
                </div>
                <template v-else>
                  <button
                    @click="startReply(msg)"
                    class="mt-2 text-blue-600 text-sm hover:underline"
                  >
                    回复
                  </button>
                </template>
              </div>
            </div>
          </div>

          <div v-else-if="activeTab === 'favorites'" class="space-y-4">
            <div v-if="favorites.length === 0" class="text-center py-12 text-gray-500">
              暂无收藏的房源
            </div>

            <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div
                v-for="house in favorites"
                :key="house.id"
                class="bg-white border rounded-lg overflow-hidden hover:shadow-md transition-shadow"
              >
                <NuxtLink :to="`/houses/${house.id}`">
                  <div class="relative h-40 bg-gray-200">
                    <img
                      v-if="house.images && house.images.length > 0"
                      :src="getImageUrl(house.images[0])"
                      :alt="house.title"
                      class="w-full h-full object-cover"
                    />
                    <div class="absolute top-2 left-2">
                      <span class="bg-red-500 text-white px-2 py-1 rounded text-sm font-bold">
                        ¥{{ house.price }}/月
                      </span>
                    </div>
                  </div>
                </NuxtLink>
                <div class="p-4">
                  <NuxtLink :to="`/houses/${house.id}`" class="font-medium text-gray-900 hover:text-blue-600">
                    {{ house.title }}
                  </NuxtLink>
                  <p class="text-gray-500 text-sm mt-1">
                    {{ house.bedrooms }}室{{ house.living_rooms }}厅 | {{ house.area }}㎡
                  </p>
                  <button
                    @click="removeFavorite(house.id)"
                    class="mt-3 text-red-600 text-sm hover:underline"
                  >
                    取消收藏
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div v-else-if="activeTab === 'sent_messages'" class="space-y-4">
            <div v-if="sentMessages.length === 0" class="text-center py-12 text-gray-500">
              暂无发出的留言
            </div>

            <div v-else class="space-y-4">
              <div
                v-for="msg in sentMessages"
                :key="msg.id"
                class="p-4 border rounded-lg"
              >
                <div class="flex justify-between items-start mb-2">
                  <div class="text-gray-500 text-sm">
                    发送给：<span class="text-gray-900">{{ msg.receiver?.nickname || msg.receiver?.phone }}</span>
                  </div>
                  <span class="text-gray-400 text-sm">
                    {{ formatDate(msg.created_at) }}
                  </span>
                </div>
                <p class="text-gray-700">{{ msg.content }}</p>
                <div class="mt-2 text-sm">
                  <span class="text-gray-500">关于房源：</span>
                  <NuxtLink
                    v-if="msg.house"
                    :to="`/houses/${msg.house.id}`"
                    class="text-blue-600 hover:underline"
                  >
                    {{ msg.house.title }}
                  </NuxtLink>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
const config = useRuntimeConfig()
const api = useApi()

const activeTab = ref(authStore.isLandlord ? 'houses' : 'favorites')
const houseStatusFilter = ref('')
const myHouses = ref<any[]>([])
const receivedMessages = ref<any[]>([])
const sentMessages = ref<any[]>([])
const favorites = ref<any[]>([])
const unreadCount = ref(0)
const replyingTo = ref<string | null>(null)
const replyContent = ref('')

const getImageUrl = (path: string) => {
  if (path.startsWith('http')) return path
  return `${config.public.apiUrl}${path}`
}

const getStatusLabel = (status: string): string => {
  const labels: Record<string, string> = {
    draft: '草稿',
    published: '已发布',
    off_shelves: '已下架',
    rented: '已出租'
  }
  return labels[status] || status
}

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

const fetchMyHouses = async () => {
  try {
    const params = new URLSearchParams()
    if (houseStatusFilter.value) {
      params.append('status', houseStatusFilter.value)
    }
    const response = await api.get<{ data: any[] }>(`/my-houses?${params.toString()}`)
    myHouses.value = response.data || []
  } catch (error) {
    console.error('Failed to fetch my houses:', error)
    myHouses.value = []
  }
}

const fetchReceivedMessages = async () => {
  try {
    const response = await api.get<{ data: any[] }>('/messages')
    receivedMessages.value = response.data || []
    unreadCount.value = receivedMessages.value.filter((m: any) => !m.is_read).length
  } catch (error) {
    console.error('Failed to fetch messages:', error)
    receivedMessages.value = []
    unreadCount.value = 0
  }
}

const fetchSentMessages = async () => {
  try {
    const response = await api.get<{ data: any[] }>('/messages/sent')
    sentMessages.value = response.data || []
  } catch (error) {
    console.error('Failed to fetch sent messages:', error)
    sentMessages.value = []
  }
}

const fetchFavorites = async () => {
  try {
    const response = await api.get<{ data: any[] }>('/favorites')
    favorites.value = response.data || []
  } catch (error) {
    console.error('Failed to fetch favorites:', error)
    favorites.value = []
  }
}

const updateHouseStatus = async (houseId: string, status: string) => {
  try {
    await api.put(`/houses/${houseId}/status`, { status })
    fetchMyHouses()
  } catch (error) {
    console.error('Failed to update house status:', error)
    alert('操作失败，请重试')
  }
}

const deleteHouse = async (houseId: string) => {
  if (!confirm('确定要删除这个房源吗？')) return

  try {
    await api.del(`/houses/${houseId}`)
    fetchMyHouses()
  } catch (error) {
    console.error('Failed to delete house:', error)
    alert('删除失败，请重试')
  }
}

const startReply = (msg: any) => {
  replyingTo.value = msg.id
  replyContent.value = ''
  if (!msg.is_read) {
    api.put(`/messages/${msg.id}/read`).catch(console.error)
  }
}

const sendReply = async (msgId: string) => {
  if (!replyContent.value.trim()) return

  try {
    await api.post(`/messages/${msgId}/reply`, {
      content: replyContent.value.trim()
    })
    replyingTo.value = null
    replyContent.value = ''
    fetchReceivedMessages()
  } catch (error) {
    console.error('Failed to send reply:', error)
    alert('发送失败，请重试')
  }
}

const removeFavorite = async (houseId: string) => {
  if (!confirm('确定要取消收藏吗？')) return

  try {
    await api.del(`/favorites/${houseId}`)
    fetchFavorites()
  } catch (error) {
    console.error('Failed to remove favorite:', error)
    alert('操作失败，请重试')
  }
}

watch(activeTab, () => {
  if (activeTab.value === 'houses') fetchMyHouses()
  if (activeTab.value === 'messages') fetchReceivedMessages()
  if (activeTab.value === 'sent_messages') fetchSentMessages()
  if (activeTab.value === 'favorites') fetchFavorites()
})

watch(houseStatusFilter, () => {
  fetchMyHouses()
})

onMounted(() => {
  if (authStore.isLandlord) fetchMyHouses()
  if (authStore.isTenant) fetchFavorites()
})
</script>
