<template>
	<section class="row c-flegal between-xs middle-xs">
    <ul class="col-xs-12 col-sm-4">
      <div>
        <span class="mv0 mh3">Copyright &copy;</span>
        <a href="http://jbrodriguez.io">Juan B. Rodriguez</a>
      </div>
    </ul>
    <ul class="col-xs-12 col-sm-4">
      <div class="center-xs">
        <span>mediaGUI &nbsp;</span>
        <span>v{{version}}</span>
      </div>
    </ul>
    <ul class="col-xs-12 col-sm-4 end-xs middle-xs">
      <a class="end-xs middle-xs mv0 mh3" href="http://jbrodriguez.io" title="jbrodriguez.io" target="_blank">
        <img src="../assets/logo-small.png" alt="Logo for Juan B. Rodriguez" />
      </a>
    </ul>
  </section>
</template>

<script>
export default {
  name: 'footer-bar',

  computed: {
    version () {
      return this.$store.getters.version
    }
  }
}
</script>

<style lang="scss" scoped>
@import "../styles/_settings.scss";

.c-flegal {

	a, img, span, div {
		display: flex;
	}

	a {
		color: $footer-anchor-color;
	}
}
</style>
