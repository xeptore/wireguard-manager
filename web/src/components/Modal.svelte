<script lang="ts">
  import pWaterfall from "p-waterfall";
  import { createEventDispatcher } from "svelte";
  import { fly } from "svelte/transition";

  let formTouched = false;
  let nameShowValidationError = false;
  let nameValidationError: string | null = "Required";
  let descriptionShowValidationError = false;
  let descriptionValidationError: string | null = null;

  $: nameIsValid = nameValidationError === null;
  $: descriptionIsValid = descriptionValidationError === null;
  $: submitEnabled = formTouched && nameIsValid && descriptionIsValid;

  function handleNameInput(v: string) {
    if (v.trim().length === 0) {
      nameValidationError = "Required";
      return;
    }

    if (v.trim().length !== v.length) {
      nameValidationError = "Invalid value";
      return;
    }

    if (v.length > 256) {
      nameValidationError = "Cannot be longer than 256 characters";
      return;
    }

    nameValidationError = null;
    return;
  }

  function handleDescriptionInput(v: string) {
    if (v.trim().length !== v.length) {
      nameValidationError = "Invalid value";
      return;
    }

    if (v.length > 10_000) {
      descriptionValidationError = "Cannot be longer than 10,000 characters";
      return;
    }

    descriptionValidationError = null;
    return;
  }

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
      class="modal-content border-none shadow-lg flex flex-col w-2/3 lg:w-1/2 xl:w-1/2 2xl:w-3/7 m-auto pointer-events-auto text-white bg-gray-800 rounded-md outline-none text-current"
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
              on:focusin|once|self|trusted|capture={() => {
                formTouched = true;
              }}
              on:focusout|once|self|trusted|capture={() => {
                nameShowValidationError = true;
              }}
              on:input|self|trusted|capture={(e) => {
                handleNameInput(e.currentTarget.value);
              }}
            />
            {#if nameShowValidationError && nameValidationError !== null}
              <p class="mt-2 text-pink-600">{nameValidationError}</p>
            {/if}
          </div>
          <div class="mb-6 last:mb-0">
            <label class="block font-bold mb-2" for="description"
              >Description</label
            >
            <input
              id="description"
              name="description"
              type="description"
              autocomplete="off"
              class="px-3 py-2 max-w-full border-gray-700 rounded w-full dark:placeholder-gray-400 focus:ring focus:ring-blue-600 focus:border-blue-600 focus:outline-none h-12 border bg-gray-800"
              placeholder="Here goes some description..."
              on:focusin|once|self|trusted|capture={() => {
                formTouched = true;
              }}
              on:focusout|once|self|trusted|capture={() => {
                descriptionShowValidationError = true;
              }}
              on:input|self|trusted|capture={(e) => {
                handleDescriptionInput(e.currentTarget.value);
              }}
            />
            {#if descriptionShowValidationError && descriptionValidationError !== null}
              <p class="mt-2 text-pink-600">{descriptionValidationError}</p>
            {/if}
          </div>
          <hr class="my-6 -mx-6 dark:border-gray-800 border-t" />
          <div class="-mb-3 text-center">
            <button
              class="w-full focus:outline-none uppercase transition-colors focus:ring duration-150 border cursor-pointer rounded-lg border-violet-600 ring-violet-700 bg-violet-500 text-white hover:bg-violet-600 mb-3 py-2 px-3"
              type="submit"
              class:disabled={!submitEnabled}
              disabled={!submitEnabled}
            >
              <span class="px-2">Create</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</div>

<style>
  button.disabled {
    cursor: not-allowed;
    background-color: dimgray;
    border-color: rgb(143, 143, 143);
  }
</style>
