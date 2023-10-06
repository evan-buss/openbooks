import {
  createAction,
  Middleware,
  MiddlewareAPI,
  PayloadAction
} from "@reduxjs/toolkit";
import { openbooksApi } from "./api";
import { deleteHistoryItem } from "./historySlice";
import {
  ConnectionResponse,
  DownloadResponse,
  MessageType,
  Notification,
  NotificationType,
  RequestMessage,
  Response,
  SearchResponse
} from "./messages";
import { addNotification } from "./notificationSlice";
import {
  removeInFlightDownload,
  setConnectionState,
  setSearchResults,
  setUsername
} from "./stateSlice";
import { AppDispatch, RootState } from "./store";
import { displayNotification, downloadFile } from "./util";
import { connectToServer } from "./connectionSlice";

// Action to send a websocket message to the server
export const sendMessage = createAction<RequestMessage>("socket/send_message");

// Web socket redux middleware.
// Listens to socket and dispatches handlers.
// Handles send_message actions by sending to socket.
export const websocketConn =
  (wsUrl: string): Middleware =>
  ({ dispatch }: MiddlewareAPI<AppDispatch, RootState>) => {
    const socket = new WebSocket(wsUrl);

    socket.onopen = () => onOpen(dispatch);
    socket.onclose = () => onClose(dispatch);
    socket.onmessage = (message) => route(dispatch, message);
    socket.onerror = () =>
      displayNotification({
        appearance: NotificationType.DANGER,
        title: "Unable to connect to server.",
        timestamp: new Date().getTime()
      });

    return (next: AppDispatch) => (action: PayloadAction<RequestMessage>) => {
      // Send Message action? Send data to the socket.
      if (sendMessage.match(action)) {
        if (socket.readyState === socket.OPEN) {
          socket.send(JSON.stringify(action.payload));
        } else {
          displayNotification({
            appearance: NotificationType.WARNING,
            title: "Server connection closed. Reload page.",
            timestamp: new Date().getTime()
          });
        }
      }

      return next(action);
    };
  };

const onOpen = (dispatch: AppDispatch): void => {
  console.log("WebSocket connected.");
  dispatch(setConnectionState(true));
  dispatch(connectToServer());
};

const onClose = (dispatch: AppDispatch): void => {
  console.log("WebSocket closed.");
  dispatch(setConnectionState(false));
};

const route = (dispatch: AppDispatch, msg: MessageEvent<string>): void => {
  const response = JSON.parse(msg.data) as Response;
  const timestamp = new Date().getTime();
  let notification: Notification = {
    ...response,
    timestamp
  };

  switch (response.type) {
    case MessageType.STATUS:
      break;
    case MessageType.CONNECT:
      dispatch(setUsername((response as ConnectionResponse).name));
      break;
    case MessageType.SEARCH:
      dispatch(setSearchResults(response as SearchResponse));
      break;
    case MessageType.DOWNLOAD:
      downloadFile((response as DownloadResponse)?.downloadPath);
      dispatch(openbooksApi.util.invalidateTags(["books"]));
      dispatch(removeInFlightDownload());
      break;
    case MessageType.RATELIMIT:
      dispatch(deleteHistoryItem());
      break;
    default:
      console.error(response);
      notification = {
        appearance: NotificationType.DANGER,
        title: "Unknown message type. See console.",
        timestamp
      };
  }

  dispatch(addNotification(notification));
  displayNotification(notification);
};
