import React from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const client = new QueryClient();

type Props = {
  children: React.ReactNode;
};

const QueryProvider = ({ children }: Props) => {
  return <QueryClientProvider client={client}>{children}</QueryClientProvider>;
};

export default QueryProvider;
