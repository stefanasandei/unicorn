export interface Permission {
  id: number;
  name: string;
}

export interface Role {
  id: string;
  name: string;
  permissions: Permission[];
  created_at: string;
  updated_at: string;
}

export interface Account {
  id: string;
  name: string;
  email: string;
  type: "user" | "bot";
  organization_id: string;
  role_id: string;
  created_at: string;
  updated_at: string;
  last_login_at?: string;
}

export interface Organization {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
  accounts?: Account[];
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  token_type: string;
  expires_at: string;
  message: string;
}

export interface CreateOrganizationRequest {
  name: string;
}

export interface CreateUserRequest {
  name: string;
  email: string;
  password: string;
  role_id: string;
}

export interface CreateRoleRequest {
  name: string;
  permissions: number[];
}

export interface Secret {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
  user_id: string;
  metadata?: string;
}

export interface SecretCreateRequest {
  name: string;
  value: string;
  metadata?: string;
}

export interface SecretUpdateRequest {
  value?: string;
  metadata?: string;
}

export interface StorageBucket {
  id: string;
  name: string;
  user_id: string;
  created_at: string;
  updated_at: string;
  files?: StorageFile[];
}

export interface StorageFile {
  id: string;
  name: string;
  size: number;
  content_type: string;
  bucket_id: string;
  created_at: string;
  updated_at: string;
}

export interface File {
  name: string;
  contents: string;
}

export interface ComputeContainerInfo {
  id: string;
  name: string;
  image: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface ComputeCreateRequest {
  name: string;
  image: string;
  command?: string[];
  environment?: Record<string, string>;
  ports?: Record<string, string>;
  volumes?: Record<string, string>;
}

export interface ProcessInfo {
  stdin?: string;
  time?: string;
  max_opened_files?: number;
  max_processes?: number;
  permissions?: Permissions;
  env?: { [key: string]: string };
  working_directory?: string;
}

export interface Permissions {
  read?: boolean;
  write?: boolean;
  network?: boolean;
}

export interface ExecutionRequest {
  runtime: {
    name: string;
    version?: string;
  };
  project: {
    entry?: string;
    files: File[];
  };
  process?: ProcessInfo;
}

export interface ProcessResult {
  stdout: string;
  stderr: string;
  output: string;
  time: number; // ms
  memory: bigint; // bytes
  exit_code: number;
}

export interface WorkerResponse {
  compile: ProcessResult;
  run: ProcessResult;
}

export type ExecutionTaskStatus = "successful" | "error" | "failed";

export interface ResponseTask {
  status: ExecutionTaskStatus;
  output: WorkerResponse;
}

// Legacy types for backward compatibility
export interface LambdaExecuteRequest {
  runtime: string;
  code: string;
  handler?: string;
  environment?: Record<string, string>;
  timeout?: number;
}

export interface LambdaExecuteResponse {
  result: string;
  logs: string;
  execution_time: number;
  memory_used: number;
}

export interface ApiError {
  error: string;
  details?: string;
  status_code: number;
  timestamp: string;
}

export interface User {
  id: string;
  name: string;
  email: string;
  role: Role;
  organization: Organization;
}

export interface AuthContext {
  user: User | null;
  token: string | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  isLoading: boolean;
}
