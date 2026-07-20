import React from 'react';
import { Wifi, WifiOff, Battery, Trash2 } from 'lucide-react';

export const DeviceList = ({ devices, onDelete }) => {
  return (
    <div className="card list-card">
      <table className="data-table">
        <thead>
          <tr>
            <th>Status</th>
            <th>Nazwa Urządzenia</th>
            <th>Typ</th>
            <th>Bateria</th>
            <th>Uptime</th>
            <th>Ostatni Odczyt</th>
            <th>Akcje</th>
          </tr>
        </thead>
        <tbody>
          {devices.map(dev => (
            <tr key={dev.id} className={dev.status === 'offline' ? 'row-offline' : ''}>
              <td>
                {dev.status === 'online' 
                  ? <div className="status-badge online"><Wifi size={14} /> Online</div>
                  : <div className="status-badge offline"><WifiOff size={14} /> Offline</div>
                }
              </td>
              <td className="fw-500">{dev.name}</td>
              <td className="text-muted">{dev.type}</td>
              <td>
                <div className="battery-indicator">
                  <Battery size={16} color={dev.battery > 20 ? 'var(--success)' : 'var(--danger)'} />
                  <span>{dev.battery}%</span>
                </div>
              </td>
              <td className="text-muted">{dev.uptime}</td>
              <td className="fw-600">
                {dev.value !== null ? `${dev.value} ${dev.unit}` : '-'}
              </td>
              <td>
                <button className="btn-icon" onClick={() => onDelete(dev.id)} title="Usuń urządzenie">
                  <Trash2 size={16} color="var(--danger)" />
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
