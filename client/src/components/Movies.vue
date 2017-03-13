<template>
	<section>
    <div class="row">
      <pager
        ref="paginator1"
        :pageCount="pageCount"
        :pageRange="5"
        :marginPages="3"
        :forcePage="forcePage"
        :containerClass="'pagination col-xs-12'"
        :clickHandler="onPaginationClick"
      />
    </div>

    <movie v-for="movie in movies" :movie="movie" :key="movie.id"></movie>

    <div class="row">
      <pager
        ref="paginator2"
        :pageCount="pageCount"
        :pageRange="5"
        :marginPages="3"
        :forcePage="forcePage"
        :containerClass="'pagination col-xs-12'"
        :clickHandler="onPaginationClick"
      />
    </div>
	</section>
</template>

<script>
// import Paginate from 'vuejs-paginate';

import Movie from './Movie';
import Pager from './Pager';
import * as types from '../store/types';

export default {
  name: 'movies',

  components: { Pager, Movie },

  data() {
    return {
      shouldScroll: false,
    };
  },

  methods: {
    onPaginationClick(pageNum) {
      this.shouldScroll = true;

      // let's keep both paginators in sync, they're zero-based
      // this.$refs.paginator1.selected = pageNum - 1;
      // this.$refs.paginator2.selected = pageNum - 1;

      const offset = Math.ceil((pageNum - 1) * this.$store.state.options.limit);
      this.$store.commit(types.SET_OFFSET, offset);
      this.$store.dispatch(types.FETCH_MOVIES);
    },

    // setPaginator(state) {
    //   console.log(`offset(${state.getters.offset})`); // eslint-disable-line
    // },
  },

  computed: {
    movies() {
      return this.$store.getters.getMovies;
    },

    pageCount() {
      return Math.ceil(this.$store.state.total / this.$store.state.options.limit);
    },

    forcePage() {
      return this.$store.state.options.offset / this.$store.state.options.limit;
    },
  },

  created() {
    // const refs = this.$refs;
    this.$store.dispatch(types.FETCH_MOVIES);
//     this.$store.watch(
//       function (state) { // eslint-disable-line
//         return state.options.offset;
//       },
//       function (newVal, oldVal) { // eslint-disable-line
//         // let's keep both paginators in sync, they're zero-based
//         refs.paginator1.selected = newVal + 1;
//         refs.paginator2.selected = newVal + 1;

// // console.log(`pag1(${refs.paginator1.selected})-old(${oldVal})-new(${newVal})`);
//       },
//     );
  },

  updated() {
    if (this.shouldScroll) {
      window.scrollTo(0, 0);
      this.shouldScroll = false;
    }
  },
};
</script>

<style>
</style>
