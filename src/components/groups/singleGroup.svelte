<script>
  import Button from "../../shared/button.svelte";
  import GroupPostOverlay from "../posts/createPost.svelte"; // NEeds prop for group post
  import EventOverlay from "./createEvent.svelte";
  import SearchBar from "../profile/searchBar.svelte";
  import Posts from "../posts/posts.svelte";
  import { leaveGroup, joinGroup, selectUser } from "../../utils";
  import { groupSelected } from "../../stores";
  import { fly, fade } from "svelte/transition";

  const getGroupPosts = () =>
    console.log("i want that new post which i created in the group");
  const getEvents = () =>
    console.log("i want that new event which i just created");

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

  let events = [
    {
      creator: "Teresa",
      title: "Kesksuve koristus Pikal tänaval",
      description:
        "Kõigepealt puhastame jõe vee ära ja siis vaatame edasi. Palun registreerida",
      date: "22.08.2024",
      RSVP: "Not Going",
    },
    {
      creator: "Reese Withoutherspoon",
      title: "Üle Viljandi järve jooks (Jeesuse või Kalevipoja mod lubatud)",
      description:
        "Võistlusel osaleda ei saa kained! Äärmisel juhul võid kasutada aineid. Start kui viina enam poest ei saa, ehk siis 22.00. Pealtvaatajad võivad olla kained",
      date: "43.27.245",
      RSVP: "Going",
    },
  ];

  export function togglePostOverlay() {
    showPostOverlay = !showPostOverlay;
    if (!showPostOverlay) {
      getGroupPosts();
    }
  }

  export function toggleEventOverlay() {
    showEventOverlay = !showEventOverlay;
    if (!showEventOverlay) {
      getEvents();
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
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- <div class="createGroupPost" on:click={togglePostOverlay}>
      Create new post to the group..
    </div> -->

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
        {#each events as event}
          <div class="singleEvent">
            <div class="eventInfo">
              <div class="eventTitle">{event.title}</div>
              <div class="eventDescription">{event.description}</div>
            </div>
            <div class="eventDate">
              <div>{event.date}</div>
              <Button type="secondary" inverse>{event.RSVP}</Button>
            </div>
          </div>
        {/each}
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
  .groupName {
    font-weight: bold;
    font-size: large;
    color: yellowgreen;
    cursor: pointer;
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
  .createGroupPost,
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

  .createGroupPost,
  .createEvent {
    display: flex;
    flex-direction: row;
    color: #555;
    border-color: greenyellow;
    padding: 8px;
    margin: 4px;
  }
  .createGroupPost:hover,
  .createEvent:hover {
    cursor: pointer;
  }
</style>
