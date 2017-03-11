<template>
  <article class="row top-xs c-minfo" :style="{ backgroundImage: 'url(' + background + ')' }">
    <div class="col-xs-12">
      <div class="row">
        <h2 class="col-xs-12 c-text">{{movie.title}} ({{movie.year}})</h2>
      </div>
      <div class="row between-xs">
        <span class="col-xs-12 col-sm-6 director">{{movie.director}}</span>
        <span class="col-xs-12 col-sm-6 end-sm">{{movie.production_countries}}</span>
      </div>
    </div>
    <!--<div class="col-xs-12">
      <div class="row between-xs">
        <span class="col-xs-12 col-sm-6 director">{{movie.director}}</span>
        <span class="col-xs-12 col-sm-6 end-sm">{{movie.production_countries}}</span>
      </div>
    </div>-->
    <!--<div class="col-xs-12">
      <div class="row between-xs">
        <span class="col-xs-12 col-sm-6 ">{{movie.actors}}</span>
        <span class="col-xs-12 col-sm-6 end-sm">{{movie.genres}}</span>
      </div>
    </div>
    <div class="col-xs-12 bottom-spacer">
      <div class="row between-xs">
        <div class="col-xs-12 col-sm-9">
          <span class="label">{{movie.resolution}}</span>
          <span class="label secondary spacer">{{movie.location}}</span>
          <span class="label">{{movie.id}}</span>
        </div>
        <div class="col-xs-12 col-sm-3 end-sm">
          {{watched}}
          <span class="label"><i class="icon-plus"></i>&nbsp;{{movie.added}}</span>
        </div>
      </div>
    </div>
    <div class="col-xs-12 bottom-spacer">
      <span>{{movie.overview}}</span>
    </div>-->

  </article>
</template>

<script>
export default {
  name: 'movie',

  props: {
    movie: {
      type: Object,
      required: true,
    },
  },

  computed: {
    watched() {
      return this.movie.count_watched > 0;
    },

    background() {
      return `/img/b${this.movie.backdrop}`;
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

	.director {
		color: $director-color;
		font-weight: bold;
	}
}

.c-text {
  color: #fff;
  text-shadow: 0 1px 0 hsla(0,0%,0%,.75), 0 0 1px hsla(0,0%,0%,.75), 0 1px 5px hsla(0,0%,0%,.75);
  background: linear-gradient(to bottom, rgba(0,0,0,0.65) 0%,rgba(0,0,0,0) 100%);
}
</style>
