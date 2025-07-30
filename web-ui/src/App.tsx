import React, { useState } from "react";
import OnboardingPage from "./pages/OnboardingPage";
import LoginPage from "./pages/LoginPage";
import DashboardPage from "./pages/DashboardPage";
import { useAuth } from "./hooks/useAuth";
import "./App.css";

export default function App() {
  const { token } = useAuth();
  const [showLogin, setShowLogin] = useState(false);

  if (token) {
    return <DashboardPage />;
  }

  return (
    <div>
      {showLogin ? (
        <>
          <LoginPage />
          <button onClick={() => setShowLogin(false)} type="button">
            Back to Onboarding
          </button>
        </>
      ) : (
        <>
          <OnboardingPage />
          <button onClick={() => setShowLogin(true)} type="button">
            Already have an account? Login
          </button>
        </>
      )}
    </div>
  );
}
