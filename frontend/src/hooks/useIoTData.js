import { useState, useEffect } from 'react';

export const useIoTData = (token) => {
  const [devices, setDevices] = useState([]);
  const [alerts, setAlerts] = useState([]);
  const [historyData, setHistoryData] = useState({});

  const getHeaders = () => ({
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  });

  const fetchInitialData = async () => {
    if (!token) return;
    try {
      const [devRes, histRes, alertRes] = await Promise.all([
        fetch('http://localhost:8080/api/devices', { headers: getHeaders() }),
        fetch('http://localhost:8080/api/history', { headers: getHeaders() }),
        fetch('http://localhost:8080/api/alerts', { headers: getHeaders() })
      ]);
      
      if (devRes.ok) setDevices(await devRes.json());
      if (histRes.ok) setHistoryData(await histRes.json());
      if (alertRes.ok) setAlerts(await alertRes.json());
    } catch (e) {
      console.error("Failed to fetch initial data", e);
    }
  };

  useEffect(() => {
    if (!token) {
      setDevices([]);
      setAlerts([]);
      setHistoryData({});
      return;
    }
    
    fetchInitialData();

    // 2. Setup WebSocket for real-time updates (pass token in query string)
    const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);
    
    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === 'update') {
        const { latest, alerts: newAlerts } = msg.data;
        
        // Update devices with latest values
        if (latest) {
          setDevices(prev => prev.map(dev => {
            if (latest[dev.id] !== undefined) {
              return { ...dev, value: latest[dev.id] };
            }
            return dev;
          }));
          
          // Append to history for charts
          setHistoryData(prev => {
             const nowTime = new Date().toLocaleTimeString([], {hour: '2-digit', minute:'2-digit', second:'2-digit'});
             const newHistory = { ...prev };
             for (const [devId, val] of Object.entries(latest)) {
                if (newHistory[devId]) {
                   newHistory[devId] = [...newHistory[devId].slice(1), { time: nowTime, value: val }];
                }
             }
             return newHistory;
          });
        }
        
        // Update alerts if provided
        if (newAlerts && newAlerts.length > 0) {
          setAlerts(newAlerts);
        }
      }
    };

    return () => ws.close();
  }, [token]);

  const addDevice = async (device) => {
    const res = await fetch('http://localhost:8080/api/devices', {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify(device)
    });
    if (res.ok) fetchInitialData();
  };

  const deleteDevice = async (id) => {
    const res = await fetch(`http://localhost:8080/api/devices/${id}`, {
      method: 'DELETE',
      headers: getHeaders()
    });
    if (res.ok) fetchInitialData();
  };

  return { devices, alerts, historyData, fetchInitialData, addDevice, deleteDevice };
};
