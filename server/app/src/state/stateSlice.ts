import { createAsyncThunk, createSlice, PayloadAction } from "@reduxjs/toolkit";
import { addHistoryItem, HistoryItem, updateHistoryItem } from "./historySlice";
import { MessageType, SearchResponse } from "./messages";
import { AppDispatch, RootState } from "./store";
import { sendMessage } from "./socketMiddleware";

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

const sendDownload = createAsyncThunk(
  "state/send_download",
  (book: string, { dispatch }) => {
    dispatch(addInFlightDownload(book));
    dispatch(
      sendMessage({
        type: MessageType.DOWNLOAD,
        payload: { book }
      })
    );
  }
);

// Send a search to the server. Add to query history and set loading.
const sendSearch = createAsyncThunk(
  "state/send_sendSearch",
  (queryString: string, { dispatch }) => {
    // Send the books search query to the server
    dispatch(
      sendMessage({
        type: MessageType.SEARCH,
        payload: {
          query: queryString
        }
      })
    );

    const timestamp = new Date().getTime();

    // Add query to item history.
    dispatch(addHistoryItem({ query: queryString, timestamp }));
    dispatch(setActiveItem({ query: queryString, timestamp: timestamp }));
  }
);

const setSearchResults = createAsyncThunk<
  Promise<void>,
  SearchResponse,
  { dispatch: AppDispatch; state: RootState }
>(
  "state/set_search_results",
  async ({ books, errors }: SearchResponse, { dispatch, getState }) => {
    const activeItem = getState().state.activeItem;
    if (activeItem === null) {
      return;
    }
    const updatedItem: HistoryItem = {
      query: activeItem.query,
      timestamp: activeItem.timestamp,
      results: books,
      errors: errors
    };

    dispatch(setActiveItem(updatedItem));
    dispatch(updateHistoryItem(updatedItem));
  }
);

export const {
  setActiveItem,
  setConnectionState,
  setUsername,
  addInFlightDownload,
  removeInFlightDownload,
  toggleSidebar
} = stateSlice.actions;

export { stateSlice, sendMessage, sendDownload, sendSearch, setSearchResults };

export default stateSlice.reducer;
