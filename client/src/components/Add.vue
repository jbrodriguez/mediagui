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

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import Component from 'vue-class-component'

import * as constant from '@/constants'

import Movie from './Movie.vue'

@Component({
	components: {
		Movie,
	},
	computed: {
		...mapGetters('domain', { movies: 'getMovies' }),
	},
})
export default class Add extends Vue {
	private title: string = ''

	private mounted() {
		this.$store.commit(constant.CLEAN_MOVIES)
	}

	private addMovie() {
		this.$store.commit(constant.CLEAN_MOVIES)
		this.$store.dispatch(constant.ADD_MOVIE, { title: this.title })
	}
}
</script>

<style lang="scss" scoped>
</style>
