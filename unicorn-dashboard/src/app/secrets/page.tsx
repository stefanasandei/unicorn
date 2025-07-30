'use client';

import React, { useEffect, useState } from 'react';
import { Layout } from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogFooter, 
  DialogHeader, 
  DialogTitle, 
  DialogTrigger 
} from '@/components/ui/dialog';
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from '@/components/ui/table';
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
} from '@/components/ui/alert-dialog';
import { Database, Plus, Eye, EyeOff, Edit, Trash2, Copy, Check } from 'lucide-react';
import { apiClient } from '@/lib/api';
import { Secret } from '@/types/api';

export default function SecretsPage() {
  const [secrets, setSecrets] = useState<Secret[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [showValues, setShowValues] = useState<Record<string, boolean>>({});
  const [copiedSecret, setCopiedSecret] = useState<string | null>(null);

  // Form states
  const [newSecretName, setNewSecretName] = useState('');
  const [newSecretValue, setNewSecretValue] = useState('');
  const [newSecretMetadata, setNewSecretMetadata] = useState('');
  const [editingSecret, setEditingSecret] = useState<Secret | null>(null);
  const [editSecretValue, setEditSecretValue] = useState('');
  const [editSecretMetadata, setEditSecretMetadata] = useState('');

  useEffect(() => {
    fetchSecrets();
  }, []);

  const fetchSecrets = async () => {
    try {
      setIsLoading(true);
      const secretsData = await apiClient.listSecrets();
      setSecrets(secretsData);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch secrets');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateSecret = async () => {
    if (!newSecretName.trim() || !newSecretValue.trim()) {
      setError('Name and value are required');
      return;
    }

    try {
      let metadata = {};
      if (newSecretMetadata.trim()) {
        try {
          metadata = JSON.parse(newSecretMetadata);
        } catch {
          setError('Invalid JSON metadata');
          return;
        }
      }

      await apiClient.createSecret({
        name: newSecretName,
        value: newSecretValue,
        metadata,
      });
      
      setNewSecretName('');
      setNewSecretValue('');
      setNewSecretMetadata('');
      fetchSecrets();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create secret');
    }
  };

  const handleUpdateSecret = async () => {
    if (!editingSecret) return;

    try {
      let metadata = {};
      if (editSecretMetadata.trim()) {
        try {
          metadata = JSON.parse(editSecretMetadata);
        } catch {
          setError('Invalid JSON metadata');
          return;
        }
      }

      await apiClient.updateSecret(editingSecret.id, {
        value: editSecretValue,
        metadata,
      });
      
      setEditingSecret(null);
      setEditSecretValue('');
      setEditSecretMetadata('');
      fetchSecrets();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to update secret');
    }
  };

  const handleDeleteSecret = async (secretId: string) => {
    try {
      await apiClient.deleteSecret(secretId);
      fetchSecrets();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to delete secret');
    }
  };

  const handleCopySecret = async (secretId: string) => {
    try {
      const secret = await apiClient.getSecret(secretId);
      await navigator.clipboard.writeText(secret.value);
      setCopiedSecret(secretId);
      setTimeout(() => setCopiedSecret(null), 2000);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to copy secret');
    }
  };

  const toggleSecretValue = (secretId: string) => {
    setShowValues(prev => ({
      ...prev,
      [secretId]: !prev[secretId]
    }));
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
                <CardTitle>Secrets</CardTitle>
                <CardDescription>
                  Store and manage encrypted secrets securely
                </CardDescription>
              </div>
              <Dialog>
                <DialogTrigger asChild>
                  <Button>
                    <Plus className="h-4 w-4 mr-2" />
                    Create Secret
                  </Button>
                </DialogTrigger>
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
                    </div>
                    <div>
                      <Label htmlFor="secret-metadata">Metadata (JSON)</Label>
                      <Input
                        id="secret-metadata"
                        value={newSecretMetadata}
                        onChange={(e) => setNewSecretMetadata(e.target.value)}
                        placeholder='{"key": "value"}'
                      />
                    </div>
                  </div>
                  <DialogFooter>
                    <Button onClick={handleCreateSecret}>Create Secret</Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
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
                          {showValues[secret.id] ? '••••••••' : '••••••••'}
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
                          {Object.keys(secret.metadata).length} keys
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
                        <Dialog>
                          <DialogTrigger asChild>
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => setEditingSecret(secret)}
                            >
                              <Edit className="h-4 w-4" />
                            </Button>
                          </DialogTrigger>
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
                                  value={secret.name}
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
                              </div>
                              <div>
                                <Label htmlFor="edit-secret-metadata">Metadata (JSON)</Label>
                                <Input
                                  id="edit-secret-metadata"
                                  value={editSecretMetadata}
                                  onChange={(e) => setEditSecretMetadata(e.target.value)}
                                  placeholder='{"key": "value"}'
                                />
                              </div>
                            </div>
                            <DialogFooter>
                              <Button onClick={handleUpdateSecret}>Update Secret</Button>
                            </DialogFooter>
                          </DialogContent>
                        </Dialog>
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
                                Are you sure you want to delete this secret? This action cannot be undone.
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
      </div>
    </Layout>
  );
} 