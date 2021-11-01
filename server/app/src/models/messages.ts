export enum MessageType {
  ERROR = 1,
  CONNECT = 2,
  SEARCH = 3,
  DOWNLOAD = 4,
  // SERVERS = 5,
  WAIT = 5,
  IRCERROR = 6
}

export interface BookResponse {
  type: MessageType;

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
