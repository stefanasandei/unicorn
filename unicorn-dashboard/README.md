# Unicorn Dashboard

A comprehensive Next.js TypeScript dashboard for managing cloud resources through the Unicorn API. Built with modern web technologies including Tailwind CSS and Shadcn/ui components.

## Features

### 🔐 Authentication & Authorization
- **Login/Logout**: Secure authentication with JWT tokens
- **Onboarding**: Complete setup flow for new organizations
- **Role-based Access Control**: Manage users and permissions
- **Organization Management**: Multi-tenant support

### 🛡️ Identity & Access Management (IAM)
- **User Management**: Create and manage user accounts
- **Role Management**: Define roles with granular permissions
- **Organization Setup**: Configure organizational structure
- **Permission System**: Read, Write, Delete permissions

### 🔒 Secrets Manager
- **Encrypted Storage**: Secure secret management
- **CRUD Operations**: Create, read, update, delete secrets
- **Metadata Support**: Add custom metadata to secrets
- **Copy to Clipboard**: Easy secret retrieval

### 📁 Storage
- **Bucket Management**: Create and manage storage buckets
- **File Operations**: Upload, download, and delete files
- **File Browser**: Navigate through bucket contents
- **File Metadata**: View file information and properties

### 🖥️ Compute
- **Container Management**: Deploy and manage Docker containers
- **Status Monitoring**: Real-time container status
- **Quick Deploy Templates**: Pre-configured container templates
- **Environment Variables**: Configure container environments

### ⚡ Lambda
- **Serverless Functions**: Execute code in various runtimes
- **Code Editor**: Built-in function editor with syntax highlighting
- **Templates**: Pre-built function templates
- **Execution Metrics**: Monitor performance and resource usage
- **Multiple Runtimes**: Node.js, Python, Go, Java support

## Tech Stack

- **Framework**: Next.js 15 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **UI Components**: Shadcn/ui
- **Icons**: Lucide React
- **HTTP Client**: Axios
- **Form Handling**: React Hook Form
- **Validation**: Zod
- **State Management**: React Context API

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn
- Unicorn API running locally (default: http://localhost:8080)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd unicorn-dashboard
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env.local
   ```
   
   Edit `.env.local`:
   ```env
   NEXT_PUBLIC_API_URL=http://localhost:8080
   ```

4. **Run the development server**
   ```bash
   npm run dev
   ```

5. **Open your browser**
   Navigate to [http://localhost:3000](http://localhost:3000)

## Project Structure

```
src/
├── app/                    # Next.js app router pages
│   ├── dashboard/         # Dashboard page
│   ├── iam/              # IAM management page
│   ├── secrets/          # Secrets manager page
│   ├── storage/          # Storage management page
│   ├── compute/          # Compute containers page
│   ├── lambda/           # Lambda functions page
│   ├── login/            # Login page
│   └── onboarding/       # Onboarding flow
├── components/            # Reusable components
│   ├── ui/               # Shadcn/ui components
│   └── Layout.tsx        # Main layout component
├── contexts/             # React contexts
│   └── AuthContext.tsx   # Authentication context
├── lib/                  # Utility libraries
│   ├── api.ts            # API client
│   └── utils.ts          # Utility functions
└── types/                # TypeScript type definitions
    └── api.ts            # API type definitions
```

## API Integration

The dashboard integrates with the Unicorn API through a centralized API client (`src/lib/api.ts`). The client handles:

- **Authentication**: Automatic token management
- **Error Handling**: Centralized error processing
- **Request/Response Interceptors**: Automatic token injection
- **Type Safety**: Full TypeScript support

### Available Endpoints

- **Authentication**: `/api/v1/login`, `/api/v1/token/refresh`
- **IAM**: `/api/v1/roles`, `/api/v1/organizations`
- **Secrets**: `/api/v1/secrets`
- **Storage**: `/api/v1/buckets`
- **Compute**: `/api/v1/compute`
- **Lambda**: `/api/v1/lambda`

## Usage

### First Time Setup

1. **Start the Unicorn API**
   ```bash
   # Navigate to unicorn-api directory
   cd ../unicorn-api
   go run cmd/main.go
   ```

2. **Access the Dashboard**
   - Open [http://localhost:3000](http://localhost:3000)
   - You'll be redirected to the onboarding page

3. **Complete Onboarding**
   - Create your organization
   - Set up the admin role
   - Create the admin user account

4. **Login and Explore**
   - Use your admin credentials to log in
   - Explore all the available services

### Managing Resources

#### IAM Management
- Create roles with specific permissions
- Add users to your organization
- Assign roles to users
- View organization details

#### Secrets Management
- Create encrypted secrets
- View secret metadata
- Copy secret values securely
- Update secret values and metadata

#### Storage Management
- Create storage buckets
- Upload files to buckets
- Download files from buckets
- Manage file metadata

#### Compute Management
- Deploy Docker containers
- Monitor container status
- Use quick deploy templates
- Configure container environments

#### Lambda Functions
- Write serverless functions
- Choose from multiple runtimes
- Use pre-built templates
- Monitor execution metrics

## Development

### Adding New Features

1. **Create API types** in `src/types/api.ts`
2. **Add API methods** in `src/lib/api.ts`
3. **Create UI components** in `src/components/`
4. **Add pages** in `src/app/`
5. **Update navigation** in `src/components/Layout.tsx`

### Styling Guidelines

- Use Tailwind CSS classes for styling
- Leverage Shadcn/ui components for consistency
- Follow the existing design patterns
- Use Lucide React icons

### State Management

- Use React Context for global state (auth, user data)
- Use local state for component-specific data
- Implement proper loading and error states

## Deployment

### Build for Production

```bash
npm run build
npm start
```

### Environment Variables

- `NEXT_PUBLIC_API_URL`: URL of the Unicorn API
- `NEXT_PUBLIC_APP_NAME`: Application name (optional)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For support and questions:
- Check the documentation
- Review the API documentation
- Open an issue on GitHub

---

Built with ❤️ using Next.js, TypeScript, and Tailwind CSS
