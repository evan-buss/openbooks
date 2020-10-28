import { createSlice, PayloadAction } from "@reduxjs/toolkit"
import { RootState } from "./store";

interface ServerState {
    servers: string[];
}

const initialState: ServerState = {
    servers: []
}

export const serverSlice = createSlice({
    name: "servers",
    initialState,
    reducers: {
        updateServers(state, action: PayloadAction<string[]>) {
            state.servers = action.payload;
        }
    }
});

export const { updateServers } = serverSlice.actions;

export const selectServers = (state: RootState) => state.servers.servers;

export default serverSlice.reducer;