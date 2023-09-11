import React from "react";
import { Provider } from "react-redux";
import { store } from "@/feature/store/store";

interface StoreProviderProps {
  children: React.ReactNode;
}

export default function StoreProvider({ children }: StoreProviderProps) {
  return <Provider store={store}>{children}</Provider>;
}
