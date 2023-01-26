<script lang="ts">
  import { onMount } from "svelte";

  let backdrop: HTMLElement | null;
  let modal: HTMLElement | null;

  function handleClick() {
    modal?.classList.add("hidden");
    backdrop?.classList.add("hidden");
  }

  onMount(() => {
    backdrop?.addEventListener("click", handleClick);

    return () => {
      backdrop?.removeEventListener("click", handleClick);
    };
  });
</script>

<div
  bind:this={backdrop}
  class="fixed top-0 left-0 w-full h-full bg-[#00000096]"
/>

<div
  class="modal fixed top-0 left-0 fade outline-none"
  bind:this={modal}
  aria-modal="true"
  role="dialog"
>
  <div
    class="modal-dialog modal-dialog-scrollable pointer-events-none fixed flex h-full"
  >
    <div
      class="modal-content border-none shadow-lg flex flex-col w-3/5 h-5/6 m-auto pointer-events-auto text-white bg-gray-700 rounded-md outline-none text-current"
    >
      <div
        class="modal-header flex flex-shrink-0 items-center justify-between p-4 border-b border-gray-500 rounded-t-md"
      >
        <h5 class="text-xl font-medium leading-normal text-white">
          Modal title
        </h5>
        <button
          type="button"
          class="btn-close box-content w-4 h-4 p-1 border-none rounded-none opacity-75 focus:shadow-none focus:outline-none focus:opacity-100 hover:opacity-75 hover:no-underline"
          data-bs-dismiss="modal"
          aria-label="Close">x</button
        >
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
          class="inline-block px-6 py-2.5 bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-purple-700 hover:shadow-lg focus:bg-purple-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 ease-in-out"
          data-bs-dismiss="modal"
        >
          Close
        </button>
        <button
          type="button"
          class="inline-block px-6 py-2.5 bg-blue-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-blue-700 hover:shadow-lg focus:bg-blue-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-blue-800 active:shadow-lg transition duration-150 ease-in-out ml-1"
        >
          Save changes
        </button>
      </div>
    </div>
  </div>
</div>
