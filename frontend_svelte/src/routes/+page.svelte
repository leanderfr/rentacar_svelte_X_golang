<style>
@import '$public/tailwind.css';

</style>

<script>
// home= quando menu 'principal' esta acionado
import Home from '$lib/components/Home.svelte';
// Datatable= usado para manipular dados de todas as tabelas da aplicacao
import Datatable from '$lib/components/Datatable.svelte';
import Booking from '$lib/components/Bookings.svelte';
import WorkgroupMenu from '$lib/components/WorkgroupMenu.svelte';
import WorkgroupLogin from '$lib/components/WorkgroupLogin.svelte';
import WorkgroupDataReset from '$lib/components/WorkgroupDataReset.svelte';
import WorkgroupCloningModal from '$lib/components/WorkgroupCloningModal.svelte';
import WorkgroupRandomChooseAnother from '$lib/components/WorkgroupRandomChooseAnother.svelte';
import LanguageChoice from '$lib/components/LanguageChoice.svelte';
import BackendMenu from '$lib/components/BackendMenu.svelte';


import ReportedError from '$lib/components/ReportedError.svelte';

import jQuery from 'jquery'

import { Spinner } from 'spin.js';
import { simplifyLanguageData, showOrHideTooltip, getWorkgroupRecordsReady, refreshTerms, slidingMessage, showCloningStatus, setBackendUrl, 
      selectedCountry, decodeHTMLEntities, logAPI, getWorkgroupFilesReady } from '$js/utils.js'
import { recordIdsVisible, Terms, UnitedStates_selected, clickedSidebarMenu, backendUrl, backendUrlGolang, backendUrlPHP, clientIp,
      soundEnabled, reportedErrorMessage, workgroupReady, isLanguageChosen, apiCallsVisible, cloningCurrentBackend   } from '$lib/stores/stores.js'
import { invalidateAll } from '$app/navigation'  

export let data


// controla se o idioma ja foi escolhido e se ja ha grupo de trabalho definido/dados preparados
let pageDotSvelte_Loaded = false

// seletor de país é na verdade um checkbox, o lado direito, que equivale a 'checked', é o lado dos Estados Unidos
// ou seja, qdo 'checked' = Estados Unidos,   qdo nao 'checked'= Brasil

// 'simplifyLanguageData' converte o array de objetos JSON recebido ao carregar +page.js, em um array associativo 
// com isso, ao mudar de idioma ou carregar a tela pela 1a vez, automaticamente
// cada termo sera exibido no idioma atual
// exemplo: Terms.welcome, Terms.available_cars, etc etc 

// atencao:  '$' aciona subscribe/unsubscribe da variable 'writeable' sem necessidade de criar os metodos
$: $Terms = simplifyLanguageData( data.languageData )

 

// se mudou idioma/país, atualiza todos os textos de tooltip (title)
$: $Terms, setTimeout(() => {  if (pageDotSvelte_Loaded) showOrHideTooltip() }, 1000);  

import { afterUpdate, onMount } from 'svelte';  

onMount(() => {
  // define que para invocar jquery, pode ser usado jquery() e jq() (dentro de componentes svelte) e $() dentro de arquivos jscript
  // o caracter '$' é de uso reservado dentro do svelte, por isso necessario dentro de componente .svelte usar jquery() ou jq()
  window.jQuery = window.jq = window.$ = jQuery;

  // verifica se ja foi escolhido idioma/país, caso nao tenha, mais abaixo sera aberta janela para escolha do idioma
  $isLanguageChosen = localStorage.getItem("rentacar_current_language") == null ? false : true

  // lê qual backend escolhido atualmente e joga para variavel writable 'backendUrl'
  setBackendUrl()   


  /* 
  jQuery-UI é usada em varios momentos: tooltip, mensagens rapidas (.show()), etc etc
  foram feitas varias tentativas de carrega lo usando 'import', mas é impossivel, pois o jQuery precisa ja ter sido carregado

  a unica maneira que funciona é carregando 'jQuery-UI' dentro de onMount (que equivaleria a jscript document.ready())
  */
  jq.getScript('js/jquery-ui.min.js').done(function()  { 

    setTimeout(() => { 
      // melhora o visual do 'title' dos botoes  , tooltip pertence à 'jquery-ui'
      showOrHideTooltip() 

      // necessario cancelar autocomplete padrao do jquery-ui, para poder carregar plugin jquery melhorado 
      jq.fn.jquery_default_autocomplete = jq.fn.autocomplete;
      delete jq.fn.autocomplete;

      // arquivos js abaixo precisam estar na pasta /static/js, caso contrario o node/sveltekit nao 
      // transfere para a pasta /build ao transpilar
      // autocomplete do jquery-ui um melhorado
      jq.getScript( "js/jquery.autocomplete.js", function() {}) 

      // auxilia digitacao valores decimais (obrigado a Diego Plentz !!)
      jq.getScript( "js/maskMONEY.js", function() {}) 

      // permite arrasto de divs diferentes mas consideradas de um mesmo grupo (obrigado a Sudheer Someshwara !!)
      jq.getScript( "js/multiDraggable.js", function() {}) 

      // auxilia digitacao datas e horas  (obrigado a Josh Bush !!)
      jq.getScript( "js/maskDateHour.js", function() {}) 


      // antes de sair da apliacao web, desliga verificacao de novas notificacoes
      window.onbeforeunload = function () {
          clearInterval(newNotificationsChecking) 
          newNotificationsChecking = null
          console.log("newNotificationsChecking killed")
      }

      // melhora title da div que informa o backend/da acesso a mudanca de backend
      jq('#divBackendChoice').tooltip({ 
        tooltipClass: 'prettierTitle_white',
        show: false,   // sem animacao ao exibir 
        hide: false,   // sem animacao ao ocultar
        position: { my: "left top", at: "left+20 top-40", collision: "flipfit" }
      })

      // idioma/país ja foi escolhido, verifica o grupo entao
      if ($isLanguageChosen)   isWorkgroupReady()

    }, 
    100);
  });


  // propriedades da animacao que sera exibida dentro de 'divLoading', sempre que houver uma requisicao ao server-side
  var opts = {
    lines: 12 // The number of lines to draw
  , length: 40 // The length of each line
  , width: 18 // The line thickness
  , radius: 42 // The radius of the inner circle
  , scale: 0.3 // Scales overall size of the spinner
  , corners: 3 // Corner roundness (0..1)
  , color: 'gray' // #rgb or #rrggbb or array of colors
  , opacity: 0.3 // Opacity of the lines
  , rotate: 0 // The rotation offset
  , direction: 1 // 1: clockwise, -1: counterclockwise
  , speed: 1 // Rounds per second
  , trail: 60 // Afterglow percentage
  , fps: 20 // Frames per second when using setTimeout() as a fallback for CSS
  , zIndex: 2e9 // The z-index (defaults to 2000000000)
  , className: 'spinner' // The CSS class to assign to the spinner
  , top: '50%' // Top position relative to parent
  , left: '50%' // Left position relative to parent
  , shadow: true // Whether to render a shadow
  , hwaccel: true // Whether to use hardware acceleration
  , position: 'absolute' // Element positioning
  }

  // para exibir/ocultar esta div, usar as funcoes: showLoadingGif()/hideLoadingGif()
  var target = document.getElementById('divLoading');
  var spinner = new Spinner(opts).spin(target);


  // prepara animacao que vai aparecer qdo usuario passar mouse sobre icone do cachorrinho, canto inferior direito
  jq('#divDoggy').mouseover(function (e) {
    jq('#divDoggy_1').show(); jq('#divDoggy_2').show(); jq('#divDoggy_3').show();
    jq('#divDoggy_1').animate({ bottom: '55px', right: '105px', zIndex: 3000 }, 200, function () {
    });

    jq('#divDoggy_2').animate({ bottom: '75px', right: '125px', zIndex: 3000 }, 200);
    jq('#divDoggy_3').animate({ bottom: '90px', right: '105px' }, 200, function () {
      jq(this).css('z-index', 2101);
    });
  });

  // usuario tirou  mouse do icone cachorro 
  jq('#divDoggy').mouseout(function (e) {
    jq('#divDoggy_1').hide(); jq('#divDoggy_2').hide(); jq('#divDoggy_3').hide();
    jq('#divDoggy_1').css('right', '80px'); jq('#divDoggy_1').css('bottom', '1px');
    jq('#divDoggy_2').css('right', '80px'); jq('#divDoggy_2').css('bottom', '1px');
    jq('#divDoggy_3').css('right', '80px'); jq('#divDoggy_3').css('bottom', '-150px');
  });


})

// botoes menu esquerdo
$: buttons_sidebar_menu = [
  { name: $Terms.itemmenu_main,
    id: 'main',
    gray_icon: 'menu_item_home_gray.svg',
    blue_icon: 'menu_item_home_blue.svg',
    division_admin: ''
    
  },


  { name: $Terms.itemmenu_booking, 
    id: 'booking',
    gray_icon: 'menu_item_bookings_gray.svg',
    blue_icon: 'menu_item_bookings_blue.svg',
    division_admin: ''
  },


  { name: $Terms.itemmenu_cars,
    id: 'cars',
    gray_icon: 'menu_item_cars_gray.svg',
    blue_icon: 'menu_item_cars_blue.svg',
    division_admin: 'yes',     // indica que este icone ja faz parte da area adminitrativa do menu

  },

  { name: $Terms.itemmenu_manufacturers, 
    id: 'manufacturers',
    gray_icon: 'menu_item_manufacturers_gray.svg',
    blue_icon: 'menu_item_manufacturers_blue.svg',
    division_admin: ''
  },
  { name: $Terms.itemmenu_languages, 
    id: 'languages',
    gray_icon: 'menu_item_terms_gray.svg',
    blue_icon: 'menu_item_terms_blue.svg',
    division_admin: '', 
  } ,

  { name: $Terms.itemmenu_workgroups, 
    id: 'workgroups',
    gray_icon: 'menu_item_workgroups_gray.svg',
    blue_icon: 'menu_item_workgroups_blue.svg',
    division_admin: '', 
  } ,


];




/************************************************************************************************************************************************************
 se a tela for redimensionada, atualiza datatable para que ocupe o espaço (largura) disponivel 
************************************************************************************************************************************************************/

const windowResize = (event) => {
  // se esta vizualizando uma datatable no momento, redimensiona colunas apos 2 segundos 
  if (typeof currentDataTable == 'object' && currentDataTable!=null)  
      setTimeout(() => {currentDataTable.forceDatatableRefresh()}, 2000);

}


/************************************************************************************************************************************************************
 apos eventual resize da tela, reposiciona (caso estejam visiveis): 
1. o menu de grupo de trabalho (caso esteja visivel)
2. a tela de notificacoes

==> para que fiquem exatamente abaixo dos seus respectivos botoes de acionamento
************************************************************************************************************************************************************/

afterUpdate(() => {

  // posiciona o menu de grupos logo abaixo do icone 'grupos' (canto superior direito)
  let btnWorkgroupMenuInfo = document.getElementById( 'btnWorkgroupMenu' ).getBoundingClientRect()     // posicao do botao que aciona menu de grupos
  // posiciona menu de grupos pouco abaixo do botao 'menu grupos'
  let divTopPosition = parseInt(btnWorkgroupMenuInfo.top, 10) + parseInt(btnWorkgroupMenuInfo.height, 10) + 10

  jq('#workgroupMenu').css("top", divTopPosition);
  jq('#workgroupMenu').css("right", 20);


  // posiciona a barra de notificacoes logo abaixo do icone 'notificacoes' (canto superior direito)
  let btnNotificationsInfo = document.getElementById( 'btnNotifications' ).getBoundingClientRect()     // posicao do botao 'notificacoes'

  // posiciona menu de grupos pouco abaixo do botao 'menu grupos'
  divTopPosition = parseInt(btnNotificationsInfo.top, 10) + parseInt(btnNotificationsInfo.height, 10) + 10

  jq('#notificationsContainer').css("top", divTopPosition);
  jq('#notificationsContainer').css("right", 20);
 
  
});

/************************************************************************************************************************************************************
************************************************************************************************************************************************************/


// controla se div com notificacoes esta visivel
let calledNotificationBar = false;

// essa variavel faz referencia à datatable atualmente carregada, se houver uma
let currentDataTable

// controla se div com menu de grupos esta visivel
let calledWorkgroupMenu = false;

// controla se div de login em grupo de trabalho deve ser exibida
let calledWorkgroupLogin = false

// controla se div de reset dos dados do grupo deve ser exibida
let calledWorkgroupReset = false
// se userConfirmedWorkgroupDataReset = true, chama backend (resetOrChooseAnotherWorkgroup) para reset dos dados do grupo
let userConfirmedWorkgroupDataReset = false

// controla se div de escolha aleatoria de outro nome de grupo deve ser exibida
let calledWorkgroupRandomChooseAnother = false
// se userConfirmedRandomChooseAnotherGroup = true, chama backend (resetOrChooseAnotherWorkgroup) para sorteio de outro nome de grupo
let userConfirmedRandomChooseAnotherGroup = false

// controla se o menu de backend esta visivel
let calledBackendMenu = false


// conexao para busca de novas notificacoes, sera executada a cada 10 segundos
// nao sera usado SSE (server side events), porque SSE espera respossa em formato texto
// e novas notificacoes devem ser recebidas em formato JSON
// e tambem porque SSE é executada a cada 3 segundos, nao ha como aumentar esse tempo infelizmente
let newNotificationsChecking = null

let notificationItens = ''


/********************************************************************************************************************************************************************************
 fetch api para busca de novas notificacoes (notifica edicao de dados feita por outras pessoas dentro do mesmo grupo => insercao, edicao, exclusao de carros, fabricantes, etc)
 parte 1= montagem da chamada a cada 10 segundos
********************************************************************************************************************************************************************************/
const mountNewNotificationsConnection = () => {

// se por acaso ja foi montada a conexao, nao monta uma 2a vez
if (newNotificationsChecking != null) return;

// é atraves dessa conexao que o navegador busca, com base no workgroup atual (localStorage.rentacar_workgroup_name)
// se ha novas notificacoes a serem exibidas no icone de notificacoes (canto superior direito) 
let lastNotificationIdRead = localStorage.getItem("rentacar_last_notification_id_read");
if (lastNotificationIdRead == null)   localStorage.setItem("rentacar_last_notification_id_read", 0);


// solicita a cada 10 segundos verificacao de novas notificacoes
newNotificationsChecking = setInterval(async function () {
  checkNewNotifications()

}, 10000);


}

/********************************************************************************************************************************************************************************
 fetch api para busca de novas notificacoes (notifica edicao de dados feita por outras pessoas dentro do mesmo grupo => insercao, edicao, exclusao de carros, fabricantes, etc)
 parte 2= chamada fetch api
********************************************************************************************************************************************************************************/

async function checkNewNotifications() {

let _currentWorkgroupName = localStorage.getItem("rentacar_workgroup_name");

if (_currentWorkgroupName == null) return  

// obtem valor de 'backendUrl', em arquivos .JS nao se pode obter valor de variabel store usando somente '$'
// a verificacao de notificacoes  é  tarefa que so o backend golang faz
// os demais backends nao fazem essa tarefa

// exemplo URL:  http://localhost/Admin/notifications,   http://localhost/Sincere/notifications, etc

let route = `${$backendUrlGolang}/` + localStorage.getItem("rentacar_workgroup_name") + `/notifications/${$clientIp}`


await fetch(route, {method: 'GET'})

.then((response) => {
  if (!response.ok) {
    throw new Error(`Read Record Metadata Err fatal= ${response.status}`);
  }
  return response.json();
})

.then((notifications) => {

  jq('#notificationItens').html('')    // div com as notificacoes
  jq('#newNotificationsAmount').css('visibility', 'hidden')    // div com a qtde novas notificacoes

  let divItemDefault = `<div class='notificationItem'>@item</div>`

  // se detectadas notificacoes, popula a div 'notificationsItens' com as novas notificacoes

  let newNotificationsAmount = 0
  notificationItens = ''

  let lastNotificationIdDisplayed = parseInt(localStorage.getItem("rentacar_last_notification_id_read"), 10)
  let notificationId = 0


  setTimeout(() => {

        Object.keys(notifications).forEach(function(key)   {

          let itemAdd = divItemDefault

          // notifications[key].description_english ou notifications[key].description_portuguese contem a notificacao que foi gerada no momento em qua houve a alteracao na base
          // exemplo: carro X foi alterado, reserva Y foi feita, etc etc

          notificationId = parseInt(notifications[key].id, 10)

          // se a notificacao sendo exibida é posterior à ultima notificacao lida, coloca em negrito, sinalizando que ela é nova, contabiliza a notificacao como nova (newNotificationsAmount++)
          if (notificationId > parseInt(localStorage.getItem("rentacar_last_notification_id_read"), 10))  {
            notifications[key].description_english = notifications[key].description_english.replaceAll('@color', 'font-weight:bold;color:black')
            notifications[key].description_portuguese = notifications[key].description_portuguese.replaceAll('@color', 'font-weight:bold;color:black')

            newNotificationsAmount++ // so considera como nova notificacao se o ID for posterior à ultima notificacao obtida 

            // contabiliza exatamente qual a ultima notificacao lida (lastNotificationIdDisplayed)
            if (notificationId > lastNotificationIdDisplayed)  
              lastNotificationIdDisplayed = notificationId   // memoriza qual a ultima notificacao ja exibida , para mais adiante remover o negrito das novas notificacoes

          // notificacao é anterior (id) à ultima notificacao visualizada, coloca em cinza, sinalizando que ja foi lida
          }  else {
            notifications[key].description_english = notifications[key].description_english.replaceAll('@color', 'color:gray')
            notifications[key].description_portuguese = notifications[key].description_portuguese.replaceAll('@color', 'color:gray')

          }

          // mostra a notificacao no idioma escolhido no momento
          if (selectedCountry()=='usa') 
            itemAdd = itemAdd.replaceAll('@item', decodeHTMLEntities(notifications[key].description_english) ) 

          if (selectedCountry()=='brazil') 
            itemAdd = itemAdd.replaceAll('@item', decodeHTMLEntities(notifications[key].description_portuguese) )  

          notificationItens += itemAdd
        })

        // exibe qtde novas notificacoes se houver pelo menos uma
        if (newNotificationsAmount > 0) {
          jq('#newNotificationsAmount').html(newNotificationsAmount)
          jq('#newNotificationsAmount').css('visibility', 'visible')
        }

        // se daqui 5 segundos o usuario aainda estiver com a barra de notificacoes aberta, marca todas as notificacoes como lidas e remove o negrito delas
        setTimeout(() => {
          if (calledNotificationBar) {
            localStorage.setItem("rentacar_last_notification_id_read", lastNotificationIdDisplayed);

            // remove o negrito de todas as notificacoes, sinalizando que todas foram visualizadas em 3 segundos
            notificationItens = notificationItens.replaceAll(`font-weight:bold;color:black`, 'color:gray')

            // remove qtde de notificacoes novas (em vermelho, canto superior direito)
            jq('#newNotificationsAmount').html('')
            jq('#newNotificationsAmount').css('visibility', 'hidden')
          }
          
        }, 5000);
      })

    
  }, 100);


}

/************************************************************************************************************************************************************
 muda idioma se houve clique na bandeira diferente da atual
************************************************************************************************************************************************************/

const flagClicked = (event) =>  {
  if ( (event.currentTarget.id=='flagBRAZIL' && $UnitedStates_selected) || 
        (event.currentTarget.id=='flagUSA' && !$UnitedStates_selected)  )  {

      localStorage.setItem("rentacar_current_language", event.currentTarget.id=='flagUSA' ? 'usa' : 'brazil')

      jq('#chkLanguageSelector').trigger('click')
    }
}

/************************************************************************************************************************************************************
 usuario pediu para acessar dados de outro grupo
************************************************************************************************************************************************************/
const enterAnotherWorkgroup = () => {

setTimeout(() => {calledWorkgroupLogin = false }, 500);
setTimeout(() => {calledWorkgroupMenu = false }, 1000);

setTimeout(() => {

    // atualiza expressoes do grupo acessado e recomeça e checar notificacoes
    jq.when(refreshTerms() ).done(function( ) {  
        let _currentWorkgroupName = localStorage.getItem("rentacar_workgroup_name");
        slidingMessage('slidingWindowMessage', $Terms.new_current_workgroup + '&nbsp;&nbsp;' + _currentWorkgroupName, 3000)    

        // fecha conexao de verificacao de novas notificacoes, do grupo atual, primeiramente
        if ( newNotificationsChecking != null) {
          clearInterval(newNotificationsChecking) 
          newNotificationsChecking = null
          console.log("newNotificationsChecking killed")
        }

        jq('#homeCarsBrowser').scrollLeft(0);  // reseta scroll browser carros Home.svelte (se estiver visivel)
        jq('#bookingCarsBrowser').scrollTop(0);  // reseta scroll browser carros Booking.svelte

        invalidateAll()   // atualiza dados na tela
        setTimeout(() => {mountNewNotificationsConnection()}, 1000);
    })        

}, 1500);
}




/************************************************************************************************************************************************************
 usuario pediu para:
1- resetar dados do grupo atual (userConfirmedWorkgroupDataReset)
2- sortear outro nome de grupo (userConfirmedRandomChooseAnotherGroup)
************************************************************************************************************************************************************/
const resetOrChooseAnotherWorkgroup = () =>  { 

// reseta os dados do grupo de trabalho atual

calledWorkgroupReset = false   // fecha janela modal de reset
calledWorkgroupMenu = false    // fecha janela modal de workgroup menu
calledWorkgroupRandomChooseAnother = false   // fecha janela de confirmacao de escolha de outro grupo

// copia os dados do grupo 'Admin' (dados default) para o grupo atual
let whatToDo = ''
let sucessMsg = ''
if ( userConfirmedWorkgroupDataReset )    {
  whatToDo = 'reset'
  sucessMsg = $Terms.workgroup_data_was_reset
}

if ( userConfirmedRandomChooseAnotherGroup )    {
  whatToDo = 'another'
  sucessMsg = $Terms.workgroup_randomly_chosen_success
}

setTimeout(() => {
    // fecha conexao de verificacao de novas notificacoes, do grupo atual,  primeiramente
    if ( newNotificationsChecking != null ) {
      clearInterval(newNotificationsChecking) 
      newNotificationsChecking = null
      console.log("newNotificationsChecking killed")
    }

    jq.when(getWorkgroupRecordsReady( whatToDo )).done(function()  {  

      // registros foram resetados e imagens clonadas com sucesso
      clearInterval(verifyCloningStatus)

      // fecha a janela de status da clonagem 
      userConfirmedWorkgroupDataReset = false
      userConfirmedRandomChooseAnotherGroup = false

      // avisa que clonou
      setTimeout(() => {
        // avisa que deu certo  e atualiza terms ing/portugues
        refreshTerms()
        slidingMessage('slidingWindowMessage', sucessMsg, 7000)      

        // reinicia verificacao de notificacoes
        setTimeout(() => { mountNewNotificationsConnection()  }, 1000);            
      }, 500);      

      jq('#homeCarsBrowser').scrollLeft(0);  // reseta scroll browser carros Home.svelte (se estiver visivel)
      jq('#bookingCarsBrowser').scrollTop(0);  // reseta scroll browser carros Booking.svelte
    })        
}, 200);


// solicita a cada 1/2 segundo qual a situacao da clonagem de registros/arquivos imagens para o grupo obtido acima
let verifyCloningStatus = setInterval(async function () {showCloningStatus()}, 500);
}



/************************************************************************************************************************************************************
 manipula teclas pressionadas por toda aplicacao
************************************************************************************************************************************************************/

const onKeyDown = (e) =>  {

  // se usuario pressionou Enter ou setas acima/abaixo, avanca o foco para o proximo/anterior campo
  if ( (e.which == 13 || e.which == 38 || e.which == 40)  && jq('.text_formFieldValue').is(':focus') )   { 
        // tive que criar a propriedade 'sequence', porque o vite js fica enchendo o saco sobre nao usar tabIndex

        let tab =  jq(':focus').attr("sequence");
        if (e.which==13 || e.which == 40)  tab++;
        else if (e.which==38)  tab--;

        e.preventDefault()

        // tive que criar a propriedade 'sequence', porque o vite js fica enchendo o saco sobre nao usar tabIndex
        jq("[sequence='"+tab+"']").focus();          
  }

  // tela de login (grupo) esta visivel 
  if (typeof jq('#workgroupLogin').attr("id") != 'undefined')  {
    if (e.which == 13) {
      jq('#btnLOGIN').trigger('click')
      return;
    }

    if (e.which == 27) {
      jq('#btnCLOSE_LOGIN_WINDOW').trigger('click')
      return;
    }

  }

  // tela de reset dos dados do grupo  esta visivel 
  // tela de sorteio de novo grupo  esta visivel 
  if (typeof jq('#workgroupDataReset').attr("id") != 'undefined')  {

    if (e.which == 27) {
      jq('#btnCLOSE_RESET_WINDOW').trigger('click')
      return;
    }
  }

  // tela de sorteio de novo grupo  esta visivel 
  if (typeof jq('#workgroupRandomChooseAnother').attr("id") != 'undefined')  {

    if (e.which == 27) {
      jq('#btnCLOSE_RANDOM_CHOOSE_WINDOW').trigger('click')
      return;
    }

  }






  // se usuario pressionou Esc  ou F2 e ha  um formulario de edicao de registro ou uma janela modal aberta
  if (e.which == 27 || e.which == 113)   { 
        let metadata_window = typeof jq('#metadataWindow').attr("id");
        let record_form = typeof jq('#recordForm').attr("id");
        let confirmation_window = typeof jq('#confirmationWindow').attr("id");
        let calendar_window = typeof jq('#calendarModalWindow').attr("id");
        let booking_form = typeof jq('#bookingRecordForm').attr("id");
  
        // metadata_window é prioridade no fechamento porque sera sempre o ultimo da pilha, sempre estar por cima de outros forms/janelas
        // a janela modal '#metadataWindow' possui botao fechar com nome diferente, para evitar conflito com botao fechar das janelas modais:
        // '#recordForm' e '#confirmationWindow'
        if (metadata_window != 'undefined' && e.which == 27 )  {
          jq('#btnCLOSE_METADATA_WINDOW').trigger('click')
          return;
        }

        // record_form, confirmation_window, calendar_window, booking_form, estao no mesmo nivel de vizualizacao, uma nunca vai estar por cima da outra
        // por isso elas podem usar o mesmo ID para botao 'fechar' (btnCLOSE)
        if (record_form != 'undefined' || confirmation_window != 'undefined' || calendar_window != 'undefined' || booking_form != 'undefined')  {
          if (e.which == 27)   jq('#btnCLOSE').trigger('click')

          // os forms abaixo possuam botao F2= salvar (F2= 113)
          if (record_form != 'undefined' || booking_form != 'undefined')  {
            if (e.which == 113)   jq('#btnSAVE').trigger('click')   // f2 foi pressionado
          }
        }
  }

  // se pressionar Enter estando com foco no input.textbox de pesquisa (canto superior direito de uma datatable)
  if (e.which == 13 && jq('.txtTABLE_SEARCHBOX').is(':focus') )  {
      currentDataTable.performDatatableSearch()
  }
}


// manipula erros genericos javascript
const handleJsGeneralError = (e) => {

// desliga verificacao de notificacoes
if (typeof newNotificationsChecking!='undefined') {
  clearInterval(newNotificationsChecking) 
  console.log("newNotificationsChecking killed")
}

$reportedErrorMessage = e.error.message
}

// manipula erros de PROMISE, gerados ao fazer api fetch no back end
// maioria das chamadas api fetch tem tratamento de erro especifico acoplado, algumas nao
const handleJsAPIFetchError = (e) => {

// desliga verificacao de notificacoes
if (typeof newNotificationsChecking!='undefined') {
  clearInterval(newNotificationsChecking) 
  console.log("newNotificationsChecking killed")
}

$reportedErrorMessage = e.reason.message
}

/************************************************************************************************************************************************************
************************************************************************************************************************************************************/
const currentWorkgroup = () => {
let _currentWorkgroupName = localStorage.getItem("rentacar_workgroup_name");
return _currentWorkgroupName.toLowerCase()
}

/************************************************************************************************************************************************************
// verifica se ha grupo de trabalho definido , localStorage.getItem("rentacar_workgroup_name")
//    se sim, exibe grupo ja conhecido
//    se nao, obtem grupo novo  
************************************************************************************************************************************************************/
const isWorkgroupReady = () => {

let _currentLanguage = localStorage.getItem("rentacar_current_language")
$UnitedStates_selected = _currentLanguage == 'usa' ? true : false

refreshTerms();

setTimeout(() => {

    // fecha eventual conexao de verificacao de novas notificacoes, caso exista
    if ( newNotificationsChecking != null ) {
      clearInterval(newNotificationsChecking) 
      newNotificationsChecking = null
      console.log("newNotificationsChecking killed")
    }

    // chama backend para criacao de grupo e preparacao (clonagem) dos dados
    // se o grupo ja esta preparado, getWorkgroupRecordsReady retorna sem fazer nada

    jq.when(getWorkgroupRecordsReady('generate')).done(function()  {  

      // clonagem dos registros e arquivos de imagem foi ok
      clearInterval(verifyCloningStatus)
      pageDotSvelte_Loaded=true   // avisa que pode continuar a preparacao da pagina
      setTimeout(() => {mountNewNotificationsConnection()  }, 1000);      
      setTimeout(() => {
        checkNewNotifications()    // faz uma 1a verificacao de novas notificacoes
      }, 1500); 

      jq('#homeCarsBrowser').scrollLeft(0);  // reseta scroll browser carros Home.svelte (se estiver visivel)
      jq('#bookingCarsBrowser').scrollTop(0);  // reseta scroll browser carros Booking.svelte
    })        


    // solicita a cada 1/2 segundo qual a situacao da clonagem de registros/arquivos imagens para o grupo obtido acima
    let verifyCloningStatus = setInterval(async function () {showCloningStatus()}, 500);
 

}, 500);
}

/************************************************************************************************************************************************************
// atualiza localStorage com o idioma atual, porque usuario alterou idioma clicando diretamente no checkbox brasil/usa
// $UnitedStates_selected reflete o estado do checkbox de alteracao de país/idioma
************************************************************************************************************************************************************/

const updateCountry = () => {
let _currentLanguage = $UnitedStates_selected ? 'usa' : 'brazil'

localStorage.setItem('rentacar_current_language', _currentLanguage )   

jq('#homeCarsBrowser').scrollLeft(0);  // reseta scroll browser carros Home.svelte (se estiver visivel)
jq('#bookingCarsBrowser').scrollTop(0);  // reseta scroll browser carros Booking.svelte

// força atualizacao da tela
setTimeout(() => {
  invalidateAll()  ; 
}, 100);

}

</script>


<!--
*************************************************************************************************************************************
*************************************************************************************************************************************
*************************************************************************************************************************************

HTML 

*************************************************************************************************************************************
*************************************************************************************************************************************
*************************************************************************************************************************************
-->

<!-- funcoes associadas a janela da aplicacao  -->
<svelte:window  on:keydown={onKeyDown} on:resize={windowResize}  on:error={handleJsGeneralError} on:unhandledrejection={handleJsAPIFetchError} />


<svelte:head>
</svelte:head>


<!-- 
usuario pediu pra mostrar menu de backend's possiveis 
se usuario alterar mesmo o backend, nao é necessario fazer nada aqui, pois o componente 'BackendMenu'
ja alterou => localStorage.getItem("rentacar_current_backend");
as chamadas API a partir de agora todas refletirarao o novo backend 
-->
{#if calledBackendMenu }  

  <div id='divBackendMenuBackdrop' class='w-full h-full  absolute flex items-center justify-center left-0 top-0 z-50 bg-[rgba(0,0,0,0.5)] cursor-pointer'  aria-hidden="true"
    on:click|self={ () => {console.log('clicou');calledBackendMenu = false} }   >  
    <BackendMenu  on:dispatchProceedBackendChange={ () => {calledBackendMenu = false} }        />
   </div> 

{/if}



<!-- menu de grupos foi acionado -->
{#if calledWorkgroupMenu}

  <div id='divWorkgroupMenuBackdrop' class='w-full h-full  absolute flex  left-0 top-0 z-40 bg-[rgba(0,0,0,0.1)] cursor-pointer'  
    on:click|self={() => { calledWorkgroupMenu = false }} aria-hidden="true"  >  
    <WorkgroupMenu  
        on:dispatchCallMenuLogin={ () => {calledWorkgroupLogin=true} }    
        on:dispatchResetWorkgroupData={ () => {calledWorkgroupReset=true} }  
        on:dispatchRandomChooseAnother={ () => {calledWorkgroupRandomChooseAnother=true} }      />
  </div> 

{/if}

<!-- login de grupo foi acionado -->
{#if calledWorkgroupLogin}

  <div id='divWorkgroupLoginBackdrop' class='w-full h-full  absolute flex items-center justify-center left-0 top-0 z-50 bg-[rgba(0,0,0,0.5)] cursor-pointer'  
    on:click|self={() => { calledWorkgroupLogin = false }} aria-hidden="true"  >  
    <WorkgroupLogin   
        on:dispatchCloseWorkgroupLogin={() =>  {calledWorkgroupLogin=false}}  
        on:dispatchAcessAnotherGroupData={() =>  {enterAnotherWorkgroup()}}  /> 
  </div> 


{/if}

<!-- tela de confirmacao de reset de grupo foi acionada -->
{#if calledWorkgroupReset}

  <div id='divWorkgroupResetBackdrop' class='w-full h-full  absolute flex items-center justify-center left-0 top-0 z-50 bg-[rgba(0,0,0,0.5)] cursor-pointer'  
    on:click|self={() => { calledWorkgroupReset = false }} aria-hidden="true"  >  
    <WorkgroupDataReset  
        on:dispatchCloseWorkgroupReset={() =>  {calledWorkgroupReset=false}}  
        on:dispatchPerformGroupDataReset={() =>  {
          userConfirmedWorkgroupDataReset = true  // avisa que usuario pediu reset
          userConfirmedRandomChooseAnotherGroup = false;   // avisa que usuario nao pediu sorteio de outro nome
          resetOrChooseAnotherWorkgroup() } 
        }       /> 
  </div> 

{/if}

<!-- tela de confirmacao de sorteio de outro nome de grupo -->
{#if calledWorkgroupRandomChooseAnother}

  <div id='divWorkgroupRandomChooseAnotherBackdrop' class='w-full h-full  absolute flex items-center justify-center left-0 top-0 z-50 bg-[rgba(0,0,0,0.5)] cursor-pointer'  
    on:click|self={() => { calledWorkgroupRandomChooseAnother = false }} aria-hidden="true"  >
    <WorkgroupRandomChooseAnother 
        on:dispatchCloseWorkgroupAnother={ () => {calledWorkgroupRandomChooseAnother=false} }  
        on:dispatchPerformWorkgroupRandomChooseAnother={() =>  {
          userConfirmedWorkgroupDataReset = false  // avisa que usuario nao pediu reset
          userConfirmedRandomChooseAnotherGroup = true;   // avisa que usuario pediu sorteio de outro nome
          resetOrChooseAnotherWorkgroup() } 
        }    />
  </div> 

{/if}


<!-- reset de dados do grupo ou sorteio aleatorio de outro nome de grupo, foi confirmado 
exibe tela de status da clonagem, pois o reset/clonagem sera acionado  -->
{#if userConfirmedWorkgroupDataReset || userConfirmedRandomChooseAnotherGroup }

  <div id='backDrop' class='w-screen h-screen  absolute flex items-center justify-center left-0 top-0 z-50 aria-hidden="true" bg-[rgba(0,0,0,0.5)]'  >  
      <WorkgroupCloningModal {userConfirmedWorkgroupDataReset}    />
  </div>  

{/if} 


<!-- barra com notificacoes foi acionada -->
{#if calledNotificationBar}

  <div id='notificationsBarBackdrop' class='w-full h-full  absolute flex left-0 top-0 z-40 bg-[rgba(0,0,0,0.1)] cursor-pointer'  
    on:click|self={() => { calledNotificationBar = false }} aria-hidden="true"  >  

        <div class="flex flex-col h-[calc(100%-90px)]  w-[750px]  bg-white border-gray-200 border-2 rounded-lg z-50 absolute"  
          style="scrollbar-color:gray white;; scrollbar-width: thin;"  id='notificationsContainer'  >  

              <div class="flex h-12 w-full  p-2 border-b-2 border-blue-400" >
                  {$Terms.notifications}
              </div>

              <div class="overflow-y-scroll  h-full" id='notificationItens' >
                {@html notificationItens}
              </div>
        </div>

  </div>

{/if}



<!-- som que sera tocado qdo algum erro exibido -->
<audio id="errorBeep" src="error_beep.mp3" preload="auto" autobuffer></audio>

<!-- contem a animacao que é carregada ao processar algo que demore -->
<div id="divLoading" class="cssLOADING_HIDE"></div>  

<div class="flex h-screen  select-none font-Roboto text-sm flex-row " >


    <!-- ************************************************************************** -->
    <!-- menu lateral esquerdo [#EDF2F9] -->
    <!-- ************************************************************************** -->
    <div class="bg-gray-100 w-[15%] h-full  "> 

        <!-- logotipo -->
        <div class="flex items-center cursor-pointer hover:bg-gray-200 border-b-2 border-gray-300 h-[105px] border-b-gray-400 "   aria-hidden="true"  >

            <div class="w-full " >
              <div class="w-full flex flex-col   align-middle py-2 h-16 border-b-2 border-gray-300" >
                  <div class="text-lg pl-3 "><span class="text-red-700 ">RENT</span> A CAR</div>
                  <div class="block  pl-3 "   aria-hidden="true"  >WebApp Demonstration</div>            
              </div>

              <!-- seletor de backend     -->
              <div class='backendSelector' title={$Terms.change_backend} id='divBackendChoice' 
                  on:click={ () => { calledBackendMenu = true } } aria-hidden="true">  
                <div style='font-size:14px;margin-top:-5px'>{$Terms.current_backend}</div>
                <div style='padding-left:10px;width:50px;height:50px;'><img  style="margin-top:0px" alt='' id='backendIcon'> </div>
              </div>
          
            </div>  
        </div>

        <!-- botoes menu lateral, parte superior -->
        <div class="pt-6"> 

          {#each buttons_sidebar_menu as {id, name, gray_icon, blue_icon, division_admin}}

              <!-- se 'division_admin'= 'yes' coloca um traço divisor entre o menu da area cliente e da area administracao -->
              {#if division_admin == 'yes'}  

                <div class="pl-5  border-b-2 border-b-gray-300 mt-10 mb-2  pb-2 mx-auto w-11/12 flex">{$Terms.itemmenu_administration}</div>

              {/if}

              <!-- opcao menu atual esta selecionada -->
              {#if id == $clickedSidebarMenu}

                <div class="btn_sidebar_selected"  aria-hidden="true" > 
                    <div class="ml-4 -mt-1 ">
                      <img src={blue_icon} alt=''>
                    </div>
                    <span class="text-blue-700">{name}</span>
                  </div>

              <!-- opcao menu atual nao esta selecionada -->
              {:else}
                <div class="btn_sidebar" on:click={() => { $clickedSidebarMenu = id }} aria-hidden="true" > 
                    <div class="ml-4 -mt-1">
                      <img src={gray_icon} alt=''>
                    </div>
                    <span class="text-gray-600">{name}</span> 
                </div>
              {/if}

          {/each}

        </div>
    </div>

    <!-- ******************************************************************************************** -->
    <!-- lado maior e direito da tela, onde o conteudo é atualizado conforme botoes menu sao clicados -->
    <!-- ******************************************************************************************** -->

    <div class="flex flex-col  w-[85%] pr-1  h-full grow overflow-hidden  " > 

        <!-- 1a linha: bem vindo e botao notificacoes, icone login/logout -->
        <div class="w-full flex flex-row items-center static px-6  h-[93px] ">    

            <!-- lado esquerdo -->
            <div class="grow text-xl font-bold">
              { $Terms.welcome }
            </div>

            <!-- lado direito -->
            <div class=" flex justify-end  items-center h-14  gap-4">

              <!-- seletor de idioma/país -->
              <div class="flex flex-row  gap-4">    
                <!-- bandeira Brasil -->
                <div class:flagClicked={! $UnitedStates_selected} class:flagUnclicked={$UnitedStates_selected} aria-hidden="true" id='flagBRAZIL' on:click={flagClicked}  
                  _originalTooltip_={$Terms.change_country} class="containsTooltip" >     
                  <img src='brazil_flag.svg' alt=''>
                </div>

                <label for="chkLanguageSelector" class="switch_language containsTooltip" _originalTooltip_={$Terms.change_country} >
                  <input id="chkLanguageSelector" type="checkbox"  bind:checked={$UnitedStates_selected} on:change={ () => {updateCountry()}}  >
                  <span class="slider_language round"></span>
                </label>

                <!-- bandeira EUA -->
                <div class:flagClicked={$UnitedStates_selected} class:flagUnclicked={! $UnitedStates_selected}  aria-hidden="true" id='flagUSA'  on:click={flagClicked}  
                  _originalTooltip_={$Terms.change_country}  class="containsTooltip" >    
                  <img src='usa_flag.svg' alt=''>                  
                </div>
              </div>

              <!-- espaçamento -->
              <div class='w-[40px]'>&nbsp;</div>

              <!-- desliga / liga exibicao das chamadas API (parte inferior da tela) -->
              <div class="w-14 hover:bg-gray-300 hover:rounded-full flex justify-center items-center cursor-pointer h-full  bg-no-repeat bg-center containsTooltip" aria-hidden="true" 
                    _originalTooltip_={$Terms.api_fetch_visible} class:bg-icon-show-api-calls={$apiCallsVisible} class:bg-icon-hide-api-calls={! $apiCallsVisible}   
                    on:click={() => ($apiCallsVisible = ! $apiCallsVisible)} id="btn_turn_sound" >
              </div>


              <!-- desliga liga som de alerta ao exibir mensagens  -->
              <div class="w-14 hover:bg-gray-300 hover:rounded-full flex justify-center items-center cursor-pointer h-full  bg-no-repeat bg-center containsTooltip" aria-hidden="true" 
                    _originalTooltip_={$Terms.turn_sound} class:bg-icon-active-sound={$soundEnabled} class:bg-icon-inactive-sound={! $soundEnabled}  
                     on:click={() => ($soundEnabled = ! $soundEnabled)} id="btn_turn_sound" >
              </div>


              <!-- desliga / liga exibicao do ID dos registros -->
              <div class="w-14 hover:bg-gray-300 hover:rounded-full flex justify-center items-center cursor-pointer h-full  bg-no-repeat bg-center containsTooltip" aria-hidden="true" 
                    _originalTooltip_={$Terms.record_ids_visible} class:bg-icon-show-record-ids={$recordIdsVisible} class:bg-icon-hide-record-ids={! $recordIdsVisible}   
                    on:click={() => ($recordIdsVisible = ! $recordIdsVisible)} id="btn_turn_sound" >
              </div>


              <!-- botao grupo atual, dá acesso ao menu de grupos   -->
              <div class="w-40 text-[16px]  hover:bg-gray-300  border-gray-600 border-[1px] rounded-xl  hover:rounded-xl flex justify-left pl-3 items-center cursor-pointer h-12  
              bg-no-repeat bg-[right_20px_center] containsTooltip bg-icon-current-group" 
                    aria-hidden="true" _originalTooltip_={$Terms.current_workgroup}  on:click={() => (calledWorkgroupMenu = ! calledWorkgroupMenu)} id="btnWorkgroupMenu" style='color:blue' >                
              </div> 

             

              <!-- notifications button -->
              <div class="w-14 hover:bg-gray-300 hover:rounded-full flex containsTooltip justify-center items-center cursor-pointer h-full bg-icon-bell  bg-no-repeat bg-center " aria-hidden="true" 
                       on:click={() => (calledNotificationBar = !calledNotificationBar)} _originalTooltip_={$Terms.notifications} id='btnNotifications' >
                <span id="newNotificationsAmount" class="text-sm relative -top-2 -right-2 bg-red-600 text-white rounded-full px-1 w-6 text-center invisible">&nbsp;</span>
              </div>


            </div>

        </div>


        <!-- 2a linha, parte principal da aplicacao, conteudo dinamico que sera alterado tempo todo: home, veiculos, etc etc -->
        <div class="w-full overflow-y-auto  overflow-x-hidden mt-3  h-full " id='divMAIN' > 

          <div class="w-full pl-6 pr-4 flex flex-col items-center h-full " > 

              <!-- div usada para reportar eventuais erros jscript e erros de chamadas API's -->
              {#if $reportedErrorMessage!=''}      
                <div id='errorBackDrop' class='w-full h-full  absolute flex items-center justify-center left-0 top-0 z-50 bg-[rgba(0,0,0,0.5)] '  >  
                  <ReportedError  />
                </div>
              {/if}


              <!-- so carrega componentes (items menu) da pagina principal, se a mesma foi totalmente carregada 
              cada componente svelte abaixo possui seu 'onMount()' e  para nao dar erro, page.svelte tem que ter sido 100% carregada -->

              {#if pageDotSvelte_Loaded}
                {#if $clickedSidebarMenu == 'main'}      
                  <Home  />  
                {/if}

                {#if $clickedSidebarMenu == 'booking'}      
                  <Booking   />
                {/if}

                <!-- toda datatable possui a coluna de botoes (delete, edit, etc), com  150 pixels 
                existe tambem a coluna de selecao de registro (checkbox) que tem 50px.. somando tudo dá 200px 
                por isso, sempre pelo menos uma das colunas precisa ter a largura de calc(nn%-200px)  -->
                {#if $clickedSidebarMenu == 'cars'}      
                  <Datatable   
                            fields={ [ {field_name: 'name', field_title: $Terms.fieldname_generic_name,  field_width:'basis-[33%] grow-0 '}, 
                                       {field_name: 'manufacturer_logo', field_title: $Terms.fieldname_car_manufacturer, field_width:'basis-[33%] grow-0 '},
                                       {field_name: 'car_image', field_title: $Terms.fieldname_image, field_width:'basis-[calc(33%-200px)] grow-0'} ] }
                            apiURL= {`${$backendUrl}/@group_info/cars`}    orderBy={'name'}   orderDirection={'asc'}   metadataTableName={$Terms.table_cars} 
                            tableTitle = {$Terms.itemmenu_cars}  tableName = 'cars'  formUse = 'car'  toDoAfterPatchOrPost=''  fieldConcatenateWithId = 'name'
                            bind:this={currentDataTable}      />
                {/if}

                {#if $clickedSidebarMenu == 'manufacturers'}      
                  <Datatable   
                            fields={ [ {field_name: 'name', field_title: $Terms.fieldname_generic_name,  field_width:'basis-[calc(60%-200px)] grow-0 '}, 
                                       {field_name: 'manufacturer_logo', field_title: $Terms.fieldname_logo, field_width:'basis-[40%] grow-0'} ] }
                            apiURL= {`${$backendUrl}/@group_info/manufacturers`}    orderBy={'name'}   orderDirection={'asc'}     metadataTableName={$Terms.table_manufacturers}   
                            tableTitle = {$Terms.itemmenu_manufacturers}   tableName='manufacturers' formUse = 'manufacturer'  toDoAfterPatchOrPost=''    fieldConcatenateWithId = 'name'
                            bind:this={currentDataTable}     /> 
                {/if}



                {#if $clickedSidebarMenu == 'languages'}      
                  <Datatable   
                            fields={ [ {field_name: 'item', field_title: $Terms.fieldname_language_item,  field_width:'basis-[calc(50%-200px)] grow-0 '}, 
                                       {field_name: 'portuguese', field_title: $Terms.portuguese_title, field_width:'basis-[25%] grow-0'} ,
                                       {field_name: 'english', field_title: $Terms.english_title, field_width:'basis-[25%] grow-0'} ] }
                            apiURL= {`${$backendUrl}/@group_info/terms`}    orderBy={'item'}   orderDirection={'asc'}        metadataTableName={$Terms.table_languages}   
                            tableTitle = {$Terms.datatable_title_expressions}   tableName='terms' formUse='term'   toDoAfterPatchOrPost='refreshTerms()'  fieldConcatenateWithId = 'item' 
                            bind:this={currentDataTable}      /> 
                {/if}

                {#if $clickedSidebarMenu == 'workgroups'}      
                    <!-- grupo admin tem acesso a mais detalhes dos grupos --> 
                    {#if currentWorkgroup() == 'admin'}      
                        <Datatable   
                                  fields={ [ {field_name: 'name', field_title: $Terms.fieldname_generic_name,  field_width:'basis-[10%] '  }, 
                                            {field_name: 'in_use', field_title: $Terms.fieldname_in_use, field_width:'basis-[10%] ', is_boolean: 'true'},
                                            {field_name: 'client_ip', field_title: $Terms.fieldname_ip, field_width:'basis-[10%] '},
                                            {field_name: 'client_city', field_title: $Terms.fieldname_city, field_width:'basis-[15%] '},
                                            {field_name: 'client_country', field_title: $Terms.fieldname_country, field_width:'basis-[15%] '},
                                            {field_name: 'client_loc', field_title: $Terms.fieldname_location, field_width:'basis-[20%] '},
                                            {field_name: 'updated_at', field_title: $Terms.fieldname_updated_at, field_width:'basis-[10%] '},
                                            {field_name: 'deleted_at', field_title: $Terms.fieldname_reset_at, field_width:'basis-[10%] '} ] }
                                  apiURL= {`${$backendUrl}/admin/workgroups`}    orderBy={'name'}   orderDirection={'asc'}        metadataTableName={$Terms.table_workgroups}   
                                  tableTitle = {$Terms.datatable_title_workgroups}   tableName='workgroups' formUse='workgroup'   toDoAfterPatchOrPost=''   fieldConcatenateWithId = 'name'
                                  bind:this={currentDataTable}    recordsReadonly=true  /> 



                    <!-- outros grupos, obtidos por usuarios comuns, tem acesso so ao nome dos grupos disponiveis  --> 
                    {:else}
                        <Datatable   
                                  fields={ [ {field_name: 'name', field_title: $Terms.fieldname_generic_name,  field_width:'basis-[100%] grow-0 '} ]} , 
                                  apiURL= {`${$backendUrl}/@group_info/workgroups`}    orderBy={'name'}   orderDirection={'asc'}        metadataTableName={$Terms.table_workgroups}   
                                  tableTitle = {$Terms.datatable_title_workgroups}   tableName='workgroups' formUse='workgroup'   toDoAfterPatchOrPost=''   fieldConcatenateWithId = 'name'
                                  bind:this={currentDataTable}    recordsReadonly=true  /> 
                    {/if}
                {/if}

              {:else}
                <div id='backDrop' class='w-screen h-screen  absolute flex items-center justify-center left-0 top-0 z-10 aria-hidden="true"' class:bg-[rgba(0,0,0,0.2)]={! $workgroupReady} >  

                    {#if ! $isLanguageChosen}
                      <LanguageChoice  on:dispatchProceedWorkgroupVerify={ () => {isWorkgroupReady()} }     />

                    {:else} 
                        {#if $workgroupReady}
                          <!-- mensagem 'carregando...'   -->
                          <div class="divLOADING_MSG">{$Terms.loading}</div>
                        {/if}

                        <!-- se é a 1a vez usuario carregando aplicacao, sistema sorteia um nome de grupo e clona registros, e arquivos de imagem dos carros, para que usuario possa alterar à vontade -->
                        {#if ! $workgroupReady}
                          <WorkgroupCloningModal />
                        {/if} 
                    {/if}

                </div>

              {/if}



          </div>
        </div>


    </div>

    <!-- mensagem rolante usada para avisos na tela principal -->
    <div  id="slidingWindowMessage" >
      &nbsp;
    </div>


</div>

<!-- div que mostra chamadas API -->
<div id='apiDisplay' class:invisible={! $apiCallsVisible}  class:visible={$apiCallsVisible}  >
</div>

<!-- icone cachorrinho canto inferior direito -->
<div class='_doggy'  id='divDoggy'></div>
<div class='_doggy_1' id='divDoggy_1'></div>
<div class='_doggy_2' id='divDoggy_2'></div>
{#if $UnitedStates_selected} 
  <div class='_doggy_3_english' id='divDoggy_3'></div>
{:else}
  <div class='_doggy_3_portuguese' id='divDoggy_3'></div>
{/if}
