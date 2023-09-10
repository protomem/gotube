import { QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import { queryClient } from "./query";

interface QueryProviderProps {
  children: React.ReactNode;
}

export default function QueryProvider({ children }: QueryProviderProps) {
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}
