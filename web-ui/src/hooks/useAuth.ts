import { useState, useEffect } from "react";
import { getCookie, setCookie, deleteCookie } from "../utils/cookies";

export function useAuth() {
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    const t = getCookie("token");
    if (t) setToken(t);
  }, []);

  function login(token: string) {
    setCookie("token", token);
    setToken(token);
    window.location.reload();
  }

  function logout() {
    deleteCookie("token");
    setToken(null);
    window.location.reload();
  }

  return { token, login, logout };
}
