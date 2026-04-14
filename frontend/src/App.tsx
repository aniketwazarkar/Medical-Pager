import React, { useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import ChatLayout from './pages/ChatLayout';
import TenantAdminDashboard from './pages/TenantAdminDashboard';
import SuperAdminDashboard from './pages/SuperAdminDashboard';
import RoleGuard from './components/RoleGuard';
import './styles/index.css';

function App() {
  useEffect(() => {
    document.documentElement.setAttribute('data-theme', 'dark');
  }, []);

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />

        {/* Any authenticated user */}
        <Route path="/" element={
          <RoleGuard>
            <ChatLayout />
          </RoleGuard>
        } />

        {/* Tenant Admin + Super Admin only */}
        <Route path="/admin" element={
          <RoleGuard allowedRoles={['tenant_admin', 'super_admin']}>
            <TenantAdminDashboard />
          </RoleGuard>
        } />

        {/* Super Admin only */}
        <Route path="/super-admin" element={
          <RoleGuard allowedRoles={['super_admin']}>
            <SuperAdminDashboard />
          </RoleGuard>
        } />

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
