import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { isDefined } from '@common/common';
import { initialMockApiJson } from '@models/editor.model';
import { IMockApi } from '@models/mockApi.model';
import { MockApiService } from '@services/mockApiService';

@Component({
  selector: 'app-api-list-item',
  templateUrl: './api-list-item.component.html',
  styleUrl: './api-list-item.component.scss'
})
export class ApiListItemComponent implements OnInit{
  @Output() click = new EventEmitter<string>()
  @Input() mockApi : IMockApi = initialMockApiJson

  constructor(
    private mockApiService : MockApiService
  ){}

  ngOnInit(): void {
    // check that the mockApi has been correctly filled with Input value
    if (this.mockApi === initialMockApiJson) {
      console.error("mock Api not correctly passed in from the parent component");
    }

  }

  onSelectedItem(){
    // check that the mockApi in the editor is saved, otherwise raise notitfication

    this.mockApiService.selectMockApi(this.mockApi)
  }
  
  mockApiResponses(): string[] {
    // create the array first, then remove items if they are not defined in the mock api
    let responses : string[] = ["GET","POST","PATCH","DELETE"]
    if (!isDefined(this.mockApi.Responses.Get)) {delete responses[responses.findIndex((name: string) => {return name==="Get"})]}
    if (!isDefined(this.mockApi.Responses.Post)) {delete responses[responses.findIndex((name: string) => {return name==="Post"})]}
    if (!isDefined(this.mockApi.Responses.Patch)) {delete responses[responses.findIndex((name: string) => {return name==="Patch"})]}
    if (!isDefined(this.mockApi.Responses.Delete)) {delete responses[responses.findIndex((name: string) => {return name==="Delete"})]}
    return responses
  }

  public onDeleteItem() {
    console.log("item removed")
    // delete mock api from the back end
    this.mockApiService.deleteAllMockApis(this.mockApi.Name).subscribe({
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
