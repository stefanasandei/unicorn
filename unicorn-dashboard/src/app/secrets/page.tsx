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
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
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
} from "lucide-react";
import { apiClient } from "@/lib/api";
import { Secret } from "@/types/api";

export default function SecretsPage() {
  const [secrets, setSecrets] = useState<Secret[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [showValues, setShowValues] = useState<Record<string, boolean>>({});
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
    } catch (err: any) {
      const errorMsg =
        err.response?.data?.message ||
        err.response?.data?.details ||
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
    } catch (err: any) {
      console.error("Create secret error:", err);
      console.error("Error response:", err.response?.data);
      const errorMsg =
        err.response?.data?.message ||
        err.response?.data?.details ||
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
    } catch (err: any) {
      const errorMsg =
        err.response?.data?.message ||
        err.response?.data?.details ||
        "Failed to update secret";
      setErrorMessage(errorMsg);
      setShowErrorDialog(true);
    }
  };

  const handleDeleteSecret = async (secretId: string) => {
    try {
      await apiClient.deleteSecret(secretId);
      fetchSecrets();
    } catch (err: any) {
      const errorMsg =
        err.response?.data?.message ||
        err.response?.data?.details ||
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
    } catch (err: any) {
      const errorMsg =
        err.response?.data?.message ||
        err.response?.data?.details ||
        "Failed to copy secret";
      setErrorMessage(errorMsg);
      setShowErrorDialog(true);
    }
  };

  const toggleSecretValue = (secretId: string) => {
    setShowValues((prev) => ({
      ...prev,
      [secretId]: !prev[secretId],
    }));
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
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Secrets Manager</h1>
          <p className="text-gray-600">
            Manage encrypted secrets for your applications
          </p>
          {apiConnected !== null && (
            <div
              className={`mt-2 text-sm ${
                apiConnected ? "text-green-600" : "text-red-600"
              }`}
            >
              API Status: {apiConnected ? "Connected" : "Not Connected"}
            </div>
          )}
        </div>

        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Secrets</CardTitle>
                <CardDescription>
                  Store and manage encrypted secrets securely
                </CardDescription>
              </div>
              <Button onClick={() => setShowCreateDialog(true)}>
                <Plus className="h-4 w-4 mr-2" />
                Create Secret
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Value</TableHead>
                  <TableHead>Metadata</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {secrets.map((secret) => (
                  <TableRow key={secret.id}>
                    <TableCell className="font-medium">{secret.name}</TableCell>
                    <TableCell>
                      <div className="flex items-center space-x-2">
                        <span className="font-mono text-sm">
                          {showValues[secret.id] ? "••••••••" : "••••••••"}
                        </span>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => toggleSecretValue(secret.id)}
                        >
                          {showValues[secret.id] ? (
                            <EyeOff className="h-4 w-4" />
                          ) : (
                            <Eye className="h-4 w-4" />
                          )}
                        </Button>
                      </div>
                    </TableCell>
                    <TableCell>
                      {secret.metadata ? (
                        <Badge variant="secondary">
                          {(() => {
                            try {
                              const parsed = JSON.parse(secret.metadata);
                              return Object.keys(parsed).length;
                            } catch {
                              return 1;
                            }
                          })()}{" "}
                          keys
                        </Badge>
                      ) : (
                        <span className="text-gray-400">None</span>
                      )}
                    </TableCell>
                    <TableCell>
                      {new Date(secret.created_at).toLocaleDateString()}
                    </TableCell>
                    <TableCell>
                      <div className="flex space-x-2">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleCopySecret(secret.id)}
                        >
                          {copiedSecret === secret.id ? (
                            <Check className="h-4 w-4" />
                          ) : (
                            <Copy className="h-4 w-4" />
                          )}
                        </Button>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => openEditDialog(secret)}
                        >
                          <Edit className="h-4 w-4" />
                        </Button>
                        <AlertDialog>
                          <AlertDialogTrigger asChild>
                            <Button variant="outline" size="sm">
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </AlertDialogTrigger>
                          <AlertDialogContent>
                            <AlertDialogHeader>
                              <AlertDialogTitle>Delete Secret</AlertDialogTitle>
                              <AlertDialogDescription>
                                Are you sure you want to delete this secret?
                                This action cannot be undone.
                              </AlertDialogDescription>
                            </AlertDialogHeader>
                            <AlertDialogFooter>
                              <AlertDialogCancel>Cancel</AlertDialogCancel>
                              <AlertDialogAction
                                onClick={() => handleDeleteSecret(secret.id)}
                                className="bg-red-600 hover:bg-red-700"
                              >
                                Delete
                              </AlertDialogAction>
                            </AlertDialogFooter>
                          </AlertDialogContent>
                        </AlertDialog>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        {/* Create Secret Dialog */}
        <Dialog open={showCreateDialog} onOpenChange={setShowCreateDialog}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create New Secret</DialogTitle>
              <DialogDescription>
                Create a new encrypted secret
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label htmlFor="secret-name">Secret Name</Label>
                <Input
                  id="secret-name"
                  value={newSecretName}
                  onChange={(e) => setNewSecretName(e.target.value)}
                  placeholder="Enter secret name"
                />
                <p className="text-sm text-gray-500 mt-1">
                  Only alphanumeric characters, hyphens, and underscores. Max 50
                  characters.
                </p>
              </div>
              <div>
                <Label htmlFor="secret-value">Secret Value</Label>
                <Input
                  id="secret-value"
                  type="password"
                  value={newSecretValue}
                  onChange={(e) => setNewSecretValue(e.target.value)}
                  placeholder="Enter secret value"
                />
                <p className="text-sm text-gray-500 mt-1">
                  Max 10KB. This will be encrypted and stored securely.
                </p>
              </div>
              <div>
                <Label htmlFor="secret-metadata">
                  Metadata (JSON - Optional)
                </Label>
                <Input
                  id="secret-metadata"
                  value={newSecretMetadata}
                  onChange={(e) => setNewSecretMetadata(e.target.value)}
                  placeholder='{"key": "value"}'
                />
                <p className="text-sm text-gray-500 mt-1">
                  Optional JSON metadata for the secret
                </p>
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={closeCreateDialog}>
                Cancel
              </Button>
              <Button onClick={handleCreateSecret}>Create Secret</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Edit Secret Dialog */}
        <Dialog open={showEditDialog} onOpenChange={setShowEditDialog}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Edit Secret</DialogTitle>
              <DialogDescription>
                Update the secret value and metadata
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label>Secret Name</Label>
                <Input
                  value={editingSecret?.name || ""}
                  disabled
                  className="bg-gray-50"
                />
              </div>
              <div>
                <Label htmlFor="edit-secret-value">Secret Value</Label>
                <Input
                  id="edit-secret-value"
                  type="password"
                  value={editSecretValue}
                  onChange={(e) => setEditSecretValue(e.target.value)}
                  placeholder="Enter new secret value"
                />
                <p className="text-sm text-gray-500 mt-1">
                  Max 10KB. Leave empty to keep current value.
                </p>
              </div>
              <div>
                <Label htmlFor="edit-secret-metadata">
                  Metadata (JSON - Optional)
                </Label>
                <Input
                  id="edit-secret-metadata"
                  value={editSecretMetadata}
                  onChange={(e) => setEditSecretMetadata(e.target.value)}
                  placeholder='{"key": "value"}'
                />
                <p className="text-sm text-gray-500 mt-1">
                  Optional JSON metadata for the secret
                </p>
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={closeEditDialog}>
                Cancel
              </Button>
              <Button onClick={handleUpdateSecret}>Update Secret</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Error Dialog */}
        <Dialog open={showErrorDialog} onOpenChange={setShowErrorDialog}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2">
                <AlertCircle className="h-5 w-5 text-red-500" />
                Error
              </DialogTitle>
            </DialogHeader>
            <div className="py-4">
              <p className="text-gray-700">{errorMessage}</p>
            </div>
            <DialogFooter>
              <Button onClick={() => setShowErrorDialog(false)}>OK</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </Layout>
  );
}
