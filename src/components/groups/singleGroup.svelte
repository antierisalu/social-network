<script>
  import Button from "../../shared/button.svelte";
  import GroupPostOverlay from "../posts/createPost.svelte"; // NEeds prop for group post
  import EventOverlay from "./createEvent.svelte";
  import SearchBar from "../profile/searchBar.svelte";
  import Posts from "../posts/posts.svelte";
  import {
    leaveGroup,
    joinGroup,
    selectUser,
    getPosts,
    getEvents,
  } from "../../utils";
  import { groupSelected, events } from "../../stores";
  import { fade } from "svelte/transition";
  let certainty = 50;
  getPosts();
  let group;
  let showPostOverlay;
  let showEventOverlay;
  $: getgroup($groupSelected);

  function getgroup(groupID) {
    console.log("attempt to get group with groupID:", groupID);
    fetch(`/getGroup`, {
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
        console.log(data);
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

  function sendRSVP(eventID, certainty) {
    console.log("attempt to send rsvp for eventID:", eventID);
    fetch(`/sendRSVP`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        RSVP: RSVP,
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        event.certainty = data;
        console.log(data);
      })
      .catch((error) => console.error(error));
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
          <div class="groupTitle">{group.title}</div>
          <!-- svelte-ignore a11y-click-events-have-key-events -->
          <div class="owner" on:click={selectUser(group.ownerID)}>
            Created by: {group.ownerName}
          </div>
          <div class="groupDescription">{group.description}</div>
        </div>
        <div class="groupImage">
          {#if group.image}
            <img src={group.image} alt="" />
          {/if}
        </div>
        <div class="rightSide">
          <Button inverse on:click={() => leaveGroup(group.id)}
            >Leave Group</Button
          >
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
            <script>
              event.certainty = 50;
            </script>
            <div class="singleEvent">
              <div class="eventInfo">
                <div class="eventTitle">{event.title}</div>
                <!-- <div class="owner" on:click={selectUser(event.ownerID)}>
                Created by: {event.ownerName}
              </div> -->
                <div class="eventDescription">{event.description}</div>
              </div>
              <div class="eventDate">
                <div>{event.date}</div>

                <div class="slideContainer">
                  <p>Slide for RSVP</p>

                  <input
                    bind:value={event.certainty}
                    type="range"
                    min="0"
                    max="100"
                    class="slider"
                    id="myRange"
                  />
                  {#if event.certainty < 20}
                    <p>Not going</p>
                  {:else if event.certainty > 80}
                    <p>Going</p>
                  {:else}
                    <p>Not Sure</p>
                  {/if}
                </div>

                <Button
                  type="secondary"
                  inverse
                  on:click={() => sendRSVP(event.id, event.certainty)}
                  >Submit</Button
                >
              </div>
            </div>
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
          <img src={group.image} alt="" />
        {/if}
      </div>
      <div class="">
        <Button type="secondary" on:click={() => joinGroup(group.id, 1)}
          >Join Group</Button
        >
      </div>
    </div>
  {/if}
</main>

<style>
  .slideContainer {
    font-weight: bold;
    font-size: large;
    accent-color: yellowgreen;
  }

  p,
  #myRange {
    margin: 0;
  }

  div {
    padding: 4px;
  }

  .singleEvent {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .eventInfo {
    width: 100%;
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
  .events,
  .singleEvent,
  .topPart,
  .createEvent,
  .groupDescription {
    border: solid 1px #555;
    border-radius: 8px;
    margin: 4px 0;
  }

  .eventDescription {
    font-size: x-small;
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
