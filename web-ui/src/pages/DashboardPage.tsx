import { useEffect, useState } from "react";
import { getAccountInfo } from "../api/user";
import { useAuth } from "../hooks/useAuth";
import LogoutButton from "../components/LogoutButton";

export default function DashboardPage() {
  const { token, logout } = useAuth();
  const [info, setInfo] = useState<{
    roleName?: string;
    orgName?: string;
  } | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!token) return;
    getAccountInfo(token)
      .then(setInfo)
      .catch(() => setError("Failed to load account info"));
  }, [token]);

  return (
    <div className="dashboard">
      <h2>Welcome to the Unicorn Admin Dashboard</h2>
      {info && (
        <>
          <p>
            Organization: <b>{info.orgName}</b>
          </p>
          <p>
            Role: <b>{info.roleName}</b>
          </p>
        </>
      )}
      {error && <div style={{ color: "red" }}>{error}</div>}
      <LogoutButton onLogout={logout} />
    </div>
  );
}
