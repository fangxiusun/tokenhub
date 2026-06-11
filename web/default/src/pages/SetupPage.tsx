import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { LanguageSwitcher } from '../components/LanguageSwitcher';
import api from '../lib/api';

interface SetupPageProps {
  onNavigate: (page: string) => void;
}

export function SetupPage({ onNavigate }: SetupPageProps) {
  const { t } = useTranslation();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [setupComplete, setSetupComplete] = useState(false);

  useEffect(() => {
    checkSetup();
  }, []);

  const checkSetup = async () => {
    try {
      const response: any = await api.get('/setup');
      if (response.success && response.setup) {
        setSetupComplete(true);
      }
    } catch (err) {
      console.error('Failed to check setup');
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (password !== confirmPassword) {
      setError(t('passwords_not_match'));
      return;
    }

    if (password.length < 8) {
      setError(t('password_min_length'));
      return;
    }

    setLoading(true);

    try {
      const response: any = await api.post('/setup', { username, password });
      if (response.success) {
        setSetupComplete(true);
      } else {
        setError(response.message || 'Setup failed');
      }
    } catch (err: any) {
      setError(err.response?.data?.message || 'Setup failed');
    } finally {
      setLoading(false);
    }
  };

  if (setupComplete) {
    return (
      <div className="auth-page">
        <LanguageSwitcher className="absolute" />
        <div className="auth-sidebar">
          <h1>{t('app_name')}</h1>
          <p>{t('setup_complete_desc')}</p>
        </div>
        <div className="auth-container">
          <div className="auth-card">
            <div className="auth-header">
              <h2>{t('setup_complete')}</h2>
              <p>{t('root_user_created')}</p>
            </div>
            <button
              onClick={() => onNavigate('login')}
              className="btn btn-primary"
              style={{ width: '100%' }}
            >
              {t('go_to_login')}
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="auth-page">
      <LanguageSwitcher className="absolute" />
      <div className="auth-sidebar">
        <h1>{t('app_name')}</h1>
        <p>{t('setup_subtitle')}</p>
        <div className="features">
          <div className="feature">
            <div className="feature-icon">
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
            </div>
            <span>{t('feature_root_all_permissions')}</span>
          </div>
          <div className="feature">
            <div className="feature-icon">
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
            </div>
            <span>{t('feature_config_admin_user')}</span>
          </div>
        </div>
      </div>
      <div className="auth-container">
        <div className="auth-card">
          <div className="auth-header">
            <h2>{t('system_setup')}</h2>
            <p>{t('create_root_admin')}</p>
          </div>

          {error && (
            <div className="alert alert-error">{error}</div>
          )}

          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label className="form-label">{t('username')}</label>
              <input
                type="text"
                className="form-input"
                placeholder={t('enter_username')}
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
              />
            </div>

            <div className="form-group">
              <label className="form-label">{t('password')}</label>
              <input
                type="password"
                className="form-input"
                placeholder={t('enter_password')}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                minLength={8}
              />
            </div>

            <div className="form-group">
              <label className="form-label">{t('confirm_password')}</label>
              <input
                type="password"
                className="form-input"
                placeholder={t('enter_confirm_password')}
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
              />
            </div>

            <button
              type="submit"
              className="btn btn-primary"
              disabled={loading}
              style={{ width: '100%' }}
            >
              {loading ? 'Creating...' : 'Create Root User'}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
