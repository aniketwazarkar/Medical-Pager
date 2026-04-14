import React, { useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useAuthStore } from './store/authStore';
import Login from './pages/Login';
import ChatLayout from './pages/ChatLayout';
import TenantAdminDashboard from './pages/TenantAdminDashboard';
import SuperAdminDashboard from './pages/SuperAdminDashboard';
import './styles/index.css';

const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  if (!isAuthenticated) return <Navigate to="/login" />;
  return <>{children}</>;
};

function App() {
  // Theme initialization
  useEffect(() => {
    // Check if dark mode is preferred or set by tenant
    document.documentElement.setAttribute('data-theme', 'dark'); // Default to dark for premium feel
  }, []);

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route 
          path="/" 
          element={
            <ProtectedRoute>
              <ChatLayout />
            </ProtectedRoute>
          } 
        />
        <Route 
          path="/admin" 
          element={
            <ProtectedRoute>
              <TenantAdminDashboard />
            </ProtectedRoute>
          } 
        />
        <Route 
          path="/super-admin" 
          element={
            <ProtectedRoute>
              <SuperAdminDashboard />
            </ProtectedRoute>
          } 
        />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
