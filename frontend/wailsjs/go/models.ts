export namespace main {

	export class GeneratePDFResult {
	    data: string;
	    pageCount: number;

	    static createFrom(source: any = {}) {
	        return new GeneratePDFResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = source["data"];
	        this.pageCount = source["pageCount"];
	    }
	}
	export class ParseResult {
	    name: string;
	    email: string;
	    format: string;
	
	    static createFrom(source: any = {}) {
	        return new ParseResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.email = source["email"];
	        this.format = source["format"];
	    }
	}
	export class TemplateInfo {
	    name: string;
	    displayName: string;
	    format: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new TemplateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.displayName = source["displayName"];
	        this.format = source["format"];
	        this.description = source["description"];
	    }
	}

}

export namespace resume {
	
	export class Award {
	    name: string;
	    // Go type: time
	    date?: any;
	    notes?: string;
	
	    static createFrom(source: any = {}) {
	        return new Award(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.date = this.convertValues(source["date"], null);
	        this.notes = source["notes"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Certification {
	    name: string;
	    issuer?: string;
	    notes?: string;
	    // Go type: time
	    date?: any;
	
	    static createFrom(source: any = {}) {
	        return new Certification(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.issuer = source["issuer"];
	        this.notes = source["notes"];
	        this.date = this.convertValues(source["date"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Certifications {
	    title?: string;
	    items: Certification[];
	
	    static createFrom(source: any = {}) {
	        return new Certifications(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.items = this.convertValues(source["items"], Certification);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Link {
	    uri: string;
	    label?: string;
	
	    static createFrom(source: any = {}) {
	        return new Link(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uri = source["uri"];
	        this.label = source["label"];
	    }
	}
	export class Location {
	    city: string;
	    state?: string;
	    province?: string;
	    country?: string;
	    remote?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Location(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.city = source["city"];
	        this.state = source["state"];
	        this.province = source["province"];
	        this.country = source["country"];
	        this.remote = source["remote"];
	    }
	}
	export class Contact {
	    name: string;
	    email: string;
	    phone?: string;
	    location?: Location;
	    links?: Link[];
	
	    static createFrom(source: any = {}) {
	        return new Contact(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.email = source["email"];
	        this.phone = source["phone"];
	        this.location = this.convertValues(source["location"], Location);
	        this.links = this.convertValues(source["links"], Link);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DateRange {
	    // Go type: time
	    start: any;
	    // Go type: time
	    end?: any;
	
	    static createFrom(source: any = {}) {
	        return new DateRange(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.start = this.convertValues(source["start"], null);
	        this.end = this.convertValues(source["end"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Degree {
	    name: string;
	    descriptions?: string[];
	
	    static createFrom(source: any = {}) {
	        return new Degree(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.descriptions = source["descriptions"];
	    }
	}
	export class Thesis {
	    title: string;
	    highlights?: string[];
	    link: Link;
	    description?: string;
	
	    static createFrom(source: any = {}) {
	        return new Thesis(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.highlights = source["highlights"];
	        this.link = this.convertValues(source["link"], Link);
	        this.description = source["description"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GPA {
	    gpa?: string;
	    max_gpa?: string;
	
	    static createFrom(source: any = {}) {
	        return new GPA(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gpa = source["gpa"];
	        this.max_gpa = source["max_gpa"];
	    }
	}
	export class Education {
	    institution: string;
	    degree: Degree;
	    specializations?: string[];
	    gpa?: GPA;
	    awards?: Award[];
	    dates: DateRange;
	    location?: Location;
	    thesis?: Thesis;
	
	    static createFrom(source: any = {}) {
	        return new Education(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.institution = source["institution"];
	        this.degree = this.convertValues(source["degree"], Degree);
	        this.specializations = source["specializations"];
	        this.gpa = this.convertValues(source["gpa"], GPA);
	        this.awards = this.convertValues(source["awards"], Award);
	        this.dates = this.convertValues(source["dates"], DateRange);
	        this.location = this.convertValues(source["location"], Location);
	        this.thesis = this.convertValues(source["thesis"], Thesis);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EducationList {
	    title?: string;
	    institutions: Education[];
	
	    static createFrom(source: any = {}) {
	        return new EducationList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.institutions = this.convertValues(source["institutions"], Education);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Experience {
	    company: string;
	    title: string;
	    employment_type?: string;
	    highlights?: string[];
	    duties?: string[];
	    notes?: string;
	    dates: DateRange;
	    location?: Location;
	    technologies?: string[];
	
	    static createFrom(source: any = {}) {
	        return new Experience(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.company = source["company"];
	        this.title = source["title"];
	        this.employment_type = source["employment_type"];
	        this.highlights = source["highlights"];
	        this.duties = source["duties"];
	        this.notes = source["notes"];
	        this.dates = this.convertValues(source["dates"], DateRange);
	        this.location = this.convertValues(source["location"], Location);
	        this.technologies = source["technologies"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ExperienceList {
	    title?: string;
	    positions: Experience[];
	
	    static createFrom(source: any = {}) {
	        return new ExperienceList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.positions = this.convertValues(source["positions"], Experience);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Language {
	    name: string;
	    proficiency?: string;
	
	    static createFrom(source: any = {}) {
	        return new Language(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.proficiency = source["proficiency"];
	    }
	}
	export class LanguageList {
	    title?: string;
	    languages: Language[];
	
	    static createFrom(source: any = {}) {
	        return new LanguageList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.languages = this.convertValues(source["languages"], Language);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Layout {
	    density?: string;
	    typography?: string;
	    header?: string;
	    sections?: string[];
	    skill_columns?: number;
	    references?: boolean;

	    static createFrom(source: any = {}) {
	        return new Layout(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.density = source["density"];
	        this.typography = source["typography"];
	        this.header = source["header"];
	        this.sections = source["sections"];
	        this.skill_columns = source["skill_columns"];
	        this.references = source["references"];
	    }
	}
	
	
	export class Project {
	    name: string;
	    link?: Link;
	    highlights?: string[];
	    dates?: DateRange;
	    technologies?: string[];
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.link = this.convertValues(source["link"], Link);
	        this.highlights = source["highlights"];
	        this.dates = this.convertValues(source["dates"], DateRange);
	        this.technologies = source["technologies"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ProjectList {
	    title?: string;
	    projects: Project[];
	
	    static createFrom(source: any = {}) {
	        return new ProjectList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.projects = this.convertValues(source["projects"], Project);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SkillCategory {
	    category: string;
	    items: string[];
	
	    static createFrom(source: any = {}) {
	        return new SkillCategory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category = source["category"];
	        this.items = source["items"];
	    }
	}
	export class Skills {
	    title?: string;
	    categories: SkillCategory[];
	
	    static createFrom(source: any = {}) {
	        return new Skills(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.categories = this.convertValues(source["categories"], SkillCategory);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Resume {
	    contact: Contact;
	    summary?: string;
	    certifications?: Certifications;
	    skills: Skills;
	    experience: ExperienceList;
	    projects?: ProjectList;
	    education: EducationList;
	    languages?: LanguageList;
	    layout?: Layout;
	
	    static createFrom(source: any = {}) {
	        return new Resume(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.contact = this.convertValues(source["contact"], Contact);
	        this.summary = source["summary"];
	        this.certifications = this.convertValues(source["certifications"], Certifications);
	        this.skills = this.convertValues(source["skills"], Skills);
	        this.experience = this.convertValues(source["experience"], ExperienceList);
	        this.projects = this.convertValues(source["projects"], ProjectList);
	        this.education = this.convertValues(source["education"], EducationList);
	        this.languages = this.convertValues(source["languages"], LanguageList);
	        this.layout = this.convertValues(source["layout"], Layout);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	export class ValidationError {
	    field: string;
	    message: string;
	    type: string;
	    value?: any;
	
	    static createFrom(source: any = {}) {
	        return new ValidationError(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.field = source["field"];
	        this.message = source["message"];
	        this.type = source["type"];
	        this.value = source["value"];
	    }
	}

}

