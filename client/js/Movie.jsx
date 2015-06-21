const 	React 			= require('react'),
		DatePicker 		= require('react-datepicker'),
		IconRating 		= require('react-icon-rating'),
		moment 			= require('moment'),
		moviesBO 		= require('./movies')

module.exports = React.createClass({
	setScore: function(score) {
		moviesBO.setMovieScore(this.props.movie, score)
	},

	getInitialState: function() {
		return {
			dateWatched: moment(),
			tmdb_id: this.props.movie.tmdb_id,
			rating: this.props.movie.rating
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

		const setStateWatched = function(date) {
			this.setState({
				dateWatched: date
			})
		}

		const setStateRating = function(data) {
			return
		}

		const fixTmdbId = function() {
			return
		}

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
						<div className="col-xs-12 col-sm-10">
							<img src={"/img/b" + movie.backdrop} />
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
						<div className="col-xs-12 col-sm-2 addon">
							<input className="addon-field" type="text" defaultValue={movie.tmdb_id}></input>
							<button className="btn btn-default">Fix</button>
						</div>
						<div className="col-xs-12 col-sm-10 addon end-sm">
							{score}
							<IconRating
								className="rspacer"
								max="10"
								currentRating={movie.score}
								toggledClassName="icon-star-filled"
								untoggledClassName="icon-star-empty"
								onChange={this.setScore} />
							<DatePicker
								key="{key}"
								placeholderText="YYYY-MM-DD"
								onChange={setStateWatched} />
						</div>
					</div>
				</div>										
			</article>				
		)	
	}
})