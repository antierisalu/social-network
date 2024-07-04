<script>
import Button from "../../shared/button.svelte";
import GroupPostOverlay from "../posts/createPost.svelte"; // NEeds prop for group post
import EventOverlay from "./createEvent.svelte";
import SearchBar from "../profile/searchBar.svelte"

let showPostOverlay
let showEventOverlay

let group = {
    owner: "LAENUHAI",
    title: "Kinnisvara müük Viljandis",
    description: "This is a group for people who would like to breathe"
}

let events = [
    {
        creator: "Teresa",
        title: "Kesksuve koristus Pikal tänaval",
        description: "Kõigepealt puhastame jõe vee ära ja siis vaatame edasi. Palun registreerida",
        date: "22.08.24",
        RSVP: "Not Going"
    },
    {   
        creator: "Reese Withoutherspoon",
        title: "Üle Viljandi järve jooks (Jeesuse või Kalevipoja mod lubatud)",
        description: "Võistlusel osaleda ei saa kained! Äärmisel juhul võid kasutada aineid. Start kui viina enam poest ei saa, ehk siis 22.00. Pealtvaatajad võivad olla kained",
        date: "43.07.245",
        RSVP: "Going"
    },
]

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
  const getGroupPosts = () => console.log("i want that new post which i created in the group")
  const getEvents = () => console.log("i want that new event which i just created")

</script>

<main>
    {#if showPostOverlay}
    <GroupPostOverlay on:close={togglePostOverlay} />
    {/if}
    {#if showEventOverlay}
    <EventOverlay on:close={toggleEventOverlay} />
    {/if}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="createPost" on:click={togglePostOverlay}>Create new post to the group..</div>

    <div class="group">
        <div class="topPart">
            <div class="leftSide">
                <div class="groupTitle">{group.title}</div>
                <div class="owner">Created by: {group.owner}</div>
                <div class="groupDescription">{group.description}</div>
            </div>
            <div class="groupImage"><img src="../postsImages/2"></div>
            <div class="rightSide">
                <Button w120 inverse>Leave Group</Button>
                <SearchBar isGroup={true} placeHolda="Invite Users" w120/>
            </div>
        </div>
        <div class="events">
            <div class="createPost" on:click={toggleEventOverlay}>Add new event..</div>
            {#each events as event}
            <div class="singleEvent">
                <div class="eventInfo">
                    <div class="eventTitle">{event.title}</div>
                    <div class="eventDescription">{event.description}</div>
                </div>
                <div class="eventDate">
                    <div>{event.date}</div>
                    <Button type="secondary" inverse w120>{event.RSVP}</Button>
                </div>
            </div>
            {/each}
        </div>
    </div>

</main>

<style>
    div {
        padding:4px;
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

    .events, .singleEvent, .topPart, .createPost, .groupDescription   {
        border: solid 1px #555;
        border-radius: 8px;
        margin: 4px 0
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

  .createPost {
    display: flex;
    flex-direction: row;
    color: #555;
    border-color: greenyellow;
    padding: 8px;
    margin: 4px;
  }
  .createPost:hover {
    cursor: pointer;
  }
    

</style>
  