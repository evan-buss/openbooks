export enum MessageType {
    ERROR = 0,
    CONNECT = 1,
    SEARCH = 2,
    DOWNLOAD = 3,
    SERVERS = 4,
    WAIT = 5,
    IRCERROR = 6
}

export interface BookResponse {
    type: MessageType,
    [x: string]: any
}

export interface BookRequest {
    type: MessageType,
    payload: any
}

export interface ServersResponse {
    type: MessageType,
    servers: string[]
}

export interface ErrorResponse {
    error: number,
    details: string
}

export interface IrcErrorResponse {
    type: MessageType,
    status: string
}

export interface WaitResponse {
    type: MessageType,
    status: string
}

export interface ConnectionResponse {
    type: MessageType,
    Status: string,
    wait: string,
}

export interface SearchRequest {
    type: MessageType,
    query: string
}

export interface SearchResponse {
    type: MessageType,
    books: BookDetail[]
}

export interface BookDetail {
    server: string,
    author: string,
    title: string,
    format: string,
    size: string,
    full: string
}

export interface DownloadResponse {
    type: MessageType,
    name: string,
    //TODO: What type to make a binary blob?
    file: any
}