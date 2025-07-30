# Minimal Onboarding Web UI

This project is a minimal Vite + React + TypeScript web UI for onboarding with a Go API (IAM mock service, Swagger schema).

## Features

- Onboarding flow: create organization and admin user (with all permissions)
- Admin login using cookies
- Clean, functional, and maintainable code

## Getting Started

1. Install dependencies:
   ```sh
   npm install
   ```
2. Start the development server:
   ```sh
   npm run dev
   ```
3. Open [http://localhost:5173](http://localhost:5173) in your browser.

## Project Structure

- `src/` — Main source code
- `src/App.tsx` — Main app logic (onboarding, login)

## Customization

- Connect to your Go API by updating the API endpoints in the code.

---

This project was bootstrapped with Vite, React, and TypeScript.
