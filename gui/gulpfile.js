const gulp = require('gulp');
const ngAnnotate = require('gulp-ng-annotate');
const concat = require('gulp-concat');
const sass = require('gulp-sass');
const template = require('gulp-template');
const livereload = require('gulp-livereload');
const browserify = require('browserify');
const source = require('vinyl-source-stream');
const buffer = require('vinyl-buffer');

let isLiveReload = process.argv.indexOf('--live-reload') !== -1 || process.argv.indexOf('--livereload') !== -1;

gulp.task('templates:index', ['js', 'sass'], function() {

    return gulp.src("./web/src/game.html")
        .pipe(template({}))
        .pipe(gulp.dest("./web/www"));
});

gulp.task('templates:rest', function() {

    return gulp.src("./web/src/templates/**/*")
        .pipe(gulp.dest("./web/www/templates"));
});

gulp.task('js:libs', function() {

    return gulp.src([
        "./web/src/lib/q/q.js",
        "./web/src/lib/angular/angular.js",
        "./web/src/lib/angular-ui-router/release/angular-ui-router.js"
    ])
        .pipe(concat('libs.js'))
        .pipe(gulp.dest('./web/www/js/'));
});

gulp.task('js:app', function() {

    return gulp.src([
        './web/src/js/**/*.js'
    ])
        .pipe(concat('app.js'))
        .pipe(ngAnnotate())
        .pipe(gulp.dest('./web/www/js/'));
});

gulp.task('sass', function() {

    return gulp.src('./web/src/sass/app.scss')
        .pipe(sass({errLogToConsole: true}))
        .pipe(gulp.dest('./web/www/css/'));
});

gulp.task('copyfonts', function() {

    return gulp.src(['./web/src/lib/bootstrap-sass/assets/fonts/**/*'])
        .pipe(gulp.dest('./web/www/fonts'));
});

gulp.task('copyimages', function() {

    return gulp.src(['./web/src/img/*', './web/src/img/**/*'])
        .pipe(gulp.dest('./web/www/img'));
});

gulp.task('copystatics', ['copyfonts', 'copyimages']);

gulp.task('watch', function() {

    if (isLiveReload) {
        livereload.listen();
    }

    gulp.watch(['./web/src/sass/**/*.scss'], ['sass:livereload']);
    gulp.watch(['./web/src/img/**/*', './web/src/font/**/*'], ['copystatics:livereload']);
    gulp.watch(['./web/src/js/**/*.js'], ['js:app:livereload']);
    gulp.watch(['./web/src/templates/**/*', './web/src/game.html'], ['templates:livereload']);
});

gulp.task('js:app:livereload', ['js:app'], function() {
    livereload.reload();
});

gulp.task('templates:livereload', ['templates'], function() {
    livereload.reload();
});

gulp.task('sass:livereload', ['sass'], function() {
    livereload.reload();
});

gulp.task('default:livereload', ['default'], function() {
    livereload.reload();
});

gulp.task('copystatics:livereload', ['copystatics'], function() {
    livereload.reload();
});

gulp.task('js', ['js:libs', 'js:app']);
gulp.task('templates', ['templates:index', 'templates:rest']);
gulp.task('default', ['copystatics', 'sass', 'templates', 'js']);
