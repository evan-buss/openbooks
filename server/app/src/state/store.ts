import { Action, configureStore, ThunkAction } from "@reduxjs/toolkit";
import throttle from "lodash/throttle";
import historyReducer from "./historySlice";
import { websocketConn } from "./socketMiddleware";
import serverReducer from "./serverSlice";
import stateReducer from "./stateSlice";
import { enableMapSet } from "immer";

const wsProtocol = window.location.protocol === "https:" ? "wss://" : "ws://";
const wsHost = process.env.NODE_ENV === "development" ? "localhost:5228" : window.location.host;
const wsUrl = `${wsProtocol}${wsHost}/ws`;

enableMapSet();

export const store = configureStore({
    reducer: {
        state: stateReducer,
        history: historyReducer,
        servers: serverReducer
    },
    middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(websocketConn(wsUrl))
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
