export enum NotificationType {
  NOTIFY,
  SUCCESS,
  WARNING,
  DANGER
}

export interface Notification {
  type: NotificationType;
  title: string;
  detail?: string;
  timestamp: number;
}

export enum MessageType {
  STATUS,
  CONNECT,
  SEARCH,
  DOWNLOAD,
  RATELIMIT
}

export interface BookResponse {
  type: MessageType;
  notify: NotificationType;
  title: string;
  detail?: string;

  [x: string]: any;
}

export interface SearchResponse {
  type: MessageType;
  books: BookDetail[];
  errors: ParseError[];
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
