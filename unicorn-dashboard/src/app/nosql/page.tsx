"use client";

import { useState } from "react";
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
  Database,
  FolderOpen,
  FileText,
  Plus,
  Search,
  Filter,
  BarChart3,
  Settings,
  Play,
  Pause,
  Trash2,
  Eye,
  Edit,
  Copy,
  Download,
  Upload,
  Sparkles,
  Zap,
  Globe,
  Shield,
  Activity,
} from "lucide-react";

interface NoSQLDatabase {
  id: string;
  name: string;
  type: "mongodb" | "cassandra" | "dynamodb" | "redis";
  status: "running" | "stopped" | "error";
  collections: number;
  documents: number;
  size: string;
  createdAt: string;
  lastAccessed: string;
}

interface FolderOpen {
  id: string;
  databaseId: string;
  name: string;
  documents: number;
  size: string;
  indexes: number;
  createdAt: string;
}

interface FileText {
  id: string;
  collectionId: string;
  data: Record<string, any>;
  createdAt: string;
  updatedAt: string;
}

export default function NoSQLPage() {
  const [databases, setDatabases] = useState<NoSQLDatabase[]>([
    {
      id: "1",
      name: "user-management",
      type: "mongodb",
      status: "running",
      collections: 5,
      documents: 12450,
      size: "2.3 GB",
      createdAt: "2024-01-10T09:15:00Z",
      lastAccessed: "2024-01-20T14:22:00Z",
    },
    {
      id: "2",
      name: "analytics-data",
      type: "cassandra",
      status: "running",
      collections: 3,
      documents: 89234,
      size: "15.7 GB",
      createdAt: "2024-01-05T16:20:00Z",
      lastAccessed: "2024-01-20T13:45:00Z",
    },
    {
      id: "3",
      name: "session-store",
      type: "redis",
      status: "running",
      collections: 1,
      documents: 1567,
      size: "156 MB",
      createdAt: "2024-01-15T10:30:00Z",
      lastAccessed: "2024-01-20T12:30:00Z",
    },
  ]);

  const [collections, setCollections] = useState<FolderOpen[]>([
    {
      id: "1",
      databaseId: "1",
      name: "users",
      documents: 5432,
      size: "856 MB",
      indexes: 3,
      createdAt: "2024-01-10T09:15:00Z",
    },
    {
      id: "2",
      databaseId: "1",
      name: "profiles",
      documents: 5432,
      size: "1.2 GB",
      indexes: 2,
      createdAt: "2024-01-10T09:16:00Z",
    },
    {
      id: "3",
      databaseId: "2",
      name: "events",
      documents: 45678,
      size: "8.9 GB",
      indexes: 5,
      createdAt: "2024-01-05T16:20:00Z",
    },
  ]);

  const [documents, setDocuments] = useState<FileText[]>([
    {
      id: "1",
      collectionId: "1",
      data: {
        _id: "507f1f77bcf86cd799439011",
        email: "john.doe@example.com",
        name: "John Doe",
        age: 30,
        createdAt: "2024-01-15T10:30:00Z",
      },
      createdAt: "2024-01-15T10:30:00Z",
      updatedAt: "2024-01-20T14:22:00Z",
    },
    {
      id: "2",
      collectionId: "1",
      data: {
        _id: "507f1f77bcf86cd799439012",
        email: "jane.smith@example.com",
        name: "Jane Smith",
        age: 28,
        createdAt: "2024-01-16T11:45:00Z",
      },
      createdAt: "2024-01-16T11:45:00Z",
      updatedAt: "2024-01-19T09:15:00Z",
    },
    {
      id: "3",
      collectionId: "3",
      data: {
        _id: "507f1f77bcf86cd799439013",
        eventType: "page_view",
        userId: "507f1f77bcf86cd799439011",
        page: "/dashboard",
        timestamp: "2024-01-20T14:22:00Z",
      },
      createdAt: "2024-01-20T14:22:00Z",
      updatedAt: "2024-01-20T14:22:00Z",
    },
  ]);

  const [showCreateDatabase, setShowCreateDatabase] = useState(false);
  const [showCreateCollection, setShowCreateCollection] = useState(false);
  const [showAddDocument, setShowAddDocument] = useState(false);
  const [selectedDatabase, setSelectedDatabase] = useState<string | null>(null);
  const [selectedCollection, setSelectedCollection] = useState<string | null>(
    null
  );
  const [newDatabaseName, setNewDatabaseName] = useState("");
  const [newDatabaseType, setNewDatabaseType] = useState<
    "mongodb" | "cassandra" | "dynamodb" | "redis"
  >("mongodb");
  const [newCollectionName, setNewCollectionName] = useState("");
  const [newDocumentData, setNewDocumentData] = useState("");

  const totalDatabases = databases.length;
  const runningDatabases = databases.filter(
    (db) => db.status === "running"
  ).length;
  const totalCollections = collections.length;
  const totalDocuments = documents.length;

  const handleCreateDatabase = () => {
    if (newDatabaseName.trim()) {
      const newDatabase: NoSQLDatabase = {
        id: Date.now().toString(),
        name: newDatabaseName.trim(),
        type: newDatabaseType,
        status: "running",
        collections: 0,
        documents: 0,
        size: "0 MB",
        createdAt: new Date().toISOString(),
        lastAccessed: new Date().toISOString(),
      };
      setDatabases([...databases, newDatabase]);
      setNewDatabaseName("");
      setNewDatabaseType("mongodb");
      setShowCreateDatabase(false);
    }
  };

  const handleCreateCollection = () => {
    if (selectedDatabase && newCollectionName.trim()) {
      const newCollection: FolderOpen = {
        id: Date.now().toString(),
        databaseId: selectedDatabase,
        name: newCollectionName.trim(),
        documents: 0,
        size: "0 MB",
        indexes: 0,
        createdAt: new Date().toISOString(),
      };
      setCollections([...collections, newCollection]);

      // Update database collection count
      setDatabases(
        databases.map((db) =>
          db.id === selectedDatabase
            ? { ...db, collections: db.collections + 1 }
            : db
        )
      );

      setNewCollectionName("");
      setShowCreateCollection(false);
    }
  };

  const handleAddDocument = () => {
    if (selectedCollection && newDocumentData.trim()) {
      try {
        const parsedData = JSON.parse(newDocumentData);
        const newDocument: FileText = {
          id: Date.now().toString(),
          collectionId: selectedCollection,
          data: parsedData,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        };
        setDocuments([...documents, newDocument]);

        // Update collection document count
        setCollections(
          collections.map((col) =>
            col.id === selectedCollection
              ? { ...col, documents: col.documents + 1 }
              : col
          )
        );

        setNewDocumentData("");
        setShowAddDocument(false);
      } catch (error) {
        alert("Invalid JSON format");
      }
    }
  };

  const toggleDatabaseStatus = (databaseId: string) => {
    setDatabases(
      databases.map((db) =>
        db.id === databaseId
          ? { ...db, status: db.status === "running" ? "stopped" : "running" }
          : db
      )
    );
  };

  const getDatabaseIcon = (type: string) => {
    switch (type) {
      case "mongodb":
        return "🍃";
      case "cassandra":
        return "🌿";
      case "dynamodb":
        return "⚡";
      case "redis":
        return "🔴";
      default:
        return "🗄️";
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "running":
        return "bg-success text-success-foreground";
      case "stopped":
        return "bg-muted text-muted-foreground";
      case "error":
        return "bg-destructive text-destructive-foreground";
      default:
        return "bg-muted text-muted-foreground";
    }
  };

  const getTypeColor = (type: string) => {
    switch (type) {
      case "mongodb":
        return "bg-green-100 text-green-800";
      case "cassandra":
        return "bg-blue-100 text-blue-800";
      case "dynamodb":
        return "bg-yellow-100 text-yellow-800";
      case "redis":
        return "bg-red-100 text-red-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  return (
    <Layout>
      <div className="space-y-6">
        {/* Header */}
        <div className="space-y-2">
          <div className="flex items-center space-x-3">
            <div className="p-2 rounded-lg bg-gradient-to-br from-primary/10 to-primary/20">
              <Database className="h-6 w-6 text-primary" />
            </div>
            <div>
              <h1 className="text-3xl font-bold text-foreground">
                NoSQL Database
              </h1>
              <p className="text-muted-foreground">
                Manage NoSQL databases and collections with ease
              </p>
            </div>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card className="border-border/50">
            <CardContent className="p-6">
              <div className="flex items-center space-x-3">
                <div className="p-2 rounded-lg bg-primary/10">
                  <Database className="h-5 w-5 text-primary" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">
                    Total Databases
                  </p>
                  <p className="text-2xl font-bold text-foreground">
                    {totalDatabases}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card className="border-border/50">
            <CardContent className="p-6">
              <div className="flex items-center space-x-3">
                <div className="p-2 rounded-lg bg-success/10">
                  <Play className="h-5 w-5 text-success" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Running</p>
                  <p className="text-2xl font-bold text-foreground">
                    {runningDatabases}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card className="border-border/50">
            <CardContent className="p-6">
              <div className="flex items-center space-x-3">
                <div className="p-2 rounded-lg bg-accent/10">
                  <FolderOpen className="h-5 w-5 text-accent" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Collections</p>
                  <p className="text-2xl font-bold text-foreground">
                    {totalCollections}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card className="border-border/50">
            <CardContent className="p-6">
              <div className="flex items-center space-x-3">
                <div className="p-2 rounded-lg bg-info/10">
                  <FileText className="h-5 w-5 text-info" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Documents</p>
                  <p className="text-2xl font-bold text-foreground">
                    {totalDocuments}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Main Content */}
        <Tabs defaultValue="databases" className="space-y-6">
          <TabsList className="bg-card border-border">
            <TabsTrigger
              value="databases"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
            >
              Databases
            </TabsTrigger>
            <TabsTrigger
              value="collections"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
            >
              Collections
            </TabsTrigger>
            <TabsTrigger
              value="documents"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
            >
              Documents
            </TabsTrigger>
            <TabsTrigger
              value="wizard"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
            >
              Database Wizard
            </TabsTrigger>
          </TabsList>

          <TabsContent value="databases" className="space-y-4">
            <div className="flex justify-between items-center">
              <h2 className="text-xl font-semibold text-foreground">
                Database Management
              </h2>
              <Dialog
                open={showCreateDatabase}
                onOpenChange={setShowCreateDatabase}
              >
                <DialogTrigger asChild>
                  <Button className="bg-primary hover:bg-primary/90">
                    <Plus className="h-4 w-4 mr-2" />
                    Create Database
                  </Button>
                </DialogTrigger>
                <DialogContent className="border-border/50">
                  <DialogHeader>
                    <DialogTitle className="text-foreground">
                      Create New Database
                    </DialogTitle>
                    <DialogDescription className="text-muted-foreground">
                      Create a new NoSQL database with your preferred type.
                    </DialogDescription>
                  </DialogHeader>
                  <div className="space-y-4">
                    <div>
                      <Label
                        htmlFor="database-name"
                        className="text-foreground"
                      >
                        Database Name
                      </Label>
                      <Input
                        id="database-name"
                        value={newDatabaseName}
                        onChange={(e) => setNewDatabaseName(e.target.value)}
                        placeholder="Enter database name..."
                        className="border-border/50 focus:border-primary focus:ring-primary/20"
                      />
                    </div>
                    <div>
                      <Label
                        htmlFor="database-type"
                        className="text-foreground"
                      >
                        Database Type
                      </Label>
                      <select
                        id="database-type"
                        value={newDatabaseType}
                        onChange={(e) =>
                          setNewDatabaseType(e.target.value as any)
                        }
                        className="w-full p-2 border border-border/50 rounded-md bg-background text-foreground focus:border-primary focus:ring-primary/20"
                      >
                        <option value="mongodb">MongoDB</option>
                        <option value="cassandra">Cassandra</option>
                        <option value="dynamodb">DynamoDB</option>
                        <option value="redis">Redis</option>
                      </select>
                    </div>
                  </div>
                  <DialogFooter>
                    <Button
                      variant="outline"
                      onClick={() => setShowCreateDatabase(false)}
                    >
                      Cancel
                    </Button>
                    <Button
                      onClick={handleCreateDatabase}
                      className="bg-primary hover:bg-primary/90"
                    >
                      Create Database
                    </Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
            </div>

            <div className="grid gap-4">
              {databases.map((database) => (
                <Card
                  key={database.id}
                  className="border-border/50 hover:shadow-theme-md transition-shadow"
                >
                  <CardContent className="p-6">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-4">
                        <div className="text-2xl">
                          {getDatabaseIcon(database.type)}
                        </div>
                        <div>
                          <div className="flex items-center space-x-2">
                            <h3 className="font-semibold text-foreground">
                              {database.name}
                            </h3>
                            <Badge className={getTypeColor(database.type)}>
                              {database.type}
                            </Badge>
                            <Badge className={getStatusColor(database.status)}>
                              {database.status}
                            </Badge>
                          </div>
                          <p className="text-sm text-muted-foreground">
                            Created{" "}
                            {new Date(database.createdAt).toLocaleDateString()}
                          </p>
                        </div>
                      </div>

                      <div className="flex items-center space-x-6">
                        <div className="text-center">
                          <p className="text-sm text-muted-foreground">
                            Collections
                          </p>
                          <p className="font-semibold text-foreground">
                            {database.collections}
                          </p>
                        </div>
                        <div className="text-center">
                          <p className="text-sm text-muted-foreground">
                            Documents
                          </p>
                          <p className="font-semibold text-foreground">
                            {database.documents}
                          </p>
                        </div>
                        <div className="text-center">
                          <p className="text-sm text-muted-foreground">Size</p>
                          <p className="font-semibold text-foreground">
                            {database.size}
                          </p>
                        </div>

                        <div className="flex items-center space-x-2">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => toggleDatabaseStatus(database.id)}
                          >
                            {database.status === "running" ? (
                              <>
                                <Pause className="h-4 w-4 mr-2" />
                                Stop
                              </>
                            ) : (
                              <>
                                <Play className="h-4 w-4 mr-2" />
                                Start
                              </>
                            )}
                          </Button>

                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => {
                              setSelectedDatabase(database.id);
                              setShowCreateCollection(true);
                            }}
                          >
                            <FolderOpen className="h-4 w-4 mr-2" />
                            Add Collection
                          </Button>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </TabsContent>

          <TabsContent value="collections" className="space-y-4">
            <div className="flex justify-between items-center">
              <h2 className="text-xl font-semibold text-foreground">
                Collection Management
              </h2>
              <Button
                variant="outline"
                onClick={() => setShowCreateCollection(true)}
                className="border-border/50"
              >
                <FolderOpen className="h-4 w-4 mr-2" />
                Create Collection
              </Button>
            </div>

            <div className="grid gap-4">
              {collections.map((collection) => (
                <Card key={collection.id} className="border-border/50">
                  <CardContent className="p-6">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-4">
                        <div className="p-2 rounded-lg bg-accent/10">
                          <FolderOpen className="h-5 w-5 text-accent" />
                        </div>
                        <div>
                          <h3 className="font-semibold text-foreground">
                            {collection.name}
                          </h3>
                          <p className="text-sm text-muted-foreground">
                            Database:{" "}
                            {
                              databases.find(
                                (db) => db.id === collection.databaseId
                              )?.name
                            }
                          </p>
                        </div>
                      </div>

                      <div className="flex items-center space-x-6">
                        <div className="text-center">
                          <p className="text-sm text-muted-foreground">
                            Documents
                          </p>
                          <p className="font-semibold text-foreground">
                            {collection.documents}
                          </p>
                        </div>
                        <div className="text-center">
                          <p className="text-sm text-muted-foreground">Size</p>
                          <p className="font-semibold text-foreground">
                            {collection.size}
                          </p>
                        </div>
                        <div className="text-center">
                          <p className="text-sm text-muted-foreground">
                            Indexes
                          </p>
                          <p className="font-semibold text-foreground">
                            {collection.indexes}
                          </p>
                        </div>

                        <div className="flex items-center space-x-2">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => {
                              setSelectedCollection(collection.id);
                              setShowAddDocument(true);
                            }}
                          >
                            <FileText className="h-4 w-4 mr-2" />
                            Add Document
                          </Button>

                          <Button variant="outline" size="sm">
                            <Eye className="h-4 w-4 mr-2" />
                            View
                          </Button>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </TabsContent>

          <TabsContent value="documents" className="space-y-4">
            <div className="flex justify-between items-center">
              <h2 className="text-xl font-semibold text-foreground">
                Document Explorer
              </h2>
              <Button
                variant="outline"
                onClick={() => setShowAddDocument(true)}
                className="border-border/50"
              >
                <FileText className="h-4 w-4 mr-2" />
                Add Document
              </Button>
            </div>

            <div className="space-y-4">
              {documents.map((document) => (
                <Card key={document.id} className="border-border/50">
                  <CardContent className="p-4">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-4">
                        <div className="p-2 rounded-lg bg-info/10">
                          <FileText className="h-5 w-5 text-info" />
                        </div>
                        <div>
                          <p className="text-sm text-muted-foreground">
                            Collection:{" "}
                            {
                              collections.find(
                                (col) => col.id === document.collectionId
                              )?.name
                            }
                          </p>
                          <div className="mt-2 p-3 bg-muted/50 rounded-md">
                            <pre className="text-sm text-foreground font-mono overflow-x-auto">
                              {JSON.stringify(document.data, null, 2)}
                            </pre>
                          </div>
                        </div>
                      </div>

                      <div className="text-right">
                        <p className="text-sm text-muted-foreground">
                          Created:{" "}
                          {new Date(document.createdAt).toLocaleString()}
                        </p>
                        <p className="text-sm text-muted-foreground">
                          Updated:{" "}
                          {new Date(document.updatedAt).toLocaleString()}
                        </p>
                        <div className="flex items-center space-x-2 mt-2">
                          <Button variant="outline" size="sm">
                            <Edit className="h-4 w-4 mr-2" />
                            Edit
                          </Button>
                          <Button variant="outline" size="sm">
                            <Copy className="h-4 w-4 mr-2" />
                            Copy
                          </Button>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </TabsContent>

          <TabsContent value="wizard" className="space-y-4">
            <h2 className="text-xl font-semibold text-foreground">
              Database Wizard
            </h2>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Card className="border-border/50">
                <CardHeader>
                  <CardTitle className="text-foreground flex items-center space-x-2">
                    <Sparkles className="h-5 w-5 text-primary" />
                    Quick Setup
                  </CardTitle>
                  <CardDescription className="text-muted-foreground">
                    Create a database with recommended settings
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-3">
                    <Button className="w-full justify-start" variant="outline">
                      <Globe className="h-4 w-4 mr-2" />
                      Web Application Database
                    </Button>
                    <Button className="w-full justify-start" variant="outline">
                      <Activity className="h-4 w-4 mr-2" />
                      Analytics Database
                    </Button>
                    <Button className="w-full justify-start" variant="outline">
                      <Shield className="h-4 w-4 mr-2" />
                      Session Store
                    </Button>
                    <Button className="w-full justify-start" variant="outline">
                      <Zap className="h-4 w-4 mr-2" />
                      Cache Database
                    </Button>
                  </div>
                </CardContent>
              </Card>

              <Card className="border-border/50">
                <CardHeader>
                  <CardTitle className="text-foreground flex items-center space-x-2">
                    <Settings className="h-5 w-5 text-accent" />
                    Advanced Configuration
                  </CardTitle>
                  <CardDescription className="text-muted-foreground">
                    Customize database settings and performance
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-foreground">
                        Auto-scaling
                      </span>
                      <Button variant="outline" size="sm">
                        Configure
                      </Button>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-foreground">
                        Backup & Recovery
                      </span>
                      <Button variant="outline" size="sm">
                        Setup
                      </Button>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-foreground">
                        Security & Encryption
                      </span>
                      <Button variant="outline" size="sm">
                        Configure
                      </Button>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-foreground">
                        Monitoring & Alerts
                      </span>
                      <Button variant="outline" size="sm">
                        Setup
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>

            <Card className="border-border/50">
              <CardHeader>
                <CardTitle className="text-foreground">
                  Database Templates
                </CardTitle>
                <CardDescription className="text-muted-foreground">
                  Choose from pre-configured database templates
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div className="p-4 border border-border/50 rounded-lg hover:shadow-theme-md transition-shadow">
                    <div className="flex items-center space-x-2 mb-2">
                      <Globe className="h-4 w-4 text-primary" />
                      <h3 className="font-medium text-foreground">
                        E-commerce
                      </h3>
                    </div>
                    <p className="text-sm text-muted-foreground mb-3">
                      Product catalog, user management, and order processing
                    </p>
                    <Button size="sm" className="w-full">
                      Use Template
                    </Button>
                  </div>

                  <div className="p-4 border border-border/50 rounded-lg hover:shadow-theme-md transition-shadow">
                    <div className="flex items-center space-x-2 mb-2">
                      <Activity className="h-4 w-4 text-success" />
                      <h3 className="font-medium text-foreground">Analytics</h3>
                    </div>
                    <p className="text-sm text-muted-foreground mb-3">
                      Event tracking, metrics, and data warehousing
                    </p>
                    <Button size="sm" className="w-full">
                      Use Template
                    </Button>
                  </div>

                  <div className="p-4 border border-border/50 rounded-lg hover:shadow-theme-md transition-shadow">
                    <div className="flex items-center space-x-2 mb-2">
                      <Shield className="h-4 w-4 text-warning" />
                      <h3 className="font-medium text-foreground">
                        IoT Platform
                      </h3>
                    </div>
                    <p className="text-sm text-muted-foreground mb-3">
                      Device management, sensor data, and real-time processing
                    </p>
                    <Button size="sm" className="w-full">
                      Use Template
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        {/* Create Collection Dialog */}
        <Dialog
          open={showCreateCollection}
          onOpenChange={setShowCreateCollection}
        >
          <DialogContent className="border-border/50">
            <DialogHeader>
              <DialogTitle className="text-foreground">
                Create New Collection
              </DialogTitle>
              <DialogDescription className="text-muted-foreground">
                Create a new collection in the selected database.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label htmlFor="collection-name" className="text-foreground">
                  Collection Name
                </Label>
                <Input
                  id="collection-name"
                  value={newCollectionName}
                  onChange={(e) => setNewCollectionName(e.target.value)}
                  placeholder="Enter collection name..."
                  className="border-border/50 focus:border-primary focus:ring-primary/20"
                />
              </div>
            </div>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => setShowCreateCollection(false)}
              >
                Cancel
              </Button>
              <Button
                onClick={handleCreateCollection}
                disabled={!selectedDatabase || !newCollectionName.trim()}
                className="bg-primary hover:bg-primary/90"
              >
                Create Collection
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Add Document Dialog */}
        <Dialog open={showAddDocument} onOpenChange={setShowAddDocument}>
          <DialogContent className="border-border/50">
            <DialogHeader>
              <DialogTitle className="text-foreground">
                Add Document
              </DialogTitle>
              <DialogDescription className="text-muted-foreground">
                Add a new document to the selected collection (JSON format).
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label htmlFor="document-data" className="text-foreground">
                  Document Data (JSON)
                </Label>
                <textarea
                  id="document-data"
                  value={newDocumentData}
                  onChange={(e) => setNewDocumentData(e.target.value)}
                  placeholder='{"key": "value", "number": 123, "array": [1, 2, 3]}'
                  className="w-full h-32 p-3 border border-border/50 rounded-md font-mono text-sm bg-background text-foreground focus:border-primary focus:ring-primary/20"
                />
              </div>
            </div>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => setShowAddDocument(false)}
              >
                Cancel
              </Button>
              <Button
                onClick={handleAddDocument}
                disabled={!selectedCollection || !newDocumentData.trim()}
                className="bg-primary hover:bg-primary/90"
              >
                Add Document
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </Layout>
  );
}
