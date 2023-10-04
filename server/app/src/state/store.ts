import { configureStore } from "@reduxjs/toolkit";
import { setupListeners } from "@reduxjs/toolkit/query/react";
import { enableMapSet } from "immer";
import throttle from "lodash/throttle";
import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import { openbooksApi } from "./api";
import historyReducer from "./historySlice";
import notificationReducer from "./notificationSlice";
import { websocketConn } from "./socketMiddleware";
import stateReducer from "./stateSlice";
import { getWebsocketURL } from "./util";

enableMapSet();

export const store = configureStore({
  reducer: {
    state: stateReducer,
    history: historyReducer,
    notifications: notificationReducer,
    [openbooksApi.reducerPath]: openbooksApi.reducer
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(
      websocketConn(getWebsocketURL().href),
      openbooksApi.middleware
    )
});

setupListeners(store.dispatch);

const saveState = (key: string, state: unknown): void => {
  try {
    const serialized = JSON.stringify(state);
    localStorage.setItem(key, serialized);
  } catch (err) {
    console.error("Error saving state:", err);
  }
};

store.subscribe(
  throttle(() => {
    saveState("history", store.getState().history.items);
    saveState("active", store.getState().state.activeItem);
  }, 1000)
);

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
export const useAppDispatch: () => AppDispatch = useDispatch;
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;
