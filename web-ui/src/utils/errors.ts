export function getErrorMessage(err: unknown, fallback: string) {
  if (typeof err === "object" && err && "response" in err) {
    const response = (err as { response?: { data?: { error?: string } } })
      .response;
    if (response && response.data && typeof response.data.error === "string") {
      return response.data.error;
    }
  }
  return fallback;
}
