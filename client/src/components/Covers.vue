<template>
	<section class="row">
    <article v-for="movie in movies" class="col-xs-12 col-sm-6 col-md-3 col-lg-2 mb3">
      <cover :movie="movie" />
    </article>
	</section>
</template>

<script>
import Cover from './Cover';
import * as types from '../store/types';

export default {
  name: 'covers',

  components: { Cover },

  created() {
    this.$store.dispatch(types.FETCH_COVER);
  },

  computed: {
    movies() {
      return this.$store.getters.getMovies;
    },
  },
};
</script>

<style>
</style>
