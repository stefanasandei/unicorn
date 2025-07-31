"use client";

import React from "react";
import { useAuth } from "@/contexts/AuthContext";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Settings,
  LogOut,
  User,
  Shield,
  Database,
  Server,
  Code,
  FileText,
  Home,
  Building,
  HardDrive,
  Activity,
  Sparkles,
} from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";

interface LayoutProps {
  children: React.ReactNode;
}

const navigation = [
  { name: "Monitoring", href: "/monitoring", icon: Activity },
  { name: "IAM", href: "/iam", icon: Shield },
  { name: "Secrets", href: "/secrets", icon: Database },
  { name: "Storage", href: "/storage", icon: FileText },
  { name: "Compute", href: "/compute", icon: Server },
  { name: "Databases", href: "/rdb", icon: HardDrive },
  { name: "Lambda", href: "/lambda", icon: Code },
];

export const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { user, logout } = useAuth();
  const pathname = usePathname();

  if (!user) {
    return <div>{children}</div>;
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="bg-card border-b border-border shadow-sm backdrop-blur-sm bg-card/80">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center space-x-3">
              <div className="flex items-center space-x-2">
                <div className="p-1.5 rounded-lg bg-gradient-to-br from-primary to-primary/80">
                  <Sparkles className="h-5 w-5 text-primary-foreground" />
                </div>
                <h1 className="text-xl font-bold bg-gradient-to-r from-foreground to-foreground/80 bg-clip-text text-transparent">
                  Unicorn Dashboard
                </h1>
              </div>
            </div>

            <div className="flex items-center space-x-4">
              <Badge variant="secondary" className="bg-secondary/50 text-secondary-foreground border-secondary/20">
                {user.organization.name}
              </Badge>

              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    className="relative h-8 w-8 rounded-full hover:bg-accent/50"
                  >
                    <Avatar className="h-8 w-8 ring-2 ring-primary/20">
                      <AvatarImage src="/avatars/01.png" alt={user.name} />
                      <AvatarFallback className="bg-primary/10 text-primary font-medium">
                        {user.name.charAt(0).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent className="w-56" align="end" forceMount>
                  <DropdownMenuLabel className="font-normal">
                    <div className="flex flex-col space-y-1">
                      <p className="text-sm font-medium leading-none text-foreground">
                        {user.name}
                      </p>
                      <p className="text-xs leading-none text-muted-foreground">
                        {user.email}
                      </p>
                    </div>
                  </DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem className="hover:bg-accent/50">
                    <User className="mr-2 h-4 w-4" />
                    <span>Profile</span>
                  </DropdownMenuItem>
                  <DropdownMenuItem className="hover:bg-accent/50">
                    <Settings className="mr-2 h-4 w-4" />
                    <span>Settings</span>
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={logout} className="hover:bg-destructive/10 text-destructive">
                    <LogOut className="mr-2 h-4 w-4" />
                    <span>Log out</span>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </div>
      </header>

      <div className="flex">
        {/* Sidebar */}
        <aside className="w-64 bg-card/50 border-r border-border backdrop-blur-sm">
          <nav className="mt-6 px-3">
            <div className="space-y-2">
              {navigation.map((item) => {
                const isActive = pathname === item.href;
                return (
                  <Link
                    key={item.name}
                    href={item.href}
                    className={`group flex items-center px-3 py-2.5 text-sm font-medium rounded-lg transition-all duration-200 ${
                      isActive
                        ? "bg-primary/10 text-primary border border-primary/20 shadow-sm"
                        : "text-muted-foreground hover:bg-accent/50 hover:text-accent-foreground"
                    }`}
                  >
                    <item.icon
                      className={`mr-3 flex-shrink-0 h-5 w-5 transition-colors ${
                        isActive
                          ? "text-primary"
                          : "text-muted-foreground group-hover:text-accent-foreground"
                      }`}
                    />
                    {item.name}
                  </Link>
                );
              })}
            </div>
          </nav>
        </aside>

        {/* Main content */}
        <main className="flex-1 bg-background/50">
          <div className="py-8">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
              {children}
            </div>
          </div>
        </main>
      </div>
    </div>
  );
};
