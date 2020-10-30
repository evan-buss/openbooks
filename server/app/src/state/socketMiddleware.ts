import { PayloadAction, Store } from "@reduxjs/toolkit";
import { BookResponse, MessageType } from "../models/messages";
import { setServers } from "./serverSlice";
import { sendMessage } from "./stateSlice";

// Web socket redux middleware. 
// Listens for messages and sends dispatches to the socket
export const websocketConn = (wsUrl: string) => {
    return (store: any) => {
        const socket = new WebSocket(wsUrl);

        socket.onopen = () => onOpen(store);
        socket.onclose = () => onClose();
        socket.onmessage = (message) => route(store, message);
        socket.onerror = (event) => console.error(event);

        return (next: any) => (action: PayloadAction<any>) => {
            if (sendMessage.match(action)) {
                socket.send(action.payload.message)
            }
            console.log(action);
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
            displayNotification(response.details)
            break;
        case MessageType.CONNECT:
            displayNotification("Welcome, connection established.");
            store.dispatch(sendMessage({ type: MessageType.SERVERS, payload: {} }));
            break;
        case MessageType.SEARCH:
            // this.loading$.next(false);
            // this.searchResults$.next(message.books as BookDetail[]);
            displayNotification("Search results returned.")
            break;
        case MessageType.DOWNLOAD:
            break;
        case MessageType.SERVERS:
            // this.servers$.next(response.servers);
            store.dispatch(setServers(response.servers))
            break;
        case MessageType.WAIT:
            console.log(response);
            displayNotification(response.status)
            break;
        case MessageType.IRCERROR:
            displayNotification("Internal Server Error. Try again.")
            break;
    }
}

const displayNotification = (message: string) => {
    console.log(message);
}

const onOpen = (store: Store): void => {
    console.log("WebSocket connected.");
    store.dispatch(sendMessage({ type: MessageType.CONNECT, payload: {} }));
}

const onClose = (): void => {
    console.log("WebSocket closed.");
}
