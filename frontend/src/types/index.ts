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

// Resume data types matching Go structs

export interface Location {
  city: string;
  state?: string;
  province?: string;
  country?: string;
  remote?: boolean;
}

export interface Link {
  uri: string;
  label?: string;
}

export interface Contact {
  name: string;
  email: string;
  phone?: string;
  location?: Location;
  links?: Link[];
}

export interface SkillCategory {
  category: string;
  items: string[];
}

export interface Skills {
  title?: string;
  categories: SkillCategory[];
}

export interface DateRange {
  start: string; // ISO 8601
  end?: string;  // ISO 8601, absent means "Present"
}

export interface Experience {
  company: string;
  title: string;
  employment_type?: string;
  highlights?: string[];
  duties?: string[];
  notes?: string;
  dates: DateRange;
  location?: Location;
  technologies?: string[];
}

export interface ExperienceList {
  title?: string;
  positions: Experience[];
}

export interface Degree {
  name: string;
  descriptions?: string[];
}

export interface GPA {
  gpa?: string;
  max_gpa?: string;
}

export interface Award {
  name: string;
  date?: string;
  notes?: string;
}

export interface Thesis {
  title: string;
  highlights?: string[];
  link: Link;
  description?: string;
}

export interface Education {
  institution: string;
  degree: Degree;
  specializations?: string[];
  gpa?: GPA;
  awards?: Award[];
  dates: DateRange;
  location?: Location;
  thesis?: Thesis;
}

export interface EducationList {
  title?: string;
  institutions: Education[];
}

export interface Project {
  name: string;
  link?: Link;
  highlights?: string[];
  dates?: DateRange;
  technologies?: string[];
}

export interface ProjectList {
  title?: string;
  projects: Project[];
}

export interface Certification {
  name: string;
  issuer?: string;
  notes?: string;
  date?: string;
}

export interface Certifications {
  title?: string;
  items: Certification[];
}

export interface Language {
  name: string;
  proficiency?: string;
}

export interface LanguageList {
  title?: string;
  languages: Language[];
}

export interface Layout {
  density?: string;
  typography?: string;
  header?: string;
  sections?: string[];
  skill_columns?: number;
}

export interface Resume {
  contact: Contact;
  summary?: string;
  certifications?: Certifications;
  skills: Skills;
  experience: ExperienceList;
  projects?: ProjectList;
  education: EducationList;
  languages?: LanguageList;
  layout?: Layout;
}

export interface ValidationError {
  field: string;
  message: string;
  type: string;
  value?: unknown;
}
