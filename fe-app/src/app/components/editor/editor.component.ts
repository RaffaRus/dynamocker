import { Component, OnInit } from '@angular/core';
import { IMockApi, dummyMockApi } from '@models/mockApi';
import { MockApiService } from '@services/mockApiService';

@Component({
  selector: 'app-editor',
  templateUrl: './editor.component.html',
  styleUrl: './editor.component.css'
})
export class EditorComponent implements OnInit{

  private selectedMockApi_ : IMockApi = dummyMockApi

  constructor(
    private mockApiService :  MockApiService
  ){}

  ngOnInit(): void {
    
    this.mockApiService.newMockApiSelectedObservable().subscribe({
      next: (mockApi : IMockApi) => {
        this.selectedMockApi_= mockApi
      },
      error: (err: Error) => {
        console.error("received error from newMockApiSelectedObservable subscription: ", err)
      }
    })
  }

}
