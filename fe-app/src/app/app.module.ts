import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';

import {MatIconModule} from '@angular/material/icon';
import {MatBadgeModule} from '@angular/material/badge';
import {MatButtonModule} from '@angular/material/button';

import { AppComponent } from './app.component';
import { BackgroundComponent } from './components/background/background.component';
import { ApiListItemComponent } from './components/api-list-item/api-list-item.component';


@NgModule({
  declarations: [ AppComponent, BackgroundComponent, ApiListItemComponent ],
  imports:      [ BrowserModule, FormsModule, MatIconModule, MatBadgeModule, MatButtonModule],
  bootstrap:    [ AppComponent ]
})
export class AppModule { }

