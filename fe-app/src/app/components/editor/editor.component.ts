import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { IMockApi, dummyMockApi } from '@models/mockApi';
import { initialMockApiJson, notificationLevel } from '@models/monaco';
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
  
  public codeEditorInstance!: monaco.editor.IStandaloneCodeEditor;
  public isSaveButtonEnabled : boolean = false

  private _selectedMockApi: IMockApi = dummyMockApi
  private _isSelectedMockApiNew: boolean = true

  constructor(
    private mockApiService: MockApiService,
    private editorService: EditorService
  ) {
    this.assignDefaultJsonSchema()
  }

  ngOnInit(): void {

    // disable save button at startup
    this.isSaveButtonEnabled = false

    // create monaco editor
    this.codeEditorInstance = monaco.editor.create(this._editorContainer.nativeElement, {
      theme: 'vs',
      wordWrap: 'on',
      wrappingIndent: 'indent',
      language: 'json',
      minimap: { enabled: false },
      automaticLayout: true,
      value: initialMockApiJson
    });

    // subscribe to "new_mock_api_selected"
    this.mockApiService.newMockApiSelectedObservable().subscribe({
      next: (mockApi: IMockApi) => {
        this._selectedMockApi = mockApi
        this._isSelectedMockApiNew = false
      },
      error: (err: Error) => {
        console.error("received error from newMockApiSelectedObservable subscription: ", err)
      }
    })

    // set the editor property 'unsavedModifications' as true
    if (this.codeEditorInstance.getModel() != null) {
      this.codeEditorInstance.getModel()!.onDidChangeContent((event) => {
        this.editorService.unsavedModifications = true
      });
    }
    
    // enable/diable save button based on the json validation schema
    if (this.codeEditorInstance.getModel() != null) {
      monaco.editor.onDidChangeMarkers( event => {
        // enable save button only if there is no error in the JSON
        if (monaco.editor.getModelMarkers({}).length === 0) {
          this.isSaveButtonEnabled = true
        } else {
          this.isSaveButtonEnabled = false
        }
      })
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
        error: (err : Error) => {
          console.log("Could not create new MockApi: "+err)
          this.notify("Could not create new MockApi: "+ err.message, notificationLevel.error)

        },
        complete: () => {
          const msg = "MockApi succesfully created"
          this.notify(msg, notificationLevel.info)
          this.editorService.unsavedModifications = false
          this.mockApiService.refreshList()
        }
      }
      )
    } else {
      this.mockApiService.patchMockApi(this._selectedMockApi).pipe().subscribe({
        error: (err : Error) => {
          console.log("Could not modify new MockApi: "+ err.message)
          this.notify("Could not modify new MockApi: "+ err.message, notificationLevel.error)

        },
        complete: () => {
          const msg = "MockAPi succesfullly modified"
          this.notify(msg, notificationLevel.info)
          this.editorService.unsavedModifications = false
          this.mockApiService.refreshList()
        }
      }
      )

    }
  }
  
  notify(msg: string, level: notificationLevel): void {
    this.jqxNotification.elementRef.nativeElement.querySelector('.jqx-notification-content').innerHTML = msg;
    this.jqxNotification.elementRef.nativeElement.querySelector('.notificationContainer').classList.add(level)
    this.jqxNotification.open();
};
}