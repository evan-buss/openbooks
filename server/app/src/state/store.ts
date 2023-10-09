import { configureStore } from "@reduxjs/toolkit";
import { setupListeners } from "@reduxjs/toolkit/query/react";
import { enableMapSet } from "immer";
import throttle from "lodash/throttle";
import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import { openbooksApi } from "./api";
import historyReducer, { setHistory } from "./historySlice";
import notificationReducer from "./notificationSlice";
import connectionReducer from "./connectionSlice";
import { websocketConn } from "./socketMiddleware";
import stateReducer, { setActiveItem } from "./stateSlice";
import { getWebsocketURL } from "./util";

enableMapSet();

const storageVersion = localStorage.getItem("version");

export const store = configureStore({
  reducer: {
    state: stateReducer,
    history: historyReducer,
    notifications: notificationReducer,
    connection: connectionReducer,
    [openbooksApi.reducerPath]: openbooksApi.reducer
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(
      websocketConn(getWebsocketURL().href),
      openbooksApi.middleware
    )
});

setupListeners(store.dispatch);

switch (storageVersion) {
  case "5.0.0":
    break;
  default:
    applyInitialMigration();
}

// We now need to associate each history item with a server.
// All existing history items are associated with the "IRC Highway" server.
function applyInitialMigration() {
  console.log("Applying initial migration to 5.0.0");

  const migrated = store.getState().history.items.map((item) => {
    return {
      ...item,
      serverName: "IRC Highway"
    };
  });

  store.dispatch(setHistory(migrated));
  store.dispatch(setActiveItem(null));
  localStorage.setItem("version", "5.0.0");
}

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
    saveState("connection", store.getState().connection.selectedServer);
  }, 1000)
);

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
export const useAppDispatch: () => AppDispatch = useDispatch;
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;
