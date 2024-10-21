export namespace main {
	
	export class Cmyk {
	    c: number;
	    m: number;
	    y: number;
	    k: number;
	
	    static createFrom(source: any = {}) {
	        return new Cmyk(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.c = source["c"];
	        this.m = source["m"];
	        this.y = source["y"];
	        this.k = source["k"];
	    }
	}
	export class Rgb {
	    r: number;
	    g: number;
	    b: number;
	
	    static createFrom(source: any = {}) {
	        return new Rgb(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.r = source["r"];
	        this.g = source["g"];
	        this.b = source["b"];
	    }
	}

}

