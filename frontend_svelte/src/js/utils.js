
import { get } from 'svelte/store'
import { UnitedStates_selected, Terms, HomeCalendar_CurrentFirstDate, soundEnabled, backendUrl, manufacturersAutocomplete, 
    imagesUrl, workgroupReady, backendUrlGolang, backendUrlPHP, clientIp, cloningCurrentBackend } from '$lib/stores/stores.js'
import { invalidateAll } from '$app/navigation'  




//**********************************************************************************************
// obtem país selecionado atualmente (baseado na variavel writable '$UnitedStates_selected' )
//**********************************************************************************************
export const selectedCountry = () => {

  // UnitedStates_selected se refere ao checkbox 'chkLanguageSelector' em page.svelte, que qdo esta ON= USA, OFF= Brazil
  let _country = 'brazil'

  // lê writeable 'UnitedStates_selected' e joga para var local 'UnitedStates_selected'
  UnitedStates_selected.subscribe((UnitedStates_selected) => {
    if (UnitedStates_selected) _country = 'usa'
  })
  return _country

}

//**********************************************************************************************
// obtem idioma selecionado atualmente (baseado na variavel writable '$UnitedStates_selected' )
//**********************************************************************************************
export const selectedLanguage = () => {

  // UnitedStates_selected se refere ao checkbox 'chkLanguageSelector' em page.svelte, que qdo esta ON= USA, OFF= Brazil
  let language = 'portuguese'

  // lê writeable 'UnitedStates_selected' e joga para var local 'UnitedStates_selected'
  UnitedStates_selected.subscribe((UnitedStates_selected) => {
    if (UnitedStates_selected) language = 'english'
  })
  return language
}




//**************************************************************************************************************************************
// transforma o array json recebido com os termos do idioma atual, em um objeto simples que possui a estrutura: item.valor, exemplo:
//
// se UnitedStates_selected = false ,    Terms.errormsg_fill_language_item = Preencha a descrição do item (mín 3 letras) 
// se UnitedStates_selected = true ,    Terms.errormsg_fill_language_item = Fill the item description (min 3 letters)
//**************************************************************************************************************************************
export const simplifyLanguageData = (jsonResponse) => {

  // converte o array associativo recebido filled_json em um objeto JSON simples
  let _language = {}
  let i
  for (i = 0; i < jsonResponse.length; i++) {

    let vlr = jsonResponse[i]

    _language[`${vlr['item']}`] = `${vlr['expression']}`
  }
  return (_language)

}

/************************************************************************************************************************************************************
todo elemento que possuir a classe 'containsTooltip' receberá um tooltip jquery-ui 
um elemento pode inicialmente começar sem tooltip/title, e em seguida , receber tooltip/title
exemplo:   botao 'cancelar pesquisa datatable', - inicialmente nao possui title/tooltip.. mas se usuario pesquisar um texto, 
o botao 'cancelar pesquisa datatable' receberá title/tooltip parecido com 'cancelar pesquisar/limpar pesquisar'
por isso, existe a propriedade '_originalTooltip_', para que o title possa ser retirado/devolvido a qq tempo
************************************************************************************************************************************************************/
export const showOrHideTooltip = () => {

  // a funcao 'showOrHideTooltip' é invocada ja na 1a pagina (+page.svelte) , e tambem é invocada em varios outros pontos da app
  // nao ha como controlar se 'jquery-ui' (responsavel por tooltip()) ja foi carregada
  // necessario testar se ja foi 

  if (typeof jq(window).tooltip == "undefined")  return   // nao carregou ainda

  // define parte visual e o conteudo (title) do title
  jq('.containsTooltip').tooltip({
    tooltipClass: 'prettierTitle_black',
    show: false,   // sem animacao ao exibir
    hide: false,   // sem animacao ao ocultar
    position: { my: "left top+15", at: "left bottom", collision: "flipfit" }
  })

  jq('.containsTooltip').each(function () {
    // usuario pode ter desativado o termo ingles/portugues do respectivo ponto da tela
    if (typeof jq(this).attr('_originalTooltip_') == 'undefined')
      jq(this).attr('title', '? ERR ?')
    else
      jq(this).attr('title', jq(this).attr('_originalTooltip_'))
  })


  // oculta e remove o conteudo do title do elemento que nao é mais para possuir title
  // exemplo: um datatable que nao possui registros selecionados,  o botao 'delete' na parte superior do datatable sera desabilitado e terá seu title desativado abaixo
  jq('.ToRemoveTooltip').each(function () {
    if (jq(this).data('ui-tooltip')) {
      jq(this).tooltip('destroy');
      jq(this).removeAttr("title");
    }
  })

}

/************************************************************************************************************************************************************
por isso, existe a propriedade '_originalTooltip_', para que o title possa ser retirado/devolvido a qq tempo
************************************************************************************************************************************************************/
export const slidingMessage = (div_id, html, time) => {

  let divOBJ = jq(`#${div_id}`)
  divOBJ.html('&nbsp;&nbsp;&nbsp;&nbsp;' + html);
  divOBJ.show("slide", { direction: "left" }, 200);


  setTimeout(function () { divOBJ.hide("slide", { direction: "right" }, 200); }, time);

  // toca beep alertando erro
  //let soundEnabled = get(soundEnabled)   // soundEnabled= writeble store, true/false, usuario define se é pra tocar beep ao mostrar msgs 

  // lê writeable 'soundEnabled' e joga para var local 'soundEnabled'
  soundEnabled.subscribe((soundEnabled_yes) => {
    if (soundEnabled_yes) {
      var thissound = document.getElementById('errorBeep');
      if (thissound['play']) thissound.play();
    }
  })

}


/************************************************************************************************************************************************************
cada input.type=text possui ao seu redor uma div que informa erro, com nome similar, por exemplo:  txtName, possui a div errorName

se detectado erro de digitacao em algum campo, a div de erro sera exibida
************************************************************************************************************************************************************/
export const showFormErrors = (errors) => {

  setTimeout(() => {

    let someError = false
    let tofocusField = ''

    // verifica se algum error detectado
    for (let i=0; i < errors.length; i++) {

      let field = errors[i]

      // 1o campo detectado com erro sera focado
      tofocusField = tofocusField == '' ? field : tofocusField

      someError = true

      // cada input.text possui um respectivo error  
      let divError = field.replace('txt', 'error')

      // exibe a msg de erro do respectivo local - cada campo tem o seu
      jq(`#${divError}`).removeClass('noerrorTextbox').removeClass('errorTextbox').addClass('errorTextbox')
    }

    let terms = get(Terms)   // Terms= writeble store, contendo todas as expressoes usadas, de cada idioma

    if (someError) {
      slidingMessage('slidingFormMessage', terms.errormsg_form_edit_error, 2000)
      jq(`#${tofocusField}`).focus()
    }

  }, 200);
}



/***********************************************************************************************************************************************************/
// rola uma div ate encontrar um determinado elemento 
// por exemplo: uma datatable em que determinado registro foi recem alterado, rola a div container do datatable ate encontrar a divTR do elemento alterado
/***********************************************************************************************************************************************************/
export const scrollUntilElementVisible = (element_id) => {

  if (typeof jq("#" + element_id).attr("id") == "undefined") return; // certifica que existe

  let container_div = jq("#" + element_id).scrollParent(); // div rolavel na qual o elemento esta contido,  se nao houver div, retorna 'window' (tela toda)

  // elemento esta "rolando" dentro de um div (a div tem que possuir um ID !!!)
  let element = document.getElementById(element_id);

  element.scrollIntoView();

  let posY = container_div.scrollTop();
  container_div.scrollTop(posY - 15); // volta um pouco pq scrollIntoView exagera
}


/***********************************************************************************************************************************************************/
// prepara texto grande para gravacao na base de dados
/***********************************************************************************************************************************************************/
export const encodeHTMLEntities = (text) => {
  return $("<textarea/>")
    .text(text)
    .html();
}

/***********************************************************************************************************************************************************/
// converte texto grande lido na base dados em texto apresentavel na tela
/***********************************************************************************************************************************************************/

export const decodeHTMLEntities = (text) => {
  return $("<textarea/>")
    .html(text)
    .text();
}

/***********************************************************************************************************************
 funcao abaixo é chamada apos usuario alterar expressao ingles/portugues, ela recarrega todas as expressoes
 e consequentemente svelte atualiza na tela 

***********************************************************************************************************************/
export async function refreshTerms() {
  let language = selectedLanguage()


  // obtem valor de 'backendUrl', em arquivos .JS nao se pode obter valor de variabel store usando somente '$'
  let _backendUrl_
  backendUrl.subscribe((value) => {
    _backendUrl_ = value
  })

  let _currentWorkgroupName = localStorage.getItem("rentacar_workgroup_name");

  // qdo a pagina esta sendo carregada pela 1a vez, nao ha idioma escolhido ainda, nesse caso, necessario usar as expressoes do grupo 'Admin' mesmo.
  if (_currentWorkgroupName == null)  _currentWorkgroupName='Admin'

console.log('e='+`${_backendUrl_}/${_currentWorkgroupName}/terms/${language}`)
  await fetch(`${_backendUrl_}/${_currentWorkgroupName}/terms/${language}`, { method: 'GET' })

  .then((response) => {

    if (!response.ok) {
      throw new Error(`Language Refresh Err Fatal= ${response.status}`);
    }
    return response.json();
  })

  .then((phrases) => {
    // Terms= variavel writable 
    Terms.set(simplifyLanguageData(phrases))
  })


}

/***********************************************************************************************************************
 faz uma janela ser arrastavel usando o getAndShowCarDetails

***********************************************************************************************************************/
export const makeWindowDraggable = (title_id, window_id) => {
  $(`#${window_id}`).draggable({ handle: `#${title_id}`, containment: '#backDrop' });
}


/***********************************************************************************************************************
 retorna hora atual para forçar carregamento de logotipo atualizado, para ignorar cache de imagem
***********************************************************************************************************************/
export const forceImgRefresh = () => {
  return new Date().getTime()
}


/********************************************************************************************************************************************************
 retorna o SVG da bandeira do país atualmente selecionado
******************************************************************************************************************************************************/
export const svgCurrentCountryFlag = () => {

  if (selectedCountry() == 'brazil')
    return `<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32"><rect x="1" y="4" width="30" height="24" rx="4" ry="4" fill="#459a45"></rect><path d="M27,4H5c-2.209,0-4,1.791-4,4V24c0,2.209,1.791,4,4,4H27c2.209,0,4-1.791,4-4V8c0-2.209-1.791-4-4-4Zm3,20c0,1.654-1.346,3-3,3H5c-1.654,0-3-1.346-3-3V8c0-1.654,1.346-3,3-3H27c1.654,0,3,1.346,3,3V24Z" opacity=".15"></path><path d="M3.472,16l12.528,8,12.528-8-12.528-8L3.472,16Z" fill="#fedf00"></path><circle cx="16" cy="16" r="5" fill="#0a2172"></circle><path d="M14,14.5c-.997,0-1.958,.149-2.873,.409-.078,.35-.126,.71-.127,1.083,.944-.315,1.951-.493,2.999-.493,2.524,0,4.816,.996,6.519,2.608,.152-.326,.276-.666,.356-1.026-1.844-1.604-4.245-2.583-6.875-2.583Z" fill="#fff"></path><path d="M27,5H5c-1.657,0-3,1.343-3,3v1c0-1.657,1.343-3,3-3H27c1.657,0,3,1.343,3,3v-1c0-1.657-1.343-3-3-3Z" fill="#fff" opacity=".2"></path></svg>`

  else
    return `<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32"><rect x="1" y="4" width="30" height="24" rx="4" ry="4" fill="#fff"></rect><path d="M1.638,5.846H30.362c-.711-1.108-1.947-1.846-3.362-1.846H5c-1.414,0-2.65,.738-3.362,1.846Z" fill="#a62842"></path><path d="M2.03,7.692c-.008,.103-.03,.202-.03,.308v1.539H31v-1.539c0-.105-.022-.204-.03-.308H2.03Z" fill="#a62842"></path><path fill="#a62842" d="M2 11.385H31V13.231H2z"></path><path fill="#a62842" d="M2 15.077H31V16.923000000000002H2z"></path><path fill="#a62842" d="M1 18.769H31V20.615H1z"></path><path d="M1,24c0,.105,.023,.204,.031,.308H30.969c.008-.103,.031-.202,.031-.308v-1.539H1v1.539Z" fill="#a62842"></path><path d="M30.362,26.154H1.638c.711,1.108,1.947,1.846,3.362,1.846H27c1.414,0,2.65-.738,3.362-1.846Z" fill="#a62842"></path><path d="M5,4h11v12.923H1V8c0-2.208,1.792-4,4-4Z" fill="#102d5e"></path><path d="M27,4H5c-2.209,0-4,1.791-4,4V24c0,2.209,1.791,4,4,4H27c2.209,0,4-1.791,4-4V8c0-2.209-1.791-4-4-4Zm3,20c0,1.654-1.346,3-3,3H5c-1.654,0-3-1.346-3-3V8c0-1.654,1.346-3,3-3H27c1.654,0,3,1.346,3,3V24Z" opacity=".15"></path><path d="M27,5H5c-1.657,0-3,1.343-3,3v1c0-1.657,1.343-3,3-3H27c1.657,0,3,1.343,3,3v-1c0-1.657-1.343-3-3-3Z" fill="#fff" opacity=".2"></path><path fill="#fff" d="M4.601 7.463L5.193 7.033 4.462 7.033 4.236 6.338 4.01 7.033 3.279 7.033 3.87 7.463 3.644 8.158 4.236 7.729 4.827 8.158 4.601 7.463z"></path><path fill="#fff" d="M7.58 7.463L8.172 7.033 7.441 7.033 7.215 6.338 6.989 7.033 6.258 7.033 6.849 7.463 6.623 8.158 7.215 7.729 7.806 8.158 7.58 7.463z"></path><path fill="#fff" d="M10.56 7.463L11.151 7.033 10.42 7.033 10.194 6.338 9.968 7.033 9.237 7.033 9.828 7.463 9.603 8.158 10.194 7.729 10.785 8.158 10.56 7.463z"></path><path fill="#fff" d="M6.066 9.283L6.658 8.854 5.927 8.854 5.701 8.158 5.475 8.854 4.744 8.854 5.335 9.283 5.109 9.979 5.701 9.549 6.292 9.979 6.066 9.283z"></path><path fill="#fff" d="M9.046 9.283L9.637 8.854 8.906 8.854 8.68 8.158 8.454 8.854 7.723 8.854 8.314 9.283 8.089 9.979 8.68 9.549 9.271 9.979 9.046 9.283z"></path><path fill="#fff" d="M12.025 9.283L12.616 8.854 11.885 8.854 11.659 8.158 11.433 8.854 10.702 8.854 11.294 9.283 11.068 9.979 11.659 9.549 12.251 9.979 12.025 9.283z"></path><path fill="#fff" d="M6.066 12.924L6.658 12.494 5.927 12.494 5.701 11.799 5.475 12.494 4.744 12.494 5.335 12.924 5.109 13.619 5.701 13.19 6.292 13.619 6.066 12.924z"></path><path fill="#fff" d="M9.046 12.924L9.637 12.494 8.906 12.494 8.68 11.799 8.454 12.494 7.723 12.494 8.314 12.924 8.089 13.619 8.68 13.19 9.271 13.619 9.046 12.924z"></path><path fill="#fff" d="M12.025 12.924L12.616 12.494 11.885 12.494 11.659 11.799 11.433 12.494 10.702 12.494 11.294 12.924 11.068 13.619 11.659 13.19 12.251 13.619 12.025 12.924z"></path><path fill="#fff" d="M13.539 7.463L14.13 7.033 13.399 7.033 13.173 6.338 12.947 7.033 12.216 7.033 12.808 7.463 12.582 8.158 13.173 7.729 13.765 8.158 13.539 7.463z"></path><path fill="#fff" d="M4.601 11.104L5.193 10.674 4.462 10.674 4.236 9.979 4.01 10.674 3.279 10.674 3.87 11.104 3.644 11.799 4.236 11.369 4.827 11.799 4.601 11.104z"></path><path fill="#fff" d="M7.58 11.104L8.172 10.674 7.441 10.674 7.215 9.979 6.989 10.674 6.258 10.674 6.849 11.104 6.623 11.799 7.215 11.369 7.806 11.799 7.58 11.104z"></path><path fill="#fff" d="M10.56 11.104L11.151 10.674 10.42 10.674 10.194 9.979 9.968 10.674 9.237 10.674 9.828 11.104 9.603 11.799 10.194 11.369 10.785 11.799 10.56 11.104z"></path><path fill="#fff" d="M13.539 11.104L14.13 10.674 13.399 10.674 13.173 9.979 12.947 10.674 12.216 10.674 12.808 11.104 12.582 11.799 13.173 11.369 13.765 11.799 13.539 11.104z"></path><path fill="#fff" d="M4.601 14.744L5.193 14.315 4.462 14.315 4.236 13.619 4.01 14.315 3.279 14.315 3.87 14.744 3.644 15.44 4.236 15.01 4.827 15.44 4.601 14.744z"></path><path fill="#fff" d="M7.58 14.744L8.172 14.315 7.441 14.315 7.215 13.619 6.989 14.315 6.258 14.315 6.849 14.744 6.623 15.44 7.215 15.01 7.806 15.44 7.58 14.744z"></path><path fill="#fff" d="M10.56 14.744L11.151 14.315 10.42 14.315 10.194 13.619 9.968 14.315 9.237 14.315 9.828 14.744 9.603 15.44 10.194 15.01 10.785 15.44 10.56 14.744z"></path><path fill="#fff" d="M13.539 14.744L14.13 14.315 13.399 14.315 13.173 13.619 12.947 14.315 12.216 14.315 12.808 14.744 12.582 15.44 13.173 15.01 13.765 15.44 13.539 14.744z"></path></svg>`

}

/********************************************************************************************************************************************************
faz busca no back end por lista de registros especifica e vincula o resultado a um input type=text
fazendo com que o texto possua autocomplete 
******************************************************************************************************************************************************/

export const prepareAutocomplete = async (tableName, error_msg, destination_textbox) => {

  // obtem valor de backendUrl, em arquivos .JS nao se pode obter valor de variabel store usando somente '$'
  let _backendUrl
  backendUrl.subscribe((value) => {
    _backendUrl = value
  })

  let _currentWorkgroupName = localStorage.getItem("rentacar_workgroup_name");

  // busca itens para o autocomplete no back end
  let apiURL = `${_backendUrl}/${_currentWorkgroupName}/${tableName}/itens_for_autocomplete/`

  const url = new URL(apiURL);
  logAPI(apiURL)

  await fetch(url, { method: "GET" })

    .then((response) => {

      if (!response.ok) {
        throw new Error(`Itens For Autocomplete Err Fatal= ${response.status}`);
      }
      return response.text()
    })

    .then((items) => {

      let regs = items.split('|')   // divisor de registros

      // cria array temporario que recebera os itens 
      // necessario criar array temporario porque o array vinculado ao input type=text é uma variavel 'writable', que so pode ser alterada via '.set()' (.set() mais abaixo)
      let arrayAttachedToAutocomplete = `${tableName}Autocomplete`    // array destino (variavel writable, compartilhada por toda aplicacao)
      let arrayTMP = []

      for (var i = 0; i < regs.length; i++) {
        var reg = regs[i].split(';');

        arrayTMP.push({ data: reg[0], value: reg[1] })    // data= id do registro         value= texto da opcao
      }

      // joga para a variavel writable vinculada ao input type=text o array criado
      // faz isso porque no momento da gravacao do form, é necessario verificar se o texto digitado existe como item, como opcao valida dentro do array
      eval(`${arrayAttachedToAutocomplete}.set(arrayTMP)`)


      // opcoes gerais do autocomplete 
      let autocomplete_options = {
        minChars: 1,
        autoSelectFirst: true,
        triggerSelectOnValidInput: true,
        lookupLimit: 100,
        maxHeight: 200,
        width: 'auto',
        zIndex: 10030,
        noCache: true,
        deferRequestBy: 0,
      }

      autocomplete_options.lookup = arrayTMP

      // inicia autocomplete no input type=text destino
      jq(`#${destination_textbox}`).autocomplete(autocomplete_options) 

      // se usuario clicar no input type=text destino, força exibicao do autocomplete sem ele ter que digitar uma letra      
      jq(`#${destination_textbox}`).on( "click", function() { 
        try { $(this).autocomplete().show() } 
        catch(err) {} 
      })


      // qdo usuario sair do campo input type=text, verificar se ele digitou vlr valido
      eval(`jq('#${destination_textbox}').on('blur', function() { getAutocompleteItemId('${destination_textbox}', ${arrayAttachedToAutocomplete}) } )   `)
    })
}


//**********************************************************************************************************************************************
// ao sair (blur) de textbox que possui autocomplete, verifica se texto digitado tem correspondencia no array vinculado 
//**********************************************************************************************************************************************
export const getAutocompleteItemId = (txtbox_id, arrayAttachedToAutocomplete) => {

  // arrayAttachedToAutocomplete se refere a uma variavel writeable (store)
  // exemplos: manufacturersAutocomplete, etc
  let arrayTMP = get(arrayAttachedToAutocomplete)

  let typed = jq(`#${txtbox_id}`).val()

  let item_id = 0    // item_id se refere ao ID do registro digitado
  let exists = false;
  jq.each(arrayTMP, function (key, option) {

    // se o texto digitado nao existe no array de opcoes
    if (option.value == typed) {
      item_id = option.data
      exists = true;
      return false;
    }
  });

  if (!exists) {
    jq(`#${txtbox_id}`).val('')     // limpa o input type= text se usuario digitou opcao inexistente
  }

  return parseInt(item_id, 10)  // retorna em formato 'int' para que seja diferenciado de campos text no momento da gravacao

}


//**********************************************************************************************************************************************
// funcoes para implementar for each, usando RANGE
//**********************************************************************************************************************************************

function* iter_range(begin, end, step) {
  // Normalize our inputs
  step = step ? step : 1;

  if (typeof (end) === 'undefined') {
    end = begin > 0 ? begin : 0;
    begin = begin < 0 ? begin : 0;
  }

  if (begin == end) {
    return;
  }

  if (begin > end) {
    step = step * -1;
  }

  for (let x = begin; x < end; x += step) {
    yield x;
  }
}

export function range(begin, end, step) {
  return Array.from(iter_range(begin, end, step));
}

//**********************************************************************************************************************************************
// retorna string em formato hora, dependendo do país selecionado no front end  e baseado em um numero de hora
// Brasil= numero 5, string= 05:00,  numero 22, string= 22:00
// USA= numero 5, string= 05:00 am,  numero 22, string= 10:00 pm
//**********************************************************************************************************************************************
export const hourFormat = (hour) => {

  let country_frontend = selectedCountry()

  let hourTMP = hour

  let am_pm = 'am'

  if (country_frontend == 'usa') {
    // converte exemplo: 23:00 para 11:00 pm
    if (hourTMP > 12) {
      hourTMP -= 12
      am_pm = 'pm'
    }
  }
  // zeros à esquerda
  hourTMP = hourTMP.toString();
  while (hourTMP.length < 2) hourTMP = "0" + hourTMP;

  if (country_frontend == 'usa')
    hourTMP += ':00 ' + am_pm
  else if (country_frontend == 'brazil')
    hourTMP += ':00'

  return (hourTMP)
}



/************************************************************************************************************************************************************
funcao disparada qdo usuario clica no card de algum carro, carrega os detalhes do carro via API e exibe nos respectivos <span> ou <img>
esta funcao serve a: Home.svelte (cards de carros), BookingForm.svelte (opcoes de carros lado direito da tela)

a funcao 'getAndShowCarDetails' pode ter sido chamada de Home.svelte ou BookingForm.svelte
************************************************************************************************************************************************************/
export async function getAndShowCarDetails(carId) {


  // em Home.svelte, os dias em que o carro selecionado esta reservado, sao marcados com um X vermelho (classe= RESERVED_DAY)
  // remove aqui todos os dias que possuem a classe RESERVED_DAY, pois o calendario sera redesenhado mais abaixo
  jq('.RESERVED_DAY').removeClass('RESERVED_DAY')


  // esconde cards que vao conter os detalhes do carro ate terminar de carregar os dados
  setTimeout(() => {
    jq('.bookingCarDetail').css('visibility', 'hidden')           // .bookingCarDetail , BookingForm.svelte
    jq('.homeCarDetail').css('visibility', 'hidden')           // .homeCarDetail , Home.svelte
  }, 1);

  // se esta exibindo detalhes dentro de Home.svelte, passa a informacao do mês selecionado atualmente (calendario lado esquerdo da tela)
  // para que back end retorne os dias do mês em que o carro esta reservado
  let isHomeVisible = typeof jq('#homeCarsBrowser').attr('id')!='undefined'

  // obtem valor de 'backendUrl', em arquivos .JS nao se pode obter valor de variavel store usando somente '$' como os arquivos .svelte podem
  let _backendUrl_
  backendUrl.subscribe((value) => {
    _backendUrl_ = value
  })

  // obtem valor de 'imagesUrl', em arquivos .JS nao se pode obter valor de variavel store usando somente '$'  como os arquivos .svelte podem
  let _imagesUrl_
  imagesUrl.subscribe((value) => {
    _imagesUrl_ = value
  })

  // obtem valor de 'imagesUrl', em arquivos .JS nao se pode obter valor de variavel store usando somente '$'  como os arquivos .svelte podem
  let _terms_
  Terms.subscribe((value) => {
    _terms_ = value
  })


  let url = ''


  // devido a assincronicidade/rapidez do svelte, a funcao atual (getAndShowCarDetails) pode ser chamada sem haver ID de carro selecionado, clicado
  // por isso necessario o IF abaixo.
  if (carId==0)  return;

  // exibe animacao que avisa pausa para processamento
  setTimeout(() => {
    showLoadingGif();
  }, 1);


  // se 'getAndShowCarDetails' foi chamado por Home.svelte 
  // calendario (Home.svelte) visivel, busca alem dos detalhes do carro, os dias do mês/ano que carro ja foi reservado, para marcar um X nos dias em que esta reservado
  if (isHomeVisible) {
      // obtem mês/ano atualmente selecionados no calendario (lado esquerdo da tela)
      let _HomeCalendar_CurrentFirstDate_
      HomeCalendar_CurrentFirstDate.subscribe((value) => {    //  HomeCalendar_CurrentFirstDate, variavel writeble (stores.js)
        _HomeCalendar_CurrentFirstDate_ = value
      })

    let firstDay = new Date(_HomeCalendar_CurrentFirstDate_.getFullYear(), _HomeCalendar_CurrentFirstDate_.getMonth(), _HomeCalendar_CurrentFirstDate_.getDate());
    let lastDay = new Date(firstDay.getFullYear(), firstDay.getMonth()+1, 0);
 
    let _firstDay = firstDay.toISOString().split('T')[0]   // yyyy-mm-dd
    let _lastDay = lastDay.toISOString().split('T')[0]

     // se calendario (Home.svelte) visivel, busca alem dos detalhes do carro, os dias do mês/ano que carro ja foi reservado

     url = new URL(`${_backendUrl_}/car/${carId}/${_firstDay}/${_lastDay}`);
  }

  // 'getAndShowCarDetails' foi chamado por BookingFomr.svelte ou Bookings.svelte
  else 
     url = new URL(`${_backendUrl_}/car/${carId}`);

  logAPI('GET',url)
  await fetch(url, { method: 'GET' }) 

    .then((response) => {

      if (!response.ok) {
        throw new Error(`Read Car  Err Fatal= ${response.status}`);
      }
      return response.json();
    })

    .then((fields) => {

      // popula os detalhes do carro  recebidos da API

      Object.keys(fields).forEach(function (key, idx, arr) {

        // joga vlr do campo recebido via API para a respectiva <span>  ou <img> ou <checkbox>, ou etc etc
        
        // se campo atual é o nome do arquivo de imagem do carro
        let url_image  = ''

        switch (key) {
            case 'car_image':
              url_image = `${_imagesUrl_}` + fields[key]
              jq('#carDetails_Picture').css('background-image', `url(${url_image})`);
              break

            // se campo atual é o nome do arquivo de imagenm do carro
            case 'manufacturer_logo':
              url_image = `${_imagesUrl_}` + fields[key]

              jq('#carDetails_ManufacturerLogo').css('background-image', `url(${url_image})`);
              break

            case 'transmission_manual':

              // fields[key] == 1, transmissao manual   (true)
              // fields[key] == 0, transmissao automatica   (false)

              // cards de carros em Home.svelte
              if (fields[key] == 0) jq('#carDetails_transmission').html(_terms_.automatic_transmission)   // cambio automatico
              if (fields[key] == 1) jq('#carDetails_transmission').html(_terms_.manual_transmission)               // cambio manual

              // cards de carros em Bookings.svelte
              if (fields[key] == 0) jq('#carDetails_transmission_short').html(_terms_.automatic_transmission_short)   // cambio automatico
              if (fields[key] == 1) jq('#carDetails_transmission_short').html(_terms_.manual_transmission_short)               // cambio manual

              break

            // se campo atual sao os dias em que o carro esta reservado
            case 'days_reserved':
              let days = fields[key].split(',')   // dias reservados vem separados por ','

              // cada dia no calendario do mes/ano corrente, possui a classe (sem CSS), 'dayOfCalendar'
              // se o dia exibido no calendario foi detectado como dia reservado do veiculo, coloca um X vermelho sobre o dia
              jq('.dayOfCalendar').each(function () {
                let dayInCalendar = parseInt(jq(this).html(), 10)

                for (let i=0; i < days.length; i++)  {
                  let reservedDay = parseInt(days[i], 10)
                  if (reservedDay ==  dayInCalendar)   jq(this).addClass('RESERVED_DAY')
                }
              })

              break

            default:
              // se campo atual é qq outro campo
              jq(`#carDetails_${key}`).html(fields[key])

        }
      })


      // exibe cards que receberam os detalhes do carro
      setTimeout(() => {
        jq('.bookingCarDetail').css('visibility', 'visible')           // .bookingCarDetail , BookingForm.svelte
        jq('.homeCarDetail').css('visibility', 'visible')           // .homeCarDetail , Home.svelte

      }, 100);

      // esconde animacao que avisa pausa para processamento
      setTimeout(() => { hideLoadingGif(); }, 500)
    })
}


/************************************************************************************************************************************************************
snake case para pascal case 
usada quando recebe campos Json e converte para input type=text
************************************************************************************************************************************************************/
const snakeToCamel = str => str.replace( /([-_]\w)/g, g => g[ 1 ].toUpperCase() );
export const snakeToPascalCase = str => {
    let camelCase = snakeToCamel( str );
    let pascalCase = camelCase[ 0 ].toUpperCase() + camelCase.substr( 1 );
    return pascalCase;
}

/************************************************************************************************************************************************************
preenche zeros à esquerda
************************************************************************************************************************************************************/
export const leadingZeros = (str, size) => {
    while (str.length < size) str = "0" + str;
    return str;
}


/************************************************************************************************************************************************************
prepara data em formato string, considerando o UTC local
usado quando envia-se datas para gravacao no back end
***********************************************************************************************************************************************************/
export const dateToIsoStringConsideringLocalUTC = (date) => {
  var tzo = -date.getTimezoneOffset(),
      dif = tzo >= 0 ? '+' : '-',
      pad = function(num) {
          return (num < 10 ? '0' : '') + num;
      };

  return date.getFullYear() +
      '-' + pad(date.getMonth() + 1) +
      '-' + pad(date.getDate()) +
      'T' + pad(date.getHours()) +
      ':' + pad(date.getMinutes()) +
      ':' + pad(date.getSeconds()) +
      dif + pad(Math.floor(Math.abs(tzo) / 60)) +
      ':' + pad(Math.abs(tzo) % 60);
}

/************************************************************************************************************************************************************
verifica se ha 'workgroup' definido, se nao houver, busca um workgroup aleatorio dentro da base de dados
e o define como grupo padrao deste computador, ate que o usuario queira muda lo

getWorkgroupRecordsReady() serve para 2 situacoes:
1- sortear, gerar um grupo, caso usuario acessando aplicacao pela 1a vez (whatToDo= generate)
2- resetar os dados do grupo atual, caso usuario ja esteja em um grupo e pediu para resetar   (whatToDo= reset)
3- sortear outro nome de grupo, usuario ja tem grupo e pediu para sortear outro nome  (whatToDo= another)
***********************************************************************************************************************************************************/

export async function getWorkgroupRecordsReady( whatToDo  = false ) {

// memoriza IP do cliente para que alteracoes feitas pelo proprio cliente nao sejam exibidas
// exibe somente alteracoes feitas por outros clientes do mesmo grupo
var xmlHttp = new XMLHttpRequest();
xmlHttp.open( "GET", 'https://ipinfo.io/json?token=08d19598882aee', false ); // false for synchronous request
xmlHttp.send( null );
let _clientInfo =  xmlHttp.responseText 

let tmp = JSON.parse(_clientInfo);
clientIp.set( tmp.ip )
 
  
let terms = get(Terms)     // obtem expressoes ingles/portugues

let _currentWorkgroupName = localStorage.getItem("rentacar_workgroup_name");

let workgroupInform

// getWorkgroupRecordsReady() foi chamada do menu de grupos, usuario pediu ou para resetar os dados ou para sortear (aplicacao) outro nome de grupo
if (whatToDo == 'reset' || whatToDo == 'another')  {
  workgroupInform = _currentWorkgroupName
}

// getWorkgroupRecordsReady() foi chamada no carregamento da aplicacao
else {
  // se workgroup ja obtido, siginifca que aplicacao foi executada anteriormente, grupo foi sorteado e dados default copiados
  if (_currentWorkgroupName!=null)  {
    jq('#btnWorkgroupMenu').html ( _currentWorkgroupName )
    refreshTerms()

    return false
  }    

  // aplicacao sendo executada pela 1a vez, pois nao ha workgroup definido
  workgroupInform = 'none'
}

workgroupReady.set(false)   // avisa que esta sem grupo por enquanto, para que a pagina nao continue carregamento

try {
    // obtem valor de 'backendUrl', em arquivos .JS nao se pode obter valor de variavel store usando somente '$'
    // a geracao de dados, clonagem de dados de workgroup é a unica tarefa que so o backend golang faz
    // os demais backends nao fazem essa tarefa
    // manipulacao de grupos e notificacoes sao feitas somente pelo backend golang
    let _backendUrlGolang_
    backendUrlGolang.subscribe((value) => {
      _backendUrlGolang_ = value 
    })

    cloningCurrentBackend.set( _backendUrlGolang_ )

    await fetch(`${_backendUrlGolang_}/workgroup/${whatToDo}/${workgroupInform}`, 
                  { method: 'POST', body: _clientInfo,  headers: { 'Content-Type': 'application/json' } })   

    .then((response) => {

      if (!response.ok) {
        jq('#btnWorkgroupMenu').html ( '* error *' )  
        throw new Error(`Workgroup Get Err Fatal= ${response.status}`);
      }
      return response.text();            // infoWorkgroup = workgroup name|workgroup id
    })

    .then((infoWorkgroup) =>  {
      if ( infoWorkgroup.indexOf('__success__')!=-1 )   {
        let _infoWorkgroup = infoWorkgroup.split('|')             // infoWorkgroup = __success__|workgroup name|workgroup id
        _currentWorkgroupName = _infoWorkgroup[1]

        // grava workgroup obtido como local storage
        localStorage.setItem("rentacar_workgroup_name", _currentWorkgroupName);
        // reseta notificacoes do grupo
        localStorage.setItem("rentacar_last_notification_id_read", 0);

      } 
      else  {
        // pode haver todo tipo de erro ao clonar os registros default
        slidingMessage('slidingWindowMessage', terms.errormsg_general_database + '&nbsp;&nbsp;Error workgroup= '+infoWorkgroup, 5000)
        throw infoWorkgroup
      }
    })


// se erro conexao, erro jscript
} catch (error) {
  jq('#btnWorkgroupMenu').html ( '* error *' )  
  slidingMessage('slidingWindowMessage', terms.errormsg_general_database + '&nbsp;&nbsp;Error workgroup= '+error.message, 5000)

  return false
}
}


/************************************************************************************************************************************************************
funcao executada a cada 1/2 segundo (setInterval)  apos iniciado processo de clonagem dos dados (grupo admin ->  novo grupo)

exibe situacao atual da clonagem dos registros/arquivos
***********************************************************************************************************************************************************/
export async function showCloningStatus() {


  let _cloningCurrentBackend_
  cloningCurrentBackend.subscribe((value) => {
    _cloningCurrentBackend_ = value 
  })

  if (_cloningCurrentBackend_=='') return;

  let _route_ = `${_cloningCurrentBackend_}/cloning_status`
  logAPI('GET', _route_)
  await fetch(_route_, { method: 'GET' })

  .then((response) => {  
      if (response.ok) return response.text()
  }) 

  .then( (cloningStatus) => {       

      // talvez o processo de clonagem nao tenha sido iniciado ainda   ( cloningStatus=='' )
      // leva uns segundos entre o inicio do processo (getWorkgroupRecordsReady) e a verificacao do seu status (showCloningStatus)
      if (cloningStatus!='')   {

          // ***************************************************************************************************************
          // cloningStatus =>  status info adaptar ao idioma front end^info adicional para concatenar|percent info|workgroup name
          // ***************************************************************************************************************

          // nome do workgroup escolhido aleatoriamente no back end (sao 150 workgroups disponiveis)
          let chosenWorkgroupName = cloningStatus.split('|')[2]
          jq('#chosenWorkgroupName').html( chosenWorkgroupName )


          // situacao atual e percentual concluido
          let termToLanguageAdjust = ''
          let textToConcatenate = ''

          if (typeof cloningStatus.split('|')[0].split('^')[1]!='undefined')    {    // status info^some text to concatante|percent info|workgroup name
            termToLanguageAdjust = cloningStatus.split('|')[0].split('^')[0]
            textToConcatenate =  '&nbsp;&nbsp;' + cloningStatus.split('|')[0].split('^')[1]
          } 
          else 
            termToLanguageAdjust = cloningStatus.split('|')[0]

          let progressBarWidth = jq('.divLOADING_BAR_OUTER').width()
          let completedPercent = parseInt(cloningStatus.split('|')[1], 10)

          let terms = get(Terms)   // Terms= writeble store, contendo todas as expressoes usadas, de cada idioma

          // devido à assincronicidade do svelte é sempre bom testar se o conteudo, a frase veio (!= undefined)
          if ( typeof terms[termToLanguageAdjust]!='undefined' ||  termToLanguageAdjust=='')  {
            let toPrint
            // se backend nao mandou nada no 1o argumento (divisor= ^), exibe só o 2o
            if ( termToLanguageAdjust=='')  
              toPrint = textToConcatenate
            else
              toPrint = terms[termToLanguageAdjust] + textToConcatenate

            if (toPrint.length > 60)
              jq('#workgroupCloningStatus').css('fontSize', '22px') 
            else 
              jq('#workgroupCloningStatus').css('fontSize', '24px') 

            jq('#workgroupCloningStatus').html( toPrint ) 

            jq('#completedPercent').html( completedPercent + '&nbsp;%' ) 
            jq('#workgroupCloningPercent').width( progressBarWidth * completedPercent / 100 )
          }
      }
  })
}

/************************************************************************************************************************************************************
le qual backend escolhido no momento e joga para variavel generica 
************************************************************************************************************************************************************/
export const setBackendUrl = ( backendToSet = '' ) => {

let _backendUrlGolang_
backendUrlGolang.subscribe((value) => {
  _backendUrlGolang_ = value 
})

let _backendUrlPHP_
backendUrlPHP.subscribe((value) => {
  _backendUrlPHP_ = value 
})

// se foi passado backend para memorizar
if (backendToSet!='')
  localStorage.setItem("rentacar_current_backend", backendToSet)

// se usuario nao optou por backend ainda, inicia usando golang 
if ( localStorage.getItem("rentacar_current_backend") == null ) 
  localStorage.setItem("rentacar_current_backend", 'golang')

let backend = localStorage.getItem("rentacar_current_backend")

if (backend=='golang')   backendUrl.set(_backendUrlGolang_)
if (backend=='php')   backendUrl.set(_backendUrlPHP_)

// icon exemplo:  golang_symbol.svg,  php_symbol.svg, etc
// icone canto superior esquerdo (seletor de backend)
jq('#backendIcon').attr('src', `${backend}_symbol.svg`)
}

/************************************************************************************************************************************************************
adiciona na div que exibe as APIS, a ultima chamada API feita 
************************************************************************************************************************************************************/
export const logAPI = (method, url) => {
  let html = `<div style='display:flex;flex-direction:row;'>`+
          `<span style='width:80px;'>${method}</span>`+
          `${url}`+
          '</div>'

  jq('#apiDisplay').append(html)
  let div = document.getElementById('apiDisplay');
  div.scrollTop = div.scrollHeight;
}

/************************************************************************************************************************************************************
aguarda segundos para executar a proxima linha
************************************************************************************************************************************************************/
const delayToContinue = ms => new Promise(res => setTimeout(res, ms));


/************************************************************************************************************************************************************
verifica se ha 'workgroup' definido, se nao houver, busca um workgroup aleatorio dentro da base de dados
e o define como grupo padrao deste computador, ate que o usuario queira muda lo

getWorkgroupRecordsReady() serve para 2 situacoes:
1- sortear, gerar um grupo, caso usuario acessando aplicacao pela 1a vez (whatToDo= generate)
2- resetar os dados do grupo atual, caso usuario ja esteja em um grupo e pediu para resetar   (whatToDo= reset)
3- sortear outro nome de grupo, usuario ja tem grupo e pediu para sortear outro nome  (whatToDo= another)
***********************************************************************************************************************************************************/

export async function getWorkgroupFilesReady() {

let terms = get(Terms)     // obtem expressoes ingles/portugues

let _currentWorkgroupName = localStorage.getItem("rentacar_workgroup_name");

let _workgroupReady_
workgroupReady.subscribe((value) => {
  _workgroupReady_ = value 
})

// devido à rapidez do svelte, talvez a funcao atual foi chamada sem necessidade (_workgroupReady_ ja esta ok)
if (_workgroupReady_) return;


try {
    // o backend golang fará a clonagem dos arquivos de imagem
    let _backendUrlGolang_
    backendUrlGolang.subscribe((value) => {
      _backendUrlGolang_ = value 
    })

    cloningCurrentBackend.set( _backendUrlGolang_ )

    // clona os arquivos do grupo 'Admin' para o grupo cujos registros foram clonados acima (_currentWorkgroupName)
    await fetch(`${_backendUrlGolang_}/workgroup/clone_image_files/${_currentWorkgroupName}`, { method: 'POST'})


    .then((response) => {
      if (!response.ok) {
        jq('#btnWorkgroupMenu').html ( '* error *' )  
        throw new Error(`Workgroup Get Err Fatal= ${response.status}`);
      }
      return response.text();            
    }) 

    .then((cloningResult) => {
      if ( cloningResult.indexOf('__success__')!=-1 )   {

          setTimeout(() => {
            workgroupReady.set(true)  
            invalidateAll()   // atualiza dados na tela
          }, 2500);        

          jq('#btnWorkgroupMenu').html ( _currentWorkgroupName )   // botao 'grupo' ,  Home.svelte
          refreshTerms()

          return true
      }
      else  {
        // pode haver todo tipo de erro ao clonar os arquibvos de imagem
        slidingMessage('slidingWindowMessage', terms.errormsg_general_database + '&nbsp;&nbsp;Error workgroup= '+cloningResult, 5000)
        throw cloningResult
      }

    })


// se erro conexao, erro jscript
} catch (error) {
  jq('#btnWorkgroupMenu').html ( '* error *' )  
  slidingMessage('slidingWindowMessage', terms.errormsg_general_database + '&nbsp;&nbsp;Error workgroup= '+error.message, 5000)

  return false
}
}

