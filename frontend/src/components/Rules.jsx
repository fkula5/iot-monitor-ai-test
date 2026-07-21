import React, { useState } from 'react';
import { Plus, Trash2, Activity } from 'lucide-react';

export const Rules = ({ rules, addRule, deleteRule, devices, permissions, userRole }) => {
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({
    deviceId: 'all',
    condition: '>',
    threshold: '',
    message: ''
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!formData.threshold || !formData.message) return;
    addRule({
      ...formData,
      threshold: parseFloat(formData.threshold)
    });
    setFormData({ deviceId: 'all', condition: '>', threshold: '', message: '' });
    setShowForm(false);
  };

  return (
    <div className="rules-container">
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '20px' }}>
        <h2>Reguły Logiczne (Dynamic Rule Engine)</h2>
        {permissions.canWriteRules && (
          <button className="btn-primary" onClick={() => setShowForm(!showForm)}>
            <Plus size={16} style={{ marginRight: '8px' }} /> Dodaj Regułę
          </button>
        )}
      </div>

      {showForm && (
        <form className="card" onSubmit={handleSubmit} style={{ marginBottom: '20px', padding: '20px' }}>
          <h3>Nowa reguła</h3>
          <div style={{ display: 'flex', gap: '15px', flexWrap: 'wrap', marginTop: '15px' }}>
            <div className="input-group" style={{ flex: '1', minWidth: '200px' }}>
              <label>Urządzenie</label>
              <select value={formData.deviceId} onChange={e => setFormData({...formData, deviceId: e.target.value})}>
                <option value="all">Wszystkie Urządzenia (all)</option>
                {devices.map(d => (
                  <option key={d.id} value={d.id}>{d.name} ({d.id})</option>
                ))}
              </select>
            </div>
            
            <div className="input-group" style={{ flex: '0.5', minWidth: '100px' }}>
              <label>Warunek</label>
              <select value={formData.condition} onChange={e => setFormData({...formData, condition: e.target.value})}>
                <option value=">">Większe niż (&gt;)</option>
                <option value="<">Mniejsze niż (&lt;)</option>
                <option value="==">Równe (==)</option>
              </select>
            </div>

            <div className="input-group" style={{ flex: '1', minWidth: '150px' }}>
              <label>Wartość (Threshold)</label>
              <input type="number" step="0.1" required value={formData.threshold} onChange={e => setFormData({...formData, threshold: e.target.value})} />
            </div>

            <div className="input-group" style={{ flex: '2', minWidth: '250px' }}>
              <label>Treść Alertu</label>
              <input type="text" placeholder="np. Zbyt wysoka temperatura!" required value={formData.message} onChange={e => setFormData({...formData, message: e.target.value})} />
            </div>
          </div>
          <div style={{ marginTop: '15px', display: 'flex', gap: '10px' }}>
            <button type="submit" className="btn-primary">Zapisz regułę</button>
            <button type="button" className="btn-secondary" onClick={() => setShowForm(false)}>Anuluj</button>
          </div>
        </form>
      )}

      <div className="grid">
        {rules.map(rule => (
          <div key={rule.id} className="card" style={{ position: 'relative' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '15px' }}>
              <div className="icon-wrapper" style={{ background: 'rgba(255,100,100,0.1)' }}>
                <Activity size={24} color="var(--danger)" />
              </div>
              <div>
                <h4 style={{ margin: 0 }}>Reguła #{rule.id}</h4>
                <small style={{ color: 'var(--text-secondary)' }}>Cel: {rule.deviceId === 'all' ? 'Wszystkie' : rule.deviceId}</small>
              </div>
            </div>
            
            <div style={{ padding: '10px', background: 'var(--bg-secondary)', borderRadius: '8px', marginBottom: '15px' }}>
              <code style={{ fontSize: '14px', color: 'var(--accent)' }}>
                IF wartość {rule.condition} {rule.threshold} THEN
              </code>
              <p style={{ margin: '5px 0 0 0', fontSize: '14px', fontWeight: 'bold' }}>"{rule.message}"</p>
            </div>

            {permissions.canWriteRules && (
              <button 
                style={{ position: 'absolute', top: '15px', right: '15px', background: 'none', border: 'none', color: 'var(--text-secondary)', cursor: 'pointer' }}
                onClick={() => deleteRule(rule.id)}
              >
                <Trash2 size={18} />
              </button>
            )}
          </div>
        ))}
        {rules.length === 0 && (
          <div className="card" style={{ gridColumn: '1 / -1', textAlign: 'center', padding: '40px', color: 'var(--text-secondary)' }}>
            Nie zdefiniowano żadnych reguł.
          </div>
        )}
      </div>
    </div>
  );
};
