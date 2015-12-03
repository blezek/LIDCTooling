var gulp = require('gulp');
var $    = require('gulp-load-plugins')();
var watch = require('gulp-watch');

var sassPaths = [
  'bower_components/foundation-sites/scss',
  'bower_components/motion-ui/src',
  'scss'
];

var jsPaths = [
  'bower_components/jquery/dist/jquery.min*',
  'bower_components/motion-ui/dist/motion-ui.min*',
  'bower_components/what-input/what-input.min.js',
  'bower_components/foundation-sites/js/foundation.min.js'
]

gulp.task('js', function() {
  gulp.src(jsPaths).pipe(gulp.dest("static/js"));
});

gulp.task('sass', function() {
  return gulp.src('scss/app.scss')
    .pipe($.sass({
      includePaths: sassPaths
    })
      .on('error', $.sass.logError))
    .pipe($.autoprefixer({
      browsers: ['last 2 versions', 'ie >= 9']
    }))
    .pipe(gulp.dest('static/css'));
});

gulp.task('default', ['sass', 'js'], function() {
});

gulp.task('watch', function() {
  gulp.watch(['scss/*'], ['sass']);
})
