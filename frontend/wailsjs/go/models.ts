export namespace models {
	
	export class LogContent {
	    lines: string[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new LogContent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lines = source["lines"];
	        this.total = source["total"];
	    }
	}
	export class LogFile {
	    id: string;
	    name: string;
	    path: string;
	    size: number;
	    // Go type: time
	    lastModified: any;
	    isOpen: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LogFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.lastModified = this.convertValues(source["lastModified"], null);
	        this.isOpen = source["isOpen"];
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
	export class SearchResult {
	    lineNumber: number;
	    content: string;
	    matches: number[];
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lineNumber = source["lineNumber"];
	        this.content = source["content"];
	        this.matches = source["matches"];
	    }
	}
	export class SystemInfo {
	    cpuUsage: number;
	    memoryUsage: number;
	    memoryUsed: number;
	    memoryTotal: number;
	    memoryAvailable: number;
	    appMemoryUsage: number;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cpuUsage = source["cpuUsage"];
	        this.memoryUsage = source["memoryUsage"];
	        this.memoryUsed = source["memoryUsed"];
	        this.memoryTotal = source["memoryTotal"];
	        this.memoryAvailable = source["memoryAvailable"];
	        this.appMemoryUsage = source["appMemoryUsage"];
	    }
	}

}

