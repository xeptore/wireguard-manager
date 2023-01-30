<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { fly } from "svelte/transition";

  const dispatch = createEventDispatcher();

  function closeModal() {
    dispatch("close");
  }

  function handleKeyPress(e: KeyboardEvent) {
    if (e.key === "Escape") {
      dispatch("close");
    }
  }
</script>

<svelte:window on:keydown={handleKeyPress} />
<div
  on:click|once={closeModal}
  aria-hidden="true"
  class="fixed top-0 left-0 w-full h-full bg-[#00000096]"
/>

<div
  class="modal fixed top-0 left-0 fade outline-none"
  aria-modal="true"
  role="dialog"
>
  <div
    in:fly={{ y: 200, duration: 300, delay: 1_000 }}
    out:fly={{ y: 200, duration: 300 }}
    class="modal-dialog modal-dialog-scrollable pointer-events-none fixed flex h-full w-full"
  >
    <div
      class="modal-content border-none shadow-lg flex flex-col w-2/3 h-5/6 m-auto pointer-events-auto text-white bg-gray-800 rounded-md outline-none text-current"
    >
      <div
        class="modal-header flex items-center justify-between p-4 border-b border-gray-500 rounded-t-md"
      >
        <h5 class="text-xl font-medium leading-normal text-white">
          New Client QR Code
        </h5>
        <button
          type="button"
          class="btn-close box-content w-4 h-4 p-1 border-none rounded-none opacity-75 focus:shadow-none focus:outline-none hover:opacity-100 hover:no-underline"
          on:click={closeModal}
          aria-label="Close"
        >
          <i class="fa-solid fa-xmark" />
        </button>
      </div>
      <div
        class="modal-body dark-scrollbars relative p-4 max-h-[76vh] overflow-x-hidden overflow-y-scroll flex h-full"
      >
        <slot />
      </div>
    </div>
  </div>
</div>
