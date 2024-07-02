<script>
    import Button from "../../shared/button.svelte";

    let groups = [
        { name: 'Time travel', private: true, Member: true, requestedToJoin: false },
        { name: 'Dark Magic', private: false, Member: false, requestedToJoin: false },
        { name: 'Hacking', private: false, Member: false, requestedToJoin: true },
        { name: 'Resurrection', private: false, Member: false, requestedToJoin: false }
    ]


    function handleClick() {

    }

</script>

<main>
    <Button type="secondary" customStyle="width:100%">Create group</Button>
    <div class="groups">
        {#each groups as group }
        <div class="singleGroup">
            <div class="groupName">{group.name}</div>
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