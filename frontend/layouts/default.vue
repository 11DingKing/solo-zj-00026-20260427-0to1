<template>
  <div class="min-h-screen bg-gray-50">
    <header class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <NuxtLink to="/" class="flex items-center space-x-2">
            <span class="text-xl font-bold text-blue-600">租房平台</span>
          </NuxtLink>
          
          <nav class="hidden md:flex items-center space-x-6">
            <NuxtLink to="/" class="text-gray-700 hover:text-blue-600">
              房源列表
            </NuxtLink>
            
            <template v-if="authStore.isAuthenticated">
              <template v-if="authStore.isLandlord">
                <NuxtLink to="/houses/create" class="text-gray-700 hover:text-blue-600">
                  发布房源
                </NuxtLink>
              </template>
              
              <NuxtLink to="/profile" class="text-gray-700 hover:text-blue-600">
                个人中心
              </NuxtLink>
              
              <button
                @click="handleLogout"
                class="text-gray-700 hover:text-red-600"
              >
                退出登录
              </button>
            </template>
            
            <template v-else>
              <NuxtLink to="/login" class="text-gray-700 hover:text-blue-600">
                登录
              </NuxtLink>
              <NuxtLink
                to="/register"
                class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
              >
                注册
              </NuxtLink>
            </template>
          </nav>
        </div>
      </div>
    </header>
    
    <main>
      <slot />
    </main>
    
    <footer class="bg-white border-t mt-auto">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div class="text-center text-gray-500 text-sm">
          <p>© 2024 租房信息发布平台. 保留所有权利.</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
const router = useRouter()

onMounted(async () => {
  const tokenCookie = useCookie('auth_token')
  if (tokenCookie.value && !authStore.user) {
    try {
      authStore.setToken(tokenCookie.value)
      await authStore.fetchUser()
    } catch (error) {
      console.error('Failed to fetch user:', error)
    }
  }
})

const handleLogout = () => {
  authStore.logout()
  router.push('/')
}
</script>
