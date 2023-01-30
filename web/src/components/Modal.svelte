<script lang="ts">
  import pWaterfall from "p-waterfall";
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

  async function handleSubmit(formElement: HTMLFormElement) {
    await pWaterfall(
      [
        async (formData) => {
          const values = Object.fromEntries(formData.entries());
          try {
            return await fetch("http://localhost:8080/peers", {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              redirect: "follow",
              referrerPolicy: "no-referrer",
              credentials: "include",
              body: JSON.stringify(values),
            });
          } catch (error) {
            console.error(error);
            throw error;
          }
        },
        async (apiRes) => {
          const x = await apiRes.text();
          dispatch("success", x);
          closeModal();
        },
      ],
      new FormData(formElement)
    );
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
    class="modal-dialog modal-dialog-scrollable pointer-events-none fixed flex h-full w-full"
  >
    <div
      class="modal-content border-none shadow-lg flex flex-col w-2/3 m-auto pointer-events-auto text-white bg-gray-800 rounded-md outline-none text-current"
    >
      <div
        class="modal-header flex items-center justify-between p-4 border-b border-gray-500 rounded-t-md"
      >
        <h5 class="text-xl font-medium leading-normal text-white">
          New Client
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
        class="modal-body dark-scrollbars relative p-4 max-h-[76vh] overflow-x-hidden overflow-y-scroll"
      >
        <form on:submit|preventDefault={(e) => handleSubmit(e.currentTarget)}>
          <div class="mb-6 last:mb-0">
            <label class="block font-bold mb-2" for="name">Name</label>
            <input
              id="name"
              name="name"
              autocomplete="name"
              class="px-3 py-2 max-w-full border-gray-700 rounded w-full dark:placeholder-gray-400 focus:ring focus:ring-blue-600 focus:border-blue-600 focus:outline-none h-12 border bg-gray-800"
              placeholder="Ali Agha"
            />
          </div>
          <div class="mb-6 last:mb-0">
            <label class="block font-bold mb-2" for="description"
              >Description</label
            >
            <div>
              <div class="relative">
                <input
                  id="description"
                  name="description"
                  type="description"
                  autocomplete="off"
                  class="px-3 py-2 max-w-full border-gray-700 rounded w-full dark:placeholder-gray-400 focus:ring focus:ring-blue-600 focus:border-blue-600 focus:outline-none h-12 border bg-gray-800"
                  placeholder="Here goes some description..."
                />
              </div>
            </div>
          </div>
          <hr class="my-6 -mx-6 dark:border-gray-800 border-t" />
          <div class="-mb-3 text-center">
            <button
              class="w-full focus:outline-none uppercase transition-colors focus:ring duration-150 border cursor-pointer rounded-lg border-violet-600 ring-violet-700 bg-violet-500 text-white hover:bg-violet-600 mb-3 py-2 px-3"
              type="submit"
            >
              <span class="px-2">Create</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</div>
