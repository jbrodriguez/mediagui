import React from 'react'

export default class CoverCard extends React.Component {
	render() {
		const movie = this.props.movie

		var watched;
		if (movie.count_watched > 0) {
			watched = (
				<div className="overlay__cover">
					<span>watched</span>
				</div>
			)
		}

		return (
			<article className="col-xs-12 col-sm-6 col-md-3 col-lg-2 bottom-spacer-large ">
				<div className="movie-cover">
					<div className="row" >
						<div className="col-xs-12">
							<div className="overlay">
								<img src={"/img/p" + movie.cover} />
								{watched}
							</div>
						</div>
					</div>
					<div className="row ">
						<div className="col-xs-12">
							<div className="movie-cover__details">
								<p className="cover-title">{movie.title}</p>
								<div className="between-xs cover-details">
									<span>{movie.year}</span>
									<span>{movie.imdb_rating}</span>
									<span>{this.hourMinute(movie.runtime)}</span>
								</div>
							</div>
						</div>
					</div>
				</div>
			</article>
		)
	}

	hourMinute(minutes) {
        var hour = Math.floor(minutes / 60)
        var minute = Math.floor(minutes % 60)

        var time = ''
        if (hour > 0) time += (hour + ":")
        if (minute >= 0) {
            if (minute <= 9) time += "0"+minute
            else time += minute
        }
        if (hour <= 0) time += "m"

        return time
	}
}