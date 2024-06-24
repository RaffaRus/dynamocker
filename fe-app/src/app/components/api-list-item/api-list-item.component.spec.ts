import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ApiListItemComponent } from './api-list-item.component';

describe('ApiListItemComponent', () => {
  let component: ApiListItemComponent;
  let fixture: ComponentFixture<ApiListItemComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ApiListItemComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ApiListItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
