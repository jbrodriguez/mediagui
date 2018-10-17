<template>
	<section v-if="total > 0">
		<div class="row between-xs middle-xs mb2">
			<div class="col-xs-12 col-sm-10 middle-xs">
				<paginate v-model="page"
						:pageCount="pageCount"
				        :pageRange="5"
				        :marginPages="3"
				        :forcePage="forcePage"
				    	:containerClass="'pagination'"
				        :clickHandler="onPaginationClick"
						:page-class="'p-li'"
						:page-link-class="'p-a'"
						:next-link-class="'p-next'"
						/>
			</div>
			<div class="col-xs-12 col-sm-2 end-xs">
				<span>TOTAL </span>
				<span class="c-total">{{ total }}</span>
			</div>
		</div>

		<movie-x v-for="movie in movies" :movie="movie" :key="movie.id" />
		<!-- <span v-for="movie in movies" :movie="movie" :key="movie.id">{{movie.title}}</span> -->

		<div class="row">
			<div class="col-xs-12 middle-xs">
				<paginate v-model="page"
						:pageCount="pageCount"
				        :pageRange="5"
				        :marginPages="3"
				        :forcePage="forcePage"
				    	:containerClass="'pagination'"
				        :clickHandler="onPaginationClick"
						:page-class="'p-li'"
						:page-link-class="'p-a'"
						/>
			</div>
		</div>
	</section>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import Component from 'vue-class-component'
import { State, Action, Getter } from 'vuex-class'

import Paginate from 'vuejs-paginate'

import { Movie, OptionsState, DomainState, RootState } from '@/types'
import * as constant from '@/constants'
import MovieX from './Movie.vue'

// const domain = namespace('domain')

@Component({
	components: {
		MovieX,
		Paginate,
	},
	computed: {
		...mapGetters('domain', { movies: 'getMovies' }),
	},
})
export default class Movies extends Vue {
	private shouldScroll: boolean = false
	private page: number = 1

	@State((state: RootState) => state.options)
	private options!: OptionsState

	@State((state: RootState) => state.domain)
	private domain!: DomainState

	private mounted() {
		this.$store.dispatch(constant.FETCH_MOVIES)
	}

	private updated() {
		if (this.shouldScroll) {
			window.scrollTo(0, 0)
			this.shouldScroll = false
		}
	}

	private onPaginationClick(pageNum: number) {
		this.shouldScroll = true

		// let's keep both paginators in sync, they're zero-based
		// this.$refs.paginator1.selected = pageNum - 1;
		// this.$refs.paginator2.selected = pageNum - 1;

		const offset = Math.ceil((pageNum - 1) * this.options.limit)
		this.$store.commit(constant.SET_OFFSET, offset)
		this.$store.dispatch(constant.FETCH_MOVIES)
	}

	get pageCount() {
		return Math.ceil(this.domain.total / this.options.limit)
	}

	get forcePage() {
		return this.options.offset / this.options.limit
	}

	get total() {
		return this.domain.total
	}

	get visible() {
		return this.domain.total > 0
	}
}

// import Paginate from 'vuejs-paginate'

// import Movie from './Movie'
// import Pager from './Pager'
// import * as types from '../store/types'

// export default {
// 	name: 'movies',

// 	components: { Paginate, Movie },

// 	data() {
// 		return {
// 			shouldScroll: false,
// 		}
// 	},

// 	methods: {
// 		onPaginationClick(pageNum) {
// 			this.shouldScroll = true

// 			// let's keep both paginators in sync, they're zero-based
// 			// this.$refs.paginator1.selected = pageNum - 1;
// 			// this.$refs.paginator2.selected = pageNum - 1;

// 			const offset = Math.ceil((pageNum - 1) * this.$store.state.options.limit)
// 			this.$store.commit(types.SET_OFFSET, offset)
// 			this.$store.dispatch(types.FETCH_MOVIES)
// 		},
// 	},

// 	computed: {
// 		movies() {
// 			return this.$store.getters.getMovies
// 		},

// 		pageCount() {
// 			return Math.ceil(this.$store.state.total / this.$store.state.options.limit)
// 		},

// 		forcePage() {
// 			return this.$store.state.options.offset / this.$store.state.options.limit
// 		},

// 		total() {
// 			return this.$store.state.total
// 		},

// 		visible() {
// 			return this.$store.state.total > 0
// 		},
// 	},

// 	created() {
// 		this.$store.dispatch(types.FETCH_MOVIES)
// 	},

// 	updated() {
// 		if (this.shouldScroll) {
// 			window.scrollTo(0, 0)
// 			this.shouldScroll = false
// 		}
// 	},
// }
</script>

<style lang="scss" scoped>
@import '../styles/variables.scss';

.c-total {
	color: $primary-header;
}
</style>
