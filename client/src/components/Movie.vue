<template>
	<article class="row top-xs c-minfo c-align" >
		<div class="col-xs-12 " >
			<div class="c-bg" :style="{ backgroundImage: 'url(' + background + ')' }">
			<div class="row">
				<div class="col-xs-12">
					<div class="row">
						<div class="col-xs-12">
							<div class="c-title-shade">
							<div class="row">
								<div class="col-xs-12 col-sm-10">
									<h2 class="c-title pt2 ml2">{{movie.title}} ({{movie.year}})</h2>
								</div>
								<div class="col-xs-12 col-sm-2 end-xs">
									<h2 class="c-title pt2 mr2">{{runtime}} | {{movie.imdb_rating}}</h2>
								</div>
							</div>
							</div>
						</div>
					</div>
				</div>

				<div class="col-xs-12">
					<div class="c-moverlay ml2">
						<img :src="cover" class="c-cover" />
						<div v-if="watched" class="c-moverlaycover">
							<span>watched</span>
						</div>
					</div>
				</div>

				<div class="col-xs-12 ph2 pv2">
					<div class="c-shaded ph2 pv2 mh2">
						<div class="row">
							<span class="col-xs-12 col-sm-6 c-text director">{{movie.director}}</span>
							<span class="col-xs-12 col-sm-6 end-sm c-text">{{movie.production_countries}}</span>
						</div>
						<div class="row between-xs mb2">
							<span class="col-xs-12 col-sm-6 c-text">{{movie.actors}}</span>
							<span class="col-xs-12 col-sm-6 end-sm c-text">{{movie.genres}}</span>
						</div>

						<div class="row between-xs">
							<div class="col-xs-12 col-sm-9">
								<span class="label">{{movie.resolution}}</span>
								<span class="label secondary spacer">{{movie.location}}</span>
								<span class="label">{{movie.id}}</span>
							</div>
							<div class="col-xs-12 col-sm-3 end-sm">
								<span v-if="watched" class="label success mv0 mh2">
									<font-awesome-icon class="mh1" icon="binoculars" />&nbsp;{{lastWatched}}
								</span>
								<span class="label">
									<font-awesome-icon class="mh1" icon="plus" />&nbsp;{{added}}
								</span>
							</div>
						</div>
						<div class="row mt2">
							<div class="col-xs-12 c-text">
								<span>{{overview}}</span>
							</div>
						</div>
						<div class="row between-xs mt2">
							<div class="col-xs-12 col-sm-4 addon">
								<input class="addon-field" type="number" min="0" step="1" v-model.number="tmdb" />
								<button class="btn btn-default mr2" @click="fixMovie">Fix</button>
								<div class="mr2">
									<label for="cbox" class="mr2">Not Dup ?</label>
									<input type="checkbox" id="cbox" v-model="duplicate" @change="setDuplicate">
								</div>
								<div v-if="loading"
									class="loading middle-xs">
									<div class="loading-bar"></div>
									<div class="loading-bar"></div>
									<div class="loading-bar"></div>
									<div class="loading-bar"></div>
								</div>
							</div>
							<div class="col-xs-12 col-sm-3 addon center-xs">
								<div v-if="watched">
									<span class="c-text">History:</span>
									<span class="label success mv0 mh2 ">{{movie.count_watched}}</span>
									<select :value="shows[shows.length-1]" class="c-select">
										<option v-for="(show, index) in shows" :key="index" :value="show">{{show}}</option>
									</select>
								</div>
							</div>
							<div class="col-xs-12 col-sm-5 addon end-sm">
								<span v-if="hasRating" class="label success mv0 mr2">{{movie.score}}</span>
								<Rating :max="10" :value="movie.score" @rating-selected="setScore" class="mr2" />
								<flat-pickr v-model="seen" :config="fpConfig" @on-change="setWatched"/>
							</div>
						</div>
					</div>
				</div>
			</div>
			</div>
		</div>
	</article>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import Component from 'vue-class-component'
import { Prop } from 'vue-property-decorator'
import { State, Action, Getter, namespace } from 'vuex-class'

import format from 'date-fns/format'
import FlatPickr from 'vue-flatpickr-component'
import 'flatpickr/dist/flatpickr.css'

import { Movie } from '@/types'
import { hourMinute } from '@/lib/utils'
import * as constant from '@/constants'

import Rating from './Rating.vue'

const en = {
	weekdays: {
		shorthand: ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'],
		longhand: ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'],
	},
	months: {
		shorthand: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
		longhand: [
			'January',
			'February',
			'March',
			'April',
			'May',
			'June',
			'July',
			'August',
			'September',
			'October',
			'November',
			'December',
		],
	},
	daysInMonth: [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31],
	firstDayOfWeek: 1,
	ordinal: (nth: number) => {
		const s = nth % 100
		if (s > 3 && s < 21) {
			return 'th'
		}
		switch (s % 10) {
			case 1:
				return 'st'
			case 2:
				return 'nd'
			case 3:
				return 'rd'
			default:
				return 'th'
		}
	},
	rangeSeparator: ' to ',
	weekAbbreviation: 'Wk',
	scrollTitle: 'Scroll to increment',
	toggleTitle: 'Click to toggle',
	amPM: ['AM', 'PM'],
	yearAriaLabel: 'Year',
}

@Component({
	components: {
		Rating,
		FlatPickr,
	},
	computed: {
		...mapGetters('domain', { movies: 'getMovies' }),
	},
})
export default class MovieX extends Vue {
	@Prop(Object)
	private movie!: Movie

	private tmdb: number = this.movie.tmdb_id
	private duplicate: boolean = this.movie.showIfDuplicate === 0
	private seen: Date = new Date()
	private fpConfig: any = {
		locale: en,
	}
	private loading: boolean = false

	private updated() {
		this.loading = false
	}

	private setScore(score: number) {
		// console.log(`score(${score})`)
		this.$store.dispatch(constant.SET_SCORE, { id: this.movie.id, score })
	}

	private setWatched(seen: Date) {
		const watched = format(seen, 'YYYY-MM-DDTHH:mm:ssZ')
		// console.log(`seen(${this.seen})-args(${seen})-watched(${watched})`)
		// console.log(`seen(${watched})`)
		this.$store.dispatch(constant.SET_WATCHED, { id: this.movie.id, watched })
	}

	private fixMovie() {
		// console.log(`typeof(${typeof this.tmdb})`)
		// if (this.tmdb !== this.movie.tmdb_id) {
		this.loading = true
		this.$store.dispatch(constant.FIX_MOVIE, { id: this.movie.id, tmdb: this.tmdb })
		// }
	}

	private setDuplicate() {
		this.loading = true
		// console.log(`this.duplicate(${this.duplicate})`)
		this.$store.dispatch(constant.SET_DUPLICATE, {
			id: this.movie.id,
			showIfDuplicate: this.duplicate ? 0 : 1,
		})
	}

	private get watched() {
		return this.movie.count_watched > 0
	}

	private get runtime() {
		return hourMinute(this.movie.runtime)
	}

	private get added() {
		// const added = parse(this.movie.added);
		return format(this.movie.added, 'MMM DD, YYYY H:mm')
	}

	private get hasRating() {
		return this.movie.score !== 0
	}

	private get lastWatched() {
		return format(this.movie.last_watched, 'MMM DD, YYYY')
	}

	private get cover() {
		return `img/p${this.movie.cover}`
	}

	private get background() {
		return `/img/b${this.movie.backdrop}`
	}

	private get overview() {
		return this.movie.overview.length > 675 ? `${this.movie.overview.slice(0, 675)} ...` : this.movie.overview
	}

	private get shows() {
		return this.movie.all_watched.split('|').map(show => format(show, 'MMM DD, YYYY'))
	}
}
</script>

<style lang="scss" scoped>
@import '../styles/variables.scss';

.c-bg {
	// height: 765px;
	background-size: cover;
}

.c-minfo {
	font-size: 0.9em;
	margin-bottom: 2em;
	border-bottom: 1px solid $article-border-color;

	&:not(last-child) {
		margin-bottom: 2em;
	}

	.director {
		color: $director-color;
		font-weight: bold;
	}
}

.c-align {
	align-content: space-between;
}

.c-title {
	color: #fff;
	text-shadow: 0 1px 0 hsla(0, 0%, 0%, 0.75), 0 0 1px hsla(0, 0%, 0%, 0.75), 0 1px 5px hsla(0, 0%, 0%, 0.75);
}

.c-title-shade {
	background: linear-gradient(to bottom, rgba(0, 0, 0, 0.65) 0%, rgba(0, 0, 0, 0) 100%);
}

.c-text {
	color: #fff;
	text-shadow: 0 1px 0 hsla(0, 0%, 0%, 0.75), 0 0 1px hsla(0, 0%, 0%, 0.75), 0 1px 5px hsla(0, 0%, 0%, 0.75);
}

.c-shaded {
	// font-size: 1.25em;
	color: white;
	background-color: rgba(0, 0, 0, 0.5); // padding: 0.25em 1em;
}

.c-cover {
	height: 400px;
	width: auto;
	opacity: 0.75;
}

.c-moverlay {
	position: relative;
	overflow: hidden;
}

.c-moverlaycover {
	position: absolute;
	left: -40px;
	top: 10px;
	background-color: rgba(170, 0, 0, 0.6);
	transform: rotate(-45deg);
	box-shadow: 0 0 10px #888;

	span {
		border: 1px solid #faa;
		color: #fff;
		display: block;
		font-size: 0.8em;
		margin: 1px 0;
		padding: 5px 50px;
		text-align: center;
		text-decoration: none;
		/* shadow */
		text-shadow: 0 0 5px #444;
	}
}

.c-select {
	width: 125px;
}
</style>
