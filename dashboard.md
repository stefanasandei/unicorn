# Unicorn Dashboard

A professional Next.js dashboard with Tailwind CSS and shadcn/ui that implements IAM and Secrets Manager from the Unicorn API.

## Features

- **Identity and Access Management (IAM)**
  - User authentication and authorization
  - Role-based access control
  - Organization management

- **Secrets Manager**
  - Securely store and manage sensitive information
  - Create, read, update, and delete secrets
  - Metadata support for better organization

- **Modern UI**
  - Built with Next.js, Tailwind CSS, and shadcn/ui
  - Responsive design for all devices
  - Dark mode support

## Tech Stack

- **Frontend**
  - Next.js 15
  - TypeScript
  - Tailwind CSS
  - shadcn/ui components
  - React Hook Form with Zod validation

- **API Integration**
  - Axios for API requests
  - JWT authentication
  - Context API for state management

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/stefanasandei/unicorn.git
   cd unicorn/dashboard
   ```

2. Install dependencies:
   ```bash
   npm install
   # or
   yarn install
   ```

3. Create a `.env.local` file in the root directory:
   ```
   NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
   ```

4. Start the development server:
   ```bash
   npm run dev
   # or
   yarn dev
   ```

5. Open [http://localhost:12000](http://localhost:12000) in your browser.

## Implementation Details

### Authentication

The dashboard uses JWT authentication with the Unicorn API. The authentication flow is as follows:

1. User logs in with email and password
2. API returns a JWT token
3. Token is stored in localStorage
4. Token is included in all subsequent API requests
5. Token is validated on each page load
6. User is redirected to login if token is invalid or expired

### IAM Implementation

The dashboard implements the following IAM features:

- **User Management**
  - Create, read, update, and delete users
  - Assign roles to users

- **Role Management**
  - Create, read, update, and delete roles
  - Define permissions for roles

- **Permission System**
  - Read (0), Write (1), and Delete (2) permissions
  - Role-based access control for UI elements and API calls

### Secrets Manager Implementation

The dashboard implements the following Secrets Manager features:

- **Secret Management**
  - Create, read, update, and delete secrets
  - Secure storage of sensitive information
  - Metadata support for better organization

## Screenshots

(Screenshots will be added once the dashboard is deployed)

## License

This project is licensed under the MIT License - see the LICENSE file for details.