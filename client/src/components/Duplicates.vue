<template>
	<section>
    <movie v-for="movie in movies" :movie="movie" :key="movie.id"></movie>
	</section>
</template>

<script>
import Movie from './Movie'
import * as types from '../store/types'

export default {
  name: 'duplicates',

  components: { Movie },

  computed: {
    movies () {
      return this.$store.getters.getMovies
    }
  },

  created () {
    this.$store.dispatch(types.FETCH_DUPLICATES)
  }
}
</script>

<style>
</style>
