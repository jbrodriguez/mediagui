<template>
	<section class="row">
    <article v-for="movie in movies" class="col-xs-12 col-sm-6 col-md-3 col-lg-2 mb3">
      <cover :movie="movie" />
    </article>
	</section>
</template>

<script>
import Cover from './Cover';

export default {
  name: 'covers',

  components: { Cover },

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
};
</script>

<style>
</style>
