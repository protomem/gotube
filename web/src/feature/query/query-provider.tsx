import { QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import { queryClient } from "@/feature/query/query";

interface QueryProviderProps {
  children: React.ReactNode;
}

export default function QueryProvider({ children }: QueryProviderProps) {
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}
