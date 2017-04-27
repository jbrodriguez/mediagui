<template>
	<div class="">
		<section class="row">
			<div class="col-xs-12 mb2 middle-xs">
				<input type="text" placeholder="Enter movie title" v-model="title">

				<button class="btn btn-default" @click="addMovie">Add Stub Movie</button>
			</div>
		</section>

		<movie v-for="movie in movies" :movie="movie" :key="movie.id"></movie>

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
			this.$store.commit(types.CLEAN_MOVIES)
			this.$store.dispatch(types.ADD_MOVIE, { title: this.title })
		},
	},
}
</script>

<style lang="scss" scoped>

</style>
