// Angular
import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';

// Material
import {MatIconModule} from '@angular/material/icon';
import {MatBadgeModule} from '@angular/material/badge';
import {MatDividerModule} from '@angular/material/divider';
import {MatButtonModule} from '@angular/material/button';

// Components
import { AppComponent } from './app.component';
import { ApiListComponent } from '@components/api-list/api-list.component';
import { ApiListItemComponent } from '@components/api-list-item/api-list-item.component';
import { BackgroundComponent } from '@components/background/background.component';
import { EditorComponent } from '@components/editor/editor.component';
import { TreeComponent } from '@components/tree/tree.component';

// Services
import { MockApiService } from '@services/mockApiService';
import { MonacoEditorService } from '@services/monacoEditorService';


@NgModule({
  declarations: [ 
    AppComponent, 
    BackgroundComponent, 
    ApiListComponent,
    ApiListItemComponent,
    EditorComponent,
    TreeComponent
  ],
  imports: [ 
    BrowserModule, 
    FormsModule, 
    MatIconModule, 
    MatDividerModule,
    MatBadgeModule, 
    MatButtonModule,
    CommonModule,
    HttpClientModule
  ],
  bootstrap:    [ 
    AppComponent 
  ],
  providers: [
    MockApiService,
    HttpClient,
    MonacoEditorService
  ]
})
export class AppModule { }

