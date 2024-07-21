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

  startingMockApiList : IMockApi[] = [
    {
      added: new Date(),
      name: "lol1",
      url: "url",
      responses : {
        delete: JSON.parse("{\"deleted1\":\"deleted1\"}"),
        get: JSON.parse("{\"get1\":\"get1\"}"),
        post: JSON.parse("{\"post1\":\"post1\"}"),
        patch: JSON.parse("{\"patch1\":\"patch1\"}"),
      },
    },
    {
      added: new Date(),
      name: "lol2",
      url: "ur2",
      responses : {
        delete: JSON.parse("{\"deleted2\":\"deleted2\"}"),
        get: JSON.parse("{\"get2\":\"get2\"}"),
        post: JSON.parse("{\"post2\":\"post2\"}"),
        patch: JSON.parse("{\"patch2\":\"patch2\"}"),
      },
    },
    {
      added: new Date(),
      name: "lol2",
      url: "ur2",
      responses : {
        delete: JSON.parse("{\"deleted2\":\"deleted2\"}"),
        get: JSON.parse("{\"get2\":\"get2\"}"),
        post: JSON.parse("{\"post2\":\"post2\"}"),
        patch: JSON.parse("{\"patch2\":\"patch2\"}"),
      },
    },
    {
      added: new Date(),
      name: "lol2",
      url: "ur2",
      responses : {
        delete: JSON.parse("{\"deleted2\":\"deleted2\"}"),
        get: JSON.parse("{\"get2\":\"get2\"}"),
        post: JSON.parse("{\"post2\":\"post2\"}"),
        patch: JSON.parse("{\"patch2\":\"patch2\"}"),
      },
    },
    {
      added: new Date(),
      name: "lol2",
      url: "ur2",
      responses : {
        delete: JSON.parse("{\"deleted2\":\"deleted2\"}"),
        get: JSON.parse("{\"get2\":\"get2\"}"),
        post: JSON.parse("{\"post2\":\"post2\"}"),
        patch: JSON.parse("{\"patch2\":\"patch2\"}"),
      },
    },
    {
      added: new Date(),
      name: "lol2",
      url: "ur2",
      responses : {
        delete: JSON.parse("{\"deleted2\":\"deleted2\"}"),
        get: JSON.parse("{\"get2\":\"get2\"}"),
        post: JSON.parse("{\"post2\":\"post2\"}"),
        patch: JSON.parse("{\"patch2\":\"patch2\"}"),
      },
    },
    {
      added: new Date(),
      name: "lol2",
      url: "ur2",
      responses : {
        delete: JSON.parse("{\"deleted2\":\"deleted2\"}"),
        get: JSON.parse("{\"get2\":\"get2\"}"),
        post: JSON.parse("{\"post2\":\"post2\"}"),
        patch: JSON.parse("{\"patch2\":\"patch2\"}"),
      },
    },
    {
      added: new Date(),
      name: "lol2",
      url: "ur2",
      responses : {
        delete: JSON.parse("{\"deleted2\":\"deleted2\"}"),
        get: JSON.parse("{\"get2\":\"get2\"}"),
        post: JSON.parse("{\"post2\":\"post2\"}"),
        patch: JSON.parse("{\"patch2\":\"patch2\"}"),
      },
    },
  ]

  // list of mock apis
  _apiList: IMockApi[] = []

  constructor(
    private mockApiService : MockApiService
  ){}

  ngOnInit(): void {

    // get the mock apis for the first time
    // this.mockApiService.getAllMockApis().subscribe({
    //   next: (value) => {
    //     this._apiList = value
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
            this._apiList = value
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
