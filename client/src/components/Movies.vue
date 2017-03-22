<template>
	<section v-if="total > 0">
    <div class="row between-xs middle-xs mb2">
      <div class="col-xs-12 col-sm-10 middle-xs">
        <pager
          :pageCount="pageCount"
          :pageRange="5"
          :marginPages="3"
          :forcePage="forcePage"
          :containerClass="'pagination'"
          :clickHandler="onPaginationClick"
        />
      </div>
      <div class="col-xs-12 col-sm-2 end-xs">
        <span>TOTAL </span>
        <span class="c-total">{{ total }}</span>
      </div>
    </div>

    <movie v-for="movie in movies" :movie="movie" :key="movie.id"></movie>

    <div class="row">
      <div class="col-xs-12 middle-xs">
        <pager
          :pageCount="pageCount"
          :pageRange="5"
          :marginPages="3"
          :forcePage="forcePage"
          :containerClass="'pagination'"
          :clickHandler="onPaginationClick"
        />
      </div>
    </div>
	</section>
</template>

<script>
// import Paginate from 'vuejs-paginate';

import Movie from './Movie'
import Pager from './Pager'
import * as types from '../store/types'

export default {
  name: 'movies',

  components: { Pager, Movie },

  data () {
    return {
      shouldScroll: false
    }
  },

  methods: {
    onPaginationClick (pageNum) {
      this.shouldScroll = true

      // let's keep both paginators in sync, they're zero-based
      // this.$refs.paginator1.selected = pageNum - 1;
      // this.$refs.paginator2.selected = pageNum - 1;

      const offset = Math.ceil((pageNum - 1) * this.$store.state.options.limit)
      this.$store.commit(types.SET_OFFSET, offset)
      this.$store.dispatch(types.FETCH_MOVIES)
    }
  },

  computed: {
    movies () {
      return this.$store.getters.getMovies
    },

    pageCount () {
      return Math.ceil(this.$store.state.total / this.$store.state.options.limit)
    },

    forcePage () {
      return this.$store.state.options.offset / this.$store.state.options.limit
    },

    total () {
      return this.$store.state.total
    },

    visible () {
      return this.$store.state.total > 0
    }
  },

  created () {
    this.$store.dispatch(types.FETCH_MOVIES)
  },

  updated () {
    if (this.shouldScroll) {
      window.scrollTo(0, 0)
      this.shouldScroll = false
    }
  }
}
</script>

<style lang="scss" scoped>
@import "../styles/_settings.scss";

.c-total {
  color: $primaryHeaderColor;
}

</style>
