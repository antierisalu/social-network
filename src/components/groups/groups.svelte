<script>
    import Button from "../../shared/button.svelte";
    import CreateGroup from "./createGroup.svelte";
    import { allGroups } from "../../stores";
    import { getGroups } from "../../utils";
    import { onMount } from "svelte";


    let groups = [
        { title: 'Dark Magic', private: false, Member: false, requestedToJoin: false },
        { title: 'Time travel', private: true, Member: true, requestedToJoin: false },
        { title: 'Hacking', private: false, Member: false, requestedToJoin: true },
        { title: 'Resurrection', private: false, Member: false, requestedToJoin: false }
    ]


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


</script>

<main>

    {#if showOverlay}
    <CreateGroup on:close={toggleOverlay} />
    {/if}

    <Button type="secondary" customStyle="width:98%" on:click={toggleOverlay}>Create group</Button>
    <div class="groups">
        {#each $allGroups as group }
        <div class="singleGroup">
            <div class="groupName">{group.title}</div>
            {#if !group.Member && !group.requestedToJoin}
                <Button type="secondary" customStyle="margin-bottom: 0; max-height:35px">Join</Button>
            {:else if !group.Member && group.requestedToJoin}
                <Button inverse={true} customStyle="margin-bottom: 0">Cancel Request</Button>
            {:else}
                <Button inverse={true} customStyle="margin-bottom: 0">Leave</Button>
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


    }


</style>