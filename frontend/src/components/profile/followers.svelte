<script>
  import { createEventDispatcher } from "svelte";
  import { userProfileData, API_URL, IMAGE_URL } from "../../stores";
  import { scale, fade, fly } from "svelte/transition";
  import Button from "../../shared/button.svelte";

  export let x;
  export let y;

  async function selectUser(userID) {
    const response = await fetch(`${API_URL}/user?id=${userID}`,{
      credentials: "include",
    });
    if (response.ok) {
      const selectedUser = await response.json();
      console.log("OUUU", selectedUser);
      userProfileData.set(selectedUser);
    } else {
      console.error("Error fetching users:", response.status);
    }
  }
  const dispatch = createEventDispatcher();
  export let followers;
  function closeOverlay() {
    dispatch("close");
  }
  console.log(x, y);
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div out:fade={{ duration: 150 }} class="overlay" on:click={closeOverlay}></div>
<div in:fly={{ x: x, y: y }} out:fade={{ duration: 150 }} class="modal">
  <div
    in:scale={{ start: 0.1, duration: 700 }}
    out:fade={{ duration: 150 }}
    class="modal todal"
  >
    <div class="modal-content">
      {#each followers as follower}
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div
          class="singleUser"
          on:click={() => {
            selectUser(follower.ID);
            closeOverlay();
          }}
        >
          <!-- svelte-ignore a11y-missing-attribute -->
          <img src={IMAGE_URL}{follower.Avatar} />
          {follower.FirstName}
          {follower.LastName}
        </div>
      {/each}
    </div>
    <div class="closeButton">
      <Button on:click={closeOverlay}>Close Overlay</Button>
    </div>
  </div>
</div>

<style>
  img {
    max-height: 20px;
  }

  .todal {
    padding: 20px;
  }

  .overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 0;
  }

  .singleUser {
    border: solid 1px #333;
    border-radius: 6px;
    cursor: pointer;
    padding: 8px;
    background: #011;
    margin: 2px;
  }

  .modal {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background-color: #011;
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
    z-index: 2;
    width: 80%;
    max-width: 500px;
    border-radius: 8px;
  }

  .modal-content {
    position: relative;
  }

  .closeButton {
    margin-top: 20px;
  }
</style>
