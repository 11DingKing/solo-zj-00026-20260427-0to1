import { defineStore } from 'pinia'

type Role = 'landlord' | 'tenant'

interface User {
  id: string
  phone: string
  role: Role
  nickname: string | null
  avatar: string | null
  created_at: string
}

interface AuthState {
  user: User | null
  token: string | null
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    token: null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    isLandlord: (state) => state.user?.role === 'landlord',
    isTenant: (state) => state.user?.role === 'tenant',
  },

  actions: {
    setToken(token: string) {
      this.token = token
      const cookie = useCookie<string>('auth_token', {
        maxAge: 60 * 60 * 24 * 7,
        path: '/',
      })
      cookie.value = token
    },

    setUser(user: User) {
      this.user = user
    },

    async login(phone: string, password: string) {
      const api = useApi()
      const response = await api.post<{ token: string; user: User }>('/auth/login', {
        phone,
        password,
      })
      
      this.setToken(response.token)
      this.setUser(response.user)
      
      return response
    },

    async register(phone: string, password: string, role: Role, nickname?: string) {
      const api = useApi()
      const response = await api.post<{ token: string; user: User }>('/auth/register', {
        phone,
        password,
        role,
        nickname,
      })
      
      this.setToken(response.token)
      this.setUser(response.user)
      
      return response
    },

    async fetchUser() {
      try {
        const api = useApi()
        const user = await api.get<User>('/auth/me')
        this.setUser(user)
        return user
      } catch (error) {
        this.logout()
        throw error
      }
    },

    logout() {
      this.user = null
      this.token = null
      const cookie = useCookie('auth_token')
      cookie.value = null
    },

    async updateProfile(data: { nickname?: string; avatar?: string }) {
      const api = useApi()
      const user = await api.put<User>('/auth/me', data)
      this.setUser(user)
      return user
    },
  },
})
