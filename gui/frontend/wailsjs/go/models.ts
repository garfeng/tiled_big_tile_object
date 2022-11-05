export namespace main {
	
	export class Param {
	    tileSize: number;
	    dstWidth: number;
	    dstHeight': number;
	    srcImages: string[];
	    dstRoot: string;
	    dstPrefix: string;
	
	    static createFrom(source: any = {}) {
	        return new Param(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tileSize = source["tileSize"];
	        this.dstWidth = source["dstWidth"];
	        this.dstHeight' = source["dstHeight'"];
	        this.srcImages = source["srcImages"];
	        this.dstRoot = source["dstRoot"];
	        this.dstPrefix = source["dstPrefix"];
	    }
	}

}

