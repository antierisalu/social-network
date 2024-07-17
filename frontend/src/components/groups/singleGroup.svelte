<script>
  import Button from "../../shared/button.svelte";
  import Button2 from "../../shared/button2.svelte";
  import GroupPostOverlay from "../posts/createPost.svelte"; // NEeds prop for group post
  import EventOverlay from "./createEvent.svelte";
  import SearchBar from "../profile/searchBar.svelte";
  import Posts from "../posts/posts.svelte";
  import Event from "./event.svelte";
  import {
    leaveGroup,
    joinGroup,
    selectUser,
    getPosts,
    getEvents,
    getGroups,
    deleteGroup,
  } from "../../utils";
  import {
    groupSelected,
    events,
    userInfo,
    API_URL,
    IMAGE_URL,
  } from "../../stores";
  import { fade } from "svelte/transition";
  let certainty = 50;
  getPosts();
  let group;
  let showPostOverlay;
  let showEventOverlay;
  $: getgroup($groupSelected);

  function getgroup(groupID) {
    fetch(`${API_URL}/getGroup`, {
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
        group = data;
        console.log(group);
      })
      .catch((error) => console.error(error));
  }

  export function togglePostOverlay() {
    showPostOverlay = !showPostOverlay;
    if (!showPostOverlay) {
      getGroupPosts();
    }
  }

  export function toggleEventOverlay() {
    showEventOverlay = !showEventOverlay;
    if (!showEventOverlay) {
      getEvents(group.id);
    }
  }
</script>

<main in:fade>
  {#if group && group.joinStatus === 1}
    {#if showPostOverlay}
      <GroupPostOverlay on:close={togglePostOverlay} />
    {/if}
    {#if showEventOverlay}
      <EventOverlay on:close={toggleEventOverlay} />
    {/if}

    <div class="group">
      <div class="topPart">
        <div class="leftSide">
          <img src="{IMAGE_URL}{group.media}" alt="" />
          <div class="groupTitle">{group.title}</div>
          <!-- svelte-ignore a11y-click-events-have-key-events -->
          <div class="owner" on:click={selectUser(group.ownerID)}>
            Created by: {group.ownerName}
          </div>
          <div class="groupDescription">{group.description}</div>
        </div>
        <div class="groupImage">
          {#if group.media.Valid === true}
            <img src="{IMAGE_URL}{group.media.String}" alt="" />
          {/if}
        </div>
        <div class="rightSide">
          {#if group.ownerID == $userInfo.id}
            <div>
              <Button2
                btnText="Delete Group"
                onClick={() => deleteGroup(group.id)}
              ></Button2>
            </div>
          {:else}
            <div>
              <Button2
                btnText="Leave Group"
                onClick={() => leaveGroup(group.id)}
              ></Button2>
            </div>
          {/if}
          <SearchBar
            groupID={group.id}
            isGroup={true}
            placeHolda="Invite Users"
            w120
          />
        </div>
      </div>
      <div class="events">
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div class="createEvent" on:click={toggleEventOverlay}>
          Add new event..
        </div>
        {#if $events}
          {#each $events as event}
            <Event {event} />
          {/each}
        {/if}
      </div>
      <div class="posts">
        <Posts posts={group.posts}></Posts>
      </div>
    </div>
  {:else if group && group.joinStatus !== 1}
    <div class="topPart">
      <div class="leftSide">
        <div class="groupTitle">{group.title}</div>
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div class="owner" on:click={selectUser(group.ownerID)}>
          Created by: {group.ownerName}
        </div>
        <div class="groupDescription">{group.description}</div>
      </div>
      <div class="groupImage">
        {#if group.image}
          <img src="{IMAGE_URL}{group.image}" alt="" />
        {/if}
      </div>
      <div class="">
        {#if group.joinStatus == -1}
          <Button type="secondary" on:click={() => joinGroup(group.id, 0)}
            >Join Group</Button
          >
        {:else if group.joinStatus == 0}
          <Button2 btnText="Cancel Request" onClick={() => leaveGroup(group.id)}
          ></Button2>
        {/if}
      </div>
    </div>
  {/if}
</main>

<style>
  div {
    padding: 4px;
  }

  .groupTitle {
    font-size: xx-large;
  }

  .owner {
    font-size: small;
  }

  .owner:hover {
    font-weight: bold;
    color: yellowgreen;
    cursor: pointer;
  }

  .events {
    background-color: rgba(103, 158, 20, 0.27);
  }

  .events,
  .topPart,
  .createEvent,
  .groupDescription {
    border: solid 1px #555;
    border-radius: 8px;
    margin: 4px 0;
  }

  .topPart {
    display: flex;
    justify-content: space-between;
  }

  .leftSide {
    max-width: 280px;
    min-width: 250px;
    border: none;
  }
  .rightSide {
    max-width: 220px;
    min-width: 130px;
    border: none;
  }

  .groupImage {
    display: flex;
    flex-grow: 1;
    justify-content: center;
    flex-shrink: 1;
    max-width: 600px;
    max-height: 600px;
    border: none;
  }
  .groupImage img {
    display: flex;
    flex-shrink: 1;
    flex-grow: 1;
    max-width: 100%;
    max-height: 100%;
  }

  .createEvent {
    background-color: #011;
    display: flex;
    flex-direction: row;
    color: #555;
    border-color: greenyellow;
    padding: 8px;
    margin: 4px;
  }

  .createEvent:hover {
    cursor: pointer;
  }
</style>
