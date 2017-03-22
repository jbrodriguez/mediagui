<template>
	<nav class="row between-xs">
    <ul class="col-xs-12 col-sm-2 center-xs">
      <li class="pv2 ph0 c-hlogo">
        <router-link to="/">mediaGUI</router-link>
      </li>
    </ul>
    <ul class="col-xs-12 col-sm-10">
      <li class="pv2 ph0 c-hmenu">
        <div class="row between-xs">
          <div class="col-xs-12 col-sm-8">
            <div class="c-hmenusection">
              <router-link to="/movies" class="mv0 mh2">MOVIES</router-link>

              <select v-model="selected" @change="changeFilter">
                <option v-for="option in filters" :value="option.value">{{option.label}}</option>
              </select>

              <input type="search" placeholder="Enter search string" :value="query" @input="updateQuery">

              <select v-model="sortBy" @change="changeSort">
                <option v-for="sort in sorts" :value="sort.value">{{sort.label}}</option>
              </select>

              <i class="fa mv0 ml2" :class="chevron" @click="changeOrder"></i>

							<span class="mv0 mh2">|</span>

              <router-link to="/import" class="mv0">IMPORT</router-link>
            </div>
          </div>
          <div class="col-xs-12 col-sm-4 end-xs">
            <div class="c-hmenusection">
              <router-link to="/settings">SETTINGS</router-link>
							<span class="mv0 mh2">|</span>
              <router-link to="/duplicates">DUPLICATES</router-link>
              <router-link to="/prune" class="mv0 mh2">PRUNE</router-link>
            </div>
          </div>
        </div>
      </li>
    </ul>
	</nav>
</template>

<script>
import debounce from 'lodash.debounce'

import * as types from '../store/types'

export default {
  name: 'header-bar',

  data () {
    return {
      selected: this.$store.state.options.filterBy,
      sortBy: this.$store.state.options.sortBy
    }
  },

  methods: {
    changeFilter (e) {
      this.selected = e.target.value
      this.$store.commit(types.SET_FILTER, e.target.value)
      this.$store.dispatch(types.FETCH_MOVIES)
    },

    changeSort (e) {
      this.sortBy = e.target.value
      this.$store.commit(types.SET_SORT, e.target.value)
      this.$store.dispatch(types.FETCH_MOVIES)
    },

    updateQuery: debounce(
      function handle (e) {
        this.$store.commit(types.SET_QUERY, e.target.value)
        this.$store.dispatch(types.FETCH_MOVIES)
      },
      750,
    ),

    changeOrder () {
      this.$store.commit(types.FLIP_ORDER)
      this.$store.dispatch(types.FETCH_MOVIES)
    }
  },

  computed: {
    filters () {
      return this.$store.state.options.filterByOptions
    },

    sorts () {
      return this.$store.state.options.sortByOptions
    },

    query () {
      return this.$store.state.options.query
    },

    chevron () {
      return this.$store.state.options.sortOrder === 'asc' ? 'fa-chevron-circle-up' : 'fa-chevron-circle-down'
    }
  }
}
</script>

<style lang="scss" scoped>
@import "../styles/_settings.scss";

.c-hlogo {
	background-color: $headerLogoBackground;

	a {
		color: $headerLogoAnchorColor;
	}
}

.c-hmenu {
	background-color: $headerMenuBackground;
}

.c-hmenusection {
	// background-color: $headerMenuBackground;

	a {
		font-size: 0.9em;
		// color: $headerMenuAnchorColor;

		&.router-link-active {
			padding-bottom: 3px;
			border-bottom: 2px solid $headerMenuAnchorColor;
		}
	}
}
</style>
