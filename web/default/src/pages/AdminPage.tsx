import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../stores/authStore';
import { LanguageSwitcher } from '../components/LanguageSwitcher';
import api from '../lib/api';

interface AdminPageProps {
  onNavigate: (page: string) => void;
}

interface UserData {
  id: number;
  username: string;
  display_name: string;
  role: number;
  status: number;
  email: string;
  group: string;
  created_at: string;
}

interface PrivilegeGroup {
  id: number;
  name: string;
  display_name: string;
  description: string;
  quota: number;
  rate_limit: number;
  status: number;
}

export function AdminPage({ onNavigate }: AdminPageProps) {
  const { t } = useTranslation();
  const { user: currentUser } = useAuthStore();
  const [activeTab, setActiveTab] = useState<'users' | 'groups'>('users');
  const [users, setUsers] = useState<UserData[]>([]);
  const [groups, setGroups] = useState<PrivilegeGroup[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreateGroup, setShowCreateGroup] = useState(false);
  const [showCreateUser, setShowCreateUser] = useState(false);
  const [showResetPassword, setShowResetPassword] = useState<number | null>(null);
  const [newPassword, setNewPassword] = useState('');
  const [newGroup, setNewGroup] = useState({ name: '', display_name: '', description: '', quota: 100000, rate_limit: 60 });
  const [newUser, setNewUser] = useState({ username: '', password: '', email: '', role: 1 });

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [usersRes, groupsRes]: any[] = await Promise.all([
        api.get('/admin/user'),
        api.get('/admin/privilege-group'),
      ]);
      if (usersRes.success) setUsers(usersRes.data?.items || []);
      if (groupsRes.success) setGroups(groupsRes.data || []);
    } catch (err) {
      console.error('Failed to fetch data');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateUser = async () => {
    if (!newUser.username || !newUser.password) {
      alert(t('username_password_required'));
      return;
    }
    if (newUser.password.length < 8) {
      alert(t('password_min_length'));
      return;
    }
    try {
      const response: any = await api.post('/admin/user', newUser);
      if (response.success) {
        setShowCreateUser(false);
        setNewUser({ username: '', password: '', email: '', role: 1 });
        fetchData();
      }
    } catch (err: any) {
      alert(err.response?.data?.message || t('create_user_failed'));
    }
  };

  const handleCreateGroup = async () => {
    try {
      const response: any = await api.post('/admin/privilege-group', newGroup);
      if (response.success) {
        setShowCreateGroup(false);
        setNewGroup({ name: '', display_name: '', description: '', quota: 100000, rate_limit: 60 });
        fetchData();
      }
    } catch (err) {
      console.error('Failed to create group');
    }
  };

  const handleDeleteGroup = async (id: number) => {
    if (!confirm(t('confirm_delete'))) return;
    try {
      await api.delete(`/admin/privilege-group/${id}`);
      fetchData();
    } catch (err) {
      console.error('Failed to delete group');
    }
  };

  const handleManageUser = async (userId: number, action: string, value?: number) => {
    try {
      await api.post('/admin/user/manage', { id: userId, action, value });
      fetchData();
    } catch (err) {
      console.error('Failed to manage user');
    }
  };

  const handleResetPassword = async (userId: number) => {
    if (!newPassword || newPassword.length < 8) {
      alert(t('password_min_length'));
      return;
    }
    try {
      await api.post('/admin/user/manage', { 
        id: userId, 
        action: 'reset_password', 
        password: newPassword 
      });
      setShowResetPassword(null);
      setNewPassword('');
      alert(t('password_reset_success'));
    } catch (err) {
      console.error('Failed to reset password');
    }
  };

  const handleSetRole = async (userId: number, newRole: number) => {
    try {
      await api.post('/admin/user/manage', { id: userId, action: 'set_role', value: newRole });
      fetchData();
    } catch (err) {
      console.error('Failed to set role');
    }
  };

  const handleSetGroup = async (userId: number, newGroup: string) => {
    try {
      await api.post('/admin/user/manage', { id: userId, action: 'set_group', password: newGroup });
      fetchData();
    } catch (err) {
      console.error('Failed to set group');
    }
  };

  const canManageUser = (targetUser: UserData) => {
    if (!currentUser) return false;
    if (currentUser.role === 100) return true; // Root can manage anyone
    if (currentUser.role === 10 && targetUser.role <= 10) return true; // Admin can manage admin and common users
    return false;
  };

  const canCreateUser = () => {
    if (!currentUser) return false;
    return currentUser.role >= 10;
  };

  const getRoleBadge = (role: number) => {
    if (role === 100) return <span style={{ background: '#dc2626', color: 'white', padding: '0.125rem 0.5rem', borderRadius: '4px', fontSize: '0.75rem' }}>{t('role_root')}</span>;
    if (role === 10) return <span style={{ background: '#2563eb', color: 'white', padding: '0.125rem 0.5rem', borderRadius: '4px', fontSize: '0.75rem' }}>{t('role_admin')}</span>;
    return <span style={{ background: '#6b7280', color: 'white', padding: '0.125rem 0.5rem', borderRadius: '4px', fontSize: '0.75rem' }}>{t('role_user')}</span>;
  };

  const getStatusBadge = (status: number) => {
    if (status === 1) return <span style={{ color: '#16a34a' }}>●</span>;
    return <span style={{ color: '#dc2626' }}>●</span>;
  };

  if (loading) {
    return (
      <div className="dashboard" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <div className="spinner" style={{ width: '2rem', height: '2rem', borderColor: '#667eea', borderTopColor: 'transparent' }} />
      </div>
    );
  }

  return (
    <div className="dashboard">
      {/* Navbar */}
      <nav className="navbar">
        <div className="navbar-brand">
          <button onClick={() => onNavigate('dashboard')} style={{ background: 'none', border: 'none', cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
            <div className="navbar-logo">
              <span>T</span>
            </div>
            <span className="navbar-title">{t('app_name')}</span>
          </button>
        </div>
        <div className="navbar-actions">
          <LanguageSwitcher />
          <button onClick={() => onNavigate('dashboard')} className="btn btn-ghost">
            {t('back_to_dashboard')}
          </button>
        </div>
      </nav>

      {/* Main Content */}
      <div className="main-content">
        {/* Tabs */}
        <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '1.5rem', borderBottom: '1px solid #e2e8f0', paddingBottom: '0.5rem' }}>
          <button
            onClick={() => setActiveTab('users')}
            style={{
              padding: '0.5rem 1rem',
              borderRadius: '8px',
              border: 'none',
              background: activeTab === 'users' ? '#667eea' : 'transparent',
              color: activeTab === 'users' ? 'white' : '#64748b',
              cursor: 'pointer',
              fontWeight: 500,
            }}
          >
            {t('users')} ({users.length})
          </button>
          <button
            onClick={() => setActiveTab('groups')}
            style={{
              padding: '0.5rem 1rem',
              borderRadius: '8px',
              border: 'none',
              background: activeTab === 'groups' ? '#667eea' : 'transparent',
              color: activeTab === 'groups' ? 'white' : '#64748b',
              cursor: 'pointer',
              fontWeight: 500,
            }}
          >
            {t('privilege_groups')} ({groups.length})
          </button>
        </div>

        {/* Users Tab */}
        {activeTab === 'users' && (
          <div>
            {canCreateUser() && (
              <div style={{ marginBottom: '1rem' }}>
                <button
                  onClick={() => setShowCreateUser(true)}
                  className="btn btn-primary"
                >
                  + {t('create_user')}
                </button>
              </div>
            )}
            <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
              <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                <thead>
                  <tr style={{ background: '#f8fafc' }}>
                    <th style={{ padding: '0.75rem 1rem', textAlign: 'left', fontSize: '0.75rem', fontWeight: 600, color: '#64748b', textTransform: 'uppercase' }}>{t('user')}</th>
                    <th style={{ padding: '0.75rem 1rem', textAlign: 'left', fontSize: '0.75rem', fontWeight: 600, color: '#64748b', textTransform: 'uppercase' }}>{t('role')}</th>
                    <th style={{ padding: '0.75rem 1rem', textAlign: 'left', fontSize: '0.75rem', fontWeight: 600, color: '#64748b', textTransform: 'uppercase' }}>{t('group')}</th>
                    <th style={{ padding: '0.75rem 1rem', textAlign: 'left', fontSize: '0.75rem', fontWeight: 600, color: '#64748b', textTransform: 'uppercase' }}>{t('status')}</th>
                    <th style={{ padding: '0.75rem 1rem', textAlign: 'right', fontSize: '0.75rem', fontWeight: 600, color: '#64748b', textTransform: 'uppercase' }}>{t('actions')}</th>
                  </tr>
                </thead>
                <tbody>
                  {users.map((user) => (
                    <tr key={user.id} style={{ borderTop: '1px solid #e2e8f0' }}>
                      <td style={{ padding: '0.75rem 1rem' }}>
                        <div>
                          <div style={{ fontWeight: 500, color: '#0f172a' }}>{user.username}</div>
                          <div style={{ fontSize: '0.75rem', color: '#64748b' }}>{user.email || t('no_email')}</div>
                        </div>
                      </td>
                      <td style={{ padding: '0.75rem 1rem' }}>
                        {canManageUser(user) && user.role !== 100 ? (
                          <select
                            value={user.role}
                            onChange={(e) => handleSetRole(user.id, parseInt(e.target.value))}
                            style={{
                              padding: '0.25rem 0.5rem',
                              borderRadius: '4px',
                              border: '1px solid #e2e8f0',
                              fontSize: '0.75rem',
                              background: 'white',
                            }}
                          >
                            <option value={1}>{t('role_user')}</option>
                            <option value={10}>{t('role_admin')}</option>
                          </select>
                        ) : (
                          getRoleBadge(user.role)
                        )}
                      </td>
                      <td style={{ padding: '0.75rem 1rem' }}>
                        {canManageUser(user) ? (
                          <select
                            value={user.group}
                            onChange={(e) => handleSetGroup(user.id, e.target.value)}
                            style={{
                              padding: '0.25rem 0.5rem',
                              borderRadius: '4px',
                              border: '1px solid #e2e8f0',
                              fontSize: '0.75rem',
                              background: 'white',
                            }}
                          >
                            {groups.map((g) => (
                              <option key={g.id} value={g.name}>
                                {g.display_name || g.name}
                              </option>
                            ))}
                          </select>
                        ) : (
                          <span style={{ background: '#f1f5f9', padding: '0.25rem 0.5rem', borderRadius: '4px', fontSize: '0.75rem' }}>
                            {user.group}
                          </span>
                        )}
                      </td>
                      <td style={{ padding: '0.75rem 1rem' }}>
                        <span style={{ display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
                          {getStatusBadge(user.status)}
                          <span style={{ fontSize: '0.75rem' }}>{user.status === 1 ? t('status_active') : t('status_disabled')}</span>
                        </span>
                      </td>
                      <td style={{ padding: '0.75rem 1rem', textAlign: 'right' }}>
                        {canManageUser(user) && (
                          <div style={{ display: 'flex', gap: '0.5rem', justifyContent: 'flex-end' }}>
                            {user.status === 1 ? (
                              <button
                                onClick={() => handleManageUser(user.id, 'disable')}
                                className="btn btn-ghost"
                                style={{ fontSize: '0.75rem', padding: '0.25rem 0.5rem' }}
                              >
                                {t('disable')}
                              </button>
                            ) : (
                              <button
                                onClick={() => handleManageUser(user.id, 'enable')}
                                className="btn btn-ghost"
                                style={{ fontSize: '0.75rem', padding: '0.25rem 0.5rem', color: '#16a34a' }}
                              >
                                {t('enable')}
                              </button>
                            )}
                            <button
                              onClick={() => setShowResetPassword(user.id)}
                              className="btn btn-ghost"
                              style={{ fontSize: '0.75rem', padding: '0.25rem 0.5rem', color: '#2563eb' }}
                            >
                              {t('reset_password')}
                            </button>
                          </div>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {/* Groups Tab */}
        {activeTab === 'groups' && (
          <div>
            <div style={{ display: 'flex', justifyContent: 'flex-end', marginBottom: '1rem' }}>
              <button
                onClick={() => setShowCreateGroup(true)}
                className="btn btn-primary"
              >
                + {t('create_group')}
              </button>
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: '1rem' }}>
              {groups.map((group) => (
                <div key={group.id} className="bg-white rounded-xl border border-gray-200 p-4">
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '0.5rem' }}>
                    <div>
                      <h3 style={{ fontWeight: 600, color: '#0f172a' }}>{group.display_name || group.name}</h3>
                      <p style={{ fontSize: '0.75rem', color: '#64748b' }}>{group.name}</p>
                    </div>
                    <span style={{ 
                      fontSize: '0.75rem', 
                      padding: '0.125rem 0.5rem', 
                      borderRadius: '4px',
                      background: group.status === 1 ? '#dcfce7' : '#f1f5f9',
                      color: group.status === 1 ? '#166534' : '#64748b'
                    }}>
                      {group.status === 1 ? t('enabled') : t('disabled')}
                    </span>
                  </div>
                  {group.description && (
                    <p style={{ fontSize: '0.875rem', color: '#64748b', marginBottom: '0.5rem' }}>{group.description}</p>
                  )}
                  <div style={{ display: 'flex', gap: '1rem', fontSize: '0.75rem', color: '#64748b' }}>
                    <span>{t('quota')}: {group.quota.toLocaleString()}</span>
                    <span>{t('rate')}: {group.rate_limit}/min</span>
                  </div>
                  <div style={{ marginTop: '0.75rem', display: 'flex', gap: '0.5rem' }}>
                    <button
                      onClick={() => handleDeleteGroup(group.id)}
                      className="btn btn-ghost"
                      style={{ fontSize: '0.75rem', padding: '0.25rem 0.5rem', color: '#dc2626' }}
                    >
                      {t('delete')}
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Create User Modal */}
        {showCreateUser && (
          <div className="modal-overlay" onClick={() => setShowCreateUser(false)}>
            <div className="modal" onClick={(e) => e.stopPropagation()} style={{ maxWidth: '400px' }}>
              <div className="modal-header">
                <h2 className="modal-title">{t('create_user')}</h2>
                <button className="modal-close" onClick={() => setShowCreateUser(false)}>×</button>
              </div>
              <div className="modal-body">
                <div className="form-group">
                  <label className="form-label">{t('username')}</label>
                  <input
                    type="text"
                    className="input"
                    value={newUser.username}
                    onChange={(e) => setNewUser({ ...newUser, username: e.target.value })}
                    placeholder={t('enter_username')}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">{t('password')}</label>
                  <input
                    type="password"
                    className="input"
                    value={newUser.password}
                    onChange={(e) => setNewUser({ ...newUser, password: e.target.value })}
                    placeholder={t('enter_password')}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">{t('email')}</label>
                  <input
                    type="email"
                    className="input"
                    value={newUser.email}
                    onChange={(e) => setNewUser({ ...newUser, email: e.target.value })}
                    placeholder={t('enter_email')}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">{t('role')}</label>
                  <select
                    className="input"
                    value={newUser.role}
                    onChange={(e) => setNewUser({ ...newUser, role: parseInt(e.target.value) })}
                  >
                    <option value={1}>{t('role_user')}</option>
                    <option value={10}>{t('role_admin')}</option>
                  </select>
                </div>
                <button
                  onClick={handleCreateUser}
                  className="btn btn-primary"
                  style={{ width: '100%', marginTop: '1rem' }}
                >
                  {t('create')}
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Create Group Modal */}
        {showCreateGroup && (
          <div className="modal-overlay" onClick={() => setShowCreateGroup(false)}>
            <div className="modal" onClick={(e) => e.stopPropagation()} style={{ maxWidth: '400px' }}>
              <div className="modal-header">
                <h2 className="modal-title">{t('create_group')}</h2>
                <button className="modal-close" onClick={() => setShowCreateGroup(false)}>×</button>
              </div>
              <div className="modal-body">
                <div className="form-group">
                  <label className="form-label">{t('group_name')}</label>
                  <input
                    type="text"
                    className="input"
                    value={newGroup.name}
                    onChange={(e) => setNewGroup({ ...newGroup, name: e.target.value })}
                    placeholder="e.g., svip, vip"
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">{t('group_display_name')}</label>
                  <input
                    type="text"
                    className="input"
                    value={newGroup.display_name}
                    onChange={(e) => setNewGroup({ ...newGroup, display_name: e.target.value })}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">{t('description')}</label>
                  <input
                    type="text"
                    className="input"
                    value={newGroup.description}
                    onChange={(e) => setNewGroup({ ...newGroup, description: e.target.value })}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">{t('quota')}</label>
                  <input
                    type="number"
                    className="input"
                    value={newGroup.quota}
                    onChange={(e) => setNewGroup({ ...newGroup, quota: parseInt(e.target.value) || 0 })}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">{t('rate_limit')}</label>
                  <input
                    type="number"
                    className="input"
                    value={newGroup.rate_limit}
                    onChange={(e) => setNewGroup({ ...newGroup, rate_limit: parseInt(e.target.value) || 0 })}
                  />
                </div>
                <button
                  onClick={handleCreateGroup}
                  className="btn btn-primary"
                  style={{ width: '100%' }}
                >
                  {t('confirm')}
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Reset Password Modal */}
        {showResetPassword && (
          <div className="modal-overlay" onClick={() => setShowResetPassword(null)}>
            <div className="modal" onClick={(e) => e.stopPropagation()} style={{ maxWidth: '400px' }}>
              <div className="modal-header">
                <h2 className="modal-title">{t('reset_password')}</h2>
                <button className="modal-close" onClick={() => setShowResetPassword(null)}>×</button>
              </div>
              <div className="modal-body">
                <div className="form-group">
                  <label className="form-label">{t('new_password')}</label>
                  <input
                    type="password"
                    className="input"
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    placeholder={t('enter_new_password')}
                  />
                </div>
                <button
                  onClick={() => handleResetPassword(showResetPassword)}
                  className="btn btn-primary"
                  style={{ width: '100%' }}
                >
                  {t('reset_password')}
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
