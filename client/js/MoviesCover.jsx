const 	React 	= require('react')



module.exports = React.createClass({
	// componentWillMount: function() {
	// 	movies.getCover()
	// },
    hourMinute: function(minutes) {
        var hour = Math.floor(minutes / 60);
        var minute = Math.floor(minutes % 60);

        var time = '';
        if (hour > 0) time += (hour + ":");
        if (minute >= 0) {
            if (minute <= 9) time += "0"+minute;
            else time += minute;
        }
        if (hour <= 0) time += "m";

        return time;
    },	

 //    truncate: function(text, length, end) {
	// 	length = length || 23;
	// 	end = end || '...';

	// 	if (text.length <= length || text.length - end.length <= length) {
	// 	    return text;
	// 	}
	// 	else {
	// 		console.log("is something here")
	// 	    return String(text).substring(0, length-end.length) + end;
	// 	}
	// },

	render: function() {
		// console.log('MoviesCover.jsx: ' + JSON.stringify(this.props, null, 4))

		const movies = this.props.movies.items

		// console.log('movies: ' + movies)

		const that = this


		if (typeof movies != 'undefined') {
			var items = movies.map(function(movie, i) {
				var watched;
				if (movie.count_watched > 0) {
					watched = (
						<div className="overlay__cover">
							<span>watched</span>
						</div>
					)
				}

				return (
					<div key={i} className="col-xs-12 col-sm-6 col-md-3 col-lg-2 bottom-spacer-large">
						<div className="overlay" >
							<img src={"/img/p" + movie.cover} />
							{watched}
							<p className="cover-title">{movie.title}</p>
							<span className="cover-details label">{movie.year} | {movie.imdb_rating} | {that.hourMinute(movie.runtime)}</span>
						</div>
					</div>				
				)
			})

			return (
				<section className="row">
					{items}
				</section>
			)

		} else {
			return null
		}

	}
})

