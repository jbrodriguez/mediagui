<template>
  <div class="">
    <section class="row">
      <div class="col-xs-12 mb2">
        <button class="btn btn-default" @click="runImport">Import</button>
      </div>
    </section>
    <section class="row">
      <div class="col-xs-12">
        <div class="c-console">
          <div class="row">
            <div class="col-xs-12">
              <p v-for="line in lines" class="console__line">{{line}}</p>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import * as types from '../store/types';

export default {
  name: 'import',

  // components: { Movie },

  methods: {
    runImport() {
      this.$store.dispatch(types.RUN_IMPORT);
    },
  },

  computed: {
    movies() {
      return this.$store.getters.getMovies;
    },

    lines() {
      return this.$store.state.lines;
    },
  },

  created() {
    this.$store.dispatch(types.FETCH_MOVIES);
  },

  updated() {
    this.$el.scrollTop = this.$el.scrollHeight;
  },
};
</script>

<style lang="scss" scoped>
@import "../styles/_settings.scss";

// .fparent {
//   position: relative
// }

// .fchild {
//   position: absolute;
//   height: 100%;
// }

// .flex {
//   display: flex;
// }

// .flx_i {
//   flex: 1;
// }

.c-console {
	display: flex;
  overflow: auto;
  // width: 100%;
  height: 750px;
  // height: 100%;
  // border: 1px solid $console-border;
  padding: 0.75em;
  color: $console-color;
  background-color: $console-bg;
  // flex: 1;

	.console__line {
		font-size: 0.8em;
		// line-height: 0.5em;
	}

  p {
    margin-bottom: 0;
    white-space: nowrap;
  }
}
</style>
