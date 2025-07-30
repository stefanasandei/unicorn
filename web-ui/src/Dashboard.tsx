import { useEffect, useState } from "react";
import { getAccountInfo } from "./api/user";

function setCookie(name: string, value: string, days = 1) {
  const expires = new Date(Date.now() + days * 864e5).toUTCString();
  document.cookie = `${name}=${encodeURIComponent(
    value
  )}; expires=${expires}; path=/`;
}

export default function Dashboard() {
  const [info, setInfo] = useState<{
    roleName?: string;
    orgName?: string;
  } | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    const match = document.cookie.match(/token=([^;]+)/);
    if (!match) return;
    getAccountInfo(match[1])
      .then(setInfo)
      .catch(() => setError("Failed to load account info"));
  }, []);

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
      <button
        style={{ position: "absolute", top: 10, right: 10 }}
        onClick={() => {
          setCookie("token", "", -1); // Remove cookie
        }}
        type="button"
      >
        Logout
      </button>
    </div>
  );
}
