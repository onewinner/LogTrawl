export namespace models {
	
	export class URLStat {
	    rank: number;
	    url: string;
	    count: number;
	    percentage: number;
	
	    static createFrom(source: any = {}) {
	        return new URLStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rank = source["rank"];
	        this.url = source["url"];
	        this.count = source["count"];
	        this.percentage = source["percentage"];
	    }
	}
	export class URLAnalysis {
	    getTop10: URLStat[];
	    postTop10: URLStat[];
	
	    static createFrom(source: any = {}) {
	        return new URLAnalysis(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.getTop10 = this.convertValues(source["getTop10"], URLStat);
	        this.postTop10 = this.convertValues(source["postTop10"], URLStat);
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
	export class IPStat {
	    rank: number;
	    ip: string;
	    count: number;
	    firstAccess: string;
	    lastAccess: string;
	
	    static createFrom(source: any = {}) {
	        return new IPStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rank = source["rank"];
	        this.ip = source["ip"];
	        this.count = source["count"];
	        this.firstAccess = source["firstAccess"];
	        this.lastAccess = source["lastAccess"];
	    }
	}
	export class IPAnalysis {
	    internalTop10: IPStat[];
	    externalTop10: IPStat[];
	
	    static createFrom(source: any = {}) {
	        return new IPAnalysis(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.internalTop10 = this.convertValues(source["internalTop10"], IPStat);
	        this.externalTop10 = this.convertValues(source["externalTop10"], IPStat);
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
	export class OverviewStats {
	    totalRequests: number;
	    uniqueIPs: number;
	    internalIPs: number;
	    externalIPs: number;
	
	    static createFrom(source: any = {}) {
	        return new OverviewStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalRequests = source["totalRequests"];
	        this.uniqueIPs = source["uniqueIPs"];
	        this.internalIPs = source["internalIPs"];
	        this.externalIPs = source["externalIPs"];
	    }
	}
	export class AnalysisResult {
	    overview: OverviewStats;
	    ipAnalysis: IPAnalysis;
	    urlAnalysis: URLAnalysis;
	
	    static createFrom(source: any = {}) {
	        return new AnalysisResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.overview = this.convertValues(source["overview"], OverviewStats);
	        this.ipAnalysis = this.convertValues(source["ipAnalysis"], IPAnalysis);
	        this.urlAnalysis = this.convertValues(source["urlAnalysis"], URLAnalysis);
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
	export class DatePattern {
	    name: string;
	    pattern: string;
	    example: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new DatePattern(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.pattern = source["pattern"];
	        this.example = source["example"];
	        this.description = source["description"];
	    }
	}
	export class DateRange {
	    // Go type: time
	    startDate?: any;
	    // Go type: time
	    endDate?: any;
	    datePattern: string;
	    totalDays: number;
	    hasDateInfo: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DateRange(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.startDate = this.convertValues(source["startDate"], null);
	        this.endDate = this.convertValues(source["endDate"], null);
	        this.datePattern = source["datePattern"];
	        this.totalDays = source["totalDays"];
	        this.hasDateInfo = source["hasDateInfo"];
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
	export class FileSplitOptions {
	    strategy: string;
	    filePath: string;
	    outputDir: string;
	    datePattern?: string;
	    daysPerFile?: number;
	    sizePerFile?: number;
	    linesPerFile?: number;
	
	    static createFrom(source: any = {}) {
	        return new FileSplitOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.strategy = source["strategy"];
	        this.filePath = source["filePath"];
	        this.outputDir = source["outputDir"];
	        this.datePattern = source["datePattern"];
	        this.daysPerFile = source["daysPerFile"];
	        this.sizePerFile = source["sizePerFile"];
	        this.linesPerFile = source["linesPerFile"];
	    }
	}
	export class FileSplitResult {
	    success: boolean;
	    message: string;
	    outputFiles: string[];
	    totalFiles: number;
	    totalSize: number;
	
	    static createFrom(source: any = {}) {
	        return new FileSplitResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.outputFiles = source["outputFiles"];
	        this.totalFiles = source["totalFiles"];
	        this.totalSize = source["totalSize"];
	    }
	}
	
	
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
	    lines?: number;
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
	        this.lines = source["lines"];
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
	export class SpecificIPURLAnalysis {
	    getTop10: URLStat[];
	    postTop10: URLStat[];
	
	    static createFrom(source: any = {}) {
	        return new SpecificIPURLAnalysis(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.getTop10 = this.convertValues(source["getTop10"], URLStat);
	        this.postTop10 = this.convertValues(source["postTop10"], URLStat);
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
	export class SpecificIPResult {
	    ip: string;
	    ipType: string;
	    totalRequests: number;
	    firstAccess: string;
	    lastAccess: string;
	    urlAnalysis: SpecificIPURLAnalysis;
	
	    static createFrom(source: any = {}) {
	        return new SpecificIPResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = source["ip"];
	        this.ipType = source["ipType"];
	        this.totalRequests = source["totalRequests"];
	        this.firstAccess = source["firstAccess"];
	        this.lastAccess = source["lastAccess"];
	        this.urlAnalysis = this.convertValues(source["urlAnalysis"], SpecificIPURLAnalysis);
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

