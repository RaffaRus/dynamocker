import { Component, OnInit } from '@angular/core';
import { MockApi, dummyMockApi } from '@models/mockApi';
import { MockApiService } from '@services/mockApiService';

@Component({
  selector: 'app-editor',
  standalone: true,
  imports: [],
  templateUrl: './editor.component.html',
  styleUrl: './editor.component.css'
})
export class EditorComponent implements OnInit{

  private selectedMockApi_ : MockApi = dummyMockApi

  constructor(
    private mockApiService :  MockApiService
  ){}

  ngOnInit(): void {
    
    this.mockApiService.newMockApiSelectedObservable().subscribe({
      next: (mockApi : MockApi) => {
        this.selectedMockApi_= mockApi
      },
      error: (err: Error) => {
        console.error("received error from newMockApiSelectedObservable subscription: ", err)
      }
    })
  }

}
