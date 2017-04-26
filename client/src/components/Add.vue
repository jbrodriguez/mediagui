<template>
	<div class="">
		<section class="row">
			<div class="col-xs-12 mb2 middle-xs">
				<input type="text"
						placeholder="Enter movie title"
						v-model="title">

				<input type="text"
						placeholder="Enter movie year"
						v-model="year">

				<input type="text"
						placeholder="Enter movie tmdb_id"
						v-model="tmdb_id">

				<button class="btn btn-default"
				        @click="addMovie">Add Stub Movie</button>
			</div>
		</section>

		<movie v-for="movie in movies"
		       :movie="movie"
		       :key="movie.id"></movie>

	</div>
</template>

<script>
import * as types from '../store/types'
import Movie from './Movie'

export default {
	name: 'add',

	components: { Movie },

	data() {
		return {
			title: '',
			year: '',
			tmdb_id: '',
		}
	},

	created() {
		this.$store.commit(types.CLEAN_MOVIES)
	},

	computed: {
		movies() {
			return this.$store.getters.getMovies
		},
	},

	methods: {
		addMovie() {
			this.$store.dispatch(types.ADD_MOVIE, { title: this.title, year: this.year, tmdb_id: this.tmdb_id })
		},
	},
}
</script>

<style lang="scss" scoped>

</style>
