import { http } from './http'

type LoginInput = {
  username: string
  password: string
}

type LoginOutput = {
  token: string
  user: {
    id: number
    username: string
    role: string
  }
}

export async function login(payload: LoginInput): Promise<LoginOutput> {
  const response = await http.post('/auth/login', payload)
  return response.data.data
}
