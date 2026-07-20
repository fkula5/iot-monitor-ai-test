import React, { useState, useEffect } from 'react';
import { LayoutDashboard, Server, Settings, AlertTriangle, LogOut, Plus, Trash2 } from 'lucide-react';
import './App.css';
import { useIoTData } from './hooks/useIoTData';
import { Dashboard } from './components/Dashboard';
import { DeviceList } from './components/DeviceList';
import { AlertsList } from './components/AlertsList';
import { Rules } from './components/Rules';

function App() {
  const [activeTab, setActiveTab] = useState('dashboard');
  const [token, setToken] = useState(localStorage.getItem('jwt_token') || null);
  const [timeRange, setTimeRange] = useState('-15m');
  const [loginForm, setLoginForm] = useState({ username: 'admin', password: 'admin123' });
  const [loginError, setLoginError] = useState('');
  const [toasts, setToasts] = useState([]);

  const { devices, alerts, rules, historyData, fetchInitialData, addDevice, deleteDevice, sendCommand, addRule, deleteRule } = useIoTData(token, timeRange);

  // Simple toast logic based on new alerts
  useEffect(() => {
    const validAlerts = alerts || [];
    if (validAlerts.length > 0) {
      const latestAlert = validAlerts[0]; // Alerts are prepended
      // Ensure we don't spam the same alert on load
      const isNew = Date.now() - (latestAlert.id / 1000000) < 5000;
      if (isNew) {
        setToasts(prev => [...prev, latestAlert]);
        setTimeout(() => setToasts(prev => prev.slice(1)), 4000);
      }
    }
  }, [alerts]);

  const getUnreadAlerts = () => (alerts || []).length;

  const handleLogin = async () => {
    setLoginError('');
    try {
      const res = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(loginForm)
      });
      if (!res.ok) {
        setLoginError('Nieprawidłowe dane logowania');
        return;
      }
      const data = await res.json();
      localStorage.setItem('jwt_token', data.token);
      setToken(data.token);
    } catch (e) {
      setLoginError('Błąd połączenia z serwerem logowania');
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('jwt_token');
    setToken(null);
  };

  if (!token) {
    return (
      <div className="login-container">
        <div className="card login-card">
          <h2>IoT Monitor</h2>
          <p className="subtitle">Zaloguj się do systemu</p>
          {loginError && <p style={{color: 'var(--danger)', marginBottom: '10px'}}>{loginError}</p>}
          <div className="input-group">
            <label>Login</label>
            <input 
              type="text" 
              value={loginForm.username}
              onChange={e => setLoginForm({...loginForm, username: e.target.value})}
            />
          </div>
          <div className="input-group">
            <label>Hasło</label>
            <input 
              type="password" 
              value={loginForm.password}
              onChange={e => setLoginForm({...loginForm, password: e.target.value})}
            />
          </div>
          <button className="btn-primary" onClick={handleLogin}>
            Zaloguj
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="app-layout">
      {/* Sidebar */}
      <aside className="sidebar">
        <div className="sidebar-header">
          <div className="logo-icon"><Server size={24} color="white" /></div>
          <h2>IoT Monitor</h2>
        </div>
        
        <nav className="sidebar-nav">
          <button className={`nav-item ${activeTab === 'dashboard' ? 'active' : ''}`} onClick={() => setActiveTab('dashboard')}>
            <LayoutDashboard size={20} />
            <span>Dashboard</span>
          </button>
          <button className={`nav-item ${activeTab === 'devices' ? 'active' : ''}`} onClick={() => setActiveTab('devices')}>
            <Server size={20} />
            <span>Urządzenia</span>
          </button>
          <button className={`nav-item ${activeTab === 'alerts' ? 'active' : ''}`} onClick={() => setActiveTab('alerts')}>
            <div style={{position: 'relative'}}>
              <AlertTriangle size={20} />
              {getUnreadAlerts() > 0 && <span style={{position:'absolute', top:'-4px', right:'-4px', background:'var(--danger)', width:'8px', height:'8px', borderRadius:'50%'}}></span>}
            </div>
            <span>Alerty</span>
          </button>
          <button className={`nav-item ${activeTab === 'rules' ? 'active' : ''}`} onClick={() => setActiveTab('rules')}>
            <Settings size={20} />
            <span>Reguły Engine</span>
          </button>
          <button className="nav-item">
            <Settings size={20} />
            <span>Ustawienia</span>
          </button>
        </nav>

        <div className="sidebar-footer">
          <button className="nav-item logout" onClick={handleLogout}>
            <LogOut size={20} />
            <span>Wyloguj</span>
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="main-content">
        <header className="top-header">
          <h1>
            {activeTab === 'dashboard' && 'Dashboard'}
            {activeTab === 'devices' && 'Zarządzanie Urządzeniami'}
            {activeTab === 'alerts' && 'Alerty Systemowe'}
          </h1>
          <div className="user-profile">
            <div className="avatar">A</div>
            <span>Admin</span>
          </div>
        </header>
        
        <div className="content-area">
          {activeTab === 'dashboard' && <Dashboard devices={devices} historyData={historyData} timeRange={timeRange} setTimeRange={setTimeRange} />}
          {activeTab === 'devices' && (
            <div>
              <div style={{display: 'flex', gap: '10px', marginBottom: '20px'}}>
                <button className="btn-primary" onClick={() => addDevice({id: 'dev-'+Date.now(), name: 'Nowy Czujnik', type: 'temperature', unit: '°C'})}>
                  <Plus size={16} style={{marginRight: '8px'}}/> Dodaj Urządzenie
                </button>
              </div>
              <DeviceList devices={devices} onDelete={deleteDevice} onCommand={sendCommand} />
            </div>
          )}
          {activeTab === 'alerts' && <AlertsList alerts={alerts} />}
          {activeTab === 'rules' && <Rules rules={rules} addRule={addRule} deleteRule={deleteRule} devices={devices} />}
        </div>
      </main>

      {/* Toasts */}
      <div style={{ position: 'fixed', bottom: '20px', right: '20px', display: 'flex', flexDirection: 'column', gap: '10px', zIndex: 9999 }}>
        {toasts.map((toast, idx) => (
          <div key={idx} style={{ background: 'var(--danger)', color: 'white', padding: '15px 20px', borderRadius: '8px', boxShadow: '0 4px 12px rgba(0,0,0,0.3)', animation: 'slideIn 0.3s ease-out' }}>
            <strong>Nowy Alert:</strong> {toast.message}
          </div>
        ))}
      </div>
    </div>
  );
}

export default App;
