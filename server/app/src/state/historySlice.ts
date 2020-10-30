import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { BookDetail } from "../models/messages";
import { RootState } from "./store";

// HistoryItem represents a single search history item
export type HistoryItem = {
    query: string;
    timestamp: number;
    results?: BookDetail[];
}

interface HistoryState {
    items: HistoryItem[];
}

const loadState = (): HistoryItem[] => {
    try {
        console.log('loading state');
        const serialized = localStorage.getItem('history');
        console.log(serialized);
        return JSON.parse(serialized!) ?? [];
    } catch (err) {
        return [];
    }
}

const initialState: HistoryState = {
    items: loadState()
}

export const historySlice = createSlice({
    name: 'history',
    initialState,
    reducers: {
        addHistoryItem: (state, action: PayloadAction<HistoryItem>) => {
            state.items = [action.payload, ...state.items].slice(0, 16);
        },
        deleteHistoryItem: (state, action: PayloadAction<number>) => {
            state.items = state.items.filter(x => x.timestamp !== action.payload);
        }
    },
})

export const { addHistoryItem, deleteHistoryItem } = historySlice.actions;

export const selectHistory = (state: RootState) => state.history.items;

export default historySlice.reducer;
