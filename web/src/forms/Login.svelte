<script lang="ts">
  let formTouched = false;
  let usernameShowValidationError = false;
  let usernameValidationError: string | null = "Required";
  let passwordShowValidationError = false;
  let passwordValidationError: string | null = "Required";

  $: usernameIsValid = usernameValidationError === null;
  $: passwordIsValid = passwordValidationError === null;
  $: submitEnabled = formTouched && usernameIsValid && passwordIsValid;

  function handleUsernameInput(v: string) {
    if (v.trim().length === 0) {
      usernameValidationError = "Required";
      return;
    }

    if (v.trim().length !== v.length) {
      usernameValidationError = "Invalid value";
      return;
    }

    if (v.length > 64) {
      usernameValidationError = "Cannot be longer than 64 characters";
      return;
    }

    usernameValidationError = null;
    return;
  }

  function handlePasswordInput(v: string) {
    if (v.trim().length === 0) {
      passwordValidationError = "Required";
      return;
    }

    if (v.length > 128) {
      passwordValidationError = "Cannot be longer than 128 characters";
      return;
    }

    passwordValidationError = null;
    return;
  }
</script>

<form action="/auth/login" method="post">
  <div class="mb-6 last:mb-0">
    <label class="block font-bold mb-2" for="username">Username</label>
    <input
      id="username"
      name="username"
      autocomplete="email"
      class="px-3 py-2 max-w-full border-gray-700 rounded w-full dark:placeholder-gray-400 focus:ring focus:ring-blue-600 focus:border-blue-600 focus:outline-none h-12 border bg-gray-800"
      placeholder="john.doe"
      on:focusin|once|self|trusted|capture={() => {
        formTouched = true;
      }}
      on:focusout|once|self|trusted|capture={() => {
        usernameShowValidationError = true;
      }}
      on:input|self|trusted|capture={(e) => {
        handleUsernameInput(e.currentTarget.value);
      }}
    />
    {#if usernameShowValidationError && usernameValidationError !== null}
      <p class="mt-2 text-pink-600">{usernameValidationError}</p>
    {/if}
  </div>
  <div class="mb-6 last:mb-0">
    <label class="block font-bold mb-2" for="password">Password</label>
    <input
      id="password"
      name="password"
      type="password"
      autocomplete="current-password"
      class="px-3 py-2 max-w-full border-gray-700 rounded w-full dark:placeholder-gray-400 focus:ring focus:ring-blue-600 focus:border-blue-600 focus:outline-none h-12 border bg-gray-800"
      placeholder="**********"
      on:focusin|once|self|trusted|capture={() => {
        formTouched = true;
      }}
      on:focusout|once|self|trusted|capture={() => {
        passwordShowValidationError = true;
      }}
      on:input|self|trusted|capture={(e) => {
        handlePasswordInput(e.currentTarget.value);
      }}
    />
    {#if passwordShowValidationError && passwordValidationError !== null}
      <p class="mt-2 text-pink-600">{passwordValidationError}</p>
    {/if}
  </div>
  <hr class="my-6 -mx-6 dark:border-gray-800 border-t" />
  <div class="-mb-3 text-center">
    <button
      class="w-full focus:outline-none uppercase transition-colors focus:ring duration-300 border cursor-pointer rounded-lg border-violet-600 ring-violet-700 bg-violet-500 text-white hover:bg-violet-600 mb-3 py-2 px-3"
      type="submit"
      class:disabled={!submitEnabled}
      disabled={!submitEnabled}
    >
      <span class="px-2">Sign in</span>
    </button>
  </div>
</form>

<style>
  button.disabled {
    cursor: not-allowed;
    background-color: dimgray;
    border-color: rgb(143, 143, 143);
  }
</style>
