import { Component } from '@angular/core';
import { UNITS } from './models/unit';
import { Card } from '../../shared/components/card/card';

@Component({
  selector: 'app-home',
  imports: [Card],
  templateUrl: './home.html',
  styleUrl: './home.scss',
})
export class Home {
  units = UNITS;
}
