import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { Thermometer, Droplets, Activity } from 'lucide-react';

export const Dashboard = ({ devices, historyData }) => {
  const activeDevices = devices.filter(d => d.status === 'online').length;
  const avgTemp = devices
    .filter(d => d.type === 'temperature' && d.status === 'online')
    .reduce((acc, curr, _, arr) => acc + curr.value / arr.length, 0);

  return (
    <div className="dashboard-view">
      <div className="metrics-grid">
        <div className="card metric-card">
          <div className="metric-icon"><Activity size={24} color="var(--primary)" /></div>
          <div>
            <h3>Urządzenia Online</h3>
            <p className="metric-value">{activeDevices} / {devices.length}</p>
          </div>
        </div>
        <div className="card metric-card">
          <div className="metric-icon"><Thermometer size={24} color="var(--warning)" /></div>
          <div>
            <h3>Średnia Temperatura</h3>
            <p className="metric-value">{avgTemp ? avgTemp.toFixed(1) : '--'} °C</p>
          </div>
        </div>
      </div>

      <div className="charts-grid">
        <div className="card chart-card">
          <h3>Magazyn A - Temperatura (°C)</h3>
          <div className="chart-wrapper">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={historyData['dev-1']}>
                <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="var(--border-color)" />
                <XAxis dataKey="time" tick={{fontSize: 12, fill: 'var(--text-muted)'}} />
                <YAxis domain={['auto', 'auto']} tick={{fontSize: 12, fill: 'var(--text-muted)'}} />
                <Tooltip contentStyle={{ borderRadius: '8px', border: 'none', boxShadow: 'var(--shadow-md)' }} />
                <Line type="monotone" dataKey="value" stroke="var(--warning)" strokeWidth={3} dot={false} isAnimationActive={false} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>
        <div className="card chart-card">
          <h3>Chłodnia Sektor B - Temperatura (°C)</h3>
          <div className="chart-wrapper">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={historyData['dev-4']}>
                <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="var(--border-color)" />
                <XAxis dataKey="time" tick={{fontSize: 12, fill: 'var(--text-muted)'}} />
                <YAxis domain={['auto', 'auto']} tick={{fontSize: 12, fill: 'var(--text-muted)'}} />
                <Tooltip contentStyle={{ borderRadius: '8px', border: 'none', boxShadow: 'var(--shadow-md)' }} />
                <Line type="monotone" dataKey="value" stroke="var(--primary)" strokeWidth={3} dot={false} isAnimationActive={false} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>
    </div>
  );
};
