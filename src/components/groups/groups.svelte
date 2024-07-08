<script>
  import Button from "../../shared/button.svelte";
  import CreateGroup from "./createGroup.svelte";
  import { allGroups, groupSelected } from "../../stores";
  import { getGroups, leaveGroup, joinGroup } from "../../utils";
  import { onMount } from "svelte";

  onMount(async () => {
    await getGroups();
    console.log($allGroups);
  });

  let showOverlay = false;

  export function toggleOverlay() {
    showOverlay = !showOverlay;
    if (!showOverlay) {
      getGroups();
    }
  }

  const openGroup = (groupID) => {
    $groupSelected = groupID;
  };
</script>

<main>
  {#if showOverlay}
    <CreateGroup on:close={toggleOverlay} />
  {/if}

  <Button type="secondary" customStyle="width:98%" on:click={toggleOverlay}
    >Create group</Button
  >
  <div class="groups">
    {#each $allGroups as group}
      <div class="singleGroup">
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div class="groupName" on:click={openGroup(group.id)}>
          {group.title}
        </div>
        {#if group.joinStatus === -1}
          <Button
            type="secondary"
            customStyle="margin-bottom: 0; max-height:35px"
            on:click={joinGroup(group.id, 1)}>Join</Button
          >
        {:else if group.joinStatus === 0}
          <Button
            inverse={true}
            customStyle="margin-bottom: 0"
            on:click={leaveGroup(group.id)}>Cancel Request</Button
          >
        {:else if group.joinStatus === 1}
          <Button
            inverse={true}
            customStyle="margin-bottom: 0"
            on:click={leaveGroup(group.id)}>Leave</Button
          >
        {/if}
      </div>
    {/each}
  </div>
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
  }
  .groupName {
    word-break: break-all;
    padding: 8px;
    cursor: pointer;
  }
</style>
