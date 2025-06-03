import { w as writable } from './index-Bt5_4Bzy.js';
import './client-CjdeEz1m.js';

const Terms = writable([]);
let UnitedStates_selected = writable(false);
let today = /* @__PURE__ */ new Date();
new Date(today.getFullYear(), today.getMonth(), 1);
new Date(today.getFullYear(), today.getMonth(), today.getDate());
const clickedSidebarMenu = writable("main");
const backendUrl = writable("");
const backendUrlGolang = writable("https://rentacarbackendgo-dd6070431ef2.herokuapp.com");
let soundEnabled = writable(true);
let recordIdsVisible = writable(false);
let apiCallsVisible = writable(false);
let clientIp = writable("");
let imagesStillLoading = writable(false);
let reportedErrorMessage = writable("");
let workgroupReady = writable(true);
let isLanguageChosen = writable(true);
const selectedLanguage = () => {
  let language = "portuguese";
  UnitedStates_selected.subscribe((UnitedStates_selected2) => {
    if (UnitedStates_selected2)
      language = "english";
  });
  return language;
};
const simplifyLanguageData = (jsonResponse) => {
  let _language = {};
  let i;
  for (i = 0; i < jsonResponse.length; i++) {
    let vlr = jsonResponse[i];
    _language[`${vlr["item"]}`] = `${vlr["expression"]}`;
  }
  return _language;
};

const load = async ({ fetch }) => {
  let _backendUrlGolang_;
  backendUrlGolang.subscribe((url) => {
    _backendUrlGolang_ = url;
  });
  let language = selectedLanguage();
  const languageRes = await fetch(`${_backendUrlGolang_}/terms/${language}`, { method: "GET" });
  const languageData = await languageRes.json();
  return {
    languageData
  };
};

var _page = /*#__PURE__*/Object.freeze({
  __proto__: null,
  load: load
});

const index = 2;
let component_cache;
const component = async () => component_cache ??= (await import('./_page.svelte-BnyvIl9O.js')).default;
const universal_id = "src/routes/+page.js";
const imports = ["_app/immutable/nodes/2.BD4TgjEL.js","_app/immutable/chunks/scheduler.BcErCovF.js","_app/immutable/chunks/entry.BAu7sNJh.js","_app/immutable/chunks/index.BHrIfvv6.js"];
const stylesheets = ["_app/immutable/assets/2.PMtF30yU.css"];
const fonts = [];

var _2 = /*#__PURE__*/Object.freeze({
  __proto__: null,
  component: component,
  fonts: fonts,
  imports: imports,
  index: index,
  stylesheets: stylesheets,
  universal: _page,
  universal_id: universal_id
});

export { Terms as T, UnitedStates_selected as U, _2 as _, clickedSidebarMenu as a, backendUrlGolang as b, clientIp as c, apiCallsVisible as d, soundEnabled as e, recordIdsVisible as f, backendUrl as g, imagesStillLoading as h, isLanguageChosen as i, reportedErrorMessage as r, simplifyLanguageData as s, workgroupReady as w };
//# sourceMappingURL=2-BshwT9cA.js.map
