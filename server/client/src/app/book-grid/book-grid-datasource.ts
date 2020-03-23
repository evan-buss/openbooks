import { DataSource } from '@angular/cdk/collections';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { map, tap } from 'rxjs/operators';
import { Observable, merge } from 'rxjs';
import { BookServiceService } from '../book-service.service';
import { BookDetail } from '../messages';


/**
 * Data source for the BookGrid view. This class should
 * encapsulate all logic for fetching and manipulating the displayed data
 * (including sorting, pagination, and filtering).
 */
export class BookGridDataSource extends DataSource<BookDetail> {
	data: BookDetail[] = [];
	paginator: MatPaginator;
	sort: MatSort;

	constructor(private books: BookServiceService) {
		super();
	}

	/**
	 * Connect this data source to the table. The table will only update when
	 * the returned stream emits new items.
	 * @returns A stream of the items to be rendered.
	 */
	connect(): Observable<BookDetail[]> {
		// Combine everything that affects the rendered data into one update
		// stream for the data-table to consume.
		const dataMutations = [
			this.books.searchResults$.asObservable().pipe(tap(data => this.data = data)),
			this.paginator.page,
			this.sort.sortChange
		];

		return merge(...dataMutations).pipe(map(() => {
			return this.getPagedData(this.getSortedData([...this.data]));
		}));
	}

	/**
	 *  Called when the table is being destroyed. Use this function, to clean up
	 * any open connections or free any held resources that were set up during connect.
	 */
	disconnect() { }

	/**
	 * Paginate the data (client-side). If you're using server-side pagination,
	 * this would be replaced by requesting the appropriate data from the server.
	 */
	private getPagedData(data: BookDetail[]) {
		const startIndex = this.paginator.pageIndex * this.paginator.pageSize;
		return data.splice(startIndex, this.paginator.pageSize);
	}

	/**
	 * Sort the data (client-side). If you're using server-side sorting,
	 * this would be replaced by requesting the appropriate data from the server.
	 */
	private getSortedData(data: BookDetail[]) {
		if (!this.sort.active || this.sort.direction === '') {
			return data;
		}

		return data.sort((a, b) => {
			const isAsc = this.sort.direction === 'asc';
			switch (this.sort.active) {
				case 'server': return compare(a.server, b.server, isAsc);
				case 'author': return compare(a.author, b.author, isAsc);
				case 'title': return compare(a.title, b.title, isAsc);
				case 'format': return compare(a.format, b.format, isAsc);
				case 'size': return compare(a.size, b.size, isAsc);
				default: return 0;
			}
		});
	}
}

/** Simple sort comparator for example ID/Name columns (for client-side sorting). */
function compare(a: string | number, b: string | number, isAsc: boolean) {
	return (a < b ? -1 : 1) * (isAsc ? 1 : -1);
}
