import { useEffect, useState } from "react";
import { getAccountInfo } from "../api/user";
import { useAuth } from "../hooks/useAuth";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Icons } from "@/components/ui/icons";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";

// Assuming these components are available or created based on shadcn/ui
// You might need to create them or adjust imports if they are custom.
// For example, NavLink would be a simple anchor tag or a router link.
const NavLink = ({ children, to, className = "", onClick }: any) => (
  <a
    href={to}
    onClick={onClick}
    className={`flex items-center gap-3 rounded-lg px-3 py-2 text-muted-foreground transition-all hover:text-primary ${className}`}
  >
    {children}
  </a>
);

export default function DashboardPage() {
  const { token, logout } = useAuth();
  const [info, setInfo] = useState<{
    roleName?: string;
    orgName?: string;
    permissions?: number[];
    name?: string;
    email?: string;
  } | null>(null);
  const [error, setError] = useState("");
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    if (!token) return;
    getAccountInfo(token)
      .then(setInfo)
      .catch(() => setError("Failed to load account info"));
  }, [token]);

  function handleCopyToken() {
    if (!token) return;
    navigator.clipboard.writeText(token);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  }

  return (
    <div className="grid min-h-screen w-full md:grid-cols-[220px_1fr] lg:grid-cols-[280px_1fr]">
      {/* Desktop Sidebar */}
      <div className="hidden border-r bg-muted/40 md:block">
        <div className="flex h-full max-h-screen flex-col gap-2">
          <div className="flex h-14 items-center border-b px-4 lg:h-[60px] lg:px-6">
            <a href="/" className="flex items-center gap-2 font-semibold">
              <Icons.eye className="h-6 w-6" />
              <span>Unicorn Admin</span>
            </a>
          </div>
          <div className="flex-1">
            <nav className="grid items-start px-2 text-sm font-medium lg:px-4">
              <NavLink to="#" className="text-primary">
                <Icons.eye className="h-4 w-4" />
                Dashboard
              </NavLink>
              {/* Add more navigation links here */}
            </nav>
          </div>
          <div className="mt-auto p-4 border-t">
            <Card>
              <CardHeader className="p-2 pt-0 md:p-4">
                <CardTitle>Account</CardTitle>
                <CardDescription>
                  View your account details and manage settings.
                </CardDescription>
              </CardHeader>
              <CardContent className="p-2 pt-0 md:p-4 md:pt-0">
                <div className="flex items-center gap-2 mb-2">
                  <Icons.user className="h-5 w-5 text-muted-foreground" />
                  <span className="font-semibold text-sm">
                    {info?.name || "User"}
                  </span>
                </div>
                <div className="flex items-center gap-2">
                  <Icons.mail className="h-4 w-4 text-muted-foreground" />
                  <span className="text-xs text-muted-foreground">
                    {info?.email}
                  </span>
                </div>
                <Separator className="my-3" />
                <Button
                  size="sm"
                  variant="outline"
                  className="w-full"
                  onClick={logout}
                >
                  <Icons.logout className="mr-2 h-4 w-4" /> Logout
                </Button>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex flex-col">
        <header className="flex h-14 items-center gap-4 border-b bg-muted/40 px-4 lg:h-[60px] lg:px-6">
          {/* Mobile Sidebar Trigger */}
          <Sheet>
            <SheetTrigger asChild>
              <Button
                variant="outline"
                size="icon"
                className="shrink-0 md:hidden"
              >
                <Icons.eye className="h-5 w-5" />
                <span className="sr-only">Toggle navigation menu</span>
              </Button>
            </SheetTrigger>
            <SheetContent side="left" className="flex flex-col">
              <nav className="grid gap-2 text-lg font-medium">
                <a
                  href="#"
                  className="flex items-center gap-2 text-lg font-semibold"
                >
                  <Icons.eye className="h-6 w-6" />
                  <span>Unicorn Admin</span>
                </a>
                <NavLink to="#" className="mx-[-0.65rem] text-primary">
                  <Icons.eye className="h-5 w-5" />
                  Dashboard
                </NavLink>
                {/* Add more mobile navigation links here */}
              </nav>
              <div className="mt-auto pt-4 border-t">
                <p className="text-sm text-muted-foreground mb-2">Account</p>
                <div className="flex items-center gap-2 mb-1">
                  <Icons.user className="h-4 w-4 text-muted-foreground" />
                  <span className="font-semibold text-sm">
                    {info?.name || "User"}
                  </span>
                </div>
                <div className="flex items-center gap-2">
                  <Icons.mail className="h-3.5 w-3.5 text-muted-foreground" />
                  <span className="text-xs text-muted-foreground">
                    {info?.email}
                  </span>
                </div>
                <Separator className="my-3" />
                <Button
                  size="sm"
                  variant="outline"
                  className="w-full"
                  onClick={logout}
                >
                  <Icons.logout className="mr-2 h-4 w-4" /> Logout
                </Button>
              </div>
            </SheetContent>
          </Sheet>

          {/* Page Title for Main Content */}
          <h1 className="text-xl font-semibold ml-auto md:ml-0">Dashboard</h1>
          {/* Optional: Add user dropdown/avatar here for more actions */}
          <div className="ml-auto flex items-center gap-2 md:hidden">
            <Icons.user className="h-6 w-6 text-muted-foreground" />
            <span className="font-medium text-sm">{info?.name || "User"}</span>
          </div>
        </header>

        <main className="flex flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6">
          <div className="flex items-center">
            <h2 className="text-lg font-semibold md:text-2xl">
              Welcome{info?.name ? `, ${info.name}` : ""}!
            </h2>
          </div>
          <div className="grid flex-1 items-start gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {error ? (
              <Alert
                variant="destructive"
                className="md:col-span-2 lg:col-span-3 xl:col-span-4"
              >
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            ) : (
              <>
                <Card className="col-span-1 md:col-span-2 lg:col-span-2">
                  <CardHeader>
                    <CardTitle>Account Information</CardTitle>
                    <CardDescription>
                      Details about your organization and role.
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div>
                      <h3 className="text-sm font-medium text-muted-foreground">
                        Organization
                      </h3>
                      <p className="text-xl font-semibold text-foreground">
                        {info?.orgName || "N/A"}
                      </p>
                    </div>
                    <div>
                      <h3 className="text-sm font-medium text-muted-foreground">
                        Role
                      </h3>
                      <p className="text-xl font-semibold text-foreground">
                        {info?.roleName?.split(":")[0] || "N/A"}:{" "}
                        {info?.permissions
                          ? JSON.stringify(info?.permissions)
                          : "N/A"}
                      </p>
                    </div>
                  </CardContent>
                </Card>

                <Card className="col-span-1 md:col-span-2 lg:col-span-2">
                  <CardHeader>
                    <CardTitle>Authentication Token</CardTitle>
                    <CardDescription>
                      Use this token for API access. Keep it secure.
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="relative flex items-center gap-2 rounded-md border bg-muted p-3 text-sm font-mono break-all">
                      <Icons.key className="h-4 w-4 text-muted-foreground" />
                      <code className="flex-1 select-all">
                        {token || "N/A"}
                      </code>
                      <Button
                        variant="secondary"
                        size="sm"
                        onClick={handleCopyToken}
                        className="ml-2 shrink-0"
                      >
                        {copied ? "Copied!" : "Copy"}
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              </>
            )}
          </div>
        </main>
      </div>
    </div>
  );
}
