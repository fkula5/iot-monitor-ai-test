import React, { useState, useEffect } from 'react';
import { Trash2, UserCog, User, ShieldAlert } from 'lucide-react';

export const UserManagement = ({ token }) => {
  const [users, setUsers] = useState([]);
  const [error, setError] = useState(null);

  const fetchUsers = async () => {
    try {
      const res = await fetch('/api/users', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      if (!res.ok) throw new Error('Brak dostępu lub błąd serwera');
      const data = await res.json();
      setUsers(data);
      setError(null);
    } catch (e) {
      setError(e.message);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, [token]);

  const changeRole = async (userId, newRole) => {
    try {
      const res = await fetch(`/api/users/${userId}/role`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ role: newRole })
      });
      if (res.ok) {
        fetchUsers();
      } else {
        alert('Nie udało się zmienić roli');
      }
    } catch (e) {
      alert('Błąd sieci');
    }
  };

  const deleteUser = async (userId) => {
    if (!window.confirm('Czy na pewno chcesz usunąć tego użytkownika?')) return;
    try {
      const res = await fetch(`/api/users/${userId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      if (res.ok) {
        fetchUsers();
      } else {
        alert('Nie udało się usunąć użytkownika');
      }
    } catch (e) {
      alert('Błąd sieci');
    }
  };

  return (
    <div className="card list-card">
      <div style={{ padding: '20px', borderBottom: '1px solid var(--border-color)', display: 'flex', alignItems: 'center', gap: '10px' }}>
        <UserCog size={24} color="var(--primary)" />
        <h2 style={{ margin: 0 }}>Zarządzanie Użytkownikami</h2>
      </div>
      
      {error && (
        <div style={{ padding: '20px', color: 'var(--danger)', display: 'flex', alignItems: 'center', gap: '10px' }}>
          <ShieldAlert size={20} />
          {error}
        </div>
      )}

      <table className="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Użytkownik</th>
            <th>Rola</th>
            <th>Akcje</th>
          </tr>
        </thead>
        <tbody>
          {users.map(u => (
            <tr key={u.id}>
              <td className="text-muted">#{u.id}</td>
              <td className="fw-500">
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                  <User size={16} color="var(--text-muted)" />
                  {u.username}
                </div>
              </td>
              <td>
                <select 
                  value={u.role} 
                  onChange={(e) => changeRole(u.id, e.target.value)}
                  style={{
                    padding: '6px 12px',
                    borderRadius: '6px',
                    border: '1px solid var(--border-color)',
                    background: u.role === 'Admin' ? 'rgba(var(--primary-rgb), 0.1)' : 'var(--surface-color)',
                    color: u.role === 'Admin' ? 'var(--primary)' : 'var(--text-primary)',
                    fontWeight: '500'
                  }}
                  disabled={u.username === 'admin'}
                >
                  <option value="Viewer">Viewer</option>
                  <option value="Admin">Admin</option>
                </select>
              </td>
              <td>
                <button 
                  className="btn-icon" 
                  onClick={() => deleteUser(u.id)} 
                  title="Usuń konto"
                  disabled={u.username === 'admin'}
                  style={{ opacity: u.username === 'admin' ? 0.3 : 1 }}
                >
                  <Trash2 size={16} color="var(--danger)" />
                </button>
              </td>
            </tr>
          ))}
          {users.length === 0 && !error && (
            <tr>
              <td colSpan="4" style={{ textAlign: 'center', padding: '30px' }} className="text-muted">
                Brak użytkowników do wyświetlenia.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
};
