import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTable } from '@angular/material/table';
import { BookGridDataSource, } from './book-grid-datasource';
import { BookService } from '../book-service.service';
import { BookDetail } from '../messages';

@Component({
	selector: 'app-book-grid',
	templateUrl: './book-grid.component.html',
	styleUrls: ['./book-grid.component.css']
})
export class BookGridComponent implements AfterViewInit, OnInit {
	@ViewChild(MatPaginator) paginator: MatPaginator;
	@ViewChild(MatSort) sort: MatSort;
	@ViewChild(MatTable) table: MatTable<BookDetail>;
	dataSource: BookGridDataSource;


	/** Columns displayed in the table. Columns IDs can be added, removed, or reordered. */
	displayedColumns = ['server', 'author', 'title', 'format', 'size', 'full'];

	// TODO: Use breakpoint observer to hide columns
	constructor(private service: BookService) { }

	ngOnInit() {
		this.dataSource = new BookGridDataSource(this.service);
	}
	applyFilter(val: string) {
		val = val.trim().toLowerCase();
		// TODO: Figure this shit out...
		// this.dataSource.filter = filterValue;
	}

	ngAfterViewInit() {
		this.dataSource.sort = this.sort;
		this.dataSource.paginator = this.paginator;
		this.table.dataSource = this.dataSource;
	}
}
