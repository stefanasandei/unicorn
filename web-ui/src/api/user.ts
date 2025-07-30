// API helpers for Unicorn IAM onboarding
import axios from "axios";

const API_BASE = "http://localhost:8080/api/v1";

export type Organization = {
  id: string;
  name: string;
};

export type Role = {
  id: string;
  name: string;
  permissions: number[];
};

export type Account = {
  id: string;
  name: string;
  email: string;
  role_id: string;
  organization_id: string;
};

export async function createOrganization(name: string) {
  const res = await axios.post(`${API_BASE}/organizations`, { name });
  return res.data.organization as Organization;
}

export async function createRole(name: string, permissions: number[]) {
  const res = await axios.post(`${API_BASE}/roles`, { name, permissions });
  return res.data.role as Role;
}

export async function createUser(
  orgId: string,
  name: string,
  email: string,
  password: string,
  roleId: string
) {
  const res = await axios.post(`${API_BASE}/organizations/${orgId}/users`, {
    name,
    email,
    password,
    role_id: roleId,
  });
  return res.data.account as Account;
}

export async function login(email: string, password: string) {
  const res = await axios.post(`${API_BASE}/login`, { email, password });
  return res.data as { token: string; token_type: string; expires_at: string };
}

export async function getAccountInfo(token: string) {
  // Validate token to get account_id and role_id
  const validate = await axios.get(`${API_BASE}/token/validate`, {
    params: { token },
  });
  const { account_id, role_id } = validate.data.claims;

  // Get all roles for the user's organization
  const rolesRes = await axios.get(`${API_BASE}/roles`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  const roles = rolesRes.data.roles as Role[];
  const role = roles.find((r) => r.id === role_id);

  // Get organization info and all users
  const orgRes = await axios.get(`${API_BASE}/organizations`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  const orgName = orgRes.data.organization_name as string;
  const users = orgRes.data.users as {
    id: string;
    name: string;
    role_id: string;
  }[];
  const user = users.find((u) => u.id === account_id);

  return { roleName: role?.name, orgName, userName: user?.name };
}
