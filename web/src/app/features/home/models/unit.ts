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
    imageSrc: '/assets/units/unit-1.webp',
    imageAlt: 'Unidad 1',
  },
  {
    slug: 'unit-2-variables',
    title: 'Variables y Tipos de Datos',
    subtitle: 'Unidad 2',
    description: 'Primitivos, variables, E/S, strings y listas.',
    imageSrc: '/assets/units/unit-2.webp',
    imageAlt: 'Unidad 2',
  },
  {
    slug: 'unit-3-funciones',
    title: 'Funciones',
    subtitle: 'Unidad 3',
    description: 'Abstracción, parámetros, alcance y modularización.',
    imageSrc: '/assets/units/unit-3.webp',
    imageAlt: 'Unidad 3',
  },
  {
    slug: 'unit-4-control',
    title: 'Estructuras de Control',
    subtitle: 'Unidad 4',
    description: 'Condicionales, iterativas y sentencias anidadas.',
    imageSrc: '/assets/units/unit-4.webp',
    imageAlt: 'Unidad 4',
  },
];
