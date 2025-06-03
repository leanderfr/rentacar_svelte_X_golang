

<script>

import { reportedErrorMessage, imagesStillLoading, Terms } from '$lib/stores/stores.js'
import { onMount } from 'svelte';  

// desativa qq animacao gif de 'carregando...'
onMount( () => {
  $imagesStillLoading=false  
  hideLoadingGif() 
} )

</script>


<div  class="flex flex-col w-[80%] bg-white relative rounded-lg h-[50%]"  >

  <div class="flex flex-row w-full h-full"  >

    <div class="items-center bg-center bg-[length:80%_80%] justify-center w-[40%] bg-no-repeat bg-icon-reported-error"  > 
    </div>

    <div class="w-[60%] text-red-800 text-[20px] flex flex-col  justify-center gap-y-11" > 
        <div>
          {$Terms.there_was_js_error_line1}
        </div>
        <div>
          {$Terms.there_was_js_error_line2}
        </div>
        <div>
          {$Terms.there_was_js_error_line3}
        </div>

        <div class="w-full flex flex-row" > 
          {$Terms.there_was_js_error_line4}=&nbsp;&nbsp;<div class="text-red-600">{$reportedErrorMessage}</div>
        </div>
    </div>

  </div>

</div> 


