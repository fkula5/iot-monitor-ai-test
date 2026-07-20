import { render, screen, fireEvent } from '@testing-library/react';
import { expect, test, vi } from 'vitest';
import { Rules } from './Rules';

test('renders rules and can add a new one', () => {
  const mockRules = [
    { id: 1, deviceId: 'all', condition: '>', threshold: 50, message: 'Test message 1' }
  ];
  
  const mockDevices = [
    { id: 'dev-1', name: 'Test Sensor' }
  ];

  const addRuleMock = vi.fn();
  const deleteRuleMock = vi.fn();

  render(
    <Rules 
      rules={mockRules} 
      devices={mockDevices} 
      addRule={addRuleMock} 
      deleteRule={deleteRuleMock} 
    />
  );

  // Assert rule is displayed
  expect(screen.getByText(/"Test message 1"/i)).toBeDefined();
  
  // Open form
  const addBtn = screen.getByText(/Dodaj Regułę/);
  fireEvent.click(addBtn);

  // Fill form
  const messageInput = screen.getByPlaceholderText('np. Zbyt wysoka temperatura!');
  fireEvent.change(messageInput, { target: { value: 'Nowy Alert' } });
  
  const thresholdInput = screen.getByDisplayValue(''); // The only empty input is threshold
  fireEvent.change(thresholdInput, { target: { value: '45.5' } });
  
  const submitBtn = screen.getByText('Zapisz regułę');
  fireEvent.click(submitBtn);

  // Add rule mock should have been called
  expect(addRuleMock).toHaveBeenCalled();
});
