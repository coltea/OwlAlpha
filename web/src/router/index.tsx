import type { ReactElement } from 'react'
import { Navigate, Route, Routes } from 'react-router-dom'
import { AppLayout } from '../layouts/AppLayout'
import { DashboardPage } from '../pages/dashboard'
import { LoginPage } from '../pages/login'
import { ReportsPage } from '../pages/reports'
import { SettingsPage } from '../pages/settings'

function RequireAuth({ children }: { children: ReactElement }) {
  const token = localStorage.getItem('owlalpha_token')
  if (!token) {
    return <Navigate to="/login" replace />
  }
  return children
}

export function AppRouter() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route
        path="/"
        element={
          <RequireAuth>
            <AppLayout />
          </RequireAuth>
        }
      >
        <Route index element={<DashboardPage />} />
        <Route path="reports" element={<ReportsPage />} />
        <Route path="settings" element={<SettingsPage />} />
      </Route>
    </Routes>
  )
}
