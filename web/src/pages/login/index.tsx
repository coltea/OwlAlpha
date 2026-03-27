import { FormEvent, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { login } from '../../services/auth'

export function LoginPage() {
  const navigate = useNavigate()
  const [username, setUsername] = useState('admin')
  const [password, setPassword] = useState('admin123456')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const onSubmit = async (event: FormEvent) => {
    event.preventDefault()
    setLoading(true)
    setError('')

    try {
      const result = await login({ username, password })
      localStorage.setItem('owlalpha_token', result.token)
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : '登录失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-page">
      <form className="login-card" onSubmit={onSubmit}>
        <div>
          <div className="eyebrow">OwlAlpha 0.1</div>
          <h1>后台登录</h1>
          <p className="muted">登录后查看 A 股每日报告、任务记录和系统配置。</p>
        </div>
        <label>
          <span>用户名</span>
          <input value={username} onChange={(e) => setUsername(e.target.value)} />
        </label>
        <label>
          <span>密码</span>
          <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        </label>
        {error ? <div className="error-text">{error}</div> : null}
        <button className="primary-button" type="submit" disabled={loading}>
          {loading ? '登录中...' : '登录'}
        </button>
      </form>
    </div>
  )
}
