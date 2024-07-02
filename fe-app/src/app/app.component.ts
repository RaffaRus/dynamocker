import { Component } from '@angular/core';
import { MonacoEditorService } from '@services/monacoEditorService';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
})

export class AppComponent {
  title = 'dynamocker';

  constructor(
    private monacoEditorService : MonacoEditorService
  ){
    monacoEditorService.load()
  }
}
