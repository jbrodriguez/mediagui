<template>
	<section class="row">
    <article v-for="movie in movies" class="col-xs-12 col-sm-6 col-md-3 col-lg-2">
      <img :src="'/img/p' + movie.cover" />
    </article>
	</section>
</template>

<script>
import * as types from '../store/types';

export default {
  name: 'cover',

  data() {
    const data = {
      // movies: this.$store.getters.getMovies,
    };

    return data;
  },

  computed: {
    movies() {
      return this.$store.getters.getMovies;
    },
  },

  created() {
    this.$store.dispatch(types.FETCH_COVER);
  },
};
</script>

<style lang="scss" scoped>
@import "../styles/_settings.scss";
</style>
