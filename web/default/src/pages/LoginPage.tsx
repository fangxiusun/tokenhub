import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../stores/authStore';
import { LanguageSwitcher } from '../components/LanguageSwitcher';
import api from '../lib/api';

interface LoginPageProps {
  onNavigate: (page: string) => void;
}

export function LoginPage({ onNavigate }: LoginPageProps) {
  const { t } = useTranslation();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuthStore();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const response: any = await api.post('/user/login', { username, password });
      if (response.success) {
        login(response.data, response.data.token);
        onNavigate('dashboard');
      } else {
        setError(response.message || t('login_failed'));
      }
    } catch (err: any) {
      setError(err.response?.data?.message || t('login_failed'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-page">
      <LanguageSwitcher className="absolute" />
      
      <div className="auth-sidebar">
        <h1>{t('app_name')}</h1>
        <p>{t('app_description')}</p>
        
        <div className="features">
          <div className="feature">
            <div className="feature-icon">
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <span>{t('feature_proxy')}</span>
          </div>
          <div className="feature">
            <div className="feature-icon">
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
            </div>
            <span>{t('feature_security')}</span>
          </div>
          <div className="feature">
            <div className="feature-icon">
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <span>{t('feature_stats')}</span>
          </div>
        </div>
      </div>

      <div className="auth-container">
        <div className="auth-card">
          <div className="auth-header">
            <h2>{t('sign_in')}</h2>
            <p>{t('login_subtitle')}</p>
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
                autoComplete="username"
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
                autoComplete="current-password"
                required
              />
            </div>

            <button
              type="submit"
              className="btn btn-primary"
              disabled={loading}
            >
              {loading ? (
                <>
                  <span className="spinner" />
                  {t('signing_in')}
                </>
              ) : (
                t('sign_in')
              )}
            </button>
          </form>

          <div className="auth-footer">
            {t('dont_have_account')}{' '}
            <button onClick={() => onNavigate('register')} className="auth-link">
              {t('sign_up')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
