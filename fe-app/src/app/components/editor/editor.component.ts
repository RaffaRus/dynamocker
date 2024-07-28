// Core
import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';

// External
import * as monaco from 'monaco-editor';

// Components
import { jqxNotificationComponent } from 'jqwidgets-ng/jqxnotification';

// Models
import { IMockApi } from '@models/mockApi.model';
import { initialMockApiJson, initialMockApiJsonString, notificationLevel } from '@models/editor.model';
import { DynamockerBackendErrorMessageI } from '@models/httpError.model';

// Services
import { EditorService } from '@services/editor';
import { MockApiService } from '@services/mockApiService';
import { HttpErrorResponse } from '@angular/common/http';


@Component({
  selector: 'app-editor',
  templateUrl: './editor.component.html',
  styleUrl: './editor.component.scss'
})
export class EditorComponent implements OnInit {
  @ViewChild('monacoDivTag', { static: true }) _editorContainer!: ElementRef;
  @ViewChild('jqxNotification', { static: false }) jqxNotification!: jqxNotificationComponent;
  
  public codeEditorInstance!: monaco.editor.IStandaloneCodeEditor;
  public isSaveButtonEnabled : boolean = true

  private _selectedMockApi: IMockApi = initialMockApiJson
  private _isSelectedMockApiNew: boolean = true

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
      minimap: { enabled: false },
      automaticLayout: true,
      value: initialMockApiJsonString
    });

    // subscribe to "new_mock_api_selected"
    this.mockApiService.newMockApiSelectedObservable().subscribe({
      next: (mockApi: IMockApi) => {
        this._selectedMockApi = mockApi
        this._isSelectedMockApiNew = mockApi == initialMockApiJson
        this.codeEditorInstance.setValue(JSON.stringify(mockApi, null, 2))
      },
      error: (err: Error) => {
        console.error("received error from newMockApiSelectedObservable subscription: ", err)
      }
    })

    // set the editor property 'unsavedModifications' as true
    if (this.codeEditorInstance.getModel() != null) {
      this.codeEditorInstance.getModel()!.onDidChangeContent((event) => {
        this.editorService.unsavedModifications = true
        this._selectedMockApi = JSON.parse(this.codeEditorInstance.getValue())
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
    // TODO: set empty strings as invalid
    // TODO: make fields different than {name|url|response}
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
          ],
          "additionalProperties": false
        }
      }
    ]
    });
  }

  onSave() {
    if (!this.isSaveButtonEnabled) {
      return
    }
    if (this._isSelectedMockApiNew) {
      this.mockApiService.postMockApi(this._selectedMockApi).pipe().subscribe({
        error: (err : HttpErrorResponse) => {
          const toNotify = err.error as DynamockerBackendErrorMessageI
          this.notify(toNotify.error_msg, notificationLevel.error)
          },
        complete: () => {
          const msg = "MockApi succesfully created"
          this.notify(msg, notificationLevel.info)
          this.editorService.unsavedModifications = false
          this.mockApiService.refreshList()
          this._isSelectedMockApiNew = false
        }
      }
    )
  } else {
    this.mockApiService.putMockApi(this._selectedMockApi).pipe().subscribe({
      error: (err : HttpErrorResponse) => {
        const toNotify = err.error as DynamockerBackendErrorMessageI
        this.notify(toNotify.error_msg, notificationLevel.error)
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
    console.log(this.jqxNotification.elementRef.nativeElement.querySelector('.jqx-notification-content'))
    this.jqxNotification.elementRef.nativeElement.querySelector('.jqx-notification-content').innerHTML = msg;
    // console.log(this.jqxNotification.elementRef.nativeElement.querySelector('.notificationContainer'))
    // this.jqxNotification.elementRef.nativeElement.querySelector('.notificationContainer').classList.add(level)
    this.jqxNotification.open();
};
}