import { Action, configureStore, ThunkAction } from "@reduxjs/toolkit";
import historyReducer from "./historySlice";
import serverReducer from "./serverSlice";
import throttle from "lodash/throttle";

export const store = configureStore({
    reducer: {
        history: historyReducer,
        servers: serverReducer
    },
});


const saveState = (key: string, state: any): void => {
    try {
        const serialized = JSON.stringify(state);
        localStorage.setItem(key, serialized);
    } catch (err) { }
}

store.subscribe(throttle(() => {
    saveState("history", store.getState().history.items);
}, 1000));

export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<
    ReturnType,
    RootState,
    unknown,
    Action<string>
>;

