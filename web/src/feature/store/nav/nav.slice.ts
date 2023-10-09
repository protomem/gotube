import { PayloadAction, createSlice } from "@reduxjs/toolkit";

export enum NavItem {
  New = "New",
  Popular = "Popular",
  Subscriptions = "Subscriptions",
}

export interface NavState {
  selectedItem: NavItem | null;
}

const initialState: NavState = {
  selectedItem: NavItem.New,
};

const navSlice = createSlice({
  name: "nav",
  initialState,
  reducers: {
    setSelectedItem: (state, { payload }: PayloadAction<NavItem | null>) => {
      state.selectedItem = payload;
    },
  },
});

export const navReducer = navSlice.reducer;
export const navActions = navSlice.actions;
