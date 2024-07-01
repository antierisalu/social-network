<script>
  import { allUsers, userProfileData } from "../../stores";

  $: users = $allUsers;

  var searchQuery = "";

  // Reactive declaration for filtered users
  $: filteredUsers = searchQuery ? searchUsers(searchQuery) : users;

  const searchUsers = (searchQuery) => {
    const [firstNameQuery, lastNameQuery] = searchQuery
      .toLowerCase()
      .trim()
      .split(" ");
    return users.filter((user) => {
      if (!firstNameQuery && !lastNameQuery) {
        return;
      }

      let firstNameMatch;
      let firstNameMatch2;
      let lastNameMatch;
      let lastNameMatch2;

      if (!lastNameQuery || (firstNameQuery && lastNameQuery)) {
        firstNameMatch = user.FirstName.toLowerCase().includes(firstNameQuery);
        firstNameMatch2 = user.LastName.toLowerCase().includes(firstNameQuery);
      }
      if (!firstNameQuery || (firstNameQuery && lastNameQuery)) {
        lastNameMatch = user.LastName.toLowerCase().includes(lastNameQuery);
        lastNameMatch2 = user.FirstName.toLowerCase().includes(lastNameQuery);
      }

      return (
        firstNameMatch || lastNameMatch || firstNameMatch2 || lastNameMatch2
      );
    });
  };

  export const selectUser = async (userID) => {
    const response = await fetch("http://localhost:8080/user?id=" + userID);
    if (response.ok) {
      const selectedUser = await response.json();
      userProfileData.set(selectedUser);
      searchQuery = "";
    } else {
      console.error("Error fetching users:", response.status);
    }
  };
</script>

<div class="searchUsers">
  <input type="search" bind:value={searchQuery} placeholder="Search users..." />
  {#if searchQuery}
    <div class="dropdown">
      {#each filteredUsers as user}
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div class="singleUser" on:click={() => selectUser(user.ID)}>
          <!-- svelte-ignore a11y-missing-attribute -->
          <img src={user.Avatar} />
          {user.FirstName}
          {user.LastName}
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  input {
    margin: 0;
    border-color: greenyellow;
    width: 98%;
    height: 100%;
  }

  .searchUsers {
    position: relative;
    margin-bottom: 8px;
  }

  .dropdown {
    position: absolute;
    width: 100%;
    z-index: 1;
    margin-top: 8px;
  }

  .singleUser {
    border: solid 1px #333;
    border-radius: 6px;
    cursor: pointer;
    padding: 8px;
    background: #011;
    margin: 2px;
  }

  img {
    max-height: 20px;
  }
</style>
