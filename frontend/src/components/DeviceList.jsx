import React from 'react';
import { Wifi, WifiOff, Battery, Trash2, Power, RotateCw, PowerOff } from 'lucide-react';

export const DeviceList = ({ devices, onDelete, onCommand, permissions, userRole }) => {
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
            {permissions.canWriteDevices && <th>Akcje</th>}
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
              {permissions.canWriteDevices && (
                <td>
                  <div style={{ display: 'flex', gap: '8px' }}>
                    <button className="btn-icon" onClick={() => onCommand(dev.id, 'TURN_ON')} title="Włącz">
                      <Power size={16} color="var(--success)" />
                    </button>
                    <button className="btn-icon" onClick={() => onCommand(dev.id, 'TURN_OFF')} title="Wyłącz">
                      <PowerOff size={16} color="var(--danger)" />
                    </button>
                    <button className="btn-icon" onClick={() => onCommand(dev.id, 'RESTART')} title="Zrestartuj">
                      <RotateCw size={16} color="var(--primary)" />
                    </button>
                    <button className="btn-icon" onClick={() => onDelete(dev.id)} title="Usuń urządzenie">
                      <Trash2 size={16} color="var(--text-muted)" />
                    </button>
                  </div>
                </td>
              )}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
