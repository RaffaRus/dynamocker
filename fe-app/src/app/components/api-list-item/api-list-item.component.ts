import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { MockApi, dummyMockApi } from '@models/mockApi';
import { isDefined } from '@common/common';
import { MockApiService } from '@services/mockApiService';
import { Observable, Subscription } from 'rxjs';
import { subscriptionLogsToBeFn } from 'rxjs/internal/testing/TestScheduler';

@Component({
  selector: 'app-api-list-item',
  standalone: true,
  templateUrl: './api-list-item.component.html',
  styleUrl: './api-list-item.component.scss'
})
export class ApiListItemComponent implements OnInit{
  @Output() click = new EventEmitter<string>()

  private mockApi_ : MockApi = dummyMockApi

  constructor(
    private mockApiService : MockApiService
  ){}

  ngOnInit(): void {
    // check that the mockApi has been correctly filled with Input value
    if (this.mockApi_ === dummyMockApi) {
      console.error("mock Api not correctly passed in from the parent component");
    }

  }

  onClick(){
    this.mockApiService.selectMockApi(this.mockApi_)
  }
  
  public onDeleteItem() {
    console.log("item removed")
    // delete mock api from the back end
    this.mockApiService.deleteAllMockApis(this.mockApi_.Name).subscribe({
      next: (value) => {
        console.log(value)
        // emit the requirement to refresh list
        this.mockApiService.refreshList()
      },
      error: (err) => {
        console.error(err)
      }
    })
  }

}
