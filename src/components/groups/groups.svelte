<script>
    import Button from "../../shared/button.svelte";
    import CreateGroup from "./createGroup.svelte";
    import { allGroups } from "../../stores";
    import { getGroups, leaveGroup } from "../../utils";
    import { onMount } from "svelte";

    onMount(async () =>{
        await getGroups();
        console.log($allGroups)
    })
    
    let showOverlay = false
    
    export function toggleOverlay() {
        showOverlay = !showOverlay;
        if (!showOverlay) {
        getGroups();
    }
  }

    const joinGroup = (groupID) => console.log("attempt to join group with groupID:", groupID)
    const openGroup = (groupID) => console.log("attempt to open group with groupID:", groupID)

</script>

<main>

    {#if showOverlay}
    <CreateGroup on:close={toggleOverlay} />
    {/if}

    <Button type="secondary" customStyle="width:98%" on:click={toggleOverlay}>Create group</Button>
    <div class="groups">
        {#each $allGroups as group }
        <div class="singleGroup">
            <div class="groupName" on:click={openGroup(group.id)}>{group.title}</div>
            {#if !group.Member && !group.requestedToJoin}
                <Button type="secondary" customStyle="margin-bottom: 0; max-height:35px" on:click={joinGroup(group.id)}>Join</Button>
            {:else if !group.Member && group.requestedToJoin}
                <Button inverse={true} customStyle="margin-bottom: 0" on:click={leaveGroup(group.id)}>Cancel Request</Button>
            {:else}
                <Button inverse={true} customStyle="margin-bottom: 0" on:click={leaveGroup(group.id)}>Leave</Button>
            {/if}
        </div>
        {/each}
    </div>
    
</main>

<style>
    .singleGroup{
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        align-items: center;
        border: solid 1px #333;
        border-radius:8px;
        padding: 8px;
        text-align: center;
        margin-bottom: 0;
        margin: 4px
    }
    .groupName {
        word-break: break-all;
        padding: 8px;
        cursor: pointer;


    }


</style>