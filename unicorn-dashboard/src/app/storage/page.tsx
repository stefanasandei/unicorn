'use client';

import React, { useEffect, useState } from 'react';
import { Layout } from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
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
import { FileText, Plus, Upload, Download, Trash2, Folder, File } from 'lucide-react';
import { apiClient } from '@/lib/api';
import { StorageBucket, File as FileType } from '@/types/api';

export default function StoragePage() {
  const [buckets, setBuckets] = useState<StorageBucket[]>([]);
  const [selectedBucket, setSelectedBucket] = useState<StorageBucket | null>(null);
  const [files, setFiles] = useState<FileType[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  // Form states
  const [newBucketName, setNewBucketName] = useState('');
  const [uploadFile, setUploadFile] = useState<File | null>(null);

  useEffect(() => {
    fetchBuckets();
  }, []);

  const fetchBuckets = async () => {
    try {
      setIsLoading(true);
      const bucketsData = await apiClient.listBuckets();
      setBuckets(bucketsData);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch buckets');
    } finally {
      setIsLoading(false);
    }
  };

  const fetchFiles = async (bucketId: string) => {
    try {
      const filesData = await apiClient.listFiles(bucketId);
      setFiles(filesData);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch files');
    }
  };

  const handleCreateBucket = async () => {
    if (!newBucketName.trim()) {
      setError('Bucket name is required');
      return;
    }

    try {
      await apiClient.createBucket(newBucketName);
      setNewBucketName('');
      fetchBuckets();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create bucket');
    }
  };

  const handleUploadFile = async () => {
    if (!uploadFile || !selectedBucket) {
      setError('Please select a file and bucket');
      return;
    }

    try {
      await apiClient.uploadFile(selectedBucket.id, uploadFile);
      setUploadFile(null);
      fetchFiles(selectedBucket.id);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to upload file');
    }
  };

  const handleDownloadFile = async (fileId: string) => {
    if (!selectedBucket) return;

    try {
      const blob = await apiClient.downloadFile(selectedBucket.id, fileId);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = files.find(f => f.id === fileId)?.name || 'download';
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to download file');
    }
  };

  const handleDeleteFile = async (fileId: string) => {
    if (!selectedBucket) return;

    try {
      await apiClient.deleteFile(selectedBucket.id, fileId);
      fetchFiles(selectedBucket.id);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to delete file');
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
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
          <h1 className="text-2xl font-bold text-gray-900">Storage</h1>
          <p className="text-gray-600">
            Manage your storage buckets and files
          </p>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 rounded-md p-4">
            <p className="text-red-800">{error}</p>
          </div>
        )}

        <Tabs defaultValue="buckets" className="space-y-4">
          <TabsList>
            <TabsTrigger value="buckets" className="flex items-center gap-2">
              <Folder className="h-4 w-4" />
              Buckets
            </TabsTrigger>
            <TabsTrigger value="files" className="flex items-center gap-2">
              <File className="h-4 w-4" />
              Files
            </TabsTrigger>
          </TabsList>

          <TabsContent value="buckets" className="space-y-4">
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle>Storage Buckets</CardTitle>
                    <CardDescription>
                      Create and manage storage buckets
                    </CardDescription>
                  </div>
                  <Dialog>
                    <DialogTrigger asChild>
                      <Button>
                        <Plus className="h-4 w-4 mr-2" />
                        Create Bucket
                      </Button>
                    </DialogTrigger>
                    <DialogContent>
                      <DialogHeader>
                        <DialogTitle>Create New Bucket</DialogTitle>
                        <DialogDescription>
                          Create a new storage bucket
                        </DialogDescription>
                      </DialogHeader>
                      <div className="space-y-4">
                        <div>
                          <Label htmlFor="bucket-name">Bucket Name</Label>
                          <Input
                            id="bucket-name"
                            value={newBucketName}
                            onChange={(e) => setNewBucketName(e.target.value)}
                            placeholder="Enter bucket name"
                          />
                        </div>
                      </div>
                      <DialogFooter>
                        <Button onClick={handleCreateBucket}>Create Bucket</Button>
                      </DialogFooter>
                    </DialogContent>
                  </Dialog>
                </div>
              </CardHeader>
              <CardContent>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Bucket Name</TableHead>
                      <TableHead>Files</TableHead>
                      <TableHead>Created</TableHead>
                      <TableHead>Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {buckets.map((bucket) => (
                      <TableRow key={bucket.id}>
                        <TableCell className="font-medium">{bucket.name}</TableCell>
                        <TableCell>
                          <Badge variant="secondary">
                            {bucket.files?.length || 0} files
                          </Badge>
                        </TableCell>
                        <TableCell>
                          {new Date(bucket.created_at).toLocaleDateString()}
                        </TableCell>
                        <TableCell>
                          <div className="flex space-x-2">
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => {
                                setSelectedBucket(bucket);
                                fetchFiles(bucket.id);
                              }}
                            >
                              View Files
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

          <TabsContent value="files" className="space-y-4">
            {selectedBucket ? (
              <Card>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div>
                      <CardTitle>Files in {selectedBucket.name}</CardTitle>
                      <CardDescription>
                        Manage files in this bucket
                      </CardDescription>
                    </div>
                    <Dialog>
                      <DialogTrigger asChild>
                        <Button>
                          <Upload className="h-4 w-4 mr-2" />
                          Upload File
                        </Button>
                      </DialogTrigger>
                      <DialogContent>
                        <DialogHeader>
                          <DialogTitle>Upload File</DialogTitle>
                          <DialogDescription>
                            Upload a file to this bucket
                          </DialogDescription>
                        </DialogHeader>
                        <div className="space-y-4">
                          <div>
                            <Label htmlFor="file-upload">Select File</Label>
                            <Input
                              id="file-upload"
                              type="file"
                              onChange={(e) => setUploadFile(e.target.files?.[0] || null)}
                            />
                          </div>
                        </div>
                        <DialogFooter>
                          <Button onClick={handleUploadFile}>Upload File</Button>
                        </DialogFooter>
                      </DialogContent>
                    </Dialog>
                  </div>
                </CardHeader>
                <CardContent>
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>File Name</TableHead>
                        <TableHead>Size</TableHead>
                        <TableHead>Type</TableHead>
                        <TableHead>Uploaded</TableHead>
                        <TableHead>Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {files.map((file) => (
                        <TableRow key={file.id}>
                          <TableCell className="font-medium">{file.name}</TableCell>
                          <TableCell>{formatFileSize(file.size)}</TableCell>
                          <TableCell>
                            <Badge variant="secondary">{file.content_type}</Badge>
                          </TableCell>
                          <TableCell>
                            {new Date(file.created_at).toLocaleDateString()}
                          </TableCell>
                          <TableCell>
                            <div className="flex space-x-2">
                              <Button
                                variant="outline"
                                size="sm"
                                onClick={() => handleDownloadFile(file.id)}
                              >
                                <Download className="h-4 w-4" />
                              </Button>
                              <AlertDialog>
                                <AlertDialogTrigger asChild>
                                  <Button variant="outline" size="sm">
                                    <Trash2 className="h-4 w-4" />
                                  </Button>
                                </AlertDialogTrigger>
                                <AlertDialogContent>
                                  <AlertDialogHeader>
                                    <AlertDialogTitle>Delete File</AlertDialogTitle>
                                    <AlertDialogDescription>
                                      Are you sure you want to delete this file? This action cannot be undone.
                                    </AlertDialogDescription>
                                  </AlertDialogHeader>
                                  <AlertDialogFooter>
                                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                                    <AlertDialogAction
                                      onClick={() => handleDeleteFile(file.id)}
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
            ) : (
              <Card>
                <CardContent className="flex items-center justify-center h-32">
                  <p className="text-gray-500">Select a bucket to view its files</p>
                </CardContent>
              </Card>
            )}
          </TabsContent>
        </Tabs>
      </div>
    </Layout>
  );
} 