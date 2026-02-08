export type AppScreen = "drop" | "gallery";
export type Theme = "light" | "dark";

export interface TemplateInfo {
  name: string;
  displayName: string;
  format: string;
  description: string;
}

export interface ParseResult {
  name: string;
  email: string;
  format: string;
}
