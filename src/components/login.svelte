<script>
  import Register from "./register.svelte";
  import Button from "../shared/button.svelte";
  import { createEventDispatcher } from "svelte";
  const dispatch = createEventDispatcher();

  export function setLoggedIn() {
    dispatch("login", {
      loginStatus: true,
    });
  }

  $: login = true;
  let user;
  let password;

  async function fetchData() {
    try {
      const response = await fetch("/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username: user, password: password }),
      });
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      setLoggedIn();
    } catch (error) {
      console.error("Error fetching data:", error.message);
    }
  }

  function switchToRegister() {
    login = !login; // This will show the Register component
  }
</script>

<div class="login">
  {#if login}
    <form on:submit|preventDefault>
      <input type="text" placeholder="E-mail" bind:value={user} required />
      <input
        type="password"
        placeholder="Password"
        required
        bind:value={password}
      />
      <Button type="secondary" on:click={fetchData}>Login</Button>
    </form>
    <div class="regBtn">
      <Button type="secondary" inverse={true} btn200px={true} on:click={switchToRegister}
        >Register Instead</Button
      >
    </div>
  {:else}
    <Register on:click={switchToRegister} />
  {/if}
</div>

<style>
  .login {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  form {
    display: flex;
    flex-direction: column;
  }

  input {
    width: 300px;
    border-radius: 6px;
    padding: 8px 12px;
  }
</style>
