import React, { useState } from "react";
import { createOrganization, createRole, createUser, login } from "../api/user";
import type { Organization, Role } from "../api/user";
import { getErrorMessage } from "../utils/errors";
import { useAuth } from "../hooks/useAuth";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
  CardFooter,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Icons } from "@/components/ui/icons";

export default function OnboardingPage() {
  const { login: doLogin } = useAuth();
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
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50 via-white to-cyan-100">
      <Card className="w-full max-w-lg shadow-2xl rounded-2xl border-0 bg-white/80 backdrop-blur-md">
        <CardHeader className="text-center pb-2">
          <Icons.user className="mx-auto mb-2 h-10 w-10 text-indigo-500" />
          <CardTitle className="text-3xl font-bold tracking-tight text-gray-900">
            Welcome to Unicorn Admin
          </CardTitle>
          <CardDescription className="text-gray-500 mt-2">
            Let's get your organization set up in a few easy steps.
          </CardDescription>
        </CardHeader>
        <CardContent>
          {error && (
            <Alert variant="destructive" className="mb-4">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}
          <div className="flex flex-col items-center">
            <div className="flex gap-2 mb-6">
              <div
                className={`h-2 w-8 rounded-full ${
                  step >= 1 ? "bg-indigo-500" : "bg-gray-200"
                }`}
              ></div>
              <div
                className={`h-2 w-8 rounded-full ${
                  step >= 2 ? "bg-indigo-500" : "bg-gray-200"
                }`}
              ></div>
              <div
                className={`h-2 w-8 rounded-full ${
                  step >= 3 ? "bg-indigo-500" : "bg-gray-200"
                }`}
              ></div>
            </div>
            {step === 1 && (
              <form onSubmit={handleOrg} className="w-full space-y-4">
                <Label htmlFor="orgName" className="text-gray-700">
                  Organization Name
                </Label>
                <Input
                  id="orgName"
                  value={orgName}
                  onChange={(e) => setOrgName(e.target.value)}
                  required
                  placeholder="e.g. Acme Inc."
                  className=""
                />
                <Button
                  type="submit"
                  className="w-full mt-2"
                  disabled={loading}
                >
                  {loading ? (
                    <Icons.spinner className="mr-2 h-4 w-4 animate-spin" />
                  ) : null}
                  Create Organization
                </Button>
              </form>
            )}
            {step === 2 && (
              <div className="w-full flex flex-col items-center gap-4">
                <div className="text-lg font-medium text-gray-700">
                  Organization created:{" "}
                  <span className="font-bold text-indigo-600">{org?.name}</span>
                </div>
                <Button
                  onClick={handleRole}
                  className="w-full"
                  disabled={loading}
                >
                  {loading ? (
                    <Icons.spinner className="mr-2 h-4 w-4 animate-spin" />
                  ) : null}
                  Create Admin Role
                </Button>
              </div>
            )}
            {step === 3 && (
              <form onSubmit={handleAdmin} className="w-full space-y-4">
                <div className="text-lg font-medium text-gray-700">
                  Admin Role created:{" "}
                  <span className="font-bold text-indigo-600">
                    {role?.name}
                  </span>
                </div>
                <Label htmlFor="adminName" className="text-gray-700">
                  Name
                </Label>
                <Input
                  id="adminName"
                  value={admin.name}
                  onChange={(e) =>
                    setAdmin((a) => ({ ...a, name: e.target.value }))
                  }
                  required
                  placeholder="Your Name"
                />
                <Label htmlFor="adminEmail" className="text-gray-700">
                  Email
                </Label>
                <Input
                  id="adminEmail"
                  type="email"
                  value={admin.email}
                  onChange={(e) =>
                    setAdmin((a) => ({ ...a, email: e.target.value }))
                  }
                  required
                  placeholder="admin@example.com"
                />
                <Label htmlFor="adminPassword" className="text-gray-700">
                  Password
                </Label>
                <Input
                  id="adminPassword"
                  type="password"
                  value={admin.password}
                  onChange={(e) =>
                    setAdmin((a) => ({ ...a, password: e.target.value }))
                  }
                  required
                  minLength={8}
                  placeholder="Create a password"
                />
                <Button
                  type="submit"
                  className="w-full mt-2"
                  disabled={loading}
                >
                  {loading ? (
                    <Icons.spinner className="mr-2 h-4 w-4 animate-spin" />
                  ) : null}
                  Create Admin User & Login
                </Button>
              </form>
            )}
          </div>
        </CardContent>
        <CardFooter className="flex flex-col items-center gap-2 pt-2">
          <span className="text-xs text-gray-400">Step {step} of 3</span>
        </CardFooter>
      </Card>
    </div>
  );
}
