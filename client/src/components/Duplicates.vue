<template>
	<section>
		<movie v-for="movie in movies" :movie="movie" :key="movie.id"></movie>
	</section>
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
export default class Duplicates extends Vue {
	public mounted() {
		this.$store.dispatch(constant.FETCH_DUPLICATES)
	}
}
</script>

<style>
</style>
