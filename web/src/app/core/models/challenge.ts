export interface Challenge {
  id: number;
  slug: string;
  category: string;
  title: string;
  description: string;
  template: string;
  test_code: string;
  hints: string;
}
