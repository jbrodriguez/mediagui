<template>
  <div class="c-console">
    <div class="row">
      <div class="col-xs-12">
        <p v-for="line in lines" class="console__line">{{line}}</p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'console',

  computed: {
    lines() {
      return this.$store.state.lines;
    },
  },

  updated() {
    this.$el.scrollTop = this.$el.scrollHeight;
  },
};
</script>

<style lang="scss" scoped>
@import "../styles/_settings.scss";

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
