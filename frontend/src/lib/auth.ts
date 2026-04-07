// Authentication utilities for TerasVPS Frontend

import { api } from './api';

export interface AuthState {
  isAuthenticated: boolean;
  user: any;
  token: string | null;
}

// Singleton auth state
let authState: AuthState = {
  isAuthenticated: false,
  user: null,
  token: null,
};

// Get auth state
export function getAuthState(): AuthState {
  return { ...authState };
}

// Set auth state
export function setAuthState(state: Partial<AuthState>) {
  authState = { ...authState, ...state };
}

// Login function
export async function login(email: string, password: string): Promise<{ success: boolean; error?: string }> {
  try {
    const response = await api.login({ email, password });

    // Set token
    api.setToken(response.data.token);

    // Update auth state
    setAuthState({
      isAuthenticated: true,
      user: response.data.user,
      token: response.data.token,
    });

    return { success: true };
  } catch (error: any) {
    return {
      success: false,
      error: error.message || 'Login failed',
    };
  }
}

// Register function
export async function register(
  username: string,
  email: string,
  password: string,
  phone?: string
): Promise<{ success: boolean; error?: string }> {
  try {
    const response = await api.register({ username, email, password, phone });

    // Set token
    api.setToken(response.data.token);

    // Update auth state
    setAuthState({
      isAuthenticated: true,
      user: response.data.user,
      token: response.data.token,
    });

    return { success: true };
  } catch (error: any) {
    return {
      success: false,
      error: error.message || 'Registration failed',
    };
  }
}

// Logout function
export async function logout(): Promise<void> {
  try {
    await api.logout();
  } catch (error) {
    // Ignore logout errors
    console.error('Logout error:', error);
  } finally {
    // Clear token
    api.removeToken();

    // Clear auth state
    setAuthState({
      isAuthenticated: false,
      user: null,
      token: null,
    });

    // Redirect to home
    if (typeof window !== 'undefined') {
      window.location.href = '/';
    }
  }
}

// Check if user is authenticated
export function checkAuth(): boolean {
  return authState.isAuthenticated;
}

// Get current user
export function getCurrentUser(): any {
  return authState.user;
}

// Get token
export function getToken(): string | null {
  return authState.token;
}

// Initialize auth state from localStorage
export function initAuthState(): void {
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;

  if (token) {
    api.setToken(token);

    // Decode JWT to get user info (simple implementation)
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      setAuthState({
        isAuthenticated: true,
        user: {
          id: payload.user_id,
          email: payload.email,
          role: payload.role,
        },
        token: token,
      });
    } catch (error) {
      // Invalid token, clear it
      api.removeToken();
      setAuthState({
        isAuthenticated: false,
        user: null,
        token: null,
      });
    }
  }
}
