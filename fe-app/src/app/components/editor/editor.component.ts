import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { IMockApi, dummyMockApi } from '@models/mockApi';
import { initialMockApiJson } from '@models/monaco';
import { MockApiService } from '@services/mockApiService';
import * as monaco from 'monaco-editor';

@Component({
  selector: 'app-editor',
  templateUrl: './editor.component.html',
  styleUrl: './editor.component.scss'
})
export class EditorComponent implements OnInit {
  @ViewChild('editorContainer', { static: true }) _editorContainer!: ElementRef;
  codeEditorInstance!: monaco.editor.IStandaloneCodeEditor;

  private selectedMockApi_: IMockApi = dummyMockApi

  constructor(
    private mockApiService: MockApiService
  ) {
    this.assignDefaultJsonSchema()
  }

  ngOnInit(): void {


    this.codeEditorInstance = monaco.editor.create(this._editorContainer.nativeElement, {
      theme: 'vs',
      wordWrap: 'on',
      wrappingIndent: 'indent',
      language: 'json',
      // minimap: { enabled: false },
      automaticLayout: true,
      value: initialMockApiJson
    });

    this.mockApiService.newMockApiSelectedObservable().subscribe({
      next: (mockApi: IMockApi) => {
        this.selectedMockApi_ = mockApi
      },
      error: (err: Error) => {
        console.error("received error from newMockApiSelectedObservable subscription: ", err)
      }
    })
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
}