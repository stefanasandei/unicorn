import { useEffect, useState } from "react";
import { getAccountInfo } from "../api/user";
import { useAuth } from "../hooks/useAuth";
import LogoutButton from "../components/LogoutButton";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Icons } from "@/components/ui/icons";

export default function DashboardPage() {
  const { token, logout } = useAuth();
  const [info, setInfo] = useState<{
    roleName?: string;
    orgName?: string;
    permissions?: number[];
    name?: string;
    email?: string;
  } | null>(null);
  const [error, setError] = useState("");
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    if (!token) return;
    getAccountInfo(token)
      .then(setInfo)
      .catch(() => setError("Failed to load account info"));
  }, [token]);

  function handleCopyToken() {
    if (!token) return;
    navigator.clipboard.writeText(token);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  }

  return (
    <div className="min-h-screen flex bg-gradient-to-br from-indigo-50 via-white to-cyan-100">
      {/* Sidebar */}
      <aside className="hidden md:flex flex-col w-64 bg-white/80 backdrop-blur-md border-r shadow-lg p-6 min-h-screen">
        <div className="flex flex-col items-center gap-2 mb-8">
          <Icons.user className="h-10 w-10 text-indigo-500" />
          <span className="font-bold text-lg text-gray-900">
            {info?.name || "User"}
          </span>
          <span className="text-xs text-gray-500">{info?.email}</span>
        </div>
        <nav className="flex flex-col gap-4 mt-4">
          <span className="text-gray-700 font-medium">Dashboard</span>
          <button
            className="flex items-center gap-2 text-gray-600 hover:text-indigo-600 transition"
            onClick={logout}
          >
            <Icons.logout className="h-5 w-5" /> Logout
          </button>
        </nav>
        <div className="mt-auto text-xs text-gray-400 pt-8">
          © {new Date().getFullYear()} Unicorn Admin
        </div>
      </aside>
      {/* Main content */}
      <main className="flex-1 p-6 md:p-12">
        <div className="max-w-3xl mx-auto grid gap-8">
          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 mb-2">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 mb-1">
                Welcome{info?.name ? `, ${info.name}` : ""}!
              </h1>
              <p className="text-gray-500 text-lg">
                Here's your account overview and authentication token.
              </p>
            </div>
            <div className="md:hidden flex items-center gap-2">
              <Icons.user className="h-7 w-7 text-indigo-500" />
              <span className="font-bold text-gray-900">
                {info?.name || "User"}
              </span>
            </div>
            <LogoutButton onLogout={logout} />
          </div>
          {error ? (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          ) : (
            <div className="grid md:grid-cols-2 gap-8">
              <Card className="shadow-xl border-0 bg-white/90">
                <CardHeader>
                  <CardTitle>Account Information</CardTitle>
                  <CardDescription>
                    Details about your organization and role.
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div>
                    <h3 className="font-medium text-gray-600">Organization</h3>
                    <p className="text-2xl font-semibold text-indigo-700">
                      {info?.orgName}
                    </p>
                  </div>
                  <div>
                    <h3 className="font-medium text-gray-600">Role</h3>
                    <p className="text-2xl font-semibold text-indigo-700">
                      {info?.roleName?.split(":")[0]}:{" "}
                      {JSON.stringify(info?.permissions)}
                    </p>
                  </div>
                </CardContent>
              </Card>
              <Card className="shadow-xl border-0 bg-white/90">
                <CardHeader>
                  <CardTitle>Authentication Token</CardTitle>
                  <CardDescription>
                    Use this token for API access. Keep it secure.
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="relative group">
                    <div className="flex items-center gap-2 rounded-lg bg-gray-50 p-3 font-mono text-sm border border-gray-200">
                      <Icons.key className="h-4 w-4 text-gray-500" />
                      <code className="flex-1 break-all select-all">
                        {token}
                      </code>
                      <button
                        type="button"
                        className="ml-2 px-2 py-1 rounded bg-indigo-100 hover:bg-indigo-200 text-indigo-700 text-xs font-medium transition"
                        onClick={handleCopyToken}
                      >
                        {copied ? "Copied!" : "Copy"}
                      </button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
