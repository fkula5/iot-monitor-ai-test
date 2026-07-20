import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import userEvent from '@testing-library/user-event';
import App from './App';

// Mock fetch globally
global.fetch = vi.fn();

describe('App - Login & Registration', () => {
  beforeEach(() => {
    localStorage.clear();
    vi.clearAllMocks();
    global.fetch.mockResolvedValue({
      ok: true,
      json: async () => ({})
    });
  });

  test('renders login form when no token is present', () => {
    render(<App />);
    expect(screen.getByText('Zaloguj się do systemu')).toBeInTheDocument();
    expect(screen.getByText('Zaloguj')).toBeInTheDocument();
  });

  test('switches to registration form when clicking register button', () => {
    render(<App />);
    const toggleBtn = screen.getByText('Nie masz konta? Zarejestruj się.');
    fireEvent.click(toggleBtn);

    expect(screen.getByText('Utwórz nowe konto')).toBeInTheDocument();
    expect(screen.getByText('Zarejestruj się')).toBeInTheDocument();
  });

  test('handles successful login', async () => {
    global.fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ token: 'mock.token.part' })
    });

    render(<App />);
    const loginBtn = screen.getByText('Zaloguj');
    fireEvent.click(loginBtn);

    await waitFor(() => {
      expect(localStorage.getItem('jwt_token')).toBe('mock.token.part');
    });
  });
});
