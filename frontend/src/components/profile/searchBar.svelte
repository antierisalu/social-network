<script>
  import {
    allUsers,
    userProfileData,
    API_URL,
    IMAGE_URL,
    userInfo,
  } from "../../stores";
  import { joinGroup } from "../../utils";
  import { sendMessage } from "../../websocket";
  import { fade, fly } from "svelte/transition";

  export let placeHolda = "Search users";
  export let w120 = false;
  export let isGroup = false;
  export let groupID;
  let groupMembers = [];

  if (isGroup) {
    placeHolda = "Search groups";
    fetch(`${API_URL}/getGroupMembers`, {
      credentials: "include",
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        groupID: groupID,
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        groupMembers = data;
      })
      .catch((error) => console.error(error));
  }

  $: users = isGroup
    ? $allUsers.filter(
        (user) => !groupMembers.map((m) => m.ID).includes(user.ID),
      )
    : $allUsers;

  var searchQuery = "";

  // Reactive declaration for filtered users
  $: filteredUsers = searchQuery ? searchUsers(searchQuery) : users;

  const searchUsers = (searchQuery) => {
    console.log(users);
    // console.log(isGroup, searchQuery)
    if (isGroup && searchQuery === " ") {
      return users;
    }
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

  export async function inviteUser(userID, groupID, event) {
    fetch(`${API_URL}/joinGroup`, {
      credentials: "include",
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        groupID: groupID,
        action: 2,
        targetID: userID,
      }),
    })
      .then((response) => {
        if (response.ok) {
          response.json().then((data) => {
            sendMessage(
              JSON.stringify({
                type: "groupInvite",
                fromid: $userInfo.id,
                groupID: groupID,
                targetID: userID,
                data: `groupInvite_${$userInfo.id}_${groupID}`,
              }),
            ); //link == groupInvite_fromid_groupid
              
            filteredUsers = filteredUsers.filter(user => user.ID !== userID)
          });
        }
      })
      .catch((error) => {
        console.error("Error inviting user to group:", error);
      });
  }

  export const selectUser = async (userID) => {
    const response = await fetch(`${API_URL}/user?id=${userID}`, {
      credentials: "include",
    });
    if (response.ok) {
      const selectedUser = await response.json();
      console.log(selectedUser);
      userProfileData.set(selectedUser);
      searchQuery = "";
    } else {
      console.error("Error fetching users:", response.status);
    }
  };
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
  class="searchUsers"
  on:click={() => {
    if (isGroup) searchQuery = " ";
  }}
>
  <input
    type="search"
    bind:value={searchQuery}
    placeholder={placeHolda}
    class:w120
  />
  {#if searchQuery}
    <div class="dropdown">
      {#each filteredUsers as user (user.ID)}
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div
          out:fly={{ y: -100, duration: 200 }}
          class="singleUser"
          on:click={(e) => {
            isGroup ? inviteUser(user.ID, groupID) : selectUser(user.ID);
          }}
        >
          <!-- svelte-ignore a11y-missing-attribute -->
          <img src="{IMAGE_URL}{user.Avatar}" />
          {user.FirstName}
          {user.LastName}
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .w120 {
    width: 120px;
  }

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
