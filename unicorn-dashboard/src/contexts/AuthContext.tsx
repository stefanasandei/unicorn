"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { apiClient } from "@/lib/api";
import { User, LoginRequest } from "@/types/api";

interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

interface AuthProviderProps {
  children: React.ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const initializeAuth = async () => {
      const storedToken = localStorage.getItem("token");
      if (storedToken) {
        try {
          const validation = await apiClient.validateToken(storedToken);
          if (validation.valid) {
            setToken(storedToken);
            // Fetch user data
            const orgData = await apiClient.getOrganizations();
            // For now, we'll create a basic user object
            // In a real app, you'd have a user endpoint
            setUser({
              id: validation.claims.account_id,
              name: "User", // This would come from a user endpoint
              email: "user@example.com", // This would come from a user endpoint
              role: {
                id: validation.claims.role_id,
                name: "User",
                permissions: [],
                created_at: "",
                updated_at: "",
              },
              organization: {
                id: "",
                name: orgData.organization_name,
                created_at: "",
                updated_at: "",
              },
            });
          } else {
            localStorage.removeItem("token");
          }
        } catch (error) {
          console.error("Token validation failed:", error);
          localStorage.removeItem("token");
        }
      }
      setIsLoading(false);
    };

    initializeAuth();
  }, []);

  const login = async (credentials: LoginRequest) => {
    try {
      const response = await apiClient.login(credentials);
      localStorage.setItem("token", response.token);
      setToken(response.token);

      // Fetch user data after login
      const orgData = await apiClient.getOrganizations();
      setUser({
        id: "user-id", // This would come from the login response or a user endpoint
        name: "User",
        email: credentials.email,
        role: {
          id: "role-id",
          name: "User",
          permissions: [],
          created_at: "",
          updated_at: "",
        },
        organization: {
          id: "",
          name: orgData.organization_name,
          created_at: "",
          updated_at: "",
        },
      });
    } catch (error) {
      console.error("Login failed:", error);
      throw error;
    }
  };

  const logout = () => {
    localStorage.removeItem("token");
    setToken(null);
    setUser(null);
    window.location.href = "/";
  };

  const value: AuthContextType = {
    user,
    token,
    login,
    logout,
    isLoading,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
