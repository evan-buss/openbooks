import { Action, configureStore, ThunkAction } from "@reduxjs/toolkit";
import throttle from "lodash/throttle";
import historyReducer from "./historySlice";
import notificationReducer from "./notificationSlice";
import { websocketConn } from "./socketMiddleware";
import stateReducer from "./stateSlice";
import { enableMapSet } from "immer";
import { getWebsocketURL } from "./util";
import { openbooksApi } from "./api";
import { setupListeners } from "@reduxjs/toolkit/dist/query";
import { useDispatch } from "react-redux";

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

const saveState = (key: string, state: any): void => {
  try {
    const serialized = JSON.stringify(state);
    localStorage.setItem(key, serialized);
  } catch (err) {}
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
