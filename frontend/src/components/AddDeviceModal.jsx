import React, { useState } from 'react';
import { X } from 'lucide-react';

export const AddDeviceModal = ({ onClose, onSave }) => {
  const [formData, setFormData] = useState({
    id: `dev-${Math.floor(Math.random() * 10000)}`,
    name: '',
    type: 'temperature',
    unit: '°C'
  });

  const handleTypeChange = (e) => {
    const type = e.target.value;
    let unit = '°C';
    if (type === 'humidity') unit = '%';
    if (type === 'pressure') unit = 'hPa';
    if (type === 'power') unit = 'W';
    setFormData({ ...formData, type, unit });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!formData.name) {
      alert("Proszę podać nazwę urządzenia.");
      return;
    }
    onSave({
      id: formData.id,
      name: formData.name,
      type: formData.type,
      unit: formData.unit,
      battery: 100, // default
      status: 'offline' // default
    });
  };

  return (
    <div className="modal-overlay" style={{
      position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
      backgroundColor: 'rgba(0,0,0,0.5)', zIndex: 1000,
      display: 'flex', alignItems: 'center', justifyContent: 'center',
      backdropFilter: 'blur(4px)'
    }}>
      <div className="card modal-content" style={{width: '100%', maxWidth: '450px', padding: '24px'}}>
        <div className="modal-header" style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px'}}>
          <h2 style={{margin: 0, fontSize: '1.25rem', fontWeight: '600'}}>Dodaj Nowy Czujnik</h2>
          <button className="btn-icon" onClick={onClose} style={{background:'transparent', border:'none', cursor:'pointer', color:'var(--text-muted)'}}><X size={20} /></button>
        </div>
        <form onSubmit={handleSubmit} style={{display: 'flex', flexDirection: 'column', gap: '16px'}}>
          <div className="input-group">
            <label style={{display: 'block', marginBottom: '8px', fontSize: '0.875rem', color: 'var(--text-muted)'}}>ID Urządzenia</label>
            <input type="text" value={formData.id} onChange={e => setFormData({...formData, id: e.target.value})} placeholder="np. dev-magazyn-1" style={{width: '100%', padding: '10px 12px', borderRadius: '8px', border: '1px solid var(--border-color)', background: 'var(--surface-color)', color: 'var(--text-primary)'}} />
          </div>
          <div className="input-group">
            <label style={{display: 'block', marginBottom: '8px', fontSize: '0.875rem', color: 'var(--text-muted)'}}>Nazwa Czujnika</label>
            <input type="text" value={formData.name} onChange={e => setFormData({...formData, name: e.target.value})} placeholder="np. Czujnik Magazyn B" autoFocus style={{width: '100%', padding: '10px 12px', borderRadius: '8px', border: '1px solid var(--border-color)', background: 'var(--surface-color)', color: 'var(--text-primary)'}} />
          </div>
          <div className="input-group">
            <label style={{display: 'block', marginBottom: '8px', fontSize: '0.875rem', color: 'var(--text-muted)'}}>Typ Sensoryczny</label>
            <select value={formData.type} onChange={handleTypeChange} style={{width: '100%', padding: '10px 12px', borderRadius: '8px', border: '1px solid var(--border-color)', background: 'var(--surface-color)', color: 'var(--text-primary)'}}>
              <option value="temperature">Temperatura</option>
              <option value="humidity">Wilgotność</option>
              <option value="pressure">Ciśnienie</option>
              <option value="power">Moc / Zasilanie</option>
              <option value="custom">Inny</option>
            </select>
          </div>
          <div className="input-group">
            <label style={{display: 'block', marginBottom: '8px', fontSize: '0.875rem', color: 'var(--text-muted)'}}>Jednostka Miary</label>
            <input type="text" value={formData.unit} onChange={e => setFormData({...formData, unit: e.target.value})} placeholder="np. °C, %, hPa" style={{width: '100%', padding: '10px 12px', borderRadius: '8px', border: '1px solid var(--border-color)', background: 'var(--surface-color)', color: 'var(--text-primary)'}} />
          </div>
          
          <div style={{display: 'flex', justifyContent: 'flex-end', gap: '12px', marginTop: '16px'}}>
            <button type="button" onClick={onClose} style={{padding: '10px 20px', borderRadius: '8px', border: '1px solid var(--border-color)', background: 'transparent', color: 'var(--text-primary)', cursor: 'pointer'}}>Anuluj</button>
            <button type="submit" className="btn-primary" style={{padding: '10px 20px', borderRadius: '8px', border: 'none', background: 'var(--primary)', color: 'white', cursor: 'pointer', fontWeight: '500'}}>Zapisz Urządzenie</button>
          </div>
        </form>
      </div>
    </div>
  );
};
