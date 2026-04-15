import React, { useEffect, useState } from 'react';
import { useAuthStore } from '../store/authStore';
import { useNavigate } from 'react-router-dom';
import { Users, Activity, Plus, X, ChevronLeft } from 'lucide-react';

interface UserData {
  id: string;
  name: string;
  email: string;
  role: string;
}

// UI-only colour palette — kept here since colour is a frontend concern.
// Extend this if new roles are added in data/roles.json.
const ROLE_COLORS: Record<string, string> = {
  super_admin:  '#a78bfa',
  tenant_admin: '#60a5fa',
  doctor:       '#34d399',
  nurse:        '#fb923c',
  staff:        '#9ca3af',
};

interface RoleDefinition {
  value: string;
  label: string;
  description: string;
  assignable: boolean;
}

const TenantAdminDashboard = () => {
  const { user, token, isAuthenticated } = useAuthStore();
  const navigate = useNavigate();
  const [users, setUsers] = useState<UserData[]>([]);
  const [roles, setRoles] = useState<RoleDefinition[]>([]);
  const [showAddModal, setShowAddModal] = useState(false);
  const [newUser, setNewUser] = useState({ name: '', email: '', password: '', role: 'doctor' });
  const [addError, setAddError] = useState('');
  const [addLoading, setAddLoading] = useState(false);
  const role = user?.role;

  const fetchUsers = () => {
    fetch('http://localhost:8080/api/v1/users', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
      .then(res => res.json())
      .then(data => { if (Array.isArray(data)) setUsers(data); })
      .catch(console.error);
  };

  // Fetch assignable roles from backend (single source of truth)
  useEffect(() => {
    if (!token) return;
    fetch('http://localhost:8080/api/v1/roles', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
      .then(res => res.json())
      .then(data => {
        if (Array.isArray(data)) {
          setRoles(data);
          // Set default role for new user form to first assignable role
          setNewUser(prev => ({ ...prev, role: data[0]?.value ?? 'doctor' }));
        }
      })
      .catch(console.error);
  }, [token]);

  useEffect(() => {
    fetchUsers();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [token]);

  const updateRole = async (userId: string, newRole: string) => {
    try {
      const res = await fetch(`http://localhost:8080/api/v1/users/${userId}/role`, {
        method: 'PUT',
        headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
        body: JSON.stringify({ role: newRole })
      });
      if (res.ok) setUsers(users.map(u => u.id === userId ? { ...u, role: newRole } : u));
    } catch (err) { console.error(err); }
  };

  const addUser = async (e: React.FormEvent) => {
    e.preventDefault();
    setAddError('');
    if (!newUser.name || !newUser.email || !newUser.password) {
      setAddError('All fields are required.'); return;
    }
    setAddLoading(true);
    try {
      const res = await fetch('http://localhost:8080/api/v1/users', {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
        body: JSON.stringify(newUser)
      });
      const data = await res.json();
      if (!res.ok) { setAddError(data.error || 'Failed to add user'); return; }
      setUsers(prev => [...prev, data]);
      setNewUser({ name: '', email: '', password: '', role: 'doctor' });
      setShowAddModal(false);
    } catch {
      setAddError('Network error');
    } finally {
      setAddLoading(false);
    }
  };

  if (!isAuthenticated) return null;

  return (
    <div style={{
      display: 'flex', flexDirection: 'column',
      height: '100vh', width: '100vw',
      backgroundColor: 'var(--bg-color)',
      overflow: 'hidden',
      position: 'fixed', top: 0, left: 0
    }}>

      {/* Header */}
      <header style={{
        padding: '0.875rem 1.5rem',
        backgroundColor: 'var(--panel-bg)',
        borderBottom: '1px solid var(--border-color)',
        display: 'flex', justifyContent: 'space-between', alignItems: 'center',
        flexShrink: 0, zIndex: 10
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
          <button
            onClick={() => navigate('/')}
            style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'var(--text-muted)', display: 'flex', alignItems: 'center', gap: '0.25rem', fontSize: '0.875rem' }}
          >
            <ChevronLeft size={16} /> Back
          </button>
          <span style={{ color: 'var(--border-color)' }}>|</span>
          <h2 style={{ margin: 0, fontSize: '1.1rem', fontWeight: 600 }}>Tenant Admin Portal</h2>
        </div>
        <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
          <button
            className="btn"
            onClick={() => setShowAddModal(true)}
            style={{ display: 'flex', alignItems: 'center', gap: '0.4rem', fontSize: '0.875rem' }}
          >
            <Plus size={15} /> Add User
          </button>
          {role === 'super_admin' && (
            <button className="btn btn-secondary" onClick={() => navigate('/super-admin')} style={{ fontSize: '0.8rem' }}>
              Super Admin →
            </button>
          )}
        </div>
      </header>

      {/* Scrollable Content */}
      <div style={{ flex: 1, overflowY: 'auto', overflowX: 'hidden', padding: '1.5rem' }}>

        {/* Stats */}
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '1rem', marginBottom: '1.5rem', maxWidth: '600px' }}>
          <div className="card" style={{ padding: '1.25rem', display: 'flex', alignItems: 'center', gap: '1rem' }}>
            <div style={{ width: 44, height: 44, borderRadius: '0.5rem', backgroundColor: 'rgba(96,165,250,0.15)', display: 'flex', alignItems: 'center', justifyContent: 'center', flexShrink: 0 }}>
              <Users size={22} color="var(--primary)" />
            </div>
            <div>
              <div style={{ fontSize: '1.5rem', fontWeight: 700 }}>{users.length}</div>
              <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>Registered Staff</div>
            </div>
          </div>
          <div className="card" style={{ padding: '1.25rem', display: 'flex', alignItems: 'center', gap: '1rem' }}>
            <div style={{ width: 44, height: 44, borderRadius: '0.5rem', backgroundColor: 'rgba(52,211,153,0.15)', display: 'flex', alignItems: 'center', justifyContent: 'center', flexShrink: 0 }}>
              <Activity size={22} color="var(--success)" />
            </div>
            <div>
              <div style={{ fontSize: '1.5rem', fontWeight: 700, color: 'var(--success)' }}>98.9%</div>
              <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>Operations Health</div>
            </div>
          </div>
        </div>

        {/* Users Table */}
        <div className="card" style={{ padding: '1.25rem' }}>
          <h3 style={{ margin: '0 0 1.25rem 0', fontSize: '1rem' }}>Staff & Role Management</h3>

          <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.875rem', tableLayout: 'fixed' }}>
            <colgroup>
              <col style={{ width: '28%' }} />
              <col style={{ width: '32%' }} />
              <col style={{ width: '20%' }} />
              <col style={{ width: '20%' }} />
            </colgroup>
            <thead>
              <tr style={{ borderBottom: '1px solid var(--border-color)', color: 'var(--text-muted)', textAlign: 'left' }}>
                <th style={{ padding: '0.6rem 0.75rem', fontWeight: 500 }}>Name</th>
                <th style={{ padding: '0.6rem 0.75rem', fontWeight: 500 }}>Email</th>
                <th style={{ padding: '0.6rem 0.75rem', fontWeight: 500 }}>Role</th>
                <th style={{ padding: '0.6rem 0.75rem', fontWeight: 500 }}>Change Role</th>
              </tr>
            </thead>
            <tbody>
              {users.map(u => (
                <tr key={u.id} style={{ borderBottom: '1px solid var(--border-color)' }}>
                  <td style={{ padding: '0.75rem', overflow: 'hidden' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.6rem' }}>
                      <div style={{
                        width: 30, height: 30, borderRadius: '50%',
                        backgroundColor: ROLE_COLORS[u.role] || 'var(--border-color)',
                        display: 'flex', alignItems: 'center', justifyContent: 'center',
                        color: '#fff', fontWeight: 600, fontSize: '0.8rem', flexShrink: 0
                      }}>
                        {u.name.charAt(0).toUpperCase()}
                      </div>
                      <span style={{ fontWeight: 500, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{u.name}</span>
                    </div>
                  </td>
                  <td style={{ padding: '0.75rem', color: 'var(--text-muted)', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{u.email}</td>
                  <td style={{ padding: '0.75rem' }}>
                    <span style={{
                      padding: '0.2rem 0.55rem', borderRadius: '999px', fontSize: '0.72rem', fontWeight: 500,
                      backgroundColor: `${ROLE_COLORS[u.role] ?? '#9ca3af'}22`,
                      color: ROLE_COLORS[u.role] ?? 'var(--text-muted)',
                      border: `1px solid ${ROLE_COLORS[u.role] ?? '#9ca3af'}44`,
                      display: 'inline-block', maxWidth: '100%', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap'
                    }}>
                      {u.role}
                    </span>
                  </td>
                  <td style={{ padding: '0.75rem' }}>
                    {u.role !== 'super_admin' ? (
                      <select
                        value={u.role}
                        onChange={e => updateRole(u.id, e.target.value)}
                        style={{
                          padding: '0.3rem 0.4rem', borderRadius: '6px', width: '100%',
                          border: '1px solid var(--border-color)',
                          backgroundColor: 'var(--bg-color)', color: 'var(--text-main)',
                          fontSize: '0.8rem', cursor: 'pointer'
                        }}
                      >
                        {roles.map(r => (
                          <option key={r.value} value={r.value}>{r.label}</option>
                        ))}
                      </select>
                    ) : (
                      <span style={{ color: 'var(--text-muted)', fontSize: '0.8rem' }}>Owner</span>
                    )}
                  </td>
                </tr>
              ))}
              {users.length === 0 && (
                <tr>
                  <td colSpan={4} style={{ padding: '2.5rem', textAlign: 'center', color: 'var(--text-muted)' }}>
                    <Users size={32} style={{ margin: '0 auto 0.75rem', display: 'block', opacity: 0.3 }} />
                    No staff yet. Use "Add User" to get started.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Add User Modal */}
      {showAddModal && (
        <div style={{ position: 'fixed', inset: 0, backgroundColor: 'rgba(0,0,0,0.65)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000 }}>
          <div className="card" style={{ width: '100%', maxWidth: '420px', margin: '1rem', padding: '1.75rem', display: 'flex', flexDirection: 'column', gap: '1rem' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <h3 style={{ margin: 0 }}>Add New Staff Member</h3>
              <button onClick={() => { setShowAddModal(false); setAddError(''); }} style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'var(--text-muted)' }}>
                <X size={20} />
              </button>
            </div>

            <form onSubmit={addUser} style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
              <div>
                <label style={{ fontSize: '0.8rem', color: 'var(--text-muted)', display: 'block', marginBottom: '0.3rem' }}>Full Name</label>
                <input className="input-control" placeholder="Dr. Jane Smith" value={newUser.name} onChange={e => setNewUser({ ...newUser, name: e.target.value })} />
              </div>
              <div>
                <label style={{ fontSize: '0.8rem', color: 'var(--text-muted)', display: 'block', marginBottom: '0.3rem' }}>Email</label>
                <input className="input-control" type="email" placeholder="jane@hospital.com" value={newUser.email} onChange={e => setNewUser({ ...newUser, email: e.target.value })} />
              </div>
              <div>
                <label style={{ fontSize: '0.8rem', color: 'var(--text-muted)', display: 'block', marginBottom: '0.3rem' }}>Initial Password</label>
                <input className="input-control" type="password" placeholder="••••••••" value={newUser.password} onChange={e => setNewUser({ ...newUser, password: e.target.value })} />
              </div>
              <div>
                <label style={{ fontSize: '0.8rem', color: 'var(--text-muted)', display: 'block', marginBottom: '0.3rem' }}>Role</label>
                <select className="input-control" value={newUser.role} onChange={e => setNewUser({ ...newUser, role: e.target.value })}>
                  {roles.map(r => (
                    <option key={r.value} value={r.value}>{r.label}</option>
                  ))}
                </select>
              </div>

              {addError && (
                <div style={{ color: 'var(--error)', fontSize: '0.8rem', padding: '0.5rem 0.75rem', backgroundColor: 'rgba(248,113,113,0.1)', borderRadius: '6px' }}>
                  {addError}
                </div>
              )}

              <div style={{ display: 'flex', gap: '0.5rem', justifyContent: 'flex-end', marginTop: '0.25rem' }}>
                <button type="button" className="btn btn-secondary" onClick={() => { setShowAddModal(false); setAddError(''); }}>Cancel</button>
                <button type="submit" className="btn" disabled={addLoading}>{addLoading ? 'Adding...' : 'Add User'}</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default TenantAdminDashboard;
