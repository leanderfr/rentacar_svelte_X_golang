

<script>

// variaveis genericas, writeables
import { Terms, backendUrlGolang  } from '$lib/stores/stores.js'
import { onMount, createEventDispatcher } from 'svelte'

import { slidingMessage, logAPI } from '$js/utils.js'

const dispatch = createEventDispatcher()

/********************************************************************************************************************************************************
 se usuario clicar no [ X ] do formulario, despacha comando para fechar janela
********************************************************************************************************************************************************/
const btnCloseFormClicked = () => {
    dispatch('dispatchCloseWorkgroupLogin')
}

/********************************************************************************************************************************************************
 usuario informou grupo existente, avisa para fechar tudo e buscar dados do grupo 
********************************************************************************************************************************************************/
const accessAnotherGroupData = () => {
    dispatch('dispatchAcessAnotherGroupData')
}



/********************************************************************************************************************************************************
********************************************************************************************************************************************************/
onMount(() => {

  setTimeout(() => {
    jq('#txtWorkgroup').focus()  
  }, 300);

})


/********************************************************************************************************************************************************
verifica se o grupo digitado existe e esta em uso
********************************************************************************************************************************************************/
async function performLogingWorkgroup() {

jq('#unavailableGroup').html('')

let wrk = jq.trim( jq('#txtWorkgroup').val() )

if (wrk.length < 4)  {
    slidingMessage('slidingWindowMessage', $Terms.errormsg_type_workgroup, 2000)    
    return
}


// necessario abrir evento assincrono para exibir div ajax loading, caso contrario navegador nao atualiza a tela
setTimeout(() => {
  showLoadingGif();   // mostra animacao 'processando..' mesmo sem tela ter sido renderizada
}, 1);

// manipulacao de grupos e notificacoes é feita somente pelo backend golang
let _route_ = `${$backendUrlGolang}/workgroup/change/${wrk}` 
logAPI('GET', _route_)

await fetch(_route_, { method: 'GET' })

.then((response) => {
  setTimeout( () => {hideLoadingGif();}, 1)

  if (!response.ok) {
    throw new Error(`Access workgroup Err fatal= ${response.status}`);
  }
  return response.text();
})

.then((infoWorkgroup) => { 

  // grupo digitado nao existe
  if (infoWorkgroup.indexOf('none')!=-1)  {
    jq('#unavailableGroup').html( $Terms.workgroup_not_available )

    setTimeout( () => {
      jq('#workgroupLogin').effect( "shake", {distance:20} );
    }, 100)
    
    return
  }

  let _infoWorkgroup = infoWorkgroup.split('|')             // infoWorkgroup = __success__|workgroup name|workgroup id
  let _currentWorkgroupName = _infoWorkgroup[1]
  let _currentWorkgroupId = _infoWorkgroup[2]

  // grava workgroup acessado como local storage
  localStorage.setItem("rentacar_workgroup_name", _currentWorkgroupName);
  
  jq('#btnWorkgroupMenu').html ( _currentWorkgroupName )   // botao 'grupo' ,  Home.svelte

  setTimeout(() => {
    accessAnotherGroupData()
  }, 100);
})

}

/****************************************************************************************************
 se usuario passar mouse sobre div login, coloca foto no campo 'grupo'
****************************************************************************************************/
const focusWorkgroupLogin = () =>  {
  if ( ! jq('#txtWorkgroup').is(":focus") )  {
    jq('#txtWorkgroup').focus();
    jq('#txtWorkgroup').select();
  }
}

</script>

<div  class="font-Roboto flex flex-col w-auto  overflow-hidden bg-white p-[4px] rounded-xl "  id='workgroupLogin' on:mouseenter={ focusWorkgroupLogin } aria-hidden="true" >

    <div  class="flex flex-col w-[900px]  border-2  relative rounded-lg " style='background-color: #edf8f4;border-color:green' >

        <div class="flex flex-col w-full h-auto pt-6 pb-10 pr-5 pl-3 " >

          <div class="text-xl pb-8" >
            {$Terms.to_access_another_workgroup_data}
          </div>

          <div class="flex flex-row justify-between pb-14 gap-4" >
            <div style='text-align: justify;text-justify: inter-word'>{$Terms.note_about_group_login}</div>
            <div><img style='min-width:320px;min-height:70px;' src='notifications_.png' alt=''> </div>
          </div>


          <div class="flex flex-col w-auto ">

            <div class="flex flex-row w-full ">
              <div class="flex py-[6px] pr-5">{$Terms.current_workgroup}:</div>

              <div class="w-[60%] pl-3">
                <div class="flex flex-col">
                  <div class="w-full" on:mouseenter={ focusWorkgroupLogin } aria-hidden="true" >
                    <input type="text" autocomplete="off" sequence="1" maxlength='50' minlength='2'  id="txtWorkgroup" class='text_formFieldValue' on:focus={(e)=>{e.target.select()}} >
                  </div>
                  <div class="w-full text-red-600 h-10 py-3" id='unavailableGroup'>
                    &nbsp;
                  </div>
                </div>
              </div>

            </div>

          </div>

        </div>

        <!-- botoes confirmar ou cancelar  -->
        <div class="flex flex-row w-auto justify-between px-6 border-t-[1px] border-t-gray-300 py-3 gap-20">
          <button  id="btnCLOSE_LOGIN_WINDOW" class="btnCANCEL" on:click={btnCloseFormClicked} >{$Terms.button_cancel}</button>

          <button  id="btnLOGIN" class="btnLOGIN" aria-hidden="true" on:click={()=>{performLogingWorkgroup()}} >Enter= {$Terms.button_login_workgroup}</button>
        </div>


    </div> 

</div>

