import React from "react";

export default function LogoutButton({ onLogout }: { onLogout: () => void }) {
  return (
    <button
      style={{ position: "absolute", top: 10, right: 10 }}
      onClick={onLogout}
      type="button"
    >
      Logout
    </button>
  );
}
