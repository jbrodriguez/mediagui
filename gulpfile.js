// 'use strict';

var gulp       	= require('gulp'),
    gutil      	= require('gulp-util'),
    // nodemon    	= require('gulp-nodemon'),
    source     	= require('vinyl-source-stream'),
    buffer     	= require('vinyl-buffer'),
    browserify 	= require('browserify'),
    babelify   	= require('babelify'),
    watchify   	= require('watchify'),
    lrload 		= require('livereactload'),
    path 		= require('path'),
    del 		= require('del'),
	strings 	= require('string'),
	exec 		= require('child_process').execSync,
	spawn 		= require('child_process').spawn,
    config		= require('./config.js')

function command(tag, cmd) {
	gutil.log(gutil.colors.blue('executing ' + cmd))
	var result = exec(cmd, {encoding: 'utf-8'})
	var output = strings(result).chompRight('\n').toString()
	gutil.log(gutil.colors.yellow('tag: [' + tag + '] ') + gutil.colors.green(output))
	return output
}

var isDebug = process.env.NODE_ENV === 'debug'
var mediagui;

gulp.task('dev', gulp.series(
		clean,
		gulp.parallel(client, server),
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

function index() {
	return gulp.src(config.index.src)
	 	.pipe(gulp.dest(config.index.dst))
}

var bundler = browserify({
	entries:      [ path.join(config.base.client, 'js/app.js') ],
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

// function start() {
// 	arg = path.join(process.cwd(), config.start.arg)
// 	cmd = path.join(process.cwd(), config.start.src, "mediagui") + " -webdir " + arg
// 	gutil.log('executing: ', cmd)
//     mediagui = exec(cmd, [' -webdir', arg])
//     // add a 'data' event listener for the spawn instance
//     mediagui.stdout.on('data', function(data) {
//     	gutil.log("sup dude:\n" + data);
//     })
//     // add an 'end' event listener to close the writeable stream
//     mediagui.stdout.on('end', function(data) {
//         gutil.log('mediagui stopped');
//     });

//     mediagui.on('error', function(data) {
// 		gutil.log(data);
//     })

//     // when the spawn child process exits, check if there were any errors and close the writeable stream
//     mediagui.on('close', function(code) {
//         if (code != 0) {
//             gutil.log('Failed: ' + code);
//         }
//     });
// }

function watch() {
    gutil.log('Watching ...')

	gulp.watch(config.watch.index, index)
	gulp.watch(config.watch.go, server)

	// start listening reload notifications
	lrload.monitor(path.join(config.watch.app, 'bundle.js'), {displayNotification: true})
}

// gulp.task('serverwatch', function() {
// 	nodemon({ script: 'server.js', ext: 'js', ignore: ['gulpfile.js', 'static/bundle.js', 'node_modules/*'] })
// 		.on('change', [])
// 		.on('restart', function () {
// 			console.log('Server restarted')
// 		})
// })

// gulp.task('watch', ['serverwatch', 'scripts'])

// gulp.task('watch', ['scripts'])





// gulp.task('copy', function() {
// 	gulp.src(config.copy.src)
// 	.pipe(gulp.dest(config.copy.dst))
// })

// gulp.task('clean', function() {
//     gutil.log('Cleaning: ' + gutil.colors.blue(config.clean.build));

//     del.sync(config.clean.build)
// })

// gulp.task('default', function(cb) {
// 	series('clean', 'copy', 'build:server', 'watch', cb)
// })