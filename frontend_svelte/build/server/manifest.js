const manifest = (() => {
function __memo(fn) {
	let value;
	return () => value ??= (value = fn());
}

return {
	appDir: "_app",
	appPath: "_app",
	assets: new Set(["brazil_flag.svg","css/jquery-ui.min.css","css/jquery-ui.min.css:Zone.Identifier","css/spin.css","css/spin.css:Zone.Identifier","css_symbol.svg","enter_arrows.png","error_beep.mp3","favicon.png","golang_symbol.svg","html_symbol.svg","javascript_symbol.svg","jquery_symbol.svg","js/jquery-ui.min.js","js/jquery.autocomplete.js","js/loadingGif.js","js/maskDateHour.js","js/maskMONEY.js","js/multiDraggable.js","loading.gif","menu_item_bookings_blue.svg","menu_item_bookings_gray.svg","menu_item_cars_blue.svg","menu_item_cars_gray.svg","menu_item_home_blue.svg","menu_item_home_gray.svg","menu_item_manufacturers_blue.svg","menu_item_manufacturers_gray.svg","menu_item_terms_blue.svg","menu_item_terms_gray.svg","menu_item_workgroups_blue.svg","menu_item_workgroups_gray.svg","nodejs_symbol.svg","none_image.png","notifications_.png","php_symbol.svg","svelte_symbol.svg","tailwind_symbol.svg","usa_flag.svg"]),
	mimeTypes: {".svg":"image/svg+xml",".css":"text/css",".png":"image/png",".mp3":"audio/mpeg",".js":"text/javascript",".gif":"image/gif"},
	_: {
		client: {"start":"_app/immutable/entry/start.CCXbWjgG.js","app":"_app/immutable/entry/app.Cwf8C1oQ.js","imports":["_app/immutable/entry/start.CCXbWjgG.js","_app/immutable/chunks/entry.BAu7sNJh.js","_app/immutable/chunks/scheduler.BcErCovF.js","_app/immutable/entry/app.Cwf8C1oQ.js","_app/immutable/chunks/scheduler.BcErCovF.js","_app/immutable/chunks/index.BHrIfvv6.js"],"stylesheets":[],"fonts":[],"uses_env_dynamic_public":false},
		nodes: [
			__memo(() => import('./chunks/0-DI0lnWpg.js')),
			__memo(() => import('./chunks/1-CsYoAkll.js')),
			__memo(() => import('./chunks/2-BshwT9cA.js').then(function (n) { return n._; }))
		],
		routes: [
			{
				id: "/",
				pattern: /^\/$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 2 },
				endpoint: null
			}
		],
		matchers: async () => {
			
			return {  };
		},
		server_assets: {}
	}
}
})();

const prerendered = new Set([]);

const base = "";

export { base, manifest, prerendered };
//# sourceMappingURL=manifest.js.map
