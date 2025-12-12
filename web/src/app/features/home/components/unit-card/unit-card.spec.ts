import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UnitCard } from './unit-card';

describe('UnitCard', () => {
  let component: UnitCard;
  let fixture: ComponentFixture<UnitCard>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [UnitCard]
    })
    .compileComponents();

    fixture = TestBed.createComponent(UnitCard);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
