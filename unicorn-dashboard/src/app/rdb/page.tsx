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
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Database,
  Plus,
  RefreshCw,
  Trash2,
  Copy,
  ExternalLink,
  HardDrive,
  Settings,
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
        return "bg-green-100 text-green-800";
      case "stopped":
        return "bg-red-100 text-red-800";
      case "starting":
        return "bg-yellow-100 text-yellow-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  const getTypeIcon = (type: RDBType) => {
    return type === "postgresql" ? "🐘" : "🐬";
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

  return (
    <Layout>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Databases</h1>
          <p className="text-gray-600">
            Manage your PostgreSQL and MySQL database instances
          </p>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 rounded-md p-4">
            <p className="text-red-800">{error}</p>
          </div>
        )}

        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Database Instances</CardTitle>
                <CardDescription>
                  Create and manage database instances with persistent storage
                </CardDescription>
              </div>
              <div className="flex space-x-2">
                <Button
                  variant="outline"
                  onClick={fetchInstances}
                  disabled={isLoading}
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
                    <Button>
                      <Plus className="h-4 w-4 mr-2" />
                      Create Database
                    </Button>
                  </DialogTrigger>
                  <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
                    <DialogHeader>
                      <DialogTitle>Create New Database Instance</DialogTitle>
                      <DialogDescription>
                        Deploy a new PostgreSQL or MySQL database instance
                      </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4">
                      <div>
                        <Label htmlFor="instance-name">
                          Instance Name (optional)
                        </Label>
                        <Input
                          id="instance-name"
                          value={newInstanceName}
                          onChange={(e) => setNewInstanceName(e.target.value)}
                          placeholder="Auto-generated if not provided"
                        />
                      </div>

                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <Label htmlFor="instance-type">Database Type</Label>
                          <Select
                            value={newInstanceType}
                            onValueChange={(value: RDBType) => {
                              setNewInstanceType(value);
                              setNewInstancePort(getDefaultPort(value));
                            }}
                          >
                            <SelectTrigger>
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
                          <Label htmlFor="instance-preset">
                            Resource Preset
                          </Label>
                          <Select
                            value={newInstancePreset}
                            onValueChange={(value: RDBPreset) =>
                              setNewInstancePreset(value)
                            }
                          >
                            <SelectTrigger>
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
                          <Label htmlFor="instance-database">
                            Database Name
                          </Label>
                          <Input
                            id="instance-database"
                            value={newInstanceDatabase}
                            onChange={(e) =>
                              setNewInstanceDatabase(e.target.value)
                            }
                            placeholder="main"
                          />
                        </div>

                        <div>
                          <Label htmlFor="instance-port">Port (optional)</Label>
                          <Input
                            id="instance-port"
                            value={newInstancePort}
                            onChange={(e) => setNewInstancePort(e.target.value)}
                            placeholder={getDefaultPort(newInstanceType)}
                          />
                        </div>
                      </div>

                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <Label htmlFor="instance-username">Username</Label>
                          <Input
                            id="instance-username"
                            value={newInstanceUsername}
                            onChange={(e) =>
                              setNewInstanceUsername(e.target.value)
                            }
                            placeholder="user"
                          />
                        </div>

                        <div>
                          <Label htmlFor="instance-password">
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
                          />
                        </div>
                      </div>

                      <div>
                        <Label htmlFor="instance-environment">
                          Environment Variables (JSON, optional)
                        </Label>
                        <Input
                          id="instance-environment"
                          value={newInstanceEnvironment}
                          onChange={(e) =>
                            setNewInstanceEnvironment(e.target.value)
                          }
                          placeholder='{"POSTGRES_INITDB_ARGS": "--encoding=UTF-8"}'
                        />
                      </div>

                      <div>
                        <div className="flex items-center justify-between mb-2">
                          <Label>Database Storage Volumes (optional)</Label>
                          <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={addVolume}
                          >
                            <Plus className="h-4 w-4 mr-2" />
                            Add Volume
                          </Button>
                        </div>
                        <p className="text-sm text-gray-600 mb-3">
                          Add persistent storage for your database. Volumes are
                          automatically mounted to the appropriate database
                          directory.
                        </p>
                        {volumes.map((volume, index) => (
                          <div
                            key={index}
                            className="grid grid-cols-2 gap-2 mb-2"
                          >
                            <div className="text-sm font-medium text-gray-700 flex items-center">
                              Volume Name
                            </div>
                            <Input
                              placeholder="e.g., data"
                              value={volume.name}
                              onChange={(e) =>
                                updateVolume(index, "name", e.target.value)
                              }
                            />
                            <div className="text-sm font-medium text-gray-700 flex items-center">
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
                            />
                            <div className="text-sm font-medium text-gray-700 flex items-center">
                              Mount Path
                            </div>
                            <div className="text-sm text-gray-500 flex items-center">
                              {newInstanceType === "postgresql"
                                ? "/var/lib/postgresql/data"
                                : "/var/lib/mysql"}
                            </div>
                            <Button
                              type="button"
                              variant="outline"
                              size="sm"
                              className="col-span-2"
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
                      >
                        {isCreating ? "Creating..." : "Create Database"}
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
                <RefreshCw className="h-6 w-6 animate-spin" />
                <span className="ml-2">Loading database instances...</span>
              </div>
            ) : !instances ? (
              <div className="text-center py-8">
                <Database className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-gray-900 mb-2">
                  No database instances
                </h3>
                <p className="text-gray-600 mb-4">
                  Create your first database instance to get started.
                </p>
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Port</TableHead>
                    <TableHead>Database</TableHead>
                    <TableHead>Volumes</TableHead>
                    <TableHead>Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {instances.map((instance) => (
                    <TableRow key={instance.id}>
                      <TableCell className="font-medium">
                        {instance.name}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center">
                          <span className="mr-2">
                            {getTypeIcon(instance.type)}
                          </span>
                          {instance.type}
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge className={getStatusColor(instance.status)}>
                          {instance.status}
                        </Badge>
                      </TableCell>
                      <TableCell>{instance.port}</TableCell>
                      <TableCell>{instance.database}</TableCell>
                      <TableCell>
                        {instance.volumes ? (
                          <div className="flex items-center">
                            <HardDrive className="h-4 w-4 mr-1" />
                            {instance.volumes.length} volume(s)
                          </div>
                        ) : (
                          <span className="text-gray-500">No volumes</span>
                        )}
                      </TableCell>
                      <TableCell>
                        <div className="flex space-x-2">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => {
                              setSelectedInstance(instance);
                              setShowConnectionDialog(true);
                            }}
                          >
                            <ExternalLink className="h-4 w-4 mr-1" />
                            Connection
                          </Button>
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleDeleteInstance(instance.id)}
                          >
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

        {/* Connection URL Dialog */}
        <Dialog
          open={showConnectionDialog}
          onOpenChange={setShowConnectionDialog}
        >
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Connection Information</DialogTitle>
              <DialogDescription>
                Use these details to connect to your database instance
              </DialogDescription>
            </DialogHeader>
            {selectedInstance && (
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label>Host</Label>
                    <Input value={selectedInstance.host} readOnly />
                  </div>
                  <div>
                    <Label>Port</Label>
                    <Input value={selectedInstance.port} readOnly />
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label>Database</Label>
                    <Input value={selectedInstance.database} readOnly />
                  </div>
                  <div>
                    <Label>Username</Label>
                    <Input value={selectedInstance.username} readOnly />
                  </div>
                </div>
                <div>
                  <Label>Connection URL</Label>
                  <div className="flex space-x-2">
                    <Input
                      value={generateConnectionUrl(selectedInstance)}
                      readOnly
                      className="font-mono text-sm"
                    />
                    <Button
                      variant="outline"
                      onClick={() =>
                        copyToClipboard(generateConnectionUrl(selectedInstance))
                      }
                    >
                      <Copy className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
                {selectedInstance.volumes && (
                  <div>
                    <Label>Volumes</Label>
                    <div className="space-y-2">
                      {selectedInstance.volumes.map((volume, index) => (
                        <div
                          key={index}
                          className="flex justify-between text-sm"
                        >
                          <span>{volume.name}</span>
                          <span>{volume.size}MB</span>
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
