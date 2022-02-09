import { createAction, createSlice, PayloadAction } from "@reduxjs/toolkit";
import { MessageType, SearchResponse } from "./messages";
import { addHistoryItem, HistoryItem, updateHistoryItem } from "./historySlice";
import { AppThunk } from "./store";

interface AppState {
  isConnected: boolean;
  activeItem: HistoryItem | null;
  username?: string;
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
  activeItem: loadActive(),
  username: undefined
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
    }
  }
});

// Action that sends a websocket message to the server
const sendMessage = createAction("socket/send_message", (message: any) => ({
  payload: { message: JSON.stringify(message) }
}));

// Send a search to the server. Add to query history and set loading.
const sendSearch =
  (queryString: string): AppThunk =>
  (dispatch) => {
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
  };

const setSearchResults =
  ({ books, errors }: SearchResponse): AppThunk =>
  (dispatch, getStore) => {
    const activeItem = getStore().state.activeItem;
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
  };
const { setActiveItem, setConnectionState, setUsername } = stateSlice.actions;

export {
  stateSlice,
  sendMessage,
  setActiveItem,
  setConnectionState,
  sendSearch,
  setSearchResults,
  setUsername
};

export default stateSlice.reducer;
