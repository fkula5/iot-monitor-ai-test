import { useState, useEffect } from 'react';

export const useIoTData = (token, timeRange = '-15m') => {
  const [devices, setDevices] = useState([]);
  const [alerts, setAlerts] = useState([]);
  const [rules, setRules] = useState([]);
  const [historyData, setHistoryData] = useState({});

  const getHeaders = () => ({
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  });

  const fetchInitialData = async () => {
    if (!token) return;
    try {
      const [devRes, histRes, alertRes, rulesRes] = await Promise.all([
        fetch('/api/devices', { headers: getHeaders() }),
        fetch(`/api/history?range=${timeRange}`, { headers: getHeaders() }),
        fetch('/api/alerts', { headers: getHeaders() }),
        fetch('/api/rules', { headers: getHeaders() })
      ]);
      if (devRes.ok) setDevices((await devRes.json()) || []);
      if (histRes.ok) setHistoryData(await histRes.json());
      if (alertRes.ok) setAlerts((await alertRes.json()) || []);
      if (rulesRes.ok) setRules((await rulesRes.json()) || []);
    } catch (e) {
      console.error("Failed to fetch initial data", e);
    }
  };

  useEffect(() => {
    if (!token) {
      setDevices([]);
      setAlerts([]);
      setRules([]);
      setHistoryData({});
      return;
    }
    
    fetchInitialData();

    // 2. Setup WebSocket for real-time updates (pass token in query string)
    const ws = new WebSocket(`ws://${window.location.host}/ws?token=${token}`);
    
    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === 'update') {
        const { latest, alerts: newAlerts, devices: newDevices } = msg.data;
        
        // Update entire device list if provided (so we get status changes like online/offline)
        if (newDevices && newDevices.length > 0) {
           setDevices(prev => newDevices.map(dev => {
             // Keep the latest value if we already had it in the UI (or if it's in `latest`)
             const val = (latest && latest[dev.id] !== undefined) ? latest[dev.id] : 
                         (prev.find(p => p.id === dev.id)?.value);
             return { ...dev, value: val };
           }));
        } else if (latest) {
          // Fallback just in case: Update devices with latest values only
          setDevices(prev => prev.map(dev => {
            if (latest[dev.id] !== undefined) {
              return { ...dev, value: latest[dev.id] };
            }
            return dev;
          }));
        }
        
        if (latest) {
          // Append to history for charts
          setHistoryData(prev => {
             const nowTime = new Date().toLocaleTimeString([], {hour: '2-digit', minute:'2-digit', second:'2-digit'});
             const newHistory = { ...prev };
             for (const [devId, val] of Object.entries(latest)) {
                if (!newHistory[devId]) {
                   newHistory[devId] = [];
                }
                const currentHistory = newHistory[devId];
                const updatedHistory = [...currentHistory, { time: nowTime, value: val }];
                if (updatedHistory.length > 50) {
                   newHistory[devId] = updatedHistory.slice(1);
                } else {
                   newHistory[devId] = updatedHistory;
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
  }, [token, timeRange]);

  const addDevice = async (device) => {
    const res = await fetch('/api/devices', {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify(device)
    });
    if (res.ok) fetchInitialData();
  };

  const deleteDevice = async (id) => {
    const res = await fetch(`/api/devices/${id}`, {
      method: 'DELETE',
      headers: getHeaders()
    });
    if (res.ok) fetchInitialData();
  };

  const sendCommand = async (id, command) => {
    const res = await fetch(`/api/devices/${id}/command`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify({ command })
    });
    return res.ok;
  };

  const addRule = async (rule) => {
    const res = await fetch('/api/rules', {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify(rule)
    });
    if (res.ok) fetchInitialData();
  };

  const deleteRule = async (id) => {
    const res = await fetch(`/api/rules/${id}`, {
      method: 'DELETE',
      headers: getHeaders()
    });
    if (res.ok) fetchInitialData();
  };

  return {
    devices, alerts, rules, historyData,
    fetchInitialData, addDevice, deleteDevice, sendCommand,
    addRule, deleteRule
  };
};
