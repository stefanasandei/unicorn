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
  FileText,
  Plus,
  Upload,
  Download,
  Trash2,
  Folder,
  File,
  HardDrive,
  Cloud,
  Database,
  AlertCircle,
  Clock,
  Zap,
} from "lucide-react";
import { apiClient } from "@/lib/api";
import { StorageBucket, StorageFile } from "@/types/api";

export default function StoragePage() {
  const [buckets, setBuckets] = useState<StorageBucket[]>([]);
  const [selectedBucket, setSelectedBucket] = useState<StorageBucket | null>(
    null
  );
  const [files, setFiles] = useState<StorageFile[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  // Form states
  const [newBucketName, setNewBucketName] = useState("");
  const [uploadFile, setUploadFile] = useState<File | null>(null);

  useEffect(() => {
    fetchBuckets();
  }, []);

  const fetchBuckets = async () => {
    try {
      setIsLoading(true);
      const bucketsData = await apiClient.listBuckets();
      setBuckets(bucketsData);
    } catch (err: unknown) {
      const error = err as { response?: { data?: { error?: string } } };
      setError(error.response?.data?.error || "Failed to fetch buckets");
    } finally {
      setIsLoading(false);
    }
  };

  const fetchFiles = async (bucketId: string) => {
    try {
      const filesData = await apiClient.listFiles(bucketId);
      setFiles(filesData);
    } catch (err: unknown) {
      const error = err as { response?: { data?: { error?: string } } };
      setError(error.response?.data?.error || "Failed to fetch files");
    }
  };

  const handleCreateBucket = async () => {
    if (!newBucketName.trim()) {
      setError("Bucket name is required");
      return;
    }

    try {
      await apiClient.createBucket(newBucketName);
      setNewBucketName("");
      fetchBuckets();
    } catch (err: unknown) {
      const error = err as { response?: { data?: { error?: string } } };
      setError(error.response?.data?.error || "Failed to create bucket");
    }
  };

  const handleUploadFile = async () => {
    if (!uploadFile || !selectedBucket) {
      setError("Please select a file and bucket");
      return;
    }

    try {
      await apiClient.uploadFile(selectedBucket.id, uploadFile);
      setUploadFile(null);
      fetchFiles(selectedBucket.id);
    } catch (err: unknown) {
      const error = err as { response?: { data?: { error?: string } } };
      setError(error.response?.data?.error || "Failed to upload file");
    }
  };

  const handleDownloadFile = async (fileId: string) => {
    if (!selectedBucket) return;

    try {
      const blob = await apiClient.downloadFile(selectedBucket.id, fileId);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = files.find((f) => f.id === fileId)?.name || "download";
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch (err: unknown) {
      const error = err as { response?: { data?: { error?: string } } };
      setError(error.response?.data?.error || "Failed to download file");
    }
  };

  const handleDeleteFile = async (fileId: string) => {
    if (!selectedBucket) return;

    try {
      await apiClient.deleteFile(selectedBucket.id, fileId);
      fetchFiles(selectedBucket.id);
    } catch (err: unknown) {
      const error = err as { response?: { data?: { error?: string } } };
      setError(error.response?.data?.error || "Failed to delete file");
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
  };

  const getFileIcon = (contentType: string) => {
    if (contentType.startsWith("image/"))
      return <FileText className="h-4 w-4 text-blue-500" />;
    if (contentType.startsWith("video/"))
      return <FileText className="h-4 w-4 text-purple-500" />;
    if (contentType.startsWith("audio/"))
      return <FileText className="h-4 w-4 text-green-500" />;
    if (contentType.includes("pdf"))
      return <FileText className="h-4 w-4 text-red-500" />;
    return <FileText className="h-4 w-4 text-muted-foreground" />;
  };

  const totalStorage = files.reduce((sum, file) => sum + file.size, 0);

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
              <HardDrive className="h-6 w-6 text-primary" />
            </div>
            <div>
              <h1 className="text-3xl font-bold text-foreground">Storage</h1>
              <p className="text-muted-foreground">
                Manage your storage buckets and files
              </p>
            </div>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">
                Total Buckets
              </CardTitle>
              <div className="p-2 rounded-lg bg-blue-500/10">
                <Folder className="h-4 w-4 text-blue-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">
                {buckets.length}
              </div>
              <p className="text-xs text-muted-foreground">Storage buckets</p>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">
                Total Files
              </CardTitle>
              <div className="p-2 rounded-lg bg-green-500/10">
                <File className="h-4 w-4 text-green-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">
                {files.length}
              </div>
              <p className="text-xs text-muted-foreground">Stored files</p>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">
                Total Storage
              </CardTitle>
              <div className="p-2 rounded-lg bg-purple-500/10">
                <Cloud className="h-4 w-4 text-purple-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">
                {formatFileSize(totalStorage)}
              </div>
              <p className="text-xs text-muted-foreground">Used storage</p>
            </CardContent>
          </Card>
        </div>

        {error && (
          <div className="bg-destructive/10 border border-destructive/50 rounded-lg p-4 flex items-center space-x-2">
            <AlertCircle className="h-4 w-4 text-destructive" />
            <p className="text-destructive">{error}</p>
          </div>
        )}

        <Tabs defaultValue="buckets" className="space-y-6">
          <TabsList className="bg-card border border-border">
            <TabsTrigger
              value="buckets"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
            >
              <Folder className="h-4 w-4 mr-2" />
              Buckets
            </TabsTrigger>
            <TabsTrigger
              value="files"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
            >
              <File className="h-4 w-4 mr-2" />
              Files
            </TabsTrigger>
          </TabsList>

          <TabsContent value="buckets" className="space-y-6">
            <Card className="border-border/50">
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle className="text-foreground">
                      Storage Buckets
                    </CardTitle>
                    <CardDescription>
                      Create and manage storage buckets
                    </CardDescription>
                  </div>
                  <Dialog>
                    <DialogTrigger asChild>
                      <Button className="bg-primary hover:bg-primary/90">
                        <Plus className="h-4 w-4 mr-2" />
                        Create Bucket
                      </Button>
                    </DialogTrigger>
                    <DialogContent className="border-border/50">
                      <DialogHeader>
                        <DialogTitle className="text-foreground">
                          Create New Bucket
                        </DialogTitle>
                        <DialogDescription>
                          Create a new storage bucket
                        </DialogDescription>
                      </DialogHeader>
                      <div className="space-y-4">
                        <div>
                          <Label
                            htmlFor="bucket-name"
                            className="text-foreground"
                          >
                            Bucket Name
                          </Label>
                          <Input
                            id="bucket-name"
                            value={newBucketName}
                            onChange={(e) => setNewBucketName(e.target.value)}
                            placeholder="Enter bucket name"
                            className="border-border/50 focus:border-primary focus:ring-primary/20"
                          />
                        </div>
                      </div>
                      <DialogFooter>
                        <Button
                          onClick={handleCreateBucket}
                          className="bg-primary hover:bg-primary/90"
                        >
                          Create Bucket
                        </Button>
                      </DialogFooter>
                    </DialogContent>
                  </Dialog>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {buckets.map((bucket) => (
                    <div
                      key={bucket.id}
                      className="border border-border/50 rounded-lg p-4 hover:shadow-md transition-all duration-200"
                    >
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-3">
                          <div className="p-2 rounded-lg bg-accent/50">
                            <Folder className="h-4 w-4 text-accent-foreground" />
                          </div>
                          <div>
                            <h3 className="font-medium text-foreground">
                              {bucket.name}
                            </h3>
                            <p className="text-sm text-muted-foreground flex items-center gap-1">
                              <Clock className="h-3 w-3" />
                              Created{" "}
                              {new Date(bucket.created_at).toLocaleDateString()}
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center space-x-3">
                          <Badge variant="outline" className="bg-secondary/20">
                            {bucket.files?.length || 0} files
                          </Badge>
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => {
                              setSelectedBucket(bucket);
                              fetchFiles(bucket.id);
                            }}
                            className="border-border/50"
                          >
                            View Files
                          </Button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="files" className="space-y-6">
            {selectedBucket ? (
              <Card className="border-border/50">
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div>
                      <CardTitle className="text-foreground">
                        Files in {selectedBucket.name}
                      </CardTitle>
                      <CardDescription>
                        Manage files in this bucket
                      </CardDescription>
                    </div>
                    <Dialog>
                      <DialogTrigger asChild>
                        <Button className="bg-primary hover:bg-primary/90">
                          <Upload className="h-4 w-4 mr-2" />
                          Upload File
                        </Button>
                      </DialogTrigger>
                      <DialogContent className="border-border/50">
                        <DialogHeader>
                          <DialogTitle className="text-foreground">
                            Upload File
                          </DialogTitle>
                          <DialogDescription>
                            Upload a file to this bucket
                          </DialogDescription>
                        </DialogHeader>
                        <div className="space-y-4">
                          <div>
                            <Label
                              htmlFor="file-upload"
                              className="text-foreground"
                            >
                              Select File
                            </Label>
                            <Input
                              id="file-upload"
                              type="file"
                              onChange={(e) =>
                                setUploadFile(e.target.files?.[0] || null)
                              }
                              className="border-border/50 focus:border-primary focus:ring-primary/20"
                            />
                          </div>
                        </div>
                        <DialogFooter>
                          <Button
                            onClick={handleUploadFile}
                            className="bg-primary hover:bg-primary/90"
                          >
                            Upload File
                          </Button>
                        </DialogFooter>
                      </DialogContent>
                    </Dialog>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {files.map((file) => (
                      <div
                        key={file.id}
                        className="border border-border/50 rounded-lg p-4 hover:shadow-md transition-all duration-200"
                      >
                        <div className="flex items-center justify-between">
                          <div className="flex items-center space-x-3">
                            <div className="p-2 rounded-lg bg-accent/50">
                              {getFileIcon(file.content_type)}
                            </div>
                            <div>
                              <h3 className="font-medium text-foreground">
                                {file.name}
                              </h3>
                              <p className="text-sm text-muted-foreground flex items-center gap-1">
                                <Clock className="h-3 w-3" />
                                Uploaded{" "}
                                {new Date(file.created_at).toLocaleDateString()}
                              </p>
                            </div>
                          </div>
                          <div className="flex items-center space-x-3">
                            <Badge
                              variant="outline"
                              className="bg-secondary/20"
                            >
                              {formatFileSize(file.size)}
                            </Badge>
                            <Badge
                              variant="secondary"
                              className="bg-accent/50 text-accent-foreground"
                            >
                              {file.content_type}
                            </Badge>
                            <div className="flex space-x-1">
                              <Button
                                variant="outline"
                                size="sm"
                                onClick={() => handleDownloadFile(file.id)}
                                className="border-border/50"
                              >
                                <Download className="h-4 w-4" />
                              </Button>
                              <AlertDialog>
                                <AlertDialogTrigger asChild>
                                  <Button
                                    variant="outline"
                                    size="sm"
                                    className="border-border/50"
                                  >
                                    <Trash2 className="h-4 w-4" />
                                  </Button>
                                </AlertDialogTrigger>
                                <AlertDialogContent className="border-border/50">
                                  <AlertDialogHeader>
                                    <AlertDialogTitle className="text-foreground">
                                      Delete File
                                    </AlertDialogTitle>
                                    <AlertDialogDescription>
                                      Are you sure you want to delete this file?
                                      This action cannot be undone.
                                    </AlertDialogDescription>
                                  </AlertDialogHeader>
                                  <AlertDialogFooter>
                                    <AlertDialogCancel>
                                      Cancel
                                    </AlertDialogCancel>
                                    <AlertDialogAction
                                      onClick={() => handleDeleteFile(file.id)}
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
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            ) : (
              <Card className="border-border/50">
                <CardContent className="flex items-center justify-center h-32">
                  <div className="text-center space-y-2">
                    <Folder className="h-8 w-8 text-muted-foreground mx-auto" />
                    <p className="text-muted-foreground">
                      Select a bucket to view its files
                    </p>
                  </div>
                </CardContent>
              </Card>
            )}
          </TabsContent>
        </Tabs>
      </div>
    </Layout>
  );
}
