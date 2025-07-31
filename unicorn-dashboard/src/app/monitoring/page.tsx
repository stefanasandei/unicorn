"use client";

import React, { useEffect, useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { Layout } from "@/components/Layout";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Activity,
  TrendingUp,
  DollarSign,
  Server,
  Database,
  FileText,
  Code,
  Eye,
  Clock,
  AlertTriangle,
} from "lucide-react";
import { apiClient } from "@/lib/api";

interface ResourceUsage {
  id: string;
  resource_type: string;
  resource_name: string;
  status: string;
  cpu_usage: number;
  memory_usage: number;
  storage_usage: number;
  network_usage: number;
  cost_per_hour: number;
  total_cost: number;
  currency: string;
  last_active_at?: string;
  resource_created_at: string;
}

interface BillingPeriod {
  id: string;
  period_start: string;
  period_end: string;
  total_cost: number;
  currency: string;
  is_paid: boolean;
  compute_cost: number;
  lambda_cost: number;
  storage_cost: number;
  rdb_cost: number;
  secret_cost: number;
}

interface UsageSummary {
  total_resources: number;
  active_resources: number;
  total_cost: number;
  currency: string;
  usage_by_type: Record<string, number>;
  cost_by_type: Record<string, number>;
}

interface MonthlyTrend {
  month: string;
  total_cost: number;
  resources: number;
}

export default function MonitoringPage() {
  const { user } = useAuth();
  const [resourceUsage, setResourceUsage] = useState<ResourceUsage[]>([]);
  const [billingHistory, setBillingHistory] = useState<BillingPeriod[]>([]);
  const [usageSummary, setUsageSummary] = useState<UsageSummary | null>(null);
  const [monthlyTrends, setMonthlyTrends] = useState<MonthlyTrend[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [activeTab, setActiveTab] = useState("overview");

  useEffect(() => {
    const fetchMonitoringData = async () => {
      try {
        setIsLoading(true);

        // Fetch active resources
        const activeResourcesResponse = await apiClient.getActiveResources();
        setResourceUsage(activeResourcesResponse || []);

        // Fetch resource usage summary
        const usageResponse = await apiClient.getResourceUsage();

        // Fetch billing history
        const billingResponse = await apiClient.getBillingHistory();
        setBillingHistory(billingResponse);

        // Fetch usage summary
        if (usageResponse.summary) {
          setUsageSummary(usageResponse.summary);
        }

        // Fetch monthly trends
        const trendsResponse = await apiClient.getMonthlyTrends(6);
        setMonthlyTrends(trendsResponse);
      } catch (error) {
        console.error("Failed to fetch monitoring data:", error);
        // Show error state instead of mock data
        setResourceUsage([]);
        setBillingHistory([]);
        setUsageSummary(null);
        setMonthlyTrends([]);
      } finally {
        setIsLoading(false);
      }
    };

    fetchMonitoringData();
  }, []);

  const getResourceTypeIcon = (type: string) => {
    switch (type) {
      case "compute":
        return <Server className="h-4 w-4" />;
      case "storage":
        return <FileText className="h-4 w-4" />;
      case "lambda":
        return <Code className="h-4 w-4" />;
      case "rdb":
        return <Database className="h-4 w-4" />;
      case "secret":
        return <Eye className="h-4 w-4" />;
      default:
        return <Activity className="h-4 w-4" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "active":
        return "bg-green-500";
      case "inactive":
        return "bg-yellow-500";
      case "deleted":
        return "bg-red-500";
      default:
        return "bg-gray-500";
    }
  };

  const formatCurrency = (amount: number, currency: string = "USD") => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: currency,
    }).format(amount);
  };

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return "0 B";
    const k = 1024;
    const sizes = ["B", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
  };

  if (isLoading) {
    return (
      <Layout>
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              Monitoring & Billing
            </h1>
            <p className="text-muted-foreground">
              Track resource usage, costs, and performance metrics
            </p>
          </div>
          <Button onClick={() => window.location.reload()}>
            <Activity className="mr-2 h-4 w-4" />
            Refresh Data
          </Button>
        </div>

        <Tabs
          value={activeTab}
          onValueChange={setActiveTab}
          className="space-y-4"
        >
          <TabsList>
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="resources">Active Resources</TabsTrigger>
            <TabsTrigger value="billing">Billing History</TabsTrigger>
            <TabsTrigger value="trends">Usage Trends</TabsTrigger>
          </TabsList>

          <TabsContent value="overview" className="space-y-4">
            {/* Usage Summary Cards */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    Total Resources
                  </CardTitle>
                  <Server className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">
                    {usageSummary?.total_resources || 0}
                  </div>
                  <p className="text-xs text-muted-foreground">
                    {usageSummary?.active_resources || 0} active
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    Total Cost
                  </CardTitle>
                  <DollarSign className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">
                    {formatCurrency(
                      usageSummary?.total_cost || 0,
                      usageSummary?.currency
                    )}
                  </div>
                  <p className="text-xs text-muted-foreground">This month</p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    Avg CPU Usage
                  </CardTitle>
                  <Activity className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">
                    {resourceUsage.length > 0
                      ? Math.round(
                          resourceUsage.reduce(
                            (sum, r) => sum + r.cpu_usage,
                            0
                          ) / resourceUsage.length
                        )
                      : 0}
                    %
                  </div>
                  <p className="text-xs text-muted-foreground">
                    Across all resources
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    Storage Used
                  </CardTitle>
                  <FileText className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">
                    {formatBytes(
                      resourceUsage.reduce((sum, r) => sum + r.storage_usage, 0)
                    )}
                  </div>
                  <p className="text-xs text-muted-foreground">Total storage</p>
                </CardContent>
              </Card>
            </div>

            {/* Cost Breakdown */}
            <Card>
              <CardHeader>
                <CardTitle>Cost Breakdown by Type</CardTitle>
                <CardDescription>
                  Distribution of costs across different resource types
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {usageSummary?.cost_by_type &&
                    Object.entries(usageSummary.cost_by_type).map(
                      ([type, cost]) => (
                        <div
                          key={type}
                          className="flex items-center justify-between"
                        >
                          <div className="flex items-center space-x-2">
                            {getResourceTypeIcon(type)}
                            <span className="capitalize">{type}</span>
                          </div>
                          <div className="text-right">
                            <div className="font-medium">
                              {formatCurrency(cost, usageSummary.currency)}
                            </div>
                            <div className="text-xs text-muted-foreground">
                              {usageSummary.usage_by_type[type] || 0} resources
                            </div>
                          </div>
                        </div>
                      )
                    )}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="resources" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>Active Resources</CardTitle>
                <CardDescription>
                  Real-time status and metrics for your resources
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {resourceUsage.map((resource) => (
                    <div key={resource.id} className="border rounded-lg p-4">
                      <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center space-x-2">
                          {getResourceTypeIcon(resource.resource_type)}
                          <span className="font-medium">
                            {resource.resource_name}
                          </span>
                          <Badge className={getStatusColor(resource.status)}>
                            {resource.status}
                          </Badge>
                        </div>
                        <div className="text-right">
                          <div className="font-medium">
                            {formatCurrency(
                              resource.total_cost,
                              resource.currency
                            )}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            ${resource.cost_per_hour}/hour
                          </div>
                        </div>
                      </div>

                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                        <div>
                          <div className="text-muted-foreground">CPU</div>
                          <div className="font-medium">
                            {resource.cpu_usage.toFixed(1)}%
                          </div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">Memory</div>
                          <div className="font-medium">
                            {formatBytes(resource.memory_usage)}
                          </div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">Storage</div>
                          <div className="font-medium">
                            {formatBytes(resource.storage_usage)}
                          </div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">Network</div>
                          <div className="font-medium">
                            {formatBytes(resource.network_usage)}
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="billing" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>Billing History</CardTitle>
                <CardDescription>
                  Monthly billing periods and cost breakdowns
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {billingHistory.map((period) => (
                    <div key={period.id} className="border rounded-lg p-4">
                      <div className="flex items-center justify-between mb-4">
                        <div>
                          <div className="font-medium">
                            {new Date(period.period_start).toLocaleDateString()}{" "}
                            - {new Date(period.period_end).toLocaleDateString()}
                          </div>
                          <div className="text-sm text-muted-foreground">
                            Billing Period
                          </div>
                        </div>
                        <div className="text-right">
                          <div className="text-2xl font-bold">
                            {formatCurrency(period.total_cost, period.currency)}
                          </div>
                          <Badge
                            variant={period.is_paid ? "default" : "secondary"}
                          >
                            {period.is_paid ? "Paid" : "Pending"}
                          </Badge>
                        </div>
                      </div>

                      <div className="grid grid-cols-2 md:grid-cols-5 gap-4 text-sm">
                        <div>
                          <div className="text-muted-foreground">Compute</div>
                          <div className="font-medium">
                            {formatCurrency(
                              period.compute_cost,
                              period.currency
                            )}
                          </div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">Lambda</div>
                          <div className="font-medium">
                            {formatCurrency(
                              period.lambda_cost,
                              period.currency
                            )}
                          </div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">Storage</div>
                          <div className="font-medium">
                            {formatCurrency(
                              period.storage_cost,
                              period.currency
                            )}
                          </div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">RDB</div>
                          <div className="font-medium">
                            {formatCurrency(period.rdb_cost, period.currency)}
                          </div>
                        </div>
                        <div>
                          <div className="text-muted-foreground">Secrets</div>
                          <div className="font-medium">
                            {formatCurrency(
                              period.secret_cost,
                              period.currency
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="trends" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>Monthly Usage Trends</CardTitle>
                <CardDescription>
                  Cost and resource usage trends over time
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {monthlyTrends.map((trend, index) => (
                    <div key={index} className="border rounded-lg p-4">
                      <div className="flex items-center justify-between">
                        <div>
                          <div className="font-medium">{trend.month}</div>
                          <div className="text-sm text-muted-foreground">
                            {trend.resources} resources
                          </div>
                        </div>
                        <div className="text-right">
                          <div className="text-xl font-bold">
                            {formatCurrency(trend.total_cost, "USD")}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            Total cost
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </Layout>
  );
}
