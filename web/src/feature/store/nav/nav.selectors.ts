import { RootState } from "@/feature/store/store";

const selectNavModule = (state: RootState) => state.nav;

export const selectSelectedItem = (state: RootState) =>
  selectNavModule(state).selectedItem;
