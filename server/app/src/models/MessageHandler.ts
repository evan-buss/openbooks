import { updateServers } from "../state/serverSlice";
// TODO: Figure out if this is bad...
import { store } from "../state/store";
import { BookRequest, BookResponse, MessageType } from "./messages";


export class MessageHandler {

    private ws: WebSocket;

    constructor() {
        const wsProtocol = window.location.protocol === "https:" ? "wss://" : "ws://";
        const wsHost = process.env.NODE_ENV === "development" ? "localhost:5228" : window.location.host;

        this.ws = new WebSocket(`${wsProtocol}${wsHost}/ws`);
        this.ws.onopen = () => this.onOpen();
        this.ws.onclose = () => this.onClose();
        this.ws.onmessage = (message) => this.route(message);
        this.ws.onerror = (event) => console.error(event);
    }

    public route(msg: MessageEvent<any>): void {
        console.log(msg);

        let response = JSON.parse(msg.data) as BookResponse;
        switch (response.type) {
            case MessageType.ERROR:
                this.displayNotification(response.details)
                break;
            case MessageType.CONNECT:
                this.displayNotification("Welcome, connection established.");
                break;
            case MessageType.SEARCH:
                // this.loading$.next(false);
                // this.searchResults$.next(message.books as BookDetail[]);
                this.displayNotification("Search results returned.")
                break;
            case MessageType.DOWNLOAD:
                break;
            case MessageType.SERVERS:
                // this.servers$.next(response.servers);
                store.dispatch(updateServers(response.servers))
                break;
            case MessageType.WAIT:
                console.log(response);
                this.displayNotification(response.status)
                break;
            case MessageType.IRCERROR:
                this.displayNotification("Internal Server Error. Try again.")
                break;
        }
    }

    private onOpen(): void {
        console.log("WebSocket connected.");
        this.ws.send(JSON.stringify({ type: MessageType.CONNECT, payload: {} } as BookRequest));

        setTimeout(() => {
            this.ws.send(JSON.stringify({ type: MessageType.SERVERS, payload: {} } as BookRequest));
        }, 5000);
    }

    private onClose(): void {
        console.log("WebSocket closed.");
    }

    private displayNotification(message: string): void {

    }

    public dispose(): void {
        console.log("Disposing websocket connection.");
        this.ws.close();
    }
}