import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NavComponent } from './nav/nav.component';
import { LayoutModule } from '@angular/cdk/layout';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';
import { BookGridComponent } from './book-grid/book-grid.component';
import { MatTableModule } from '@angular/material/table';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatSortModule } from '@angular/material/sort';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { FlexLayoutModule } from '@angular/flex-layout';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { BookServiceService } from './book-service.service';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';

@NgModule({
	declarations: [
		AppComponent,
		NavComponent,
		BookGridComponent,
	],
	imports: [
		BrowserModule,
		BrowserAnimationsModule,
		FormsModule,
		FlexLayoutModule,
		LayoutModule,
		MatExpansionModule,
		MatToolbarModule,
		MatButtonModule,
		MatFormFieldModule,
		MatInputModule,
		MatSidenavModule,
		MatIconModule,
		ReactiveFormsModule,
		MatListModule,
		MatTableModule,
		MatSnackBarModule,
		MatPaginatorModule,
		MatSortModule
	],
	providers: [BookServiceService],
	bootstrap: [AppComponent]
})
export class AppModule { }
