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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
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
  Database,
  Plus,
  RefreshCw,
  Trash2,
  Copy,
  ExternalLink,
  HardDrive,
  Settings,
  AlertCircle,
  Clock,
  Zap,
  Server,
  Key,
} from "lucide-react";
import { apiClient } from "@/lib/api";
import {
  RDBInstanceInfo,
  RDBCreateRequest,
  RDBType,
  RDBPreset,
} from "@/types/api";
import { toast } from "sonner";

export default function RDBPage() {
  const [instances, setInstances] = useState<RDBInstanceInfo[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isCreating, setIsCreating] = useState(false);
  const [error, setError] = useState("");

  // Form states
  const [newInstanceName, setNewInstanceName] = useState("");
  const [newInstanceType, setNewInstanceType] = useState<RDBType>("postgresql");
  const [newInstancePreset, setNewInstancePreset] =
    useState<RDBPreset>("micro");
  const [newInstanceDatabase, setNewInstanceDatabase] = useState("main");
  const [newInstanceUsername, setNewInstanceUsername] = useState("user");
  const [newInstancePassword, setNewInstancePassword] = useState("");
  const [newInstancePort, setNewInstancePort] = useState("");
  const [newInstanceEnvironment, setNewInstanceEnvironment] = useState("");

  // Volume states
  const [volumes, setVolumes] = useState<
    Array<{
      name: string;
      size: number;
    }>
  >([]);

  // Connection URL dialog
  const [showConnectionDialog, setShowConnectionDialog] = useState(false);
  const [selectedInstance, setSelectedInstance] =
    useState<RDBInstanceInfo | null>(null);

  useEffect(() => {
    fetchInstances();
  }, []);

  const fetchInstances = async () => {
    try {
      setIsLoading(true);
      setError("");
      const instancesData = await apiClient.listRDB();
      setInstances(instancesData);
    } catch (err: unknown) {
      console.error("RDB fetch error:", err);
      const error = err as {
        response?: { data?: { error?: string; details?: string } };
        message?: string;
      };
      const errorMessage =
        error.response?.data?.error ||
        error.response?.data?.details ||
        error.message ||
        "Failed to fetch RDB instances";
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateInstance = async () => {
    if (!newInstanceDatabase.trim()) {
      setError("Database name is required");
      return;
    }

    if (!newInstanceUsername.trim()) {
      setError("Username is required");
      return;
    }

    // Validate volume sizes (minimum 1MB, maximum 100GB)
    for (const volume of volumes) {
      if (volume.size < 1 || volume.size > 100000) {
        setError("Volume size must be between 1MB and 100GB");
        return;
      }
    }

    try {
      setIsCreating(true);
      setError("");

      let environment = {};
      if (newInstanceEnvironment.trim()) {
        try {
          environment = JSON.parse(newInstanceEnvironment);
        } catch {
          setError("Invalid environment variables JSON format");
          return;
        }
      }

      const request: RDBCreateRequest = {
        name: newInstanceName.trim() || undefined,
        type: newInstanceType,
        preset: newInstancePreset,
        database: newInstanceDatabase,
        username: newInstanceUsername,
        password: newInstancePassword || undefined,
        port: newInstancePort || undefined,
        environment:
          Object.keys(environment).length > 0 ? environment : undefined,
        volumes: volumes.length > 0 ? volumes : undefined,
      };

      const newInstance = await apiClient.createRDB(request);
      setInstances([...instances, newInstance]);

      // Reset form
      setNewInstanceName("");
      setNewInstanceType("postgresql");
      setNewInstancePreset("micro");
      setNewInstanceDatabase("main");
      setNewInstanceUsername("user");
      setNewInstancePassword("");
      setNewInstancePort("");
      setNewInstanceEnvironment("");
      setVolumes([]);

      toast.success("Database instance created successfully!");

      window.location.reload();
    } catch (err: unknown) {
      console.error("RDB creation error:", err);
      const error = err as {
        response?: { data?: { error?: string; details?: string } };
        message?: string;
      };
      const errorMessage =
        error.response?.data?.error ||
        error.response?.data?.details ||
        error.message ||
        "Failed to create RDB instance";
      setError(errorMessage);
    } finally {
      setIsCreating(false);
    }
  };

  const handleDeleteInstance = async (instanceId: string) => {
    if (
      !confirm(
        "Are you sure you want to delete this database instance? This action cannot be undone."
      )
    ) {
      return;
    }

    try {
      await apiClient.deleteRDB(instanceId);
      setInstances(instances.filter((instance) => instance.id !== instanceId));
      toast.success("Database instance deleted successfully!");
    } catch (err: unknown) {
      console.error("RDB deletion error:", err);
      const error = err as {
        response?: { data?: { error?: string; details?: string } };
        message?: string;
      };
      const errorMessage =
        error.response?.data?.error ||
        error.response?.data?.details ||
        error.message ||
        "Failed to delete RDB instance";
      setError(errorMessage);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case "running":
        return "bg-green-500 text-white";
      case "stopped":
        return "bg-red-500 text-white";
      case "starting":
        return "bg-yellow-500 text-white";
      default:
        return "bg-muted text-muted-foreground";
    }
  };

  const getTypeIcon = (type: RDBType) => {
    return type === "postgresql" ? "🐘" : "🐬";
  };

  const getTypeColor = (type: RDBType) => {
    return type === "postgresql" ? "text-blue-500" : "text-orange-500";
  };

  const generateConnectionUrl = (instance: RDBInstanceInfo) => {
    // Note: Password is not included in the API response for security
    // Users should use the password they set during creation
    const password = "[YOUR_PASSWORD]"; // Placeholder for security
    if (instance.type === "postgresql") {
      return `postgresql://${instance.username}:${password}@${instance.host}:${instance.port}/${instance.database}`;
    } else {
      return `mysql://${instance.username}:${password}@${instance.host}:${instance.port}/${instance.database}`;
    }
  };

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      toast.success("Connection URL copied to clipboard!");
    } catch (err) {
      toast.error("Failed to copy to clipboard");
    }
  };

  const addVolume = () => {
    setVolumes([...volumes, { name: "", size: 1024 }]); // Default to 1GB in MB
  };

  const removeVolume = (index: number) => {
    setVolumes(volumes.filter((_, i) => i !== index));
  };

  const updateVolume = (
    index: number,
    field: string,
    value: string | number
  ) => {
    const newVolumes = [...volumes];
    newVolumes[index] = { ...newVolumes[index], [field]: value };
    setVolumes(newVolumes);
  };

  const getDefaultPort = (type: RDBType) => {
    return type === "postgresql" ? "5432" : "3306";
  };

  const runningInstances =
    instances?.filter((i) => i.status?.toLowerCase() === "running")?.length ||
    0;
  const totalInstances = instances?.length || 0;

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
              <Database className="h-6 w-6 text-primary" />
            </div>
            <div>
              <h1 className="text-3xl font-bold text-foreground">Databases</h1>
              <p className="text-muted-foreground">
                Manage your PostgreSQL and MySQL database instances
              </p>
            </div>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">
                Total Instances
              </CardTitle>
              <div className="p-2 rounded-lg bg-blue-500/10">
                <Database className="h-4 w-4 text-blue-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">
                {totalInstances}
              </div>
              <p className="text-xs text-muted-foreground">
                Database instances
              </p>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">
                Running
              </CardTitle>
              <div className="p-2 rounded-lg bg-green-500/10">
                <Server className="h-4 w-4 text-green-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">
                {runningInstances}
              </div>
              <p className="text-xs text-muted-foreground">Active instances</p>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">
                Storage
              </CardTitle>
              <div className="p-2 rounded-lg bg-purple-500/10">
                <HardDrive className="h-4 w-4 text-purple-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">
                Persistent
              </div>
              <p className="text-xs text-muted-foreground">Volumes attached</p>
            </CardContent>
          </Card>
        </div>

        {error && (
          <div className="bg-destructive/10 border border-destructive/50 rounded-lg p-4 flex items-center space-x-2">
            <AlertCircle className="h-4 w-4 text-destructive" />
            <p className="text-destructive">{error}</p>
          </div>
        )}

        <Card className="border-border/50">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-foreground">
                  Database Instances
                </CardTitle>
                <CardDescription>
                  Create and manage database instances with persistent storage
                </CardDescription>
              </div>
              <div className="flex space-x-2">
                <Button
                  variant="outline"
                  onClick={fetchInstances}
                  disabled={isLoading}
                  className="border-border/50"
                >
                  <RefreshCw
                    className={`h-4 w-4 mr-2 ${
                      isLoading ? "animate-spin" : ""
                    }`}
                  />
                  Refresh
                </Button>
                <Dialog>
                  <DialogTrigger asChild>
                    <Button className="bg-primary hover:bg-primary/90">
                      <Plus className="h-4 w-4 mr-2" />
                      Create Database
                    </Button>
                  </DialogTrigger>
                  <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto border-border/50">
                    <DialogHeader>
                      <DialogTitle className="text-foreground">
                        Create New Database Instance
                      </DialogTitle>
                      <DialogDescription>
                        Deploy a new PostgreSQL or MySQL database instance
                      </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4">
                      <div>
                        <Label
                          htmlFor="instance-name"
                          className="text-foreground"
                        >
                          Instance Name (optional)
                        </Label>
                        <Input
                          id="instance-name"
                          value={newInstanceName}
                          onChange={(e) => setNewInstanceName(e.target.value)}
                          placeholder="Auto-generated if not provided"
                          className="border-border/50 focus:border-primary focus:ring-primary/20"
                        />
                      </div>

                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <Label
                            htmlFor="instance-type"
                            className="text-foreground"
                          >
                            Database Type
                          </Label>
                          <Select
                            value={newInstanceType}
                            onValueChange={(value: RDBType) => {
                              setNewInstanceType(value);
                              setNewInstancePort(getDefaultPort(value));
                            }}
                          >
                            <SelectTrigger className="border-border/50">
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="postgresql">
                                <span className="flex items-center">
                                  🐘 PostgreSQL
                                </span>
                              </SelectItem>
                              <SelectItem value="mysql">
                                <span className="flex items-center">
                                  🐬 MySQL
                                </span>
                              </SelectItem>
                            </SelectContent>
                          </Select>
                        </div>

                        <div>
                          <Label
                            htmlFor="instance-preset"
                            className="text-foreground"
                          >
                            Resource Preset
                          </Label>
                          <Select
                            value={newInstancePreset}
                            onValueChange={(value: RDBPreset) =>
                              setNewInstancePreset(value)
                            }
                          >
                            <SelectTrigger className="border-border/50">
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="micro">
                                Micro (0.5 CPU, 512MB RAM)
                              </SelectItem>
                              <SelectItem value="small">
                                Small (1 CPU, 1GB RAM)
                              </SelectItem>
                              <SelectItem value="medium">
                                Medium (2 CPU, 2GB RAM)
                              </SelectItem>
                            </SelectContent>
                          </Select>
                        </div>
                      </div>

                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <Label
                            htmlFor="instance-database"
                            className="text-foreground"
                          >
                            Database Name
                          </Label>
                          <Input
                            id="instance-database"
                            value={newInstanceDatabase}
                            onChange={(e) =>
                              setNewInstanceDatabase(e.target.value)
                            }
                            placeholder="main"
                            className="border-border/50 focus:border-primary focus:ring-primary/20"
                          />
                        </div>

                        <div>
                          <Label
                            htmlFor="instance-port"
                            className="text-foreground"
                          >
                            Port (optional)
                          </Label>
                          <Input
                            id="instance-port"
                            value={newInstancePort}
                            onChange={(e) => setNewInstancePort(e.target.value)}
                            placeholder={getDefaultPort(newInstanceType)}
                            className="border-border/50 focus:border-primary focus:ring-primary/20"
                          />
                        </div>
                      </div>

                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <Label
                            htmlFor="instance-username"
                            className="text-foreground"
                          >
                            Username
                          </Label>
                          <Input
                            id="instance-username"
                            value={newInstanceUsername}
                            onChange={(e) =>
                              setNewInstanceUsername(e.target.value)
                            }
                            placeholder="user"
                            className="border-border/50 focus:border-primary focus:ring-primary/20"
                          />
                        </div>

                        <div>
                          <Label
                            htmlFor="instance-password"
                            className="text-foreground"
                          >
                            Password (optional)
                          </Label>
                          <Input
                            id="instance-password"
                            type="password"
                            value={newInstancePassword}
                            onChange={(e) =>
                              setNewInstancePassword(e.target.value)
                            }
                            placeholder="Auto-generated if not provided"
                            className="border-border/50 focus:border-primary focus:ring-primary/20"
                          />
                        </div>
                      </div>

                      <div>
                        <Label
                          htmlFor="instance-environment"
                          className="text-foreground"
                        >
                          Environment Variables (JSON, optional)
                        </Label>
                        <Input
                          id="instance-environment"
                          value={newInstanceEnvironment}
                          onChange={(e) =>
                            setNewInstanceEnvironment(e.target.value)
                          }
                          placeholder='{"POSTGRES_INITDB_ARGS": "--encoding=UTF-8"}'
                          className="border-border/50 focus:border-primary focus:ring-primary/20"
                        />
                      </div>

                      <div>
                        <div className="flex items-center justify-between mb-2">
                          <Label className="text-foreground">
                            Database Storage Volumes (optional)
                          </Label>
                          <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={addVolume}
                            className="border-border/50"
                          >
                            <Plus className="h-4 w-4 mr-2" />
                            Add Volume
                          </Button>
                        </div>
                        <p className="text-sm text-muted-foreground mb-3">
                          Add persistent storage for your database. Volumes are
                          automatically mounted to the appropriate database
                          directory.
                        </p>
                        {volumes.map((volume, index) => (
                          <div
                            key={index}
                            className="grid grid-cols-2 gap-2 mb-2"
                          >
                            <div className="text-sm font-medium text-foreground flex items-center">
                              Volume Name
                            </div>
                            <Input
                              placeholder="e.g., data"
                              value={volume.name}
                              onChange={(e) =>
                                updateVolume(index, "name", e.target.value)
                              }
                              className="border-border/50 focus:border-primary focus:ring-primary/20"
                            />
                            <div className="text-sm font-medium text-foreground flex items-center">
                              Size (MB)
                            </div>
                            <Input
                              type="number"
                              placeholder="1024"
                              min="1"
                              max="100000"
                              value={volume.size}
                              onChange={(e) =>
                                updateVolume(
                                  index,
                                  "size",
                                  parseInt(e.target.value) || 1024
                                )
                              }
                              className="border-border/50 focus:border-primary focus:ring-primary/20"
                            />
                            <div className="text-sm font-medium text-foreground flex items-center">
                              Mount Path
                            </div>
                            <div className="text-sm text-muted-foreground flex items-center">
                              {newInstanceType === "postgresql"
                                ? "/var/lib/postgresql/data"
                                : "/var/lib/mysql"}
                            </div>
                            <Button
                              type="button"
                              variant="outline"
                              size="sm"
                              className="col-span-2 border-border/50"
                              onClick={() => removeVolume(index)}
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </div>
                        ))}
                      </div>
                    </div>
                    <DialogFooter>
                      <Button
                        onClick={handleCreateInstance}
                        disabled={isCreating}
                        className="bg-primary hover:bg-primary/90"
                      >
                        {isCreating ? (
                          <div className="flex items-center space-x-2">
                            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-primary-foreground"></div>
                            <span>Creating...</span>
                          </div>
                        ) : (
                          "Create Database"
                        )}
                      </Button>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="flex items-center justify-center py-8">
                <RefreshCw className="h-6 w-6 animate-spin text-primary" />
                <span className="ml-2 text-muted-foreground">
                  Loading database instances...
                </span>
              </div>
            ) : !instances ? (
              <div className="text-center py-8">
                <Database className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="text-lg font-medium text-foreground mb-2">
                  No database instances
                </h3>
                <p className="text-muted-foreground mb-4">
                  Create your first database instance to get started.
                </p>
              </div>
            ) : (
              <div className="space-y-4">
                {instances && instances.length > 0 ? (
                  instances.map((instance) => (
                    <div
                      key={instance.id}
                      className="border border-border/50 rounded-lg p-4 hover:shadow-md transition-all duration-200"
                    >
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-3">
                          <div className="p-2 rounded-lg bg-accent/50">
                            <Database className="h-4 w-4 text-accent-foreground" />
                          </div>
                          <div>
                            <h3 className="font-medium text-foreground">
                              {instance.name}
                            </h3>
                            <p className="text-sm text-muted-foreground flex items-center gap-1">
                              <Clock className="h-3 w-3" />
                              Created{" "}
                              {new Date(
                                instance.created_at
                              ).toLocaleDateString()}
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center space-x-3">
                          <Badge variant="outline" className="bg-secondary/20">
                            <span className="mr-1">
                              {getTypeIcon(instance.type)}
                            </span>
                            {instance.type}
                          </Badge>
                          <Badge className={getStatusColor(instance.status)}>
                            {instance.status}
                          </Badge>
                          <div className="flex space-x-1">
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => {
                                setSelectedInstance(instance);
                                setShowConnectionDialog(true);
                              }}
                              className="border-border/50"
                            >
                              <ExternalLink className="h-4 w-4 mr-1" />
                              Connection
                            </Button>
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => handleDeleteInstance(instance.id)}
                              className="border-border/50 text-destructive hover:text-destructive hover:bg-destructive/10"
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </div>
                        </div>
                      </div>
                      <div className="mt-3 grid grid-cols-3 gap-4 text-sm">
                        <div className="p-2 rounded-lg bg-accent/20">
                          <div className="text-muted-foreground">Port</div>
                          <div className="font-medium text-foreground">
                            {instance.port}
                          </div>
                        </div>
                        <div className="p-2 rounded-lg bg-accent/20">
                          <div className="text-muted-foreground">Database</div>
                          <div className="font-medium text-foreground">
                            {instance.database}
                          </div>
                        </div>
                        <div className="p-2 rounded-lg bg-accent/20">
                          <div className="text-muted-foreground">Volumes</div>
                          <div className="font-medium text-foreground flex items-center">
                            <HardDrive className="h-3 w-3 mr-1" />
                            {instance.volumes ? instance.volumes.length : 0}
                          </div>
                        </div>
                      </div>
                    </div>
                  ))
                ) : (
                  <div className="text-center py-8">
                    <Database className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
                    <h3 className="text-lg font-medium text-foreground mb-2">
                      No database instances found
                    </h3>
                    <p className="text-muted-foreground">
                      Create your first database instance to get started.
                    </p>
                  </div>
                )}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Connection URL Dialog */}
        <Dialog
          open={showConnectionDialog}
          onOpenChange={setShowConnectionDialog}
        >
          <DialogContent className="border-border/50">
            <DialogHeader>
              <DialogTitle className="text-foreground">
                Connection Information
              </DialogTitle>
              <DialogDescription>
                Use these details to connect to your database instance
              </DialogDescription>
            </DialogHeader>
            {selectedInstance && (
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label className="text-foreground">Host</Label>
                    <Input
                      value={selectedInstance.host}
                      readOnly
                      className="bg-muted/50 border-border/50"
                    />
                  </div>
                  <div>
                    <Label className="text-foreground">Port</Label>
                    <Input
                      value={selectedInstance.port}
                      readOnly
                      className="bg-muted/50 border-border/50"
                    />
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label className="text-foreground">Database</Label>
                    <Input
                      value={selectedInstance.database}
                      readOnly
                      className="bg-muted/50 border-border/50"
                    />
                  </div>
                  <div>
                    <Label className="text-foreground">Username</Label>
                    <Input
                      value={selectedInstance.username}
                      readOnly
                      className="bg-muted/50 border-border/50"
                    />
                  </div>
                </div>
                <div>
                  <Label className="text-foreground">Connection URL</Label>
                  <div className="flex space-x-2">
                    <Input
                      value={generateConnectionUrl(selectedInstance)}
                      readOnly
                      className="font-mono text-sm bg-muted/50 border-border/50"
                    />
                    <Button
                      variant="outline"
                      onClick={() =>
                        copyToClipboard(generateConnectionUrl(selectedInstance))
                      }
                      className="border-border/50"
                    >
                      <Copy className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
                {selectedInstance.volumes && (
                  <div>
                    <Label className="text-foreground">Volumes</Label>
                    <div className="space-y-2">
                      {selectedInstance.volumes.map((volume, index) => (
                        <div
                          key={index}
                          className="flex justify-between text-sm p-2 rounded-lg bg-accent/20"
                        >
                          <span className="text-foreground">{volume.name}</span>
                          <span className="text-muted-foreground">
                            {volume.size}MB
                          </span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            )}
          </DialogContent>
        </Dialog>
      </div>
    </Layout>
  );
}
