export enum NotificationType {
  NOTIFY,
  SUCCESS,
  WARNING,
  DANGER
}

export type RequestMessage =
  | { type: MessageType.CONNECT; payload: ConnectPayload }
  | { type: MessageType.DOWNLOAD; payload: DownloadPayload }
  | { type: MessageType.SEARCH; payload: SearchPayload };

export type ResponseMessage = { type: MessageType.STATUS; payload: string };

interface ConnectPayload {
  address: string;
  enableTLS: boolean;
  channel: string;
}

interface DownloadPayload {
  book: string;
}

interface SearchPayload {}

export enum MessageType {
  STATUS,
  CONNECT,
  SEARCH,
  DOWNLOAD,
  RATELIMIT
}

// Notification is used to show a UI toast notification the user.
export interface Notification {
  appearance: NotificationType;
  title: string;
  detail?: string;
  timestamp: number;
}

// Response is received from websocket requests
export interface Response extends Omit<Notification, "timestamp"> {
  type: MessageType;
}

// ConnectionResponse is received after successful IRC connection
export interface ConnectionResponse extends Response {
  name: string;
}

// SearchResponse is received after search results are received and parsed.
export interface SearchResponse extends Response {
  books: BookDetail[];
  errors: ParseError[];
}

// DownloadResponse is received after file is downloaded from IRC and ready for
// user download.
export interface DownloadResponse extends Response {
  downloadPath?: string;
}

export interface BookDetail {
  server: string;
  author: string;
  title: string;
  format: string;
  size: string;
  full: string;
}

export interface ParseError {
  error: string;
  line: string;
}
