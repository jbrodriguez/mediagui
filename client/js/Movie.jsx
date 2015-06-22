const 	React 			= require('react'),
		DatePicker 		= require('react-datepicker'),
		IconRating 		= require('react-icon-rating'),
		moment 			= require('moment'),
		moviesBO 		= require('./movies')

module.exports = React.createClass({
	setScore: function(score) {
		moviesBO.setMovieScore(this.props.movie, score)
	},

	setWatched: function(watched) {
		moviesBO.setMovieWatched(this.props.movie, watched)
	},

	setTmdbId: function(e) {
		this.tmdb_id = e.target.value
	},

	fixMovie: function() {
		if (this.tmdb_id) {
			this.setState({ loading: true })
			moviesBO.fixMovie(this.props.movie, parseInt(this.tmdb_id))
		}
	},

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

	getInitialState: function() {
		return {
			loading: false
		}
	},

	render: function() {
		const movie = this.props.movie
		const key = this.props.key
		// const options = this.props.options

		// var that = this;

		const setStateTmdbId = function(data) {
			return
		}

		const fixTmdbId = function() {
			return
		}

				// <span className="label success spacer"><i className="icon-watched"></i>&nbsp;{moment.utc(movie.last_watched).local().format('MMM DD, YYYY')}</span>


		const saveWatched = function() {
		}
		
		var watched;
		if (movie.last_watched != '') {
			watched = (
				<span className="label success spacer"><i className="icon-watched"></i>&nbsp;{moment(movie.last_watched).format('MMM DD, YYYY')}</span>
			)
		}

		var score;
		if (movie.score != 0) {
			score = (
				<span className="label success rspacer">{movie.score}</span>
			)
		}

		var loading;
		if (this.state.loading) {
			loading = (
				<div className="loading middle-xs">
					<div className="loading-bar"></div>
					<div className="loading-bar"></div>
					<div className="loading-bar"></div>
					<div className="loading-bar"></div>
				</div>
			)
		}


		// console.log('movie.score: id('+movie.id+')-score('+movie.score+')')

		return (
			<article>
				<div className="col-xs-12">
					<h2>{movie.title} ({movie.year})</h2>
				</div>
				<div className="col-xs-12">
					<div className="row moviep-images">
						<div className="col-xs-12 col-sm-2">
							<img src={"/img/p" + movie.cover} />
						</div>
						<div className="col-xs-12 col-sm-10 backOver">
							<img src={"/img/b" + movie.backdrop} />
							<div className="row between-xs backOver-wrap">
								<div className="col-xs-6">
									<span>{this.hourMinute(movie.runtime)}</span>
								</div>
								<div className="col-xs-6 end-xs">
									<span>{movie.imdb_rating}</span>
								</div>
							</div>
						</div>
					</div>
				</div>
				<div className="col-xs-12">
					<div className="row between-xs">
						<span className="col-xs-12 col-sm-6 director">{movie.director}</span>
						<span className="col-xs-12 col-sm-6 end-sm">{movie.production_countries}</span>
					</div>
				</div>
				<div className="col-xs-12">
					<div className="row between-xs">
						<span className="col-xs-12 col-sm-6 ">{movie.actors}</span>
						<span className="col-xs-12 col-sm-6 end-sm">{movie.genres}</span>
					</div>
				</div>
				<div className="col-xs-12 bottom-spacer">
					<div className="row between-xs">
						<div className="col-xs-12 col-sm-9">
							<span className="label">{movie.resolution}</span>
							<span className="label secondary spacer">{movie.location}</span>
						</div>
						<div className="col-xs-12 col-sm-3 end-sm">
							{watched}							
							<span className="label"><i className="icon-plus"></i>&nbsp;{moment(movie.added).format('MMM DD, YYYY H:mm')}</span>
						</div>
					</div>
				</div>
				<div className="col-xs-12 bottom-spacer">
					<span>{movie.overview}</span>
				</div>
				<div className="col-xs-12 bottom-spacer-large">
					<div className="row between-xs">
						<div className="col-xs-12 col-sm-3 addon">
							<input className="addon-field" type="text" defaultValue={movie.tmdb_id} onChange={this.setTmdbId}></input>
							<button className="btn btn-default rspacer" onClick={this.fixMovie}>Fix</button>
							{loading}
						</div>
						<div className="col-xs-12 col-sm-9 addon end-sm">
							{score}
							<IconRating
								className="rspacer"
								max="10"
								currentRating={movie.score}
								toggledClassName="icon-star-full"
								untoggledClassName="icon-star-empty"
								onChange={this.setScore} />
							<DatePicker
								key="{key}"
								placeholderText="YYYY-MM-DD"
								onChange={this.setWatched} />
						</div>
					</div>
				</div>										
			</article>				
		)	
	}
})