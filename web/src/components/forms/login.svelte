<script lang="ts">
  let error: string | null = null;
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
    />
    {error !== null ? error : ""}
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
  <div class="-mb-3 text-center">
    <button
      class="w-full focus:outline-none uppercase transition-colors focus:ring duration-150 border cursor-pointer rounded-lg border-violet-600 ring-violet-700 bg-violet-500 text-white hover:bg-violet-600 mb-3 py-2 px-3"
      type="submit"
    >
      <span class="px-2">Sign in</span>
    </button>
  </div>
</form>
