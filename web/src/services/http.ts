import axios from 'axios'

export const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '/api/v1',
})

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('owlalpha_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})
