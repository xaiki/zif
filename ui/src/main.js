import Vue from "vue"
import VueRouter from "vue-router"
Vue.use(VueRouter);

import App from "./components/app.vue"
import Home from "./components/home.vue"
import Settings from "./components/settings.vue"
import Stream from "./components/stream-view.vue"

const router = new VueRouter({
	mode: "history",
	base: __dirname,
	routes: [
		{ path: "/", component: Home },
		{ path: "/settings", component: Settings},
		{ path: "/stream/:ih/:title/:size/:filecount", name: "stream", component: Stream }
	]
});

new Vue({
	router,
	render: h => h(App)
}).$mount("#app");

// Is this needed?
router.push("/");
