import { ReactNode, useEffect } from "react";
import { useAuthStore } from "@/domain/stores/auth";

interface Props {
  children: ReactNode;
}

export default function StoresProvider({ children }: Props) {
  const authLoadLocalData = useAuthStore((state) => state.loadLocalData);
  useEffect(() => {
    authLoadLocalData();
  }, [authLoadLocalData]);

  return <>{children}</>;
}
