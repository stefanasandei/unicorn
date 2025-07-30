import React, { useState } from "react";
import { createOrganization, createRole, createUser, login } from "./api/user";
import type { Organization, Role } from "./api/user";

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

export default function Onboarding({
  onLogin,
}: {
  onLogin: (token: string) => void;
}) {
  const [step, setStep] = useState(1);
  const [orgName, setOrgName] = useState("");
  const [org, setOrg] = useState<Organization | null>(null);
  const [role, setRole] = useState<Role | null>(null);
  const [admin, setAdmin] = useState({ name: "", email: "", password: "" });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  // Step 1: Create Organization
  async function handleOrg(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      const o = await createOrganization(orgName);
      setOrg(o);
      setStep(2);
    } catch (err) {
      setError(getErrorMessage(err, "Failed to create organization"));
    } finally {
      setLoading(false);
    }
  }

  // Step 2: Create Admin Role
  async function handleRole() {
    setLoading(true);
    setError("");
    try {
      const r = await createRole(`admin:${org?.id}`, [0, 1, 2]);
      setRole(r);
      setStep(3);
    } catch (err) {
      setError(getErrorMessage(err, "Failed to create role"));
    } finally {
      setLoading(false);
    }
  }

  // Step 3: Create Admin User
  async function handleAdmin(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      if (!org || !role) throw new Error("Missing org or role");
      await createUser(
        org.id,
        admin.name,
        admin.email,
        admin.password,
        role.id
      );
      // Login
      const res = await login(admin.email, admin.password);
      setCookie("token", res.token);
      onLogin(res.token);
    } catch (err) {
      setError(getErrorMessage(err, "Failed to create admin user"));
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="onboarding">
      <h2>Onboarding</h2>
      {error && <div style={{ color: "red" }}>{error}</div>}
      {step === 1 && (
        <form onSubmit={handleOrg}>
          <label>
            Organization Name:
            <input
              value={orgName}
              onChange={(e) => setOrgName(e.target.value)}
              required
            />
          </label>
          <button type="submit" disabled={loading}>
            Create Organization
          </button>
        </form>
      )}
      {step === 2 && (
        <div>
          <p>
            Organization created: <b>{org?.name}</b>
          </p>
          <button onClick={handleRole} disabled={loading}>
            Create Admin Role
          </button>
        </div>
      )}
      {step === 3 && (
        <form onSubmit={handleAdmin}>
          <p>
            Admin Role created: <b>{role?.name}</b>
          </p>
          <label>
            Name:
            <input
              value={admin.name}
              onChange={(e) =>
                setAdmin((a) => ({ ...a, name: e.target.value }))
              }
              required
            />
          </label>
          <label>
            Email:
            <input
              type="email"
              value={admin.email}
              onChange={(e) =>
                setAdmin((a) => ({ ...a, email: e.target.value }))
              }
              required
            />
          </label>
          <label>
            Password:
            <input
              type="password"
              value={admin.password}
              onChange={(e) =>
                setAdmin((a) => ({ ...a, password: e.target.value }))
              }
              required
              minLength={8}
            />
          </label>
          <button type="submit" disabled={loading}>
            Create Admin User & Login
          </button>
        </form>
      )}
    </div>
  );
}
