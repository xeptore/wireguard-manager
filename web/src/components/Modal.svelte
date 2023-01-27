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
    in:fly={{ y: 200, duration: 300 }}
    out:fly={{ y: 200, duration: 300 }}
    class="modal-dialog modal-dialog-scrollable pointer-events-none fixed flex h-full"
  >
    <div
      class="modal-content border-none shadow-lg flex flex-col w-2/3 h-5/6 m-auto pointer-events-auto text-white bg-gray-800 rounded-md outline-none text-current"
    >
      <div
        class="modal-header flex items-center justify-between p-4 border-b border-gray-500 rounded-t-md"
      >
        <h5 class="text-xl font-medium leading-normal text-white">
          Modal title
        </h5>
        <button
          type="button"
          class="btn-close box-content w-4 h-4 p-1 border-none rounded-none opacity-75 focus:shadow-none focus:outline-none hover:opacity-100 hover:no-underline"
          on:click={closeModal}
          aria-label="Close"
        >
          <i class="fa-solid fa-xmark align-text-top" />
        </button>
      </div>
      <div
        class="modal-body dark-scrollbars relative p-4 max-h-[76vh] overflow-hidden overflow-y-scroll"
      >
        <p>
          This is some placeholder content to show the scrolling behavior for
          modals. We use repeated line breaks to demonstrate how content can
          exceed minimum inner height, thereby showing inner scrolling. When
          content becomes longer than the predefined max-height of modal,
          content will be cropped and scrollable within the modal.
        </p>
        <br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br
        /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br
        /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br
        /><br /><br /><br /><br />
        <p>This content should appear at the bottom after you scroll.</p>
      </div>
      <div
        class="modal-footer flex flex-shrink-0 flex-wrap items-center justify-end p-4 border-t border-gray-500 rounded-b-md"
      >
        <button
          type="button"
          class="inline-block px-6 py-2.5 bg-transparent ring-2 ring-blue-600 text-blue-500 hover:bg-blue-600 hover:text-white ring-inset text-sm rounded-lg shadow-md hover:shadow-lg focus:shadow-lg focus:outline-none focus:ring-0 active:shadow-lg transition duration-150 ease-in-out"
          on:click={closeModal}
        >
          Cancel
        </button>
        <button
          type="button"
          class="inline-block px-6 py-2.5 text-white text-sm rounded-lg shadow-md hover:shadow-lg ring-blue-500 bg-blue-600 hover:bg-blue-700 focus:shadow-lg focus:outline-none focus:ring-0 active:shadow-lg transition duration-150 ease-in-out ml-1"
        >
          Create
        </button>
      </div>
    </div>
  </div>
</div>
