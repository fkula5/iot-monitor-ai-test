import { render, screen } from '@testing-library/react';
import { expect, test, vi } from 'vitest';
import { Dashboard } from './Dashboard';

// Mock recharts because it uses ResizeObserver which jsdom doesn't support by default
vi.mock('recharts', () => {
  return {
    LineChart: () => <div>LineChart</div>,
    Line: () => <div>Line</div>,
    XAxis: () => <div>XAxis</div>,
    YAxis: () => <div>YAxis</div>,
    CartesianGrid: () => <div>CartesianGrid</div>,
    Tooltip: () => <div>Tooltip</div>,
    ResponsiveContainer: ({ children }) => <div>{children}</div>,
  };
});

test('renders dashboard with correct device counts', () => {
  const mockDevices = [
    { id: '1', name: 'Dev 1', type: 'temperature', status: 'online', value: 20 },
    { id: '2', name: 'Dev 2', type: 'humidity', status: 'offline', value: 50 },
    { id: '3', name: 'Dev 3', type: 'temperature', status: 'online', value: 30 }
  ];

  render(
    <Dashboard 
      devices={mockDevices} 
      historyData={{}} 
      timeRange="-15m" 
      setTimeRange={() => {}} 
    />
  );

  // 2 online devices out of 3 total
  expect(screen.getByText('2 / 3')).toBeDefined();
  
  // Average temp of online temperature devices: (20 + 30) / 2 = 25
  expect(screen.getByText('25.0 °C')).toBeDefined();
});
