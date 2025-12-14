export interface Unit {
  slug: string;
  title: string;
  subtitle?: string;
  description: string;
  imageSrc?: string;
  imageAlt?: string;
}

export const UNITS: Unit[] = [
  {
    slug: 'unit-1-intro',
    title: 'Introducción a la Programación',
    subtitle: 'Unidad 1',
    description: 'Aprende los fundamentos de Python',
    imageSrc: '/assets/units/unit-1.png',
    imageAlt: 'Unidad 1',
  },
];
