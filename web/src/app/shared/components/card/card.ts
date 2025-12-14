import { Component, input } from '@angular/core';
import { RouterLink } from '@angular/router';

type RouterLinkTarget = string | any[];

@Component({
  selector: 'app-card',
  standalone: true,
  imports: [RouterLink],
  templateUrl: './card.html',
  styleUrl: './card.scss',
})
export class Card {
  title = input.required<string>();
  subtitle = input<string>();
  description = input.required<string>();
  link = input.required<RouterLinkTarget>();
  imageSrc = input<string>();
  imageAlt = input<string>('');
}
