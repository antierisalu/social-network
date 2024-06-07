<script>
  import { fade, slide } from "svelte/transition";
  import Register from "./register.svelte";
  import {
    loggedIn,
    authError,
    displayUserAuthError,
    userInfo,
  } from "../stores";
  import Button from "../shared/button.svelte";
  import { updateSessionToken } from "../utils";
  let errorString = "";
  $: errorString = $authError;

  export function setLoggedIn() {
    loggedIn.set(true);
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
        displayUserAuthError("Invalid Credentials");
        throw new Error("Network response was not ok");
      }

      const data = await response.json();
      console.log("NOH", data);
      userInfo.set(data);
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

<div class="login" in:fade>
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
          if (user && password) {
            fetchData();
          } else {
            displayUserAuthError("Please fill in all the fields");
          }
        }}>Login</Button
      >
    </form>
    {#if errorString != ""}
      <div class="error" transition:slide>
        <Button
          type="primary"
          customStyle="width:300px; min-height: 35px; cursor: default; pointer-events: none;"
          >{errorString}</Button
        >
      </div>
    {/if}
    <div class="regBtn">
      <Button
        type="secondary"
        inverse={true}
        customStyle="width:200px"
        on:click={switchToRegister}>Register Instead</Button
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
