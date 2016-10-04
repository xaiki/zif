import Vue from "vue"
import VueRouter from "vue-router"
Vue.use(VueRouter);

import App from "./components/app.vue"
import Home from "./components/home.vue"
import Resolve from "./components/resolve.vue"

const router = new VueRouter({
	mode: "history",
	base: __dirname,
	routes: [
		{ path: "/", component: Home },
		{ path: "/resolve", component: Resolve }
	]
});

new Vue({
	router,
	render: h => h(App)
})t("#app");

// Is this needed?
router.push("/");
