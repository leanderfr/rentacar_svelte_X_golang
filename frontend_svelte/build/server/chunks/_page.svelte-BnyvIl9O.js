import { a as subscribe, b as set_store_value } from './utils-ByOLlbjM.js';
import { c as create_ssr_component, a as add_attribute, e as escape, b as each, v as validate_component, d as add_classes, f as createEventDispatcher } from './ssr-CTylFKGs.js';
import { s as simplifyLanguageData, T as Terms, U as UnitedStates_selected, r as reportedErrorMessage, c as clientIp, b as backendUrlGolang, i as isLanguageChosen, a as clickedSidebarMenu, d as apiCallsVisible, e as soundEnabled, f as recordIdsVisible, g as backendUrl, w as workgroupReady, h as imagesStillLoading } from './2-BshwT9cA.js';
import 'jquery';
import 'moment';
import 'spin.js';
import './client-CjdeEz1m.js';
import './index-Bt5_4Bzy.js';
import './exports-DuWZopOC.js';

const WorkgroupCloningModal = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $Terms, $$unsubscribe_Terms;
  $$unsubscribe_Terms = subscribe(Terms, (value) => $Terms = value);
  let { userConfirmedWorkgroupDataReset } = $$props;
  if ($$props.userConfirmedWorkgroupDataReset === void 0 && $$bindings.userConfirmedWorkgroupDataReset && userConfirmedWorkgroupDataReset !== void 0)
    $$bindings.userConfirmedWorkgroupDataReset(userConfirmedWorkgroupDataReset);
  $$unsubscribe_Terms();
  return `<div class="w-full h-full absolute flex left-0 top-0 z-40 bg-[rgba(0,0,0,0.1)] cursor-pointer justify-center items-center"><div class="flex flex-col w-[50%] bg-white relative rounded-lg h-48 border-blue-600 border-[2px] "><div class="flex flex-col h-36 justify-center items-center gap-4" data-svelte-h="svelte-d3q258"><div class="text-2xl " style="color:blue" id="workgroupCloningStatus"> </div> <div class="divLOADING_BAR_OUTER "><div class="divLOADING_BAR_INNER justify-center " id="workgroupCloningPercent"></div></div></div> <div class="flex flex-row justify-between "><div class="flex flex-row text-2xl pl-3">${userConfirmedWorkgroupDataReset ? `<div class="flex justify-start h-4 text-left"> ${escape($Terms.workgroup_data_resetting)}</div>` : `<div class="flex justify-start h-4 text-left"> ${escape($Terms.chosen_workgroup)}</div>`} <div class="flex justify-start h-4 text-left pl-4" style="color:blue" id="chosenWorkgroupName" data-svelte-h="svelte-dzbl8c"> </div></div> <div class="flex flex-row text-2xl pr-5" style="color:blue" id="completedPercent" data-svelte-h="svelte-1yff8gp"> </div></div></div></div>`;
});
const LanguageChoice = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $$unsubscribe_isLanguageChosen;
  $$unsubscribe_isLanguageChosen = subscribe(isLanguageChosen, (value) => value);
  let proceedTranslation;
  createEventDispatcher();
  proceedTranslation = "Prosseguir";
  $$unsubscribe_isLanguageChosen();
  return `<div class="w-full h-full absolute flex left-0 top-0 z-40 bg-[rgba(0,0,0,0.2)] cursor-pointer justify-center items-center"><div class="flex flex-col w-[500px] bg-white relative rounded-lg h-48 border-gray-700 border-[1px] justify-around items-center px-7 py-5 "><div class="flex flex-row w-full bg-white gap-4"> <div id="flagBRAZIL" class="${["divAsButton", ""].join(" ").trim()}" aria-hidden="true" data-svelte-h="svelte-1qmab2c"><div class="text-lg">Português</div> <svg xmlns="http://www.w3.org/2000/svg" width="62" height="62" viewBox="0 0 32 32"><rect x="1" y="4" width="30" height="24" rx="4" ry="4" fill="#459a45"></rect><path d="M27,4H5c-2.209,0-4,1.791-4,4V24c0,2.209,1.791,4,4,4H27c2.209,0,4-1.791,4-4V8c0-2.209-1.791-4-4-4Zm3,20c0,1.654-1.346,3-3,3H5c-1.654,0-3-1.346-3-3V8c0-1.654,1.346-3,3-3H27c1.654,0,3,1.346,3,3V24Z" opacity=".15"></path><path d="M3.472,16l12.528,8,12.528-8-12.528-8L3.472,16Z" fill="#fedf00"></path><circle cx="16" cy="16" r="5" fill="#0a2172"></circle><path d="M14,14.5c-.997,0-1.958,.149-2.873,.409-.078,.35-.126,.71-.127,1.083,.944-.315,1.951-.493,2.999-.493,2.524,0,4.816,.996,6.519,2.608,.152-.326,.276-.666,.356-1.026-1.844-1.604-4.245-2.583-6.875-2.583Z" fill="#fff"></path><path d="M27,5H5c-1.657,0-3,1.343-3,3v1c0-1.657,1.343-3,3-3H27c1.657,0,3,1.343,3,3v-1c0-1.657-1.343-3-3-3Z" fill="#fff" opacity=".2"></path></svg></div>  <div id="flagBRAZIL" class="${["divAsButton", ""].join(" ").trim()}" aria-hidden="true" data-svelte-h="svelte-189gyyh"><div class="text-lg">English</div> <svg xmlns="http://www.w3.org/2000/svg" width="62" height="62" viewBox="0 0 32 32"><rect x="1" y="4" width="30" height="24" rx="4" ry="4" fill="#fff"></rect><path d="M1.638,5.846H30.362c-.711-1.108-1.947-1.846-3.362-1.846H5c-1.414,0-2.65,.738-3.362,1.846Z" fill="#a62842"></path><path d="M2.03,7.692c-.008,.103-.03,.202-.03,.308v1.539H31v-1.539c0-.105-.022-.204-.03-.308H2.03Z" fill="#a62842"></path><path fill="#a62842" d="M2 11.385H31V13.231H2z"></path><path fill="#a62842" d="M2 15.077H31V16.923000000000002H2z"></path><path fill="#a62842" d="M1 18.769H31V20.615H1z"></path><path d="M1,24c0,.105,.023,.204,.031,.308H30.969c.008-.103,.031-.202,.031-.308v-1.539H1v1.539Z" fill="#a62842"></path><path d="M30.362,26.154H1.638c.711,1.108,1.947,1.846,3.362,1.846H27c1.414,0,2.65-.738,3.362-1.846Z" fill="#a62842"></path><path d="M5,4h11v12.923H1V8c0-2.208,1.792-4,4-4Z" fill="#102d5e"></path><path d="M27,4H5c-2.209,0-4,1.791-4,4V24c0,2.209,1.791,4,4,4H27c2.209,0,4-1.791,4-4V8c0-2.209-1.791-4-4-4Zm3,20c0,1.654-1.346,3-3,3H5c-1.654,0-3-1.346-3-3V8c0-1.654,1.346-3,3-3H27c1.654,0,3,1.346,3,3V24Z" opacity=".15"></path><path d="M27,5H5c-1.657,0-3,1.343-3,3v1c0-1.657,1.343-3,3-3H27c1.657,0,3,1.343,3,3v-1c0-1.657-1.343-3-3-3Z" fill="#fff" opacity=".2"></path><path fill="#fff" d="M4.601 7.463L5.193 7.033 4.462 7.033 4.236 6.338 4.01 7.033 3.279 7.033 3.87 7.463 3.644 8.158 4.236 7.729 4.827 8.158 4.601 7.463z"></path><path fill="#fff" d="M7.58 7.463L8.172 7.033 7.441 7.033 7.215 6.338 6.989 7.033 6.258 7.033 6.849 7.463 6.623 8.158 7.215 7.729 7.806 8.158 7.58 7.463z"></path><path fill="#fff" d="M10.56 7.463L11.151 7.033 10.42 7.033 10.194 6.338 9.968 7.033 9.237 7.033 9.828 7.463 9.603 8.158 10.194 7.729 10.785 8.158 10.56 7.463z"></path><path fill="#fff" d="M6.066 9.283L6.658 8.854 5.927 8.854 5.701 8.158 5.475 8.854 4.744 8.854 5.335 9.283 5.109 9.979 5.701 9.549 6.292 9.979 6.066 9.283z"></path><path fill="#fff" d="M9.046 9.283L9.637 8.854 8.906 8.854 8.68 8.158 8.454 8.854 7.723 8.854 8.314 9.283 8.089 9.979 8.68 9.549 9.271 9.979 9.046 9.283z"></path><path fill="#fff" d="M12.025 9.283L12.616 8.854 11.885 8.854 11.659 8.158 11.433 8.854 10.702 8.854 11.294 9.283 11.068 9.979 11.659 9.549 12.251 9.979 12.025 9.283z"></path><path fill="#fff" d="M6.066 12.924L6.658 12.494 5.927 12.494 5.701 11.799 5.475 12.494 4.744 12.494 5.335 12.924 5.109 13.619 5.701 13.19 6.292 13.619 6.066 12.924z"></path><path fill="#fff" d="M9.046 12.924L9.637 12.494 8.906 12.494 8.68 11.799 8.454 12.494 7.723 12.494 8.314 12.924 8.089 13.619 8.68 13.19 9.271 13.619 9.046 12.924z"></path><path fill="#fff" d="M12.025 12.924L12.616 12.494 11.885 12.494 11.659 11.799 11.433 12.494 10.702 12.494 11.294 12.924 11.068 13.619 11.659 13.19 12.251 13.619 12.025 12.924z"></path><path fill="#fff" d="M13.539 7.463L14.13 7.033 13.399 7.033 13.173 6.338 12.947 7.033 12.216 7.033 12.808 7.463 12.582 8.158 13.173 7.729 13.765 8.158 13.539 7.463z"></path><path fill="#fff" d="M4.601 11.104L5.193 10.674 4.462 10.674 4.236 9.979 4.01 10.674 3.279 10.674 3.87 11.104 3.644 11.799 4.236 11.369 4.827 11.799 4.601 11.104z"></path><path fill="#fff" d="M7.58 11.104L8.172 10.674 7.441 10.674 7.215 9.979 6.989 10.674 6.258 10.674 6.849 11.104 6.623 11.799 7.215 11.369 7.806 11.799 7.58 11.104z"></path><path fill="#fff" d="M10.56 11.104L11.151 10.674 10.42 10.674 10.194 9.979 9.968 10.674 9.237 10.674 9.828 11.104 9.603 11.799 10.194 11.369 10.785 11.799 10.56 11.104z"></path><path fill="#fff" d="M13.539 11.104L14.13 10.674 13.399 10.674 13.173 9.979 12.947 10.674 12.216 10.674 12.808 11.104 12.582 11.799 13.173 11.369 13.765 11.799 13.539 11.104z"></path><path fill="#fff" d="M4.601 14.744L5.193 14.315 4.462 14.315 4.236 13.619 4.01 14.315 3.279 14.315 3.87 14.744 3.644 15.44 4.236 15.01 4.827 15.44 4.601 14.744z"></path><path fill="#fff" d="M7.58 14.744L8.172 14.315 7.441 14.315 7.215 13.619 6.989 14.315 6.258 14.315 6.849 14.744 6.623 15.44 7.215 15.01 7.806 15.44 7.58 14.744z"></path><path fill="#fff" d="M10.56 14.744L11.151 14.315 10.42 14.315 10.194 13.619 9.968 14.315 9.237 14.315 9.828 14.744 9.603 15.44 10.194 15.01 10.785 15.44 10.56 14.744z"></path><path fill="#fff" d="M13.539 14.744L14.13 14.315 13.399 14.315 13.173 13.619 12.947 14.315 12.216 14.315 12.808 14.744 12.582 15.44 13.173 15.01 13.765 15.44 13.539 14.744z"></path></svg></div></div> <div class="w-full flex justify-end pr-3 pt-7 "><button class="${[
    "btnCONTINUE ",
    " invisible"
  ].join(" ").trim()}" aria-hidden="true" style="max-width:100px;">${escape(proceedTranslation)}</button></div></div></div>`;
});
const ReportedError = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $$unsubscribe_imagesStillLoading;
  let $Terms, $$unsubscribe_Terms;
  let $reportedErrorMessage, $$unsubscribe_reportedErrorMessage;
  $$unsubscribe_imagesStillLoading = subscribe(imagesStillLoading, (value) => value);
  $$unsubscribe_Terms = subscribe(Terms, (value) => $Terms = value);
  $$unsubscribe_reportedErrorMessage = subscribe(reportedErrorMessage, (value) => $reportedErrorMessage = value);
  $$unsubscribe_imagesStillLoading();
  $$unsubscribe_Terms();
  $$unsubscribe_reportedErrorMessage();
  return `<div class="flex flex-col w-[80%] bg-white relative rounded-lg h-[50%]"><div class="flex flex-row w-full h-full"><div class="items-center bg-center bg-[length:80%_80%] justify-center w-[40%] bg-no-repeat bg-icon-reported-error" data-svelte-h="svelte-141h4xl"></div> <div class="w-[60%] text-red-800 text-[20px] flex flex-col justify-center gap-y-11"><div>${escape($Terms.there_was_js_error_line1)}</div> <div>${escape($Terms.there_was_js_error_line2)}</div> <div>${escape($Terms.there_was_js_error_line3)}</div> <div class="w-full flex flex-row">${escape($Terms.there_was_js_error_line4)}=  <div class="text-red-600">${escape($reportedErrorMessage)}</div></div></div></div></div>`;
});
const css = {
  code: "@import '$public/tailwind.css';",
  map: null
};
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let buttons_sidebar_menu;
  let $UnitedStates_selected, $$unsubscribe_UnitedStates_selected;
  let $reportedErrorMessage, $$unsubscribe_reportedErrorMessage;
  let $Terms, $$unsubscribe_Terms;
  let $$unsubscribe_clientIp;
  let $$unsubscribe_backendUrlGolang;
  let $isLanguageChosen, $$unsubscribe_isLanguageChosen;
  let $clickedSidebarMenu, $$unsubscribe_clickedSidebarMenu;
  let $apiCallsVisible, $$unsubscribe_apiCallsVisible;
  let $soundEnabled, $$unsubscribe_soundEnabled;
  let $recordIdsVisible, $$unsubscribe_recordIdsVisible;
  let $$unsubscribe_backendUrl;
  let $workgroupReady, $$unsubscribe_workgroupReady;
  $$unsubscribe_UnitedStates_selected = subscribe(UnitedStates_selected, (value) => $UnitedStates_selected = value);
  $$unsubscribe_reportedErrorMessage = subscribe(reportedErrorMessage, (value) => $reportedErrorMessage = value);
  $$unsubscribe_Terms = subscribe(Terms, (value) => $Terms = value);
  $$unsubscribe_clientIp = subscribe(clientIp, (value) => value);
  $$unsubscribe_backendUrlGolang = subscribe(backendUrlGolang, (value) => value);
  $$unsubscribe_isLanguageChosen = subscribe(isLanguageChosen, (value) => $isLanguageChosen = value);
  $$unsubscribe_clickedSidebarMenu = subscribe(clickedSidebarMenu, (value) => $clickedSidebarMenu = value);
  $$unsubscribe_apiCallsVisible = subscribe(apiCallsVisible, (value) => $apiCallsVisible = value);
  $$unsubscribe_soundEnabled = subscribe(soundEnabled, (value) => $soundEnabled = value);
  $$unsubscribe_recordIdsVisible = subscribe(recordIdsVisible, (value) => $recordIdsVisible = value);
  $$unsubscribe_backendUrl = subscribe(backendUrl, (value) => value);
  $$unsubscribe_workgroupReady = subscribe(workgroupReady, (value) => $workgroupReady = value);
  let { data } = $$props;
  if ($$props.data === void 0 && $$bindings.data && data !== void 0)
    $$bindings.data(data);
  $$result.css.add(css);
  let $$settled;
  let $$rendered;
  let previous_head = $$result.head;
  do {
    $$settled = true;
    $$result.head = previous_head;
    set_store_value(Terms, $Terms = simplifyLanguageData(data.languageData), $Terms);
    {
      setTimeout(
        () => {
        },
        1e3
      );
    }
    buttons_sidebar_menu = [
      {
        name: $Terms.itemmenu_main,
        id: "main",
        gray_icon: "menu_item_home_gray.svg",
        blue_icon: "menu_item_home_blue.svg",
        division_admin: ""
      },
      {
        name: $Terms.itemmenu_booking,
        id: "booking",
        gray_icon: "menu_item_bookings_gray.svg",
        blue_icon: "menu_item_bookings_blue.svg",
        division_admin: ""
      },
      {
        name: $Terms.itemmenu_cars,
        id: "cars",
        gray_icon: "menu_item_cars_gray.svg",
        blue_icon: "menu_item_cars_blue.svg",
        division_admin: "yes"
        // indica que este icone ja faz parte da area adminitrativa do menu
      },
      {
        name: $Terms.itemmenu_manufacturers,
        id: "manufacturers",
        gray_icon: "menu_item_manufacturers_gray.svg",
        blue_icon: "menu_item_manufacturers_blue.svg",
        division_admin: ""
      },
      {
        name: $Terms.itemmenu_languages,
        id: "languages",
        gray_icon: "menu_item_terms_gray.svg",
        blue_icon: "menu_item_terms_blue.svg",
        division_admin: ""
      },
      {
        name: $Terms.itemmenu_workgroups,
        id: "workgroups",
        gray_icon: "menu_item_workgroups_gray.svg",
        blue_icon: "menu_item_workgroups_blue.svg",
        division_admin: ""
      }
    ];
    $$rendered = `   ${$$result.head += `<!-- HEAD_undefined_START --><!-- HEAD_undefined_END -->`, ""}  ${``}  ${``}  ${``}  ${``}  ${``}  ${``}  ${``}  <audio id="errorBeep" src="error_beep.mp3" preload="auto" autobuffer></audio>  <div id="divLoading" class="cssLOADING_HIDE"></div> <div class="flex h-screen select-none font-Roboto text-sm flex-row ">   <div class="bg-gray-100 w-[15%] h-full "> <div class="flex items-center cursor-pointer hover:bg-gray-200 border-b-2 border-gray-300 h-[105px] border-b-gray-400 " aria-hidden="true"><div class="w-full "><div class="w-full flex flex-col align-middle py-2 h-16 border-b-2 border-gray-300" data-svelte-h="svelte-1aryqd0"><div class="text-lg pl-3 "><span class="text-red-700 ">RENT</span> A CAR</div> <div class="block pl-3 " aria-hidden="true">WebApp Demonstration</div></div>  <div class="backendSelector"${add_attribute("title", $Terms.change_backend, 0)} id="divBackendChoice" aria-hidden="true"><div style="font-size:14px;margin-top:-5px">${escape($Terms.current_backend)}</div> <div style="padding-left:10px;width:50px;height:50px;" data-svelte-h="svelte-1htkj2g"><img style="margin-top:0px" alt="" id="backendIcon"></div></div></div></div>  <div class="pt-6">${each(buttons_sidebar_menu, ({ id, name, gray_icon, blue_icon, division_admin }) => {
      return ` ${division_admin == "yes" ? `<div class="pl-5 border-b-2 border-b-gray-300 mt-10 mb-2 pb-2 mx-auto w-11/12 flex">${escape($Terms.itemmenu_administration)}</div>` : ``}  ${id == $clickedSidebarMenu ? `<div class="btn_sidebar_selected" aria-hidden="true"><div class="ml-4 -mt-1 "><img${add_attribute("src", blue_icon, 0)} alt=""></div> <span class="text-blue-700">${escape(name)}</span></div> ` : `<div class="btn_sidebar" aria-hidden="true"><div class="ml-4 -mt-1"><img${add_attribute("src", gray_icon, 0)} alt=""></div> <span class="text-gray-600">${escape(name)}</span> </div>`}`;
    })}</div></div>    <div class="flex flex-col w-[85%] pr-1 h-full grow overflow-hidden "> <div class="w-full flex flex-row items-center static px-6 h-[93px] "> <div class="grow text-xl font-bold">${escape($Terms.welcome)}</div>  <div class="flex justify-end items-center h-14 gap-4"> <div class="flex flex-row gap-4"> <div aria-hidden="true" id="flagBRAZIL"${add_attribute("_originaltooltip_", $Terms.change_country, 0)} class="${[
      "containsTooltip",
      (!$UnitedStates_selected ? "flagClicked" : "") + " " + ($UnitedStates_selected ? "flagUnclicked" : "")
    ].join(" ").trim()}"><img src="brazil_flag.svg" alt=""></div> <label for="chkLanguageSelector" class="switch_language containsTooltip"${add_attribute("_originaltooltip_", $Terms.change_country, 0)}><input id="chkLanguageSelector" type="checkbox"${add_attribute("checked", $UnitedStates_selected, 1)}> <span class="slider_language round"></span></label>  <div aria-hidden="true" id="flagUSA"${add_attribute("_originaltooltip_", $Terms.change_country, 0)} class="${[
      "containsTooltip",
      ($UnitedStates_selected ? "flagClicked" : "") + " " + (!$UnitedStates_selected ? "flagUnclicked" : "")
    ].join(" ").trim()}"><img src="usa_flag.svg" alt=""></div></div>  <div class="w-[40px]" data-svelte-h="svelte-1qoax4b"> </div>  <div class="${[
      "w-14 hover:bg-gray-300 hover:rounded-full flex justify-center items-center cursor-pointer h-full bg-no-repeat bg-center containsTooltip",
      ($apiCallsVisible ? "bg-icon-show-api-calls" : "") + " " + (!$apiCallsVisible ? "bg-icon-hide-api-calls" : "")
    ].join(" ").trim()}" aria-hidden="true"${add_attribute("_originaltooltip_", $Terms.api_fetch_visible, 0)} id="btn_turn_sound"></div>  <div class="${[
      "w-14 hover:bg-gray-300 hover:rounded-full flex justify-center items-center cursor-pointer h-full bg-no-repeat bg-center containsTooltip",
      ($soundEnabled ? "bg-icon-active-sound" : "") + " " + (!$soundEnabled ? "bg-icon-inactive-sound" : "")
    ].join(" ").trim()}" aria-hidden="true"${add_attribute("_originaltooltip_", $Terms.turn_sound, 0)} id="btn_turn_sound"></div>  <div class="${[
      "w-14 hover:bg-gray-300 hover:rounded-full flex justify-center items-center cursor-pointer h-full bg-no-repeat bg-center containsTooltip",
      ($recordIdsVisible ? "bg-icon-show-record-ids" : "") + " " + (!$recordIdsVisible ? "bg-icon-hide-record-ids" : "")
    ].join(" ").trim()}" aria-hidden="true"${add_attribute("_originaltooltip_", $Terms.record_ids_visible, 0)} id="btn_turn_sound"></div>  <div class="w-40 text-[16px] hover:bg-gray-300 border-gray-600 border-[1px] rounded-xl hover:rounded-xl flex justify-left pl-3 items-center cursor-pointer h-12 bg-no-repeat bg-[right_20px_center] containsTooltip bg-icon-current-group" aria-hidden="true"${add_attribute("_originaltooltip_", $Terms.current_workgroup, 0)} id="btnWorkgroupMenu" style="color:blue"></div>  <div class="w-14 hover:bg-gray-300 hover:rounded-full flex containsTooltip justify-center items-center cursor-pointer h-full bg-icon-bell bg-no-repeat bg-center " aria-hidden="true"${add_attribute("_originaltooltip_", $Terms.notifications, 0)} id="btnNotifications"><span id="newNotificationsAmount" class="text-sm relative -top-2 -right-2 bg-red-600 text-white rounded-full px-1 w-6 text-center invisible" data-svelte-h="svelte-1ut2onc"> </span></div></div></div>  <div class="w-full overflow-y-auto overflow-x-hidden mt-3 h-full " id="divMAIN"><div class="w-full pl-6 pr-4 flex flex-col items-center h-full "> ${$reportedErrorMessage != "" ? `<div id="errorBackDrop" class="w-full h-full absolute flex items-center justify-center left-0 top-0 z-50 bg-[rgba(0,0,0,0.5)] ">${validate_component(ReportedError, "ReportedError").$$render($$result, {}, {}, {})}</div>` : ``}  ${`<div id="backDrop" class="${[
      "w-screen h-screen absolute flex items-center justify-center left-0 top-0 z-10 aria-hidden=&quot;true&quot;",
      !$workgroupReady ? "bg-[rgba(0,0,0,0.2)]" : ""
    ].join(" ").trim()}">${!$isLanguageChosen ? `${validate_component(LanguageChoice, "LanguageChoice").$$render($$result, {}, {}, {})}` : `${$workgroupReady ? ` <div class="divLOADING_MSG">${escape($Terms.loading)}</div>` : ``}  ${!$workgroupReady ? `${validate_component(WorkgroupCloningModal, "WorkgroupCloningModal").$$render($$result, {}, {}, {})}` : ``}`}</div>`}</div></div></div>  <div id="slidingWindowMessage" data-svelte-h="svelte-g9lyjz"> </div></div>  <div id="apiDisplay"${add_classes(((!$apiCallsVisible ? "invisible" : "") + " " + ($apiCallsVisible ? "visible" : "")).trim())} data-svelte-h="svelte-1a5bv32"></div>  <div class="_doggy" id="divDoggy"></div> <div class="_doggy_1" id="divDoggy_1"></div> <div class="_doggy_2" id="divDoggy_2"></div> ${$UnitedStates_selected ? `<div class="_doggy_3_english" id="divDoggy_3"></div>` : `<div class="_doggy_3_portuguese" id="divDoggy_3"></div>`}`;
  } while (!$$settled);
  $$unsubscribe_UnitedStates_selected();
  $$unsubscribe_reportedErrorMessage();
  $$unsubscribe_Terms();
  $$unsubscribe_clientIp();
  $$unsubscribe_backendUrlGolang();
  $$unsubscribe_isLanguageChosen();
  $$unsubscribe_clickedSidebarMenu();
  $$unsubscribe_apiCallsVisible();
  $$unsubscribe_soundEnabled();
  $$unsubscribe_recordIdsVisible();
  $$unsubscribe_backendUrl();
  $$unsubscribe_workgroupReady();
  return $$rendered;
});

export { Page as default };
//# sourceMappingURL=_page.svelte-BnyvIl9O.js.map
