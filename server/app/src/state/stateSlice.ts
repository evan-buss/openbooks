import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { api } from "./api";

export interface AppState {
  isConnected: boolean;
  isSidebarOpen: boolean;
  username?: string;
  inFlightDownloads: string[];
}

const initialState: AppState = {
  isConnected: false,
  isSidebarOpen: true,
  username: undefined,
  inFlightDownloads: []
};

export const stateSlice = createSlice({
  name: "state",
  initialState,
  reducers: {
    setConnectionState(state, action: PayloadAction<boolean>) {
      state.isConnected = action.payload;
    },
    setUsername(state, action: PayloadAction<string>) {
      state.username = action.payload;
    },
    removeInFlightDownload(state) {
      state.inFlightDownloads.shift();
    },
    toggleSidebar(state) {
      state.isSidebarOpen = !state.isSidebarOpen;
    }
  },
  extraReducers: (builder) => {
    builder.addMatcher(
      api.endpoints.download.matchFulfilled,
      (state, action) => {
        const book = action.meta.arg.originalArgs;
        state.inFlightDownloads.push(book);
      }
    );
  }
});

export const {
  setConnectionState,
  setUsername,
  removeInFlightDownload,
  toggleSidebar
} = stateSlice.actions;

export default stateSlice.reducer;
