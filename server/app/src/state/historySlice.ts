import {
  createAsyncThunk,
  createSelector,
  createSlice,
  PayloadAction
} from "@reduxjs/toolkit";
import { BookDetail, ParseError } from "./messages";
import { setActiveItem } from "./stateSlice";
import { AppDispatch, RootState } from "./store";
import { selectActiveIrcServer } from "./connectionSlice";

// HistoryItem represents a single search history item
export interface HistoryItem {
  serverName: string;
  query: string;
  timestamp: number;
  results?: BookDetail[];
  errors?: ParseError[];
}

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
      const pendingItemIndex = state.items.findIndex(
        (x) => x.timestamp === action.payload.timestamp
      );
      state.items = [
        ...state.items.slice(0, pendingItemIndex),
        action.payload,
        ...state.items.slice(pendingItemIndex + 1)
      ];
    },
    setHistory: (state, action: PayloadAction<HistoryItem[]>) => {
      state.items = action.payload;
    }
  }
});

// Delete an item from history. Clear current item and loading state if deleting active search
const deleteHistoryItem = createAsyncThunk<
  Promise<void>,
  number | undefined,
  { dispatch: AppDispatch; state: RootState }
>("history/delete_item", async (timeStamp, { dispatch, getState }) => {
  if (timeStamp === undefined) {
    dispatch(setActiveItem(null));
    const toRemove = getState().history.items.at(0)?.timestamp;
    if (toRemove) {
      dispatch(historySlice.actions.deleteByTimetamp(toRemove));
    }
    return;
  }

  const activeItem = getState().state.activeItem;
  if (activeItem?.timestamp === timeStamp) {
    dispatch(setActiveItem(null));
  }

  dispatch(historySlice.actions.deleteByTimetamp(timeStamp));
});

const selectAll = (state: RootState) => state.history.items;
const selectHistory = createSelector(
  [selectAll, selectActiveIrcServer],
  (history, server) => history.filter((x) => x.serverName === server.name)
);

export const { addHistoryItem, updateHistoryItem, setHistory } =
  historySlice.actions;
export { deleteHistoryItem, selectHistory };
export default historySlice.reducer;
