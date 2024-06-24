import { Component } from '@angular/core';

@Component({
  selector: 'app-api-list-item',
  standalone: true,
  templateUrl: './api-list-item.component.html',
  styleUrl: './api-list-item.component.scss'
})
export class ApiListItemComponent {
  
  public deleteItem() {
    console.log("delete item")
  }
}

