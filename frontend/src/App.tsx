import { Navigate, Route, Routes } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'
import ProtectedRoute from '@/routes/ProtectedRoute'
import AppLayout from '@/components/AppLayout'
import LoginPage from '@/pages/LoginPage'
import DashboardPage from '@/pages/DashboardPage'
import CustomersPage from '@/pages/CustomersPage'
import ConstructionSitesPage from '@/pages/ConstructionSitesPage'
import DriversPage from '@/pages/DriversPage'
import TrucksPage from '@/pages/TrucksPage'
import ConcreteGradesPage from '@/pages/ConcreteGradesPage'
import ComponentCatalogPage from '@/pages/ComponentCatalogPage'
import WaitingMixesPage from '@/pages/WaitingMixesPage'
import WeighingTicketsPage from '@/pages/WeighingTicketsPage'
import WeighingStatsPage from '@/pages/WeighingStatsPage'
import StationConfigPage from '@/pages/StationConfigPage'
import StationRunPage from '@/pages/StationRunPage'

// Định nghĩa toàn bộ route của ứng dụng.
export default function App() {
  const token = useAuthStore((s) => s.token)

  return (
    <Routes>
      {/* Trang đăng nhập: nếu đã đăng nhập thì chuyển vào dashboard. */}
      <Route
        path="/login"
        element={token ? <Navigate to="/" replace /> : <LoginPage />}
      />

      {/* Các route yêu cầu xác thực, nằm trong layout chung. */}
      <Route
        element={
          <ProtectedRoute>
            <AppLayout />
          </ProtectedRoute>
        }
      >
        <Route path="/" element={<DashboardPage />} />
        <Route path="/customers" element={<CustomersPage />} />
        <Route path="/construction-sites" element={<ConstructionSitesPage />} />
        <Route path="/drivers" element={<DriversPage />} />
        <Route path="/mixer-trucks" element={<TrucksPage />} />
        <Route path="/concrete-grades" element={<ConcreteGradesPage />} />
        <Route path="/component-catalog" element={<ComponentCatalogPage />} />
        <Route path="/waiting-mixes" element={<WaitingMixesPage />} />
        <Route path="/weighing-tickets" element={<WeighingTicketsPage />} />
        <Route path="/weighing-stats" element={<WeighingStatsPage />} />
        <Route path="/station-config" element={<StationConfigPage />} />
      </Route>

      {/* Trang chạy trạm trộn: hiển thị toàn màn hình, KHÔNG dùng AppLayout
          (nên sidebar không xuất hiện). Vẫn yêu cầu đăng nhập. */}
      <Route
        path="/run"
        element={
          <ProtectedRoute>
            <StationRunPage />
          </ProtectedRoute>
        }
      />

      {/* Mọi đường dẫn khác quay về trang chủ. */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
