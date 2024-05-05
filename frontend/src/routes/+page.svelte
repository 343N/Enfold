<script>
  import {Heading, A} from "flowbite-svelte"
  import { GetDirectory } from "$lib/wailsjs/go/main/App";
	import { Delay } from "$lib/delay";
  let drives = [{Name:"Grabbing drive information..."}];
  (async ()=>{
    await Delay(3000)
    // drives = await GetDirectory("/", false)
    console.log(drives)
  })()
</script>

<Heading class="text-gray-100 font-light">Enfold</Heading>
<hr class="w-100 h-1 mx-auto my-1 bg-gray-100 border-0 rounded md:my-4 dark:bg-gray-700">
{#each drives as d}
  <A><Heading tag="h5" class="text-gray-100 font-light">{d.Name}</Heading></A>
  <br>
{/each}

<style>
</style>
