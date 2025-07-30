import { api, HydrateClient } from "@/trpc/server";
import { Button } from "@/components/ui/button";

export default async function Home() {
  const hello = await api.post.hello({ text: "from tRPC" });

  return (
    <HydrateClient>
      <Button>{JSON.stringify(hello)}</Button>
    </HydrateClient>
  );
}
