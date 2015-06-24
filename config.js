var client = "client/";
var server = "server/";
var build = "target/build/"
var assets = "assets/"

// var distTar = "unbalance";
// var dist = "./" + distTar + "/";
// var release = "./release";

// var sources = {
// 		styles: "./src/styles/",
// 		images: "./src/images/",
// 		scripts: "./src/scripts/",
// 		svg: "./src/svg/",
// 		cache: "./src/cache/",
// 		tools: "./src/tools/"
// };

// var	staging = {
// 		root: "./staging/",
// 		styles: "./staging/css/",
// 		images: "./staging/img/",
// 		scripts: "./staging/js/",
// };

module.exports = {
	base: {
		client: client
	},
	
	clean: {
		build: build
		// staging: staging.root,
		// dist: dist,
		// release: release
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
		bin: "/Volumes/Users/kayak/code/bin/",
		src: server,
		dst: "../" + build
	},

	start: {
		src: build,
		arg: build
	},

	styles: {
		vendors: client + "vendor/**/*.css",
		src: assets + "styles/styles.scss",
		dst: build + "app/"
	},	

	images: {
		cache: "staging/img/",
		src: assets + "images/*",
		dst: build + "app/"
	},

	fonts: {
		src: assets + "fonts/*",
		dst: build + "app/fonts/"
	},


	watch: {
		app: build,
		index: client + "index.html",
		go: server + "**/*.go",
		styles: assets + "styles/*.scss",
		images: assets + "images/*",
		fonts: assets + "fonts/*"
	},



	// tools: {
	// 	src: sources.tools + '*',
	// 	dst: dist
	// },

	// templates: {
	// 	src: client + "app/**/*.html",
	// 	dst: sources.cache
	// },

	// images: {
	// 	cache: sources.cache,
	// 	src: sources.images + "*",
	// 	dst: staging.images
	// },

	// svg: {
	// 	src: sources.svg + "*.svg",
	// 	dst: staging.images
	// },

	// fingerprint: {
	// 	revFilter: "**/*.{css,js,jpg,png,svg}",
	// 	index: "index.html",

	// 	src: [
	// 		staging.root + "**/*.{css,js,jpg,png,svg}",
	// 		client + "index.html"
	// 	],
	// 	dst: dist
	// },

	// reference: {
	// 	ext: [
	// 		".html",
	// 		".js"
	// 	],
	// 	src: dist + "**/*.{html,js}",
	// 	dst: dist
	// },

	// publish: {
	// 	src: dist,
	// 	dst: "/boot/custom/unbalance"
	// },

	// release: {
	// 	src: distTar
	// }
}