// General
import { Component, OnInit } from '@angular/core';

// Model
import { MockApi } from '@models/mockApi';

// Service
import { MockApiService } from '@services/mockApiService';
import { isDefined } from '@common/common';

@Component({
  selector: 'app-api-list',
  standalone: true,
  imports: [],
  templateUrl: './api-list.component.html',
  styleUrl: './api-list.component.css'
})

export class ApiListComponent implements OnInit{ 

  // list of mock apis
  apiList_: MockApi[] = []

  constructor(
    private mockApiService : MockApiService
  ){}

  ngOnInit(): void {

    // get the mock apis for the first time
    this.mockApiService.getAllMockApis().subscribe({
      next: (value) => {
        this.apiList_ = value
      },
      error: (err) => {
        console.error(err)
      }
    })

    // subscribe to the refresh list observable
    this.mockApiService.refreshListObservable().subscribe({
      next: () => {
        // if a new event has been emitted, then get againg the mock api list
        this.mockApiService.getAllMockApis().subscribe({
          next: (value) => {
            this.apiList_ = value
          },
          error: (err) => {
            console.error(err)
          }
        })
      },
      error: (err) => {
        console.error(err)
      }
    })

  } 

  // new Item
  newItem() {
    // create an empty element, add it to the list and open it in the editor
  }

}
