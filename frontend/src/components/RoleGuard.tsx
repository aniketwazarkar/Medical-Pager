import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';

interface RoleGuardProps {
  children: React.ReactNode;
  // Which roles are allowed to access this route.
  // If empty, only authentication is required (any logged-in user).
  allowedRoles?: string[];
  // Where to redirect if unauthorized. Defaults to '/'.
  redirectTo?: string;
}

/**
 * RoleGuard is the single centralized place where frontend routing access rules live.
 * It reads the role from the persisted JWT user object — the backend enforces the actual
 * API-level authorization via middleware/roles.go (the real security layer).
 */
const RoleGuard = ({ children, allowedRoles = [], redirectTo = '/' }: RoleGuardProps) => {
  const { isAuthenticated, user } = useAuthStore();

  // Not logged in → send to login
  if (!isAuthenticated) return <Navigate to="/login" replace />;

  // If allowedRoles specified, enforce them for UX routing
  if (allowedRoles.length > 0) {
    const userRole = user?.role ?? '';
    if (!allowedRoles.includes(userRole)) {
      return <Navigate to={redirectTo} replace />;
    }
  }

  return <>{children}</>;
};

export default RoleGuard;
