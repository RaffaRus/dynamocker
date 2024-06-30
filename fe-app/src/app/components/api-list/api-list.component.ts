// General
import { Component, OnInit, setTestabilityGetter } from '@angular/core';

// Model
import { IMockApi } from '@models/mockApi';

// Service
import { MockApiService } from '@services/mockApiService';

@Component({
  selector: 'app-api-list',
  templateUrl: './api-list.component.html',
  styleUrl: './api-list.component.css'
})

export class ApiListComponent implements OnInit{ 

  startingMockApiList : IMockApi[] = [
    {
      Added: new Date(),
      LastModified: new Date(),
      Name: "lol1",
      URL: "url",
      Responses : {
        Delete: JSON.parse("{\"delete1\":\"delete1\"}"),
        Get: JSON.parse("{\"get1\":\"get1\"}"),
        Post: JSON.parse("{\"post1\":\"post1\"}"),
        Patch: JSON.parse("{\"patch1\":\"patch1\"}"),
      },
    },
    {
      Added: new Date(),
      LastModified: new Date(),
      Name: "lol2",
      URL: "ur2",
      Responses : {
        Delete: JSON.parse("{\"delete2\":\"delete2\"}"),
        Get: JSON.parse("{\"get2\":\"get2\"}"),
        Post: JSON.parse("{\"post2\":\"post2\"}"),
        Patch: JSON.parse("{\"patch2\":\"patch2\"}"),
      },
    },
  ]

  // list of mock apis
  apiList_: IMockApi[] = this.startingMockApiList

  constructor(
    private mockApiService : MockApiService
  ){}

  ngOnInit(): void {

    // get the mock apis for the first time
    // this.mockApiService.getAllMockApis().subscribe({
    //   next: (value) => {
    //     this.apiList_ = value
    //   },
    //   error: (err) => {
    //     console.error(err)
    //   }
    // })

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
