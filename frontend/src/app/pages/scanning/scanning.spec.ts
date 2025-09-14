import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Scanning } from './scanning';

describe('Scanning', () => {
  let component: Scanning;
  let fixture: ComponentFixture<Scanning>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Scanning]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Scanning);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
