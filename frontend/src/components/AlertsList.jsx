import React from 'react';
import { AlertCircle, AlertTriangle } from 'lucide-react';

export const AlertsList = ({ alerts }) => {
  return (
    <div className="card list-card">
      <div className="alerts-container">
        {alerts.length === 0 ? (
          <p className="text-muted text-center" style={{padding: '2rem'}}>Brak aktywnych alertów.</p>
        ) : (
          alerts.map(alert => (
            <div key={alert.id} className={`alert-item ${alert.type}`}>
              <div className="alert-icon">
                {alert.type === 'error' ? <AlertCircle size={20} /> : <AlertTriangle size={20} />}
              </div>
              <div className="alert-content">
                <strong>{alert.deviceId}</strong>
                <p>{alert.message}</p>
                <span className="alert-time">{new Date(alert.timestamp).toLocaleString()}</span>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};
