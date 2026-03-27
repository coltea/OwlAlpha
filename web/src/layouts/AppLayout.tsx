import { NavLink, Outlet, useNavigate } from 'react-router-dom'

export function AppLayout() {
  const navigate = useNavigate()

  const logout = () => {
    localStorage.removeItem('owlalpha_token')
    navigate('/login')
  }

  return (
    <div className="shell">
      <aside className="sidebar">
        <div className="sidebar-main">
          <div>
            <div className="brand">OwlAlpha</div>
            <p className="muted">A 股每日报告后台</p>
          </div>
          <nav className="nav">
            <NavLink to="/">仪表盘</NavLink>
            <NavLink to="/reports">报告中心</NavLink>
            <NavLink to="/settings">系统设置</NavLink>
          </nav>
        </div>
        <button className="ghost-button" onClick={logout}>退出登录</button>
      </aside>
      <main className="content">
        <Outlet />
      </main>
    </div>
  )
}
