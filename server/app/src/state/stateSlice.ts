import { createAction, createSlice, PayloadAction } from "@reduxjs/toolkit";
import { RootState } from "./store";

interface AppState {
    isLoading: boolean;
}

const initialState: AppState = {
    isLoading: false
}

const stateSlice = createSlice({
    name: "state",
    initialState,
    reducers: {
        setLoading(state, action: PayloadAction<boolean>) {
            state.isLoading = action.payload;
        },
    },
});

const sendMessage = createAction(
    "socket/send_message",
    (message: any) => ({
        payload: { message: JSON.stringify(message) }
    })
);

const { setLoading } = stateSlice.actions;

const selectIsLoading = (state: RootState) => state.state.isLoading;

export default stateSlice.reducer;

export {
    stateSlice,
    sendMessage,
    setLoading,
    selectIsLoading
};
