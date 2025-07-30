import React, { useState } from "react";
import { login } from "./api/user";

function setCookie(name: string, value: string, days = 1) {
  const expires = new Date(Date.now() + days * 864e5).toUTCString();
  document.cookie = `${name}=${encodeURIComponent(
    value
  )}; expires=${expires}; path=/`;
}

function getErrorMessage(err: unknown, fallback: string) {
  if (typeof err === "object" && err && "response" in err) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const response = (err as { response?: { data?: { error?: string } } })
      .response;
    if (response && response.data && typeof response.data.error === "string") {
      return response.data.error;
    }
  }
  return fallback;
}

export default function LoginPage({
  onLogin,
}: {
  onLogin: (token: string) => void;
}) {
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
      setCookie("token", res.token);
      onLogin(res.token);
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
