import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../stores/authStore';
import { LanguageSwitcher } from '../components/LanguageSwitcher';
import { SettingsDialog } from '../components/SettingsDialog';
import api from '../lib/api';

interface DashboardPageProps {
  onNavigate: (page: string) => void;
}

interface UserData {
  id: number;
  username: string;
  display_name: string;
  role: number;
  email: string;
  quota: number;
  used_quota: number;
  request_count: number;
  group: string;
}

export function DashboardPage({ onNavigate }: DashboardPageProps) {
  const { t } = useTranslation();
  const { logout } = useAuthStore();
  const [userData, setUserData] = useState<UserData | null>(null);
  const [loading, setLoading] = useState(true);
  const [settingsOpen, setSettingsOpen] = useState(false);

  useEffect(() => {
    fetchUserData();
  }, []);

  const fetchUserData = async () => {
    try {
      const response: any = await api.get('/user/self');
      if (response.success) {
        setUserData(response.data);
      }
    } catch (err) {
      console.error('Failed to fetch user data');
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    try {
      await api.post('/user/logout');
    } catch (err) {
      // Ignore error
    }
    logout();
    onNavigate('login');
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
          <div className="navbar-logo">
            <span>T</span>
          </div>
          <span className="navbar-title">{t('app_name')}</span>
        </div>
        <div className="navbar-actions">
          <LanguageSwitcher />
          <div className="user-avatar">
            {(userData?.display_name || userData?.username || '').charAt(0).toUpperCase()}
          </div>
          <span style={{ fontSize: '0.875rem', fontWeight: 500, color: '#374151' }}>
            {userData?.display_name || userData?.username}
          </span>
          <button onClick={handleLogout} className="btn btn-ghost">
            {t('sign_out')}
          </button>
        </div>
      </nav>

      {/* Main Content */}
      <div className="main-content">
        {/* Welcome Banner */}
        <div className="welcome-banner">
          <h2>{t('welcome_back')}, {userData?.display_name || userData?.username}!</h2>
          <p>{t('dashboard_description')}</p>
        </div>

        {/* Stats Grid */}
        <div className="stats-grid">
          <div className="stat-card">
            <div className="stat-icon indigo">
              <svg width="24" height="24" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div className="stat-value">{(userData?.quota || 0).toLocaleString()}</div>
            <div className="stat-label">{t('quota')}</div>
            <div className="stat-sublabel">{t('used')}: {(userData?.used_quota || 0).toLocaleString()}</div>
          </div>

          <div className="stat-card">
            <div className="stat-icon green">
              <svg width="24" height="24" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <div className="stat-value">{(userData?.request_count || 0).toLocaleString()}</div>
            <div className="stat-label">{t('requests')}</div>
            <div className="stat-sublabel">{t('total_api_calls')}</div>
          </div>

          <div className="stat-card">
            <div className="stat-icon purple">
              <svg width="24" height="24" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
            </div>
            <div className="stat-value">{userData?.username}</div>
            <div className="stat-label">{t('account')}</div>
            <div className="stat-sublabel">{userData?.email || t('not_set')}</div>
          </div>
        </div>

        {/* Quick Actions */}
        <h3 className="section-title">{t('quick_actions')}</h3>
        <div className="actions-grid">
          <div className="action-card">
            <div className="action-icon blue">
              <svg width="22" height="22" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
              </svg>
            </div>
            <div className="action-title">{t('api_keys')}</div>
            <div className="action-desc">{t('manage_api_keys')}</div>
          </div>

          <div className="action-card">
            <div className="action-icon green">
              <svg width="22" height="22" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <div className="action-title">{t('usage_logs')}</div>
            <div className="action-desc">{t('view_usage_history')}</div>
          </div>

          <div className="action-card" onClick={() => setSettingsOpen(true)}>
            <div className="action-icon orange">
              <svg width="22" height="22" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </div>
            <div className="action-title">{t('settings')}</div>
            <div className="action-desc">{t('account_settings')}</div>
          </div>

          {userData && userData.role >= 10 && (
            <div className="action-card" onClick={() => onNavigate('admin')}>
              <div className="action-icon indigo">
                <svg width="22" height="22" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                </svg>
              </div>
              <div className="action-title">{t('admin_panel')}</div>
              <div className="action-desc">{t('manage_users_channels')}</div>
            </div>
          )}
        </div>
      </div>

      <SettingsDialog
        open={settingsOpen}
        onClose={() => setSettingsOpen(false)}
        userData={userData}
      />
    </div>
  );
}
