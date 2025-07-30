import React, { useState } from "react";
import { createOrganization, createRole, createUser, login } from "../api/user";
import type { Organization, Role } from "../api/user";
import { getErrorMessage } from "../utils/errors";
import { useAuth } from "../hooks/useAuth";
import LogoutButton from "../components/LogoutButton";

export default function OnboardingPage() {
  const { login: doLogin, logout } = useAuth();
  const [step, setStep] = useState(1);
  const [orgName, setOrgName] = useState("");
  const [org, setOrg] = useState<Organization | null>(null);
  const [role, setRole] = useState<Role | null>(null);
  const [admin, setAdmin] = useState({ name: "", email: "", password: "" });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

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
      const res = await login(admin.email, admin.password);
      doLogin(res.token);
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
      <LogoutButton onLogout={logout} />
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
