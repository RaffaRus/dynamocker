// Angular
import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { HttpClient, provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';

// Material
import {MatIconModule} from '@angular/material/icon';
import {MatBadgeModule} from '@angular/material/badge';
import {MatDividerModule} from '@angular/material/divider';
import {MatButtonModule} from '@angular/material/button';

// Components
import { AppComponent } from './app.component';
import { ApiListComponent } from '@components/api-list/api-list.component';
import { ApiListItemComponent } from '@components/api-list-item/api-list-item.component';
import { EditorComponent } from '@components/editor/editor.component';
import { TreeComponent } from '@components/tree/tree.component';
import { jqxButtonModule } from 'jqwidgets-ng/jqxbuttons';
import { jqxNotificationModule } from 'jqwidgets-ng/jqxnotification';

// Services
import { MockApiService } from '@services/mockApiService';
import { EditorService } from '@services/editor';

@NgModule({
    declarations: [
        AppComponent,
        ApiListComponent,
        ApiListItemComponent,
        EditorComponent,
        TreeComponent
    ],
    bootstrap: [
        AppComponent
    ],
    imports: [
        BrowserModule,
        FormsModule,
        MatIconModule,
        MatDividerModule,
        MatBadgeModule,
        MatButtonModule,
        CommonModule,
        jqxNotificationModule,
        jqxButtonModule
    ],
    providers: [
        MockApiService,
        EditorService,
        HttpClient,
        provideHttpClient(withInterceptorsFromDi())
    ] })

export class AppModule { }

