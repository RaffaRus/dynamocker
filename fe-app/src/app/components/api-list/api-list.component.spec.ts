import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ApiListComponent } from './api-list.component';

describe('ApiListComponent', () => {
  let component: ApiListComponent;
  let fixture: ComponentFixture<ApiListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ApiListComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ApiListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
