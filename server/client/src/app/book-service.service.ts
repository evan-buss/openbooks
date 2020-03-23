import { Injectable } from '@angular/core';
import { webSocket } from 'rxjs/webSocket';
import { BookResponse, BookRequest, MessageTypes, ErrorResponse, SearchRequest, BookDetail } from './messages';
import { BehaviorSubject, Observable } from 'rxjs';
import { MatSnackBar } from '@angular/material/snack-bar';
import { delay } from 'rxjs/operators';
var dummyData = require('src/assets/dummy.json');

@Injectable({
	providedIn: 'root'
})
export class BookServiceService {

	private connection$ = webSocket("ws://localhost:5228/ws");
	public servers$: BehaviorSubject<string[]> = new BehaviorSubject<string[]>(["Please Refresh"]);;
	public searchResults$ = new BehaviorSubject<BookDetail[]>([]);

	constructor(private snackbar: MatSnackBar) {
		console.log("constructor");
		this.connection$.subscribe(
			(msg) => this.messageRouter(msg),
			err => console.log("err: " + err),
			() => this.sendNotification("Connection closed")
		);
		this.servers$.next(["test", "wtf"]);
		this.connection$.next(<BookRequest>{ type: MessageTypes.CONNECT, payload: {} })
	}

	public refreshServers(): void {
		this.connection$.next(<BookRequest>{ type: MessageTypes.SERVERS, payload: {} })
	}

	public searchBook(search: string) {
		// this.searchResults$.next(dummyData);
		this.connection$.next({ type: MessageTypes.SEARCH, payload: { query: search } });
	}

	private messageRouter(msg: any): void {
		console.log(msg);
		let message = msg as BookResponse;
		switch (message.type) {
			case MessageTypes.ERROR:
				this.sendNotification(message.details)
				break;
			case MessageTypes.CONNECT:
				this.sendNotification("Welcome, connection established.");
				break;
			case MessageTypes.SEARCH:
				this.searchResults$.next(message.books as BookDetail[]);
				this.sendNotification("Search results returned.")
				break;
			case MessageTypes.DOWNLOAD:
				break;
			case MessageTypes.SERVERS:
				this.servers$.next(message.servers);
				break;
			case MessageTypes.WAIT:
				console.log(message);
				this.sendNotification(message.status)
				break;
			case MessageTypes.IRCERROR:
				this.sendNotification("Internal Server Error. Try again.")
				break;
		}
	}

	private sendNotification(message: string): void {
		this.snackbar.open(message, 'Dismiss');
	}
}