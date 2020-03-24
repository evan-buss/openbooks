import { Component } from '@angular/core';
import { BreakpointObserver, Breakpoints } from '@angular/cdk/layout';
import { Observable, BehaviorSubject } from 'rxjs';
import { map, shareReplay } from 'rxjs/operators';
import { BookService } from '../book-service.service';

@Component({
	selector: 'app-nav',
	templateUrl: './nav.component.html',
	styleUrls: ['./nav.component.css']
})
export class NavComponent {

	isHandset$: Observable<boolean> = this.breakpointObserver.observe(Breakpoints.Handset)
		.pipe(
			map(result => result.matches),
			shareReplay()
		);

	public get serverStream(): BehaviorSubject<string[]> {
		return this.books.servers$;
	}

	refreshServers() {
		return this.books.refreshServers();
	}

	constructor(private breakpointObserver: BreakpointObserver, public books: BookService) { }
}
