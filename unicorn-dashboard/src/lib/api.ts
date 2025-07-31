import axios, { AxiosInstance, AxiosResponse } from "axios";
import {
  LoginRequest,
  LoginResponse,
  CreateOrganizationRequest,
  CreateUserRequest,
  CreateRoleRequest,
  Secret,
  SecretCreateRequest,
  SecretUpdateRequest,
  StorageBucket,
  StorageFile,
  ComputeContainerInfo,
  ComputeCreateRequest,
  LambdaExecuteRequest,
  LambdaExecuteResponse,
  Role,
  Organization,
  ApiError,
  ExecutionRequest,
  ResponseTask,
  RDBInstanceInfo,
  RDBCreateRequest,
} from "@/types/api";

class ApiClient {
  private client: AxiosInstance;
  private baseURL: string;

  constructor() {
    this.baseURL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    this.client = axios.create({
      baseURL: this.baseURL,
      headers: {
        "Content-Type": "application/json",
      },
    });

    // Add request interceptor to include auth token
    this.client.interceptors.request.use((config) => {
      const token = localStorage.getItem("token");
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    // Add response interceptor to handle errors
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          localStorage.removeItem("token");
          window.location.href = "/login";
        }
        return Promise.reject(error);
      }
    );
  }

  // Auth endpoints
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await this.client.post<LoginResponse>(
      "/api/v1/login",
      credentials
    );
    return response.data;
  }

  async refreshToken(token: string): Promise<LoginResponse> {
    const response = await this.client.post<LoginResponse>(
      "/api/v1/token/refresh",
      { token }
    );
    return response.data;
  }

  async validateToken(token: string): Promise<{
    valid: boolean;
    claims: { account_id: string; role_id: string; exp: number };
  }> {
    const response = await this.client.get(
      `/api/v1/token/validate?token=${token}`
    );
    return response.data;
  }

  // Organization setup endpoints
  async createOrganization(
    data: CreateOrganizationRequest
  ): Promise<{ organization: Organization; message: string }> {
    const response = await this.client.post("/api/v1/organizations", data);
    return response.data;
  }

  async createRole(
    data: CreateRoleRequest
  ): Promise<{ role: Role; message: string }> {
    const response = await this.client.post("/api/v1/roles", data);
    return response.data;
  }

  async getOrganizations(): Promise<{
    organization_name: string;
    users: Array<{
      id: string;
      name: string;
      email?: string;
      role_id?: string;
    }>;
  }> {
    const response = await this.client.get("/api/v1/organizations");
    return response.data;
  }

  async getCurrentUser(): Promise<{
    account_id: string;
    organization_id: string;
    role_id: string;
  }> {
    // Get account details from the backend
    const response = await this.client.get("/api/v1/accounts/me");
    return {
      account_id: response.data.account.id,
      organization_id: response.data.account.organization_id,
      role_id: response.data.account.role_id,
    };
  }

  async createUser(
    orgId: string,
    data: CreateUserRequest
  ): Promise<{
    account: { id: string; name: string; email: string };
    message: string;
  }> {
    const response = await this.client.post(
      `/api/v1/organizations/${orgId}/users`,
      data
    );
    return response.data;
  }

  // IAM endpoints
  async getRoles(): Promise<{ roles: Role[] }> {
    const response = await this.client.get("/api/v1/roles");
    return response.data;
  }

  // Secrets endpoints
  async listSecrets(): Promise<Secret[]> {
    const response = await this.client.get("/api/v1/secrets");
    return response.data;
  }

  async createSecret(data: SecretCreateRequest): Promise<Secret> {
    const response = await this.client.post("/api/v1/secrets", data);
    return response.data;
  }

  async getSecret(id: string): Promise<Secret & { value: string }> {
    const response = await this.client.get(`/api/v1/secrets/${id}`);
    const data = response.data;

    // Handle inconsistent metadata response from API
    // ReadSecret returns metadata as map, but ListSecrets returns it as string
    let metadata: string | undefined;
    if (data.metadata) {
      if (typeof data.metadata === "string") {
        metadata = data.metadata;
      } else {
        // Convert map to JSON string
        metadata = JSON.stringify(data.metadata);
      }
    }

    return {
      ...data,
      metadata,
    };
  }

  async updateSecret(id: string, data: SecretUpdateRequest): Promise<void> {
    await this.client.put(`/api/v1/secrets/${id}`, data);
  }

  async deleteSecret(id: string): Promise<void> {
    await this.client.delete(`/api/v1/secrets/${id}`);
  }

  // Storage endpoints
  async listBuckets(): Promise<StorageBucket[]> {
    const response = await this.client.get("/api/v1/buckets");
    return response.data;
  }

  async createBucket(name: string): Promise<StorageBucket> {
    const response = await this.client.post("/api/v1/buckets", { name });
    return response.data;
  }

  async listFiles(bucketId: string): Promise<StorageFile[]> {
    const response = await this.client.get(`/api/v1/buckets/${bucketId}/files`);
    return response.data;
  }

  async uploadFile(bucketId: string, file: File): Promise<StorageFile> {
    const formData = new FormData();
    formData.append("file", file);
    const response = await this.client.post(
      `/api/v1/buckets/${bucketId}/files`,
      formData,
      {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      }
    );
    return response.data;
  }

  async downloadFile(bucketId: string, fileId: string): Promise<Blob> {
    const response = await this.client.get(
      `/api/v1/buckets/${bucketId}/files/${fileId}`,
      {
        responseType: "blob",
      }
    );
    return response.data;
  }

  async deleteFile(bucketId: string, fileId: string): Promise<void> {
    await this.client.delete(`/api/v1/buckets/${bucketId}/files/${fileId}`);
  }

  // Compute endpoints
  async listCompute(): Promise<ComputeContainerInfo[]> {
    const response = await this.client.get("/api/v1/compute/list");
    return response.data;
  }

  async createCompute(
    data: ComputeCreateRequest
  ): Promise<ComputeContainerInfo> {
    const response = await this.client.post("/api/v1/compute/create", data);
    return response.data;
  }

  async deleteCompute(id: string): Promise<void> {
    await this.client.delete(`/api/v1/compute/${id}`);
  }

  // Lambda endpoints
  async executeLambda(
    data: LambdaExecuteRequest
  ): Promise<LambdaExecuteResponse> {
    const response = await this.client.post("/api/v1/lambda/execute", data);
    return response.data;
  }

  async testLambda(data: LambdaExecuteRequest): Promise<LambdaExecuteResponse> {
    const response = await this.client.post("/api/v1/lambda/test", data);
    return response.data;
  }

  // Execution endpoints (sandbox microservice)
  async executeCode(data: ExecutionRequest): Promise<ResponseTask> {
    const response = await this.client.post("/api/v1/lambda/execute", data);
    return response.data;
  }

  async getRuntimes(): Promise<
    Array<{ language: string; versions: string[] }>
  > {
    return [
      { language: "python3", versions: ["3.12", "3.11", "3.9"] },
      { language: "go", versions: ["1.20", "1.19"] },
    ];
  }

  // Debug endpoint
  async getDebugToken(): Promise<{
    token: string;
    token_type: string;
    message: string;
  }> {
    const response = await this.client.get("/api/v1/debug/token");
    return response.data;
  }

  // RDB endpoints
  async listRDB(): Promise<RDBInstanceInfo[]> {
    const response = await this.client.get("/api/v1/rdb/list");
    return response.data;
  }

  async createRDB(data: RDBCreateRequest): Promise<RDBInstanceInfo> {
    const response = await this.client.post("/api/v1/rdb/create", data);
    return response.data;
  }

  async deleteRDB(id: string): Promise<void> {
    await this.client.delete(`/api/v1/rdb/${id}`);
  }

  // Monitoring endpoints
  async getResourceUsage(): Promise<any> {
    const response = await this.client.get("/api/v1/monitoring/usage");
    return response.data;
  }

  async getBillingHistory(): Promise<any[]> {
    const response = await this.client.get("/api/v1/monitoring/billing");
    return response.data;
  }

  async getMonthlyTrends(months: number = 6): Promise<any[]> {
    const response = await this.client.get(
      `/api/v1/monitoring/trends?months=${months}`
    );
    return response.data;
  }

  async getActiveResources(): Promise<any[]> {
    const response = await this.client.get(
      "/api/v1/monitoring/resources/active"
    );
    return response.data;
  }
}

export const apiClient = new ApiClient();
