
<script>

// variaveis genericas, writeables
import { Terms  } from '$lib/stores/stores.js'
import { onMount } from 'svelte';  

export let userConfirmedWorkgroupDataReset 

</script>

<div class='w-full h-full  absolute flex  left-0 top-0 z-40 bg-[rgba(0,0,0,0.1)] cursor-pointer justify-center items-center'    >  

    <div class="flex flex-col w-[80%] bg-white relative rounded-lg  h-48 border-blue-600 border-[2px] "  > 
        <div class="flex flex-col h-36 justify-center items-center gap-4"  > 

          <div style="font-size:24px;color:blue"   id='workgroupCloningStatus' >&nbsp;</div>

          <div class="divLOADING_BAR_OUTER " >
              <div class="divLOADING_BAR_INNER justify-center " id='workgroupCloningPercent' >
              </div>
          </div>
        </div>


        <div class="flex flex-row justify-between "    >
          <div class="flex flex-row text-2xl pl-3"    >
            {#if (userConfirmedWorkgroupDataReset) }
                <div class="flex justify-start h-4  text-left">&nbsp;{$Terms.workgroup_data_resetting}</div>
            {:else}
              <div class="flex justify-start h-4  text-left">&nbsp;{$Terms.chosen_workgroup}</div>
            {/if}
            <div class="flex justify-start h-4  text-left pl-4" style='color:blue' id='chosenWorkgroupName'>&nbsp;</div>
          </div>
          <div class="flex flex-row text-2xl pr-5"  style='color:blue' id='completedPercent'>&nbsp;

          </div>
        </div>

    </div>

</div>
