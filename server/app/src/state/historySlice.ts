import { createSelector, createSlice, PayloadAction } from "@reduxjs/toolkit";
import { api } from "./api";
import { BookDetail, ParseError, SearchResponse } from "./messages";
import { RootState } from "./store";

// HistoryItem represents a single search history item
export type HistoryItem = {
  query: string;
  timestamp: number;
  results?: BookDetail[];
  errors?: ParseError[];
};

export interface HistoryState {
  active: number | undefined;
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
  active: undefined,
  items: loadState()
};

export const historySlice = createSlice({
  name: "history",
  initialState,
  reducers: {
    setActiveItem(state, action: PayloadAction<HistoryItem | null>) {
      state.active = action.payload?.timestamp;
    },
    removeResults: (state, action: PayloadAction<number>) => {
      const timestamp = action.payload;

      if (state.active === timestamp) {
        state.active = undefined;
      }

      state.items = state.items.filter((x) => x.timestamp !== timestamp);
    },
    resultsReceived: (state, action: PayloadAction<SearchResponse>) => {
      state.items[0].results = action.payload.books;
      state.items[0].errors = action.payload.errors;
    }
  },
  extraReducers: (builder) => {
    builder.addMatcher(api.endpoints.search.matchFulfilled, (state, action) => {
      const timestamp = new Date().getTime();
      const query = action.meta.arg.originalArgs;

      state.active = timestamp;
      state.items = [{ query, timestamp }, ...state.items].slice(0, 16);
    });
  }
});

export const selectActiveItem = createSelector(
  (state: RootState) => state.history.items,
  (state: RootState) => state.history.active,
  (items, active) => items.find((x) => x.timestamp === active)
);

export const { setActiveItem, removeResults, resultsReceived } =
  historySlice.actions;

export default historySlice.reducer;
