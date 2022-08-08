import { AxiosError } from "axios";
import React, { createContext, useContext, useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import {
  getMeQueryKey,
  me,
  useLogIn,
  useLogOut,
  useMe,
} from "../generated/auth/auth";
import {
  ControllersHTTPError,
  ControllersLogInRequest,
  ControllersUserInfo,
} from "../generated/models";

interface AuthContextValue {
  user: ControllersUserInfo | null;
  loading: boolean;
  logIn: (data: ControllersLogInRequest) => Promise<void>;
  logOut: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue>(null!);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<ControllersUserInfo | null>(null);
  const [loading, setLoading] = useState(true);

  const { data: meData } = useMe();
  const login = useLogIn({ mutation: { mutationKey: getMeQueryKey() } });
  const logout = useLogOut({ mutation: { mutationKey: getMeQueryKey() } });

  const navigate = useNavigate();
  const { state: navState } = useLocation();
  const state = navState as Record<string, string>;

  useEffect(() => {
    const loadUser = async () => {
      try {
        const { data: userData } = await me();
        setUser(userData);
        navigate(state?.path || "/");
      } finally {
        setLoading(false);
      }
    };

    loadUser();
  }, []);

  const value: AuthContextValue = {
    user,
    loading,
    logIn: async (data: ControllersLogInRequest) => {
      try {
        await login.mutateAsync({ data: data });
        setUser(meData?.data ?? null);
        navigate(state?.path || "/");
      } catch (error) {
        if (error instanceof AxiosError<ControllersHTTPError>)
          throw error.response?.data.message;
      }
    },
    logOut: async () => {
      await logout.mutateAsync();
      navigate("/login");
    },
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export default function useAuth() {
  return useContext(AuthContext);
}
