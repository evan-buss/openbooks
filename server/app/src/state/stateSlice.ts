import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { HistoryItem } from "./historySlice";

interface AppState {
  // TODO: Need to differentiate websocket connection status from IRC connection status
  isConnected: boolean;
  isSidebarOpen: boolean;
  activeItem: HistoryItem | null;
  username?: string;
  inFlightDownloads: string[];
}

const loadActive = (): HistoryItem | null => {
  try {
    return JSON.parse(localStorage.getItem("active")!) ?? null;
  } catch (err) {
    return null;
  }
};

const initialState: AppState = {
  isConnected: false,
  isSidebarOpen: true,
  activeItem: loadActive(),
  username: undefined,
  inFlightDownloads: []
};

const stateSlice = createSlice({
  name: "state",
  initialState,
  reducers: {
    setActiveItem(state, action: PayloadAction<HistoryItem | null>) {
      state.activeItem = action.payload;
    },
    setConnectionState(state, action: PayloadAction<boolean>) {
      state.isConnected = action.payload;
    },
    setUsername(state, action: PayloadAction<string>) {
      state.username = action.payload;
    },
    addInFlightDownload(state, action: PayloadAction<string>) {
      state.inFlightDownloads.push(action.payload);
    },
    removeInFlightDownload(state) {
      state.inFlightDownloads.shift();
    },
    toggleSidebar(state) {
      state.isSidebarOpen = !state.isSidebarOpen;
    }
  }
});

export const {
  setActiveItem,
  setConnectionState,
  setUsername,
  addInFlightDownload,
  removeInFlightDownload,
  toggleSidebar
} = stateSlice.actions;

export default stateSlice.reducer;
