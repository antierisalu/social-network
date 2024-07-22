<script>
  import Button from "../../shared/button.svelte";
  import SlowButton from "../../shared/button2.svelte";
  import CreateGroup from "./createGroup.svelte";
  import { allGroups, groupSelected, userInfo } from "../../stores";
  import { getGroups, leaveGroup, joinGroup, getEvents } from "../../utils";
  import { onMount } from "svelte";

  onMount(async () => {
    await getGroups();
  });

  let showOverlay = false;

  export function toggleOverlay() {
    showOverlay = !showOverlay;
    if (!showOverlay) {
      getGroups();
    }
  }

  const openGroup = (groupID) => {
    $groupSelected = 0;
    $groupSelected = groupID;
    getEvents($groupSelected);
  };
</script>

<main>
  {#if showOverlay}
    <CreateGroup on:close={toggleOverlay} />
  {/if}

  <Button type="secondary" customStyle="width:98%" on:click={toggleOverlay}
    >Create group</Button
  >
  {#if $allGroups}
    <div class="groups">
      {#each $allGroups as group}
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div class="singleGroup" on:click={openGroup(group.id)}>
          <div class="groupName">
            {group.title}
          </div>
          {#if group.ownerID != $userInfo.id}
            {#if group.joinStatus == -1}
              <Button
                type="secondary"
                customStyle="margin-bottom: 0; max-height:35px"
                on:click={joinGroup(group.id, 0)}>Join</Button
              >
            {:else if group.joinStatus == 2}
              <Button
                type="secondary"
                customStyle="margin-bottom: 0; max-height:55px"
                on:click={joinGroup(group.id, 1)}>Accept Request</Button
              >
              <!-- 0 == request to join -->
            {:else if group.joinStatus == 0}
              <SlowButton
                btnText="Cancel Request"
                onClick={() => leaveGroup(group.id)}
              ></SlowButton>
            {:else if group.joinStatus == 1}
              <SlowButton
                btnText="Leave Group"
                onClick={() => leaveGroup(group.id)}
              ></SlowButton>
            {/if}
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</main>

<style>
  .singleGroup {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    align-items: center;
    border: solid 1px #333;
    border-radius: 8px;
    padding: 8px;
    text-align: center;
    margin-bottom: 0;
    margin: 4px;
    min-height: 70px;
  }
  .groupName {
    word-break: break-all;
    padding: 8px;
    cursor: pointer;
  }
</style>
