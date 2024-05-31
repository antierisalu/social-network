<script>
  import Register from "./register.svelte";
  import Button from "../shared/button.svelte";
  import { updateSessionToken } from "../utils";
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
        body: JSON.stringify({ email: user, password: password }),
      });
      if (!response.ok) {
        console.log(response);
        throw new Error("Network response was not ok");
      }

      const data = await response.json();
      updateSessionToken(data.token, data.expires);
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
      <input type="email" placeholder="E-mail" bind:value={user} required />
      <input
        type="password"
        placeholder="Password"
        required
        bind:value={password}
      />
      <Button
        type="secondary"
        on:click={() => {
          if (user && password) fetchData();
        }}>Login</Button
      >
    </form>
    <div class="regBtn">
      <Button type="" btn200px={true} on:click={switchToRegister}
        >Register Instead</Button
      >
    </div>
  {:else}
    <Register on:click={switchToRegister} on:login={setLoggedIn} />
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
