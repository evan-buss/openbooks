import { Action, configureStore, ThunkAction } from "@reduxjs/toolkit";
import throttle from "lodash/throttle";
import historyReducer from "./historySlice";
import { websocketConn } from "./socketMiddleware";
import serverReducer from "./serverSlice";
import stateReducer from "./stateSlice";
import { enableMapSet } from "immer";

const websocketURL = new URL(window.location.href + "ws")
if (websocketURL.protocol.startsWith("https")) {
    websocketURL.protocol = websocketURL.protocol.replace("https", "wss");
} else {
    websocketURL.protocol = websocketURL.protocol.replace("http", "ws");
}

enableMapSet();

export const store = configureStore({
    reducer: {
        state: stateReducer,
        history: historyReducer,
        servers: serverReducer
    },
    middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(websocketConn(websocketURL.href))
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
