<script>
import { Terms, backendUrl, imagesUrl,staticImagesUrl,  manufacturersAutocomplete, clientIp, imagesStillLoading, recordIdsVisible  } from '$lib/stores/stores.js'
import { refreshTerms, slidingMessage, scrollUntilElementVisible, selectedCountry, showOrHideTooltip, 
      prepareAutocomplete, getAutocompleteItemId, logAPI, forceImgRefresh,  makeWindowDraggable, svgCurrentCountryFlag, 
      encodeHTMLEntities, showFormErrors, snakeToPascalCase, decodeHTMLEntities } from '$js/utils.js'


import TermForm from '$edit_forms/TermForm.svelte'
import ManufacturerForm from '$edit_forms/ManufacturerForm.svelte'
import CarForm from '$edit_forms/CarForm.svelte'

import jQuery from 'jquery'

import ConfirmationWindow from '$lib/components/ConfirmationWindow.svelte'
import RecordMetadataWindow from '$lib/components/RecordMetadataWindow.svelte'

import { setContext } from 'svelte';



// contem os dados lidos via API
let dataList

// controla se é para exibir a dica que fica logo abaixo do input.text de pesquisa (txtTableSearchText)
let showTipSearchbox = false

// controla se é para mostrar somente registros ativos ou inativos
let showOnlyActiveOrInactiveRecords = ''

// memoriza qual ultimo registro foi inserido ou atualizado, para poder destaca lo apos refresh do datatable  
let lastRecordIdUpdated = 0

// se forceDatatableReload aumentada, faz recarga da tabela
export let forceDatatableReload = 0

// formHttpMethodApply pode ser: POST (novo registro), PATCH (editar registro) ou DELETE - o mesmo form de edicao é usado todas as operacoes
let formHttpMethodApply = ''

// se calledEditForm= form de edicao foi chamado pelo usuario
let calledEditForm = false

// se showDeleteSelectedRecordsConfirmation= janelinha modal perguntando algo foi chamada
let showDeleteSelectedRecordsConfirmation = false

// contem o ID do registro que sera manipulado
let currentFormRecordId 

// contem ID's dos registros selecionados (usando o checkbox, 1a coluna, de cada divTR)
let selectedRecords = []
let deleteMessageWarning = ''

// formUse = recebe qual o tipo de registro sendo manipulado: car, manufacturer, term, etc etc
export let formUse = ''

// recebe alguma funcao jscript que deve ser feita apos perform_crud executado com sucesso
export let toDoAfterPatchOrPost = ''    

// tabelas com 'height' maiores para que caiba um logotipo ou imagem
let tablesWithBiggerRow = ['manufacturers', 'cars']

// se showRecordMetadata= svelte exibe tela com metadados 
let showRecordMetadata = false

// campo cujo ID deve ID do registro deve ser exibido  à sua esquerda
export let fieldConcatenateWithId = ''

// cada chamada API envia nome/id do grupo para que o usuario só manipule dados do seu grupo
let _currentWorkgroupName = localStorage.getItem("rentacar_workgroup_name");

// controla se a tabela ja foi carregada a 1a vez
let datatableReady = false


/************************************************************************************************************************************************************
usuario pediu para ver metadata (data criacao, id criador, etc) do registro de reserva
************************************************************************************************************************************************************/
export async function fillRecordMetadataWindow (record_id) {

  // necessario abrir evento assincrono para exibir div ajax loading, caso contrario navegador nao atualiza a tela
  setTimeout(() => {
    showLoadingGif();   // mostra animacao 'processando..' mesmo sem tela ter sido renderizada
  }, 1);

  let country = selectedCountry()

  // necessario passar o país para que as datas (created_at, etc) retornem no formato usado naquele país
  let url = `${$backendUrl}/${tableName}/record_metadata/${record_id}/${country}/`
console.log('url='+url)
  logAPI('GET', url)
  
  logAPI('GET', url)
  await fetch(url, {method: 'GET'})

  .then((response) => {
    if (!response.ok) {
      throw new Error(`Read Record Metadata Err fatal= ${response.status}`);
    }
    return response.json();
  })

  .then((result) => {
    setTimeout( () => {hideLoadingGif();}, 500)

    let fields = result

    // joga campo json recebido para seu respectivo <span>
    // exemplo:  json: created_at  ==>   <span id='info_created_at'>
    Object.keys(fields).forEach(function(key, idx, arr)   {
      jq(`#info_${key}`).html( fields[key] )      
    })
  })
}

/********************************************************************************************************************************************************
 valida dados do formulario e se tudo ok, tenta gravar
********************************************************************************************************************************************************/
const performSaveRecord = () =>  {

  // remove qq msg erro que foi exibida previamente
  jq('.errorTextbox').removeClass('errorTextbox').addClass('noerrorTextbox')      // cada div que exibe erro, e possui a classe 'errorTextbox'

  let errors = []

  var formData = new FormData(); 
  formData.append('country', selectedCountry())    // alguns tabelas precisam do país selecionado atualmente - exemplo: tabela carros, grava/separa os carros por país 

  let field_value

  jq('#recordForm').find('input[type=text], textarea, input[type=file], input[type=checkbox]').each(function(e) {	 
    let originalId = jq(this).attr('id')

    let inputType = ''
    if ( typeof jq(this).attr('type')!='undefined' ) inputType = jq(this).attr('type').toLowerCase()

    let elementType
    if ( typeof jq(this).prop("tagName").toLowerCase()!='undefined' ) elementType = jq(this).prop("tagName").toLowerCase()

    if ( elementType =='textarea' )  field_value = encodeHTMLEntities(jq(this).val())


    // remove 'txt/chk' do Id campo para facilitar leitura parametro no back-end
    let snakecaseID = jq(this).attr('id').replace('txt','').replace('chk','')    

    // converte pascal case para snake case, exemplo:  ManufacturerName ==>  manufacturer_name
    snakecaseID = snakecaseID.split(/\.?(?=[A-Z])/).join('_').toLowerCase();   

    let max = parseInt(jq(this).prop("maxlength"), 10)
    let min = parseInt(jq(this).attr("minlength"), 10)   // inventado (contrived)

    // se campo atual é text ou textarea e o texto digitado maior / menor do que limites 
    if ( inputType == 'text' || elementType=='textarea' )  { 
       field_value = jq(this).val()                      // input type= text / textarea

      if (field_value.length>max || field_value.length<min) errors.push( jq(this).attr('id') )
    }

    // campos com autocomplete necessario tratamento diferente
    if (originalId == 'txtManufacturerName')  {
      let manufacturer_id = getAutocompleteItemId('txtManufacturerName', manufacturersAutocomplete) 
      formData.append('manufacturer_id', manufacturer_id)
     }

    // montando body...
    // exemplos:
    // body.name = 'type name'
    // body.doors = 3 , etc etc

    // input type=file 
    if ( inputType=='file' )  { 

      // toda imagem possui um respectivo botao de upload , exemplo:  <img id='imgCarImage'>      upload é feito pelo botao <input type=file id='fileCarImage' >
      // relacao:   imgCarImage  ===>   fileCarImage

      // se nao foi feito upload de imagem no formulario, verifica se ja ha imagem exibida no elemento <img>, ou seja, nao exige upload 
      if (typeof jq(this)[0].files[0]=='undefined') {          // se nao foi feito upload
        let matchingImage = originalId.replace('file','img')    

        // 'src' ainda possui a imagem em branco, carregada durante a montagem do formulario (http), informa erro, pois usuario nao escolheu nova imagem
        if ( jq(`#${matchingImage}`).attr('src').toLowerCase().indexOf('none_image')!=-1 )   {
          errors.push( 'imageError' )

        // 'src' possui uma imagem carregada durante a montagem do formulario (http), registro esta sendo editado, avisa back end (bypass_image_upload)
        // nao foi escolhida imagem, mas nao é necessario pois o registro esta sendo editado e ja existe imagem definida e gravada
        } else {
          formData.append('bypass_image_upload', true);
        }

      // uma imagem foi escolhida, necessario fazer validacao no back end
      } else  {
        formData.append('bypass_image_upload', false);
        formData.append('chosen_image_file', jq(this)[0].files[0]);   // backend usa o mesmo ID mas vai saber diferenciar fabricante/carro pela rota de gravacao
      }
    }

    else if ( inputType=='checkbox' )  { 
      formData.append(`${snakecaseID}`, jq(this).is(":checked"));

    }  

    // campo text, textarea
    else 
      formData.append(`${snakecaseID}`, field_value);

  })

  if (errors.length>0) {
    // alguns tipos de erro nao possuem <div> especifica para exibi los,  exibir mensagem de erro rolante quando usuario nesse caso
    if ( errors.includes('imageError') )  {
      slidingMessage('slidingFormMessage', $Terms.errormsg_missing_image, 2000)    // imagem nao foi escolhida
      return;
    }

    // erros que possuem <div error> destino
    showFormErrors(errors) ;  
    return;
  }

  _currentWorkgroupName = localStorage.getItem("rentacar_workgroup_name");

  let route = `${_currentWorkgroupName}/`

  if (formHttpMethodApply=='POST') 
    route += `${formUse}`        //  car, manufacturer, user, etc etc
  if (formHttpMethodApply=='PATCH') 
    route += `${formUse}/${currentFormRecordId}` 

  // informa IP do cliente para que nao exiba notificacoes para quem fez a alteracao, 
  // somente para quem pertence ao grupo mas nao fez a alteracao
  formData.append('client_ip', $clientIp)   

  performCrudOperation(formData, formHttpMethodApply, route) 
}

/****************************************************************************************************
usuario clicou no botao 'metadados'
****************************************************************************************************/
const invokeMetadataWindow = () =>  {
  showRecordMetadata=true
}



/****************************************************************************************************
 se usuario passar mouse sobre textbox de pesquisa, coloca foco sobre ele, se nao tiver colocado ainda
****************************************************************************************************/
const focusDatatableSearchBox = () =>  {
  if ( ! jq('#txtTableSearchText').is(":focus") )  {
    jq('#txtTableSearchText').focus();
    jq('#txtTableSearchText').select();
  }
}



/****************************************************************************************************
 reseta filtro aplicado ao datatable
****************************************************************************************************/
const resetDatatableFilter = () => {
  jQuery('#txtTableSearchText').val('');
  showOnlyActiveOrInactiveRecords='';
  textToSearch=''

  forceDatatableReload++
}




/********************************************************************************************************************************************************
prepara funcionalidades iniciais dos formularios invocados por Datatable.svelte e dispara funcao que carrega os dados do registro sendo editado
********************************************************************************************************************************************************/
const getFormPopulatedAndReady = async () => { 

  $imagesStillLoading = true    // mostra animacao gif 'processnado'.. durante o carregamento de imagens
 

  // alguns forms de edicao exibem a bandeira do país atualmente selecionado atraves da div 'divCountryFlag'  
  // caso esteja inserindo registro (POST) mostra bandeira do país selecionado no front
  if (formHttpMethodApply=='POST')    
    jq('#divCountryFlag').html( svgCurrentCountryFlag() )

  // aguarda 300 milisegundos para que os ajustes acima fiquem prontos
  setTimeout(() => {

      // se vai manipular um registro ja existente (!=POST), carrega seus dados primeiro
      // POST = novo registro
      // != POST= PATCH, DELETE
      if (formHttpMethodApply!='POST')    {
        // 'formUse'= car, language, etc  (get_car, get_language, etc)
        // putFocusInFirstInputText_AndOthersParticularitiesOfTheForm() sera executada dentro de performCrudOperation
        performCrudOperation(null, 'GET', `${formUse}/${currentFormRecordId}`)    
      }


      // vai adicionar registro , campos em branco , coloca foco no 1o campo 
      if (formHttpMethodApply=='POST')    {
        // alguns forms possuem imagens (carro, logotipo fabricante, etc),  joga imagem em branco em <img> do form  pois esta inserindo registro
        setTimeout(() => {
          jq("#recordForm img").each(function() {	 
            jq(this).attr('src', `${$staticImagesUrl}/none_image.png`)  
            jq(this).removeClass('imageLoadedInForm')
            jq(this).addClass('imageLoadedInForm')

          })
        }, 300);

        putFocusInFirstInputText_AndOthersParticularitiesOfTheForm()   

      }
  }, 300);

}


/************************************************************************************************************************************************************
recebe JSON que foi obtido via FETCH API e popula elementos <input type=text, checkbox> , <img> e <span>

cada chave JSON diz respeito a um '<input type=text ou checkbox>' ou '<span>' ou '<img>' destino

para saber qual é o elemento html destino, remove-se os strings txt/img e chk do ID do elemento html e verifica se se 
ha uma chave JSON respectiva

exemplo:
{name: 'John', age: 34}   input type=text: txtJohn,  vai receber o valor:  'John' 
                          input type=text: txtAge,  vai receber o valor:  '34'  

************************************************************************************************************************************************************/
export async function populateFormDestinationFields(pairs) {

  Object.keys(pairs).forEach(function (key, idx, arr) {

    // as linhas abaixo verificam se a chave (key) json recebida possui um respectivo campo input type=text,input type=checkbox ou <img>
    // e se possuir, joga o valor recebido para o respectivo campo 

    // o JSON recebido pode fazer referencia a um campo input=text ou <img> ou input=checkbox
    // exemplo: json recebido:  'name' => equivale a input type=text id='txtName'
    // exemplo: json recebido:  'car_image' => equivale a <img id='imgCarImage'>

    let field_id_text = 'txt' + snakeToPascalCase(key) 
    let field_id_img = 'img' + snakeToPascalCase(key)
    let field_id_checkbox = 'chk' + snakeToPascalCase(key)

    let field_id

    // obtem valor de 'imagesUrl', em arquivos .JS nao se pode obter valor de variavel store usando somente '$' 
    let _imagesUrl_
    imagesUrl.subscribe((value) => {
      _imagesUrl_ = value
    })

    // se campo destino existe
    if (typeof jq(`#${field_id_text}`).attr('id') != 'undefined' || typeof jq(`#${field_id_img}`).attr('id') != 'undefined' || typeof jq(`#${field_id_checkbox}`).attr('id') != 'undefined') {
      // qdo form usado para 'delete', os campos sao <span>
      // qdo form usado para edicao, insercao, campos sao input.text ou textarea

      if (typeof jq('#' + field_id_text).attr('id') != 'undefined') field_id = field_id_text
      if (typeof jq('#' + field_id_img).attr('id') != 'undefined') field_id = field_id_img
      if (typeof jq('#' + field_id_checkbox).attr('id') != 'undefined') field_id = field_id_checkbox

      let tagType = jq(`#${field_id}`).prop("tagName").toLowerCase()   // <img> ou <input>
      let inputType = jq(`#${field_id}`).prop("type")  // txt, textarea, checkbox, radio

      // input.text e textarea sao usados qdo form usado para edicao
      if (inputType == 'text') jq(`#${field_id}`).val(pairs[key])
      if (inputType == 'textarea') jq(`#${field_id}`).val(decodeHTMLEntities(pairs[key]))

      // aguarda 1/2 segundo  para carregar a imagem para que o gif animado seja exibido e dê impressao ao usuario que esta carregando
      if (tagType == 'img') {
        setTimeout(() => {
          jq(`#${field_id}`).attr('src', _imagesUrl_ + pairs[key] + '?' + forceImgRefresh())
        }, 1);
      }

      if (inputType == 'checkbox') {
        jq(`#${field_id}`).prop('checked', (pairs[key] == 0 ? false : true))
      }

      // span e div sao usados qdo form read only
      // span é usado para campos pequenos e div é usada para textos grandes, descritivos
      if (tagType == 'span') jq(`#${field_id}`).html(pairs[key])
      if (tagType == 'div') jq(`#${field_id}`).html(decodeHTMLEntities(pairs[key].replace(/\n/g, '<br/>')))   // para textarea readonly é usado '<div>'
    }
  })

  $imagesStillLoading = false   // avisa que carregou eventual <img> do formulario atual

}




/********************************************************************************************************************************************************
 exclui o registro atual sendo exibido no form edicao
********************************************************************************************************************************************************/
const performDeleteRecord = () =>  {
  // exemplo:  cars/delete/41  ,  manufactures/delete/87
  performCrudOperation(null, 'DELETE', `${tableName}/delete/${currentFormRecordId}`)  
}


/****************************************************************************************************
 apagar registros selecionados
****************************************************************************************************/
const deleteSelectedRecords = (action) => {

if (selectedRecords.length==0)  return;


// usuario clicou no botao excluir registros selecionado (topo da datatable), abre janela de confirmacao
if (action=='ask') {
  showDeleteSelectedRecordsConfirmation=true
}

else if (action=='delete')  { 
  closeCurrentModalWindow()

  let ids = selectedRecords.toString();      // converte o array selectedRecords em string para que possa ser inserido no mysql

  // exemplo:  cars/batch_delete/41  ,  manufactures/batch_delete/87
  performCrudOperation(null, 'DELETE', `${tableName}/batch_delete/${ids}`)  
}
}

/****************************************************************************************************
 abre form para novo registro
****************************************************************************************************/
const newRecord = () => {
  setTimeout(() => {showLoadingGif() }, 1);

  formHttpMethodApply = 'POST'
  calledEditForm=true
}

/****************************************************************************************************
 exibe somente registors ativos ou todos
****************************************************************************************************/
const showActiveOrInactive = (which) => {

  // se ja esta ativo, desativa 
  if (showOnlyActiveOrInactiveRecords == which) showOnlyActiveOrInactiveRecords='' 

  // nao esta ativo, ativa 
  else 
    showOnlyActiveOrInactiveRecords = which

  forceDatatableReload++
}

/****************************************************************************************************
 transform um registro em ativo ou inativo, dependendo do seu status atual
****************************************************************************************************/
const turnActiveOrInactive = (_id_) => {

  // 'formUse'= car, language, etc  =>  api url parecida com:   /car/status/43   , /manufacturer/status/21,   etc
  performCrudOperation(null, 'PATCH', `${tableName}/status/${_id_}`)   
}

/****************************************************************************************************
 abre form para edicao de registro
****************************************************************************************************/
const editRecord = (event, _id_) =>  {
  event.stopPropagation();
  formHttpMethodApply = 'PATCH'
  currentFormRecordId = _id_

  setTimeout(() => {showLoadingGif() }, 1);
  calledEditForm=true
}

/****************************************************************************************************
 exclui determinado registro
****************************************************************************************************/
const deleteRecord = (event, _id_) => {
  event.stopPropagation();
  formHttpMethodApply = 'DELETE'
  currentFormRecordId = _id_

  calledEditForm=true
}


/****************************************************************************************************
define os campos que serao populados pelo componente superior (+page.svelte)
****************************************************************************************************/
export let fields, apiURL, orderBy, orderDirection, tableName, tableTitle, metadataTableName, recordsReadonly

/****************************************************************************************************
todos os campos acima sao 'string',  a excecao é o campo 'fields', que deve ser um objeto com a mesma 
estrutura do exemplo abaixo: field_name, field_title, etc
****************************************************************************************************

[ {field_name: 'name', field_title: $Terms.fieldname_car_name,  field_width:'w-[15%]'}, 
  {field_name: 'manufacturer_name', field_title: $Terms.fieldname_car_manufacturer, field_width:'w-[85%]'}      ]  


*****************************************************************************************************/

// legenda no canto inferior da datatable, varia de acordo com país selecionado
let baseboardMessage = ''

let divTBODY  // div que contem os registros (DIV's) da tabela

// contem o eventual texto que sera pesquisado caso usuario preencha algo em .txtTABLE_SEARCHBOX e pressione Enter 
let textToSearch


/************************************************************************************************************************************************************
 marca ou desmarca registros da datatable, é acionado pelo icone no topo superior esquerdo (marcar/desmarcar registros)
************************************************************************************************************************************************************/

const uncheckOrCheckAllRecords = () => {

  let selectedRecords_tmp = []  
 
  // percorre todos os registros e marca ou desmaca seu checkbox dependendo se ha regs marcados ja 
  jq('.divTBODY').find("input[type=checkbox]").each(function () {  

    let divTR_ID = jq(this).attr('id').match(/\d/g).join("")  // remove letras do ID, obtem só o numero da divTR

    if ( selectedRecords.length >0 )  jq(this).prop("checked", false) 
    if ( selectedRecords.length==0 )  {
      jq(this).prop("checked", true) 
      selectedRecords_tmp.push( divTR_ID )
    }

  })

  // atualiza mensagem de exclusao em lote para singular ou plural
  deleteMessageWarning = selectedRecords.length==1 ? $Terms.warning_delete_one_selected_record_window : $Terms.warning_delete_many_selected_record_window
  selectedRecords = selectedRecords_tmp
}



/************************************************************************************************************************************************************
 funcao que exibe seta para baixo/para cima na coluna atualmente selecionada (para organizacao) e baseado na direcao da organizacao (asc/desc)
************************************************************************************************************************************************************/
const tableReorder = (field_name) => {
  if (orderBy == field_name)  orderDirection = orderDirection=='asc' ? 'desc' : 'asc'
  else {orderBy = field_name; orderDirection='asc'}
  forceDatatableReload++
}

/************************************************************************************************************************************************************
 funcao executada qdo usuario clica em uma divTR (uma TR da table, digamos assim)
  a variavel 'selected_recods' contem o ID dos registros selecionados
************************************************************************************************************************************************************/
const divTRClicked = (event) => {

  // ID da divTR clicada
  let divTRClicked = jq('#'+event.currentTarget.id)
  let divTRClicked_count_id = event.currentTarget.id.match(/\d/g).join("")   // remove letras obtem so o ID (numero) do registro

  // o checkbox que acompanha cada divTR sera marcado/desmarcado dependnedo do seu estado atual
  jq('#chkTR_SELECTOR'+divTRClicked_count_id).prop('checked', ! jq('#chkTR_SELECTOR'+divTRClicked_count_id).prop('checked'))

  // verifica se item clicado ja estava clicado
  let clicked_already_selected = selectedRecords.indexOf(divTRClicked_count_id)

  // necessario criar uma variabel temporaria (selectedRecords_tmp) trabalha la, e depois atualizar na variavel correta (selectedRecords)
  // pq se nao fizer isso, nao havera reatividade nos botoes (excluir regs, etc) que se baseiam na qtde regs selecionados

  // o comando 'selectedRecords.push'  nao dispara a reatividade do svelte

  let selectedRecords_tmp = selectedRecords

  // remove da lista de selecionados, ID do registro clicado que ja estava selecionado
  if (clicked_already_selected!=-1)
      selectedRecords_tmp = selectedRecords_tmp.filter(function(e) { return e !== divTRClicked_count_id })
    
  // remove destaque da divTR selecionada por ultimo
  jq('.divTR_SELECTED').removeClass('divTR_SELECTED').addClass('divTR')

  divTRClicked.removeClass('divTR').addClass('divTR_SELECTED')  

  if (clicked_already_selected==-1) {
    selectedRecords_tmp.push(divTRClicked_count_id)
  }

  selectedRecords = selectedRecords_tmp

  // atualiza mensagem alerta exclusao em lote para singular ou plural  
  deleteMessageWarning = selectedRecords.length==1 ? $Terms.warning_delete_one_selected_record_window : $Terms.warning_delete_many_selected_record_window

  setTimeout( () => {showOrHideTooltip();}, 100)   // refaz/destroi tooltip de botoes dependendo qtde regs selecionados
}


/************************************************************************************************************************************************************
 busca dados via API e monta tabela com dados recebidos
************************************************************************************************************************************************************/
const mountTable = async () =>  {

  // se carregando a datatble pela 1a vez, reseta campo pesquisa, que pode ter vidno preenchdio de outra tela
  if (! datatableReady)  jq('.txtTABLE_SEARCHBOX').val('')
  datatableReady = true

  textToSearch=''
  //if (loaded_first_time)  jq('.txtTABLE_SEARCHBOX').val('')  txtTableSearchText
  selectedRecords = []  

  // contem os json com os registros que serao exibidos
  dataList = null

  // necessario abrir evento assincrono para exibir div ajax loading, caso contrario navegador nao atualiza a tela
  setTimeout(() => {showLoadingGif();}, 1);


  // pelo fato da legenda que fica na parte inferior da databable conter elemento html a frase sera definida em codigo, e nao buscada na tabela 'terms'
  if (selectedCountry()=='usa') 
    baseboardMessage = `Legend:&nbsp;&nbsp;&nbsp;<div style="width:10px;height:15px;margin-top:-3px;background-color:red;">&nbsp;</span></div><div>=&nbsp;&nbsp;disabled</div>`
  else   
    baseboardMessage = `Legenda:&nbsp;&nbsp;&nbsp;<div style="width:10px;height:15px;margin-top:-3px;background-color:red;">&nbsp;</span></div><div>=&nbsp;&nbsp;inativo</div>`

 
  // so considera campo para pesquisa se tiver minimo 3 caracteres
  textToSearch = ''
  if (jq.trim( jq('.txtTABLE_SEARCHBOX').val()).length>=3 ) {
    textToSearch = jq.trim( jq('.txtTABLE_SEARCHBOX').val() )
  }

  // chama API que retornará lista de registros
  // apiURL foi definida pelo componente que chamou 'Datatable.svelte'

  // +page.svelte teve que enviar apiURL com '@group_info porque page.svelte roda no servidor e nao tem acesso à: localStorage.getItem("rentacar_workgroup_name");

  /*  
  apURl exemplos:  
      http://localhost:5173/Admin/manufacturers?country=usa&order_by=name&order_direction=asc&search_txt=&only_active_or_inactive_records=
      http://localhost:5173/Discreet/cars?country=usa&order_by=name&order_direction=asc&search_txt=&only_active_or_inactive_records=
      http://localhost:5173/Mirth/terms?country=brazil&order_by=name&order_direction=asc&search_txt=&only_active_or_inactive_records=
  */

  _currentWorkgroupName = localStorage.getItem("rentacar_workgroup_name");

  let _apiURL = apiURL.replaceAll('@group_info', `${_currentWorkgroupName}`)

  const url = new URL( _apiURL  );
  url.searchParams.set('country', selectedCountry());
  url.searchParams.set('order_by', orderBy);
  url.searchParams.set('order_direction', orderDirection);
  url.searchParams.set('search_txt', textToSearch);
  url.searchParams.set('only_active_or_inactive_records', showOnlyActiveOrInactiveRecords);

  logAPI('GET', url)
  await fetch(url, {method: "GET"} )

  .then((response) => {

    if (!response.ok) {
      setTimeout( () => {hideLoadingGif();}, 100) 
      slidingMessage('slidingWindowMessage', $Terms.errormsg_general_database, 4000)    
      return
    }
    return response.json();
  })
  
  .then((tableData) => {

    setTimeout( () => {hideLoadingGif();}, 500)   // esconde div 'processando'

    // melhora o css/visual do 'title' dos botoes  
    setTimeout(() => { showOrHideTooltip() }, 100);

    dataList = tableData    

    // se a URL retornou algum registro 
    if (dataList!=null)   {

      // necesario esquecer qual a divTR foi selecionada por ultimo
      jq('.divTR_SELECTED').removeClass('divTR_SELECTED').addClass('divTR')

      // posiciona DIV no topo apos atualizacao dos dados
      setTimeout(() => {
        divTBODY.scroll({ top: 0, behavior: 'auto' })
      }, 100);


      // qdo browser exibe scrollbar vertical na div de registros (divTBODY), as colunas da linha que contem os titulos 
      // nao ficam paralelas com as colunas de dados, devido à colocacao do scroll no canto direito 
      // o codigo abaixo aumenta a div de dados, para que as colunas TH voltem a ficar paralelas com as colunas de dados
      setTimeout(() => {
          
          // iguala a linha TH com as linhas de dados
          let width = parseInt( jq('.divTBODY').css('width').replace('px',''), 10 );
          jq('.divTH_CONTAINER').css('width', width +'px');

          // se browser colocou scrollbar vertical devido a qtde de registros , corrige a distorcao de larguras
          // existe uma margem de seguranca de 5 pixels para o calculo
          if ( jq('.divTBODY').get(0).scrollHeight > jq('.divTBODY').innerHeight()+5 )  
            jq('.divTH_CONTAINER').css('width',width - 18 +'px');


      }, 100);

      setTimeout( () => {showOrHideTooltip();}, 100)   // refaz/destroi tooltip de botoes dependendo qtde regs selecionados

      // se apos update da datatable e,
      // se houve realmente uma gravacao, tenta rolar a div container do datatable ate o registro recem alterado/inserido e tenta coloca lo em destaque
      if (lastRecordIdUpdated!=0)    {
        setTimeout( () => {
            scrollUntilElementVisible('divTR'+lastRecordIdUpdated)  
            jq('#divTR'+lastRecordIdUpdated).removeClass('divTR').addClass('divTR_SELECTED')  
            lastRecordIdUpdated = 0 
        }, 500)   
      }
      return dataList
    } 

    setTimeout(() => {
      jq('.txtTABLE_SEARCHBOX').val('')
    }, 1);


    
  })
}

/************************************************************************************************************************************************************
 se por exemplo, usuario redimensionar a tela (page.svelte.windowResize()), a datatable ficara bagunçada, a funcao abaixo qdo acionada força sua atualizacao
************************************************************************************************************************************************************/
export function forceDatatableRefresh() {
  forceDatatableReload++
}

/************************************************************************************************************************************************************
funcao abaixo é acionada qdo usuario pressionou Enter com o campo de pesquisa com foco (page.svelte.onKeyDown())
************************************************************************************************************************************************************/
export function performDatatableSearch() {
  if ( jq.trim(jq('.txtTABLE_SEARCHBOX').val()).length < 3 )  {
    slidingMessage('slidingWindowMessage', $Terms.msg_datatable_minimum_length, 2000)   
    return
  }
  forceDatatableReload++
}


/*****************************************************************************************************************************************************************
// funcao acionada qdo usuario clica na div '#backDrop' ou clica em um dos botoes 'Fechar Form', e com isso fecha (sem gravar) o form de edicao que estiver aberto 
// funcao abaixo serve tanto para form quanto janelas modais de alerta, confirmacao
*****************************************************************************************************************************************************************/
export const closeCurrentModalWindow = () => {
  calledEditForm = false    // fecha form edicao
  showDeleteSelectedRecordsConfirmation = false   // fecha tela confirmacao exclusao
}

/*****************************************************************************************************************************************************************
// fehca janela modal com metadados do registro
*****************************************************************************************************************************************************************/
export const closeMetadataWindow = () => {
  showRecordMetadata = false 
}


/************************************************************************************************************************************************************
// funcao acionada após gravacao do registro no backend
************************************************************************************************************************************************************/
export const closeCurrentModalWindowAndRefreshDatatable = (id) => {

  closeCurrentModalWindow();
  // no caso de exclusao registro, nao sera passado um ID de reg para buscar
  // no caso de edicao registro, sera passado um ID de reg
  //if (typeof event.detail.record_id!=0)   lastRecordIdUpdated = event.detail.record_id; 
  lastRecordIdUpdated = id; 

  forceDatatableReload++;
}  

/************************************************************************************************************************************************************
// funcao acionada qdo necessario atualizar datatable
************************************************************************************************************************************************************/
export const refreshDatatable = () => {
  lastRecordIdUpdated = currentFormRecordId; 
  forceDatatableReload++;
}



/**********************************************************************************************
a funcao abaixo insere, edita (update), apaga (delete), muda o status 'ativo/inativo'
de qualquer  registro/tabela usado na aplicacao
**********************************************************************************************/

const performCrudOperation = async (body, formHttpMethodApply, route ) => {

  // mostra animacao indicando pausa para processamento
  setTimeout(() => {showLoadingGif();}, 1);

  let successmsg
  let errormsg

  let record_id_put_focus = ''

  let changingStatus = formHttpMethodApply=='PATCH' && route.indexOf('/status/')!=-1     // rota para mudar status do registro (ativo/inativo)
  let batchDeleting = formHttpMethodApply=='DELETE' && route.indexOf('/batch_delete/')!=-1    // rota para apagar N registros

  // prepara as msgs de erro ou sucesso, para cada tipo de operacao 
  switch (formHttpMethodApply) {

    case 'GET':
      // se leitura ocorrer ok, nao precisa mostrar aviso, nao existe 'Terms.successmsg_record_read'
      errormsg = $Terms.errormsg_reading_record
      break

    case 'POST':
      successmsg= $Terms.successmsg_record_inserted
      errormsg = $Terms.errormsg_saving_record
      break

    case 'PATCH':
      if ( changingStatus ) { 
        successmsg= $Terms.success_msg_status_changed
        errormsg = $Terms.error_msg_status_changed

      } else {   // edicao de registro
        successmsg= $Terms.successmsg_record_updated
        errormsg = $Terms.errormsg_saving_record
      }
      break

    case 'DELETE':
      if ( batchDeleting ) { 
        successmsg= $Terms.successmsg_selected_record_deleted
        errormsg = $Terms.errormsg_deleting_selected_record

      } else {  //    exclusao de 1 registro
        successmsg= $Terms.successmsg_record_deleted
        errormsg = $Terms.errormsg_deleting_record
      }
      break
  }

  try {
    const url = new URL(`${$backendUrl}/${route}`);

    let options

    // metodo 'GET' nao possui body, parametros ja vieram na URL (route)
    if (formHttpMethodApply=='GET')      
      options = {method: 'GET'}

    // metodo DELETE e metodo PATCH com rota= '/status/' nao possuem body, os parametros sao passados via URL  
    // metodos POST e PATCH possuem body (form-data com dados para gravacao)
    else {
      if ( body != null )  {
        // infelizmente o PHP, ate a versao 8.2, nao trabalha bem com metodo PATCH, necessario mudar para 'POST'
        // a funcao php, request_parse_body(), que promete ler body qdo metodo= PATCH() nao funciona na versao PHP do backend (7.4)
        // request_parse_body() so funciona na versao 8.4

        // os demais backends manipulam PATCH sem problemas

        let backend = localStorage.getItem("rentacar_current_backend")
        if (backend=='php' && formHttpMethodApply=='PATCH')   
          options = {method: 'POST', body}
        else
          options = {method: formHttpMethodApply, body}
      }

      else
        // metodo PATCH com parametros passados na URL mesmo
        options = {method: formHttpMethodApply}     
    }
    
    // chama backend para GET ou POST
    //logAPI( `<div style='width:30px'>${formHttpMethodApply}</div>${url}` )
    logAPI( formHttpMethodApply, url )

    await fetch(url, options)

    .then((response) => {
      if (formHttpMethodApply=='GET') return response.json()   // read= retorno é json contento campos do registro
      else return response.text()     // delete, patch, post = retorno é um texto informando se operacao teve sucesso ou nao
    })                      

    // se executou backend com sucesso
    .then((result) => {

      setTimeout(() => {
        hideLoadingGif(); 
      },  1);

      // sucesso ao ler dados do registro
      if (formHttpMethodApply=='GET') { 
          let pairs = result
          // preenche os campos do formulario com o conteudo JSON recebido
          populateFormDestinationFields(pairs)

          // alguns forms de edicao exibem a bandeira do país atualmente selecionado atraves da div 'divCountryFlag'
          jq('#divCountryFlag').html( svgCurrentCountryFlag() )

          putFocusInFirstInputText_AndOthersParticularitiesOfTheForm();          
      }

      // != GET= sucesso ao gravar dados do registro (update, insert, delete, delete_selected)
      else {
        let text = result

console.log('erro= '+result)

        // se gravou ok
        if (text.indexOf('__success__')!=-1)  {

          slidingMessage('slidingWindowMessage', successmsg, 2000)   

          record_id_put_focus = ''    

          if ( formHttpMethodApply=='POST' )  {
              // qdo formHttpMethodApply= 'POST', backend retorna o ID do registro recem inserido na tabela, o texto retornado= __success__|new_id
              // por exemplo:  __success__|543  

              let new_record_id = ''
              if (formHttpMethodApply=='POST')  {
                let details = text.split('|')
                new_record_id = details[1]
              }

              record_id_put_focus = new_record_id    // result=>   ok|new_record_id
          }

          if ( formHttpMethodApply=='PATCH')   {

            // ultimo parametro na URL é o ID cuja registro foi alterado, exemplo:  /car/status/54 , /manufacturer/31,  etc
            let params = route.split('/') 
            record_id_put_focus = params[ params.length-1 ]   

          }

          // fecha form edicao e avisa para colocar foco , destacar o registro recem inserido ou recem alterado 
          closeCurrentModalWindowAndRefreshDatatable(record_id_put_focus)

          // se usuario alterou algum dado que necessita uma funcao a mais apos sua alteracao, executa esta funcao
          // exemplo: usuario alterou alguma expressao ingles/portugues, necessario recarregar todas as expressoes para que svelte atualize a expressao alterada na tela
          if (toDoAfterPatchOrPost != '')  eval(toDoAfterPatchOrPost)

        // backend retornou problema
        } else  {
          // se formulario nem foi aberto, usa a div maior (window, slidingWindowMessage) para informar erro
          // se formulario aberto, usa a div interna do formulario (slidingFormMessage)

          let divSliding = (formHttpMethodApply=='GET' || changingStatus || batchDeleting) ? 'slidingWindowMessage' : 'slidingFormMessage'

          slidingMessage(divSliding, errormsg + '&nbsp;&nbsp;=> '+text, 4000)   
          return('error')
        }
      }
    })

    // se retornou erro, erro execucao da API
    .catch(function (error)   { 

      setTimeout(() => {hideLoadingGif();}, 1); 

      // se GET, fecha form pois nao foi possivel ler campos
      // GET, changingStatus e batchDeleting usa div rolante na janela principal

      if (formHttpMethodApply=='GET' || changingStatus || batchDeleting) { 
        setTimeout(() => {closeCurrentModalWindow()}, 5500); 
        slidingMessage('slidingWindowMessage', $Terms.errormsg_general_database + '&nbsp;&nbsp;Error position 2= '+error.message, 5000)   
      }
      // erro gravacao/exclusao, avisa dentro do proprio form o erro e mantem form aberto
      // usa div rolante no form aberto para informar erro
      else {
        slidingMessage('slidingFormMessage', errormsg + '&nbsp;&nbsp;Error: situation 2= '+error.message, 5000)   
        return('error')
      }

    });


  // se erro conexao, erro jscript
  } catch (error) {

    setTimeout(() => {hideLoadingGif();}, 1);

    // se GET, fecha form pois nao foi possivel ler campos
    // GET, changingStatus e batchDeleting usa div rolante na janela principal

    if (formHttpMethodApply=='GET' || changingStatus || batchDeleting) { 
      setTimeout(() => {closeCurrentModalWindow()}, 5500); 
      slidingMessage('slidingWindowMessage', $Terms.errormsg_general_database + '&nbsp;&nbsp;Error position 3= '+error.message, 5000)
    }
    // erro gravacao/exclusao, avisa dentro do proprio form o erro e mantem form aberto
    // usa div rolante no form aberto para informar erro
    else {
      slidingMessage('slidingFormMessage', errormsg + '&nbsp;&nbsp;Error situation 3= '+error.message+'-'+formHttpMethodApply, 5000)    
      return('error')
    }
  }

}

/************************************************************************************************************************************************************
coloca foco no 1o input text da div que contem os campos digitaveis
e executa algumas acoes especificas de cada formulario
************************************************************************************************************************************************************/
export const putFocusInFirstInputText_AndOthersParticularitiesOfTheForm = () => {

  makeWindowDraggable('divWINDOW_TOP', 'recordForm')

  // tableName é uma variavel do tipo 'export' que foi preenchida ao chamar 'Datatable.svelte'
  // form de carro
  if (tableName=='cars' ) {
    jq('#txtRentalPrice').maskMoney({showSymbol:true, symbol:'R$ ', thousands:'.', decimal:',', symbolStay: true});   

    prepareAutocomplete('manufacturers', $Terms.error_preparing_autocomplete, 'txtManufacturerName') 
  } 


  setTimeout(() => {
    let firstInputText = jq('#recordForm').find('input[type=text],textarea,select').filter(':visible:first');
    firstInputText.focus();  
    firstInputText.select();          
    showOrHideTooltip();
    hideLoadingGif()   // esconde div 'carregando...'  caso nao tenha sido escondida ainda
  }, 300);

}
/************************************************************************************************************************************************************
 se alguma das variaveis abaixo for alterada, dispara 'mountTable' e recarrega datatable
 a linha abaixo tb faz disparar a funcao 'mountTable' na inicializacao do componente 'Datatable.svelte',  
  o svelte considera que estas variaveis foram alteradas qdo componete é iniciado -  por isso nao ha necessidade de criar 'onMount'
  alias,  se for criada funcao onMount, sera carregada 2x ao iniciar o componente
************************************************************************************************************************************************************/

$: $Terms, forceDatatableReload, $recordIdsVisible,  mountTable()

// disponibiliza a outros componentes funcoes como: 'performCrudOperation', 'putFocusInFirstInputText_AndOthersParticularitiesOfTheForm', etc etc
// todos os forms de edicao (language, car, manufacture, etc) usam as funcoes abaixo
setContext('offerCrudOperationsToChildComponent', { 
  performCrudOperation,  
  performSaveRecord,
  putFocusInFirstInputText_AndOthersParticularitiesOfTheForm, 
  performDeleteRecord, 
  fillRecordMetadataWindow,
  invokeMetadataWindow,
  getFormPopulatedAndReady,
});


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

<!-- o evento onMount do svelte é disparado 2x qdo um componente é chamado, a unica solucao é carregar a funcao inicial usando <body> -->


<!-- usuario pediu para editar ou inserir registro -->
{#if calledEditForm}      

  <!-- backDrop é uma div fumê que ficara por trás do formulario de edicao 
  backDrop impede que o usuario clique em elementos que estão atrás do formulario de edicao   
  qdo backDrop for clicada, fecha o form edicao -->

  <div id='backDrop' class='w-full h-full  absolute flex items-center justify-center left-0 top-0 z-10 bg-[rgba(0,0,0,0.5)]' on:click|self={closeCurrentModalWindow} aria-hidden="true"  >  

      <!-- formHttpMethodApply pode ser POST, PATCH ou DELETE  --> 
      {#if (tableName=='terms')} 
        <TermForm  on:dispatchCloseForm={closeCurrentModalWindow}  {formHttpMethodApply}  />

      {:else if (tableName=='cars')}
        <CarForm  on:dispatchCloseForm={closeCurrentModalWindow} {formHttpMethodApply}  />

      {:else if (tableName=='manufacturers')}
        <ManufacturerForm  on:dispatchCloseForm={closeCurrentModalWindow}   {formHttpMethodApply}  />


      {/if}

  </div>

{/if}      



<!-- usuario pediu para exibir metadados do registro -->
{#if showRecordMetadata}      

  <!-- backDropMetadata é a div fumê que ficara por trás da janela modal de  metadados
  backDropMetadata impede que o usuario clique em elementos que estão atrás da janela modal
  qdo backDrop for clicada, fecha o form edicao -->

  <div id='backDropMetadata' class='w-full h-full  absolute flex items-center justify-center left-0 top-0 z-10 bg-[rgba(0,0,0,0.5)]' on:click|self={closeMetadataWindow} aria-hidden="true"  >  

      <RecordMetadataWindow  on:dispatchCloseForm={closeMetadataWindow}  {currentFormRecordId}  {metadataTableName}  />   

  </div>

{/if}      


<!-- usuario pediu para excluir registros selecionados -->
{#if showDeleteSelectedRecordsConfirmation}      

  <!-- backDrop é a div fumê que ficara por tras da janela modal de confirmacao exclusao
  backDrop impede que o usuario clique em elementos que estão atrás da janela
  qdo backDrop for clicada, fecha a janela de confirmacao abaixo -->

  <div id='backDrop' class='w-full h-full  absolute flex items-center justify-center left-0 top-0 z-10 bg-[rgba(0,0,0,0.4)]' on:click|self={closeCurrentModalWindow} aria-hidden="true"  >  

    <ConfirmationWindow  on:dispatchCloseWindow={closeCurrentModalWindow} on:dispatchDeleteSelected={ () => deleteSelectedRecords('delete')} 
        button_confirm_message = {$Terms.datatable_delete_selected_records}     html_message={ deleteMessageWarning }  />   

  </div>


{/if}      




<div class="divTABLE_CONTAINER font-Roboto" > 

  <!--  titulo da tabela -->
  <div class="divTABLE_TITLE" >  
    {tableTitle}
  </div>

  <!-- icones para manipulacao de registros, parte superior da datatable -->
  <div class="divTABLE_ICONS" >  

      <div class="flex flex-row " >  
          <!-- botao marcar/desmarcar registros selecionados -->
          {#if ! recordsReadonly}
              <div id='btnCheckUncheckRecords' class='containsTooltip' class:btnTABLE_UNCHECK_RECORDS={selectedRecords.length!=0}  _originalTooltip_={ $Terms.datatable_checkbox_all_records }
                  on:click={uncheckOrCheckAllRecords} aria-hidden="true" class:btnTABLE_CHECK_RECORDS={selectedRecords.length==0}  >
              </div> 
          {/if}

          <!-- espacamento horizontal -->
          <div class="w-7" >&nbsp;</div>  

          <!--  textbox para pesquisa -->
          <div class="flex flex-col " >  
            <!-- nao posso vincular o textbox abaixo com a variavel 'textToSearch' porque cada letra digitada vai acionar a pesquisa -->
            <input type="text" class='txtTABLE_SEARCHBOX'  id='txtTableSearchText'  autocomplete="off" 
                on:focus={() => {showTipSearchbox=true}} on:blur={() => {showTipSearchbox=false}} on:mouseenter={ focusDatatableSearchBox } >

            <div class="flex flex-row pt-1 {showTipSearchbox ? 'visible' : 'invisible'} " >  
                <span class="text-blue-900 font-bold">Enter</span><span class="text-black">= {$Terms.verb_search}</span>
            </div>
          </div>

          <!-- botao para cancelar filtro por texto e cancelar filtro por ativo/inativo -->
          <div id='btnResetTextTableFilter' class:containsTooltip={textToSearch!='' || showOnlyActiveOrInactiveRecords!=''} 
              _originalTooltip_={ $Terms.datatable_cancel_filter }  class:ToRemoveTooltip={textToSearch=='' &&  showOnlyActiveOrInactiveRecords==''}
              class:btnTABLE_CANCEL_FILTER_INACTIVE={textToSearch=='' && showOnlyActiveOrInactiveRecords==''} 
              class:btnTABLE_CANCEL_FILTER_ACTIVE={textToSearch!='' || showOnlyActiveOrInactiveRecords!=''} 
              on:click={ () => {resetDatatableFilter()} }  aria-hidden="true">
          </div> 
      </div>


      <div class="flex flex-row " >  

          <!-- botao somente inativos -->
          <div  class='containsTooltip' class:btnTABLE_ONLY_INACTIVE_RECORDS_ON={showOnlyActiveOrInactiveRecords=='inactive'} 
              class:btnTABLE_ONLY_INACTIVE_RECORDS_OFF={showOnlyActiveOrInactiveRecords!='inactive'}
            _originalTooltip_={$Terms.only_inactive_records}   on:click={ () => {showActiveOrInactive('inactive')} } aria-hidden="true">&nbsp;
          </div> 

          <!-- botao somente ativos -->
          <div class='containsTooltip' class:btnTABLE_ONLY_ACTIVE_RECORDS_ON={showOnlyActiveOrInactiveRecords=='active'} 
              class:btnTABLE_ONLY_ACTIVE_RECORDS_OFF={showOnlyActiveOrInactiveRecords!='active'}
            _originalTooltip_={$Terms.only_active_records}   on:click={ () => {showActiveOrInactive('active')}  } aria-hidden="true">&nbsp;
          </div> 

          <!-- espacamento horizontal -->
          <div class="w-7" >&nbsp;</div>  


          <!-- a tabela workgroup é reaonly, sao dados manipulados so pela aplicacao, nao deve permitir edicao -->
          {#if ! recordsReadonly}

              <!-- botao novo registro -->
              <div  class='containsTooltip btnTABLE_NEW_RECORD' _originalTooltip_={ $Terms.datatable_new_record }   on:click={newRecord} aria-hidden="true"></div> 

              <!-- espacamento horizontal -->
              <div class="w-7" >&nbsp;</div>  

              <!-- botao excluir registros selecionados -->
              <div id='btnDeleteSelectedRecords' class:containsTooltip={selectedRecords.length>0} _originalTooltip_={ $Terms.title_delete_selected_records }  
                  class:ToRemoveTooltip={selectedRecords.length==0}  class:btnTABLE_DELETE_RECORDS_ACTIVE={selectedRecords.length>0} 
                  class:btnTABLE_DELETE_RECORDS_INACTIVE={selectedRecords.length==0} 
                  on:click={ () => {deleteSelectedRecords('ask')}} aria-hidden="true"> 
              </div> 
          {/if}
      </div>


  </div>
 


  <!-- TITULOS DAS COLUNAS

  - se clicar no titulo da coluna, altera ordem organizacao e recarrega tabela -->

  <div class="divTH_CONTAINER"  > 

    <!-- titulo coluna 'checkbox', checkbox selecao -->
    {#if ! recordsReadonly}
      <div class="divTH basis-[50px] grow-0 shrink-0 " >&nbsp;</div>
    {/if}

    {#each fields as {field_name, field_title, field_width}} 

      {#if field_name == orderBy}
        {#if orderDirection == 'asc'}
          <div on:click={ tableReorder(field_name) } class="divTH_SELECTED_ASC {field_width} " aria-hidden="true" >{field_title}</div> 
        {:else}  
          <div on:click={ tableReorder(field_name) } class="divTH_SELECTED_DESC {field_width}" aria-hidden="true">{field_title}</div> 
        {/if}

      {:else}
          <div on:click={ tableReorder(field_name) } class="divTH {field_width} " aria-hidden="true">{field_title}</div> 
      {/if}
    {/each}

    <!-- titulo coluna 'acoes', checkbox selecao, edit e delete -->
    <!-- a tabela workgroup é readonly, sao dados manipulados so pela aplicacao, nao deve permitir edicao -->
    {#if ! recordsReadonly}
        <div class="divTH basis-[150px] grow-0 shrink-0" >&nbsp;</div>
    {/if}  

  </div> 

  <!-- CORPO DA TABELA (LINHAS) -->

  <!-- se ja recebeu os dados da funcao assincrona (mountTable), exibe --> 

  <div class="divTBODY" bind:this={divTBODY}  > 

      <!-- nao recebeu dados ainda, exibe 'carregando...'    -->
      {#if ! dataList}
          <div class="divLOADING_MSG">{$Terms.loading}</div>

      <!-- recebeu dados, monta datatable   -->
      {:else}


          <!-- loop atraves das linhas do datalist -->
          {#each dataList as record, record_count}

            <!-- datatable de fabricantes ou carros, altura da linha é maior do que as demais, é o padrao -->

            <!-- tableName => cars, manufactures, etc etc -->

            <div class="divTR" class:h-16={ tablesWithBiggerRow.includes(tableName) }    
          class:h-9={! tablesWithBiggerRow.includes(tableName)} on:click={divTRClicked} aria-hidden="true" id="divTR{record.id}" >  

              <!-- checkbox que seleciona a divTR --> 
              {#if ! recordsReadonly}
                <div class="divTD_ACTION" class:h-16={ tablesWithBiggerRow.includes(tableName) }  class:h-9={ ! tablesWithBiggerRow.includes(tableName) }   >
                  <input type="checkbox" id="chkTR_SELECTOR{record.id}" checked={jQuery.inArray(record.id.toString(), selectedRecords)!=-1} class='w-full h-3 '>
                </div>
              {/if}

              <!-- 'fields' contem os campos que sera exibidos (dentre todos os campos da tabela) e a largura de suas colunas -->
              {#each fields as field}

                    <!-- colunas que são imagens, ao inves de texto simples: 
                        manufacturer_logo (logotipo do fabricante), 
                        car_image (imagem do veiculo) 
                    -->
                    {#if field.field_name=='car_image' || field.field_name=='manufacturer_logo'  }

                      <div class="divTD_NOT_CLICKABLE  {field.field_width}">
                        <img src="{ $imagesUrl + dataList[record_count][field.field_name] + '?'+forceImgRefresh()}" class="w-[60px] h-full" alt="img" />
                      </div> 

                    <!-- 
                        colunas texto normal 

                      registros desativados (dataList[record_count]['active']==false) aparecerão em vermelho
                    -->
                    {:else}

                        <!-- se usuario pediu para ver IDs dos registros ($recordIdsVisible) e a coluna atual é a coluna principal (fieldConcatenateWithId) -->

                        {#if fieldConcatenateWithId == field.field_name && $recordIdsVisible}

                          <div class="divTD_NOT_CLICKABLE  {field.field_width}" 
                                class:line-through={dataList[record_count]['active']==false}
                                class:text-[red]={dataList[record_count]['active']==false} >
                              <span class='min-w-12 w-12' 
                                class:text-red-300={dataList[record_count]['active']==false}  
                                class:text-gray-400={dataList[record_count]['active']==true} >{record.id}</span>: 
                              { dataList[record_count][field.field_name] }
                          </div>

                        <!-- exibicao normal, sem ID do registro -->

                        {:else}
                            <!-- se o campo nao é boolean, exibe seu conteudo  -->
                            {#if ! field.is_boolean  }
                              <div class="divTD_NOT_CLICKABLE  {field.field_width}" 
                                  class:text-[red]={dataList[record_count]['active']==false} 
                                  class:line-through={dataList[record_count]['active']==false} >
                                  { dataList[record_count][field.field_name] }
                              </div>

                            <!-- se o campo é boolean, converte para sim/nao,  yes/no  -->
                            {:else} 

                              {#if dataList[record_count][field.field_name] == true }
                                <div class="divTD_NOT_CLICKABLE  {field.field_width}" class:text-[red]={dataList[record_count]['active']==false}>{ $Terms.boolean_true }</div>
                              {:else}
                                <!--
                                <div class="divTD_NOT_CLICKABLE  {field.field_width}">{ $Terms.boolean_false }</div>

                                usar traco (-) ao inves de 'não', para facilitar leitura da coluna
                                -->
                                <div class="divTD_NOT_CLICKABLE  {field.field_width} text-[red]">-</div>
                              {/if}

                            {/if}

                        {/if}

                    {/if}

              {/each}

              <!-- botoes acao= editar, excluir registro  --> 

              <!-- a tabela workgroup é readonly, sao dados manipulados so pela aplicacao, nao deve permitir edicao -->
              {#if ! recordsReadonly}
                  <div class="basis-[150px] grow-0 flex flex-row " >  

                      <!-- se o registro inativo, desativa os botoes de acao, deixa so o botao 'ativar' funcionando -->
                      {#if dataList[record_count]['active']==0 }

                          <div class="divTD_ACTION" >&nbsp;</div>
                          <div class="divTD_ACTION" >&nbsp;</div>

                          <div class="divTD_ACTION" class:h-16={ tablesWithBiggerRow.includes(tableName) }  class:h-9={ ! tablesWithBiggerRow.includes(tableName) } 
                                on:click|preventDefault={ () => turnActiveOrInactive(record.id)} aria-hidden="true"  >
                            <div class="divTD_INACTIVE_BUTTON" >&nbsp;</div>
                          </div>

                      <!-- registro esta  ativo, ativa todas as opcoes de acao  -->
                      {:else}

                          <div class="divTD_ACTION" class:h-16={tablesWithBiggerRow.includes(tableName)}  class:h-9={! tablesWithBiggerRow.includes(tableName)} 
                                    on:click|preventDefault={ () => editRecord(event, record.id)} aria-hidden="true"  >  
                            <div class="divTD_EDIT_BUTTON " >&nbsp;</div>
                          </div>
                          <div class="divTD_ACTION" class:h-16={ tablesWithBiggerRow.includes(tableName) }  class:h-9={ ! tablesWithBiggerRow.includes(tableName) }
                                    on:click|preventDefault={ () => deleteRecord(event, record.id)} aria-hidden="true"  > 
                            <div class="divTD_DELETE_BUTTON" >&nbsp;</div>
                          </div>
                          <div class="divTD_ACTION" class:h-16={ tablesWithBiggerRow.includes(tableName) }  class:h-9={ ! tablesWithBiggerRow.includes(tableName) }
                                    on:click|preventDefault={ () => turnActiveOrInactive(record.id)} aria-hidden="true"  >  
                            <div class="divTD_ACTIVE_BUTTON" >&nbsp;</div>
                          </div>

                      {/if}
                  </div>
              {/if}

            </div>

          <!-- lista registros vazia -->
          {:else}
            <div class="divTABLE_NO_RECORDS">{$Terms.no_records}</div>
          {/each}

      {/if}

    </div> 

    <div class="divTBASEBOARD"  > 
      {@html baseboardMessage}
    </div>

</div>


