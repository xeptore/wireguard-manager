<script lang="ts">
  import Modal from "../Modal.svelte";
  import PeerQrCodeModal from "../PeerQRCodeModal.svelte";

  let showCreateClientModal = false;
  let showQRCodeModal = false;
  let qrCodeSvg: string | null = null;

  function openCreateClientModal() {
    showCreateClientModal = true;
  }

  async function handleSuccess(e: CustomEvent<string>) {
    const response = await fetch(`/peers/${e.detail}`);
    const content = await response.text();
    qrCodeSvg = content;
    showQRCodeModal = true;
  }
</script>

{#if showCreateClientModal}
  <Modal
    on:close|once={() => (showCreateClientModal = false)}
    on:success={handleSuccess}
  />
{/if}
{#if showQRCodeModal}
  <PeerQrCodeModal on:close|once={() => (showQRCodeModal = false)}>
    <img
      class="w-full h-full"
      src={`data:image/svg+xml;base64,${qrCodeSvg}`}
      alt=""
    />
  </PeerQrCodeModal>
{/if}
<section class="pt-6 mb-6 flex items-center justify-between">
  <div class="flex items-center justify-start">
    <span class="inline-flex justify-center items-center w-6 h-6 mr-2">
      <svg viewBox="0 0 24 24" width="20" height="20" class="inline-block">
        <path
          fill="currentColor"
          d="M16 17V19H2V17S2 13 9 13 16 17 16 17M12.5 7.5A3.5 3.5 0 1 0 9 11A3.5 3.5 0 0 0 12.5 7.5M15.94 13A5.32 5.32 0 0 1 18 17V19H22V17S22 13.37 15.94 13M15 4A3.39 3.39 0 0 0 13.07 4.59A5 5 0 0 1 13.07 10.41A3.39 3.39 0 0 0 15 11A3.5 3.5 0 0 0 15 4Z"
        />
      </svg>
    </span>
    <h1 class="leading-tight text-2xl">Clients</h1>
  </div>
  <button
    type="button"
    on:click={openCreateClientModal}
    class="ring-blue-500 bg-blue-500 text-white hover:bg-blue-600 p-2 px-3 focus:outline-none transition-all duration-200 focus:ring-1 hover:ring-1 cursor-pointer rounded-lg"
  >
    <i class="fa-solid fa-plus" />
    <span class="ml-1">Create New</span>
  </button>
</section>
