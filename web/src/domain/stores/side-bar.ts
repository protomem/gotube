import { create } from "zustand";

interface SideBarStore {
  isOpen: boolean;
  open: () => void;
  close: () => void;
}

export const useSideBarStore = create<SideBarStore>((set) => ({
  isOpen: true,
  open: () => set({ isOpen: true }),
  close: () => set({ isOpen: false }),
}));
