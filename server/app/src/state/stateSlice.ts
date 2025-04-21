import {
  createAction,
  createAsyncThunk,
  createSlice,
  PayloadAction
} from "@reduxjs/toolkit";
import { addHistoryItem, HistoryItem, updateHistoryItem } from "./historySlice";
import { MessageType, SearchResponse } from "./messages";
import { AppDispatch, RootState } from "./store";

interface AppState {
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
    console.error("Failed to load active item from localStorage:", err);
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

// Action that sends a websocket message to the server
const sendMessage = createAction("socket/send_message", (message: any) => ({
  payload: { message: JSON.stringify(message) }
}));

const sendDownload = createAsyncThunk(
  "state/send_download",
  (
    payload: { book: string; author?: string; title?: string },
    { dispatch }
  ) => {
    dispatch(addInFlightDownload(payload.book));
    dispatch(
      sendMessage({
        type: MessageType.DOWNLOAD,
        payload: {
          book: payload.book,
          author: payload.author,
          title: payload.title,
        },
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
