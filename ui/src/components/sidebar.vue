<template>
	<div class="ui visible inverted left vertical sidebar menu icon very thin zifmenu">
		<router-link class="item" to="/"><i class="large home icon"></i></router-link>
		<router-link class="item" to="/subscribe"><i class="large add user icon"></i></router-link>
		<router-link class="item" to="/settings"><i class="large setting icon"></i></router-link>
	</div>
</template>

<style lang="sass">
	@import '../colours';

	.zifmenu {
		background-color: $darkBackground !important;
	}
</style>
