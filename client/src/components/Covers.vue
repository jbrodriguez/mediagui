<template>
	<section class="row">
		<article v-for="movie in movies" :key="movie.id" class="col-xs-12 col-sm-6 col-md-3 col-lg-2 mb3">
			<cover :movie="movie" />
		</article>
	</section>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import Component from 'vue-class-component'
import { State, Action, Getter, namespace } from 'vuex-class'

import { Movie } from '@/types'
import * as constant from '@/constants'

import Cover from './Cover.vue'

const domain = namespace('domain')

@Component({
	components: {
		Cover,
	},
	computed: {
		...mapGetters('domain', { movies: 'getMovies' }),
	},
})
export default class Covers extends Vue {
	public mounted() {
		this.$store.dispatch(constant.FETCH_MOVIES)
		// this.listMovies()
	}
}
</script>

<style>
</style>