export namespace main {
	
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

