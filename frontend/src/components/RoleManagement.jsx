import React, { useState, useEffect } from 'react';
import { Trash2, ShieldCheck, Plus, Check, X } from 'lucide-react';

export const RoleManagement = ({ token }) => {
  const [roles, setRoles] = useState([]);
  const [error, setError] = useState(null);
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    canWriteDevices: false,
    canWriteRules: false,
    canManageUsers: false
  });

  const fetchRoles = async () => {
    try {
      const res = await fetch('/api/roles', {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      if (!res.ok) throw new Error('Brak dostępu lub błąd serwera');
      const data = await res.json();
      setRoles(data);
      setError(null);
    } catch (e) {
      setError(e.message);
    }
  };

  useEffect(() => {
    fetchRoles();
  }, [token]);

  const handleCreate = async (e) => {
    e.preventDefault();
    if (!formData.name) return;
    
    try {
      const res = await fetch('/api/roles', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
      });
      if (res.ok) {
        fetchRoles();
        setShowForm(false);
        setFormData({ name: '', canWriteDevices: false, canWriteRules: false, canManageUsers: false });
      } else {
        const d = await res.json();
        alert(d.error || 'Nie udało się utworzyć roli');
      }
    } catch (e) {
      alert('Błąd sieci');
    }
  };

  const deleteRole = async (roleName) => {
    if (!window.confirm(`Czy na pewno chcesz usunąć rolę ${roleName}?`)) return;
    try {
      const res = await fetch(`/api/roles/${roleName}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      });
      if (res.ok) {
        fetchRoles();
      } else {
        const d = await res.json();
        alert(d.error || 'Nie udało się usunąć roli');
      }
    } catch (e) {
      alert('Błąd sieci');
    }
  };

  return (
    <div className="card list-card" style={{ marginTop: '20px' }}>
      <div style={{ padding: '20px', borderBottom: '1px solid var(--border-color)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
          <ShieldCheck size={24} color="var(--primary)" />
          <h2 style={{ margin: 0 }}>Zarządzanie Rolami (RBAC)</h2>
        </div>
        <button className="btn-primary" onClick={() => setShowForm(!showForm)}>
          <Plus size={16} style={{ marginRight: '8px' }} /> Nowa Rola
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleCreate} style={{ padding: '20px', borderBottom: '1px solid var(--border-color)', background: 'var(--bg-secondary)' }}>
          <div style={{ display: 'flex', gap: '20px', flexWrap: 'wrap', alignItems: 'flex-end' }}>
            <div className="input-group" style={{ flex: '1', minWidth: '200px', marginBottom: 0 }}>
              <label>Nazwa Roli (np. Technik)</label>
              <input 
                type="text" 
                required 
                value={formData.name} 
                onChange={e => setFormData({...formData, name: e.target.value})} 
              />
            </div>
            
            <div style={{ display: 'flex', gap: '15px', flexWrap: 'wrap', marginBottom: '10px' }}>
              <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
                <input 
                  type="checkbox" 
                  checked={formData.canWriteDevices} 
                  onChange={e => setFormData({...formData, canWriteDevices: e.target.checked})} 
                /> 
                Zarządzanie Urządzeniami (Włącz/Wyłącz)
              </label>
              <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
                <input 
                  type="checkbox" 
                  checked={formData.canWriteRules} 
                  onChange={e => setFormData({...formData, canWriteRules: e.target.checked})} 
                /> 
                Edycja Reguł
              </label>
              <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
                <input 
                  type="checkbox" 
                  checked={formData.canManageUsers} 
                  onChange={e => setFormData({...formData, canManageUsers: e.target.checked})} 
                /> 
                Zarządzanie Kontami
              </label>
            </div>
            
            <div style={{ display: 'flex', gap: '10px' }}>
              <button type="submit" className="btn-primary">Utwórz</button>
              <button type="button" className="btn-secondary" onClick={() => setShowForm(false)}>Anuluj</button>
            </div>
          </div>
        </form>
      )}

      {error && <div style={{ padding: '20px', color: 'var(--danger)' }}>{error}</div>}

      <table className="data-table">
        <thead>
          <tr>
            <th>Nazwa Roli</th>
            <th style={{ textAlign: 'center' }}>Urządzenia (W/W)</th>
            <th style={{ textAlign: 'center' }}>Reguły</th>
            <th style={{ textAlign: 'center' }}>Użytkownicy</th>
            <th>Akcje</th>
          </tr>
        </thead>
        <tbody>
          {roles.map(r => (
            <tr key={r.id}>
              <td className="fw-600">{r.name}</td>
              <td style={{ textAlign: 'center' }}>
                {r.canWriteDevices ? <Check size={18} color="var(--success)"/> : <X size={18} color="var(--danger)"/>}
              </td>
              <td style={{ textAlign: 'center' }}>
                {r.canWriteRules ? <Check size={18} color="var(--success)"/> : <X size={18} color="var(--danger)"/>}
              </td>
              <td style={{ textAlign: 'center' }}>
                {r.canManageUsers ? <Check size={18} color="var(--success)"/> : <X size={18} color="var(--danger)"/>}
              </td>
              <td>
                <button 
                  className="btn-icon" 
                  onClick={() => deleteRole(r.name)} 
                  title="Usuń rolę"
                  disabled={r.name === 'Admin' || r.name === 'Viewer'}
                  style={{ opacity: (r.name === 'Admin' || r.name === 'Viewer') ? 0.3 : 1 }}
                >
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
