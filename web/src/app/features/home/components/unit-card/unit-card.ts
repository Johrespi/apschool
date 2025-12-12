import { Component, input } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import {Unit} from '../../models/unit';

@Component({
  selector: 'app-unit-card',
  imports: [MatCardModule],
  templateUrl: './unit-card.html',
  styleUrl: './unit-card.scss',
})
export class UnitCard {
  unit = input.required<Unit>();
}
