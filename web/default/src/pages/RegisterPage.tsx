import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { LanguageSwitcher } from '../components/LanguageSwitcher';
import api from '../lib/api';

interface RegisterPageProps {
  onNavigate: (page: string) => void;
}

export function RegisterPage({ onNavigate }: RegisterPageProps) {
  const { t } = useTranslation();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

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
      const response: any = await api.post('/user/register', {
        username,
        email,
        password,
      });
      if (response.success) {
        onNavigate('login');
      } else {
        setError(response.message || t('register_failed'));
      }
    } catch (err: any) {
      setError(err.response?.data?.message || t('register_failed'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-page">
      <LanguageSwitcher className="absolute" />
      
      <div className="auth-sidebar">
        <h1>{t('app_name')}</h1>
        <p>{t('register_subtitle')}</p>
        
        <div className="features">
          <div className="feature">
            <div className="feature-icon">
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
            </div>
            <span>{t('feature_free_register')}</span>
          </div>
          <div className="feature">
            <div className="feature-icon">
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
            </div>
            <span>{t('feature_security')}</span>
          </div>
          <div className="feature">
            <div className="feature-icon">
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064" />
              </svg>
            </div>
            <span>{t('feature_multi_model')}</span>
          </div>
        </div>
      </div>

      <div className="auth-container">
        <div className="auth-card">
          <div className="auth-header">
            <h2>{t('sign_up')}</h2>
            <p>{t('register_form_subtitle')}</p>
          </div>

          {error && (
            <div className="alert alert-error">{error}</div>
          )}

          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label className="label">{t('username')}</label>
              <input
                type="text"
                className="input"
                placeholder={t('enter_username')}
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                autoComplete="username"
                required
                minLength={3}
                maxLength={20}
              />
            </div>

            <div className="form-group">
              <label className="label">{t('email')}</label>
              <input
                type="email"
                className="input"
                placeholder={t('enter_email')}
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                autoComplete="email"
              />
            </div>

            <div className="form-group">
              <label className="label">{t('password')}</label>
              <input
                type="password"
                className="input"
                placeholder={t('enter_password')}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                autoComplete="new-password"
                required
                minLength={8}
              />
            </div>

            <div className="form-group">
              <label className="label">{t('confirm_password')}</label>
              <input
                type="password"
                className="input"
                placeholder={t('enter_confirm_password')}
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                autoComplete="new-password"
                required
              />
            </div>

            <button
              type="submit"
              className="btn btn-primary"
              disabled={loading}
              style={{ width: '100%' }}
            >
              {loading ? (
                <>
                  <span className="spinner" />
                  {t('signing_up')}
                </>
              ) : (
                t('sign_up')
              )}
            </button>
          </form>

          <div className="auth-footer">
            {t('already_have_account')}{' '}
            <button onClick={() => onNavigate('login')} className="link">
              {t('sign_in')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
