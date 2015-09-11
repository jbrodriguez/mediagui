import React from 'react'
import DatePicker from 'react-datepicker'
import IconRating from 'react-icon-rating'
import moment from 'moment'
import { hourMinute } from '../lib/utils'

export default class MovieCard extends React.Component {
	constructor() {
		super()

		this.setScore = this.setScore.bind(this)
		this.setWatched = this.setWatched.bind(this)
		this.setTmdbId = this.setTmdbId.bind(this)
		this.fixMovie = this.fixMovie.bind(this)

		this.state = {
			loading: false
		}
	}

	render() {
		const movie = this.props.movie
		const key = this.props.key
		
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

		var watched_ribbon;
		if (movie.count_watched > 0) {
			watched_ribbon = (
				<div className="overlay__cover">
					<span>watched</span>
				</div>
			)
		}

		return (
			<article className="movie-info">
				<div className="col-xs-12">
					<h2>{movie.title} ({movie.year})</h2>
				</div>
				<div className="col-xs-12">
					<div className="row moviep-images">
						<div className="col-xs-12 col-sm-2">
							<div className="overlay">
								<img src={"/img/p" + movie.cover} />
								{watched_ribbon}
							</div>
						</div>
						<div className="col-xs-12 col-sm-10 overlay">
								<img src={"/img/b" + movie.backdrop} />
								<div className="row between-xs overlay__backdrop">
									<div className="col-xs-6">
										<span>{hourMinute(movie.runtime)}</span>
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
								selected={moment()}
								placeholderText="YYYY-MM-DD"
								onChange={this.setWatched} />
						</div>
					</div>
				</div>										
			</article>				
		)	

	}

	setScore(score) {
		// console.log('score: ', score)
		// console.log('this.props.movie: ', this.props.movie)
		this.props.actions.movies.setMovieScore(this.props.movie, score)
	}

	setWatched(watched) {
		this.props.actions.movies.setMovieWatched(this.props.movie, watched)
	}

	setTmdbId(e) {
		this.tmdb_id = e.target.value
	}

	fixMovie() {
		if (this.tmdb_id) {
			this.setState({ loading: true })
			this.props.actions.movies.fixMovie(this.props.movie, parseInt(this.tmdb_id))
		}
	}	
}