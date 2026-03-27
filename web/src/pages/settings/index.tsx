import { useEffect, useMemo, useState } from 'react'
import axios from 'axios'
import {
  checkSettings,
  fetchSettings,
  fetchSupportedModels,
  saveSettings,
  type SettingsForm,
} from '../../services/settings'

const defaultSettings: SettingsForm = {
  baseUrl: 'https://api.openai.com/v1',
  apiKey: '',
  model: 'gpt-4o-mini',
}

function fingerprint(form: SettingsForm) {
  return JSON.stringify({
    baseUrl: form.baseUrl.trim(),
    apiKey: form.apiKey.trim(),
    model: form.model.trim(),
  })
}

function isValidUrl(value: string) {
  try {
    const url = new URL(value)
    return Boolean(url.protocol === 'http:' || url.protocol === 'https:')
  } catch {
    return false
  }
}

export function SettingsPage() {
  const [form, setForm] = useState<SettingsForm>(defaultSettings)
  const [modelOptions, setModelOptions] = useState<string[]>([defaultSettings.model])
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [checking, setChecking] = useState(false)
  const [saving, setSaving] = useState(false)
  const [fetchingModels, setFetchingModels] = useState(false)
  const [checkedFingerprint, setCheckedFingerprint] = useState('')

  const validation = useMemo(() => {
    if (!form.baseUrl.trim()) {
      return '请填写 Base URL。'
    }
    if (!isValidUrl(form.baseUrl.trim())) {
      return 'Base URL 必须是有效的 http 或 https 地址。'
    }
    if (!form.apiKey.trim()) {
      return '请填写 API Key。'
    }
    if (!form.model.trim()) {
      return '请选择模型。'
    }
    return ''
  }, [form])

  const canSave = checkedFingerprint !== '' && checkedFingerprint === fingerprint(form) && !saving && !checking

  useEffect(() => {
    const load = async () => {
      try {
        const saved = await fetchSettings()
        const nextForm: SettingsForm = {
          baseUrl: saved.baseUrl || defaultSettings.baseUrl,
          apiKey: saved.apiKey || defaultSettings.apiKey,
          model: saved.model || defaultSettings.model,
        }
        setForm(nextForm)
        setModelOptions((current) => Array.from(new Set([...current, nextForm.model].filter(Boolean))))
      } catch (loadError) {
        setError(getErrorMessage(loadError, '加载模型配置失败。'))
      } finally {
        setLoading(false)
      }
    }

    void load()
  }, [])

  const updateField = <K extends keyof SettingsForm>(key: K, value: SettingsForm[K]) => {
    setForm((current) => ({ ...current, [key]: value }))
    setMessage('')
    setError('')
    setCheckedFingerprint('')
  }

  const handleFetchModels = async () => {
    if (!form.baseUrl.trim()) {
      setError('请填写 Base URL。')
      setMessage('')
      return
    }
    if (!isValidUrl(form.baseUrl.trim())) {
      setError('Base URL 必须是有效的 http 或 https 地址。')
      setMessage('')
      return
    }
    if (!form.apiKey.trim()) {
      setError('请填写 API Key。')
      setMessage('')
      return
    }

    setFetchingModels(true)
    setError('')
    setMessage('')

    try {
      const models = await fetchSupportedModels({
        baseUrl: form.baseUrl.trim(),
        apiKey: form.apiKey.trim(),
      })
      setModelOptions(models)
      if (!models.includes(form.model.trim())) {
        updateField('model', models[0] ?? '')
      }
      setMessage(`已拉取 ${models.length} 个可用模型。`)
    } catch (fetchError) {
      setError(getErrorMessage(fetchError, '拉取模型列表失败。'))
    } finally {
      setFetchingModels(false)
    }
  }

  const handleCheck = async () => {
    if (validation) {
      setError(validation)
      setMessage('')
      return
    }

    setChecking(true)
    setError('')
    setMessage('')

    try {
      const successMessage = await checkSettings({
        baseUrl: form.baseUrl.trim(),
        apiKey: form.apiKey.trim(),
        model: form.model.trim(),
      })
      setCheckedFingerprint(fingerprint(form))
      setMessage(successMessage)
      setModelOptions((current) => Array.from(new Set([...current, form.model.trim()].filter(Boolean))))
    } catch (checkError) {
      setCheckedFingerprint('')
      setError(getErrorMessage(checkError, '配置检查失败。'))
    } finally {
      setChecking(false)
    }
  }

  const handleSave = async () => {
    if (validation) {
      setError(validation)
      setMessage('')
      return
    }

    if (checkedFingerprint !== fingerprint(form)) {
      setError('请先对当前配置执行检查，检查成功后才能保存。')
      setMessage('')
      return
    }

    setSaving(true)
    setError('')
    setMessage('')

    try {
      const saved = await saveSettings({
        baseUrl: form.baseUrl.trim(),
        apiKey: form.apiKey.trim(),
        model: form.model.trim(),
      })
      const nextForm: SettingsForm = {
        baseUrl: saved.baseUrl,
        apiKey: saved.apiKey,
        model: saved.model,
      }
      setForm(nextForm)
      setCheckedFingerprint(fingerprint(nextForm))
      setModelOptions((current) => Array.from(new Set([...current, nextForm.model].filter(Boolean))))
      setMessage('模型配置已保存到数据库。')
    } catch (saveError) {
      setCheckedFingerprint('')
      setError(getErrorMessage(saveError, '保存模型配置失败。'))
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return <p>加载配置中...</p>
  }

  return (
    <div>
      <div className="page-header">
        <div>
          <div className="eyebrow">Settings</div>
          <h1>系统设置</h1>
          <p className="muted">配置 OpenAI 兼容接口的 URL、Key 和模型，先支持前端检查与本地保存。</p>
        </div>
      </div>

      <section className="panel settings-panel">
        <div className="settings-heading">
          <div>
            <h3>OpenAI 兼容接口</h3>
            <p className="muted">检查会实际发送模型请求，保存后写入数据库，模型列表需手动点击拉取。</p>
          </div>
          <span className="badge">数据库配置</span>
        </div>

        <div className="form-grid">
          <label className="field">
            <span>Base URL</span>
            <input
              type="url"
              value={form.baseUrl}
              onChange={(event) => updateField('baseUrl', event.target.value)}
              placeholder="https://api.openai.com/v1"
            />
          </label>

          <label className="field">
            <span>API Key</span>
            <input
              type="password"
              value={form.apiKey}
              onChange={(event) => updateField('apiKey', event.target.value)}
              placeholder="sk-..."
            />
          </label>

          <label className="field">
            <span>模型</span>
            <div className="field-with-action">
              <select value={form.model} onChange={(event) => updateField('model', event.target.value)}>
              {modelOptions.map((option) => (
                <option key={option} value={option}>
                  {option}
                </option>
              ))}
              </select>
              <button type="button" className="ghost-button inline-button" onClick={handleFetchModels} disabled={fetchingModels}>
                {fetchingModels ? '拉取中...' : '拉取模型列表'}
              </button>
            </div>
          </label>
        </div>

        {error ? <p className="error-text">{error}</p> : null}
        {message ? <p className="success-text">{message}</p> : null}

        <div className="action-row">
          <button type="button" className="ghost-button" onClick={handleCheck} disabled={checking}>
            {checking ? '检查中...' : '检查'}
          </button>
          <button type="button" className="primary-button" onClick={handleSave} disabled={!canSave}>
            {saving ? '保存中...' : '保存'}
          </button>
        </div>
      </section>
    </div>
  )
}

function getErrorMessage(error: unknown, fallback: string) {
  if (axios.isAxiosError(error)) {
    const responseMessage = error.response?.data?.message
    if (typeof responseMessage === 'string' && responseMessage.trim()) {
      return responseMessage
    }

    const dataMessage = error.response?.data?.data?.message
    if (typeof dataMessage === 'string' && dataMessage.trim()) {
      return dataMessage
    }
  }

  if (error instanceof Error && error.message.trim()) {
    return error.message
  }

  return fallback
}
