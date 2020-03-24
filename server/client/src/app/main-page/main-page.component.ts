import { Component, OnInit } from '@angular/core';
import { BookService } from '../book-service.service';
import { FormGroup, FormControl, Validators } from '@angular/forms';

@Component({
	selector: 'app-main-page',
	templateUrl: './main-page.component.html',
	styleUrls: ['./main-page.component.css']
})
export class MainPageComponent implements OnInit {

	searchForm = new FormGroup({
		search: new FormControl("", Validators.required)
	});

	constructor(private service: BookService) { }

	ngOnInit(): void {
	}

	searchBook() {
		this.service.searchBook(this.searchForm.value.search);
	}
}
