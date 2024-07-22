// General
import { Component, OnInit } from '@angular/core';

// Model
import { IMockApi } from '@models/mockApi.model';
import { initialMockApiJson } from '@models/editor.model';

// Service
import { MockApiService } from '@services/mockApiService';

@Component({
  selector: 'app-api-list',
  templateUrl: './api-list.component.html',
  styleUrl: './api-list.component.scss'
})

export class ApiListComponent implements OnInit{ 

  // list of mock apis
  _apiList: IMockApi[] = []

  constructor(
    private mockApiService : MockApiService
  ){}

  ngOnInit(): void {

    // get the mock apis for the first time
    this.mockApiService.getAllMockApis().subscribe({
      next: (value) => {
        this._apiList = value
      },
      error: (err) => {
        console.error(err)
      }
    })

    // subscribe to the refresh list observable
    this.mockApiService.refreshListObservable().subscribe({
      next: () => {
        console.log("requested refresh mock api list")
        // if a new event has been emitted, then get againg the mock api list
        this.mockApiService.getAllMockApis().subscribe({
          next: (value) => {
            console.log(value)
            this._apiList = value
            console.log(this._apiList)
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
    // TODO: add notification with confirmation if smt in the editor is modified and not saved
    // create an empty element, add it to the list and open it in the editor
    let newMockApiModel = initialMockApiJson
    this.mockApiService.selectMockApi(newMockApiModel)
  }

}
