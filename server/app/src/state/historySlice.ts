import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { BookDetail, ParseError } from "./messages";
import { setActiveItem } from "./stateSlice";
import { AppThunk, RootState } from "./store";

// HistoryItem represents a single search history item
type HistoryItem = {
  query: string;
  timestamp: number;
  results?: BookDetail[];
  errors?: ParseError[];
};

interface HistoryState {
  items: HistoryItem[];
}

const loadState = (): HistoryItem[] => {
  try {
    return JSON.parse(localStorage.getItem("history")!) ?? [];
  } catch (err) {
    return [];
  }
};

const initialState: HistoryState = {
  items: loadState()
};

export const historySlice = createSlice({
  name: "history",
  initialState,
  reducers: {
    addHistoryItem: (state, action: PayloadAction<HistoryItem>) => {
      state.items = [action.payload, ...state.items].slice(0, 16);
    },
    deleteByTimetamp: (state, action: PayloadAction<number>) => {
      state.items = state.items.filter((x) => x.timestamp !== action.payload);
    },
    updateHistoryItem: (state, action: PayloadAction<HistoryItem>) => {
      var pendingItemIndex = state.items.findIndex(
        (x) => x.timestamp === action.payload.timestamp
      );
      state.items = [
        ...state.items.slice(0, pendingItemIndex),
        action.payload,
        ...state.items.slice(pendingItemIndex + 1)
      ];
    }
  }
});

// Delete an item from history. Clear current item and loading state if deleting active search
const deleteHistoryItem =
  (timeStamp?: number): AppThunk =>
  (dispatch, getStore) => {
    if (timeStamp === undefined) {
      dispatch(setActiveItem(null));
      const toRemove = getStore().history.items.at(0)?.timestamp;
      if (toRemove) {
        dispatch(historySlice.actions.deleteByTimetamp(toRemove));
      }
      return;
    }

    const activeItem = getStore().state.activeItem;
    if (activeItem?.timestamp === timeStamp) {
      dispatch(setActiveItem(null));
    }

    dispatch(historySlice.actions.deleteByTimetamp(timeStamp));
  };

const { addHistoryItem, updateHistoryItem } = historySlice.actions;

const selectHistory = (state: RootState) => state.history.items;

export type { HistoryItem };
export { deleteHistoryItem, addHistoryItem, updateHistoryItem, selectHistory };

export default historySlice.reducer;
