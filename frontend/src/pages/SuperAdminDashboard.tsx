import React, { useEffect, useState } from 'react';
import { useAuthStore } from '../store/authStore';
import { Navigate, useNavigate } from 'react-router-dom';
import { Settings, Plus, LayoutGrid } from 'lucide-react';

interface TenantData {
  id: string;
  name: string;
  domain: string;
}

const SuperAdminDashboard = () => {
  const { user, token, isAuthenticated } = useAuthStore();
  const navigate = useNavigate();
  const [tenants, setTenants] = useState<TenantData[]>([]);
  const role = user?.role;

  useEffect(() => {
    if (role === 'super_admin') {
      fetch('http://localhost:8080/api/v1/tenants', {
        headers: { 'Authorization': `Bearer ${token}` }
      })
        .then(res => res.json())
        .then(data => {
          if (Array.isArray(data)) setTenants(data);
        })
        .catch(console.error);
    }
  }, [role, token]);

  if (!isAuthenticated) return <Navigate to="/login" />;
  if (role !== 'super_admin') return <Navigate to="/" />;

  return (
    <div className="app-container" style={{ display: 'flex', flexDirection: 'column', backgroundColor: 'var(--bg-color)', overflowY: 'auto' }}>
      <header style={{ padding: '1rem 2rem', backgroundColor: 'var(--panel-bg)', borderBottom: '1px solid var(--border-color)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2>System Super Admin</h2>
        <button className="btn btn-secondary" onClick={() => navigate('/')}>Back to Hub</button>
      </header>

      <div style={{ padding: '2rem', maxWidth: '1200px', margin: '0 auto', width: '100%' }}>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '1.5rem', marginBottom: '2rem' }}>
          <div className="card" style={{ padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', color: 'var(--primary)' }}>
              <LayoutGrid size={24} />
              <span style={{ fontSize: '1.5rem', fontWeight: 'bold' }}>{tenants.length}</span>
            </div>
            <div style={{ color: 'var(--text-muted)' }}>Registered Tenants</div>
          </div>
          <div className="card" style={{ padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', color: 'var(--warning)' }}>
              <Settings size={24} />
            </div>
            <div style={{ color: 'var(--text-muted)' }}>Global Configurations</div>
          </div>
        </div>

        <div className="card" style={{ padding: '1.5rem' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
            <h3>Manage Tenants</h3>
            <button className="btn" style={{ display: 'flex', alignItems: 'center', gap: '0.25rem' }}><Plus size={16} /> Add Tenant</button>
          </div>

          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr style={{ borderBottom: '1px solid var(--border-color)', textAlign: 'left', color: 'var(--text-muted)' }}>
                <th style={{ padding: '0.75rem' }}>Tenant ID</th>
                <th style={{ padding: '0.75rem' }}>Hospital / Name</th>
                <th style={{ padding: '0.75rem' }}>Domain</th>
                <th style={{ padding: '0.75rem' }}>Action</th>
              </tr>
            </thead>
            <tbody>
              {tenants.map(t => (
                <tr key={t.id} style={{ borderBottom: '1px solid var(--border-color)' }}>
                  <td style={{ padding: '0.75rem', fontFamily: 'monospace', fontSize: '0.85rem' }}>{t.id}</td>
                  <td style={{ padding: '0.75rem', fontWeight: 'bold' }}>{t.name || 'Unnamed Hospital'}</td>
                  <td style={{ padding: '0.75rem' }}>{t.domain || 'N/A'}</td>
                  <td style={{ padding: '0.75rem' }}>
                    <button className="btn btn-secondary" style={{ padding: '0.25rem 0.75rem', fontSize: '0.85rem' }}>Manage</button>
                  </td>
                </tr>
              ))}
              {tenants.length === 0 && (
                <tr>
                  <td colSpan={4} style={{ padding: '2rem', textAlign: 'center', color: 'var(--text-muted)' }}>No tenants provisioned.</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default SuperAdminDashboard;
