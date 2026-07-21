import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { DeviceList } from './DeviceList';

const mockDevices = [
  { id: 'dev-1', name: 'Sensor A', status: 'online', type: 'temperature', battery: 85, uptime: '10h' }
];

describe('DeviceList Component', () => {
  test('renders device table with data', () => {
    const permissions = { canWriteDevices: false };
    render(<DeviceList devices={mockDevices} permissions={permissions} userRole="Viewer" />);
    
    expect(screen.getByText('Sensor A')).toBeInTheDocument();
    expect(screen.getByText('temperature')).toBeInTheDocument();
  });

  test('hides action buttons if canWriteDevices is false', () => {
    const permissions = { canWriteDevices: false };
    render(<DeviceList devices={mockDevices} permissions={permissions} userRole="Viewer" />);
    
    expect(screen.queryByTitle('Włącz')).not.toBeInTheDocument();
    expect(screen.queryByTitle('Wyłącz')).not.toBeInTheDocument();
  });

  test('shows action buttons if canWriteDevices is true', () => {
    const permissions = { canWriteDevices: true };
    render(<DeviceList devices={mockDevices} permissions={permissions} userRole="Admin" />);
    
    expect(screen.getByTitle('Włącz')).toBeInTheDocument();
    expect(screen.getByTitle('Wyłącz')).toBeInTheDocument();
  });
});
