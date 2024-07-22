<script>
    import { createEventDispatcher } from "svelte";
    export let notification;
    const dispatch = createEventDispatcher();
</script>

<div class="action-buttons">
    {#if notification.type === "followRequest" || notification.type === "groupRequest" || notification.type === "groupInvite"}
        <button
            on:click={() => {
                dispatch("updateremove", {
                    action: 1,
                    type: notification.type,
                });
            }}
        >
            Accept
        </button>
        <button
            on:click={() => {
                dispatch("updateremove", {
                    action: -1,
                    type: notification.type,
                });
            }}
        >
            Decline
        </button>
    {/if}
</div>

<style>
    button {
        border: none;
        border-radius: 4px;
        cursor: pointer;
        background-color: green;
        color: white;
    }

    button:hover {
        background-color: #ddd;
        color: black;
    }

    .action-buttons {
        display: flex;
        justify-content: center;
        gap: 10px;
    }

    .action-buttons button {
        margin: 0;
    }
</style>
