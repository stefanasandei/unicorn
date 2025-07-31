"use client";

import React, { useEffect, useState } from "react";
import { Layout } from "@/components/Layout";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Shield,
  Users,
  Building,
  Plus,
  Edit,
  Trash2,
  UserCheck,
  Key,
  Crown,
  Settings,
  AlertCircle,
} from "lucide-react";
import { apiClient } from "@/lib/api";
import { Role, Organization } from "@/types/api";

export default function IAMPage() {
  const [roles, setRoles] = useState<Role[]>([]);
  const [organization, setOrganization] = useState<{
    organization_name: string;
    users: Array<{
      id: string;
      name: string;
      email: string;
      role_id?: string;
    }>;
  } | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  // Form states
  const [newRoleName, setNewRoleName] = useState("");
  const [newRolePermissions, setNewRolePermissions] = useState<number[]>([]);
  const [newUserName, setNewUserName] = useState("");
  const [newUserEmail, setNewUserEmail] = useState("");
  const [newUserPassword, setNewUserPassword] = useState("");
  const [selectedRoleId, setSelectedRoleId] = useState("");

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setIsLoading(true);
      const [rolesData, orgData] = await Promise.all([
        apiClient.getRoles(),
        apiClient.getOrganizations(),
      ]);
      setRoles(rolesData.roles);
      setOrganization(orgData);
    } catch (err: unknown) {
      const error = err as { response?: { data?: { error?: string } } };
      setError(error.response?.data?.error || "Failed to fetch IAM data");
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateRole = async () => {
    if (!newRoleName.trim()) {
      setError("Role name is required");
      return;
    }

    try {
      await apiClient.createRole({
        name: newRoleName,
        permissions: newRolePermissions,
      });
      setNewRoleName("");
      setNewRolePermissions([]);
      fetchData();
    } catch (err: unknown) {
      const error = err as { response?: { data?: { error?: string } } };
      setError(error.response?.data?.error || "Failed to create role");
    }
  };

  const handleCreateUser = async () => {
    if (
      !newUserName.trim() ||
      !newUserEmail.trim() ||
      !newUserPassword.trim() ||
      !selectedRoleId
    ) {
      setError("All fields are required");
      return;
    }

    try {
      // In a real implementation, you'd need the organization ID
      const orgId = "placeholder-org-id";
      await apiClient.createUser(orgId, {
        name: newUserName,
        email: newUserEmail,
        password: newUserPassword,
        role_id: selectedRoleId,
      });
      setNewUserName("");
      setNewUserEmail("");
      setNewUserPassword("");
      setSelectedRoleId("");
      fetchData();
    } catch (err: unknown) {
      const error = err as { response?: { data?: { error?: string } } };
      setError(error.response?.data?.error || "Failed to create user");
    }
  };

  const permissions = [
    { id: 0, name: "Read", description: "View resources", icon: "👁️" },
    {
      id: 1,
      name: "Write",
      description: "Create and modify resources",
      icon: "✏️",
    },
    { id: 2, name: "Delete", description: "Delete resources", icon: "🗑️" },
  ];

  const getRoleIcon = (roleName: string) => {
    if (roleName.toLowerCase().includes("admin"))
      return <Crown className="h-4 w-4 text-yellow-500" />;
    if (roleName.toLowerCase().includes("user"))
      return <UserCheck className="h-4 w-4 text-blue-500" />;
    return <Shield className="h-4 w-4 text-primary" />;
  };

  if (isLoading) {
    return (
      <Layout>
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="space-y-8">
        {/* Header */}
        <div className="space-y-2">
          <div className="flex items-center space-x-3">
            <div className="p-2 rounded-lg bg-gradient-to-br from-primary/10 to-primary/20">
              <Shield className="h-6 w-6 text-primary" />
            </div>
            <div>
              <h1 className="text-3xl font-bold text-foreground">
                Identity & Access Management
              </h1>
              <p className="text-muted-foreground">
                Manage roles, users, and permissions for your organization
              </p>
            </div>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">
                Total Roles
              </CardTitle>
              <div className="p-2 rounded-lg bg-blue-500/10">
                <Shield className="h-4 w-4 text-blue-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">
                {roles.length}
              </div>
              <p className="text-xs text-muted-foreground">
                Active roles in system
              </p>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">
                Total Users
              </CardTitle>
              <div className="p-2 rounded-lg bg-green-500/10">
                <Users className="h-4 w-4 text-green-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">
                {organization?.users.length || 0}
              </div>
              <p className="text-xs text-muted-foreground">Active users</p>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">
                Organization
              </CardTitle>
              <div className="p-2 rounded-lg bg-purple-500/10">
                <Building className="h-4 w-4 text-purple-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-lg font-bold text-foreground truncate">
                {organization?.organization_name || "N/A"}
              </div>
              <p className="text-xs text-muted-foreground">Current org</p>
            </CardContent>
          </Card>
        </div>

        {error && (
          <div className="bg-destructive/10 border border-destructive/50 rounded-lg p-4 flex items-center space-x-2">
            <AlertCircle className="h-4 w-4 text-destructive" />
            <p className="text-destructive">{error}</p>
          </div>
        )}

        <Tabs defaultValue="roles" className="space-y-6">
          <TabsList className="bg-card border border-border">
            <TabsTrigger
              value="roles"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
            >
              <Shield className="h-4 w-4 mr-2" />
              Roles
            </TabsTrigger>
            <TabsTrigger
              value="users"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
            >
              <Users className="h-4 w-4 mr-2" />
              Users
            </TabsTrigger>
            <TabsTrigger
              value="organization"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
            >
              <Building className="h-4 w-4 mr-2" />
              Organization
            </TabsTrigger>
          </TabsList>

          <TabsContent value="roles" className="space-y-6">
            <Card className="border-border/50">
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle className="text-foreground">
                      Roles & Permissions
                    </CardTitle>
                    <CardDescription>
                      Manage roles and their associated permissions
                    </CardDescription>
                  </div>
                  <Dialog>
                    <DialogTrigger asChild>
                      <Button className="bg-primary hover:bg-primary/90">
                        <Plus className="h-4 w-4 mr-2" />
                        Create Role
                      </Button>
                    </DialogTrigger>
                    <DialogContent className="border-border/50">
                      <DialogHeader>
                        <DialogTitle className="text-foreground">
                          Create New Role
                        </DialogTitle>
                        <DialogDescription>
                          Create a new role with specific permissions
                        </DialogDescription>
                      </DialogHeader>
                      <div className="space-y-4">
                        <div>
                          <Label
                            htmlFor="role-name"
                            className="text-foreground"
                          >
                            Role Name
                          </Label>
                          <Input
                            id="role-name"
                            value={newRoleName}
                            onChange={(e) => setNewRoleName(e.target.value)}
                            placeholder="Enter role name"
                            className="border-border/50 focus:border-primary focus:ring-primary/20"
                          />
                        </div>
                        <div>
                          <Label className="text-foreground">Permissions</Label>
                          <div className="space-y-3 mt-3">
                            {permissions.map((permission) => (
                              <label
                                key={permission.id}
                                className="flex items-center space-x-3 p-3 rounded-lg border border-border/50 hover:bg-accent/30 transition-colors cursor-pointer"
                              >
                                <input
                                  type="checkbox"
                                  checked={newRolePermissions.includes(
                                    permission.id
                                  )}
                                  onChange={(e) => {
                                    if (e.target.checked) {
                                      setNewRolePermissions([
                                        ...newRolePermissions,
                                        permission.id,
                                      ]);
                                    } else {
                                      setNewRolePermissions(
                                        newRolePermissions.filter(
                                          (id) => id !== permission.id
                                        )
                                      );
                                    }
                                  }}
                                  className="rounded border-border/50 focus:ring-primary/20"
                                />
                                <span className="text-lg">
                                  {permission.icon}
                                </span>
                                <div>
                                  <span className="text-sm font-medium text-foreground">
                                    {permission.name}
                                  </span>
                                  <p className="text-xs text-muted-foreground">
                                    {permission.description}
                                  </p>
                                </div>
                              </label>
                            ))}
                          </div>
                        </div>
                      </div>
                      <DialogFooter>
                        <Button
                          onClick={handleCreateRole}
                          className="bg-primary hover:bg-primary/90"
                        >
                          Create Role
                        </Button>
                      </DialogFooter>
                    </DialogContent>
                  </Dialog>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {roles.map((role) => (
                    <div
                      key={role.id}
                      className="border border-border/50 rounded-lg p-4 hover:shadow-md transition-all duration-200"
                    >
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-3">
                          <div className="p-2 rounded-lg bg-accent/50">
                            {getRoleIcon(role.name)}
                          </div>
                          <div>
                            <h3 className="font-medium text-foreground">
                              {role.name}
                            </h3>
                            <p className="text-sm text-muted-foreground">
                              Created{" "}
                              {new Date(role.created_at).toLocaleDateString()}
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center space-x-2">
                          <Button
                            variant="outline"
                            size="sm"
                            className="border-border/50"
                          >
                            <Edit className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="outline"
                            size="sm"
                            className="border-border/50"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </div>
                      <div className="mt-3 flex flex-wrap gap-2">
                        {role.permissions.map((permission) => (
                          <Badge
                            key={`perm:${permission.id}`}
                            variant="secondary"
                            className="bg-accent/50 text-accent-foreground"
                          >
                            {permission.name}
                          </Badge>
                        ))}
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="users" className="space-y-6">
            <Card className="border-border/50">
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle className="text-foreground">Users</CardTitle>
                    <CardDescription>
                      Manage users in your organization
                    </CardDescription>
                  </div>
                  <Dialog>
                    <DialogTrigger asChild>
                      <Button className="bg-primary hover:bg-primary/90">
                        <Plus className="h-4 w-4 mr-2" />
                        Create User
                      </Button>
                    </DialogTrigger>
                    <DialogContent className="border-border/50">
                      <DialogHeader>
                        <DialogTitle className="text-foreground">
                          Create New User
                        </DialogTitle>
                        <DialogDescription>
                          Create a new user account
                        </DialogDescription>
                      </DialogHeader>
                      <div className="space-y-4">
                        <div>
                          <Label
                            htmlFor="user-name"
                            className="text-foreground"
                          >
                            Name
                          </Label>
                          <Input
                            id="user-name"
                            value={newUserName}
                            onChange={(e) => setNewUserName(e.target.value)}
                            placeholder="Enter user name"
                            className="border-border/50 focus:border-primary focus:ring-primary/20"
                          />
                        </div>
                        <div>
                          <Label
                            htmlFor="user-email"
                            className="text-foreground"
                          >
                            Email
                          </Label>
                          <Input
                            id="user-email"
                            type="email"
                            value={newUserEmail}
                            onChange={(e) => setNewUserEmail(e.target.value)}
                            placeholder="Enter user email"
                            className="border-border/50 focus:border-primary focus:ring-primary/20"
                          />
                        </div>
                        <div>
                          <Label
                            htmlFor="user-password"
                            className="text-foreground"
                          >
                            Password
                          </Label>
                          <Input
                            id="user-password"
                            type="password"
                            value={newUserPassword}
                            onChange={(e) => setNewUserPassword(e.target.value)}
                            placeholder="Enter user password"
                            className="border-border/50 focus:border-primary focus:ring-primary/20"
                          />
                        </div>
                        <div>
                          <Label
                            htmlFor="user-role"
                            className="text-foreground"
                          >
                            Role
                          </Label>
                          <select
                            id="user-role"
                            value={selectedRoleId}
                            onChange={(e) => setSelectedRoleId(e.target.value)}
                            className="w-full p-2 border border-border/50 rounded-md bg-background text-foreground focus:border-primary focus:ring-primary/20"
                          >
                            <option value="">Select a role</option>
                            {roles.map((role) => (
                              <option key={role.id} value={role.id}>
                                {role.name}
                              </option>
                            ))}
                          </select>
                        </div>
                      </div>
                      <DialogFooter>
                        <Button
                          onClick={handleCreateUser}
                          className="bg-primary hover:bg-primary/90"
                        >
                          Create User
                        </Button>
                      </DialogFooter>
                    </DialogContent>
                  </Dialog>
                </div>
              </CardHeader>
              <CardContent>
                {organization && (
                  <div className="space-y-4">
                    {organization.users.map((user) => (
                      <div
                        key={user.id}
                        className="border border-border/50 rounded-lg p-4 hover:shadow-md transition-all duration-200"
                      >
                        <div className="flex items-center justify-between">
                          <div className="flex items-center space-x-3">
                            <div className="p-2 rounded-lg bg-accent/50">
                              <UserCheck className="h-4 w-4 text-accent-foreground" />
                            </div>
                            <div>
                              <h3 className="font-medium text-foreground">
                                {user.name}
                              </h3>
                              <p className="text-sm text-muted-foreground">
                                {user.email}
                              </p>
                            </div>
                          </div>
                          <div className="flex items-center space-x-3">
                            <Badge
                              variant="outline"
                              className="bg-secondary/20"
                            >
                              {user.role_id}
                            </Badge>
                            <div className="flex space-x-2">
                              <Button
                                variant="outline"
                                size="sm"
                                className="border-border/50"
                              >
                                <Edit className="h-4 w-4" />
                              </Button>
                              <Button
                                variant="outline"
                                size="sm"
                                className="border-border/50"
                              >
                                <Trash2 className="h-4 w-4" />
                              </Button>
                            </div>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="organization" className="space-y-6">
            <Card className="border-border/50">
              <CardHeader>
                <CardTitle className="text-foreground">
                  Organization Details
                </CardTitle>
                <CardDescription>
                  Information about your organization
                </CardDescription>
              </CardHeader>
              <CardContent>
                {organization && (
                  <div className="space-y-6">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div className="p-4 rounded-lg bg-accent/20 border border-accent/30">
                        <div className="flex items-center space-x-2 mb-2">
                          <Building className="h-5 w-5 text-primary" />
                          <h3 className="font-medium text-foreground">
                            Organization Name
                          </h3>
                        </div>
                        <p className="text-muted-foreground">
                          {organization.organization_name}
                        </p>
                      </div>
                      <div className="p-4 rounded-lg bg-accent/20 border border-accent/30">
                        <div className="flex items-center space-x-2 mb-2">
                          <Users className="h-5 w-5 text-primary" />
                          <h3 className="font-medium text-foreground">
                            Total Users
                          </h3>
                        </div>
                        <p className="text-muted-foreground">
                          {organization.users.length} users
                        </p>
                      </div>
                    </div>
                    <div className="p-4 rounded-lg bg-primary/10 border border-primary/20">
                      <div className="flex items-center space-x-2 mb-2">
                        <Settings className="h-5 w-5 text-primary" />
                        <h3 className="font-medium text-foreground">
                          Organization Settings
                        </h3>
                      </div>
                      <p className="text-sm text-muted-foreground">
                        Manage your organization&apos;s settings, billing, and
                        preferences from the organization dashboard.
                      </p>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </Layout>
  );
}
