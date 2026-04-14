import React, { useEffect, useState } from 'react';
import { useAuthStore } from '../store/authStore';
import { Navigate, useNavigate } from 'react-router-dom';
import { Users, Activity } from 'lucide-react';

interface UserData {
  id: string;
  name: string;
  email: string;
  role: string;
}

const TenantAdminDashboard = () => {
  const { user, token, isAuthenticated } = useAuthStore();
  const navigate = useNavigate();
  const [users, setUsers] = useState<UserData[]>([]);
  const role = user?.role;

  useEffect(() => {
    if (role === 'tenant_admin') {
      fetch('http://localhost:8080/api/v1/users', {
        headers: { 'Authorization': `Bearer ${token}` }
      })
        .then(res => res.json())
        .then(data => {
          if (Array.isArray(data)) setUsers(data);
        })
        .catch(console.error);
    }
  }, [role, token]);

  const updateRole = async (userId: string, newRole: string) => {
    try {
      const res = await fetch(`http://localhost:8080/api/v1/users/${userId}/role`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ role: newRole })
      });
      if (res.ok) {
        setUsers(users.map(u => u.id === userId ? { ...u, role: newRole } : u));
      }
    } catch (err) {
      console.error(err);
    }
  };

  if (!isAuthenticated) return <Navigate to="/login" />;
  if (role !== 'tenant_admin') return <Navigate to="/" />;

  return (
    <div className="app-container" style={{ display: 'flex', flexDirection: 'column', backgroundColor: 'var(--bg-color)', overflowY: 'auto' }}>
      <header style={{ padding: '1rem 2rem', backgroundColor: 'var(--panel-bg)', borderBottom: '1px solid var(--border-color)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2>Tenant Management Portal</h2>
        <button className="btn btn-secondary" onClick={() => navigate('/')}>Back to Chat</button>
      </header>

      <div style={{ padding: '2rem', maxWidth: '1200px', margin: '0 auto', width: '100%' }}>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '1.5rem', marginBottom: '2rem' }}>
          <div className="card" style={{ padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', color: 'var(--primary)' }}>
              <Users size={24} />
              <span style={{ fontSize: '1.5rem', fontWeight: 'bold' }}>{users.length}</span>
            </div>
            <div style={{ color: 'var(--text-muted)' }}>Registered Staff</div>
          </div>
          <div className="card" style={{ padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', color: 'var(--success)' }}>
              <Activity size={24} />
              <span style={{ fontSize: '1.5rem', fontWeight: 'bold' }}>98.9%</span>
            </div>
            <div style={{ color: 'var(--text-muted)' }}>Operations Health</div>
          </div>
        </div>

        <div className="card" style={{ padding: '1.5rem' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
            <h3>Manage Roles</h3>
          </div>

          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr style={{ borderBottom: '1px solid var(--border-color)', textAlign: 'left', color: 'var(--text-muted)' }}>
                <th style={{ padding: '0.75rem' }}>Name</th>
                <th style={{ padding: '0.75rem' }}>Email</th>
                <th style={{ padding: '0.75rem' }}>Role</th>
                <th style={{ padding: '0.75rem' }}>Action</th>
              </tr>
            </thead>
            <tbody>
              {users.map(u => (
                <tr key={u.id} style={{ borderBottom: '1px solid var(--border-color)' }}>
                  <td style={{ padding: '0.75rem' }}>{u.name}</td>
                  <td style={{ padding: '0.75rem' }}>{u.email}</td>
                  <td style={{ padding: '0.75rem' }}>
                    <span style={{ padding: '0.25rem 0.5rem', borderRadius: '4px', backgroundColor: 'var(--hover-bg)', fontSize: '0.85rem' }}>
                      {u.role}
                    </span>
                  </td>
                  <td style={{ padding: '0.75rem' }}>
                    <select 
                      value={u.role} 
                      onChange={(e) => updateRole(u.id, e.target.value)}
                      style={{ padding: '0.25rem', borderRadius: '4px', border: '1px solid var(--border-color)', backgroundColor: 'var(--panel-bg)', color: 'var(--text)' }}
                    >
                      <option value="doctor">Doctor</option>
                      <option value="nurse">Nurse</option>
                      <option value="staff">Staff</option>
                      <option value="tenant_admin">Tenant Admin</option>
                    </select>
                  </td>
                </tr>
              ))}
              {users.length === 0 && (
                <tr>
                  <td colSpan={4} style={{ padding: '2rem', textAlign: 'center', color: 'var(--text-muted)' }}>No users found.</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default TenantAdminDashboard;
