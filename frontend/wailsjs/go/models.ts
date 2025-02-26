export namespace main {
	
	export class DownloadProgressInfo {
	    progress: number;
	    downloadSpeed: number;
	    totalSize: number;
	    completedSize: number;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new DownloadProgressInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.progress = source["progress"];
	        this.downloadSpeed = source["downloadSpeed"];
	        this.totalSize = source["totalSize"];
	        this.completedSize = source["completedSize"];
	        this.status = source["status"];
	    }
	}

}

