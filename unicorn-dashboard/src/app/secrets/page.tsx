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
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import {
  Database,
  Plus,
  Eye,
  EyeOff,
  Edit,
  Trash2,
  Copy,
  Check,
  AlertCircle,
  Lock,
  Shield,
  Key,
  Zap,
  Clock
} from "lucide-react";
import { apiClient } from "@/lib/api";
import { Secret } from "@/types/api";

export default function SecretsPage() {
  const [secrets, setSecrets] = useState<Secret[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [showValues, setShowValues] = useState<Record<string, boolean>>({});
  const [revealedSecrets, setRevealedSecrets] = useState<
    Record<string, string>
  >({});
  const [loadingSecrets, setLoadingSecrets] = useState<Record<string, boolean>>(
    {}
  );
  const [copiedSecret, setCopiedSecret] = useState<string | null>(null);

  // Form states
  const [newSecretName, setNewSecretName] = useState("");
  const [newSecretValue, setNewSecretValue] = useState("");
  const [newSecretMetadata, setNewSecretMetadata] = useState("");
  const [editingSecret, setEditingSecret] = useState<Secret | null>(null);
  const [editSecretValue, setEditSecretValue] = useState("");
  const [editSecretMetadata, setEditSecretMetadata] = useState("");

  // Dialog states
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const [showEditDialog, setShowEditDialog] = useState(false);
  const [showErrorDialog, setShowErrorDialog] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");
  const [apiConnected, setApiConnected] = useState<boolean | null>(null);

  useEffect(() => {
    fetchSecrets();
    testAPIConnection();
  }, []);

  const testAPIConnection = async () => {
    try {
      const response = await fetch("http://localhost:8080/health");
      console.log("API health check response:", response.status);
      setApiConnected(response.ok);
    } catch (err) {
      console.error("API health check failed:", err);
      setApiConnected(false);
    }
  };

  const fetchSecrets = async () => {
    try {
      setIsLoading(true);
      const secretsData = await apiClient.listSecrets();
      setSecrets(secretsData);
    } catch (err: unknown) {
      const error = err as { response?: { data?: { message?: string; details?: string } } };
      const errorMsg =
        error.response?.data?.message ||
        error.response?.data?.details ||
        "Failed to fetch secrets";
      setErrorMessage(errorMsg);
      setShowErrorDialog(true);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateSecret = async () => {
    if (!newSecretName.trim() || !newSecretValue.trim()) {
      setErrorMessage("Name and value are required");
      setShowErrorDialog(true);
      return;
    }

    // Check authentication
    const token = localStorage.getItem("token");
    if (!token) {
      setErrorMessage("You must be logged in to create secrets");
      setShowErrorDialog(true);
      return;
    }

    console.log("Token exists:", !!token);
    console.log("Token length:", token.length);

    // Validate secret name format
    const nameRegex = /^[a-zA-Z0-9_-]+$/;
    if (!nameRegex.test(newSecretName.trim())) {
      setErrorMessage(
        "Secret name can only contain alphanumeric characters, hyphens, and underscores"
      );
      setShowErrorDialog(true);
      return;
    }

    if (newSecretName.trim().length > 50) {
      setErrorMessage("Secret name must be less than 50 characters");
      setShowErrorDialog(true);
      return;
    }

    // Validate secret value length
    if (newSecretValue.length > 10000) {
      setErrorMessage("Secret value must be less than 10KB");
      setShowErrorDialog(true);
      return;
    }

    try {
      // Validate JSON metadata if provided
      if (newSecretMetadata.trim()) {
        try {
          JSON.parse(newSecretMetadata);
        } catch {
          setErrorMessage("Invalid JSON metadata format");
          setShowErrorDialog(true);
          return;
        }
      }

      const requestData = {
        name: newSecretName,
        value: newSecretValue,
        ...(newSecretMetadata.trim() && { metadata: newSecretMetadata.trim() }),
      };

      console.log("Creating secret with data:", requestData);

      await apiClient.createSecret(requestData);

      setNewSecretName("");
      setNewSecretValue("");
      setNewSecretMetadata("");
      setShowCreateDialog(false);
      fetchSecrets();
    } catch (err: unknown) {
      console.error("Create secret error:", err);
      console.error("Error response:", (err as any).response?.data);
      const error = err as { response?: { data?: { message?: string; details?: string } } };
      const errorMsg =
        error.response?.data?.message ||
        error.response?.data?.details ||
        "Failed to create secret";
      setErrorMessage(errorMsg);
      setShowErrorDialog(true);
    }
  };

  const handleUpdateSecret = async () => {
    if (!editingSecret) return;

    // Validate secret value length if provided
    if (editSecretValue && editSecretValue.length > 10000) {
      setErrorMessage("Secret value must be less than 10KB");
      setShowErrorDialog(true);
      return;
    }

    try {
      // Validate JSON metadata if provided
      if (editSecretMetadata.trim()) {
        try {
          JSON.parse(editSecretMetadata);
        } catch {
          setErrorMessage("Invalid JSON metadata format");
          setShowErrorDialog(true);
          return;
        }
      }

      await apiClient.updateSecret(editingSecret.id, {
        value: editSecretValue,
        ...(editSecretMetadata.trim() && {
          metadata: editSecretMetadata.trim(),
        }),
      });

      setEditingSecret(null);
      setEditSecretValue("");
      setEditSecretMetadata("");
      setShowEditDialog(false);
      fetchSecrets();
    } catch (err: unknown) {
      const error = err as { response?: { data?: { message?: string; details?: string } } };
      const errorMsg =
        error.response?.data?.message ||
        error.response?.data?.details ||
        "Failed to update secret";
      setErrorMessage(errorMsg);
      setShowErrorDialog(true);
    }
  };

  const handleDeleteSecret = async (secretId: string) => {
    try {
      await apiClient.deleteSecret(secretId);
      fetchSecrets();
    } catch (err: unknown) {
      const error = err as { response?: { data?: { message?: string; details?: string } } };
      const errorMsg =
        error.response?.data?.message ||
        error.response?.data?.details ||
        "Failed to delete secret";
      setErrorMessage(errorMsg);
      setShowErrorDialog(true);
    }
  };

  const handleCopySecret = async (secretId: string) => {
    try {
      const secret = await apiClient.getSecret(secretId);
      await navigator.clipboard.writeText(secret.value);
      setCopiedSecret(secretId);
      setTimeout(() => setCopiedSecret(null), 2000);
    } catch (err: unknown) {
      const error = err as { response?: { data?: { message?: string; details?: string } } };
      const errorMsg =
        error.response?.data?.message ||
        error.response?.data?.details ||
        "Failed to copy secret";
      setErrorMessage(errorMsg);
      setShowErrorDialog(true);
    }
  };

  const toggleSecretValue = async (secretId: string) => {
    const isCurrentlyShown = showValues[secretId];

    if (isCurrentlyShown) {
      // Hide the secret
      setShowValues((prev) => ({
        ...prev,
        [secretId]: false,
      }));
      setRevealedSecrets((prev) => {
        const newState = { ...prev };
        delete newState[secretId];
        return newState;
      });
    } else {
      // Show the secret - fetch the actual value
      try {
        setLoadingSecrets((prev) => ({
          ...prev,
          [secretId]: true,
        }));

        const secret = await apiClient.getSecret(secretId);
        setRevealedSecrets((prev) => ({
          ...prev,
          [secretId]: secret.value,
        }));
        setShowValues((prev) => ({
          ...prev,
          [secretId]: true,
        }));
      } catch (err: unknown) {
        const error = err as { response?: { data?: { message?: string; details?: string } } };
        const errorMsg =
          error.response?.data?.message ||
          error.response?.data?.details ||
          "Failed to reveal secret";
        setErrorMessage(errorMsg);
        setShowErrorDialog(true);

        // Reset the toggle if there was an error
        setShowValues((prev) => ({
          ...prev,
          [secretId]: false,
        }));
      } finally {
        setLoadingSecrets((prev) => ({
          ...prev,
          [secretId]: false,
        }));
      }
    }
  };

  const openEditDialog = (secret: Secret) => {
    setEditingSecret(secret);
    setEditSecretValue("");
    setEditSecretMetadata("");
    setShowEditDialog(true);
  };

  const closeEditDialog = () => {
    setEditingSecret(null);
    setEditSecretValue("");
    setEditSecretMetadata("");
    setShowEditDialog(false);
  };

  const closeCreateDialog = () => {
    setNewSecretName("");
    setNewSecretValue("");
    setNewSecretMetadata("");
    setShowCreateDialog(false);
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
              <Lock className="h-6 w-6 text-primary" />
            </div>
            <div>
              <h1 className="text-3xl font-bold text-foreground">Secrets Manager</h1>
              <p className="text-muted-foreground">
                Manage encrypted secrets for your applications
              </p>
            </div>
          </div>
          {apiConnected !== null && (
            <div className="flex items-center space-x-2">
              <div className={`w-2 h-2 rounded-full ${apiConnected ? 'bg-green-500' : 'bg-red-500'}`}></div>
              <span className={`text-sm ${apiConnected ? 'text-green-600' : 'text-red-600'}`}>
                API Status: {apiConnected ? "Connected" : "Not Connected"}
              </span>
            </div>
          )}
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">Total Secrets</CardTitle>
              <div className="p-2 rounded-lg bg-blue-500/10">
                <Database className="h-4 w-4 text-blue-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">{secrets.length}</div>
              <p className="text-xs text-muted-foreground">Encrypted secrets</p>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">With Metadata</CardTitle>
              <div className="p-2 rounded-lg bg-green-500/10">
                <Key className="h-4 w-4 text-green-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">
                {secrets.filter(s => s.metadata).length}
              </div>
              <p className="text-xs text-muted-foreground">Secrets with metadata</p>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">Security</CardTitle>
              <div className="p-2 rounded-lg bg-purple-500/10">
                <Shield className="h-4 w-4 text-purple-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">AES-256</div>
              <p className="text-xs text-muted-foreground">Encryption standard</p>
            </CardContent>
          </Card>
        </div>

        <Card className="border-border/50">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-foreground">Secrets</CardTitle>
                <CardDescription>
                  Store and manage encrypted secrets securely
                </CardDescription>
              </div>
              <Button onClick={() => setShowCreateDialog(true)} className="bg-primary hover:bg-primary/90">
                <Plus className="h-4 w-4 mr-2" />
                Create Secret
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {secrets.map((secret) => (
                <div key={secret.id} className="border border-border/50 rounded-lg p-4 hover:shadow-md transition-all duration-200">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-3">
                      <div className="p-2 rounded-lg bg-accent/50">
                        <Lock className="h-4 w-4 text-accent-foreground" />
                      </div>
                      <div>
                        <h3 className="font-medium text-foreground">{secret.name}</h3>
                        <p className="text-sm text-muted-foreground flex items-center gap-1">
                          <Clock className="h-3 w-3" />
                          Created {new Date(secret.created_at).toLocaleDateString()}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center space-x-2">
                      <Badge variant="outline" className="bg-secondary/20">
                        {secret.metadata ? (
                          (() => {
                            try {
                              const parsed = JSON.parse(secret.metadata);
                              return `${Object.keys(parsed).length} keys`;
                            } catch {
                              return "1 key";
                            }
                          })()
                        ) : (
                          "No metadata"
                        )}
                      </Badge>
                      <div className="flex space-x-1">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleCopySecret(secret.id)}
                          className="border-border/50"
                        >
                          {copiedSecret === secret.id ? (
                            <Check className="h-4 w-4 text-green-500" />
                          ) : (
                            <Copy className="h-4 w-4" />
                          )}
                        </Button>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => openEditDialog(secret)}
                          className="border-border/50"
                        >
                          <Edit className="h-4 w-4" />
                        </Button>
                        <AlertDialog>
                          <AlertDialogTrigger asChild>
                            <Button variant="outline" size="sm" className="border-border/50">
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </AlertDialogTrigger>
                          <AlertDialogContent className="border-border/50">
                            <AlertDialogHeader>
                              <AlertDialogTitle className="text-foreground">Delete Secret</AlertDialogTitle>
                              <AlertDialogDescription>
                                Are you sure you want to delete this secret?
                                This action cannot be undone.
                              </AlertDialogDescription>
                            </AlertDialogHeader>
                            <AlertDialogFooter>
                              <AlertDialogCancel>Cancel</AlertDialogCancel>
                              <AlertDialogAction
                                onClick={() => handleDeleteSecret(secret.id)}
                                className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                              >
                                Delete
                              </AlertDialogAction>
                            </AlertDialogFooter>
                          </AlertDialogContent>
                        </AlertDialog>
                      </div>
                    </div>
                  </div>
                  <div className="mt-3 flex items-center space-x-2">
                    <span className="font-mono text-sm bg-muted/50 px-2 py-1 rounded">
                      {loadingSecrets[secret.id]
                        ? "Loading..."
                        : showValues[secret.id]
                        ? revealedSecrets[secret.id] || "••••••••"
                        : "••••••••"}
                    </span>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => toggleSecretValue(secret.id)}
                      disabled={loadingSecrets[secret.id]}
                      className="text-muted-foreground hover:text-foreground"
                    >
                      {loadingSecrets[secret.id] ? (
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-primary"></div>
                      ) : showValues[secret.id] ? (
                        <EyeOff className="h-4 w-4" />
                      ) : (
                        <Eye className="h-4 w-4" />
                      )}
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Create Secret Dialog */}
        <Dialog open={showCreateDialog} onOpenChange={setShowCreateDialog}>
          <DialogContent className="border-border/50">
            <DialogHeader>
              <DialogTitle className="text-foreground">Create New Secret</DialogTitle>
              <DialogDescription>
                Create a new encrypted secret
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label htmlFor="secret-name" className="text-foreground">Secret Name</Label>
                <Input
                  id="secret-name"
                  value={newSecretName}
                  onChange={(e) => setNewSecretName(e.target.value)}
                  placeholder="Enter secret name"
                  className="border-border/50 focus:border-primary focus:ring-primary/20"
                />
                <p className="text-sm text-muted-foreground mt-1">
                  Only alphanumeric characters, hyphens, and underscores. Max 50 characters.
                </p>
              </div>
              <div>
                <Label htmlFor="secret-value" className="text-foreground">Secret Value</Label>
                <Input
                  id="secret-value"
                  type="password"
                  value={newSecretValue}
                  onChange={(e) => setNewSecretValue(e.target.value)}
                  placeholder="Enter secret value"
                  className="border-border/50 focus:border-primary focus:ring-primary/20"
                />
                <p className="text-sm text-muted-foreground mt-1">
                  Max 10KB. This will be encrypted and stored securely.
                </p>
              </div>
              <div>
                <Label htmlFor="secret-metadata" className="text-foreground">
                  Metadata (JSON - Optional)
                </Label>
                <Input
                  id="secret-metadata"
                  value={newSecretMetadata}
                  onChange={(e) => setNewSecretMetadata(e.target.value)}
                  placeholder='{"key": "value"}'
                  className="border-border/50 focus:border-primary focus:ring-primary/20"
                />
                <p className="text-sm text-muted-foreground mt-1">
                  Optional JSON metadata for the secret
                </p>
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={closeCreateDialog} className="border-border/50">
                Cancel
              </Button>
              <Button onClick={handleCreateSecret} className="bg-primary hover:bg-primary/90">
                Create Secret
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Edit Secret Dialog */}
        <Dialog open={showEditDialog} onOpenChange={setShowEditDialog}>
          <DialogContent className="border-border/50">
            <DialogHeader>
              <DialogTitle className="text-foreground">Edit Secret</DialogTitle>
              <DialogDescription>
                Update the secret value and metadata
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label className="text-foreground">Secret Name</Label>
                <Input
                  value={editingSecret?.name || ""}
                  disabled
                  className="bg-muted/50 border-border/50"
                />
              </div>
              <div>
                <Label htmlFor="edit-secret-value" className="text-foreground">Secret Value</Label>
                <Input
                  id="edit-secret-value"
                  type="password"
                  value={editSecretValue}
                  onChange={(e) => setEditSecretValue(e.target.value)}
                  placeholder="Enter new secret value"
                  className="border-border/50 focus:border-primary focus:ring-primary/20"
                />
                <p className="text-sm text-muted-foreground mt-1">
                  Max 10KB. Leave empty to keep current value.
                </p>
              </div>
              <div>
                <Label htmlFor="edit-secret-metadata" className="text-foreground">
                  Metadata (JSON - Optional)
                </Label>
                <Input
                  id="edit-secret-metadata"
                  value={editSecretMetadata}
                  onChange={(e) => setEditSecretMetadata(e.target.value)}
                  placeholder='{"key": "value"}'
                  className="border-border/50 focus:border-primary focus:ring-primary/20"
                />
                <p className="text-sm text-muted-foreground mt-1">
                  Optional JSON metadata for the secret
                </p>
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={closeEditDialog} className="border-border/50">
                Cancel
              </Button>
              <Button onClick={handleUpdateSecret} className="bg-primary hover:bg-primary/90">
                Update Secret
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Error Dialog */}
        <Dialog open={showErrorDialog} onOpenChange={setShowErrorDialog}>
          <DialogContent className="border-border/50">
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2 text-foreground">
                <AlertCircle className="h-5 w-5 text-destructive" />
                Error
              </DialogTitle>
            </DialogHeader>
            <div className="py-4">
              <p className="text-foreground">{errorMessage}</p>
            </div>
            <DialogFooter>
              <Button onClick={() => setShowErrorDialog(false)} className="bg-primary hover:bg-primary/90">
                OK
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </Layout>
  );
}
