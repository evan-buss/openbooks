import { AnyAction, PayloadAction, Store } from "@reduxjs/toolkit";
import { BookDetail, BookResponse, MessageType } from "../models/messages";
import { displayNotification, NotificationType } from "./notifications";
import { setServers } from "./serverSlice";
import { sendMessage, setConnectionState, setSearchResults } from "./stateSlice";

// Web socket redux middleware. 
// Listens to socket and dispatches handlers. 
// Handles send_message actions by sending to socket.
export const websocketConn = (wsUrl: string): any => {
    return (store: Store) => {
        console.log(store);
        const socket = new WebSocket(wsUrl);

        socket.onopen = () => onOpen(store);
        socket.onclose = () => onClose(store);
        socket.onmessage = (message) => route(store, message);
        socket.onerror = (event) => console.error(event);

        return (next: any) => (action: PayloadAction<any>) => {
            // Send Message action? Send data to the socket.
            if (sendMessage.match(action)) {
                if (socket.readyState === socket.OPEN) {
                    socket.send(action.payload.message)
                } else {
                    displayNotification(NotificationType.WARNING, "Server connection closed. Reload page.");
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
            displayNotification(NotificationType.SUCCESS, "Welcome, connection established.");
            store.dispatch(sendMessage({ type: MessageType.SERVERS, payload: {} }));
            break;
        case MessageType.SEARCH:
            console.log("search results")
            store.dispatch((setSearchResults(response.books as BookDetail[]) as unknown) as AnyAction);
            displayNotification(NotificationType.SUCCESS, "Search results received.")
            break;
        case MessageType.DOWNLOAD:
            displayNotification(NotificationType.SUCCESS, "Book file received.", response.name);
            saveByteArray(response.name, response.file);
            break;
        case MessageType.SERVERS:
            store.dispatch(setServers(response.servers))
            break;
        case MessageType.WAIT:
            displayNotification(NotificationType.SUCCESS, response.status)
            break;
        case MessageType.IRCERROR:
            displayNotification(NotificationType.DANGER, "IRC Error. Try again.")
            break;
        default:
            console.error(response);
            displayNotification(NotificationType.DANGER, "Unknown error type. See console.")
    }
}

const onOpen = (store: Store): void => {
    console.log("WebSocket connected.");
    store.dispatch(setConnectionState(true));
    store.dispatch(sendMessage({ type: MessageType.CONNECT, payload: {} }));
}

// saveByteArray creates a link and download popup for the returned file
function saveByteArray(fileName: string, byte: Blob) {
    let link = document.createElement('a');
    link.href = `data:application/octet-stream;base64,${byte}`;
    link.download = fileName;
    link.click();
};

const onClose = (store: Store): void => {
    console.log("WebSocket closed.");
    store.dispatch(setConnectionState(false));
}
