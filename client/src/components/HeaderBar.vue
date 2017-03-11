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
              <router-link to="/movies" class="mv0 mh3">MOVIES</router-link>

              <select v-model="selected" @change="changeFilter">
                <option v-for="option in filters" :value="option.value">{{option.label}}</option>
              </select>
            </div>
          </div>
        </div>
      </li>
    </ul>
	</nav>
</template>

<script>
import * as types from '../store/types';

export default {
  name: 'header-bar',

  data() {
    return {
      selected: this.$store.state.options.filterBy,
    };
  },

  methods: {
    changeFilter(e) {
      this.selected = e.target.value;
      this.$store.commit(types.SET_FILTER, e.target.value);
    },
  },

  computed: {
    filters() {
      return this.$store.state.options.filterByOptions;
    },
  },
};
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
