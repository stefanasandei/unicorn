'use client';

import React, { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { Layout } from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { 
  Shield, 
  Database, 
  FileText, 
  Server, 
  Code, 
  Users, 
  Building,
  Activity,
  TrendingUp,
  Sparkles,
  ArrowUpRight,
  Clock
} from 'lucide-react';
import Link from 'next/link';

interface DashboardStats {
  secrets: number;
  buckets: number;
  files: number;
  containers: number;
  users: number;
}

export default function DashboardPage() {
  const { user } = useAuth();
  const [stats, setStats] = useState<DashboardStats>({
    secrets: 0,
    buckets: 0,
    files: 0,
    containers: 0,
    users: 0,
  });
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        // In a real implementation, you'd fetch these from the API
        // For now, we'll use mock data
        setStats({
          secrets: 12,
          buckets: 5,
          files: 24,
          containers: 3,
          users: 8,
        });
      } catch (error) {
        console.error('Failed to fetch stats:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchStats();
  }, []);

  const quickActions = [
    {
      title: 'Create Secret',
      description: 'Store encrypted secrets',
      icon: Database,
      href: '/secrets',
      gradient: 'from-blue-500 to-blue-600',
      bgGradient: 'from-blue-500/10 to-blue-600/10',
    },
    {
      title: 'Upload File',
      description: 'Upload files to storage',
      icon: FileText,
      href: '/storage',
      gradient: 'from-green-500 to-green-600',
      bgGradient: 'from-green-500/10 to-green-600/10',
    },
    {
      title: 'Deploy Container',
      description: 'Deploy compute containers',
      icon: Server,
      href: '/compute',
      gradient: 'from-purple-500 to-purple-600',
      bgGradient: 'from-purple-500/10 to-purple-600/10',
    },
    {
      title: 'Run Lambda',
      description: 'Execute serverless functions',
      icon: Code,
      href: '/lambda',
      gradient: 'from-orange-500 to-orange-600',
      bgGradient: 'from-orange-500/10 to-orange-600/10',
    },
  ];

  const recentActivity = [
    {
      id: 1,
      action: 'Created secret',
      resource: 'database-password',
      time: '2 minutes ago',
      type: 'secret',
      icon: Database,
      color: 'text-blue-500',
    },
    {
      id: 2,
      action: 'Uploaded file',
      resource: 'config.json',
      time: '5 minutes ago',
      type: 'file',
      icon: FileText,
      color: 'text-green-500',
    },
    {
      id: 3,
      action: 'Deployed container',
      resource: 'web-app',
      time: '10 minutes ago',
      type: 'container',
      icon: Server,
      color: 'text-purple-500',
    },
    {
      id: 4,
      action: 'Executed lambda',
      resource: 'data-processor',
      time: '15 minutes ago',
      type: 'lambda',
      icon: Code,
      color: 'text-orange-500',
    },
  ];

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
              <Sparkles className="h-6 w-6 text-primary" />
            </div>
            <div>
              <h1 className="text-3xl font-bold text-foreground">Dashboard</h1>
              <p className="text-muted-foreground">
                Welcome back, {user?.name}. Here's what's happening with your resources.
              </p>
            </div>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">Secrets</CardTitle>
              <div className="p-2 rounded-lg bg-blue-500/10">
                <Database className="h-4 w-4 text-blue-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">{stats.secrets}</div>
              <p className="text-xs text-muted-foreground flex items-center gap-1">
                <TrendingUp className="h-3 w-3 text-green-500" />
                +2 from last week
              </p>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">Storage Buckets</CardTitle>
              <div className="p-2 rounded-lg bg-green-500/10">
                <FileText className="h-4 w-4 text-green-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">{stats.buckets}</div>
              <p className="text-xs text-muted-foreground flex items-center gap-1">
                <TrendingUp className="h-3 w-3 text-green-500" />
                +1 from last week
              </p>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">Files</CardTitle>
              <div className="p-2 rounded-lg bg-purple-500/10">
                <FileText className="h-4 w-4 text-purple-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">{stats.files}</div>
              <p className="text-xs text-muted-foreground flex items-center gap-1">
                <TrendingUp className="h-3 w-3 text-green-500" />
                +5 from last week
              </p>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">Containers</CardTitle>
              <div className="p-2 rounded-lg bg-orange-500/10">
                <Server className="h-4 w-4 text-orange-500" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">{stats.containers}</div>
              <p className="text-xs text-muted-foreground flex items-center gap-1">
                <TrendingUp className="h-3 w-3 text-green-500" />
                +1 from last week
              </p>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-all duration-200 border-border/50">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-foreground">Users</CardTitle>
              <div className="p-2 rounded-lg bg-primary/10">
                <Users className="h-4 w-4 text-primary" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-foreground">{stats.users}</div>
              <p className="text-xs text-muted-foreground flex items-center gap-1">
                <TrendingUp className="h-3 w-3 text-green-500" />
                +2 from last week
              </p>
            </CardContent>
          </Card>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Quick Actions */}
          <Card className="border-border/50">
            <CardHeader>
              <CardTitle className="text-foreground">Quick Actions</CardTitle>
              <CardDescription>
                Common tasks to get you started
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                {quickActions.map((action) => (
                  <Link key={action.title} href={action.href}>
                    <Card className="hover:shadow-lg transition-all duration-200 cursor-pointer border-border/50 group">
                      <CardContent className="p-4">
                        <div className="flex items-center space-x-3">
                          <div className={`p-2.5 rounded-lg bg-gradient-to-br ${action.bgGradient} group-hover:scale-105 transition-transform`}>
                            <action.icon className={`h-5 w-5 bg-gradient-to-br ${action.gradient} bg-clip-text text-transparent`} />
                          </div>
                          <div className="flex-1">
                            <h3 className="font-medium text-foreground group-hover:text-primary transition-colors">{action.title}</h3>
                            <p className="text-sm text-muted-foreground">
                              {action.description}
                            </p>
                          </div>
                          <ArrowUpRight className="h-4 w-4 text-muted-foreground group-hover:text-primary transition-colors" />
                        </div>
                      </CardContent>
                    </Card>
                  </Link>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Recent Activity */}
          <Card className="border-border/50">
            <CardHeader>
              <CardTitle className="text-foreground">Recent Activity</CardTitle>
              <CardDescription>
                Latest actions across your resources
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {recentActivity.map((activity) => (
                  <div key={activity.id} className="flex items-center space-x-4 p-3 rounded-lg hover:bg-accent/30 transition-colors">
                    <div className="flex-shrink-0">
                      <div className={`p-2 rounded-lg bg-accent/50`}>
                        <activity.icon className={`h-4 w-4 ${activity.color}`} />
                      </div>
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-foreground">
                        {activity.action}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        {activity.resource}
                      </p>
                    </div>
                    <div className="flex-shrink-0 flex items-center gap-2">
                      <Clock className="h-3 w-3 text-muted-foreground" />
                      <Badge variant="secondary" className="text-xs bg-secondary/50">
                        {activity.time}
                      </Badge>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Organization Info */}
        <Card className="border-border/50">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-foreground">
              <div className="p-2 rounded-lg bg-primary/10">
                <Building className="h-5 w-5 text-primary" />
              </div>
              Organization Information
            </CardTitle>
            <CardDescription>
              Details about your organization and permissions
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <h3 className="font-medium text-foreground">Organization</h3>
                <p className="text-sm text-muted-foreground">{user?.organization.name}</p>
              </div>
              <div className="space-y-2">
                <h3 className="font-medium text-foreground">Your Role</h3>
                <Badge variant="outline" className="bg-secondary/20">{user?.role.name}</Badge>
              </div>
              <div className="space-y-2">
                <h3 className="font-medium text-foreground">Permissions</h3>
                <div className="flex flex-wrap gap-1">
                  {user?.role.permissions.map((permission) => (
                    <Badge key={permission.id} variant="secondary" className="bg-accent/50 text-accent-foreground">
                      {permission.name}
                    </Badge>
                  ))}
                </div>
              </div>
              <div className="space-y-2">
                <h3 className="font-medium text-foreground">Account</h3>
                <p className="text-sm text-muted-foreground">{user?.email}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </Layout>
  );
} 