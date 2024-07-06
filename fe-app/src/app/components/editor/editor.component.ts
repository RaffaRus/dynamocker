import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { IMockApi, dummyMockApi } from '@models/mockApi';
import { initialMockApiJson } from '@models/monaco';
import { EditorService } from '@services/editor';
import { MockApiService } from '@services/mockApiService';
import { jqxNotificationComponent } from 'jqwidgets-ng/jqxnotification';
import * as monaco from 'monaco-editor';

@Component({
  selector: 'app-editor',
  templateUrl: './editor.component.html',
  styleUrl: './editor.component.scss'
})
export class EditorComponent implements OnInit {
  @ViewChild('monacoDivTag', { static: true }) _editorContainer!: ElementRef;
  @ViewChild('jqxNotification', { static: false }) jqxNotification!: jqxNotificationComponent;
  
  quotes: string[] =
  [
      'I\'m gonna make him an offer he can\'t refuse.', 'Toto, I\'ve got a feeling we\'re not in Kansas anymore.',
      'You talkin\' to me?', 'Bond. James Bond.', 'I\'ll be back.', 'Round up the usual suspects.',
      'I\'m the king of the world!', 'A martini. Shaken, not stirred.',
      'May the Force be with you.',
      'Well, nobody\'s perfect.'
  ];
  notificationClick(): void {
      this.jqxNotification.elementRef.nativeElement.querySelector('.jqx-notification-content').innerHTML = this.quotes[Math.round(Math.random() * this.quotes.length - 1)];
      this.jqxNotification.open();
  };


  codeEditorInstance!: monaco.editor.IStandaloneCodeEditor;

  private _selectedMockApi: IMockApi = dummyMockApi
  private _isSelectedMockApiNew: boolean = true
  private _unsavedModifications : boolean = false

  constructor(
    private mockApiService: MockApiService,
    private editorService: EditorService
  ) {
    this.assignDefaultJsonSchema()
  }

  ngOnInit(): void {


    // create monaco editor
    this.codeEditorInstance = monaco.editor.create(this._editorContainer.nativeElement, {
      theme: 'vs',
      wordWrap: 'on',
      wrappingIndent: 'indent',
      language: 'json',
      // minimap: { enabled: false },
      automaticLayout: true,
      value: initialMockApiJson
    });

    // subscribe to modifications of the json in the monaco editor
    // this.editorService.unsavedModifications().subscribe({
    //   next: (unsaved : boolean) => {
    //     this._unsavedModifications = unsaved
    //   },
    //   error: (err: Error) => {
    //     console.error("received error from unsavedModifications subscription: ", err)
    //   }
    // })

    // subscribe to "new_mock_api_selected"
    this.mockApiService.newMockApiSelectedObservable().subscribe({
      next: (mockApi: IMockApi) => {
        this._selectedMockApi = mockApi
        this._isSelectedMockApiNew = false
        this._unsavedModifications = false
      },
      error: (err: Error) => {
        console.error("received error from newMockApiSelectedObservable subscription: ", err)
      }
    })

    // set the editor property 'unsavedModifications' as true
    if (this.codeEditorInstance.getModel() != null) {
      let model = this.codeEditorInstance.getModel() as monaco.editor.ITextModel
      model.onDidChangeContent((event) => {
        this.editorService.unsavedModifications = true
      });
    }
  }

  // used to assign our custom schema for the json file corresponding to the MockApi schema
  assignDefaultJsonSchema() {
    monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
      validate: true,
      schemas: [{
        uri: 'http://localhost:4200/mock-api-schema.json',
        fileMatch: ["*"],
        schema: {
          type: 'object',
          properties : {
            name: {
              type: 'string'
            },
            url : {
              type: 'string'
            },
            responses : {
              type: 'object',
              properties: {
                get: {
                  type: 'object'
                },
                post: {
                  type: 'object'
                },
                delete: {
                  type: 'object'
                },
                patch: {
                  type: 'object'
                },
              },
              minProperties : 1
            },
          },
          required : [
            "name",
            "url",
            "responses"
          ]
        }
      }
    ]
    });
  }

  onSave() {
    if (this._isSelectedMockApiNew) {
      this.mockApiService.postMockApi(this._selectedMockApi).pipe().subscribe({
        next: (asd) => {

        },
        error: (err) => {
          console.log(err)


        },
        complete: () => {
          this.editorService.unsavedModifications = false
          this.mockApiService.refreshList()
        }
        
      }
      )
    } else {
      this.mockApiService.patchMockApi(this._selectedMockApi).pipe().subscribe({
        next: (asd) => {

        },
        error: (err) => {
          console.log(err)
          // notification(err)

        },
        complete: () => {
          this.editorService.unsavedModifications = false
          this.mockApiService.refreshList()
        }
        
      }
      )

    }
  }
}