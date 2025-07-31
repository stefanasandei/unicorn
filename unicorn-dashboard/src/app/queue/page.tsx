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
  MessageSquare,
  Send,
  Download,
  Trash2,
  Plus,
  Clock,
  AlertCircle,
  CheckCircle,
  XCircle,
  BarChart3,
  Settings,
  Play,
  Pause,
  RotateCcw,
} from "lucide-react";

interface Queue {
  id: string;
  name: string;
  status: "active" | "paused" | "error";
  messageCount: number;
  processedCount: number;
  failedCount: number;
  createdAt: string;
  lastMessageAt: string;
}

interface Message {
  id: string;
  queueId: string;
  content: string;
  status: "pending" | "processing" | "completed" | "failed";
  priority: "low" | "normal" | "high";
  createdAt: string;
  processedAt?: string;
  retryCount: number;
}

export default function QueuePage() {
  const [queues, setQueues] = useState<Queue[]>([
    {
      id: "1",
      name: "email-notifications",
      status: "active",
      messageCount: 156,
      processedCount: 1243,
      failedCount: 12,
      createdAt: "2024-01-15T10:30:00Z",
      lastMessageAt: "2024-01-20T14:22:00Z",
    },
    {
      id: "2",
      name: "image-processing",
      status: "active",
      messageCount: 89,
      processedCount: 567,
      failedCount: 3,
      createdAt: "2024-01-10T09:15:00Z",
      lastMessageAt: "2024-01-20T13:45:00Z",
    },
    {
      id: "3",
      name: "data-sync",
      status: "paused",
      messageCount: 0,
      processedCount: 234,
      failedCount: 0,
      createdAt: "2024-01-05T16:20:00Z",
      lastMessageAt: "2024-01-19T11:30:00Z",
    },
  ]);

  const [messages, setMessages] = useState<Message[]>([
    {
      id: "1",
      queueId: "1",
      content:
        '{"type": "welcome", "email": "user@example.com", "template": "welcome-v1"}',
      status: "completed",
      priority: "normal",
      createdAt: "2024-01-20T14:22:00Z",
      processedAt: "2024-01-20T14:22:15Z",
      retryCount: 0,
    },
    {
      id: "2",
      queueId: "1",
      content:
        '{"type": "password-reset", "email": "admin@example.com", "token": "abc123"}',
      status: "pending",
      priority: "high",
      createdAt: "2024-01-20T14:21:00Z",
      retryCount: 0,
    },
    {
      id: "3",
      queueId: "2",
      content:
        '{"imageId": "img_123", "operations": ["resize", "compress"], "format": "webp"}',
      status: "processing",
      priority: "normal",
      createdAt: "2024-01-20T14:20:00Z",
      retryCount: 1,
    },
  ]);

  const [showCreateQueue, setShowCreateQueue] = useState(false);
  const [showSendMessage, setShowSendMessage] = useState(false);
  const [selectedQueue, setSelectedQueue] = useState<string | null>(null);
  const [newQueueName, setNewQueueName] = useState("");
  const [newMessageContent, setNewMessageContent] = useState("");
  const [newMessagePriority, setNewMessagePriority] = useState<
    "low" | "normal" | "high"
  >("normal");

  const totalQueues = queues.length;
  const activeQueues = queues.filter((q) => q.status === "active").length;
  const totalMessages = queues.reduce((sum, q) => sum + q.messageCount, 0);
  const totalProcessed = queues.reduce((sum, q) => sum + q.processedCount, 0);

  const handleCreateQueue = () => {
    if (newQueueName.trim()) {
      const newQueue: Queue = {
        id: Date.now().toString(),
        name: newQueueName.trim(),
        status: "active",
        messageCount: 0,
        processedCount: 0,
        failedCount: 0,
        createdAt: new Date().toISOString(),
        lastMessageAt: new Date().toISOString(),
      };
      setQueues([...queues, newQueue]);
      setNewQueueName("");
      setShowCreateQueue(false);
    }
  };

  const handleSendMessage = () => {
    if (selectedQueue && newMessageContent.trim()) {
      const newMessage: Message = {
        id: Date.now().toString(),
        queueId: selectedQueue,
        content: newMessageContent.trim(),
        status: "pending",
        priority: newMessagePriority,
        createdAt: new Date().toISOString(),
        retryCount: 0,
      };
      setMessages([...messages, newMessage]);

      // Update queue message count
      setQueues(
        queues.map((q) =>
          q.id === selectedQueue
            ? {
                ...q,
                messageCount: q.messageCount + 1,
                lastMessageAt: new Date().toISOString(),
              }
            : q
        )
      );

      setNewMessageContent("");
      setNewMessagePriority("normal");
      setShowSendMessage(false);
    }
  };

  const toggleQueueStatus = (queueId: string) => {
    setQueues(
      queues.map((q) =>
        q.id === queueId
          ? { ...q, status: q.status === "active" ? "paused" : "active" }
          : q
      )
    );
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case "active":
        return <CheckCircle className="h-4 w-4 text-success" />;
      case "paused":
        return <Pause className="h-4 w-4 text-warning" />;
      case "error":
        return <XCircle className="h-4 w-4 text-destructive" />;
      default:
        return <AlertCircle className="h-4 w-4 text-muted-foreground" />;
    }
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case "high":
        return "bg-destructive text-destructive-foreground";
      case "normal":
        return "bg-primary text-primary-foreground";
      case "low":
        return "bg-muted text-muted-foreground";
      default:
        return "bg-muted text-muted-foreground";
    }
  };

  const getMessageStatusColor = (status: string) => {
    switch (status) {
      case "completed":
        return "bg-success text-success-foreground";
      case "processing":
        return "bg-primary text-primary-foreground";
      case "pending":
        return "bg-muted text-muted-foreground";
      case "failed":
        return "bg-destructive text-destructive-foreground";
      default:
        return "bg-muted text-muted-foreground";
    }
  };

  return (
    <Layout>
      <div className="space-y-6">
        {/* Header */}
        <div className="space-y-2">
          <div className="flex items-center space-x-3">
            <div className="p-2 rounded-lg bg-gradient-to-br from-primary/10 to-primary/20">
              <MessageSquare className="h-6 w-6 text-primary" />
            </div>
            <div>
              <h1 className="text-3xl font-bold text-foreground">
                Queue Service
              </h1>
              <p className="text-muted-foreground">
                Manage message queues and monitor processing
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
                  <MessageSquare className="h-5 w-5 text-primary" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Total Queues</p>
                  <p className="text-2xl font-bold text-foreground">
                    {totalQueues}
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
                  <p className="text-sm text-muted-foreground">Active Queues</p>
                  <p className="text-2xl font-bold text-foreground">
                    {activeQueues}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card className="border-border/50">
            <CardContent className="p-6">
              <div className="flex items-center space-x-3">
                <div className="p-2 rounded-lg bg-accent/10">
                  <Send className="h-5 w-5 text-accent" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">
                    Pending Messages
                  </p>
                  <p className="text-2xl font-bold text-foreground">
                    {totalMessages}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card className="border-border/50">
            <CardContent className="p-6">
              <div className="flex items-center space-x-3">
                <div className="p-2 rounded-lg bg-info/10">
                  <BarChart3 className="h-5 w-5 text-info" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Processed</p>
                  <p className="text-2xl font-bold text-foreground">
                    {totalProcessed}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Main Content */}
        <Tabs defaultValue="queues" className="space-y-6">
          <TabsList className="bg-card border-border">
            <TabsTrigger
              value="queues"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
            >
              Queues
            </TabsTrigger>
            <TabsTrigger
              value="messages"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
            >
              Messages
            </TabsTrigger>
            <TabsTrigger
              value="monitoring"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
            >
              Monitoring
            </TabsTrigger>
          </TabsList>

          <TabsContent value="queues" className="space-y-4">
            <div className="flex justify-between items-center">
              <h2 className="text-xl font-semibold text-foreground">
                Queue Management
              </h2>
              <Dialog open={showCreateQueue} onOpenChange={setShowCreateQueue}>
                <DialogTrigger asChild>
                  <Button className="bg-primary hover:bg-primary/90">
                    <Plus className="h-4 w-4 mr-2" />
                    Create Queue
                  </Button>
                </DialogTrigger>
                <DialogContent className="border-border/50">
                  <DialogHeader>
                    <DialogTitle className="text-foreground">
                      Create New Queue
                    </DialogTitle>
                    <DialogDescription className="text-muted-foreground">
                      Create a new message queue for processing tasks.
                    </DialogDescription>
                  </DialogHeader>
                  <div className="space-y-4">
                    <div>
                      <Label htmlFor="queue-name" className="text-foreground">
                        Queue Name
                      </Label>
                      <Input
                        id="queue-name"
                        value={newQueueName}
                        onChange={(e) => setNewQueueName(e.target.value)}
                        placeholder="Enter queue name..."
                        className="border-border/50 focus:border-primary focus:ring-primary/20"
                      />
                    </div>
                  </div>
                  <DialogFooter>
                    <Button
                      variant="outline"
                      onClick={() => setShowCreateQueue(false)}
                    >
                      Cancel
                    </Button>
                    <Button
                      onClick={handleCreateQueue}
                      className="bg-primary hover:bg-primary/90"
                    >
                      Create Queue
                    </Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
            </div>

            <div className="grid gap-4">
              {queues.map((queue) => (
                <Card
                  key={queue.id}
                  className="border-border/50 hover:shadow-theme-md transition-shadow"
                >
                  <CardContent className="p-6">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-4">
                        <div className="flex items-center space-x-2">
                          {getStatusIcon(queue.status)}
                          <div>
                            <h3 className="font-semibold text-foreground">
                              {queue.name}
                            </h3>
                            <p className="text-sm text-muted-foreground">
                              Created{" "}
                              {new Date(queue.createdAt).toLocaleDateString()}
                            </p>
                          </div>
                        </div>
                      </div>

                      <div className="flex items-center space-x-6">
                        <div className="text-center">
                          <p className="text-sm text-muted-foreground">
                            Pending
                          </p>
                          <p className="font-semibold text-foreground">
                            {queue.messageCount}
                          </p>
                        </div>
                        <div className="text-center">
                          <p className="text-sm text-muted-foreground">
                            Processed
                          </p>
                          <p className="font-semibold text-foreground">
                            {queue.processedCount}
                          </p>
                        </div>
                        <div className="text-center">
                          <p className="text-sm text-muted-foreground">
                            Failed
                          </p>
                          <p className="font-semibold text-foreground">
                            {queue.failedCount}
                          </p>
                        </div>

                        <div className="flex items-center space-x-2">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => toggleQueueStatus(queue.id)}
                          >
                            {queue.status === "active" ? (
                              <>
                                <Pause className="h-4 w-4 mr-2" />
                                Pause
                              </>
                            ) : (
                              <>
                                <Play className="h-4 w-4 mr-2" />
                                Resume
                              </>
                            )}
                          </Button>

                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => {
                              setSelectedQueue(queue.id);
                              setShowSendMessage(true);
                            }}
                          >
                            <Send className="h-4 w-4 mr-2" />
                            Send Message
                          </Button>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </TabsContent>

          <TabsContent value="messages" className="space-y-4">
            <div className="flex justify-between items-center">
              <h2 className="text-xl font-semibold text-foreground">
                Message History
              </h2>
              <Button
                variant="outline"
                onClick={() => setShowSendMessage(true)}
                className="border-border/50"
              >
                <Send className="h-4 w-4 mr-2" />
                Send Message
              </Button>
            </div>

            <div className="space-y-4">
              {messages.map((message) => (
                <Card key={message.id} className="border-border/50">
                  <CardContent className="p-4">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-4">
                        <Badge
                          className={getMessageStatusColor(message.status)}
                        >
                          {message.status}
                        </Badge>
                        <Badge className={getPriorityColor(message.priority)}>
                          {message.priority}
                        </Badge>
                        <div>
                          <p className="text-sm text-muted-foreground">
                            Queue:{" "}
                            {queues.find((q) => q.id === message.queueId)?.name}
                          </p>
                          <p className="text-sm text-foreground font-mono">
                            {message.content.length > 100
                              ? `${message.content.substring(0, 100)}...`
                              : message.content}
                          </p>
                        </div>
                      </div>

                      <div className="text-right">
                        <p className="text-sm text-muted-foreground">
                          {new Date(message.createdAt).toLocaleString()}
                        </p>
                        {message.retryCount > 0 && (
                          <p className="text-xs text-warning">
                            Retries: {message.retryCount}
                          </p>
                        )}
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </TabsContent>

          <TabsContent value="monitoring" className="space-y-4">
            <h2 className="text-xl font-semibold text-foreground">
              Queue Monitoring
            </h2>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Card className="border-border/50">
                <CardHeader>
                  <CardTitle className="text-foreground">
                    Queue Performance
                  </CardTitle>
                  <CardDescription className="text-muted-foreground">
                    Processing metrics for all queues
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  {queues.map((queue) => (
                    <div key={queue.id} className="space-y-2">
                      <div className="flex justify-between items-center">
                        <span className="text-sm font-medium text-foreground">
                          {queue.name}
                        </span>
                        <span className="text-sm text-muted-foreground">
                          {queue.processedCount} processed
                        </span>
                      </div>
                      <div className="w-full bg-muted rounded-full h-2">
                        <div
                          className="bg-primary h-2 rounded-full transition-all"
                          style={{
                            width: `${
                              queue.processedCount > 0
                                ? Math.min(
                                    100,
                                    (queue.processedCount /
                                      (queue.processedCount +
                                        queue.messageCount)) *
                                      100
                                  )
                                : 0
                            }%`,
                          }}
                        />
                      </div>
                    </div>
                  ))}
                </CardContent>
              </Card>

              <Card className="border-border/50">
                <CardHeader>
                  <CardTitle className="text-foreground">
                    Recent Activity
                  </CardTitle>
                  <CardDescription className="text-muted-foreground">
                    Latest queue operations
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-3">
                  {messages.slice(0, 5).map((message) => (
                    <div
                      key={message.id}
                      className="flex items-center space-x-3"
                    >
                      <div
                        className={`w-2 h-2 rounded-full ${
                          message.status === "completed"
                            ? "bg-success"
                            : message.status === "processing"
                            ? "bg-primary"
                            : message.status === "failed"
                            ? "bg-destructive"
                            : "bg-muted"
                        }`}
                      />
                      <div className="flex-1">
                        <p className="text-sm text-foreground">
                          Message {message.id.substring(0, 8)}...
                        </p>
                        <p className="text-xs text-muted-foreground">
                          {queues.find((q) => q.id === message.queueId)?.name} •{" "}
                          {message.status}
                        </p>
                      </div>
                      <span className="text-xs text-muted-foreground">
                        {new Date(message.createdAt).toLocaleTimeString()}
                      </span>
                    </div>
                  ))}
                </CardContent>
              </Card>
            </div>
          </TabsContent>
        </Tabs>

        {/* Send Message Dialog */}
        <Dialog open={showSendMessage} onOpenChange={setShowSendMessage}>
          <DialogContent className="border-border/50">
            <DialogHeader>
              <DialogTitle className="text-foreground">
                Send Message
              </DialogTitle>
              <DialogDescription className="text-muted-foreground">
                Send a new message to a queue.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label htmlFor="queue-select" className="text-foreground">
                  Select Queue
                </Label>
                <select
                  id="queue-select"
                  value={selectedQueue || ""}
                  onChange={(e) => setSelectedQueue(e.target.value)}
                  className="w-full p-2 border border-border/50 rounded-md bg-background text-foreground focus:border-primary focus:ring-primary/20"
                >
                  <option value="">Select a queue...</option>
                  {queues.map((queue) => (
                    <option key={queue.id} value={queue.id}>
                      {queue.name}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <Label htmlFor="message-priority" className="text-foreground">
                  Priority
                </Label>
                <select
                  id="message-priority"
                  value={newMessagePriority}
                  onChange={(e) =>
                    setNewMessagePriority(
                      e.target.value as "low" | "normal" | "high"
                    )
                  }
                  className="w-full p-2 border border-border/50 rounded-md bg-background text-foreground focus:border-primary focus:ring-primary/20"
                >
                  <option value="low">Low</option>
                  <option value="normal">Normal</option>
                  <option value="high">High</option>
                </select>
              </div>

              <div>
                <Label htmlFor="message-content" className="text-foreground">
                  Message Content
                </Label>
                <textarea
                  id="message-content"
                  value={newMessageContent}
                  onChange={(e) => setNewMessageContent(e.target.value)}
                  placeholder="Enter message content (JSON recommended)..."
                  className="w-full h-32 p-3 border border-border/50 rounded-md font-mono text-sm bg-background text-foreground focus:border-primary focus:ring-primary/20"
                />
              </div>
            </div>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => setShowSendMessage(false)}
              >
                Cancel
              </Button>
              <Button
                onClick={handleSendMessage}
                disabled={!selectedQueue || !newMessageContent.trim()}
                className="bg-primary hover:bg-primary/90"
              >
                <Send className="h-4 w-4 mr-2" />
                Send Message
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </Layout>
  );
}
