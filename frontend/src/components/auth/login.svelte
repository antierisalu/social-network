<script>
  import { fade, slide } from "svelte/transition";
  import Register from "./register.svelte";
  import {
    loggedIn,
    authError,
    displayUserAuthError,
    userInfo,
    userProfileData,
    API_URL
  } from "../../stores";
  import Button from "../../shared/button.svelte";
  import { updateSessionToken, fetchUsers } from "../../utils";
  let errorString = "";
  $: errorString = $authError;

  export function setLoggedIn() {
    loggedIn.set(true);
    fetchUsers();
    fetchNotifications();
  }

  $: login = true;
  let user;
  let password;

  async function fetchData() {
    try {
      const response = await fetch(`${API_URL}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: user, password: password }),
      });
      if (!response.ok) {
        let text = await response.text();
        console.log(response, text);

        displayUserAuthError(text);
        throw new Error("Network response was not ok: " + text);
      }

      const data = await response.json();
      console.log("NOH", data);
      userInfo.set(data);
      userProfileData.set(data);
      updateSessionToken(data.session, 24);
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
        w120={false}
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
          w120={false}
          customStyle="width:300px; min-height: 35px; cursor: default; pointer-events: none;"
          >{errorString}</Button
        >
      </div>
    {/if}
    <div class="regBtn">
      <Button
        type="secondary"
        inverse
        w120={false}
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
