import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { BackgroundComponent } from './components/background/background.component';
import { EditorComponent } from './components/editor/editor.component';
import { ApiListItemComponent } from './components/api-list-item/api-list-item.component';
import { ApiListComponent } from './components/api-list/api-list.component';

@Component({
  selector: 'app-root',
  standalone: true,
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
  imports: [
    RouterOutlet,
    BackgroundComponent,
    EditorComponent,
    ApiListItemComponent,
    ApiListComponent
  ]
})

export class AppComponent {
  title = 'dynamocker';
}
