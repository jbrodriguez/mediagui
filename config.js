var client = "src/client/";
var server = "src/server/";
var build = "target/build/"

var distTar = "unbalance";
var dist = "./" + distTar + "/";
var release = "./release";

var sources = {
		styles: "./src/styles/",
		images: "./src/images/",
		scripts: "./src/scripts/",
		svg: "./src/svg/",
		cache: "./src/cache/",
		tools: "./src/tools/"
};

var	staging = {
		root: "./staging/",
		styles: "./staging/css/",
		images: "./staging/img/",
		scripts: "./staging/js/",
};

module.exports = {
	base: {
		client: client
	},
	
	clean: {
		build: build,
		staging: staging.root,
		dist: dist,
		release: release
	},

	index: {
		src: client + "index.html",
		dst: build
	},

	app: {
		src: build,
		dst: build + "app/"
	},	

	build: {
		bin: "/Volumes/Users/kayak/code/bin/"
		src: server,
		dst: build
	},

	start: {
		src: build
	},

	watch: {
		app: build,
		index: client + "index.html",
		go: server + "**/*.go"
	},








	tools: {
		src: sources.tools + '*',
		dst: dist
	},

	templates: {
		src: client + "app/**/*.html",
		dst: sources.cache
	},

	styles: {
		vendors: client + "vendor/**/*.css",
		src: sources.styles + "styles.scss",
		dst: staging.styles
	},

	images: {
		cache: sources.cache,
		src: sources.images + "*",
		dst: staging.images
	},

	svg: {
		src: sources.svg + "*.svg",
		dst: staging.images
	},

	fingerprint: {
		revFilter: "**/*.{css,js,jpg,png,svg}",
		index: "index.html",

		src: [
			staging.root + "**/*.{css,js,jpg,png,svg}",
			client + "index.html"
		],
		dst: dist
	},

	reference: {
		ext: [
			".html",
			".js"
		],
		src: dist + "**/*.{html,js}",
		dst: dist
	},

	publish: {
		src: dist,
		dst: "/boot/custom/unbalance"
	},

	release: {
		src: distTar
	}
}