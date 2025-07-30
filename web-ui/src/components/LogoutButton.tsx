import React from "react";
import { Button } from "@/components/ui/button";
import { Icons } from "@/components/ui/icons";

export default function LogoutButton({ onLogout }: { onLogout: () => void }) {
  return (
    <Button variant="ghost" size="sm" onClick={onLogout} className="gap-2">
      <Icons.logout className="h-4 w-4" />
      Logout
    </Button>
  );
}
