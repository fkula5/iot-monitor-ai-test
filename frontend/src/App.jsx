import React, { useState, useEffect } from 'react';
import { LayoutDashboard, Server, Settings, AlertTriangle, LogOut, Plus, Trash2 } from 'lucide-react';
import './App.css';
import { useIoTData } from './hooks/useIoTData';
import { Dashboard } from './components/Dashboard';
import { DeviceList } from './components/DeviceList';
import { AlertsList } from './components/AlertsList';

function App() {
  const [activeTab, setActiveTab] = useState('dashboard');
  const [token, setToken] = useState(localStorage.getItem('jwt_token') || null);
  const [timeRange, setTimeRange] = useState('-15m');
  const [loginForm, setLoginForm] = useState({ username: 'admin', password: 'admin123' });
  const [loginError, setLoginError] = useState('');

  const { devices, alerts, historyData, fetchInitialData, addDevice, deleteDevice } = useIoTData(token, timeRange);

  const getUnreadAlerts = () => alerts.length;

  const handleLogin = async () => {
    setLoginError('');
    try {
      const res = await fetch('http://localhost:8080/api/auth/login', {
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
              <DeviceList devices={devices} onDelete={deleteDevice} />
            </div>
          )}
          {activeTab === 'alerts' && <AlertsList alerts={alerts} />}
        </div>
      </main>
    </div>
  );
}

export default App;
