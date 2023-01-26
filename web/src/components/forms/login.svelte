<script lang="ts">
  let error: string | null = null;

  function handleSubmit(e: Event) {
    const formData = new FormData(e.target as HTMLFormElement);
    const values = Object.fromEntries(formData.entries());
    console.log(JSON.stringify(values, null, 2));

    error = 'invalid message'
    setTimeout(() => {
      location.assign('/');
    }, 3000);
  }
</script>

<form on:submit|preventDefault={handleSubmit}>
  <div class="mb-6 last:mb-0">
    <label class="block font-bold mb-2" for="username">Username</label>
    <input
      id="username"
      name="username"
      autocomplete="email"
      class="px-3 py-2 max-w-full border-gray-700 rounded w-full dark:placeholder-gray-400 focus:ring focus:ring-blue-600 focus:border-blue-600 focus:outline-none h-12 border bg-gray-800"
      placeholder="john.doe"
    />
    {error !== null ? error : ''}
  </div>
  <div class="mb-6 last:mb-0">
    <label class="block font-bold mb-2" for="password">Password</label>
    <div>
      <div class="relative">
        <input
          id="password"
          name="password"
          type="password"
          autocomplete="current-password"
          class="px-3 py-2 max-w-full border-gray-700 rounded w-full dark:placeholder-gray-400 focus:ring focus:ring-blue-600 focus:border-blue-600 focus:outline-none h-12 border bg-gray-800"
          placeholder="**********"
        />
      </div>
    </div>
  </div>
  <hr class="my-6 -mx-6 dark:border-gray-800 border-t" />
  <div class="flex items-center justify-center undefined -mb-3 flex-wrap">
    <button
      class="inline-flex justify-center items-center whitespace-nowrap focus:outline-none transition-colors focus:ring duration-150 border cursor-pointer rounded border-blue-600 dark:border-blue-500 ring-blue-300 dark:ring-blue-700 bg-blue-600 dark:bg-blue-500 text-white hover:bg-blue-700 hover:border-blue-700 hover:dark:bg-blue-600 hover:dark:border-blue-600 mr-3 last:mr-0 mb-3 undefined py-2 px-3"
      type="submit"
    >
      <span class="px-2">Login</span>
    </button>
  </div>
</form>
