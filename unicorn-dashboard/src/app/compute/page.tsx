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
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Server,
  Plus,
  Play,
  Square,
  Activity,
  RefreshCw,
  Trash2,
  Cpu,
  Zap,
  Clock,
  AlertCircle,
  Database,
  Globe,
  Sparkles,
} from "lucide-react";
import { apiClient } from "@/lib/api";
import { ComputeContainerInfo, ComputeCreateRequest } from "@/types/api";

export default function ComputePage() {
  const [containers, setContainers] = useState<ComputeContainerInfo[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isCreating, setIsCreating] = useState(false);
  const [error, setError] = useState("");

  // Form states
  const [newContainerName, setNewContainerName] = useState("");
  const [newContainerImage, setNewContainerImage] = useState("");
  const [newContainerCommand, setNewContainerCommand] = useState("");
  const [newContainerEnvironment, setNewContainerEnvironment] = useState("");
  const [newContainerPort, setNewContainerPort] = useState("80");
  const [newContainerHostPort, setNewContainerHostPort] = useState("8080");

  useEffect(() => {
    fetchContainers();
  }, []);

  const fetchContainers = async () => {
    try {
      setIsLoading(true);
      setError(""); // Clear any previous errors
      const containersData = await apiClient.listCompute();
      setContainers(containersData);
    } catch (err: unknown) {
      console.error("Container fetch error:", err);
      const error = err as {
        response?: { data?: { error?: string; details?: string } };
        message?: string;
      };
      const errorMessage =
        error.response?.data?.error ||
        error.response?.data?.details ||
        error.message ||
        "Failed to fetch containers";
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateContainer = async () => {
    if (!newContainerImage.trim()) {
      setError("Image is required");
      return;
    }

    try {
      setIsCreating(true);
      setError("");

      let environment = {};
      if (newContainerEnvironment.trim()) {
        try {
          environment = JSON.parse(newContainerEnvironment);
        } catch {
          setError("Invalid JSON environment");
          return;
        }
      }

      // Use user-defined ports
      const exposePort = newContainerPort;
      const ports: Record<string, string> = {};
      ports[exposePort] = newContainerHostPort;

      const request: ComputeCreateRequest = {
        name: newContainerName.trim() || undefined,
        image: newContainerImage,
        environment,
        preset: "micro", // Default to micro preset
        expose_port: exposePort,
        ports,
      };

      if (newContainerCommand.trim()) {
        request.command = newContainerCommand.split(" ");
      }

      await apiClient.createCompute(request);

      setNewContainerName("");
      setNewContainerImage("");
      setNewContainerCommand("");
      setNewContainerEnvironment("");
      setNewContainerPort("80");
      setNewContainerHostPort("8080");
      setError(""); // Clear any previous errors
      fetchContainers();
    } catch (err: unknown) {
      console.error("Container creation error:", err);
      const error = err as {
        response?: { data?: { error?: string; details?: string } };
        message?: string;
      };
      const errorMessage =
        error.response?.data?.error ||
        error.response?.data?.details ||
        error.message ||
        "Failed to create container";
      setError(errorMessage);
    } finally {
      setIsCreating(false);
    }
  };

  const handleDeleteContainer = async (containerId: string) => {
    if (!confirm("Are you sure you want to delete this container?")) {
      return;
    }

    try {
      await apiClient.deleteCompute(containerId);
      fetchContainers(); // Refresh the list
    } catch (err: unknown) {
      console.error("Container deletion error:", err);
      const error = err as {
        response?: { data?: { error?: string; details?: string } };
        message?: string;
      };
      const errorMessage =
        error.response?.data?.error ||
        error.response?.data?.details ||
        error.message ||
        "Failed to delete container";
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

  const getStatusIcon = (status: string) => {
    switch (status.toLowerCase()) {
      case "running":
        return <Play className="h-3 w-3" />;
      case "stopped":
        return <Square className="h-3 w-3" />;
      case "starting":
        return <RefreshCw className="h-3 w-3 animate-spin" />;
      default:
        return <Activity className="h-3 w-3" />;
    }
  };

  const runningContainers =
    containers?.filter((c) => c.status?.toLowerCase() === "running")?.length ||
    0;
  const totalContainers = containers?.length || 0;

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
              <Cpu className="h-6 w-6 text-primary" />
            </div>
            <div>
              <h1 className="text-3xl font-bold text-foreground">Compute</h1>
              <p className="text-muted-foreground">
                Manage your compute containers and resources
              </p>
            </div>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">
                Total Containers
              </CardTitle>
              <div className="p-2 rounded-lg bg-blue-500/10">
                <Server className="h-4 w-4 text-blue-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">
                {totalContainers}
              </div>
              <p className="text-xs text-muted-foreground">
                Deployed containers
              </p>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">
                Running
              </CardTitle>
              <div className="p-2 rounded-lg bg-green-500/10">
                <Play className="h-4 w-4 text-green-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">
                {runningContainers}
              </div>
              <p className="text-xs text-muted-foreground">Active containers</p>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">
                Performance
              </CardTitle>
              <div className="p-2 rounded-lg bg-purple-500/10">
                <Zap className="h-4 w-4 text-purple-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">Micro</div>
              <p className="text-xs text-muted-foreground">Default preset</p>
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
                <CardTitle className="text-foreground">Containers</CardTitle>
                <CardDescription>
                  Deploy and manage compute containers
                </CardDescription>
              </div>
              <div className="flex space-x-2">
                <Button
                  variant="outline"
                  onClick={fetchContainers}
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
                      Deploy Container
                    </Button>
                  </DialogTrigger>
                  <DialogContent className="border-border/50">
                    <DialogHeader>
                      <DialogTitle className="text-foreground">
                        Deploy New Container
                      </DialogTitle>
                      <DialogDescription>
                        Deploy a new compute container
                      </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4">
                      <div>
                        <Label
                          htmlFor="container-name"
                          className="text-foreground"
                        >
                          Container Name (optional)
                        </Label>
                        <Input
                          id="container-name"
                          value={newContainerName}
                          onChange={(e) => setNewContainerName(e.target.value)}
                          placeholder="Auto-generated if not provided"
                          className="border-border/50 focus:border-primary focus:ring-primary/20"
                        />
                      </div>
                      <div>
                        <Label
                          htmlFor="container-image"
                          className="text-foreground"
                        >
                          Docker Image
                        </Label>
                        <Input
                          id="container-image"
                          value={newContainerImage}
                          onChange={(e) => setNewContainerImage(e.target.value)}
                          placeholder="e.g., nginx:latest"
                          className="border-border/50 focus:border-primary focus:ring-primary/20"
                        />
                      </div>
                      <div>
                        <Label
                          htmlFor="container-command"
                          className="text-foreground"
                        >
                          Command (optional)
                        </Label>
                        <Input
                          id="container-command"
                          value={newContainerCommand}
                          onChange={(e) =>
                            setNewContainerCommand(e.target.value)
                          }
                          placeholder="e.g., nginx -g 'daemon off;'"
                          className="border-border/50 focus:border-primary focus:ring-primary/20"
                        />
                      </div>
                      <div>
                        <Label
                          htmlFor="container-environment"
                          className="text-foreground"
                        >
                          Environment Variables (JSON)
                        </Label>
                        <Input
                          id="container-environment"
                          value={newContainerEnvironment}
                          onChange={(e) =>
                            setNewContainerEnvironment(e.target.value)
                          }
                          placeholder='{"NODE_ENV": "production"}'
                          className="border-border/50 focus:border-primary focus:ring-primary/20"
                        />
                      </div>
                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <Label
                            htmlFor="container-port"
                            className="text-foreground"
                          >
                            Container Port
                          </Label>
                          <Input
                            id="container-port"
                            type="number"
                            value={newContainerPort}
                            onChange={(e) =>
                              setNewContainerPort(e.target.value)
                            }
                            placeholder="80"
                            className="border-border/50 focus:border-primary focus:ring-primary/20"
                          />
                        </div>
                        <div>
                          <Label
                            htmlFor="container-host-port"
                            className="text-foreground"
                          >
                            Host Port
                          </Label>
                          <Input
                            id="container-host-port"
                            type="number"
                            value={newContainerHostPort}
                            onChange={(e) =>
                              setNewContainerHostPort(e.target.value)
                            }
                            placeholder="8080"
                            className="border-border/50 focus:border-primary focus:ring-primary/20"
                          />
                        </div>
                      </div>
                    </div>
                    <DialogFooter>
                      <Button
                        onClick={handleCreateContainer}
                        disabled={isCreating}
                        className="bg-primary hover:bg-primary/90"
                      >
                        {isCreating ? (
                          <div className="flex items-center space-x-2">
                            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-primary-foreground"></div>
                            <span>Creating...</span>
                          </div>
                        ) : (
                          "Deploy Container"
                        )}
                      </Button>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {containers && containers.length > 0 ? (
                containers.map((container) => (
                  <div
                    key={container.id}
                    className="border border-border/50 rounded-lg p-4 hover:shadow-md transition-all duration-200"
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-3">
                        <div className="p-2 rounded-lg bg-accent/50">
                          <Server className="h-4 w-4 text-accent-foreground" />
                        </div>
                        <div>
                          <h3 className="font-medium text-foreground">
                            {container.name}
                          </h3>
                          <p className="text-sm text-muted-foreground flex items-center gap-1">
                            <Clock className="h-3 w-3" />
                            Created{" "}
                            {new Date(
                              container.created_at
                            ).toLocaleDateString()}
                          </p>
                        </div>
                      </div>
                      <div className="flex items-center space-x-3">
                        <Badge variant="outline" className="bg-secondary/20">
                          {container.image}
                        </Badge>
                        <Badge className={getStatusColor(container.status)}>
                          <div className="flex items-center space-x-1">
                            {getStatusIcon(container.status)}
                            <span>{container.status}</span>
                          </div>
                        </Badge>
                        <div className="flex space-x-1">
                          <Button
                            variant="outline"
                            size="sm"
                            className="border-border/50"
                          >
                            <Play className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="outline"
                            size="sm"
                            className="border-border/50"
                          >
                            <Square className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="outline"
                            size="sm"
                            className="border-border/50"
                          >
                            <Activity className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleDeleteContainer(container.id)}
                            className="border-border/50 text-destructive hover:text-destructive hover:bg-destructive/10"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </div>
                    </div>
                  </div>
                ))
              ) : (
                <div className="text-center py-8">
                  <Server className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
                  <h3 className="text-lg font-medium text-foreground mb-2">
                    No containers found
                  </h3>
                  <p className="text-muted-foreground">
                    Create your first container to get started.
                  </p>
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        {/* Quick Deploy Templates */}
        <Card className="border-border/50">
          <CardHeader>
            <CardTitle className="text-foreground">
              Quick Deploy Templates
            </CardTitle>
            <CardDescription>
              Common container configurations for quick deployment
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <Card
                className="cursor-pointer hover:shadow-md transition-all duration-200 border-border/50"
                onClick={() => {
                  setNewContainerImage("nginx:latest");
                  setNewContainerName("web-server");
                }}
              >
                <CardContent className="p-4">
                  <div className="flex items-center space-x-3">
                    <div className="p-2 rounded-lg bg-gradient-to-br from-blue-500 to-blue-600">
                      <Globe className="h-5 w-5 text-white" />
                    </div>
                    <div>
                      <h3 className="font-medium text-foreground">
                        Web Server
                      </h3>
                      <p className="text-sm text-muted-foreground">
                        nginx:latest
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card
                className="cursor-pointer hover:shadow-md transition-all duration-200 border-border/50"
                onClick={() => {
                  setNewContainerImage("postgres:13");
                  setNewContainerName("database");
                }}
              >
                <CardContent className="p-4">
                  <div className="flex items-center space-x-3">
                    <div className="p-2 rounded-lg bg-gradient-to-br from-green-500 to-green-600">
                      <Database className="h-5 w-5 text-white" />
                    </div>
                    <div>
                      <h3 className="font-medium text-foreground">Database</h3>
                      <p className="text-sm text-muted-foreground">
                        postgres:13
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card
                className="cursor-pointer hover:shadow-md transition-all duration-200 border-border/50"
                onClick={() => {
                  setNewContainerImage("redis:alpine");
                  setNewContainerName("cache");
                }}
              >
                <CardContent className="p-4">
                  <div className="flex items-center space-x-3">
                    <div className="p-2 rounded-lg bg-gradient-to-br from-purple-500 to-purple-600">
                      <Zap className="h-5 w-5 text-white" />
                    </div>
                    <div>
                      <h3 className="font-medium text-foreground">Cache</h3>
                      <p className="text-sm text-muted-foreground">
                        redis:alpine
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </CardContent>
        </Card>
      </div>
    </Layout>
  );
}
