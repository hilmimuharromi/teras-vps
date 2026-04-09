// API Client for TerasVPS Backend
// Uses fetch API with JWT authentication

import { apiTracker } from "./api-call-tracker";

const API_BASE_URL = import.meta.env.API_BASE_URL || "/api/v1";

// Types
export interface User {
  id: number;
  username: string;
  email: string;
  phone?: string;
  role: "customer" | "admin";
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  success: boolean;
  message: string;
  data: {
    token: string;
    user: User;
  };
}

export interface ErrorResponse {
  success: false;
  error: {
    code: string;
    message: string;
  };
}

// API Client Class
class APIClient {
  private baseURL: string;
  private token: string | null = null;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
    // Load token from localStorage
    if (typeof window !== "undefined") {
      this.token = localStorage.getItem("token");
    }
  }

  setToken(token: string) {
    this.token = token;
    if (typeof window !== "undefined") {
      localStorage.setItem("token", token);
    }
  }

  removeToken() {
    this.token = null;
    if (typeof window !== "undefined") {
      localStorage.removeItem("token");
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {},
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;

    const headers: HeadersInit = {
      "Content-Type": "application/json",
      ...options.headers,
    };

    // Always get the latest token from localStorage
    const token =
      typeof window !== "undefined"
        ? localStorage.getItem("token")
        : this.token;

    // Add Authorization header if token exists
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    // Track API call for debugging
    apiTracker.track(response, options.method || "GET", endpoint);

    // Log call ID in development
    if (import.meta.env?.DEV) {
      const callID = response.headers.get("X-API-Call-ID");
      if (callID) {
        console.log(
          `🔍 API [${callID.substring(0, 8)}] ${options.method || "GET"} ${endpoint} → ${response.status}`,
        );
      }
    }

    const data = await response.json();

    // Handle 401 Unauthorized - token expired or invalid
    if (response.status === 401 && typeof window !== "undefined") {
      // Clear invalid token
      localStorage.removeItem("token");
      // Redirect to login
      window.location.href = "/auth/login";
      throw new Error("Authentication required");
    }

    if (!response.ok) {
      throw new Error(data.error?.message || "An error occurred");
    }

    return data as T;
  }

  // Auth methods
  async register(data: {
    username: string;
    email: string;
    password: string;
    phone?: string;
  }): Promise<AuthResponse> {
    return this.request<AuthResponse>("/auth/register", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async login(data: {
    email: string;
    password: string;
  }): Promise<AuthResponse> {
    return this.request<AuthResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async logout(): Promise<{ success: boolean; message: string }> {
    return this.request("/auth/logout", {
      method: "POST",
    });
  }

  async getMe(): Promise<{ success: boolean; data: { user: User } }> {
    return this.request("/auth/me");
  }

  // User methods
  async getProfile(): Promise<{ success: boolean; data: { user: User } }> {
    return this.request("/user/profile");
  }

  async updateProfile(data: {
    phone?: string;
  }): Promise<{ success: boolean; message: string; data: { user: User } }> {
    return this.request("/user/profile", {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async changePassword(data: {
    old_password: string;
    new_password: string;
  }): Promise<{ success: boolean; message: string }> {
    return this.request("/user/password", {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }
}

// Export singleton instance
export const api = new APIClient(API_BASE_URL);

// Helper function to check if user is authenticated
export function isAuthenticated(): boolean {
  if (typeof window !== "undefined") {
    return !!localStorage.getItem("token");
  }
  return false;
}

// Helper function to get current user
export async function getCurrentUser(): Promise<User | null> {
  if (!isAuthenticated()) {
    return null;
  }

  try {
    const response = await api.getMe();
    return response.data.user;
  } catch (error) {
    return null;
  }
}
