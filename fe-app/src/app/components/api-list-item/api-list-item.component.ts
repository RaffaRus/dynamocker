import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { ResourceObject } from '@models/mockApi.model';
import { MockApiService } from '@services/mockApiService';
import { remove } from "lodash";

@Component({
  selector: 'app-api-list-item',
  templateUrl: './api-list-item.component.html',
  styleUrl: './api-list-item.component.scss'
})
export class ApiListItemComponent implements OnInit{
  @Output() click = new EventEmitter<string>()
  @Input() resObj : ResourceObject = new ResourceObject()

  constructor(
    private mockApiService : MockApiService
  ){}

  ngOnInit(): void {
    // check that the mockApi has been correctly filled with Input value
    if (this.resObj === new ResourceObject()) {
      console.error("mock Api not correctly passed in from the parent component");
    }

  }

  onSelectedItem(){
    // TODO: check that the mockApi in the editor is saved, otherwise raise notitfication
    let newResObj = new ResourceObject()
    newResObj.data = this.resObj.data
    newResObj.id = this.resObj.id
    newResObj.type = this.resObj.type
    this.mockApiService.selectMockApi(newResObj)
  }
  
  mockApiResponses(): string[] {
    // create the array first, then remove items if they are not defined in the mock api
    let responses : string[] = ["GET","POST","PATCH","DELETE","PUT"]
    if (JSON.stringify(this.resObj.data.responses.get) === '{}' ||
        this.resObj.data.responses.get === undefined) {remove(responses, m => {return m === "GET"})}
    if (JSON.stringify(this.resObj.data.responses.post) === '{}' ||
        this.resObj.data.responses.post === undefined) {remove(responses, m => {return m === "POST"})}
    if (JSON.stringify(this.resObj.data.responses.patch) === '{}' ||
        this.resObj.data.responses.patch === undefined) {remove(responses, m => {return m === "PATCH"})}
    if (JSON.stringify(this.resObj.data.responses.delete) === '{}' ||
        this.resObj.data.responses.delete === undefined) {remove(responses, m => {return m ==="DELETE"})}
    if (JSON.stringify(this.resObj.data.responses.put) === '{}' ||
        this.resObj.data.responses.put === undefined) {remove(responses, m => {return m ==="PUT"})} else {console.log(JSON.stringify(this.resObj.data.responses.put))}
    return responses
  }

  public onDeleteItem() {
    // delete mock api from the back end
    this.mockApiService.deleteMockApi(this.resObj.id).subscribe({
      next: (_) => {
        // TODO: create service to manage notifications and add notification about the delete response
        // emit the requirement to refresh list
        this.mockApiService.refreshList()
      },
      error: (err) => {
        console.error(err)
      }
    })
  }

}
