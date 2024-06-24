import { Component, OnInit } from '@angular/core';
import { MockApi } from '@models/mockApi';

@Component({
  selector: 'app-api-list',
  standalone: true,
  imports: [],
  templateUrl: './api-list.component.html',
  styleUrl: './api-list.component.css'
})
export class ApiListComponent implements OnInit {

  private apiList_ []MockApi

  constructor(
    private MockApiService
  ): {
    apiList_ = this.MockApiService.GetAllMockApis()

  }


}
