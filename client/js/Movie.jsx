const 	React 			= require('react'),
		DatePicker 		= require('react-datepicker'),
		moment 			= require('moment')

module.exports = React.createClass({
	getInitialState: function() {
		return {
			date: moment(),
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

		const setStateWatched = function(data) {
			return
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

		return (
			<article key={key}>
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
				<div className="col-xs-12">
					<span>{movie.overview}</span>
				</div>
				<div className="col-xs-12 bottom-spacer-large">
					<div className="row between-xs">
						<div className="col-xs-12 col-sm-10 top-xs date-picker">
							<input type="text"></input>
							<button className="btn btn-default">Fix</button>
							<DatePicker
								key="example1"
								selected={this.state.date}
								onChange={handleWatched}
							/>
						</div>
						<div className="col-xs-12 col-sm-2 end-sm top-xs">
							<span className="label"><i className="icon-watched"></i></span>
						</div>
					</div>
				</div>										
			</article>				
		)	
	}
})