"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { apiClient } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Building, User, Shield, ArrowRight, Check } from "lucide-react";
import Link from "next/link";

interface OnboardingStep {
  id: string;
  title: string;
  description: string;
  completed: boolean;
}

export default function OnboardingPage() {
  const [currentStep, setCurrentStep] = useState(0);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  // Form data
  const [organizationName, setOrganizationName] = useState("");
  const [adminName, setAdminName] = useState("");
  const [adminEmail, setAdminEmail] = useState("");
  const [adminPassword, setAdminPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  // Store actual IDs returned from API
  const [organizationId, setOrganizationId] = useState("");
  const [roleId, setRoleId] = useState("");

  const steps: OnboardingStep[] = [
    {
      id: "organization",
      title: "Create Organization",
      description: "Set up your organization",
      completed: false,
    },
    {
      id: "role",
      title: "Create Admin Role",
      description: "Create the default admin role",
      completed: false,
    },
    {
      id: "admin",
      title: "Create Admin User",
      description: "Create the default admin user",
      completed: false,
    },
  ];

  const handleCreateOrganization = async () => {
    if (!organizationName.trim()) {
      setError("Organization name is required");
      return;
    }

    setIsLoading(true);
    setError("");

    try {
      const response = await apiClient.createOrganization({
        name: organizationName,
      });
      setOrganizationId(response.organization.id);
      steps[0].completed = true;
      setCurrentStep(1);
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to create organization");
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateRole = async () => {
    setIsLoading(true);
    setError("");

    try {
      const response = await apiClient.createRole({
        name: `admin:${organizationName}`,
        permissions: [0, 1, 2], // Read, Write, Delete
      });
      setRoleId(response.role.id);
      steps[1].completed = true;
      setCurrentStep(2);
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to create admin role");
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateAdmin = async () => {
    if (!adminName.trim() || !adminEmail.trim() || !adminPassword.trim()) {
      setError("All fields are required");
      return;
    }

    if (adminPassword !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    if (adminPassword.length < 8) {
      setError("Password must be at least 8 characters long");
      return;
    }

    if (!organizationId || !roleId) {
      setError("Organization and role must be created first");
      return;
    }

    setIsLoading(true);
    setError("");

    try {
      await apiClient.createUser(organizationId, {
        name: adminName,
        email: adminEmail,
        password: adminPassword,
        role_id: roleId,
      });

      steps[2].completed = true;

      // Redirect to login
      router.push("/login?message=onboarding-complete");
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to create admin user");
    } finally {
      setIsLoading(false);
    }
  };

  const renderStep = () => {
    switch (currentStep) {
      case 0:
        return (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Building className="h-5 w-5" />
                Create Organization
              </CardTitle>
              <CardDescription>
                Set up your organization to get started with Unicorn Dashboard
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {error && (
                <Alert variant="destructive">
                  <AlertDescription>{error}</AlertDescription>
                </Alert>
              )}

              <div className="space-y-2">
                <Label htmlFor="org-name">Organization Name</Label>
                <Input
                  id="org-name"
                  value={organizationName}
                  onChange={(e) => setOrganizationName(e.target.value)}
                  placeholder="Enter your organization name"
                />
              </div>

              <Button
                onClick={handleCreateOrganization}
                disabled={isLoading}
                className="w-full"
              >
                {isLoading ? "Creating..." : "Create Organization"}
              </Button>
            </CardContent>
          </Card>
        );

      case 1:
        return (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Shield className="h-5 w-5" />
                Create Admin Role
              </CardTitle>
              <CardDescription>
                Create the default admin role with full permissions
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {error && (
                <Alert variant="destructive">
                  <AlertDescription>{error}</AlertDescription>
                </Alert>
              )}

              <div className="p-4 bg-blue-50 rounded-lg">
                <h4 className="font-medium text-blue-900">
                  Admin Role Permissions
                </h4>
                <ul className="mt-2 text-sm text-blue-700 space-y-1">
                  <li>• Read access to all resources</li>
                  <li>• Write access to all resources</li>
                  <li>• Delete access to all resources</li>
                </ul>
              </div>

              <Button
                onClick={handleCreateRole}
                disabled={isLoading}
                className="w-full"
              >
                {isLoading ? "Creating..." : "Create Admin Role"}
              </Button>
            </CardContent>
          </Card>
        );

      case 2:
        return (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <User className="h-5 w-5" />
                Create Admin User
              </CardTitle>
              <CardDescription>
                Create the default admin user with full permissions
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {error && (
                <Alert variant="destructive">
                  <AlertDescription>{error}</AlertDescription>
                </Alert>
              )}

              <div className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="admin-name">Admin Name</Label>
                  <Input
                    id="admin-name"
                    value={adminName}
                    onChange={(e) => setAdminName(e.target.value)}
                    placeholder="Enter admin name"
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="admin-email">Admin Email</Label>
                  <Input
                    id="admin-email"
                    type="email"
                    value={adminEmail}
                    onChange={(e) => setAdminEmail(e.target.value)}
                    placeholder="Enter admin email"
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="admin-password">Admin Password</Label>
                  <Input
                    id="admin-password"
                    type="password"
                    value={adminPassword}
                    onChange={(e) => setAdminPassword(e.target.value)}
                    placeholder="Enter admin password"
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="confirm-password">Confirm Password</Label>
                  <Input
                    id="confirm-password"
                    type="password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    placeholder="Confirm admin password"
                  />
                </div>
              </div>

              <Button
                onClick={handleCreateAdmin}
                disabled={isLoading}
                className="w-full"
              >
                {isLoading ? "Creating..." : "Create Admin User"}
              </Button>
            </CardContent>
          </Card>
        );

      default:
        return null;
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <div className="mx-auto h-12 w-12 flex items-center justify-center rounded-full bg-blue-100">
            <Building className="h-6 w-6 text-blue-600" />
          </div>
          <h2 className="mt-6 text-3xl font-extrabold text-gray-900">
            Welcome to Unicorn Dashboard
          </h2>
          <p className="mt-2 text-sm text-gray-600">
            Set up your organization and get started
          </p>
        </div>

        {/* Progress Steps */}
        <div className="flex items-center justify-between">
          {steps.map((step, index) => (
            <div key={step.id} className="flex items-center">
              <div
                className={`flex items-center justify-center w-8 h-8 rounded-full border-2 ${
                  index <= currentStep
                    ? "bg-blue-600 border-blue-600 text-white"
                    : "bg-white border-gray-300 text-gray-500"
                }`}
              >
                {step.completed ? (
                  <Check className="h-4 w-4" />
                ) : (
                  <span className="text-sm font-medium">{index + 1}</span>
                )}
              </div>
              {index < steps.length - 1 && (
                <div
                  className={`w-16 h-0.5 ${
                    index < currentStep ? "bg-blue-600" : "bg-gray-300"
                  }`}
                />
              )}
            </div>
          ))}
        </div>

        {renderStep()}

        <div className="text-center">
          <p className="text-sm text-gray-600">
            Already have an account?{" "}
            <Link
              href="/login"
              className="font-medium text-blue-600 hover:text-blue-500"
            >
              Sign in
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
