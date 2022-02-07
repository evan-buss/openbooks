import { AnyAction, PayloadAction, Store } from "@reduxjs/toolkit";
import {
  BookResponse,
  MessageType,
  NotificationType,
  SearchResponse
} from "./messages";
import { displayNotification, downloadFile } from "./util";
import {
  sendMessage,
  setConnectionState,
  setSearchResults,
  setUsername
} from "./stateSlice";
import { addNotification } from "./notificationSlice";
import { openbooksApi } from "./api";
import { deleteHistoryItem } from "./historySlice";

// Web socket redux middleware.
// Listens to socket and dispatches handlers.
// Handles send_message actions by sending to socket.
export const websocketConn = (wsUrl: string): any => {
  return (store: Store) => {
    const socket = new WebSocket(wsUrl);

    socket.onopen = () => onOpen(store);
    socket.onclose = () => onClose(store);
    socket.onmessage = (message) => route(store, message);
    socket.onerror = (event) => console.error(event);

    return (next: any) => (action: PayloadAction<any>) => {
      // Send Message action? Send data to the socket.
      if (sendMessage.match(action)) {
        if (socket.readyState === socket.OPEN) {
          socket.send(action.payload.message);
        } else {
          displayNotification({
            type: NotificationType.WARNING,
            title: "Server connection closed. Reload page.",
            timestamp: new Date().getTime()
          });
        }
      }

      return next(action);
    };
  };
};

const onOpen = (store: Store): void => {
  console.log("WebSocket connected.");
  store.dispatch(setConnectionState(true));
  store.dispatch(sendMessage({ type: MessageType.CONNECT, payload: {} }));
};

const onClose = (store: Store): void => {
  console.log("WebSocket closed.");
  store.dispatch(setConnectionState(false));
};

const route = (store: Store, msg: MessageEvent<any>): void => {
  const getNotif = (): Notification => {
    let response = JSON.parse(msg.data) as BookResponse;
    const timestamp = new Date().getTime();
    const notification: Notification = {
      type: response.NotificationType,
      title: response.title,
      detail: response.detail,
      timestamp
    };
    switch (response.type) {
      case MessageType.STATUS:
        return notification;
      case MessageType.CONNECT:
        store.dispatch(setUsername(response.name));
        return notification;
      case MessageType.SEARCH:
        store.dispatch(
          setSearchResults(
            response as unknown as SearchResponse
          ) as unknown as AnyAction
        );
        return notification;
      case MessageType.DOWNLOAD:
        downloadFile(response.downloadLink);
        store.dispatch(openbooksApi.util.invalidateTags(["books"]));
        return notification;
      case MessageType.RATELIMIT:
        store.dispatch(deleteHistoryItem() as unknown as AnyAction);
        return notification;
      default:
        console.error(response);
        return {
          type: NotificationType.DANGER,
          title: "Unknown message type. See console.",
          timestamp
        };
    }
  };

  const notif = getNotif();
  store.dispatch(addNotification(notif));
  displayNotification(notif);
};
