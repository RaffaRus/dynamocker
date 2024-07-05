import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';
import { AppModule } from './app/app.module';
import * as monaco from 'monaco-editor';
// or import * as monaco from 'monaco-editor/esm/vs/editor/editor.api';
// if shipping only a subset of the features & languages is desired

platformBrowserDynamic().bootstrapModule(AppModule).catch(err => console.error(err));

monaco.editor.create(document.getElementById('editorContainer') as HTMLElement, {
    value: 'console.log("Hello, world")',
	language: 'javascript'
});
