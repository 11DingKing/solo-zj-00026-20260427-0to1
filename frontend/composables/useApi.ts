export const useApi = () => {
  const config = useRuntimeConfig()
  const token = useCookie<string>('auth_token')

  const request = async <T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> => {
    const url = `${config.public.apiUrl}/api${endpoint}`
    
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    }

    if (token.value) {
      headers['Authorization'] = `Bearer ${token.value}`
    }

    const response = await fetch(url, {
      ...options,
      headers,
    })

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: '请求失败' }))
      throw new Error(error.error || `HTTP error! status: ${response.status}`)
    }

    if (response.status === 204) {
      return {} as T
    }

    return response.json()
  }

  const get = <T>(endpoint: string) => request<T>(endpoint, { method: 'GET' })

  const post = <T>(endpoint: string, data?: any) => request<T>(endpoint, {
    method: 'POST',
    body: data ? JSON.stringify(data) : undefined,
  })

  const put = <T>(endpoint: string, data?: any) => request<T>(endpoint, {
    method: 'PUT',
    body: data ? JSON.stringify(data) : undefined,
  })

  const del = <T>(endpoint: string) => request<T>(endpoint, { method: 'DELETE' })

  return { get, post, put, del, request }
}
