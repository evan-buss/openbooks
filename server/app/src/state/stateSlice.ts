import { createAction, createSlice, PayloadAction } from "@reduxjs/toolkit";
import { BookDetail, MessageType } from "../models/messages";
import { addHistoryItem, HistoryItem, updateHistoryItem } from "./historySlice";
import { AppThunk } from "./store";

interface AppState {
    isConnected: boolean;
    serverFilters: string[];
    activeItem: HistoryItem | null;
}

const initialState: AppState = {
    isConnected: false,
    serverFilters: [],
    activeItem: null
}

const stateSlice = createSlice({
    name: "state",
    initialState,
    reducers: {
        toggleServerFilter(state, action: PayloadAction<string>) {
            const serverName = action.payload;
            const filterSet = new Set(state.serverFilters);
            if (filterSet.has(serverName)) {
                filterSet.delete(serverName);
            } else {
                filterSet.add(serverName);
            }
            state.serverFilters = [...filterSet];
        },
        setActiveItem(state, action: PayloadAction<HistoryItem | null>) {
            state.activeItem = action.payload;
        },
        setConnectionState(state, action: PayloadAction<boolean>) {
            state.isConnected = action.payload;
        }
    },
});

// Action that sends a websocket message to the server
const sendMessage = createAction(
    "socket/send_message",
    (message: any) => ({
        payload: { message: JSON.stringify(message) }
    })
);

// Send a search to the server. Add to query history and set loading.
const sendSearch = (queryString: string): AppThunk => dispatch => {
    // Send the books search query to the server
    dispatch(sendMessage({
        type: MessageType.SEARCH,
        payload: {
            query: queryString
        }
    }));

    const timestamp = new Date().getTime();

    // Add query to item history.
    dispatch(addHistoryItem({ query: queryString, timestamp }));
    dispatch(setActiveItem({ query: queryString, timestamp: timestamp }));
}

const setSearchResults = (results: BookDetail[]): AppThunk => (dispatch, getStore) => {
    const activeItem = getStore().state.activeItem;
    if (activeItem === null) {
        return;
    }
    const updatedItem: HistoryItem = {
        query: activeItem.query,
        timestamp: activeItem.timestamp,
        results: results
    }

    dispatch(setActiveItem(updatedItem));
    dispatch(updateHistoryItem(updatedItem))
}
const { toggleServerFilter, setActiveItem, setConnectionState } = stateSlice.actions;

export {
    stateSlice,
    sendMessage,
    toggleServerFilter,
    setActiveItem,
    setConnectionState,
    sendSearch,
    setSearchResults,
};

export default stateSlice.reducer;
