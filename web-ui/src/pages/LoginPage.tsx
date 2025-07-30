import React, { useState } from "react";
import { login } from "../api/user";
import { getErrorMessage } from "../utils/errors";
import { useAuth } from "../hooks/useAuth";
import LogoutButton from "../components/LogoutButton";

export default function LoginPage() {
  const { login: doLogin, logout } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleLogin(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      const res = await login(email, password);
      doLogin(res.token);
    } catch (err) {
      setError(getErrorMessage(err, "Login failed"));
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="login-page">
      <h2>Login</h2>
      {error && <div style={{ color: "red" }}>{error}</div>}
      <LogoutButton onLogout={logout} />
      <form onSubmit={handleLogin}>
        <label>
          Email:
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </label>
        <label>
          Password:
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength={8}
          />
        </label>
        <button type="submit" disabled={loading}>
          Login
        </button>
      </form>
    </div>
  );
}
