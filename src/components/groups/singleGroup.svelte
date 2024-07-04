<script>
import Button from "../../shared/button.svelte";
import PostOverlay from "../posts/createPost.svelte";

let showOverlay

let group = {
    owner: "0wn3r",
    title: "Air Group",
    description: "This is a group for people who likes to breathe"
}

let events = [
    {
        title: "Lets breathe",
        RSVP: "Not Going"
    },
    {
        title: "Lets not breathe",
        RSVP: "Going"
    },
]

export function toggleOverlay() {
    showOverlay = !showOverlay;
    if (!showOverlay) {
      getPosts();
    }
  }
</script>

<main>
    {#if showOverlay}
    <PostOverlay on:close={toggleOverlay} />
    {/if}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="createPost" on:click={toggleOverlay}>Create new post to the group..</div>


    <div class="group">
        <div class="topPart">
            <div class="leftSide">
                <div class="groupTitle">{group.title}</div>
                <div class="groupDescription">{group.description}</div>
                <div class="owner">{group.owner}</div>
            </div>
            <div class="groupImage"><img src="../avatars/default.png"></div>
            <div class="rightSide">
                <Button>Leave</Button>
                <Button type="secondary">Invite Users</Button>
            </div>
        </div>
        <div class="events">
            <Button type="secondary">Create Event</Button>
            {#each events as event}
            <div class="singleEvent">
                <div class="eventTitle">{event.title}</div>
                <Button type="secondary" inverse>{event.RSVP}</Button>
            </div>
            {/each}
        </div>
    </div>

</main>

<style>
    .singleEvent {
        display: flex;
        justify-content: space-between;
    }

    .topPart {
        display: flex;
        justify-content: space-between;
    }
    .rightSide, .leftSide {
        width: 150px
    }

    .groupImage img{
        max-width: 100%;
    }

  .createPost {
    display: flex;
    flex-direction: row;
    color: #555;
    border: solid 1px greenyellow;
    border-radius: 8px;
    padding: 8px;
    margin: 4px;
  }
  .createPost:hover {
    cursor: pointer;
  }
    

</style>
  