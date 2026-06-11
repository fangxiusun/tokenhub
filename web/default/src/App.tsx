import React, { useState, useEffect } from 'react';
import { useAuthStore } from './stores/authStore';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';
import { DashboardPage } from './pages/DashboardPage';
import { AdminPage } from './pages/AdminPage';
import { SetupPage } from './pages/SetupPage';
import api from './lib/api';

function App() {
  const { isAuthenticated, user } = useAuthStore();
  const [page, setPage] = useState(isAuthenticated ? 'dashboard' : 'login');
  const [setupRequired, setSetupRequired] = useState<boolean | null>(null);

  useEffect(() => {
    checkSetup();
  }, []);

  const checkSetup = async () => {
    try {
      const response: any = await api.get('/setup');
      if (response.success && !response.setup) {
        setSetupRequired(true);
        setPage('setup');
      } else {
        setSetupRequired(false);
      }
    } catch (err) {
      setSetupRequired(false);
    }
  };

  const handleNavigate = (newPage: string) => {
    setPage(newPage);
  };

  // Show loading while checking setup
  if (setupRequired === null) {
    return (
      <div className="auth-page" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <div className="spinner" style={{ width: '2rem', height: '2rem', borderColor: '#667eea', borderTopColor: 'transparent' }} />
      </div>
    );
  }

  // Show setup page if needed
  if (page === 'setup') {
    return <SetupPage onNavigate={handleNavigate} />;
  }

  if (page === 'register') {
    return <RegisterPage onNavigate={handleNavigate} />;
  }

  if (!isAuthenticated || page === 'login') {
    if (isAuthenticated && page === 'login') {
      return <DashboardPage onNavigate={handleNavigate} />;
    }
    return <LoginPage onNavigate={handleNavigate} />;
  }

  if (page === 'admin' && user && user.role >= 10) {
    return <AdminPage onNavigate={handleNavigate} />;
  }

  return <DashboardPage onNavigate={handleNavigate} />;
}

export default App;
