<template>
  <article class="row top-xs c-minfo c-align" :style="{ backgroundImage: 'url(' + background + ')' }">
    <div class="col-xs-12">
        <div class="row">
          <h2 class="col-xs-12 col-sm-10 c-title pt2">{{movie.title}} ({{movie.year}})</h2>
          <h2 class="col-xs-12 col-sm-2 c-title end-xs pt2">{{runtime}} | {{movie.imdb_rating}}</h2>
        </div>
        <div class="row between-xs">
          <span class="col-xs-12 col-sm-6 c-text">{{movie.director}}</span>
          <span class="col-xs-12 col-sm-6 end-sm c-text">{{movie.production_countries}}</span>
        </div>
        <div class="row between-xs">
          <span class="col-xs-12 col-sm-6 c-text">{{movie.actors}}</span>
          <span class="col-xs-12 col-sm-6 end-sm c-text">{{movie.genres}}</span>
        </div>
    </div>

    <div class="col-xs-12">
        <div class="c-moverlay">
          <img :src="cover" class="c-cover" />
          <div v-if="watched" class="c-moverlaycover">
            <span>watched</span>
          </div>
        </div>
    </div>

    <div class="col-xs-12 mb2">
      <div class="row between-xs">
        <div class="col-xs-12 col-sm-9">
          <span class="label">{{movie.resolution}}</span>
          <span class="label secondary spacer">{{movie.location}}</span>
          <span class="label">{{movie.id}}</span>
        </div>
        <div class="col-xs-12 col-sm-3 end-sm">
          <span v-if="watched" class="label success mv0 mh2"><i class="fa fa-binoculars mh1"></i>&nbsp;{{lastWatched}}</span>
          <span class="label"><i class="fa fa-plus mh1"></i>&nbsp;{{added}}</span>
        </div>
      </div>
      <div class="row between-xs ph2 mt2">
        <div class="col-xs-12 c-text c-shaded">
          <span>{{movie.overview}}</span>
        </div>
        </span>
      </div>
					<div class="row between-xs middle-xs mt2">
						<div class="col-xs-12 col-sm-2 addon">
							<input class="addon-field" type="number" min="0" step="1" v-model.number="tmdb" />
							<button class="btn btn-default mr2" @click="fixMovie">Fix</button>
							<div v-if="loading" class="loading middle-xs">
                <div class="loading-bar"></div>
                <div class="loading-bar"></div>
                <div class="loading-bar"></div>
                <div class="loading-bar"></div>
							</div>
						</div>
						<div class="col-xs-12 col-sm-4 addon end-xs">
							<div v-if="watched">
                <span class="c-text">History:</span>
                <span class="label success mv0 mh2 ">{{movie.count_watched}}</span>
                <select :value="shows[shows.length-1]" class="c-select">
                  <option v-for="show in shows" :value="show">{{show}}</option>
                </select>
							</div>
						</div>
						<div class="col-xs-12 col-sm-6 addon end-sm">
							<span v-if="hasRating" class="label success mv0 mr2">{{movie.score}}</span>
							<Rating :max="10" :value="movie.score" @rating-selected="setScore" class="mr2" />
              <!--<datepicker v-model="seen" monday-first @selected="setWatched"></datepicker>>-->
              <VueFlatpickr v-model="seen" :options="fpOptions" />
						</div>
					</div>

    </div>
  </article>
</template>

<script>
import format from 'date-fns/format';
// import Datepicker from 'vuejs-datepicker';
import VueFlatpickr from 'vue-flatpickr';
import 'vue-flatpickr/theme/base16_flat.css';

import * as types from '../store/types';
import { hourMinute } from '../lib/utils';
import Rating from './Rating';

const en = {
  weekdays: {
    shorthand: ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'],
    longhand: [
      'Sunday', 'Monday', 'Tuesday', 'Wednesday',
      'Thursday', 'Friday', 'Saturday',
    ],
  },
  months: {
    shorthand: [
      'Jan', 'Feb', 'Mar', 'Apr',
      'May', 'Jun', 'Jul', 'Aug',
      'Sep', 'Oct', 'Nov', 'Dec',
    ],
    longhand: [
      'January', 'February', 'March', 'April',
      'May', 'June', 'July', 'August',
      'September', 'October', 'November', 'December',
    ],
  },
  daysInMonth: [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31],
  firstDayOfWeek: 1,
  ordinal: (nth) => {
    const s = nth % 100;
    if (s > 3 && s < 21) return 'th';
    switch (s % 10) {
      case 1: return 'st';
      case 2: return 'nd';
      case 3: return 'rd';
      default: return 'th';
    }
  },
  rangeSeparator: ' to ',
  weekAbbreviation: 'Wk',
  scrollTitle: 'Scroll to increment',
  toggleTitle: 'Click to toggle',
};


export default {
  name: 'movie',

  data() {
    return {
      tmdb: this.movie.tmdb_id,
      seen: new Date().toString(),
      fpOptions: {
        onValueUpdate: null,
        onChange: (selectedDates, dateStr, instance) => this.setWatched(selectedDates, dateStr, instance), // eslint-disable-line
        locale: en,
      },
      loading: false,
    };
  },

  props: {
    movie: {
      type: Object,
      required: true,
    },
  },

  components: { Rating, VueFlatpickr },

  methods: {
    setScore(score) {
      this.$store.dispatch(types.SET_SCORE, { id: this.movie.id, score });
    },

    setWatched() {
      // console.log(`seen(${this.seen})`); // eslint-disable-line
      const watched = format(this.seen, 'YYYY-MM-DDTHH:mm:ssZ');
      // console.log(`seen(${lastWatched})`); // eslint-disable-line
      this.$store.dispatch(types.SET_WATCHED, { id: this.movie.id, watched });
    },

    fixMovie() {
      // console.log(`typeof(${typeof this.tmdb})`); // eslint-disable-line
      if (this.tmdb !== this.movie.tmdb_id) {
        this.$store.dispatch(types.FIX_MOVIE, { id: this.movie.id, tmdb: this.tmdb });
      }
    },
  },

  computed: {
    watched() {
      return this.movie.count_watched > 0;
    },

    runtime() {
      return hourMinute(this.movie.runtime);
    },

    added() {
      // const added = parse(this.movie.added);
      return format(this.movie.added, 'MMM DD, YYYY H:mm');
    },

    hasRating() {
      return this.movie.score !== 0;
    },

    lastWatched() {
      return format(this.movie.last_watched, 'MMM DD, YYYY');
    },

    cover() {
      return `img/p${this.movie.cover}`;
    },

    background() {
      return `/img/b${this.movie.backdrop}`;
    },

    shows() {
      return this.movie.all_watched.split('|').map(show => format(show, 'MMM DD, YYYY'));
    },
  },
};
</script>

<style lang="scss" scoped>
@import "../styles/_settings.scss";

.c-minfo {
	font-size: 0.9em;
	margin-bottom: 2em;
	border-bottom: 1px solid $article-border-color;

  height: 765px;
  background-size: cover;

	// .director {
	// 	color: $director-color;
	// 	font-weight: bold;
	// }
}

.c-align {
  align-content: space-between;
}

.c-title {
  color: #fff;
  text-shadow: 0 1px 0 hsla(0,0%,0%,.75), 0 0 1px hsla(0,0%,0%,.75), 0 1px 5px hsla(0,0%,0%,.75);
  background: linear-gradient(to bottom, rgba(0,0,0,0.65) 0%,rgba(0,0,0,0) 100%);
}

.c-text {
  color: #fff;
  text-shadow: 0 1px 0 hsla(0,0%,0%,.75), 0 0 1px hsla(0,0%,0%,.75), 0 1px 5px hsla(0,0%,0%,.75);
}

.c-shaded {
	// font-size: 1.25em;
	color: white;
	background-color: rgba(0, 0, 0, 0.4);
	padding: 0.25em 1em;
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
