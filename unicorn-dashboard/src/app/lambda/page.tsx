"use client";

import React, { useState, useEffect } from "react";
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
import { Checkbox } from "@/components/ui/checkbox";
import {
  Code,
  Play,
  Zap,
  Clock,
  Cpu,
  AlertCircle,
  Settings,
} from "lucide-react";
import { apiClient } from "@/lib/api";
import { ExecutionRequest, ResponseTask, File, Permissions } from "@/types/api";

export default function LambdaPage() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [result, setResult] = useState<ResponseTask | null>(null);
  const [runtimes, setRuntimes] = useState<
    Array<{ language: string; versions: string[] }>
  >([]);

  // Form states
  const [runtimeName, setRuntimeName] = useState("python3");
  const [runtimeVersion, setRuntimeVersion] = useState("3.12");
  const [entryCode, setEntryCode] = useState("");
  const [files, setFiles] = useState<File[]>([]);
  const [stdin, setStdin] = useState("");
  const [timeout, setTimeout] = useState("500ms");
  const [environment, setEnvironment] = useState("");

  // Granular control states
  const [maxOpenedFiles, setMaxOpenedFiles] = useState<number>(10);
  const [maxProcesses, setMaxProcesses] = useState<number>(5);
  const [permissions, setPermissions] = useState<Permissions>({
    read: false,
    write: false,
    network: false,
  });

  // Dialog states
  const [showErrorDialog, setShowErrorDialog] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");

  useEffect(() => {
    fetchRuntimes();
    loadPreset();
  }, []);

  const fetchRuntimes = async () => {
    try {
      const runtimeData = await apiClient.getRuntimes();
      setRuntimes(runtimeData);
    } catch (err: unknown) {
      console.error("Failed to fetch runtimes:", err);
    }
  };

  const loadPreset = () => {
    // Load the preset from the documentation
    setRuntimeName("python3");
    setRuntimeVersion("3.12");
    setEntryCode("import utils\nprint(utils.add(int(input()),2))");
    setFiles([
      {
        name: "utils.py",
        contents: "def add(a, b):\n\treturn a + b",
      },
    ]);
    setStdin("2");
    setTimeout("500ms");
    setEnvironment('{"PROD": "false"}');
  };

  const loadGoPreset = () => {
    setRuntimeName("go");
    setRuntimeVersion("1.20");
    setEntryCode(`package main

import (
	"fmt"
)

func main() {
	var n int
	fmt.Scanln(&n)
	fmt.Println(n + 2)
}`);
    setFiles([
      {
        name: "utils.go",
        contents: `package main

func multiply(a, b int) int {
	return a * b
}`,
      },
    ]);
    setStdin("5");
    setTimeout("1s");
    setEnvironment('{"DEBUG": "true"}');
  };

  const handleExecute = async () => {
    if (!entryCode.trim()) {
      setErrorMessage("Entry code is required");
      setShowErrorDialog(true);
      return;
    }

    setIsLoading(true);
    setError("");
    setResult(null);

    try {
      let env = {};
      if (environment.trim()) {
        try {
          env = JSON.parse(environment);
        } catch {
          setErrorMessage("Invalid JSON environment");
          setShowErrorDialog(true);
          return;
        }
      }

      const request: ExecutionRequest = {
        runtime: {
          name: runtimeName,
          version: runtimeVersion,
        },
        project: {
          entry: entryCode,
          files: files,
        },
        process: {
          stdin: stdin || undefined,
          time: timeout || undefined,
          env: Object.keys(env).length > 0 ? env : undefined,
          max_opened_files: maxOpenedFiles,
          max_processes: maxProcesses,
          permissions: permissions,
        },
      };

      console.log("Executing code with request:", request);

      const response = await apiClient.executeCode(request);
      setResult(response);
    } catch (err: unknown) {
      console.error("Execute error:", err);
      let errorMsg = "Failed to execute code";
      if (err && typeof err === "object" && "response" in err) {
        const response = (
          err as {
            response?: { data?: { message?: string; details?: string } };
          }
        ).response;
        errorMsg =
          response?.data?.message || response?.data?.details || errorMsg;
      }
      setErrorMessage(errorMsg);
      setShowErrorDialog(true);
    } finally {
      setIsLoading(false);
    }
  };

  const addFile = () => {
    setFiles([...files, { name: "", contents: "" }]);
  };

  const updateFile = (
    index: number,
    field: "name" | "contents",
    value: string
  ) => {
    const newFiles = [...files];
    newFiles[index] = { ...newFiles[index], [field]: value };
    setFiles(newFiles);
  };

  const removeFile = (index: number) => {
    setFiles(files.filter((_, i) => i !== index));
  };

  const getRuntimeVersions = () => {
    const runtime = runtimes.find((r) => r.language === runtimeName);
    return runtime?.versions || [];
  };

  return (
    <Layout>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Code Execution</h1>
          <p className="text-muted-foreground">
            Execute code in isolated environments with various runtimes
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Code Editor */}
          <Card>
            <CardHeader>
              <CardTitle>Code Editor</CardTitle>
              <CardDescription>
                Write and configure your code execution
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="runtime">Runtime</Label>
                  <select
                    id="runtime"
                    value={runtimeName}
                    onChange={(e) => setRuntimeName(e.target.value)}
                    className="w-full p-2 border border-border/50 rounded-md bg-background text-foreground focus:border-primary focus:ring-primary/20"
                  >
                    {runtimes.map((rt) => (
                      <option key={rt.language} value={rt.language}>
                        {rt.language}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <Label htmlFor="version">Version</Label>
                  <select
                    id="version"
                    value={runtimeVersion}
                    onChange={(e) => setRuntimeVersion(e.target.value)}
                    className="w-full p-2 border border-border/50 rounded-md bg-background text-foreground focus:border-primary focus:ring-primary/20"
                  >
                    {getRuntimeVersions().map((version) => (
                      <option key={version} value={version}>
                        {version}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <div>
                <Label htmlFor="timeout">Timeout</Label>
                <Input
                  id="timeout"
                  value={timeout}
                  onChange={(e) => setTimeout(e.target.value)}
                  placeholder="500ms"
                />
                <p className="text-sm text-muted-foreground mt-1">
                  Time limit for execution (e.g., &quot;500ms&quot;,
                  &quot;5s&quot;)
                </p>
              </div>

              <div>
                <Label htmlFor="stdin">Standard Input</Label>
                <Input
                  id="stdin"
                  value={stdin}
                  onChange={(e) => setStdin(e.target.value)}
                  placeholder="Enter input for your program"
                />
              </div>

              <div>
                <Label htmlFor="environment">
                  Environment Variables (JSON)
                </Label>
                <Input
                  id="environment"
                  value={environment}
                  onChange={(e) => setEnvironment(e.target.value)}
                  placeholder='{"PROD": "false"}'
                />
              </div>

              {/* Advanced Settings */}
              <div className="space-y-4">
                <div className="flex items-center space-x-2">
                  <Settings className="h-4 w-4" />
                  <Label className="text-sm font-medium">
                    Advanced Settings
                  </Label>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label htmlFor="maxOpenedFiles">Max Opened Files</Label>
                    <Input
                      id="maxOpenedFiles"
                      type="number"
                      value={maxOpenedFiles}
                      onChange={(e) =>
                        setMaxOpenedFiles(parseInt(e.target.value) || 10)
                      }
                      min="1"
                      max="100"
                    />
                  </div>
                  <div>
                    <Label htmlFor="maxProcesses">Max Processes</Label>
                    <Input
                      id="maxProcesses"
                      type="number"
                      value={maxProcesses}
                      onChange={(e) =>
                        setMaxProcesses(parseInt(e.target.value) || 5)
                      }
                      min="1"
                      max="20"
                    />
                  </div>
                </div>

                <div>
                  <Label className="text-sm font-medium mb-2 block">
                    Permissions
                  </Label>
                  <div className="space-y-2">
                    <div className="flex items-center space-x-2">
                      <Checkbox
                        id="readPermission"
                        checked={permissions.read}
                        onCheckedChange={(checked: boolean | "indeterminate") =>
                          setPermissions({
                            ...permissions,
                            read: checked === true,
                          })
                        }
                      />
                      <Label htmlFor="readPermission" className="text-sm">
                        Read Access
                      </Label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <Checkbox
                        id="writePermission"
                        checked={permissions.write}
                        onCheckedChange={(checked: boolean | "indeterminate") =>
                          setPermissions({
                            ...permissions,
                            write: checked === true,
                          })
                        }
                      />
                      <Label htmlFor="writePermission" className="text-sm">
                        Write Access
                      </Label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <Checkbox
                        id="networkPermission"
                        checked={permissions.network}
                        onCheckedChange={(checked: boolean | "indeterminate") =>
                          setPermissions({
                            ...permissions,
                            network: checked === true,
                          })
                        }
                      />
                      <Label htmlFor="networkPermission" className="text-sm">
                        Network Access
                      </Label>
                    </div>
                  </div>
                  <p className="text-xs text-muted-foreground mt-2">
                    Control what your code can access. Read allows file reading,
                    Write allows file creation/modification, Network allows HTTP
                    requests.
                  </p>
                </div>
              </div>

              <div>
                <Label htmlFor="entry">Entry Code</Label>
                <textarea
                  id="entry"
                  value={entryCode}
                  onChange={(e) => setEntryCode(e.target.value)}
                  className="w-full h-32 p-3 border border-border/50 rounded-md font-mono text-sm bg-background text-foreground focus:border-primary focus:ring-primary/20"
                  placeholder="// Write your main code here..."
                />
              </div>

              <div>
                <div className="flex items-center justify-between mb-2">
                  <Label>Additional Files</Label>
                  <Button variant="outline" size="sm" onClick={addFile}>
                    Add File
                  </Button>
                </div>
                <div className="space-y-2">
                  {files.map((file, index) => (
                    <div
                      key={index}
                      className="border border-border/50 rounded-md p-3 bg-card"
                    >
                      <div className="flex items-center justify-between mb-2">
                        <Input
                          placeholder="filename.py"
                          value={file.name}
                          onChange={(e) =>
                            updateFile(index, "name", e.target.value)
                          }
                          className="w-1/2"
                        />
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => removeFile(index)}
                        >
                          Remove
                        </Button>
                      </div>
                      <textarea
                        placeholder="File contents..."
                        value={file.contents}
                        onChange={(e) =>
                          updateFile(index, "contents", e.target.value)
                        }
                        className="w-full h-20 p-2 border border-border/50 rounded-md font-mono text-sm bg-background text-foreground focus:border-primary focus:ring-primary/20"
                      />
                    </div>
                  ))}
                </div>
              </div>

              <div className="flex space-x-2">
                <Button
                  variant="outline"
                  onClick={loadPreset}
                  className="flex-1"
                >
                  Load Python Template
                </Button>
                <Button
                  variant="outline"
                  onClick={loadGoPreset}
                  className="flex-1"
                >
                  Load Go Template
                </Button>
              </div>

              <Button
                onClick={handleExecute}
                disabled={isLoading}
                className="w-full"
              >
                <Play className="h-4 w-4 mr-2" />
                {isLoading ? "Executing..." : "Execute Code"}
              </Button>
            </CardContent>
          </Card>

          {/* Results */}
          <Card>
            <CardHeader>
              <CardTitle>Execution Results</CardTitle>
              <CardDescription>
                Code execution output and metrics
              </CardDescription>
            </CardHeader>
            <CardContent>
              {result ? (
                <div className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div className="flex items-center space-x-2">
                      <Clock className="h-4 w-4 text-primary" />
                      <span className="text-sm font-medium text-foreground">
                        Compile Time:
                      </span>
                      <span className="text-sm text-muted-foreground">
                        {result.output.compile.time}ms
                      </span>
                    </div>
                    <div className="flex items-center space-x-2">
                      <Cpu className="h-4 w-4 text-success" />
                      <span className="text-sm font-medium text-foreground">
                        Run Time:
                      </span>
                      <span className="text-sm text-muted-foreground">
                        {result.output.run.time}ms
                      </span>
                    </div>
                    <div className="flex items-center space-x-2">
                      <Code className="h-4 w-4 text-accent" />
                      <span className="text-sm font-medium text-foreground">
                        Status:
                      </span>
                      <Badge
                        variant={
                          result.status === "successful"
                            ? "default"
                            : "destructive"
                        }
                      >
                        {result.status}
                      </Badge>
                    </div>
                  </div>

                  <Tabs defaultValue="run" className="space-y-4">
                    <TabsList>
                      <TabsTrigger value="run">Run Output</TabsTrigger>
                      <TabsTrigger value="compile">Compile Output</TabsTrigger>
                    </TabsList>

                    <TabsContent value="run">
                      <div className="space-y-2">
                        <div>
                          <Label className="text-sm font-medium text-foreground">
                            Standard Output
                          </Label>
                          <div className="bg-muted/50 p-3 rounded-md">
                            <pre className="text-sm overflow-x-auto text-foreground">
                              {result.output.run.stdout || "(empty)"}
                            </pre>
                          </div>
                        </div>
                        <div>
                          <Label className="text-sm font-medium text-foreground">
                            Standard Error
                          </Label>
                          <div className="bg-muted/50 p-3 rounded-md">
                            <pre className="text-sm overflow-x-auto text-foreground">
                              {result.output.run.stderr || "(empty)"}
                            </pre>
                          </div>
                        </div>
                        <div className="flex items-center space-x-2">
                          <span className="text-sm font-medium text-foreground">
                            Exit Code:
                          </span>
                          <Badge
                            variant={
                              result.output.run.exit_code === 0
                                ? "default"
                                : "destructive"
                            }
                          >
                            {result.output.run.exit_code}
                          </Badge>
                        </div>
                      </div>
                    </TabsContent>

                    <TabsContent value="compile">
                      <div className="space-y-2">
                        <div>
                          <Label className="text-sm font-medium text-foreground">
                            Standard Output
                          </Label>
                          <div className="bg-muted/50 p-3 rounded-md">
                            <pre className="text-sm overflow-x-auto text-foreground">
                              {result.output.compile.stdout || "(empty)"}
                            </pre>
                          </div>
                        </div>
                        <div>
                          <Label className="text-sm font-medium text-foreground">
                            Standard Error
                          </Label>
                          <div className="bg-muted/50 p-3 rounded-md">
                            <pre className="text-sm overflow-x-auto text-foreground">
                              {result.output.compile.stderr || "(empty)"}
                            </pre>
                          </div>
                        </div>
                        <div className="flex items-center space-x-2">
                          <span className="text-sm font-medium text-foreground">
                            Exit Code:
                          </span>
                          <Badge
                            variant={
                              result.output.compile.exit_code === 0
                                ? "default"
                                : "destructive"
                            }
                          >
                            {result.output.compile.exit_code}
                          </Badge>
                        </div>
                      </div>
                    </TabsContent>
                  </Tabs>
                </div>
              ) : (
                <div className="text-center text-muted-foreground py-8">
                  <Code className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
                  <p>Execute code to see results here</p>
                </div>
              )}
            </CardContent>
          </Card>
        </div>

        {/* Error Dialog */}
        <Dialog open={showErrorDialog} onOpenChange={setShowErrorDialog}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2 text-foreground">
                <AlertCircle className="h-5 w-5 text-destructive" />
                Error
              </DialogTitle>
            </DialogHeader>
            <div className="py-4">
              <p className="text-foreground">{errorMessage}</p>
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
