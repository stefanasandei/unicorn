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
import { Shield, Users, Building, Plus, Edit, Trash2 } from "lucide-react";
import { apiClient } from "@/lib/api";
import { Role, Organization } from "@/types/api";

export default function IAMPage() {
  const [roles, setRoles] = useState<Role[]>([]);
  const [organization, setOrganization] = useState<{
    organization_name: string;
    users: any[];
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
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to fetch IAM data");
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
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to create role");
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
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to create user");
    }
  };

  const permissions = [
    { id: 0, name: "Read", description: "View resources" },
    { id: 1, name: "Write", description: "Create and modify resources" },
    { id: 2, name: "Delete", description: "Delete resources" },
  ];

  if (isLoading) {
    return (
      <Layout>
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">
            Identity & Access Management
          </h1>
          <p className="text-gray-600">
            Manage roles, users, and permissions for your organization
          </p>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 rounded-md p-4">
            <p className="text-red-800">{error}</p>
          </div>
        )}

        <Tabs defaultValue="roles" className="space-y-4">
          <TabsList>
            <TabsTrigger value="roles" className="flex items-center gap-2">
              <Shield className="h-4 w-4" />
              Roles
            </TabsTrigger>
            <TabsTrigger value="users" className="flex items-center gap-2">
              <Users className="h-4 w-4" />
              Users
            </TabsTrigger>
            <TabsTrigger
              value="organization"
              className="flex items-center gap-2"
            >
              <Building className="h-4 w-4" />
              Organization
            </TabsTrigger>
          </TabsList>

          <TabsContent value="roles" className="space-y-4">
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle>Roles</CardTitle>
                    <CardDescription>
                      Manage roles and their permissions
                    </CardDescription>
                  </div>
                  <Dialog>
                    <DialogTrigger asChild>
                      <Button>
                        <Plus className="h-4 w-4 mr-2" />
                        Create Role
                      </Button>
                    </DialogTrigger>
                    <DialogContent>
                      <DialogHeader>
                        <DialogTitle>Create New Role</DialogTitle>
                        <DialogDescription>
                          Create a new role with specific permissions
                        </DialogDescription>
                      </DialogHeader>
                      <div className="space-y-4">
                        <div>
                          <Label htmlFor="role-name">Role Name</Label>
                          <Input
                            id="role-name"
                            value={newRoleName}
                            onChange={(e) => setNewRoleName(e.target.value)}
                            placeholder="Enter role name"
                          />
                        </div>
                        <div>
                          <Label>Permissions</Label>
                          <div className="space-y-2 mt-2">
                            {permissions.map((permission) => (
                              <label
                                key={permission.id}
                                className="flex items-center space-x-2"
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
                                />
                                <span className="text-sm">
                                  <strong>{permission.name}</strong> -{" "}
                                  {permission.description}
                                </span>
                              </label>
                            ))}
                          </div>
                        </div>
                      </div>
                      <DialogFooter>
                        <Button onClick={handleCreateRole}>Create Role</Button>
                      </DialogFooter>
                    </DialogContent>
                  </Dialog>
                </div>
              </CardHeader>
              <CardContent>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Role Name</TableHead>
                      <TableHead>Permissions</TableHead>
                      <TableHead>Created</TableHead>
                      <TableHead>Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {roles.map((role) => (
                      <TableRow key={role.id}>
                        <TableCell className="font-medium">
                          {role.name}
                        </TableCell>
                        <TableCell>
                          <div className="flex flex-wrap gap-1">
                            {role.permissions.map((permission) => (
                              <Badge
                                key={`perm:${permission.id}`}
                                variant="secondary"
                              >
                                {permission.name}
                              </Badge>
                            ))}
                          </div>
                        </TableCell>
                        <TableCell>
                          {new Date(role.created_at).toLocaleDateString()}
                        </TableCell>
                        <TableCell>
                          <div className="flex space-x-2">
                            <Button variant="outline" size="sm">
                              <Edit className="h-4 w-4" />
                            </Button>
                            <Button variant="outline" size="sm">
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="users" className="space-y-4">
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle>Users</CardTitle>
                    <CardDescription>
                      Manage users in your organization
                    </CardDescription>
                  </div>
                  <Dialog>
                    <DialogTrigger asChild>
                      <Button>
                        <Plus className="h-4 w-4 mr-2" />
                        Create User
                      </Button>
                    </DialogTrigger>
                    <DialogContent>
                      <DialogHeader>
                        <DialogTitle>Create New User</DialogTitle>
                        <DialogDescription>
                          Create a new user account
                        </DialogDescription>
                      </DialogHeader>
                      <div className="space-y-4">
                        <div>
                          <Label htmlFor="user-name">Name</Label>
                          <Input
                            id="user-name"
                            value={newUserName}
                            onChange={(e) => setNewUserName(e.target.value)}
                            placeholder="Enter user name"
                          />
                        </div>
                        <div>
                          <Label htmlFor="user-email">Email</Label>
                          <Input
                            id="user-email"
                            type="email"
                            value={newUserEmail}
                            onChange={(e) => setNewUserEmail(e.target.value)}
                            placeholder="Enter user email"
                          />
                        </div>
                        <div>
                          <Label htmlFor="user-password">Password</Label>
                          <Input
                            id="user-password"
                            type="password"
                            value={newUserPassword}
                            onChange={(e) => setNewUserPassword(e.target.value)}
                            placeholder="Enter user password"
                          />
                        </div>
                        <div>
                          <Label htmlFor="user-role">Role</Label>
                          <select
                            id="user-role"
                            value={selectedRoleId}
                            onChange={(e) => setSelectedRoleId(e.target.value)}
                            className="w-full p-2 border border-gray-300 rounded-md"
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
                        <Button onClick={handleCreateUser}>Create User</Button>
                      </DialogFooter>
                    </DialogContent>
                  </Dialog>
                </div>
              </CardHeader>
              <CardContent>
                {organization && (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Name</TableHead>
                        <TableHead>Role</TableHead>
                        <TableHead>Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {organization.users.map((user: any) => (
                        <TableRow key={user.id}>
                          <TableCell className="font-medium">
                            {user.name}
                          </TableCell>
                          <TableCell>
                            <Badge variant="outline">{user.role_id}</Badge>
                          </TableCell>
                          <TableCell>
                            <div className="flex space-x-2">
                              <Button variant="outline" size="sm">
                                <Edit className="h-4 w-4" />
                              </Button>
                              <Button variant="outline" size="sm">
                                <Trash2 className="h-4 w-4" />
                              </Button>
                            </div>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="organization" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>Organization Details</CardTitle>
                <CardDescription>
                  Information about your organization
                </CardDescription>
              </CardHeader>
              <CardContent>
                {organization && (
                  <div className="space-y-4">
                    <div>
                      <h3 className="font-medium text-gray-900">
                        Organization Name
                      </h3>
                      <p className="text-gray-600">
                        {organization.organization_name}
                      </p>
                    </div>
                    <div>
                      <h3 className="font-medium text-gray-900">Total Users</h3>
                      <p className="text-gray-600">
                        {organization.users.length}
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
