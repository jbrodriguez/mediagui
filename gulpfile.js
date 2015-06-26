// 'use strict';

var gulp       		= require('gulp'),
    gutil      		= require('gulp-util'),
    // nodemon    	= require('gulp-nodemon'),
    source     		= require('vinyl-source-stream'),
    buffer     		= require('vinyl-buffer'),
    browserify 		= require('browserify'),
    babelify   		= require('babelify'),
    watchify   		= require('watchify'),
    lrload 			= require('livereactload'),
    path 			= require('path'),
    del 			= require('del'),
    sleep			= require('sleep'),
	strings 		= require('string'),
	exec 			= require('child_process').execSync,
	spawn 			= require('child_process').spawn,

	sass 			= require('gulp-sass');
	autoprefixer 	= require('gulp-autoprefixer');
	concat 			= require('gulp-concat');
	minifyCss 		= require('gulp-minify-css');
	bytediff 		= require('gulp-bytediff');
	plumber 		= require('gulp-plumber');

	imagemin = require('gulp-imagemin');
	cache = require('gulp-cache');

    config			= require('./config.js')


var isDebug = process.env.NODE_ENV === 'debug'
var mediagui;


gulp.task('client', gulp.series(client))
gulp.task('styles', gulp.series(styles))
gulp.task('publish', gulp.series(
		clean,
		gulp.parallel(client_release, server, styles, images, fonts),
		publish
	)
)

gulp.task('dev', gulp.series(
		clean,
		gulp.parallel(client, server, styles, images, fonts),
		link,
		watch
	)
)

gulp.task('default', gulp.series('dev'))

function clean(done) {
    gutil.log('Cleaning: ' + gutil.colors.blue(config.clean.build))

    del.sync(config.clean.build)

    done()
}

function client(done) {
	index()
	app()

	done()
}

function client_release(done) {
	index()
	app_release()

	done()
}

function index() {
	return gulp.src(config.index.src)
	 	.pipe(gulp.dest(config.index.dst))
}

var bundler = browserify({
	entries:      [ path.join(config.base.client, 'js/app.js') ],
	extension: 	  [ "jsx" ],
	transform:    isDebug ? [ babelify, lrload ] : [ babelify ],
	debug:        isDebug,
	cache:        {},
	packageCache: {},
	fullPaths:    true // for watchify
})

function app() {
	// start JS file watching and rebundling with watchify
	var watcher = watchify(bundler)

	rebundle()

	watcher
		.on('error', gutil.log)
		.on('update', rebundle)

	function rebundle() {
		gutil.log('Update JavaScript bundle')
		watcher
			.bundle()
			.on('error', gutil.log)
			.pipe(source('bundle.js'))
			.pipe(buffer())
			.pipe(gulp.dest(config.app.dst))
	}
}

function app_release() {
	gutil.log("antes de app_release")


	bundler
		.bundle()
		.pipe(source('bundle.js'))
		// .pipe(buffer())
		.pipe(gulp.dest(config.app.dst))

	gutil.log("luego de app_release")
}

function server(done) {
	command('ls', 'ls -al /Volumes/Users/kayak/code/src/jbrodriguez/mediagui/target')

	// stop()
	build()

	done()
}

function stop() {
	command('kill9', 'pkill mediagui')
}

function build() {
	var version = command('version', 'cat VERSION')
	var count = command('count', 'git rev-list HEAD --count')
	var hash = command('hash', 'git rev-parse --short HEAD')

	gutil.log('\n src: ' + config.build.src + '\n dst: ' + config.build.dst)
	command('build', 'cd server && ' + config.build.bin + 'gom build -ldflags \"-X main.Version ' + version + '-' + count + '.' + hash + '\" -v -o ' + path.join(config.build.dst, 'mediagui') + ' main.go && cd ..')
}

function styles() {
    gutil.log('Bundling, minifying, and copying the app\'s css');

    return gulp.src(config.styles.src)
        .pipe(plumber())
		.pipe(sass())
        .pipe(concat('app.min.css')) // Before bytediff or after
        .pipe(autoprefixer('last 2 version', '> 5%'))
        .pipe(bytediff.start())
        .pipe(minifyCss({processImport: false}))
        .pipe(bytediff.stop(bytediffFormatter))
        //        .pipe(plug.concat('all.min.css')) // Before bytediff or after
        .pipe(plumber.stop())
        .pipe(gulp.dest(config.styles.dst));
}

function images() {
    gutil.log('Compressing, caching, and copying images ');

    gutil.log('cache: ' + gutil.colors.green(config.images.cache));
    gutil.log('src: ' + gutil.colors.green(config.images.src));
    gutil.log('dst: ' + gutil.colors.green(config.images.dst));

    var custom = new cache.Cache({ tmpDir: config.images.cache, cacheDirName: '' })

    return gulp
		.src(config.images.src)
        .pipe(cache(imagemin({optimizationLevel: 3}), {fileCache: custom, name: ''}))
        .pipe(gulp.dest(config.images.dst))	
}

function fonts() {
	return gulp
		.src(config.fonts.src)
		.pipe(gulp.dest(config.fonts.dst))
}

function link(done) {
	var home = process.env[(process.platform == 'win32') ? 'USERPROFILE' : 'HOME']
	var img = path.join(home, '.mediagui', 'web', 'img')

	gutil.log('\n src: ' + config.build.src + '\n dst: ' + config.build.dst)
	command('link', 'cd target/build && ln -s ' + img + ' img')

	done()
}

function watch() {
    gutil.log('Watching ...')

	gulp.watch(config.watch.index, index)
	gulp.watch(config.watch.go, server)
	gulp.watch(config.watch.styles, styles)
	gulp.watch(config.watch.images, images)
	gulp.watch(config.watch.fonts, fonts)

	// start listening reload notifications
	lrload.monitor(path.join(config.watch.app, 'bundle.js'), {displayNotification: true})
}


function publish(done) {
	var home = process.env[(process.platform == 'win32') ? 'USERPROFILE' : 'HOME']

    // const app = path.join(config.publish.src, config.publish.app, "**/*")
    const app = path.join(config.publish.src, config.publish.app)
    const index = path.join(config.publish.src, config.publish.index)
    const bin = path.join(config.publish.src, "mediagui")

    const dst = path.join(home, ".mediagui", "web")
	const binDst = path.join(home, "bin")

    const delAppDst = path.join(dst, config.publish.app)
    const delIndexDst = path.join(dst, config.publish.index)
    const delBinDst = path.join(binDst, "mediagui")

	gutil.log("app: ", app)
	gutil.log("index: ", index)
	gutil.log("bin: ", bin)

	gutil.log("dst: ", dst)
	gutil.log("binDst: ", binDst)

	gutil.log("delAppDst: ", delAppDst)
	gutil.log("delIndexDst: ", delIndexDst)
	gutil.log("delBinDst: ", delBinDst)

    del.sync(delAppDst, {force: true})
    del.sync(delIndexDst, {force: true})
    del.sync(delBinDst, {force: true})

	gulp.src(bin).pipe(gulp.dest(binDst))
	gulp.src(index).pipe(gulp.dest(dst))
	gulp.src(app).pipe(gulp.dest(dst))

	// gulp
	// 	.src(
	// 		path.join(config.publish.src, "app", "bundle.js")
	// 	)
	// 	.pipe(
	// 		gulp.dest(
	// 			path.join(home, ".mediagui", "web", "app")
	// 		)
	// 	)

	done()
}

// HELPERS
function bytediffFormatter(data) {
    var difference = (data.savings > 0) ? ' smaller.' : ' larger.'
    return data.fileName + ' went from ' +
        (data.startSize / 1000).toFixed(2) + ' kB to ' + (data.endSize / 1000).toFixed(2) + ' kB' +
        ' and is ' + formatPercent(1 - data.percent, 2) + '%' + difference
}

function formatPercent(num, precision) {
    return (num * 100).toFixed(precision)
}

function command(tag, cmd) {
	gutil.log(gutil.colors.blue('executing ' + cmd))
	var result = exec(cmd, {encoding: 'utf-8'})
	var output = strings(result).chompRight('\n').toString()
	gutil.log(gutil.colors.yellow('tag: [' + tag + '] ') + gutil.colors.green(output))
	return output
}
