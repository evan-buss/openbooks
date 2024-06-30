import { Middleware, ThunkDispatch, UnknownAction } from "@reduxjs/toolkit";
import { api } from "./api";
import { resultsReceived } from "./historySlice";
import {
  ConnectionResponse,
  DownloadResponse,
  MessageType,
  Notification,
  NotificationType,
  Response,
  SearchResponse
} from "./messages";
import { addNotification } from "./notificationSlice";
import {
  removeInFlightDownload,
  setConnectionState,
  setUsername
} from "./stateSlice";
import { displayNotification, downloadFile } from "./util";

// SSE redux middleware.
// Listens to events and dispatches handlers.
export const sseMiddleware =
  (
    sseUrl: string
  ): Middleware<
    Record<string, never>,
    unknown,
    ThunkDispatch<unknown, unknown, UnknownAction>
  > =>
  ({ dispatch }) => {
    const eventSource = new EventSource(sseUrl, {
      withCredentials: true
    });

    eventSource.onopen = () => onOpen(dispatch);
    eventSource.onmessage = (message) => route(dispatch, message);
    eventSource.onerror = () => {
      dispatch(setConnectionState(false));
      displayNotification({
        appearance: NotificationType.DANGER,
        title: "Unable to connect to server.",
        timestamp: new Date().getTime()
      });
    };

    return (next) => (action) => next(action);
  };

const onOpen = (
  dispatch: ThunkDispatch<unknown, unknown, UnknownAction>
): void => {
  console.log("Event Source connected.");
  dispatch(setConnectionState(true));
};

const route = (
  dispatch: ThunkDispatch<unknown, unknown, UnknownAction>,
  msg: MessageEvent<string>
): void => {
  const getNotif = (): Notification => {
    const response = JSON.parse(msg.data) as Response;
    const timestamp = new Date().getTime();
    const notification: Notification = {
      ...response,
      timestamp
    };

    switch (response.type) {
      case MessageType.STATUS:
        return notification;
      case MessageType.CONNECT:
        dispatch(setUsername((response as ConnectionResponse).name));
        return notification;
      case MessageType.SEARCH:
        dispatch(resultsReceived(response as SearchResponse));
        return notification;
      case MessageType.DOWNLOAD:
        downloadFile((response as DownloadResponse)?.downloadPath);
        dispatch(api.util.invalidateTags(["books"]));
        dispatch(removeInFlightDownload());
        return notification;
      case MessageType.RATELIMIT:
        return notification;
      default:
        console.error(response);
        return {
          appearance: NotificationType.DANGER,
          title: "Unknown message type. See console.",
          timestamp
        };
    }
  };

  const notif = getNotif();
  dispatch(addNotification(notif));
  displayNotification(notif);
};
