import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { UNITS } from './models/unit';
import { UnitCard } from './components/unit-card/unit-card';

@Component({
  selector: 'app-home',
  imports: [RouterLink, UnitCard],
  templateUrl: './home.html',
  styleUrl: './home.scss',
})
export class Home {
  units = UNITS
}
