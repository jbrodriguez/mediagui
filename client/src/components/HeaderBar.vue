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
            <span class="mh3">MOVIES</span>
          </div>
        </div>
      </li>
    </ul>
	</nav>
</template>

<script>
export default {
  name: 'header-bar',
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
</style>
