export enum MessageTypes {
	ERROR = 0,
	CONNECT = 1,
	SEARCH = 2,
	DOWNLOAD = 3,
	SERVERS = 4,
	WAIT = 5,
	IRCERROR = 6
}

export interface BookResponse {
	type: MessageTypes,
	[x: string]: any
}

export interface BookRequest {
	type: MessageTypes,
	payload: any
}

export interface ServersResponse {
	type: MessageTypes,
	servers: string[]
}

export interface ErrorResponse {
	error: number,
	details: string
}

export interface IrcErrorResponse {
	type: MessageTypes,
	status: string
}

export interface WaitResponse {
	type: MessageTypes,
	status: string
}

export interface ConnectionResponse {
	type: MessageTypes,
	Status: string,
	wait: string,
}

export interface SearchRequest {
	type: MessageTypes,
	query: string
}

export interface SearchResponse {
	type: MessageTypes,
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
	type: MessageTypes,
	name: string,
	//TODO: What type to make a binary blob?
	file: any
}