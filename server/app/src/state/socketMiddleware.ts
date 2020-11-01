import { AnyAction, PayloadAction, Store } from "@reduxjs/toolkit";
import { BookDetail, BookResponse, MessageType } from "../models/messages";
import { displayNotification, NotificationType } from "./notifications";
import { setServers } from "./serverSlice";
import { sendMessage, setSearchResults } from "./stateSlice";

// Web socket redux middleware. 
// Listens for messages and sends dispatches to the socket
export const websocketConn = (wsUrl: string): any => {
    return (store: Store) => {
        console.log(store);
        const socket = new WebSocket(wsUrl);

        socket.onopen = () => onOpen(store);
        socket.onclose = () => onClose();
        socket.onmessage = (message) => route(store, message);
        socket.onerror = (event) => console.error(event);

        return (next: any) => (action: PayloadAction<any>) => {
            // Send Message action? Send data to the socket.
            if (sendMessage.match(action)) {
                if (socket.readyState === socket.OPEN) {
                    socket.send(action.payload.message)
                } else {
                    socket.send("not sending because socket is not open.");
                }
            }

            return next(action);
        }
    }
}

const route = (store: Store, msg: MessageEvent<any>): void => {
    console.log(msg);

    let response = JSON.parse(msg.data) as BookResponse;
    switch (response.type) {
        // TODO: How to get the message type with typed properties
        case MessageType.ERROR:
            displayNotification(NotificationType.DANGER, response.details)
            break;
        case MessageType.CONNECT:
            displayNotification(NotificationType.NOTIFY, "Welcome, connection established.");
            store.dispatch(sendMessage({ type: MessageType.SERVERS, payload: {} }));
            break;
        case MessageType.SEARCH:
            console.log("search results")
            store.dispatch((setSearchResults(response.books as BookDetail[]) as unknown) as AnyAction);
            displayNotification(NotificationType.WARNING, "Search results returned.")
            break;
        case MessageType.DOWNLOAD:
            break;
        case MessageType.SERVERS:
            // this.servers$.next(response.servers);
            store.dispatch(setServers(response.servers))
            break;
        case MessageType.WAIT:
            console.log(response);
            displayNotification(NotificationType.WARNING, response.status)
            break;
        case MessageType.IRCERROR:
            displayNotification(NotificationType.WARNING, "Internal Server Error. Try again.")
            break;
    }
}

const onOpen = (store: Store): void => {
    console.log("WebSocket connected.");
    store.dispatch(sendMessage({ type: MessageType.CONNECT, payload: {} }));
}

const onClose = (): void => {
    console.log("WebSocket closed.");
}
