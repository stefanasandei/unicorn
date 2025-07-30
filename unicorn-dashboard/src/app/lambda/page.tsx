'use client';

import React, { useState } from 'react';
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
import { Code, Play, Zap, Clock, Cpu } from 'lucide-react';
import { apiClient } from '@/lib/api';
import { LambdaExecuteRequest, LambdaExecuteResponse } from '@/types/api';

export default function LambdaPage() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [result, setResult] = useState<LambdaExecuteResponse | null>(null);

  // Form states
  const [runtime, setRuntime] = useState('nodejs18.x');
  const [code, setCode] = useState('');
  const [handler, setHandler] = useState('index.handler');
  const [environment, setEnvironment] = useState('');
  const [timeout, setTimeout] = useState('30');

  const runtimes = [
    { value: 'nodejs18.x', label: 'Node.js 18.x' },
    { value: 'python3.9', label: 'Python 3.9' },
    { value: 'go1.x', label: 'Go 1.x' },
    { value: 'java11', label: 'Java 11' },
  ];

  const templates = [
    {
      name: 'Hello World',
      description: 'Simple function that returns a greeting',
      runtime: 'nodejs18.x',
      code: `exports.handler = async (event) => {
  const response = {
    statusCode: 200,
    body: JSON.stringify({
      message: 'Hello from Lambda!',
      event: event
    }),
  };
  return response;
};`,
    },
    {
      name: 'Data Processing',
      description: 'Process JSON data and return transformed result',
      runtime: 'python3.9',
      code: `import json

def lambda_handler(event, context):
    try:
        # Parse input data
        data = event.get('data', {})
        
        # Process the data
        processed_data = {
            'processed': True,
            'input': data,
            'timestamp': context.get_remaining_time_in_millis()
        }
        
        return {
            'statusCode': 200,
            'body': json.dumps(processed_data)
        }
    except Exception as e:
        return {
            'statusCode': 500,
            'body': json.dumps({'error': str(e)})
        }`,
    },
    {
      name: 'API Handler',
      description: 'Handle HTTP requests and return JSON responses',
      runtime: 'nodejs18.x',
      code: `exports.handler = async (event) => {
  const method = event.httpMethod;
  const path = event.path;
  
  let response = {
    statusCode: 200,
    headers: {
      'Content-Type': 'application/json',
      'Access-Control-Allow-Origin': '*'
    },
    body: JSON.stringify({
      message: 'API response',
      method: method,
      path: path,
      timestamp: new Date().toISOString()
    })
  };
  
  return response;
};`,
    },
  ];

  const handleExecute = async () => {
    if (!code.trim()) {
      setError('Code is required');
      return;
    }

    setIsLoading(true);
    setError('');
    setResult(null);

    try {
      let env = {};
      if (environment.trim()) {
        try {
          env = JSON.parse(environment);
        } catch {
          setError('Invalid JSON environment');
          return;
        }
      }

      const request: LambdaExecuteRequest = {
        runtime,
        code,
        handler,
        environment: env,
        timeout: parseInt(timeout),
      };

      const response = await apiClient.executeLambda(request);
      setResult(response);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to execute Lambda function');
    } finally {
      setIsLoading(false);
    }
  };

  const handleTest = async () => {
    if (!code.trim()) {
      setError('Code is required');
      return;
    }

    setIsLoading(true);
    setError('');
    setResult(null);

    try {
      let env = {};
      if (environment.trim()) {
        try {
          env = JSON.parse(environment);
        } catch {
          setError('Invalid JSON environment');
          return;
        }
      }

      const request: LambdaExecuteRequest = {
        runtime,
        code,
        handler,
        environment: env,
        timeout: parseInt(timeout),
      };

      const response = await apiClient.testLambda(request);
      setResult(response);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to test Lambda function');
    } finally {
      setIsLoading(false);
    }
  };

  const loadTemplate = (template: any) => {
    setRuntime(template.runtime);
    setCode(template.code);
    setHandler('index.handler');
    setEnvironment('');
    setTimeout('30');
  };

  return (
    <Layout>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Lambda</h1>
          <p className="text-gray-600">
            Execute serverless functions with various runtimes
          </p>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 rounded-md p-4">
            <p className="text-red-800">{error}</p>
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Function Editor */}
          <Card>
            <CardHeader>
              <CardTitle>Function Editor</CardTitle>
              <CardDescription>
                Write and configure your Lambda function
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="runtime">Runtime</Label>
                  <select
                    id="runtime"
                    value={runtime}
                    onChange={(e) => setRuntime(e.target.value)}
                    className="w-full p-2 border border-gray-300 rounded-md"
                  >
                    {runtimes.map((rt) => (
                      <option key={rt.value} value={rt.value}>
                        {rt.label}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <Label htmlFor="timeout">Timeout (seconds)</Label>
                  <Input
                    id="timeout"
                    type="number"
                    value={timeout}
                    onChange={(e) => setTimeout(e.target.value)}
                    min="1"
                    max="900"
                  />
                </div>
              </div>

              <div>
                <Label htmlFor="handler">Handler</Label>
                <Input
                  id="handler"
                  value={handler}
                  onChange={(e) => setHandler(e.target.value)}
                  placeholder="index.handler"
                />
              </div>

              <div>
                <Label htmlFor="environment">Environment Variables (JSON)</Label>
                <Input
                  id="environment"
                  value={environment}
                  onChange={(e) => setEnvironment(e.target.value)}
                  placeholder='{"NODE_ENV": "production"}'
                />
              </div>

              <div>
                <Label htmlFor="code">Function Code</Label>
                <textarea
                  id="code"
                  value={code}
                  onChange={(e) => setCode(e.target.value)}
                  className="w-full h-64 p-3 border border-gray-300 rounded-md font-mono text-sm"
                  placeholder="// Write your Lambda function code here..."
                />
              </div>

              <div className="flex space-x-2">
                <Button onClick={handleExecute} disabled={isLoading}>
                  <Play className="h-4 w-4 mr-2" />
                  {isLoading ? 'Executing...' : 'Execute'}
                </Button>
                <Button variant="outline" onClick={handleTest} disabled={isLoading}>
                  <Zap className="h-4 w-4 mr-2" />
                  Test
                </Button>
              </div>
            </CardContent>
          </Card>

          {/* Templates */}
          <Card>
            <CardHeader>
              <CardTitle>Code Templates</CardTitle>
              <CardDescription>
                Quick start with pre-built function templates
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {templates.map((template) => (
                  <Card key={template.name} className="cursor-pointer hover:shadow-md transition-shadow">
                    <CardContent className="p-4">
                      <div className="flex items-center justify-between">
                        <div>
                          <h3 className="font-medium">{template.name}</h3>
                          <p className="text-sm text-muted-foreground">
                            {template.description}
                          </p>
                          <Badge variant="secondary" className="mt-1">
                            {template.runtime}
                          </Badge>
                        </div>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => loadTemplate(template)}
                        >
                          Use Template
                        </Button>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Results */}
        {result && (
          <Card>
            <CardHeader>
              <CardTitle>Execution Results</CardTitle>
              <CardDescription>
                Function execution output and metrics
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                <div className="flex items-center space-x-2">
                  <Clock className="h-4 w-4 text-blue-500" />
                  <span className="text-sm font-medium">Execution Time:</span>
                  <span className="text-sm">{result.execution_time}ms</span>
                </div>
                <div className="flex items-center space-x-2">
                  <Cpu className="h-4 w-4 text-green-500" />
                  <span className="text-sm font-medium">Memory Used:</span>
                  <span className="text-sm">{result.memory_used}MB</span>
                </div>
                <div className="flex items-center space-x-2">
                  <Code className="h-4 w-4 text-purple-500" />
                  <span className="text-sm font-medium">Status:</span>
                  <Badge variant="secondary">Success</Badge>
                </div>
              </div>

              <Tabs defaultValue="result" className="space-y-4">
                <TabsList>
                  <TabsTrigger value="result">Result</TabsTrigger>
                  <TabsTrigger value="logs">Logs</TabsTrigger>
                </TabsList>

                <TabsContent value="result">
                  <div className="bg-gray-50 p-4 rounded-md">
                    <pre className="text-sm overflow-x-auto">
                      {result.result}
                    </pre>
                  </div>
                </TabsContent>

                <TabsContent value="logs">
                  <div className="bg-gray-50 p-4 rounded-md">
                    <pre className="text-sm overflow-x-auto">
                      {result.logs}
                    </pre>
                  </div>
                </TabsContent>
              </Tabs>
            </CardContent>
          </Card>
        )}
      </div>
    </Layout>
  );
} 