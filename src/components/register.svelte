<script>
  import Button from "../shared/button.svelte";
  import { updateSessionToken } from "../utils";
  import { createEventDispatcher } from "svelte";
  const dispatch = createEventDispatcher();

  let userData = {
    email: "",
    firstName: "",
    lastName: "",
    dateOfBirth: "",
    avatar: "",
    nickName: "",
    aboutMe: "",
    password: "",
    passwordConfirm: "",
  };

  //send register info to backend
  async function registerUser(registerInfo) {
    try {
      const response = await fetch("/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(registerInfo),
      });

      if (!response.ok) {
        const errorMessage = await response.text();
        throw new Error(errorMessage);
      }

      const data = await response.json();
      console.log("REGISTER:", data);
      updateSessionToken(data.token, data.expires);
      dispatch("login", {
        loginStatus: true,
      });
    } catch (error) {
      console.error("Error registering user:", error.message);
    }
  }
</script>

<div class="register">
  <form on:submit|preventDefault>
    <input
      type="email"
      placeholder="E-mail *"
      bind:value={userData.email}
      required
    />
    <input
      type="text"
      placeholder="First Name *"
      bind:value={userData.firstName}
      required
    />
    <input
      type="text"
      placeholder="Last Name *"
      bind:value={userData.lastName}
      required
    />
    <input
      type="date"
      placeholder="Date of Birth *"
      bind:value={userData.dateOfBirth}
      required
    />
    <div class="fakeInput">
      <label for="avatar">Avatar less than 100kb</label>
      <input
        type="file"
        id="avatar"
        class="hidden"
        accept="image/png, image/jpeg"
        bind:value={userData.avatar}
      />
    </div>
    <input type="text" placeholder="Nickname" bind:value={userData.nickName} />
    <input type="text" placeholder="About Me" bind:value={userData.aboutMe} />
    <input
      type="password"
      placeholder="Password"
      required
      bind:value={userData.password}
    />
    <input
      type="password"
      placeholder="Confirm Password"
      required
      bind:value={userData.passwordConfirm}
    />
    <Button
      type="secondary"
      on:click={() => {
        if (
          userData.email &&
          userData.firstName &&
          userData.lastName &&
          userData.dateOfBirth &&
          userData.password &&
          userData.passwordConfirm
        ) {
          registerUser(userData);
        }
      }}
    >
      Register
    </Button>
  </form>
  <Button type="" btn200px={true} on:click>Login Instead</Button>
</div>

<style>
  .fakeInput {
    display: inline-block;
    padding: 8px 10px;
    border: 1px solid #ccc;
    border-radius: 6px;
    background-color: #f5f5f5;
    margin-bottom: 8px;
    cursor: pointer;
  }

  .fakeInput label {
    cursor: pointer;
    text-align: left;
    margin: 1;
    opacity: 0.7;
  }

  .hidden {
    display: none;
  }

  .register {
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
