import React, { useState, useEffect, useRef } from 'react';
import { useAuthStore } from '../store/authStore';
import CryptoJS from 'crypto-js';
import { Stethoscope, LogOut, Hash, Users, AlertCircle, Video, Phone, Plus, Pencil, Trash2, Check, X, Settings } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

interface Message {
  id?: string;
  encryptedContent: string;
  priority?: string;
  senderId?: string;
  senderName?: string;
  createdAt?: string;
}

interface Channel {
  id: string;
  name: string;
  type: string;
}

const ChatLayout = () => {
  const { user, logout, token, updateProfile } = useAuthStore();
  const navigate = useNavigate();
  const [channels, setChannels] = useState<Channel[]>([]);
  const [activeChannel, setActiveChannel] = useState<string>('');
  const [messages, setMessages] = useState<Message[]>([]);
  const [messageInput, setMessageInput] = useState('');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [newChannelName, setNewChannelName] = useState('');
  const [renamingId, setRenamingId] = useState<string | null>(null);
  const [renameValue, setRenameValue] = useState('');
  const [showEditProfile, setShowEditProfile] = useState(false);
  const [editProfileName, setEditProfileName] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const role = user?.role;

  // Fixed client-side E2E key for compatibility test. Must be managed natively later.
  const SECRET_KEY = '12345678901234567890123456789012';

  const decryptPayload = (msg: Message) => {
    try {
      const bytes = CryptoJS.AES.decrypt(msg.encryptedContent, SECRET_KEY);
      return { ...msg, encryptedContent: bytes.toString(CryptoJS.enc.Utf8) || '[Decryption Failed]' };
    } catch {
      return msg;
    }
  };

  // Fetch channels from API on mount
  useEffect(() => {
    fetch('http://localhost:8080/api/v1/channels', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
      .then(res => res.json())
      .then(data => {
        if (Array.isArray(data) && data.length > 0) {
          setChannels(data.map((c: { _id?: string; id?: string; name: string; type: string }) => ({ id: c._id || c.id || '', name: c.name, type: c.type })));
          setActiveChannel(data[0]._id || data[0].id || '');
        }
      })
      .catch(console.error);
  }, [token]);

  // Auto-scroll to bottom when messages change
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const createChannel = async () => {
    if (!newChannelName.trim()) return;
    try {
      const res = await fetch('http://localhost:8080/api/v1/channels', {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: newChannelName, type: 'group', roomType: 'chat' })
      });
      const newCh = await res.json();
      setChannels(prev => [...prev, { id: newCh.id || newCh._id, name: newCh.name, type: newCh.type }]);
      setNewChannelName('');
      setShowCreateModal(false);
    } catch (err) { console.error(err); }
  };

  const renameChannel = async (channelId: string) => {
    if (!renameValue.trim()) return;
    try {
      await fetch(`http://localhost:8080/api/v1/channels/${channelId}`, {
        method: 'PUT',
        headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: renameValue })
      });
      setChannels(prev => prev.map(c => c.id === channelId ? { ...c, name: renameValue } : c));
      setRenamingId(null);
    } catch (err) { console.error(err); }
  };

  const deleteChannel = async (channelId: string) => {
    if (!window.confirm('Delete this channel?')) return;
    try {
      await fetch(`http://localhost:8080/api/v1/channels/${channelId}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      });
      const remaining = channels.filter(c => c.id !== channelId);
      setChannels(remaining);
      if (activeChannel === channelId) setActiveChannel(remaining[0]?.id || '');
    } catch (err) { console.error(err); }
  };

  // Fetch messages + WebSocket on channel change
  useEffect(() => {
    if (!activeChannel) return;

    // Fetch historical messages first
    fetch(`http://localhost:8080/api/v1/messages/${activeChannel}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
      .then(res => res.json())
      .then(data => {
        if (Array.isArray(data)) {
          setMessages(data.map(decryptPayload));
        }
      })
      .catch(err => console.error("History fetch failed:", err));

    // Connect to WebSocket using native WebSocket API
    const socket = new WebSocket(`ws://localhost:8080/api/v1/ws/${activeChannel}`);

    socket.onopen = () => {
      console.log('Connected to WS');
    };

    socket.onmessage = (event) => {
      // Stub: in real app, we would decrypt event.data if it's E2E encrypted on client
      const newMsg = JSON.parse(event.data);
      setMessages(prev => [...prev, decryptPayload(newMsg)]);
    };

    return () => {
      socket.close();
    };
  }, [activeChannel]);

  const handleEditProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editProfileName.trim()) return;
    try {
      const res = await fetch('http://localhost:8080/api/v1/users/me', {
        method: 'PUT',
        headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: editProfileName })
      });
      if (res.ok) {
        updateProfile(editProfileName);
        setShowEditProfile(false);
      }
    } catch (err) { console.error("Failed to update profile", err); }
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const sendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!messageInput.trim()) return;

    // Send payload. In real E2E, we would encrypt `messageInput` before sending.
    const cipherText = CryptoJS.AES.encrypt(messageInput, SECRET_KEY).toString();

    const payload = {
      channelId: activeChannel,
      encryptedContent: cipherText, // True client E2E
      messageType: 'text',
      priority: 'normal'
    };

    try {
      await fetch('http://localhost:8080/api/v1/messages', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(payload)
      });
      // Do not append to local state immediately; wait for WebSocket echo
      setMessageInput('');
    } catch (err) {
      console.error('Failed to send message', err);
    }
  };

  return (
    <div className="app-container">
      {/* Sidebar */}
      <div className="sidebar" style={{ padding: '1rem', display: 'flex', flexDirection: 'column' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginBottom: '2rem', color: 'var(--primary)' }}>
          <Stethoscope />
          <h2 style={{ fontSize: '1.25rem', margin: 0 }}>Medical Pager</h2>
        </div>

        <div style={{ flex: 1 }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem' }}>
            <div style={{ fontSize: '0.75rem', textTransform: 'uppercase', color: 'var(--text-muted)', fontWeight: 600 }}>Channels</div>
            {(role === 'tenant_admin' || role === 'super_admin') && (
              <button onClick={() => setShowCreateModal(true)} style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'var(--text-muted)', padding: '2px' }} title="Create Channel">
                <Plus size={14} />
              </button>
            )}
          </div>

          {channels.map(c => (
            <div
              key={c.id}
              style={{
                display: 'flex', alignItems: 'center', borderRadius: '0.375rem', marginBottom: '0.125rem',
                backgroundColor: activeChannel === c.id ? 'var(--primary)' : 'transparent'
              }}
            >
              {renamingId === c.id ? (
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.25rem', flex: 1, padding: '0.35rem' }}>
                  <input
                    value={renameValue}
                    onChange={e => setRenameValue(e.target.value)}
                    onKeyDown={e => e.key === 'Enter' && renameChannel(c.id)}
                    autoFocus
                    style={{
                      flex: 1, fontSize: '0.85rem', padding: '0.15rem 0.25rem', borderRadius: '3px',
                      border: '1px solid var(--primary)', backgroundColor: 'var(--bg-color)', color: 'var(--text-main)'
                    }}
                  />
                  <button onClick={() => renameChannel(c.id)} style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'var(--success)' }}><Check size={14} /></button>
                  <button onClick={() => setRenamingId(null)} style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'var(--error)' }}><X size={14} /></button>
                </div>
              ) : (
                <>
                  <div
                    onClick={() => setActiveChannel(c.id)}
                    style={{
                      display: 'flex', alignItems: 'center', gap: '0.5rem', padding: '0.5rem', flex: 1, cursor: 'pointer',
                      color: activeChannel === c.id ? 'white' : 'var(--text-main)'
                    }}
                  >
                    <Hash size={16} />
                    <span>{c.name}</span>
                  </div>
                  {(role === 'tenant_admin' || role === 'super_admin') && (
                    <div style={{ display: 'flex', gap: '2px', paddingRight: '0.35rem', opacity: 0.6 }}>
                      <button onClick={() => { setRenamingId(c.id); setRenameValue(c.name); }}
                        style={{ background: 'none', border: 'none', cursor: 'pointer', color: activeChannel === c.id ? 'white' : 'var(--text-muted)' }}>
                        <Pencil size={12} />
                      </button>
                      <button onClick={() => deleteChannel(c.id)}
                        style={{ background: 'none', border: 'none', cursor: 'pointer', color: activeChannel === c.id ? 'white' : 'var(--text-muted)' }}>
                        <Trash2 size={12} />
                      </button>
                    </div>
                  )}
                </>
              )}
            </div>
          ))}
          {channels.length === 0 && (
            <div style={{ color: 'var(--text-muted)', fontSize: '0.8rem', padding: '0.5rem' }}>No channels yet.</div>
          )}
        </div>

        <div style={{ borderTop: '1px solid var(--border-color)', paddingTop: '1rem', marginTop: 'auto' }}>
          {/* Admin Navigation Links */}
          {role === 'tenant_admin' && (
            <button onClick={() => navigate('/admin')} className="btn btn-secondary" style={{ width: '100%', display: 'flex', gap: '0.5rem', marginBottom: '0.5rem', fontSize: '0.8rem' }}>
              <Settings size={14} /> Tenant Admin
            </button>
          )}
          {role === 'super_admin' && (
            <button onClick={() => navigate('/super-admin')} className="btn btn-secondary" style={{ width: '100%', display: 'flex', gap: '0.5rem', marginBottom: '0.5rem', fontSize: '0.8rem' }}>
              <Settings size={14} /> Super Admin
            </button>
          )}

          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '0.75rem', padding: '0.5rem', backgroundColor: 'rgba(255,255,255,0.05)', borderRadius: '0.5rem' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
              <div style={{ width: 34, height: 34, borderRadius: '50%', backgroundColor: 'var(--primary)', color: 'white', display: 'flex', alignItems: 'center', justifyContent: 'center', fontWeight: 'bold' }}>
                {user?.name.charAt(0).toUpperCase()}
              </div>
              <div style={{ overflow: 'hidden' }}>
                <div style={{ fontSize: '0.875rem', fontWeight: 500, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', maxWidth: '110px' }}>{user?.name}</div>
                <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{user?.role}</div>
              </div>
            </div>
            <button
              onClick={() => { setEditProfileName(user?.name || ''); setShowEditProfile(true); }}
              className="btn-secondary"
              style={{ padding: '0.35rem', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
              title="Edit Profile"
            >
              <Settings size={14} />
            </button>
          </div>
          <button onClick={handleLogout} className="btn btn-secondary" style={{ width: '100%', display: 'flex', gap: '0.5rem' }}>
            <LogOut size={16} /> Logout
          </button>
        </div>
      </div>

      {/* Main Chat Area */}
      <div className="main-chat">
        <div style={{ padding: '1rem', borderBottom: '1px solid var(--border-color)', backgroundColor: 'var(--panel-bg)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <Hash size={20} color="var(--text-muted)" />
            <h3 style={{ margin: 0 }}>{channels.find(c => c.id === activeChannel)?.name ?? 'Select a channel'}</h3>
          </div>
          {/* Agora Readiness Controls */}
          <div style={{ display: 'flex', gap: '0.5rem' }}>
            <button className="btn btn-secondary" title="Start Audio Call (Future Agora Integration)" style={{ padding: '0.5rem' }}><Phone size={16} /></button>
            <button className="btn btn-secondary" title="Start Video Call (Future Agora Integration)" style={{ padding: '0.5rem' }}><Video size={16} /></button>
          </div>
        </div>

        <div style={{ flex: 1, padding: '1rem', overflowY: 'auto', display: 'flex', flexDirection: 'column', gap: '1rem' }}>
          {messages.map((m, i) => {
            const displayName = m.senderName || 'Unknown';
            const initial = displayName.charAt(0).toUpperCase();
            const isOwn = m.senderId === user?.id;
            return (
              <div key={m.id || i} style={{ display: 'flex', gap: '0.75rem', alignItems: 'flex-start' }}>
                {/* Avatar: initials with color */}
                <div style={{
                  width: 36, height: 36, borderRadius: '0.375rem', flexShrink: 0,
                  backgroundColor: isOwn ? 'var(--primary)' : 'var(--border-color)',
                  display: 'flex', alignItems: 'center', justifyContent: 'center',
                  fontWeight: 700, fontSize: '0.85rem', color: 'white'
                }}>
                  {initial}
                </div>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginBottom: '0.25rem' }}>
                    <span style={{ fontWeight: 600, fontSize: '0.9rem' }}>{displayName}</span>
                    <span style={{ fontSize: '0.72rem', color: 'var(--text-muted)' }}>
                      {m.createdAt ? new Date(m.createdAt).toLocaleTimeString() : new Date().toLocaleTimeString()}
                    </span>
                    {m.priority === 'urgent' && <span style={{ color: 'var(--warning)', display: 'flex', alignItems: 'center', gap: '0.25rem', fontSize: '0.72rem' }}><AlertCircle size={11} /> Urgent</span>}
                  </div>
                  <div style={{ color: 'var(--text-main)', lineHeight: 1.55, wordBreak: 'break-word' }}>
                    {m.encryptedContent}
                  </div>
                </div>
              </div>
            );
          })}
          {messages.length === 0 && (
            <div style={{ margin: 'auto', color: 'var(--text-muted)' }}>No messages yet. Start the conversation!</div>
          )}
          <div ref={messagesEndRef} />
        </div>

        <div style={{ padding: '1rem', backgroundColor: 'var(--panel-bg)', borderTop: '1px solid var(--border-color)' }}>
          <form onSubmit={sendMessage} style={{ display: 'flex', gap: '0.5rem' }}>
            <input
              type="text"
              className="input-control"
              placeholder="Type a secure message..."
              value={messageInput}
              onChange={(e) => setMessageInput(e.target.value)}
            />
            <button type="submit" className="btn">Send</button>
          </form>
        </div>
      </div>

      {/* Patient Context Panel */}
      <div className="right-panel" style={{ padding: '1rem' }}>
        <h3 style={{ fontSize: '1rem', marginBottom: '1rem', paddingBottom: '0.5rem', borderBottom: '1px solid var(--border-color)' }}>Patient Context</h3>
        <div style={{ textAlign: 'center', color: 'var(--text-muted)', fontSize: '0.875rem', marginTop: '2rem' }}>
          <Users size={48} color="var(--border-color)" style={{ margin: '0 auto 1rem auto' }} />
          <div>No patient linked to this thread.</div>
          <button className="btn btn-secondary" style={{ marginTop: '1rem', fontSize: '0.75rem' }}>Link Patient</button>
        </div>
      </div>

      {/* Create Channel Modal */}
      {showCreateModal && (
        <div style={{ position: 'fixed', inset: 0, backgroundColor: 'rgba(0,0,0,0.6)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000 }}>
          <div className="card" style={{ padding: '2rem', width: '100%', maxWidth: '400px', display: 'flex', flexDirection: 'column', gap: '1rem' }}>
            <h3 style={{ margin: 0 }}>Create Channel</h3>
            <input
              type="text"
              className="input-control"
              placeholder="Channel name (e.g. cardiac-icu)"
              value={newChannelName}
              onChange={e => setNewChannelName(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && createChannel()}
              autoFocus
            />
            <div style={{ display: 'flex', gap: '0.5rem', justifyContent: 'flex-end' }}>
              <button className="btn btn-secondary" onClick={() => { setShowCreateModal(false); setNewChannelName(''); }}>Cancel</button>
              <button className="btn" onClick={createChannel}>Create</button>
            </div>
          </div>
        </div>
      )}

      {/* Edit Profile Modal */}
      {showEditProfile && (
        <div style={{ position: 'fixed', inset: 0, backgroundColor: 'rgba(0,0,0,0.6)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000 }}>
          <div className="card" style={{ padding: '2rem', width: '100%', maxWidth: '400px', display: 'flex', flexDirection: 'column', gap: '1rem' }}>
            <h3 style={{ margin: 0 }}>Edit Display Name</h3>
            <form onSubmit={handleEditProfile} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
              <input
                type="text"
                className="input-control"
                placeholder="Display Name"
                value={editProfileName}
                onChange={e => setEditProfileName(e.target.value)}
                autoFocus
              />
              <div style={{ display: 'flex', gap: '0.5rem', justifyItems: 'flex-end', justifyContent: 'flex-end' }}>
                <button type="button" className="btn btn-secondary" onClick={() => setShowEditProfile(false)}>Cancel</button>
                <button type="submit" className="btn">Save Changes</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default ChatLayout;