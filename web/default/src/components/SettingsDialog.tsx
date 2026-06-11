import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import api from '../lib/api';

interface SettingsDialogProps {
  open: boolean;
  onClose: () => void;
  userData: {
    id: number;
    username: string;
    display_name: string;
    email: string;
    tfa_enabled: boolean;
    passkey_enabled: boolean;
  } | null;
}

export function SettingsDialog({ open, onClose, userData }: SettingsDialogProps) {
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState<'profile' | 'security'>('profile');
  
  // Profile state
  const [displayName, setDisplayName] = useState('');
  const [email, setEmail] = useState('');
  
  // Password state
  const [originalPassword, setOriginalPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  
  // 2FA state
  const [tfaEnabled, setTfaEnabled] = useState(false);
  const [tfaCode, setTfaCode] = useState('');
  const [tfaSecret, setTfaSecret] = useState('');
  const [tfaQrCode, setTfaQrCode] = useState('');
  
  // Passkey state
  const [passkeyEnabled, setPasskeyEnabled] = useState(false);
  
  // UI state
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (userData) {
      setDisplayName(userData.display_name || '');
      setEmail(userData.email || '');
      setTfaEnabled(userData.tfa_enabled || false);
      setPasskeyEnabled(userData.passkey_enabled || false);
    }
  }, [userData]);

  if (!open) return null;

  const handleSaveProfile = async () => {
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      const response: any = await api.put('/user/self', {
        display_name: displayName,
        email: email,
      });
      if (response.success) {
        setSuccess(t('profile_updated'));
      } else {
        setError(response.message || t('update_failed'));
      }
    } catch (err: any) {
      setError(err.response?.data?.message || t('update_failed'));
    } finally {
      setLoading(false);
    }
  };

  const handleChangePassword = async () => {
    setError('');
    setSuccess('');

    if (!originalPassword) {
      setError(t('original_password_required'));
      return;
    }

    if (newPassword.length < 8) {
      setError(t('password_min_length'));
      return;
    }

    if (newPassword !== confirmPassword) {
      setError(t('passwords_not_match'));
      return;
    }

    setLoading(true);

    try {
      const response: any = await api.put('/user/self', {
        original_password: originalPassword,
        password: newPassword,
      });
      if (response.success) {
        setSuccess(t('password_updated'));
        setOriginalPassword('');
        setNewPassword('');
        setConfirmPassword('');
      } else {
        setError(response.message || t('update_failed'));
      }
    } catch (err: any) {
      setError(err.response?.data?.message || t('update_failed'));
    } finally {
      setLoading(false);
    }
  };

  const handleEnable2FA = async () => {
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      const response: any = await api.post('/user/2fa/enable');
      if (response.success) {
        setTfaSecret(response.data.secret);
        // Generate QR code image URL using QR code API
        const qrCodeUrl = `https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=${encodeURIComponent(response.data.qr_code)}`;
        setTfaQrCode(qrCodeUrl);
      } else {
        setError(response.message || t('update_failed'));
      }
    } catch (err: any) {
      setError(err.response?.data?.message || t('update_failed'));
    } finally {
      setLoading(false);
    }
  };

  const handleVerify2FA = async () => {
    setError('');
    setSuccess('');

    if (!tfaCode || tfaCode.length !== 6) {
      setError(t('enter_2fa_code'));
      return;
    }

    setLoading(true);

    try {
      const response: any = await api.post('/user/2fa/verify', { code: tfaCode });
      if (response.success) {
        setTfaEnabled(true);
        setTfaSecret('');
        setTfaQrCode('');
        setTfaCode('');
        setSuccess(t('2fa_enabled'));
      } else {
        setError(response.message || t('update_failed'));
      }
    } catch (err: any) {
      setError(err.response?.data?.message || t('update_failed'));
    } finally {
      setLoading(false);
    }
  };

  const handleDisable2FA = async () => {
    setError('');
    setSuccess('');

    if (!tfaCode || tfaCode.length !== 6) {
      setError(t('enter_2fa_code'));
      return;
    }

    setLoading(true);

    try {
      const response: any = await api.post('/user/2fa/disable', { code: tfaCode });
      if (response.success) {
        setTfaEnabled(false);
        setTfaCode('');
        setSuccess(t('2fa_disabled'));
      } else {
        setError(response.message || t('update_failed'));
      }
    } catch (err: any) {
      setError(err.response?.data?.message || t('update_failed'));
    } finally {
      setLoading(false);
    }
  };

  const handleEnablePasskey = async () => {
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      const response: any = await api.post('/user/passkey/enable');
      if (response.success) {
        setPasskeyEnabled(true);
        setSuccess(t('passkey_enabled'));
      } else {
        setError(response.message || t('update_failed'));
      }
    } catch (err: any) {
      setError(err.response?.data?.message || t('update_failed'));
    } finally {
      setLoading(false);
    }
  };

  const handleDisablePasskey = async () => {
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      const response: any = await api.delete('/user/passkey');
      if (response.success) {
        setPasskeyEnabled(false);
        setSuccess(t('passkey_disabled'));
      } else {
        setError(response.message || t('update_failed'));
      }
    } catch (err: any) {
      setError(err.response?.data?.message || t('update_failed'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()} style={{ maxWidth: '560px' }}>
        <div className="modal-header">
          <h2 className="modal-title">{t('settings')}</h2>
          <button className="modal-close" onClick={onClose}>
            <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div className="modal-body">
          {error && (
            <div className="alert alert-error" style={{ marginBottom: '1rem' }}>{error}</div>
          )}

          {success && (
            <div className="alert" style={{ 
              marginBottom: '1rem', 
              padding: '0.75rem 1rem', 
              background: '#f0fdf4', 
              color: '#166534', 
              borderRadius: '8px',
              border: '1px solid #bbf7d0'
            }}>{success}</div>
          )}

          {/* Tabs */}
          <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '1.5rem', borderBottom: '1px solid #e2e8f0', paddingBottom: '0.5rem' }}>
            <button
              onClick={() => setActiveTab('profile')}
              style={{
                padding: '0.5rem 1rem',
                borderRadius: '8px',
                border: 'none',
                background: activeTab === 'profile' ? '#667eea' : 'transparent',
                color: activeTab === 'profile' ? 'white' : '#64748b',
                cursor: 'pointer',
                fontWeight: 500,
                fontSize: '0.875rem',
              }}
            >
              {t('profile')}
            </button>
            <button
              onClick={() => setActiveTab('security')}
              style={{
                padding: '0.5rem 1rem',
                borderRadius: '8px',
                border: 'none',
                background: activeTab === 'security' ? '#667eea' : 'transparent',
                color: activeTab === 'security' ? 'white' : '#64748b',
                cursor: 'pointer',
                fontWeight: 500,
                fontSize: '0.875rem',
              }}
            >
              {t('two_factor_auth')}
            </button>
          </div>

          {activeTab === 'profile' && (
            <>
              {/* Profile Section */}
              <div className="modal-section">
                <div className="modal-section-title">{t('profile')}</div>
                <div className="form-group">
                  <label className="form-label">{t('display_name')}</label>
                  <input
                    type="text"
                    className="form-input"
                    value={displayName}
                    onChange={(e) => setDisplayName(e.target.value)}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">{t('email')}</label>
                  <input
                    type="email"
                    className="form-input"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                </div>
                <button
                  onClick={handleSaveProfile}
                  className="btn btn-primary"
                  disabled={loading}
                  style={{ width: '100%', marginTop: '0.5rem' }}
                >
                  {loading ? (
                    <>
                      <span className="spinner" />
                      {t('saving')}
                    </>
                  ) : (
                    t('save')
                  )}
                </button>
              </div>

              {/* Password Section */}
              <div className="modal-section">
                <div className="modal-section-title">{t('change_password')}</div>
                <div className="form-group">
                  <label className="form-label">{t('current_password')}</label>
                  <input
                    type="password"
                    className="form-input"
                    value={originalPassword}
                    onChange={(e) => setOriginalPassword(e.target.value)}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">{t('new_password')}</label>
                  <input
                    type="password"
                    className="form-input"
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">{t('confirm_password')}</label>
                  <input
                    type="password"
                    className="form-input"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                  />
                </div>
                <button
                  onClick={handleChangePassword}
                  className="btn btn-primary"
                  disabled={loading}
                  style={{ width: '100%', marginTop: '0.5rem' }}
                >
                  {t('change_password')}
                </button>
              </div>
            </>
          )}

          {activeTab === 'security' && (
            <>
              {/* 2FA Section */}
              <div className="modal-section">
                <div className="modal-section-title" style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                  <span>{t('two_factor_auth')}</span>
                  <span style={{ 
                    fontSize: '0.75rem', 
                    padding: '0.25rem 0.5rem', 
                    borderRadius: '4px',
                    background: tfaEnabled ? '#dcfce7' : '#f1f5f9',
                    color: tfaEnabled ? '#166534' : '#64748b'
                  }}>
                    {tfaEnabled ? t('enabled') : t('disabled')}
                  </span>
                </div>
                
                {!tfaEnabled && !tfaQrCode && (
                  <button
                    onClick={handleEnable2FA}
                    className="btn btn-primary"
                    disabled={loading}
                    style={{ width: '100%' }}
                  >
                    {t('enable_2fa')}
                  </button>
                )}

                {tfaQrCode && (
                  <div style={{ textAlign: 'center' }}>
                    <p style={{ fontSize: '0.875rem', color: '#64748b', marginBottom: '1rem' }}>
                      Scan this QR code with your authenticator app
                    </p>
                    <div style={{ 
                      background: 'white', 
                      padding: '1rem', 
                      borderRadius: '8px', 
                      display: 'inline-block',
                      marginBottom: '1rem'
                    }}>
                      <img src={tfaQrCode} alt="2FA QR Code" style={{ width: '150px', height: '150px' }} />
                    </div>
                    {tfaSecret && (
                      <div style={{ marginBottom: '1rem' }}>
                        <p style={{ fontSize: '0.75rem', color: '#64748b', marginBottom: '0.25rem' }}>
                          Or enter this secret manually:
                        </p>
                        <code style={{ 
                          background: '#f1f5f9', 
                          padding: '0.25rem 0.5rem', 
                          borderRadius: '4px',
                          fontSize: '0.75rem',
                          wordBreak: 'break-all'
                        }}>
                          {tfaSecret}
                        </code>
                      </div>
                    )}
                    <div className="form-group">
                      <label className="form-label">{t('enter_2fa_code')}</label>
                      <input
                        type="text"
                        className="form-input"
                        placeholder="000000"
                        value={tfaCode}
                        onChange={(e) => setTfaCode(e.target.value)}
                        maxLength={6}
                      />
                    </div>
                    <button
                      onClick={handleVerify2FA}
                      className="btn btn-primary"
                      disabled={loading}
                      style={{ width: '100%' }}
                    >
                      {t('verify_2fa')}
                    </button>
                  </div>
                )}

                {tfaEnabled && (
                  <div>
                    <p style={{ fontSize: '0.875rem', color: '#64748b', marginBottom: '1rem' }}>
                      Enter your 2FA code to disable
                    </p>
                    <div className="form-group">
                      <input
                        type="text"
                        className="form-input"
                        placeholder="000000"
                        value={tfaCode}
                        onChange={(e) => setTfaCode(e.target.value)}
                        maxLength={6}
                      />
                    </div>
                    <button
                      onClick={handleDisable2FA}
                      className="btn btn-outline"
                      disabled={loading}
                      style={{ width: '100%', color: '#dc2626', borderColor: '#fecaca' }}
                    >
                      {t('disable_2fa')}
                    </button>
                  </div>
                )}
              </div>

              {/* Passkey Section */}
              <div className="modal-section">
                <div className="modal-section-title" style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                  <span>{t('passkey')}</span>
                  <span style={{ 
                    fontSize: '0.75rem', 
                    padding: '0.25rem 0.5rem', 
                    borderRadius: '4px',
                    background: passkeyEnabled ? '#dcfce7' : '#f1f5f9',
                    color: passkeyEnabled ? '#166534' : '#64748b'
                  }}>
                    {passkeyEnabled ? t('enabled') : t('disabled')}
                  </span>
                </div>
                
                {!passkeyEnabled ? (
                  <button
                    onClick={handleEnablePasskey}
                    className="btn btn-primary"
                    disabled={loading}
                    style={{ width: '100%' }}
                  >
                    {t('enable_passkey')}
                  </button>
                ) : (
                  <button
                    onClick={handleDisablePasskey}
                    className="btn btn-outline"
                    disabled={loading}
                    style={{ width: '100%', color: '#dc2626', borderColor: '#fecaca' }}
                  >
                    {t('disable_passkey')}
                  </button>
                )}
              </div>

              {/* OAuth Section */}
              <div className="modal-section">
                <div className="modal-section-title">{t('oauth_login')}</div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
                  <button 
                    className="btn btn-outline" 
                    style={{ width: '100%', justifyContent: 'flex-start', gap: '0.75rem' }}
                    onClick={() => alert('GitHub OAuth will be available soon')}
                  >
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                    </svg>
                    {t('github')}
                  </button>
                  <button 
                    className="btn btn-outline" 
                    style={{ width: '100%', justifyContent: 'flex-start', gap: '0.75rem' }}
                    onClick={() => alert('Discord OAuth will be available soon')}
                  >
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M20.317 4.3698a19.7913 19.7913 0 00-4.8851-1.5152.0741.0741 0 00-.0785.0371c-.211.3753-.4447.8648-.6083 1.2495-1.8447-.2762-3.68-.2762-5.4868 0-.1636-.3933-.4058-.8742-.6177-1.2495a.077.077 0 00-.0785-.037 19.7363 19.7363 0 00-4.8852 1.515.0699.0699 0 00-.0321.0277C.5334 9.0458-.319 13.5799.0992 18.0578a.0824.0824 0 00.0312.0561c2.0528 1.5076 4.0413 2.4228 6.0016 3.0294a.0777.0777 0 00.0842-.0276c.4616-.6304.8731-1.2952 1.226-1.9942a.076.076 0 00-.0416-.1057c-.6528-.2476-1.2743-.5495-1.8722-.8923a.077.077 0 01-.0076-.1277c.1258-.0943.2517-.1923.3718-.2914a.0743.0743 0 01.0776-.0105c3.9278 1.7933 8.18 1.7933 12.0614 0a.0739.0739 0 01.0785.0095c.1202.099.246.1981.3728.2924a.077.077 0 01-.0066.1276 12.2986 12.2986 0 01-1.873.8914.0766.0766 0 00-.0407.1067c.3604.698.7719 1.3628 1.225 1.9932a.076.076 0 00.0842.0286c1.961-.6067 3.9495-1.5219 6.0023-3.0294a.077.077 0 00.0842-.0277c.3563-.6972.7713-1.362 1.225-1.9942a.076.076 0 00-.0417-.1067c-.598-.3428-1.22-.6478-1.8722-.8914a.077.077 0 01-.0077-.1277c.1202-.0992.2458-.1983.3716-.2925a.0743.0743 0 01.0776-.0105c2.6193 1.0302 5.1976 1.0302 7.8039 0 .262.0092.5231.0172.784.0172a.0765.0765 0 01.0785.0277c.3543.6972.7741 1.3616 1.1245 1.9942a.076.076 0 00.0842.0286c1.9617-.6067 3.9504-1.5219 6.0023-3.0294a.077.077 0 00.0842-.0276c.4616-.6304.8731-1.2952 1.226-1.9942a.0766.0766 0 00-.0417-.1067c-.598-.3428-1.22-.6478-1.8722-.8914a.077.077 0 01-.0077-.1277c.1202-.0992.2458-.1983.3716-.2925a.0743.0743 0 01.0776-.0105c2.6193 1.0302 5.1976 1.0302 7.8039 0 .262.0092.5231.0172.784.0172a.0765.0765 0 01.0785.0277c.3543.6972.7741 1.3616 1.1245 1.9942a.076.076 0 00.0842.0286c1.961-.6067 3.9495-1.5219 6.0023-3.0294a.077.077 0 00.0842-.0276c.3563-.6972.7713-1.362 1.225-1.9942a.0766.0766 0 00-.0417-.1067c-.598-.3428-1.22-.6478-1.8722-.8914a.077.077 0 01-.0077-.1277c.1202-.0992.2458-.1983.3716-.2925a.0743.0743 0 01.0776-.0105c2.6193 1.0302 5.1976 1.0302 7.8039 0a.0739.0739 0 01.0785.0095c.1202.099.246.1981.3728.2924a.077.077 0 01-.0066.1276z"/>
                    </svg>
                    {t('discord')}
                  </button>
                  <button 
                    className="btn btn-outline" 
                    style={{ width: '100%', justifyContent: 'flex-start', gap: '0.75rem' }}
                    onClick={() => alert('Google OAuth will be available soon')}
                  >
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M12.48 10.92v3.28h7.84c-.24 1.84-.853 3.187-1.787 4.133-1.147 1.147-2.933 2.4-6.053 2.4-4.827 0-8.6-3.893-8.6-8.72s3.773-8.72 8.6-8.72c2.6 0 4.507 1.027 5.907 2.347l2.307-2.307C18.747 1.44 16.133 0 12.48 0 5.867 0 .307 5.387.307 12s5.56 12 12.173 12c3.573 0 6.267-1.173 8.373-3.36 2.16-2.16 2.84-5.213 2.84-7.667 0-.76-.053-1.467-.173-2.053H12.48z"/>
                    </svg>
                    {t('google')}
                  </button>
                </div>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
