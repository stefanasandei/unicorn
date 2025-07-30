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
import { Server, Plus, Play, Square, Activity } from 'lucide-react';
import { apiClient } from '@/lib/api';
import { ComputeContainerInfo, ComputeCreateRequest } from '@/types/api';

export default function ComputePage() {
  const [containers, setContainers] = useState<ComputeContainerInfo[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  // Form states
  const [newContainerName, setNewContainerName] = useState('');
  const [newContainerImage, setNewContainerImage] = useState('');
  const [newContainerCommand, setNewContainerCommand] = useState('');
  const [newContainerEnvironment, setNewContainerEnvironment] = useState('');

  useEffect(() => {
    fetchContainers();
  }, []);

  const fetchContainers = async () => {
    try {
      setIsLoading(true);
      const containersData = await apiClient.listCompute();
      setContainers(containersData);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch containers');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateContainer = async () => {
    if (!newContainerName.trim() || !newContainerImage.trim()) {
      setError('Name and image are required');
      return;
    }

    try {
      let environment = {};
      if (newContainerEnvironment.trim()) {
        try {
          environment = JSON.parse(newContainerEnvironment);
        } catch {
          setError('Invalid JSON environment');
          return;
        }
      }

      const request: ComputeCreateRequest = {
        name: newContainerName,
        image: newContainerImage,
        environment,
      };

      if (newContainerCommand.trim()) {
        request.command = newContainerCommand.split(' ');
      }

      await apiClient.createCompute(request);
      
      setNewContainerName('');
      setNewContainerImage('');
      setNewContainerCommand('');
      setNewContainerEnvironment('');
      fetchContainers();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create container');
    }
  };

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'running':
        return 'bg-green-100 text-green-800';
      case 'stopped':
        return 'bg-red-100 text-red-800';
      case 'starting':
        return 'bg-yellow-100 text-yellow-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
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
          <h1 className="text-2xl font-bold text-gray-900">Compute</h1>
          <p className="text-gray-600">
            Manage your compute containers and resources
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
                <CardTitle>Containers</CardTitle>
                <CardDescription>
                  Deploy and manage compute containers
                </CardDescription>
              </div>
              <Dialog>
                <DialogTrigger asChild>
                  <Button>
                    <Plus className="h-4 w-4 mr-2" />
                    Deploy Container
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Deploy New Container</DialogTitle>
                    <DialogDescription>
                      Deploy a new compute container
                    </DialogDescription>
                  </DialogHeader>
                  <div className="space-y-4">
                    <div>
                      <Label htmlFor="container-name">Container Name</Label>
                      <Input
                        id="container-name"
                        value={newContainerName}
                        onChange={(e) => setNewContainerName(e.target.value)}
                        placeholder="Enter container name"
                      />
                    </div>
                    <div>
                      <Label htmlFor="container-image">Docker Image</Label>
                      <Input
                        id="container-image"
                        value={newContainerImage}
                        onChange={(e) => setNewContainerImage(e.target.value)}
                        placeholder="e.g., nginx:latest"
                      />
                    </div>
                    <div>
                      <Label htmlFor="container-command">Command (optional)</Label>
                      <Input
                        id="container-command"
                        value={newContainerCommand}
                        onChange={(e) => setNewContainerCommand(e.target.value)}
                        placeholder="e.g., nginx -g 'daemon off;'"
                      />
                    </div>
                    <div>
                      <Label htmlFor="container-environment">Environment Variables (JSON)</Label>
                      <Input
                        id="container-environment"
                        value={newContainerEnvironment}
                        onChange={(e) => setNewContainerEnvironment(e.target.value)}
                        placeholder='{"NODE_ENV": "production"}'
                      />
                    </div>
                  </div>
                  <DialogFooter>
                    <Button onClick={handleCreateContainer}>Deploy Container</Button>
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
                  <TableHead>Image</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {containers.map((container) => (
                  <TableRow key={container.id}>
                    <TableCell className="font-medium">{container.name}</TableCell>
                    <TableCell>
                      <Badge variant="outline">{container.image}</Badge>
                    </TableCell>
                    <TableCell>
                      <Badge className={getStatusColor(container.status)}>
                        {container.status}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {new Date(container.created_at).toLocaleDateString()}
                    </TableCell>
                    <TableCell>
                      <div className="flex space-x-2">
                        <Button variant="outline" size="sm">
                          <Play className="h-4 w-4" />
                        </Button>
                        <Button variant="outline" size="sm">
                          <Square className="h-4 w-4" />
                        </Button>
                        <Button variant="outline" size="sm">
                          <Activity className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        {/* Quick Deploy Templates */}
        <Card>
          <CardHeader>
            <CardTitle>Quick Deploy Templates</CardTitle>
            <CardDescription>
              Common container configurations for quick deployment
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <Card className="cursor-pointer hover:shadow-md transition-shadow">
                <CardContent className="p-4">
                  <div className="flex items-center space-x-3">
                    <div className="p-2 rounded-lg bg-blue-500">
                      <Server className="h-5 w-5 text-white" />
                    </div>
                    <div>
                      <h3 className="font-medium">Web Server</h3>
                      <p className="text-sm text-muted-foreground">
                        nginx:latest
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card className="cursor-pointer hover:shadow-md transition-shadow">
                <CardContent className="p-4">
                  <div className="flex items-center space-x-3">
                    <div className="p-2 rounded-lg bg-green-500">
                      <Server className="h-5 w-5 text-white" />
                    </div>
                    <div>
                      <h3 className="font-medium">Database</h3>
                      <p className="text-sm text-muted-foreground">
                        postgres:13
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card className="cursor-pointer hover:shadow-md transition-shadow">
                <CardContent className="p-4">
                  <div className="flex items-center space-x-3">
                    <div className="p-2 rounded-lg bg-purple-500">
                      <Server className="h-5 w-5 text-white" />
                    </div>
                    <div>
                      <h3 className="font-medium">Cache</h3>
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