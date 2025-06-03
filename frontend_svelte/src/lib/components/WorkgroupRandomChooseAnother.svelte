

<script>

// variaveis genericas, writeables
import { Terms   } from '$lib/stores/stores.js'
import { createEventDispatcher } from 'svelte'

const dispatch = createEventDispatcher()


/********************************************************************************************************************************************************
 se usuario clicar no [ X ] do formulario, despacha comando para fechar janela
********************************************************************************************************************************************************/
const btnCloseFormClicked = () => {
    dispatch('dispatchCloseWorkgroupAnother')
}

/********************************************************************************************************************************************************
 usuario confirmou que quer o sorteio de outro nome de grupo 
********************************************************************************************************************************************************/
const performRandomChooseAnother = () => {
    dispatch('dispatchPerformWorkgroupRandomChooseAnother')
}


</script>

<div  class="font-Rroboto flex flex-col w-auto  overflow-hidden bg-white p-[4px] rounded-xl "  id='workgroupRandomChooseAnother'  >

    <div  class="flex flex-col w-[900px]  border-2  relative rounded-lg " style='background-color: #e6e6ff;border-color:blue' >

        <div class="flex flex-col w-full h-auto pt-6 pb-10 pr-5 pl-3 " >

          <div class="text-xl pb-8" >
            {$Terms.to_choose_another_workgroup_randomly}
          </div>

          <div class="flex flex-row justify-between pb-14 gap-4" >
            <div style='text-align: justify;text-justify: inter-word;height:50px;' >
                {$Terms.note_about_random_choose_another_group}

            </div>
          </div>


        </div>

        <!-- botoes confirmar ou cancelar  -->
        <div class="flex flex-row w-auto justify-between px-6 border-t-[1px] border-t-gray-300 py-3 gap-20">
          <button  id="btnCLOSE_RANDOM_CHOOSE_WINDOW" class="btnCANCEL" on:click={btnCloseFormClicked} >{$Terms.button_cancel}</button>

          <button  id="btnCONTINUE" class="btnCONTINUE" aria-hidden="true" on:click={()=>{performRandomChooseAnother()}} >{$Terms.button_yes_continue}</button>
        </div>


    </div> 

</div>

