import { http } from './http'

export type SettingsForm = {
  baseUrl: string
  apiKey: string
  model: string
}

export type SavedSettings = SettingsForm & {
  id?: number
  checkedAt?: string
}

export async function fetchSettings(): Promise<SavedSettings> {
  const response = await http.get('/settings/openai')
  return response.data.data.config
}

export async function checkSettings(payload: SettingsForm): Promise<string> {
  const response = await http.post('/settings/openai/check', payload)
  return response.data.data.message
}

export async function fetchSupportedModels(payload: Pick<SettingsForm, 'baseUrl' | 'apiKey'>): Promise<string[]> {
  const response = await http.post('/settings/openai/models', payload)
  return response.data.data.models
}

export async function saveSettings(payload: SettingsForm): Promise<SavedSettings> {
  const response = await http.post('/settings/openai', payload)
  return response.data.data.config
}
